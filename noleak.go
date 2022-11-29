package noleak

import (
	"fmt"
	"runtime"
	"strings"
	"testing"
	"time"
)

const checkTimeout = 1 * time.Second

func CheckMain(m *testing.M) (code int) {
	snapshot := routines()
	code = m.Run()
	active := snapshot.stillActiveAfter(time.Now().Add(checkTimeout))
	if len(active) > 0 {
		fmt.Println(active)
		code = 1
	}
	return
}

func Check(t *testing.T) {
	t.Helper()
	snapshot := routines()
	t.Cleanup(func() {
		t.Helper()
		active := snapshot.stillActiveAfter(time.Now().Add(checkTimeout))
		if len(active) > 0 {
			t.Error(active)
		}
	})
}

type goroutines map[string]string

func (gs goroutines) String() string {
	var b strings.Builder
	fmt.Fprintf(&b, "noleak: %d active", len(gs))
	for _, g := range gs {
		b.WriteString("\n\n")
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
