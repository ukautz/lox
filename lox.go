/*
This package provides a lock which can timeout and a convenience container
supporting named locks.
*/
package lox

import (
	"fmt"
	"time"
)

var (
	ErrorTimeout = fmt.Errorf("Timeout has been reached")
	ErrorNotLocked = fmt.Errorf("Unlock on not locked lock")
)

// Lock is a channel based locking which can timeout
type Lock struct {
	l      chan bool
	locked bool
}

// NewLock creates a new lock
func NewLock() *Lock {
	return &Lock{
		l:      make(chan bool, 1),
		locked: false,
	}
}

// Lock blocks until mutex is available or timeout is reached, in which case it
// returns the `ErrorTimeout` error.
// If timeout is zero than it behaves just like `sync.mutex.Lock()` and returns
// nil.
func (this *Lock) Lock(timeout time.Duration) error {
	if timeout > 0 {
		fail := time.After(timeout)
		select {
		case this.l <- true:
			this.locked = true
			return nil
		case <-fail:
			return ErrorTimeout
		}

	} else {
		this.l <- true
		this.locked = true
		return nil
	}
}

// Locked returns bool whether this Lock is currently locked or not.
func (this *Lock) Locked() bool {
	return this.locked
}

// Unlock either removes the lock and returns nil or returns the
// `ErrorNotLocked` error of no lock was there.
func (this *Lock) Unlock() error {
	if this.locked {
		this.locked = false
		<-this.l
		return nil
	} else {
		return ErrorNotLocked
	}
}

