package s3_test

import (
	awss3 "std-library/sdk/aws/s3"
	"testing"
)

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
