package demo

import (
	"context"
	"github.com/odycenter/std-library/demo/test"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strconv"
)

type HelloController struct {
	test.TestServiceServer
	callCount int
}

func (s *HelloController) SayHello(ctx context.Context, in *test.HelloRequest) (*test.HelloReply, error) {
	s.callCount++
	if s.callCount <= 2 {
		return nil, status.Error(codes.Unavailable, "Service unavailable"+strconv.Itoa(s.callCount))
	}
	s.callCount = 0
	return &test.HelloReply{Message: "Hello " + in.Name}, nil
}
