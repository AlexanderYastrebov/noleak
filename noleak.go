package noleak

import (
	"fmt"
	"runtime"
	"sort"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

const checkTimeout = 5 * time.Second

// CheckMain prints active goroutines after all tests end.
// It returns result of m.Run() or non-zero if there are active goroutines.
//
// Example:
//
//	func TestMain(m *testing.M) {
//		os.Exit(noleak.CheckMain(m))
//	}
func CheckMain(m *testing.M) (code int) {
	return checkMainFunc(m.Run)
}

// CheckMainFunc prints active goroutines after m ends.
// It returns result of m or non-zero if there are active goroutines.
//
// Example:
//
//	func TestMain(m *testing.M) {
//		os.Exit(noleak.CheckMainFunc(func() int {
//			code := m.Run()
//
//			// perform cleanup
//			...
//
//			return code
//		}))
//	}
func CheckMainFunc(m func() int) (code int) {
	return checkMainFunc(m)
}

func checkMainFunc(m func() int) (code int) {
	snapshot := routines()
	code = m()
	active := snapshot.stillActiveAfter(time.Now().Add(checkTimeout))
	if len(active) > 0 {
		code = 1

		fmt.Println(active)
		if pkg := callerPackage(2); pkg != "" {
			if tests := findTests(pkg, active); len(tests) > 0 {
				fmt.Println("leaked from:")
				for _, name := range tests {
					fmt.Print("\t")
					fmt.Println(name)
				}
			}
		}
		fmt.Println("Try setting GODEBUG=tracebackancestors=10 (or more), see https://pkg.go.dev/runtime#hdr-Environment_Variables")
	}
	return
}

// Check reports test error if there are active goroutines after test ends.
//
// Example:
//
//	func TestLeak(t *testing.T) {
//		noleak.Check(t)
//
//		...
//	}
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

var bufferSize int64 = 1024

func stack() string {
	for {
		bs := atomic.LoadInt64(&bufferSize)
		buf := make([]byte, bs)
		n := runtime.Stack(buf, true)
		if n < len(buf) {
			return string(buf[:n])
		}
		atomic.CompareAndSwapInt64(&bufferSize, bs, 2*bs)
	}
}

func callerPackage(skip int) string {
	if pc, _, _, ok := runtime.Caller(skip + 1); ok {
		if f := runtime.FuncForPC(pc); f != nil {
			pkg, _ := splitPackageAndName(f.Name())
			return pkg
		}
	}
	return ""
}

func splitPackageAndName(name string) (string, string) {
	lastSlashAt := strings.LastIndexByte(name, '/')
	if lastSlashAt == -1 {
		lastSlashAt = 0
	}
	dotAt := strings.IndexByte(name[lastSlashAt:], '.')
	if dotAt != -1 {
		return name[:lastSlashAt+dotAt], name[lastSlashAt+dotAt+1:]
	}
	return "", ""
}

func findTests(pkg string, gs goroutines) []string {
	internalPkg := strings.TrimSuffix(pkg, "_test")
	externalPkg := internalPkg + "_test"

	matches := make(map[string]struct{})
	for _, g := range gs {
		for _, line := range strings.Split(g, "\n") {
			if strings.HasPrefix(line, "\t") {
				continue
			}
			fpkg, fname := splitPackageAndName(line)
			if (fpkg == internalPkg || fpkg == externalPkg) && strings.HasPrefix(fname, "Test") {
				fname := strings.TrimSuffix(fname, "(...)")
				fname, _, _ = strings.Cut(fname, ".") // trim nested .func1
				if fname != "TestMain" {
					matches[fname] = struct{}{}
				}
			}
		}
	}

	tests := make([]string, 0, len(matches))
	for name := range matches {
		tests = append(tests, name)
	}
	sort.Strings(tests)

	return tests
}
