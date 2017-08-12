package syncx

import (
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestBarrier(t *testing.T) {
	b := NewBarrier(2, nil)

	var p int64
	step := make(chan int, 10)

	for i := 1; i <= 2; i++ {
		go func() {
			for i := 1; i <= 5; i++ {
				<-step
				atomic.StoreInt64(&p, int64(i))
				b.SignalAndWait()
			}
		}()
	}

	for i := 1; i <= 10; i++ {
		step <- i
		time.Sleep(1 * time.Microsecond)
		assert.Equal(t, atomic.LoadInt64(&p), int64((i/2)+(i%2)))
	}
}

func TestBarrier_Add_PanicsIfCountGoesNegative(t *testing.T) {
	b := NewBarrier(3, nil)
	assert.Panics(t, func() { b.Add(-4) })
}

func TestNewBarrier(t *testing.T) {
	b := NewBarrier(3, nil)
	assert.Equal(t, 1, b.phase)
	assert.Equal(t, 3, b.participants)
}

func TestNewBarrier_PanicsIfCountLessThan0(t *testing.T) {
	NewBarrier(0, nil)
	assert.Panics(t, func() { NewBarrier(-1, nil) })
}
