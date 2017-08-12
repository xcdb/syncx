package syncx

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

//TestManualResetEvent_1 ensures that Wait blocks if not signaled
func TestManualResetEvent_1(t *testing.T) {
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

//TestManualResetEvent_2 ensures that Wait doesnt block if state is signaled
func TestManualResetEvent_2(t *testing.T) {
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

//TestManualResetEvent_3 ensures that all waiting goroutines are awoken by a single call to Signal
func TestManualResetEvent_3(t *testing.T) {
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

func TestNewManualResetEvent_SetsState(t *testing.T) {
	state := []bool{false, true}
	for _, b := range state {
		e := NewManualResetEvent(b)
		assert.Equal(t, b, e.state)
	}
}

func TestManualResetEvent_Signal_SetsStateToTrue(t *testing.T) {
	e := NewManualResetEvent(false)
	e.Signal()
	assert.True(t, e.state)
}

func TestManualResetEvent_Reset_SetsStateToFalse(t *testing.T) {
	e := NewManualResetEvent(true)
	e.Reset()
	assert.False(t, e.state)
}

func TestManualResetEvent_Wait_LeavesStateAsTrue(t *testing.T) {
	e := NewManualResetEvent(true)
	e.Wait()
	assert.True(t, e.state)
}
