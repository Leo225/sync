package sync

import "time"

type Options struct {
	KeyPrefix   string
	LockTimeout time.Duration
	WaitRetry   time.Duration
}

type OptionFunc func(o *Options)

func newOption(opts ...OptionFunc) *Options {
	opt := &Options{
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
	return func(o *Options) {
		o.KeyPrefix = prefix
	}
}

func LockTimeout(duration time.Duration) OptionFunc {
	return func(o *Options) {
		o.LockTimeout = duration
	}
}

func WaitRetry(duration time.Duration) OptionFunc {
	return func(o *Options) {
		o.WaitRetry = duration
	}
}
