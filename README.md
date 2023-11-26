# cronet-go

Cronet is the Chromium network stack made available as a library. Cronet takes advantage of multiple
technologies that reduce the latency and increase the throughput of the network requests.

The Cronet Library handles the requests of apps used by millions of people on a daily basis, such as YouTube, Google
App, Google Photos, and Maps - Navigation & Transit.

## Build
### Cronet binaries

#### Precompiled
See [sleeyax/cronet-binaries](https://github.com/sleeyax/cronet-binaries/releases) for prebuilt cronet binaries.

#### Manual compilation

1. Build cronet locally by following the [instructions for Desktop builds](https://chromium.googlesource.com/chromium/src/+/master/components/cronet/build_instructions.md#desktop-builds-targets-the-current-os) from Google.
2. Copy all shared libraries from `~/chromium/src/out/Cronet/*.so` to `./build/<os>/<arch>`.
3. Rename `libcronet.x.x.x.x.so` to `libcronet.so`.
4. (optional) If you have a `sysroot` folder (e.g. for cross-compilation), copy it to `./build/<os>/<arch>/sysroot`.

### Package
Before you begin, make sure you've followed the instructions above and have the `prebuild` and `build` utilities compiled:

```bash
$ go build -v -o prebuild ./cmd/prebuild
$ go build -v -o gobuild ./cmd/build
```

Then, download the required toolchain for your system:
    
```bash
$ ./prebuild
```

You only need to do these steps once per project.

Finally, build the `example` program with the `build` utility:

```bash
# dynamic linking
$  ./gobuild -v -o example ./cmd/exmample/main.go
# static linking (if enabled)
$ ./gobuild -v -o example -tags cronet_static -trimpath -ldflags "-s -w -buildid=" ./example
$ ./build/llvm/bin/llvm-strip ./example
```

By default, `gobuild` wraps the `go build` command. To run `go test` instead, specify `test` as the first argument:

```bash
$ ./gobuild test ./...
```

See the [GitHub workflow](./.github/workflows/debug.yml) for more details.
