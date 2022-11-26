package noleak_test

import (
	"testing"

	"github.com/AlexanderYastrebov/noleak"
)

func TestMain(m *testing.M) {
	noleak.CheckMain(m)
}
