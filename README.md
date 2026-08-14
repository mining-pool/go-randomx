# go-randomx

Go cgo bindings for [RandomX](https://github.com/tevador/RandomX) **v1.2.1** —
light-mode VMs for pool-side share verification, `randomx_get_flags()` exposure,
and a batch `Search`.

## Install

```go
import randomx "github.com/mining-pool/go-randomx"
```

```bash
go get github.com/mining-pool/go-randomx
CGO_ENABLED=1 go build   # links the prebuilt static lib for your platform
```

No `build.sh`, no system install, no `replace` directive. The RandomX static lib
is **prebuilt per platform** and committed under `lib/librandomx_<goos>_<goarch>.a`;
`randomx.go` links the right one via `#cgo <goos>,<goarch> LDFLAGS` selectors, so
the library ships inside the module and `go build` just works.

Supported: `linux/amd64`, `linux/arm64`, `darwin/amd64`, `darwin/arm64`. (A C++
toolchain is still required to compile the small cgo glue and link.)

## Regenerating the static libs

The prebuilt libs are produced by
[`.github/workflows/build-libs.yml`](.github/workflows/build-libs.yml): it builds
RandomX v1.2.1 with cmake (portable, no `-march=native`) on a native runner for
each platform and commits the results. Run it from the Actions tab
(**workflow_dispatch**) after bumping the RandomX version. To add a platform, add
a matrix entry and a matching `#cgo` selector in `randomx.go`.

## License

RandomX is BSD-licensed (tevador/RandomX). See [LICENSE](LICENSE).
