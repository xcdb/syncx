package syncx

import (
	"context"
	"runtime"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewAutoResetEvent(t *testing.T) {
	e1 := NewAutoResetEvent(false)
	assertNotSignalled(t, e1)

	e2 := NewAutoResetEvent(true)
	assertSignalled(t, e2)
}

func TestAutoResetEvent_Signal(t *testing.T) {
	e := NewAutoResetEvent(false)
	e.Signal()
	assertSignalled(t, e)
}

func TestAutoResetEvent_Reset(t *testing.T) {
	e := NewAutoResetEvent(true)
	e.Reset()
	assertNotSignalled(t, e)
}

func TestAutoResetEvent_Wait_resetsSignal(t *testing.T) {
	e := NewAutoResetEvent(true)
	e.Wait()
	assertNotSignalled(t, e)
}

func TestAutoResetEvent_WaitContext_resetsSignal(t *testing.T) {
	e := NewAutoResetEvent(true)
	e.WaitContext(context.Background())
	assertNotSignalled(t, e)
}

//ensures that Signal must be called once per goroutine
func TestAutoResetEvent_Signal_wakeOne(t *testing.T) {
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

func TestAutoResetEvent_Wait_nonsignalled(t *testing.T) {
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

func TestAutoResetEvent_Wait_signalled(t *testing.T) {
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

func TestAutoResetEvent_WaitContext_nonsignalled(t *testing.T) {
	e := NewAutoResetEvent(false)

	step := make(chan int, 1)
	go func() {
		step <- 1
		err := e.WaitContext(context.Background())
		assert.Nil(t, err)
		step <- 2
	}()

	<-step //1
	e.Signal()
	<-step //2
}

func TestAutoResetEvent_WaitContext_signalled(t *testing.T) {
	e := NewAutoResetEvent(true)

	step := make(chan int, 1)
	go func() {
		step <- 1
		err := e.WaitContext(context.Background())
		assert.Nil(t, err)
		step <- 2
	}()

	<-step //1
	//e.Signal()
	<-step //2
}

func TestAutoResetEvent_WaitContext_returnsCtxErrWhenCtxDone(t *testing.T) {
	e := NewAutoResetEvent(false)

	ctx, cancel := context.WithCancel(context.Background())

	step := make(chan int, 1)
	go func() {
		step <- 1
		err := e.WaitContext(ctx)
		assert.NotNil(t, err)
		assert.Equal(t, ctx.Err(), err)
		step <- 2
	}()

	<-step //1
	cancel()
	<-step //2
}

//...

//Warning: assertSignalled can potentially return w to non-signalled
func assertSignalled(t *testing.T, w WaitHandle, msgAndArgs ...interface{}) {
	select {
	case <-w.ch():
		return
	default:
		assert.Fail(t, "", msgAndArgs)
	}
}

//Warning: assertSignalled can potentially return w to non-signalled
func assertNotSignalled(t *testing.T, w WaitHandle, msgAndArgs ...interface{}) {
	select {
	case <-w.ch():
		assert.Fail(t, "", msgAndArgs)
	default:
		return
	}
}
