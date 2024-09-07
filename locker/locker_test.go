package locker_test

import (
	"fmt"
	"github.com/odycenter/std-library/locker"
	"testing"
)

func TestLocker(t *testing.T) {
	opt := &locker.Option{
		Url:        []string{"127.0.0.1:6379"},
		UseCluster: false,
		Locker:     locker.Redis,
	}
	locker.Init(opt)
	if locker.Lock("locker_test_1") {
		fmt.Println("Locked")
	}
	if locker.Lock("locker_test_1") {
		fmt.Println("Locked")
	} else {
		fmt.Println("Failed")
	}
	locker.Unlock("locker_test_1")
	fmt.Println("Unlock")

	//opt = &locker.Option{
	//	Url:        []string{"127.0.0.1:2379"},
	//	UseCluster: false,
	//	Locker:     locker.Etcd,
	//}
	//locker.Init(opt)
	//if locker.Lock("/locker_test_1") {
	//	fmt.Println("Locked")
	//}
	//if locker.Lock("/locker_test_1") {
	//	fmt.Println("Locked")
	//} else {
	//	fmt.Println("Failed")
	//}
	//locker.Unlock("/locker_test_1")
	//fmt.Println("Unlock")
}
