package grpc_test

import (
	"fmt"
	"testing"

	"github.com/odycenter/std-library/grpc"
)

func TestGRPC(t *testing.T) {
	err := grpc.Register("", "https://www.google.com", &grpc.DefaultOptions)
	if err != nil {
		fmt.Println(err)
		return
	}
	conn, err := grpc.Get("")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()
	//client := conn.Conn()
}
