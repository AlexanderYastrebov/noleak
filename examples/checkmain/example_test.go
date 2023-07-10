package checkmain

import (
	"testing"
	"time"
)

type leaky struct {
	ch chan time.Duration
}

func newLeaky() *leaky {
	return &leaky{
		ch: make(chan time.Duration),
	}
}

func (l *leaky) run() {
	go func() {
		time.Sleep(<-l.ch)
	}()
}

func (l *leaky) done() {
	l.ch <- 10 * time.Millisecond
}

var globalLeak = newLeaky()

func TestLeakGlobal(t *testing.T) {
	globalLeak.run()

	//globalLeak.done()
}

func TestLeak(t *testing.T) {
	l1 := newLeaky()
	l1.run()
	//l1.done()
}

func TestLeakNested(t *testing.T) {
	t.Run("nested", func(t *testing.T) {
		l1 := newLeaky()
		l1.run()
		//l1.done()
	})
}

func TestLeakInner(t *testing.T) {
	doLeak(5)
}

func TestLeakNestedInner(t *testing.T) {
	t.Run("nested", func(t *testing.T) {
		doLeak(5)
	})
}

func doLeak(n int) {
	if n == 0 {
		l1 := newLeaky()
		l1.run()
		//l1.done()
	} else {
		doLeak(n - 1)
	}
}
