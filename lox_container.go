package lox
import (
	"time"
	"sync"
	"fmt"
)

var (
	ErrorLockNotExisting = fmt.Errorf("Lock does not exist")
)

// Locks is a container for named locks capable of timing out
type Locks struct {
	lock  *sync.Mutex
	locks map[string]*Lock
	counts map[string]int
}

// NewLocks creates new container holding named locks
func NewLocks() *Locks {
	return &Locks{
		lock:  new(sync.Mutex),
		locks: make(map[string]*Lock),
		counts: make(map[string]int),
	}
}

// Run a function within lock or return error if lock cannot be achieved.
// See `Lock` method for the timeout parameter.
func (this *Locks) Run(name string, run func(), timeouts ...time.Duration) error {
	if err := this.Lock(name, timeouts...); err != nil {
		return err
	} else {
		defer this.Unlock(name)
		run()
		return nil
	}
}

// Lock tries to get a lock for a given name.
// The (first) timeout is the "try timeout" after which Lock returns with an error, if no lock could be gained.
func (this *Locks) Lock(name string, timeouts ...time.Duration) error {
	try := time.Duration(0)
	if l := len(timeouts); l > 0 {
		try = timeouts[0]
	}
	t := this.fetch(name)
	return t.Lock(try)
}

// Unlock removes lock on existing lock (which can result in `ErrorNotLocked`)
// or returns `ErrorLockNotExisting` error.
func (this *Locks) Unlock(name string) error {
	this.lock.Lock()
	defer this.lock.Unlock()
	if l := this.locks[name]; l != nil {
		return l.Unlock()
	} else {
		return ErrorLockNotExisting
	}
}

func (this *Locks) fetch(name string) *Lock {
	this.lock.Lock()
	defer this.lock.Unlock()
	if t := this.locks[name]; t == nil {
		this.locks[name] = NewLock()
	}
	return this.locks[name]
}