package checkmaincleanup

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
