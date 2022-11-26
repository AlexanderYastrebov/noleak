package noleak

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"testing"
)

func CheckMain(m *testing.M) {
	before := routines()
	code := m.Run()
	active := routines().subtract(before)
	if len(active) > 0 {
		fmt.Printf("%d still active:\n%s", len(active), active.String())
		code = 1
	}
	os.Exit(code)
}

func Check(t *testing.T) {
	t.Helper()
	before := routines()
	t.Cleanup(func() {
		t.Helper()
		//time.Sleep(100 * time.Millisecond)

		active := routines().subtract(before)
		if len(active) > 0 {
			t.Errorf("%d still active:\n%s", len(active), active.String())
		}
	})
}

type goroutines map[string]string

func (a goroutines) String() string {
	var b strings.Builder
	for _, g := range a {
		if b.Len() > 0 {
			b.WriteString("\n\n")
		}
		b.WriteString(g)
	}
	return b.String()
}

func (a goroutines) subtract(b goroutines) goroutines {
	for k := range b {
		delete(a, k)
	}
	return a
}

func routines() goroutines {
	result := make(map[string]string)
	for _, g := range strings.Split(stack(), "\n\n") {
		header, _, _ := strings.Cut(g, "\n")
		// goroutine 8 [chan receive]:
		// goroutine 8 [runnable]:
		id, _, _ := strings.Cut(header, "[")
		result[id] = g
	}
	return result
}

func stack() string {
	buf := make([]byte, 1024)
	for {
		n := runtime.Stack(buf, true)
		if n < len(buf) {
			return string(buf[:n])
		}
		buf = make([]byte, 2*len(buf))
	}
}
