package async_test

import (
	"context"
	"github.com/stretchr/testify/assert"
	"std-library/app/async"
	actionlog "std-library/app/log"
	"std-library/app/log/consts/logKey"
	"std-library/app/web/errors"
	"sync"
	"testing"
)

func TestRunFunc(t *testing.T) {
	var a int
	var b int
	wg := sync.WaitGroup{}
	async.RunFunc(nil, func(ctx context.Context) {
		a = 3
		actionlog.Context(&ctx, "key", "value2")
	}, &wg)

	ctx := context.WithValue(context.Background(), logKey.Action, "root_action")
	ctx = context.WithValue(ctx, logKey.Id, "id-abcd")
	async.RunFuncWithName(&ctx, "child-action", func(ctx context.Context) {
		b = 5
		actionlog.Context(&ctx, "key", "value")
		actionlog.Context(&ctx, "key", "value2")
	}, &wg)

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
