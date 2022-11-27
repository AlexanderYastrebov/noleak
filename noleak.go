package noleak

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"testing"
	"time"
)

var (
	checkTimeout = 1 * time.Second
)

func CheckMain(m *testing.M) {
	snapshot := routines()
	code := m.Run()
	active := snapshot.stillActiveAfter(time.Now().Add(checkTimeout))
	if len(active) > 0 {
		fmt.Printf("%d still active:\n%s", len(active), active.String())
		code = 1
	}
	os.Exit(code)
}

func Check(t *testing.T) {
	t.Helper()
	snapshot := routines()
	t.Cleanup(func() {
		t.Helper()

		active := snapshot.stillActiveAfter(time.Now().Add(checkTimeout))
		if len(active) > 0 {
			t.Errorf("%d still active:\n%s", len(active), active.String())
		}
	})
}

type goroutines map[string]string

func (gs goroutines) String() string {
	var b strings.Builder
	for _, g := range gs {
		if b.Len() > 0 {
			b.WriteString("\n\n")
		}
		b.WriteString(g)
	}
	return b.String()
}

func (gs goroutines) stillActiveAfter(deadline time.Time) goroutines {
	for time.Now().Before(deadline) {
		if len(gs.stillActive()) == 0 {
			return nil
		}
		time.Sleep(100 * time.Millisecond)
	}
	return gs.stillActive()
}

func (gs goroutines) stillActive() goroutines {
	active := routines()
	for id := range gs {
		delete(active, id)
	}
	return active
}

func routines() goroutines {
	gs := make(map[string]string)
	for _, g := range strings.Split(stack(), "\n\n") {
		header, _, _ := strings.Cut(g, "\n")
		// goroutine 8 [chan receive]:
		// goroutine 8 [runnable]:
		id, _, _ := strings.Cut(header, "[")
		gs[id] = g
	}
	return gs
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
