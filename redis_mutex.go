package sync

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisSync struct {
	Opt   *Options
	Redis *redis.Client
}

type redisMutex struct {
	sessionID  string
	key        string
	opt        *Options
	redis      *redis.Client
	lockCtx    context.Context
	lockCancel context.CancelFunc
}

func NewRedisSync(client *redis.Client, opts ...OptionFunc) *RedisSync {
	return &RedisSync{
		Redis: client,
		Opt:   newOption(opts...),
	}
}

func (s *RedisSync) NewMutex(key string, opts ...OptionFunc) Mutexer {
	opt := *s.Opt
	for _, o := range opts {
		o(&opt)
	}

	rm := &redisMutex{
		sessionID: sessionID(),
		key:       key,
		opt:       &opt,
		redis:     s.Redis,
	}
	return rm
}

const (
	luaRelease = `if redis.call("get", KEYS[1]) == ARGV[1] then return redis.call("del", KEYS[1]) else return 0 end`
)

func (m *redisMutex) Lock() (err error) {
	var flag bool
	lockName := m.lockName()
	flag, err = m.redis.SetNX(context.Background(), lockName, m.sessionID, m.opt.LockTimeout).Result()
	if err != nil {
		return
	}
	if !flag {
		err = ErrLockFailed
		return
	}
	m.lockCtx, m.lockCancel = context.WithCancel(context.Background())
	go func() {
		t := time.NewTimer(m.opt.WaitRetry)
		defer t.Stop()

		for {
			select {
			case <-t.C:
				_ = m.redis.Expire(context.Background(), lockName, m.opt.WaitRetry)
				t.Reset(m.opt.WaitRetry)
			//TODO: delay failed
			//fmt.Println("trigger delay key")
			case <-m.lockCtx.Done():
				//fmt.Println("unlock trigger context done, end for loop")
				return

			}
		}
	}()
	return
}

func (m *redisMutex) Unlock() (err error) {
	var flag bool
	lockName := m.lockName()
	flag, err = m.redis.Eval(context.Background(), luaRelease, []string{lockName}, m.sessionID).Bool()
	if err != nil {
		return
	}
	if !flag {
		err = ErrUnlockFailed
		return
	}
	if m.lockCancel != nil {
		m.lockCancel()
	}
	return
}

func (m *redisMutex) lockName() string {
	return m.opt.KeyPrefix + m.key
}
