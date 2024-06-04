package async

import "context"

type Task interface {
	Execute(ctx context.Context)
}

type internalTask struct {
	process func(ctx context.Context)
}

func (i *internalTask) Execute(ctx context.Context) {
	i.process(ctx)
}
