package demo_test

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	_ "net/http/pprof"
	"std-library/demo/test"
	"std-library/grpc"
	"testing"
	"time"
)

func TestGrpcClientRetry(t *testing.T) {
	opt := &grpc.Option{
		MaxIdle:              8,
		MaxActive:            64,
		MaxConcurrentStreams: 64,
		RecycleDur:           600,
		Reuse:                true,
		Logger:               grpc.Logger{Open: true},
	}
	err := grpc.Register("local", "127.0.0.1:80", opt)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	conn, err := grpc.Get("local")
	c := test.NewTestServiceClient(conn.Conn())

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	r, err := c.SayHello(ctx, &test.HelloRequest{Name: "world"})
	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			switch st.Code() {
			case codes.Unavailable:
				log.Printf("Service unavailable, but retry was attempted")
			default:
				log.Printf("Unexpected error: %v", err)
			}
		}
		return
	}
	log.Printf("Greeting: %s", r.GetMessage())
	return
}
