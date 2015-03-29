package lox

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
	"sync"
)

func TestLockWithoutTimeout(t *testing.T) {
	Convey("When locking a new lock without timeout", t, func() {
		l := NewLock()
		err := l.Lock(0)
		Convey("No error is returned", func() {
			So(err, ShouldBeNil)
			Convey("Locking state can be seen", func() {
				So(l.Locked(), ShouldEqual, true)
				Convey("Unlocking returns no error", func() {
					So(l.Unlock(), ShouldBeNil)
					Convey("Locking state is removed", func() {
						So(l.Locked(), ShouldEqual, false)
						Convey("Further unlocking returns error", func() {
							So(l.Unlock(), ShouldEqual, ErrorNotLocked)
						})
					})
				})
			})
		})
		Convey("Locking is mutex", func() {
			var wg sync.WaitGroup
			wg.Add(2)
			start := time.Now()
			var end time.Time
			go func() {
				defer wg.Done()
				time.Sleep(time.Duration(100) * time.Millisecond)
				l.Unlock()
			}()
			go func() {
				defer wg.Done()
				l.Lock(0)
				end = time.Now()
			}()

			wg.Wait()
			diff := end.Sub(start).Seconds() * 1000
			So(diff, ShouldBeBetweenOrEqual, 99, 105)
		})
	})
}

func TestLockWithTimeout(t *testing.T) {
	Convey("When locking a new lock with timeout", t, func() {
		l := NewLock()
		err := l.Lock(time.Duration(20) * time.Millisecond)
		Convey("No error is returned", func() {
			So(err, ShouldBeNil)
			Convey("Locking state can be seen", func() {
				So(l.Locked(), ShouldEqual, true)
				Convey("Unlocking returns no error", func() {
					So(l.Unlock(), ShouldBeNil)
					Convey("Locking state is removed", func() {
						So(l.Locked(), ShouldEqual, false)
						Convey("Further unlocking returns error", func() {
							So(l.Unlock(), ShouldEqual, ErrorNotLocked)
						})
					})
				})
			})
		})
		Convey("Locking the same errors after timeout", func() {
			So(l.Lock(time.Duration(10) * time.Millisecond), ShouldEqual, ErrorTimeout)
		})
	})
}

