package syncx

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

//TestSemaphore ensures that Semaphore correctly limits the number of concurrent goroutines
func TestSemaphore(t *testing.T) {
	s := NewSemaphore(3)

	entered := make(chan bool, 5)
	release := make(chan bool, 5)
	for i := 1; i <= 5; i++ {
		go func() {
			s.Wait()
			entered <- true
			<-release
			s.Release()
		}()
	}

	<-entered
	<-entered
	<-entered

	for i := 1; i <= 2; i++ {
		time.Sleep(1 * time.Microsecond) //yield so that if extra goroutines are not being blocked, they have opportunity to run
		select {
		case <-entered:
			assert.FailNow(t, "Semaphore not limiting concurrency as expected")
		default:
		}
		release <- true //allow a single slot to be released from the Semaphore
		<-entered
	}
}

func TestNewSemaphore(t *testing.T) {
	count := []int{1, 42, 1024}
	for _, c := range count {
		s := NewSemaphore(c)
		assert.Equal(t, c, len(s.c))
	}
}

func TestNewSemaphore_panics(t *testing.T) {
	assert.Panics(t, func() { NewSemaphore(0) })
	assert.Panics(t, func() { NewSemaphore(-1) })
}

func TestSemaphore_WaitContext_returnsCtxErrWhenCtxDone(t *testing.T) {
	s := NewSemaphore(3)
	s.Wait()
	s.Wait()
	s.Wait()

	ctx, cancel := context.WithCancel(context.Background())

	step := make(chan int, 1)
	go func() {
		step <- 1
		err := s.WaitContext(ctx)
		assert.NotNil(t, err)
		assert.Equal(t, ctx.Err(), err)
		step <- 2
	}()

	<-step //1
	cancel()
	<-step //2
}
