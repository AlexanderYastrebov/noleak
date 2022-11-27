//go:build example

package noleak_test

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

func (l *leaky) doneAfter(d time.Duration) {
	l.ch <- d
}

func (l *leaky) done() {
	l.doneAfter(10 * time.Millisecond)
}

func TestWetDisabled(t *testing.T) {
	//noleak.Check(t)

	l1 := newLeaky()
	l2 := newLeaky()

	l1.run()
	l2.run()

	//l1.done()
	//l2.done()
}

func TestWet(t *testing.T) {
	noleak.Check(t)

	l1 := newLeaky()
	l2 := newLeaky()

	l1.run()
	l2.run()

	l1.done()
	//l2.done()
}

func TestDry(t *testing.T) {
	noleak.Check(t)

	l1 := newLeaky()
	l2 := newLeaky()

	l1.run()
	l2.run()

	l1.done()
	l2.done()
}
