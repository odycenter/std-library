package s3_test

import (
	"fmt"
	"io"
	"math/rand"
	"os"
	"runtime"
	awss3 "std-library/sdk/aws/s3"
	"testing"
	"time"
)

var r *rand.Rand

func init() {
	source := rand.NewSource(time.Now().UnixNano())
	r = rand.New(source)
}

func TestNew(t *testing.T) {
	t.Run("TestAWSS3", func(t *testing.T) {
		op := &awss3.Opt{
			Region:    "http://xxx.xxx.xx",
			Bucket:    "Test",
			AccessKey: "oijoeofjvpjoej^skjdflk",
			SecretKey: "89xzdup890hO*H()08U)y9",
		}
		got, err := awss3.New(op)
		if err != nil {
			return
		}
		err = got.CreateBucket()
		if err != nil {
			return
		}
	})
}

func TestUploadFileMemoryLeak(t *testing.T) {
	opt := &awss3.Opt{
		Region:    "us-west-2",
		Bucket:    "test-bucket",
		AccessKey: "test-access-key",
		SecretKey: "test-secret-key",
	}
	s3Client, err := awss3.New(opt)
	if err != nil {
		t.Fatalf("Failed to create S3 client: %v", err)
	}

	runtime.GC()
	printMemUsage("Initial memory usage")

	for i := 0; i < 100; i++ {
		testFile := createLargeTestFile(t, generateRandomFileName(), 10*1024*1024) // 10MB
		defer os.Remove(testFile)

		err := s3Client.UploadFile(fmt.Sprintf("test_key_%d", i), testFile, nil)
		if err != nil {
			t.Fatalf("Failed to upload file: %v", err)
		}

		if i%10 == 0 {
			runtime.GC()
			printMemUsage(fmt.Sprintf("After %d uploads", i+1))
		}
	}

	runtime.GC()
	printMemUsage("Final memory usage")
}

func generateRandomFileName() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	length := 10
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[r.Intn(len(charset))]
	}
	return string(b) + ".txt"
}

func createLargeTestFile(t *testing.T, fileName string, size int) string {
	file, err := os.Create(fileName)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	defer file.Close()

	if _, err := io.Copy(file, io.LimitReader(newRandomDataReader(), int64(size))); err != nil {
		t.Fatalf("Failed to write to test file: %v", err)
	}

	return fileName
}

func printMemUsage(msg string) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("%s:\n", msg)
	fmt.Printf("Alloc = %v MiB", bToMb(m.Alloc))
	fmt.Printf("\tTotalAlloc = %v MiB", bToMb(m.TotalAlloc))
	fmt.Printf("\tSys = %v MiB", bToMb(m.Sys))
	fmt.Printf("\tNumGC = %v\n", m.NumGC)
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}

type randomDataReader struct{}

func newRandomDataReader() io.Reader {
	return &randomDataReader{}
}

func (r *randomDataReader) Read(p []byte) (n int, err error) {
	for i := range p {
		p[i] = byte(i % 256)
	}
	return len(p), nil
}
