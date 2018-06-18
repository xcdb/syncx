package syncx_test

import (
	"fmt"
	"strings"
	"sync"

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

	var wg sync.WaitGroup
	wg.Add(3)

	for i := 1; i <= 3; i++ {
		go func() {
			for _, c := range strings.Split("ABCDE", "") {
				fmt.Print(c)
				b.SignalAndWait()
			}
			wg.Done()
		}()
	}

	wg.Wait()

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

func ExampleWaitAny() {
	e1 := syncx.NewAutoResetEvent(false)
	e2 := syncx.NewManualResetEvent(false)

	//start a bunch of goroutines
	for i := 1; i <= 5; i++ {
		go func() {
			//...
			syncx.WaitAny(e1, e2)
			//...
		}()
	}

	//signal the AutoResetevent, releasing a single goroutine
	e1.Signal()

	//release another
	e1.Signal()

	//signal the ManualResetEvent, releasing the remaining goroutines
	e2.Signal()

	//...

}

func ExampleWaitAny_returnValue() {
	e := syncx.NewManualResetEvent(false)
	s := syncx.NewSemaphore(3)

	//start a bunch of goroutines
	//initially, only 3 can make progress at a time
	for i := 1; i <= 5; i++ {
		go func() {
			//...
			ix := syncx.WaitAny(e, s)
			//...

			//if we consumed a resource from the semaphone, we need to release it
			if ix == 1 {
				s.Release()
			}
		}()
	}

	//signal the event, releasing all remaining goroutines
	e.Signal()

	//...

}

func ExampleWaitAll() {
	e1 := syncx.NewAutoResetEvent(false)
	e2 := syncx.NewManualResetEvent(false)
	s := syncx.NewSemaphore(3)

	//start a bunch of goroutines
	for i := 1; i <= 5; i++ {
		go func() {
			//...
			syncx.WaitAll(e1, e2, s)
			//...
			s.Release()
		}()
	}

	//signal the AutoResetEvent (nothing is released as all 3 waits are required)
	e1.Signal()

	//signal the ManualResetEvent, releasing a single goroutine
	e2.Signal()

	//signal the AutoResetEvent, releasing another goroutine
	e1.Signal()

	//...

}

func ExampleWaitAll_semantics() {
	e1 := syncx.NewAutoResetEvent(false)
	e2 := syncx.NewManualResetEvent(false)

	//start a bunch of goroutines
	for i := 1; i <= 5; i++ {
		go func() {
			//...
			syncx.WaitAll(e1, e2)
			//...
		}()
	}

	//signal the ManualResetEvent
	e2.Signal()

	//set the ManualResetEvent to non-signalled
	e2.Reset()

	//as the current implementation semantics are 'all have been signalled' this still releases a goroutine
	//this behaviour should not be depended upon
	e1.Signal()

	//...

}
