package noleak

import (
	"runtime"
	"strings"
	"testing"
)

func Check(t *testing.T) {
	t.Helper()
	before := routines()
	t.Cleanup(func() {
		t.Helper()
		//time.Sleep(100 * time.Millisecond)

		after := routines()
		for header := range before {
			delete(after, header)
		}
		if len(after) == 0 {
			return
		}
		var stacks strings.Builder
		for _, g := range after {
			if stacks.Len() > 0 {
				stacks.WriteString("\n\n")
			}
			stacks.WriteString(g)
		}
		t.Errorf("%d still active:\n%s", len(after), stacks.String())
	})
}

func routines() map[string]string {
	result := make(map[string]string)
	for _, g := range strings.Split(stack(), "\n\n") {
		header, _, _ := strings.Cut(g, "\n")
		result[header] = g
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
