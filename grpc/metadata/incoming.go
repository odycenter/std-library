package metadata

import (
	"context"
	"google.golang.org/grpc/metadata"
	"time"
)

// GetValues 从incoming 的 context 获取 metadata
func GetValues(ctx context.Context) (metadata.MD, bool) {
	return metadata.FromIncomingContext(ctx)
}

// Get 从incoming 的 context 获取 metadata 中的 key 值的values
func Get(ctx context.Context, key string) []string {
	return metadata.ValueFromIncomingContext(ctx, key)
}

type IN struct {
	ctx    context.Context
	cancel context.CancelFunc
}

// NewIncoming 创建接收用的包含metadata的context
func NewIncoming(ctx ...context.Context) *IN {
	ctx = append(ctx, context.Background())
	return &IN{ctx: ctx[0]}
}

// WithCancel 附加取消func
func (i *IN) WithCancel() *IN {
	ctx, cancel := context.WithCancel(i.ctx)
	i.ctx = ctx
	i.cancel = cancel
	return i
}

// WithTimeout 附加超时
func (i *IN) WithTimeout(dur time.Duration) *IN {
	ctx, cancel := context.WithTimeout(i.ctx, dur)
	i.ctx = ctx
	i.cancel = cancel
	return i
}

// WithDeadline 附加到期时间
func (i *IN) WithDeadline(dt time.Time) *IN {
	ctx, cancel := context.WithDeadline(i.ctx, dt)
	i.ctx = ctx
	i.cancel = cancel
	return i
}

// Ctx 返回组装好的gRPC metadata context
func (i *IN) Ctx(md metadata.MD) (context.Context, context.CancelFunc) {
	return metadata.NewIncomingContext(i.ctx, md), i.cancel
}
