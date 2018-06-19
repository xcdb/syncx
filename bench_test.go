package syncx

import (
	"reflect"
	"testing"
)

var ix int
var e []*ManualResetEvent
var w []WaitHandle

func init() {
	e = make([]*ManualResetEvent, 8)
	w = make([]WaitHandle, 8)
	for i := 0; i < 8; i++ {
		e[i] = NewManualResetEvent(false)
		w[i] = e[i]
	}
}

func BenchmarkWaitAny__select1(b *testing.B) {
	e[0].Signal()
	for n := 0; n < b.N; n++ {
		select {
		case <-e[0].ch():
			ix = 0
		}
	}
}
func BenchmarkWaitAny__select2(b *testing.B) {
	e[1].Signal()
	for n := 0; n < b.N; n++ {
		select {
		case <-e[0].ch():
			ix = 0
		case <-e[1].ch():
			ix = 1
		}
	}
}
func BenchmarkWaitAny__select3(b *testing.B) {
	e[2].Signal()
	for n := 0; n < b.N; n++ {
		select {
		case <-e[0].ch():
			ix = 0
		case <-e[1].ch():
			ix = 1
		case <-e[2].ch():
			ix = 2
		}
	}
}
func BenchmarkWaitAny__select4(b *testing.B) {
	e[3].Signal()
	for n := 0; n < b.N; n++ {
		select {
		case <-e[0].ch():
			ix = 0
		case <-e[1].ch():
			ix = 1
		case <-e[2].ch():
			ix = 2
		case <-e[3].ch():
			ix = 3
		}
	}
}
func BenchmarkWaitAny__select8(b *testing.B) {
	e[7].Signal()
	for n := 0; n < b.N; n++ {
		select {
		case <-e[0].ch():
			ix = 0
		case <-e[1].ch():
			ix = 1
		case <-e[2].ch():
			ix = 2
		case <-e[3].ch():
			ix = 3
		case <-e[4].ch():
			ix = 4
		case <-e[5].ch():
			ix = 5
		case <-e[6].ch():
			ix = 6
		case <-e[7].ch():
			ix = 7
		}
	}
}

func BenchmarkWaitAny_reflect1(b *testing.B) {
	e[0].Signal()
	for n := 0; n < b.N; n++ {
		cases := make([]reflect.SelectCase, 1)
		for i := 0; i < len(cases); i++ {
			cases[i] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(e[i].ch())}
		}
		ix, _, _ = reflect.Select(cases)
	}
}
func BenchmarkWaitAny_reflect2(b *testing.B) {
	e[1].Signal()
	for n := 0; n < b.N; n++ {
		cases := make([]reflect.SelectCase, 2)
		for i := 0; i < len(cases); i++ {
			cases[i] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(e[i].ch())}
		}
		ix, _, _ = reflect.Select(cases)
	}
}
func BenchmarkWaitAny_reflect3(b *testing.B) {
	e[2].Signal()
	for n := 0; n < b.N; n++ {
		cases := make([]reflect.SelectCase, 3)
		for i := 0; i < len(cases); i++ {
			cases[i] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(e[i].ch())}
		}
		ix, _, _ = reflect.Select(cases)
	}
}
func BenchmarkWaitAny_reflect4(b *testing.B) {
	e[3].Signal()
	for n := 0; n < b.N; n++ {
		cases := make([]reflect.SelectCase, 4)
		for i := 0; i < len(cases); i++ {
			cases[i] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(e[i].ch())}
		}
		ix, _, _ = reflect.Select(cases)
	}
}
func BenchmarkWaitAny_reflect8(b *testing.B) {
	e[3].Signal()
	for n := 0; n < b.N; n++ {
		cases := make([]reflect.SelectCase, 8)
		for i := 0; i < len(cases); i++ {
			cases[i] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(e[i].ch())}
		}
		ix, _, _ = reflect.Select(cases)
	}
}

func BenchmarkWaitAny____impl1(b *testing.B) {
	e[0].Signal()
	for n := 0; n < b.N; n++ {
		ix = WaitAny(w[0:1]...)
	}
}
func BenchmarkWaitAny____impl2(b *testing.B) {
	e[1].Signal()
	for n := 0; n < b.N; n++ {
		ix = WaitAny(w[0:2]...)
	}
}
func BenchmarkWaitAny____impl3(b *testing.B) {
	e[2].Signal()
	for n := 0; n < b.N; n++ {
		ix = WaitAny(w[0:3]...)
	}
}
func BenchmarkWaitAny____impl4(b *testing.B) {
	e[3].Signal()
	for n := 0; n < b.N; n++ {
		ix = WaitAny(w[0:4]...)
	}
}
func BenchmarkWaitAny____impl8(b *testing.B) {
	e[3].Signal()
	for n := 0; n < b.N; n++ {
		ix = WaitAny(w[0:8]...)
	}
}
