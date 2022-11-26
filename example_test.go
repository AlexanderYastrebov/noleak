//go:build example

package noleak_test

import (
	"testing"

	"github.com/AlexanderYastrebov/noleak"
)

type leaky struct {
	ch chan struct{}
}

func newLeaky() *leaky {
	return &leaky{
		ch: make(chan struct{}),
	}
}

func (l leaky) run() {
	go func() {
		<-l.ch
	}()
}

func (l leaky) done() {
	close(l.ch)
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
