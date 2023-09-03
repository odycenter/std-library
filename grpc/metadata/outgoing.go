package metadata

import (
	"context"
	"google.golang.org/grpc/metadata"
	"time"
)

type OUT struct {
	ctx    context.Context
	cancel context.CancelFunc
	md     metadata.MD
}

// NewOutgoing 创建发送用的包含metadata的context
func NewOutgoing(ctx ...context.Context) *OUT {
	ctx = append(ctx, context.Background())
	return &OUT{ctx: ctx[0]}
}

// WithCancel 附加取消func
func (t *OUT) WithCancel() *OUT {
	ctx, cancel := context.WithCancel(t.ctx)
	t.ctx = ctx
	t.cancel = cancel
	return t
}

// WithTimeout 附加超时
func (t *OUT) WithTimeout(dur time.Duration) *OUT {
	ctx, cancel := context.WithTimeout(t.ctx, dur)
	t.ctx = ctx
	t.cancel = cancel
	return t
}

// WithDeadline 附加到期时间
func (t *OUT) WithDeadline(dt time.Time) *OUT {
	ctx, cancel := context.WithDeadline(t.ctx, dt)
	t.ctx = ctx
	t.cancel = cancel
	return t
}

// SetMap 根据给定的键值映射创建MD 。
// 键中只允许使用以下 ASCII 字符：
// 数字：0-9
// 大写字母：AZ（标准化为小写）
// 小写字母：az
// 特殊字符： -_.
// 大写字母会自动转换为小写字母。
// 以“grpc-”开头的密钥仅供 grpc 内部使用，如果在元数据中设置，可能会导致错误。
func (t *OUT) SetMap(m map[string]string) *OUT {
	t.md = metadata.New(m)
	return t
}

// SetPairs
// 组织一个由键、值映射形成的MD ...如果 len(kv) 为奇数，则 Pairs 会出现panic。
// 键中只允许使用以下 ASCII 字符：
// 数字：0-9
// 大写字母：AZ（标准化为小写）
// 小写字母：az
// 特殊字符： -_.
// 大写字母会自动转换为小写字母。
// 以“grpc-”开头的密钥仅供 grpc 内部使用，如果在元数据中设置，可能会导致错误。
func (t *OUT) SetPairs(kv ...string) *OUT {
	t.md = metadata.Pairs(kv...)
	return t
}

// Join 将任意数量的 md 加入到单个MD中。
// 每个键的值的顺序由包含这些值的 mds 呈现给 Join 的顺序确定。
func (t *OUT) Join(mds ...metadata.MD) *OUT {
	t.md = metadata.Join(mds...)
	return t
}

// Ctx 返回组装好的gRPC metadata context
func (t *OUT) Ctx() (context.Context, context.CancelFunc) {
	return metadata.NewOutgoingContext(t.ctx, t.md), t.cancel
}
