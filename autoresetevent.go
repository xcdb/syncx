package syncx

import "sync"

//AutoResetEvent notifies a waiting goroutine that an event has occurred.
//
//Once it has been signaled, AutoResetEvent remains signaled until a single waiting goroutine is awoken, and then automatically returns to the non-signaled state.
//
//There is no guarantee that every call to Signal will wake a waiting goroutine.
//If two calls are too close together, so that the second call occurs before a goroutine has awoken, it is as if the second call did not happen.
//Also, if Set is called when there are no waiting goroutines, and e is already signaled, the call has no effect.
type AutoResetEvent struct {
	state bool
	cond  *sync.Cond
}

//NewAutoResetEvent returns a new AutoResetEvent with initial state s
func NewAutoResetEvent(s bool) *AutoResetEvent {
	cond := sync.NewCond(&sync.Mutex{})
	return &AutoResetEvent{
		state: s,
		cond:  cond,
	}
}

//Signal sets the state of e to signaled, waking a waiting goroutine.
func (e *AutoResetEvent) Signal() {
	e.cond.L.Lock()
	e.state = true
	e.cond.L.Unlock()
	e.cond.Signal()
}

//Reset sets the state of e to nonsignaled.
func (e *AutoResetEvent) Reset() {
	e.cond.L.Lock()
	e.state = false
	e.cond.L.Unlock()
}

//Wait suspends execution of the calling goroutine until e receives a signal.
func (e *AutoResetEvent) Wait() {
	e.cond.L.Lock()
	for e.state == false {
		e.cond.Wait()
	}
	e.state = false
	e.cond.L.Unlock()
}
