package httpclient

import (
	"net/http"
	"testing"
)

func TestDefaultHttpClient(t *testing.T) {
	defer http.DefaultClient.CloseIdleConnections()

	rsp, err := http.Get("https://pkg.go.dev/github.com/AlexanderYastrebov/noleak")
	if err != nil {
		t.Fatal(err)
	}
	rsp.Body.Close()
}

func TestHttpClient(t *testing.T) {
	client := &http.Client{}
	defer client.CloseIdleConnections()

	rsp, err := client.Get("https://pkg.go.dev/github.com/AlexanderYastrebov/noleak")
	if err != nil {
		t.Fatal(err)
	}
	rsp.Body.Close()
}
