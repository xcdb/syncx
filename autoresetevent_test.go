package syncx

import (
	"runtime"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

//TestAutoResetEvent_1 ensures that Wait blocks if not signaled
func TestAutoResetEvent_1(t *testing.T) {
	e := NewAutoResetEvent(false)

	step := make(chan int, 1)
	go func() {
		step <- 1
		e.Wait()
		step <- 2
	}()

	<-step //1
	e.Signal()
	<-step //2
}

//TestAutoResetEvent_2 ensures that Wait doesnt block if state is signaled
func TestAutoResetEvent_2(t *testing.T) {
	e := NewAutoResetEvent(true)

	step := make(chan int, 1)
	go func() {
		step <- 1
		e.Wait()
		step <- 2
	}()

	<-step //1
	//e.Signal()
	<-step //2
}

//TestAutoResetEvent_3 ensures that Signal must be called once per goroutine
func TestAutoResetEvent_3(t *testing.T) {
	runtime.GOMAXPROCS(4)
	var c int64

	e := NewAutoResetEvent(false)

	done := make(chan bool, 3)
	for i := 1; i <= 3; i++ {
		go func() {
			e.Wait()
			atomic.AddInt64(&c, 1)
			done <- true
		}()
	}

	for i := 1; i <= 3; i++ {
		e.Signal()
		time.Sleep(1 * time.Microsecond)
		assert.Equal(t, int64(i), c)
	}

	for i := 1; i <= 3; i++ {
		<-done
	}
}

//TODO: rename these to match implementation

func TestNewAutoResetEvent_SetsState(t *testing.T) {
	state := []bool{false, true}
	for _, b := range state {
		e := NewAutoResetEvent(b)
		s := len(e.c) == 1
		assert.Equal(t, b, s)
	}
}

func TestAutoResetEvent_Signal_SetsStateToTrue(t *testing.T) {
	e := NewAutoResetEvent(false)
	e.Signal()
	s := len(e.c) == 1
	assert.True(t, s)
}

func TestAutoResetEvent_Reset_SetsStateToFalse(t *testing.T) {
	e := NewAutoResetEvent(true)
	e.Reset()
	s := len(e.c) == 1
	assert.False(t, s)
}

func TestAutoResetEvent_Wait_SetsStateToFalse(t *testing.T) {
	e := NewAutoResetEvent(true)
	e.Wait()
	s := len(e.c) == 1
	assert.False(t, s)
}
