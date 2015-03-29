package lox

import (
	. "github.com/smartystreets/goconvey/convey"
	"sync"
	"testing"
	"time"
)


func TestLocksWithoutTimeout(t *testing.T) {
	Convey("With a new locks container using no timeouts", t, func() {
		c := NewLocks()
		err := c.Lock("foo")
		Convey("No error is returned", func() {
			So(err, ShouldBeNil)
			Convey("Unlocking returns no error", func() {
				So(c.Unlock("foo"), ShouldBeNil)
				Convey("Unlocking again returns error", func() {
					So(c.Unlock("foo"), ShouldEqual, ErrorNotLocked)
				})
			})
		})
		Convey("Unlocking a non existing lock returns error", func() {
			So(c.Unlock("bar"), ShouldEqual, ErrorLockNotExisting)
		})
		Convey("Locking the same named lock is mutex", func() {
			start := time.Now()
			errors := []error{}
			var end time.Time
			var wg sync.WaitGroup
			wg.Add(2)
			go func() {
				defer wg.Done()
				time.Sleep(time.Duration(100) * time.Millisecond)
				if err := c.Unlock("foo"); err != nil {
					errors = append(errors, err)
				}
			}()
			go func() {
				defer wg.Done()
				if err := c.Lock("foo"); err != nil {
					errors = append(errors, err)
				} else {
					defer c.Unlock("name")
					end = time.Now()
				}
			}()

			wg.Wait()
			So(len(errors), ShouldEqual, 0)
			diff := end.Sub(start).Seconds() * 1000
			So(diff, ShouldBeBetweenOrEqual, 99, 105)
		})
	})
}

func TestLocksWithTimeout(t *testing.T) {
	Convey("With a new locks container using timeouts", t, func() {
		c := NewLocks()
		err := c.Lock("foo", time.Duration(100) * time.Millisecond)
		Convey("No error is returned", func() {
			So(err, ShouldBeNil)
			Convey("Unlocking returns no error", func() {
				So(c.Unlock("foo"), ShouldBeNil)
				Convey("Unlocking again returns error", func() {
					So(c.Unlock("foo"), ShouldEqual, ErrorNotLocked)
				})
			})
		})
		Convey("Locking the same named lock returns error after timeout", func() {
			start := time.Now()
			err := c.Lock("foo", time.Duration(100) * time.Millisecond)
			end := time.Now()

			So(err, ShouldEqual, ErrorTimeout)
			diff := end.Sub(start).Seconds() * 1000
			So(diff, ShouldBeBetweenOrEqual, 99, 105)
		})
	})
}

func TestLocksRun(t *testing.T) {
	Convey("With a new locks container running in lock", t, func() {
		c := NewLocks()
		Convey("Running on not existing runs and returns no error", func() {
			i := 0
			err := c.Run("foo", func() {
				i = 1
			})
			So(err, ShouldBeNil)
			So(i, ShouldEqual, 1)
		})
		Convey("Running on existing lock without timeout waits after lock is free", func() {
			c.Lock("foo")
			i := 0
			start := time.Now()
			var end time.Time
			go func() {
				defer c.Unlock("foo")
				time.Sleep(time.Duration(100) * time.Millisecond)
			}()
			err := c.Run("foo", func() {
				i = 1
				end = time.Now()
			})

			So(err, ShouldBeNil)
			So(i, ShouldEqual, 1)
			diff := end.Sub(start).Seconds() * 1000
			So(diff, ShouldBeBetweenOrEqual, 99, 105)
		})
		Convey("Running on existing lock with insufficient timeout does not run and returns error", func() {
			c.Lock("foo")
			i := 0
			go func() {
				defer c.Unlock("foo")
				time.Sleep(time.Duration(100) * time.Millisecond)
			}()
			err := c.Run("foo", func() {
				i = 1
			}, time.Duration(10) * time.Millisecond)

			So(err, ShouldEqual, ErrorTimeout)
			So(i, ShouldEqual, 0)
		})
		Convey("Running on existing lock with sufficient timeout runs and returns no error", func() {
			c.Lock("foo")
			i := 0
			start := time.Now()
			var end time.Time
			go func() {
				defer c.Unlock("foo")
				time.Sleep(time.Duration(20) * time.Millisecond)
			}()
			err := c.Run("foo", func() {
				i = 1
				end = time.Now()
			}, time.Duration(50) * time.Millisecond)

			So(err, ShouldBeNil)
			So(i, ShouldEqual, 1)
			diff := end.Sub(start).Seconds() * 1000
			So(diff, ShouldBeBetweenOrEqual, 19, 25)
		})
	})
}


