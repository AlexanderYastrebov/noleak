package checkmain_test // note _test suffix

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

func TestLeakFromExternalPackage(t *testing.T) {
	l1 := newLeaky()
	l2 := newLeaky()

	l1.run()
	l2.run()

	l1.done()
	//l2.done()
}
