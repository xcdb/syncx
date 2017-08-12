package syncx

import "sync"

//Semaphore limits the number of goroutines that can access a resource or pool of resources concurrently.
type Semaphore struct {
	count int
	cond  *sync.Cond
}

//NewSemaphore returns a new Semaphore with count c
//
//It panics if c is less than 1.
func NewSemaphore(c int) *Semaphore {
	if c < 1 {
		panic("syncx: NewSemaphore c is less than 1")
	}
	cond := sync.NewCond(&sync.Mutex{})
	return &Semaphore{
		count: c,
		cond:  cond,
	}
}

//Wait suspends execution of the calling goroutine until it can enter s.
func (s *Semaphore) Wait() {
	s.cond.L.Lock()
	for s.count <= 0 {
		s.cond.Wait()
	}
	s.count--
	s.cond.L.Unlock()
}

//Release exits s, waking a waiting goroutine.
func (s *Semaphore) Release() {
	s.cond.L.Lock()
	s.count++
	s.cond.L.Unlock()
	s.cond.Signal()
}
