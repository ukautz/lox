[![Build Status](https://travis-ci.org/ukautz/lox.svg?branch=master)](https://travis-ci.org/ukautz/lox)

This package provides a lock which can timeout and a convenience container
supporting named locks.

Documentation
-------------

GoDoc can be [found here](http://godoc.org/github.com/ukautz/lox)

Example
-------

``` go

fails := []error{}
ok := 0
var wg sync.WaitGroup
rand.Seed(int64(time.Now().Nanosecond()))
c := lox.NewLocks()

// each of the Runs wait concurrently for 100ms
timeout := time.Duration(100) * time.Millisecond
for i := 0; i < 100; i++ {
    wg.Add(1)
    go func() {
        defer wg.Done()

        // unlock the append in any case.. just will return error, which
        // we do not catch, if the run does not return an error
        defer c.Unlock("append")

        // try running .. wait at most for 100 ms
        err := c.Run("foo", func() {
            // sleeping on average 5 ms
            time.Sleep(time.Duration(rand.Intn(10)) * time.Millisecond)
            ok++
        }, timeout)

        // run can fail (= not run)
        if err != nil {
            // it's ok to put the lock in the if clause, since the above
            // unlock does not panic if it was not locked
            c.Lock("append")

            fails = append(fails, err)
        }
    }()
}

wg.Wait()

// expect about 15-25 or so to run because of an (approximate) average
// of 5 ms sleep time and each waits for 100 ms concurrently
fmt.Printf("Failed: %d\nOk: %d\n", len(fails), ok)
```