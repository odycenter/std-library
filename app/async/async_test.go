package async_test

import (
	"context"
	"github.com/odycenter/std-library/app/async"
	actionlog "github.com/odycenter/std-library/app/log"
	"github.com/odycenter/std-library/app/web/errors"
	"github.com/stretchr/testify/assert"
	"log/slog"
	"sync"
	"testing"
)

func TestRunFunc(t *testing.T) {
	async.RunFunc(nil, func(ctx context.Context) {
		actionlog.Context(&ctx, "key", "value2")
		slog.InfoContext(ctx, "info message", "key", "value")
	})
	var a int
	var b int
	wg := sync.WaitGroup{}
	async.RunFunc(nil, func(ctx context.Context) {
		a = 3
		actionlog.Context(&ctx, "key", "value2")
	}, &wg)

	async.RunFuncWithName(nil, "root-action", func(ctx context.Context) {
		async.RunFuncWithName(&ctx, "child-action", func(ctx context.Context) {
			b = 5
			actionlog.Context(&ctx, "key", "value")
			actionlog.Context(&ctx, "key", "value2")
		}, &wg)
	})

	wg.Wait()
	assert.Equal(t, 8, a+b)
}

func TestRunFuncFailed(t *testing.T) {
	var a int
	var b int
	wg := sync.WaitGroup{}
	async.RunFunc(nil, func(ctx context.Context) {
		actionlog.Context(&ctx, "key", "value")
		errors.Internal("failed", "ERROR_CODE")
	}, &wg)

	async.RunFunc(nil, func(ctx context.Context) {
		b = 5
	}, &wg)

	wg.Wait()
	assert.Equal(t, 5, a+b)
}
