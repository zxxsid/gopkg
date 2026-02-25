## Cursor Cloud specific instructions

This is a pure Go library (`github.com/zxxsid/gopkg`) with no external dependencies, no servers, and no databases. The only requirement is **Go >= 1.25.7**.

### Build / Lint / Test

```
go build ./...
go vet ./...
go test ./...
```

There are currently no test files in the repository (`image/` package has no `*_test.go` files), so `go test` reports `[no test files]`.

### Notes

- No Docker, Makefile, or setup scripts exist; the Go toolchain is the only dependency.
- The library uses only Go standard library packages (`image`, `io`); there is no `go.sum` file because there are zero third-party dependencies.
- To test the library interactively, create a small Go program that imports `github.com/zxxsid/gopkg/image` with a `replace` directive pointing to `/workspace`.
