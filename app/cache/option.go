package cache

type Options struct {
	OnError func(err error)
}

type Option func(*Options)

func OnError(onError func(err error)) Option {
	return func(o *Options) {
		o.OnError = onError
	}
}
