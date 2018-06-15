package syncx_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/xcdb/syncx"
)

func TestWaitAny(t *testing.T) {
	e1 := syncx.NewAutoResetEvent(false)
	e2 := syncx.NewManualResetEvent(false)
	e3 := syncx.NewSemaphore(3)

	step := make(chan int, 1)
	go func() {
		step <- 1
		syncx.WaitAny(e1, e2, e3)
		step <- 2
	}()

	<-step //1
	<-step //2
}

func TestWaitAll(t *testing.T) {
	e1 := syncx.NewAutoResetEvent(false)
	e2 := syncx.NewManualResetEvent(false)
	e3 := syncx.NewSemaphore(3)

	step := make(chan int, 1)
	go func() {
		step <- 1
		syncx.WaitAll(e1, e2, e3)
		step <- 2
	}()

	<-step //1
	e1.Signal()
	e1.Signal() //checking that each only counts once
	select {
	case <-step:
		assert.Fail(t, "e2 not signalled")
	default:
	}
	e2.Signal()
	<-step //2
}
