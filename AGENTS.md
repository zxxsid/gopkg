## Cursor Cloud specific instructions

This is a Go image processing library (`github.com/zxxsid/gopkg`). The only system requirement is **Go >= 1.25.7**. The sole external dependency is `golang.org/x/image`.

### Build / Lint / Test

```
go build ./...
go vet ./...
go test -v ./image/...
go test -bench=. -benchmem ./image/...
```

### Notes

- No Docker, Makefile, or setup scripts exist; the Go toolchain is the only dependency.
- To test the library interactively, create a small Go program that imports `github.com/zxxsid/gopkg/image` with a `replace` directive pointing to `/workspace`.
