package noleak_test

import (
	"net/http"
	"testing"

	"github.com/AlexanderYastrebov/noleak"
)

func TestDefaultHttpClient(t *testing.T) {
	noleak.Check(t)

	defer http.DefaultClient.CloseIdleConnections()

	rsp, err := http.Get("https://pkg.go.dev/github.com/AlexanderYastrebov/noleak")
	if err != nil {
		t.Fatal(err)
	}
	rsp.Body.Close()
}

func TestHttpClient(t *testing.T) {
	noleak.Check(t)

	client := &http.Client{}
	defer client.CloseIdleConnections()

	rsp, err := client.Get("https://pkg.go.dev/github.com/AlexanderYastrebov/noleak")
	if err != nil {
		t.Fatal(err)
	}
	rsp.Body.Close()
}
