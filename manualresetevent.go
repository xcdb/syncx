package syncx

import "sync"

//ManualResetEvent notifies one or more waiting goroutines that an event has occurred.
//
//Once it has been signaled, ManualResetEvent remains signaled until it is manually reset.
//When signaled, all waiting goroutines are released, and all calls to Wait return immediately.
type ManualResetEvent struct {
	state bool
	l     *sync.RWMutex
	cond  *sync.Cond
}

//NewManualResetEvent returns a new ManualResetEvent with initial state s
func NewManualResetEvent(s bool) *ManualResetEvent {
	l := &sync.RWMutex{}
	cond := sync.NewCond(l.RLocker())
	return &ManualResetEvent{
		state: s,
		l:     l,
		cond:  cond,
	}
}

//Signal sets the state of e to signaled, waking one or more waiting goroutines.
func (e *ManualResetEvent) Signal() {
	e.l.Lock()
	e.state = true
	e.l.Unlock()
	e.cond.Broadcast()
}

//Reset sets the state of e to nonsignaled.
func (e *ManualResetEvent) Reset() {
	e.l.Lock()
	e.state = false
	e.l.Unlock()
}

//Wait suspends execution of the calling goroutine until e receives a signal.
func (e *ManualResetEvent) Wait() {
	e.cond.L.Lock()
	for e.state == false {
		e.cond.Wait()
	}
	e.cond.L.Unlock()
}
