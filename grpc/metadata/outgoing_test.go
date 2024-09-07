package metadata_test

import (
	"github.com/odycenter/std-library/grpc/metadata"
	"testing"
	"time"
)

func TestNewOutgoing(t *testing.T) {
	metadata.NewOutgoing().SetPairs("a", "1234").Ctx()
	metadata.NewOutgoing().SetMap(map[string]string{"b": "5678"}).Ctx()
	metadata.NewOutgoing().WithCancel().SetPairs("a", "1234").Ctx()
	metadata.NewOutgoing().WithCancel().SetMap(map[string]string{"b": "5678"}).Ctx()
	metadata.NewOutgoing().WithTimeout(30*time.Second).SetPairs("a", "1234").Ctx()
	metadata.NewOutgoing().WithTimeout(30 * time.Second).SetMap(map[string]string{"b": "5678"}).Ctx()
	metadata.NewOutgoing().WithDeadline(time.Now().Add(30*time.Second)).SetPairs("a", "1234").Ctx()
	metadata.NewOutgoing().WithDeadline(time.Now().Add(30 * time.Second)).SetMap(map[string]string{"b": "5678"}).Ctx()
}
