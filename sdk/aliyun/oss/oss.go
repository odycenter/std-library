// Package oss 阿里云OSS对象存储服务SDK封装
package oss

import (
	"errors"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"os"
)

var (
	ErrOpt = errors.New("sdk option nil")
)

// Opt 配置
type Opt struct {
	Endpoint        string     // 填写Bucket对应的Endpoint，以华东1（杭州）为例，填写为https://sdk-cn-hangzhou.aliyuncs.com。其它Region请按实际情况填写。
	Bucket          string     // 填写存储空间名称
	AccessKeyID     string     // 账号
	AccessKeySecret string     // 密码
	Option          *ossOption // Option对象
}

type ossOption struct {
	CreateOption            []oss.Option
	DeleteOption            []oss.Option
	PutObjectFromFileOption []oss.Option
	PutObjectOption         []oss.Option
	DeleteObjectOption      []oss.Option
	ListObjectsOption       []oss.Option
	ListObjectsV2Option     []oss.Option
}

func (o *Opt) SetCreateOption(op ...oss.Option) {
	if o.Option == nil {
		o.Option = new(ossOption)
	}
	o.Option.CreateOption = op
}
func (o *Opt) SetDeleteOption(op ...oss.Option) {
	if o.Option == nil {
		o.Option = new(ossOption)
	}
	o.Option.DeleteOption = op
}
func (o *Opt) SetPutObjectFromFileOption(op ...oss.Option) {
	if o.Option == nil {
		o.Option = new(ossOption)
	}
	o.Option.PutObjectFromFileOption = op
}
func (o *Opt) SetPutObjectOption(op ...oss.Option) {
	if o.Option == nil {
		o.Option = new(ossOption)
	}
	o.Option.PutObjectOption = op
}
func (o *Opt) SetDeleteObjectOption(op ...oss.Option) {
	if o.Option == nil {
		o.Option = new(ossOption)
	}
	o.Option.DeleteObjectOption = op
}
func (o *Opt) SetListObjectsOption(op ...oss.Option) {
	if o.Option == nil {
		o.Option = new(ossOption)
	}
	o.Option.ListObjectsOption = op
}
func (o *Opt) SetListObjectsV2Option(op ...oss.Option) {
	if o.Option == nil {
		o.Option = new(ossOption)
	}
	o.Option.ListObjectsV2Option = op
}

func (o *Opt) GetCreateOption() []oss.Option {
	if o.Option == nil {
		return nil
	}
	return o.Option.CreateOption
}
func (o *Opt) GetDeleteOption() []oss.Option {
	if o.Option == nil {
		return nil
	}
	return o.Option.DeleteOption
}
func (o *Opt) GetPutObjectFromFileOption() []oss.Option {
	if o.Option == nil {
		return nil
	}
	return o.Option.PutObjectFromFileOption
}
func (o *Opt) GetPutObjectOption() []oss.Option {
	if o.Option == nil {
		return nil
	}
	return o.Option.PutObjectOption
}
func (o *Opt) GetDeleteObjectOption() []oss.Option {
	if o.Option == nil {
		return nil
	}
	return o.Option.DeleteObjectOption
}
func (o *Opt) GetListObjectsOption() []oss.Option {
	if o.Option == nil {
		return nil
	}
	return o.Option.ListObjectsOption
}
func (o *Opt) GetListObjectsV2Option() []oss.Option {
	if o.Option == nil {
		return nil
	}
	return o.Option.ListObjectsV2Option
}

// Oss 使用ali oss管理文件上传 2022-10-06
type Oss struct {
	bucket string
	opt    *Opt
	cli    *oss.Client
}

// New 创建aliOss客户端
func New(opt *Opt) (*Oss, error) {
	if opt == nil {
		return nil, ErrOpt
	}
	cli, err := oss.New(opt.Endpoint, opt.AccessKeyID, opt.AccessKeySecret)
	if err != nil {
		return nil, err
	}
	return &Oss{opt.Bucket, opt, cli}, nil
}

func (a *Oss) GetCli() *oss.Client {
	return a.cli
}

// CreateBucket 创建桶
func (a *Oss) CreateBucket() error {
	err := a.cli.CreateBucket(a.bucket, a.opt.GetCreateOption()...)
	if err != nil {
		return err
	}
	return nil
}

// DeleteBucket 删除桶
func (a *Oss) DeleteBucket() error {
	err := a.cli.DeleteBucket(a.bucket, a.opt.GetDeleteOption()...)
	if err != nil {
		return err
	}
	return nil
}

// UploadFile
// - 上传指定文件
// - key 就对应不同的文件夹名字比如 a/b/c/d
// - path 具体的文件路径
func (a *Oss) UploadFile(key, path string) error {
	bucket, err := a.cli.Bucket(a.bucket)
	if err != nil {
		return err
	}
	err = bucket.PutObjectFromFile(key, path, a.opt.GetPutObjectFromFileOption()...)
	if err != nil {
		return err
	}
	return nil
}

// UploadFileSteam
// - 传入指定文件通过this.GetFile()
// - key 就对应不同的文件夹名字比如 a/b/c/d
// - stream 打开的文件流
func (a *Oss) UploadFileSteam(key string, stream *os.File) error {
	bucket, err := a.cli.Bucket(a.bucket)
	if err != nil {
		return err
	}
	err = bucket.PutObject(key, stream, a.opt.GetPutObjectOption()...)
	if err != nil {
		return err
	}
	return nil
}

// DeleteObject 删除文件
func (a *Oss) DeleteObject(key string) error {
	bucket, err := a.cli.Bucket(a.bucket)
	if err != nil {
		return err
	}
	err = bucket.DeleteObject(key, a.opt.GetDeleteObjectOption()...)
	if err != nil {
		return err
	}
	return nil
}

// DeleteObjects 删除多文件
func (a *Oss) DeleteObjects(key ...string) (*oss.DeleteObjectsResult, error) {
	bucket, err := a.cli.Bucket(a.bucket)
	if err != nil {
		return nil, err
	}
	delResult, err := bucket.DeleteObjects(key, a.opt.GetDeleteObjectOption()...)
	if err != nil {
		return nil, err
	}
	return &delResult, nil
}

// ListObjects 文件列表
func (a *Oss) ListObjects() (*oss.ListObjectsResult, error) {
	bucket, err := a.cli.Bucket(a.bucket)
	if err != nil {
		return nil, err
	}
	listResult, err := bucket.ListObjects(a.opt.GetListObjectsOption()...)
	if err != nil {
		return nil, err
	}
	return &listResult, nil
}

// ListObjectsV2 文件列表V2
func (a *Oss) ListObjectsV2() (*oss.ListObjectsResultV2, error) {
	bucket, err := a.cli.Bucket(a.bucket)
	if err != nil {
		return nil, err
	}
	listResult, err := bucket.ListObjectsV2(a.opt.GetListObjectsV2Option()...)
	if err != nil {
		return nil, err
	}
	return &listResult, nil
}
