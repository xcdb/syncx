package syncx_test

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/xcdb/syncx"
)

func ExampleAutoResetEvent() {
	//create event in a non-signaled state
	a := syncx.NewAutoResetEvent(false)

	//start a bunch of goroutines
	for i := 1; i <= 3; i++ {
		go func() {
			//...
			a.Wait()
			//...
		}()
	}

	//...

	//signal the event and release a goroutine
	a.Signal()

	//release another
	a.Signal()

	//release another
	a.Signal()

	//...

}

func ExampleAutoResetEvent_WaitContext() {
	ctx, cancel := context.WithCancel(context.Background())
	ready := sync.WaitGroup{}
	ready.Add(3)
	done := make(chan bool, 3)

	//create event in a non-signaled state
	a := syncx.NewAutoResetEvent(false)

	//start a bunch of goroutines
	for i := 1; i <= 3; i++ {
		go func() {
			ready.Done()
			//...
			err := a.WaitContext(ctx)
			if err == nil {
				fmt.Print("Signalled\n")
			} else if err == context.Canceled { //ctx.Err()
				fmt.Print("Cancelled\n")
			}
			//...
			done <- true
		}()
	}

	ready.Wait()

	//signal the event and release a goroutine
	a.Signal()

	//release another
	a.Signal()

	<-done
	<-done

	//cancel the context
	cancel()

	<-done

	// Output:
	// Signalled
	// Signalled
	// Cancelled
}

func ExampleManualResetEvent() {
	//create event in a non-signaled state
	m := syncx.NewManualResetEvent(false)

	//start a bunch of goroutines
	for i := 1; i <= 3; i++ {
		go func() {
			//...
			m.Wait()
			//...
		}()
	}

	//...

	//signal the event and release all goroutines
	m.Signal()

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
	a := syncx.NewAutoResetEvent(false)
	m := syncx.NewManualResetEvent(false)

	//start a bunch of goroutines
	for i := 1; i <= 5; i++ {
		go func() {
			//...
			syncx.WaitAny(a, m)
			//...
		}()
	}

	//signal the AutoResetevent, releasing a single goroutine
	a.Signal()

	//release another
	a.Signal()

	//signal the ManualResetEvent, releasing the remaining goroutines
	m.Signal()

	//...

}

func ExampleWaitAny_returnValue() {
	m := syncx.NewManualResetEvent(false)
	s := syncx.NewSemaphore(3)

	//start a bunch of goroutines
	//initially, only 3 can make progress at a time
	for i := 1; i <= 5; i++ {
		go func() {
			//...
			ix := syncx.WaitAny(m, s)
			//...

			//if we consumed a resource from the semaphone, we need to release it
			if ix == 1 {
				s.Release()
			}
		}()
	}

	//signal the event, releasing all remaining goroutines
	m.Signal()

	//...

}

func ExampleWaitAll() {
	a := syncx.NewAutoResetEvent(false)
	m := syncx.NewManualResetEvent(false)
	s := syncx.NewSemaphore(3)

	//start a bunch of goroutines
	for i := 1; i <= 5; i++ {
		go func() {
			//...
			syncx.WaitAll(a, m, s)
			//...
			s.Release()
		}()
	}

	//signal the AutoResetEvent (nothing is released as all 3 waits are required)
	a.Signal()

	//signal the ManualResetEvent, releasing a single goroutine
	m.Signal()

	//signal the AutoResetEvent, releasing another goroutine
	a.Signal()

	//...

}

func ExampleWaitAll_semantics() {
	a := syncx.NewAutoResetEvent(false)
	m := syncx.NewManualResetEvent(false)

	//start a bunch of goroutines
	for i := 1; i <= 5; i++ {
		go func() {
			//...
			syncx.WaitAll(a, m)
			//...
		}()
	}

	//signal the ManualResetEvent
	m.Signal()

	//set the ManualResetEvent to non-signalled
	m.Reset()

	//as the current implementation semantics are 'all have been signalled' this still releases a goroutine
	//this behaviour should not be depended upon
	a.Signal()

	//...

}
