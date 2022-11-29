# Simple goroutine leak detector

It compares the list of active goroutines before and after the test and raises an error on mismatch.

See and run [example_test.go](example_test.go):

```sh
GODEBUG="tracebackancestors=1" go test . -tags=example
```

## Credits

* net/http/main_test.go
* https://github.com/uber-go/goleak
* https://github.com/fortytw2/leaktest
