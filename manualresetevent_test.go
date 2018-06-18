package syncx

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewManualResetEvent(t *testing.T) {
	e1 := NewManualResetEvent(false)
	assertNotSignalled(t, e1)

	e2 := NewManualResetEvent(true)
	assertSignalled(t, e2)
}

func TestManualResetEvent_Signal(t *testing.T) {
	e := NewManualResetEvent(false)
	e.Signal()
	assertSignalled(t, e)
}

func TestManualResetEvent_Reset(t *testing.T) {
	e := NewManualResetEvent(true)
	e.Reset()
	assertNotSignalled(t, e)
}

func TestManualResetEvent_Wait_retainsSignal(t *testing.T) {
	e := NewManualResetEvent(true)
	e.Wait()
	assertSignalled(t, e)
}

func TestManualResetEvent_WaitContext_retainsSignal(t *testing.T) {
	e := NewManualResetEvent(true)
	e.WaitContext(context.Background())
	assertSignalled(t, e)
}

//ensures that all waiting goroutines are awoken by a single call to Signal
func TestManualResetEvent_Signal_wakeAll(t *testing.T) {
	e := NewManualResetEvent(false)

	done := make(chan bool, 3)
	for i := 1; i <= 3; i++ {
		go func() {
			e.Wait()
			done <- true
		}()
	}

	e.Signal()

	for i := 1; i <= 3; i++ {
		<-done
	}
}

func TestManualResetEvent_Wait_nonsignalled(t *testing.T) {
	e := NewManualResetEvent(false)

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

func TestManualResetEvent_Wait_signalled(t *testing.T) {
	e := NewManualResetEvent(true)

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

func TestManualResetEvent_WaitContext_nonsignalled(t *testing.T) {
	e := NewManualResetEvent(false)

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

func TestManualResetEvent_WaitContext_signalled(t *testing.T) {
	e := NewManualResetEvent(true)

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

func TestManualResetEvent_WaitContext_returnsCtxErrWhenCtxDone(t *testing.T) {
	e := NewManualResetEvent(false)

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
