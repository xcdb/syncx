package syncx_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/xcdb/syncx"
)

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
		assert.Fail(t, "shouldn't be signalled")
	default:
	}
	e2.Signal()
	<-step //2
}

func TestWaitAll_consumesSingleResource(t *testing.T) {
	a := syncx.NewAutoResetEvent(false)
	s := syncx.NewSemaphore(2)
	s.Wait()
	s.Wait()

	step := make(chan int, 2)
	go func() {
		step <- 1
		syncx.WaitAll(a, s)
		step <- 2
	}()

	<-step //1
	s.Release()
	s.Release()

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*10)
	defer cancel()
	err := s.WaitContext(ctx)
	assert.Nil(t, err, "WaitAll consumed 2 resources from Semaphore")

	a.Signal()
	<-step //2
}

func TestWaitAll_panicsWhenTooMany(t *testing.T) {
	ws := make([]syncx.WaitHandle, 9, 9)
	for i := 0; i < len(ws); i++ {
		ws[i] = syncx.NewAutoResetEvent(false)
	}
	assert.Panics(t, func() { syncx.WaitAll(ws...) })
}

func TestWaitAllContext_panicsWhenTooMany(t *testing.T) {
	ws := make([]syncx.WaitHandle, 9, 9)
	for i := 0; i < len(ws); i++ {
		ws[i] = syncx.NewAutoResetEvent(false)
	}
	assert.Panics(t, func() { syncx.WaitAllContext(context.Background(), ws...) })
}

func TestWaitAll_returnsTrueWhenEmpty(t *testing.T) {
	b := syncx.WaitAll()
	assert.True(t, b)
}

func TestWaitAllContext_returnsTrueWhenEmpty(t *testing.T) {
	b, err := syncx.WaitAllContext(context.Background())
	assert.True(t, b)
	assert.Nil(t, err)
}

func TestWaitAll_returnsTrue(t *testing.T) {
	for l := 1; l <= 8; l++ {
		ws := make([]syncx.WaitHandle, l, l)
		for i := 0; i < l; i++ {
			ws[i] = syncx.NewManualResetEvent(true)
		}
		b := syncx.WaitAll(ws...)
		assert.True(t, b)
	}
}

func TestWaitAllContext_returnsTrue(t *testing.T) {
	for l := 1; l <= 8; l++ {
		ws := make([]syncx.WaitHandle, l, l)
		for i := 0; i < l; i++ {
			ws[i] = syncx.NewManualResetEvent(true)
		}
		b, err := syncx.WaitAllContext(context.Background(), ws...)
		assert.True(t, b)
		assert.Nil(t, err)
	}
}

func TestWaitAllContext_returnsFalseAndCtxErrWhenCancelled(t *testing.T) {
	for l := 1; l <= 8; l++ {
		ws := make([]syncx.WaitHandle, l, l)
		for i := 0; i < l; i++ {
			ws[i] = syncx.NewManualResetEvent(false)
		}
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		b, err := syncx.WaitAllContext(ctx, ws...)
		assert.False(t, b)
		assert.NotNil(t, err)
		assert.Equal(t, ctx.Err(), err)
	}
}
