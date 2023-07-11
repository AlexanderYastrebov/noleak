# Simple goroutine leak detector

It compares the list of active goroutines before and after the test and reports an error on mismatch.

To enable detector add a [main_test.go](examples/checkmain/main_test.go) file to the package.

To check a specific test only add `noleak.Check(t)` [at the top](examples/check/example_test.go).

Tests using http.Client or http.DefaultClient [should close idle connections](examples/httpclient/example_test.go).

See and run all [examples](examples):

```sh
GODEBUG=tracebackancestors=1 go test ./examples/...
```

Setting `GODEBUG=tracebackancestors=N` extends tracebacks with the stacks at
which goroutines were created, where N limits the number of ancestor goroutines to
report, see https://pkg.go.dev/runtime.


## Credits

* net/http/main_test.go
* https://github.com/uber-go/goleak
* https://github.com/fortytw2/leaktest
