package sync

import "time"

type Option struct {
	KeyPrefix   string
	LockTimeout time.Duration
	WaitRetry   time.Duration
}

type OptionFunc func(o *Option)

func newOption(opts ...OptionFunc) *Option {
	opt := &Option{
		KeyPrefix:   "synclock: ",
		LockTimeout: 20 * time.Second,
		WaitRetry:   6 * time.Second,
	}

	for _, o := range opts {
		o(opt)
	}
	return opt
}

func KeyPrefix(prefix string) OptionFunc {
	return func(o *Option) {
		o.KeyPrefix = prefix
	}
}

func LockTimeout(duration time.Duration) OptionFunc {
	return func(o *Option) {
		o.LockTimeout = duration
	}
}

func WaitRetry(duration time.Duration) OptionFunc {
	return func(o *Option) {
		o.WaitRetry = duration
	}
}
