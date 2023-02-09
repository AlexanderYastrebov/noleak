//go:build example

package noleak_test

import (
	"os"
	"testing"
	"time"

	"github.com/AlexanderYastrebov/noleak"
)

var globalLeak = newLeaky()

func TestMain(m *testing.M) {
	if _, ok := os.LookupEnv("CLEANUP"); ok {
		testMainWithCleanup(m)
	} else {
		testMain(m)
	}
}

func testMain(m *testing.M) {
	os.Exit(noleak.CheckMain(m))
}

func testMainWithCleanup(m *testing.M) {
	os.Exit(noleak.CheckMainFunc(func() int {
		code := m.Run()

		// perform cleanup
		globalLeak.done()

		return code
	}))
}

// Detected by noleak.CheckMain(m)
func TestLeakUnchecked(t *testing.T) {
	l1 := newLeaky()
	l2 := newLeaky()

	l1.run()
	l2.run()

	//l1.done()
	//l2.done()
}

func TestLeakGlobal(t *testing.T) {
	globalLeak.run()

	//globalLeak.done()
}

// Detected by noleak.Check(t) and noleak.CheckMain(m)
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
