package syncx_test

import (
	"fmt"
	"strings"

	"github.com/xcdb/syncx"
)

func ExampleAutoResetEvent() {
	//create event in a non-signaled state
	e := syncx.NewAutoResetEvent(false)

	//start a bunch of goroutines
	for i := 1; i <= 3; i++ {
		go func() {
			//...
			e.Wait()
			//...
		}()
	}

	//...

	//signal the event and release a goroutine
	e.Signal()

	//release another
	e.Signal()

	//release another
	e.Signal()

	//...

}

func ExampleManualResetEvent() {
	//create event in a non-signaled state
	e := syncx.NewManualResetEvent(false)

	//start a bunch of goroutines
	for i := 1; i <= 3; i++ {
		go func() {
			//...
			e.Wait()
			//...
		}()
	}

	//...

	//signal the event and release all goroutines
	e.Signal()

	//...

}

func ExampleBarrier() {
	b := syncx.NewBarrier(3, func() { fmt.Print(" ") })

	done := make(chan bool, 3)
	for i := 1; i <= 3; i++ {
		go func() {
			for _, c := range strings.Split("ABCDE", "") {
				fmt.Print(c)
				b.SignalAndWait()
			}
			done <- true
		}()
	}

	for i := 1; i <= 3; i++ {
		<-done
	}

	// Output:
	// AAA BBB CCC DDD EEE
}

func ExampleBarrier_Add() {
	b := syncx.NewBarrier(3, func() { fmt.Print(" ") })

	done := make(chan bool, 3)
	for i := 1; i <= 3; i++ {
		go func(i int) {
			for _, c := range strings.Split("ABCDE", "") {
				fmt.Print(c)
				b.SignalAndWait()
				if b.Phase() == i*2 { //drop a goroutine at phases 2 and 4
					b.Add(-1)
					break
				}
			}
			done <- true
		}(i)
	}

	for i := 1; i <= 3; i++ {
		<-done
	}

	// Output:
	// AAA BB CC D E
}

func ExampleSemaphore() {
	//create semaphore able to satisfy 3 concurrent goroutines
	s := syncx.NewSemaphore(3)

	//start a bunch of goroutines
	//only 3 can make progress at a time
	for i := 1; i <= 5; i++ {
		go func() {
			s.Wait()
			//...
			s.Release()
		}()
	}

	//...

}
