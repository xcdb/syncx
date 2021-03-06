package syncx_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/xcdb/syncx"
)

func TestWaitAny(t *testing.T) {
	e1 := syncx.NewAutoResetEvent(false)
	e2 := syncx.NewManualResetEvent(false)
	e3 := syncx.NewSemaphore(1)
	e3.Wait()

	step := make(chan int, 1)
	go func() {
		step <- 1
		syncx.WaitAny(e1, e2, e3)
		step <- 2
	}()

	<-step //1
	select {
	case <-step:
		assert.Fail(t, "shouldn't be signalled")
	default:
	}
	e3.Release()
	<-step //2
}

func TestWaitAny_panicsWhenTooMany(t *testing.T) {
	ws := make([]syncx.WaitHandle, 9, 9)
	for i := 0; i < len(ws); i++ {
		ws[i] = syncx.NewAutoResetEvent(false)
	}
	assert.Panics(t, func() { syncx.WaitAny(ws...) })
}

func TestWaitAnyContext_panicsWhenTooMany(t *testing.T) {
	ws := make([]syncx.WaitHandle, 9, 9)
	for i := 0; i < len(ws); i++ {
		ws[i] = syncx.NewAutoResetEvent(false)
	}
	assert.Panics(t, func() { syncx.WaitAnyContext(context.Background(), ws...) })
}

func TestWaitAny_returnsNegative1WhenEmpty(t *testing.T) {
	ix := syncx.WaitAny()
	assert.Equal(t, -1, ix)
}

func TestWaitAnyContext_returnsNegative1WhenEmpty(t *testing.T) {
	ix, err := syncx.WaitAnyContext(context.Background())
	assert.Equal(t, -1, ix)
	assert.Nil(t, err)
}

func TestWaitAny_returnsIndexThatSatisfiedWait(t *testing.T) {
	for l := 1; l <= 8; l++ {
		es := make([]*syncx.AutoResetEvent, l, l)
		ws := make([]syncx.WaitHandle, l, l)
		for i := 0; i < l; i++ {
			es[i] = syncx.NewAutoResetEvent(false)
			ws[i] = es[i]
		}
		for j := 0; j < l; j++ {
			es[j].Signal()
			ix := syncx.WaitAny(ws...)
			assert.Equal(t, j, ix)
		}
	}
}

func TestWaitAnyContext_returnsIndexThatSatisfiedWait(t *testing.T) {
	for l := 1; l <= 8; l++ {
		es := make([]*syncx.AutoResetEvent, l, l)
		ws := make([]syncx.WaitHandle, l, l)
		for i := 0; i < l; i++ {
			es[i] = syncx.NewAutoResetEvent(false)
			ws[i] = es[i]
		}
		for j := 0; j < l; j++ {
			es[j].Signal()
			ix, err := syncx.WaitAnyContext(context.Background(), ws...)
			assert.Equal(t, j, ix)
			assert.Nil(t, err)
		}
	}
}

func TestWaitAnyContext_returnsNegative1AndCtxErrWhenCancelled(t *testing.T) {
	for l := 1; l <= 8; l++ {
		ws := make([]syncx.WaitHandle, l, l)
		for i := 0; i < l; i++ {
			ws[i] = syncx.NewAutoResetEvent(false)
		}
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		ix, err := syncx.WaitAnyContext(ctx, ws...)
		assert.Equal(t, -1, ix)
		assert.NotNil(t, err)
		assert.Equal(t, ctx.Err(), err)
	}
}
