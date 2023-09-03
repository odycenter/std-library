// Package s3 亚马逊S3对象存储服务SDK封装
package s3

import (
	"context"
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"io"
	"os"
	"path/filepath"
	"strings"
)

var (
	ErrOpt = errors.New("s3 option nil")
)

type Opt struct {
	Region    string // 地区
	Bucket    string // 文件库
	AccessKey string // 账号
	SecretKey string // 密码
}

type S3 struct {
	bucket string
	opt    *Opt
	sess   *session.Session
}

func New(opt *Opt) (*S3, error) {
	if opt == nil {
		return nil, ErrOpt
	}
	sess, err := session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials(opt.AccessKey, opt.SecretKey, ""),
		Region:      aws.String(opt.Region),
	})
	if err != nil {
		return nil, err
	}
	return &S3{bucket: opt.Bucket, opt: opt, sess: sess}, nil
}

// CreateBucket 创建桶
func (s *S3) CreateBucket() error {
	prov := s3.New(s.sess)
	_, err := prov.CreateBucket(&s3.CreateBucketInput{Bucket: &s.bucket})
	if err != nil {
		return err
	}
	err = prov.WaitUntilBucketExists(&s3.HeadBucketInput{Bucket: aws.String(s.bucket)})
	return err
}

// DeleteBucket 删除桶
func (s *S3) DeleteBucket() error {
	prov := s3.New(s.sess)
	_, err := prov.DeleteBucket(&s3.DeleteBucketInput{Bucket: aws.String(s.bucket)})
	if err != nil {
		return err
	}
	err = prov.WaitUntilBucketNotExists(&s3.HeadBucketInput{Bucket: aws.String(s.bucket)})
	if err != nil {
		return err
	}
	return nil
}

// UploadFile
// - 上传指定文件
// - key 就对应不同的文件夹名字比如 a/b/c/d
// - filenamePath 具体的文件路径
func (s *S3) UploadFile(key string, filePath string, body io.Reader) error {
	uploader := s3manager.NewUploader(s.sess)
	ctxType := ContextType(key)
	uploadInput := &s3manager.UploadInput{
		Bucket:      aws.String(s.bucket),
		ContentType: aws.String(ctxType),
		Key:         aws.String(key),
	}
	if body == nil {
		fs, err := os.Open(filePath)
		if err != nil {
			return err
		}
		defer fs.Close()
		uploadInput.Body = fs
	} else {
		uploadInput.Body = body
	}
	_, err := uploader.Upload(uploadInput)
	return err
}

// UploadFileSteam
// - 传入指定文件通过this.GetFile()
// - key 就对应不同的文件夹名字比如 a/b/c/d
// - fileStream 打开的文件流
func (s *S3) UploadFileSteam(key string, fileStream *os.File) error {
	uploader := s3manager.NewUploader(s.sess)
	_, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: &s.bucket,
		Key:    &key,
		Body:   fileStream,
	})
	return err
}

// UploadDir
// - 传入指定文件通过this.GetFile()
// - key 就对应不同的文件夹名字比如 a/b/c/d
// - fileStream 打开的文件流
func (s *S3) UploadDir(key string, dir string) error {
	uploader := s3manager.NewUploader(s.sess)
	err := uploader.UploadWithIterator(context.TODO(), newDirectoryIterator(s.bucket, dir, key))
	return err
}

// DownFile
// - 下载指定文件
// - key 就对应不同的文件夹名字比如 a/b/c/d
// - filenamePath 具体的文件路径
func (s *S3) DownFile(key string, filenamePath string) error {
	file, err := os.Create(filenamePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = s3manager.NewDownloader(s.sess).Download(file, &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})
	return err
}

// GetObject
// - 下载指定文件
// - key 就对应不同的文件夹名字比如 a/b/c/d
// - filenamePath 具体的文件路径
func (s *S3) GetObject(key string, filenamePath string) (*s3.GetObjectOutput, error) {
	return s3.New(s.sess).GetObject(&s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})
}

// 迭代文件夹全部内容
func newDirectoryIterator(bucket, dir, key string) s3manager.BatchUploadIterator {
	var paths []string
	_ = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			paths = append(paths, path)
		}
		return nil
	})

	return &DirectoryIterator{
		filePaths:    paths,
		bucket:       bucket,
		keyStartPath: key,
	}
}

// DirectoryIterator represents an iterator of a specified directory
type DirectoryIterator struct {
	filePaths []string
	bucket    string
	next      struct {
		path string
		f    *os.File
	}
	err          error
	keyStartPath string
}

func (di *DirectoryIterator) UploadObject() s3manager.BatchUploadObject {
	f := di.next.f

	// 这里兼容编码
	key := di.next.path
	if di.keyStartPath != "" {
		key = strings.Replace(key, di.keyStartPath, "", 1)
	}

	contextType := ContextType(key)
	return s3manager.BatchUploadObject{
		Object: &s3manager.UploadInput{
			Bucket:      &di.bucket,
			Key:         &key,
			Body:        f,
			ContentType: &contextType,
		},
		After: func() error {
			return f.Close()
		},
	}
}

// Next returns whether next file exists
func (di *DirectoryIterator) Next() bool {
	if len(di.filePaths) == 0 {
		di.next.f = nil
		return false
	}

	f, err := os.Open(di.filePaths[0])
	di.err = err
	di.next.f = f
	di.next.path = di.filePaths[0]
	di.filePaths = di.filePaths[1:]

	return di.Err() == nil
}

// Err returns error of DirectoryIterator
func (di *DirectoryIterator) Err() error {
	return di.err
}
