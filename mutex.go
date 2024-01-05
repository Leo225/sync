package sync

import "errors"

var (
	ErrLockFailed   = errors.New("lock failed")
	ErrUnlockFailed = errors.New("unlock failed")
)

type Mutexer interface {
	Lock() (err error)
	Unlock() (err error)
}
