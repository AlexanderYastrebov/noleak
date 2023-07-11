package checkmaincleanup

import (
	"os"
	"testing"

	"github.com/AlexanderYastrebov/noleak"
)

func TestMain(m *testing.M) {
	os.Exit(noleak.CheckMainFunc(func() int {
		code := m.Run()

		// perform cleanup
		globalLeak.done()

		return code
	}))
}
