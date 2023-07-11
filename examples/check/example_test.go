package check

import (
	"testing"
	"time"

	"github.com/AlexanderYastrebov/noleak"
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

func TestLeak(t *testing.T) {
	noleak.Check(t)

	l1 := newLeaky()
	l2 := newLeaky()

	l1.run()
	l2.run()

	l1.done()
	//l2.done()
}

func TestNoLeak(t *testing.T) {
	noleak.Check(t)

	l1 := newLeaky()
	l2 := newLeaky()

	l1.run()
	l2.run()

	l1.done()
	l2.done()
}
