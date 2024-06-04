package oss_test

import (
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	oss2 "std-library/sdk/aliyun/oss"
	"testing"
)

func TestNew(t *testing.T) {
	t.Run("TestAliOSS", func(t *testing.T) {
		op := &oss2.Opt{
			Endpoint:        "https://oss-cn-hangzhou.aliyuncs.com",
			Bucket:          "Test",
			AccessKeyID:     "GIUHh2ujluhql29knl12",
			AccessKeySecret: "adq29udjkp2",
		}
		opt := oss.Callback("sda")
		op.SetDeleteOption(opt)
		got, err := oss2.New(op)
		if err != nil {
			return
		}
		err = got.CreateBucket()
		if err != nil {
			return
		}
	})
}
