# Simple goroutine leak detector

It compares the list of active goroutines before and after the test and reports an error on mismatch.

See and run [example_test.go](example_test.go):

```sh
GODEBUG=tracebackancestors=1 go test . -tags=example
```

Setting `GODEBUG=tracebackancestors=N` extends tracebacks with the stacks at
which goroutines were created, where N limits the number of ancestor goroutines to
report, see https://pkg.go.dev/runtime.

Tests using http.Client or http.DefaultClient should close idle connections, see [httpclient_test.go](httpclient_test.go).

## Credits

* net/http/main_test.go
* https://github.com/uber-go/goleak
* https://github.com/fortytw2/leaktest
