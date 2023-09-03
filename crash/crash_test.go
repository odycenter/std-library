package crash_test

import (
	"fmt"
	"testing"

	"github.com/odycenter/std-library/crash"
)

func TestTry(t *testing.T) {
	crash.Try(func() {
		panic("have a panic!")
		//panic(fmt.Errorf("have a panic!"))
	}).Catch(nil, func(err error) {
		fmt.Println(err)
	}).Finally(func() {
		fmt.Println("finally 1")
	}, func() {
		fmt.Println("finally 2")
	})
	fmt.Println("still exec")
}
