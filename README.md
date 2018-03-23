# MOTKI CLI

Command `motki` contains interactive command-line tools for managing EVE Online character and corporation assets and industrial processes.

[![GoDoc](https://godoc.org/github.com/motki/cli?status.svg)](https://godoc.org/github.com/motki/cli)


## Getting started

Download the [latest pre-built `motki` binary](https://github.com/motki/cli/releases/latest) for your platform.

> Note that the default configuration connects to `motki.org:18443` using SSL.

#### Install with `go get`

Alternatively, you can install the MOTKI CLI with `go get`.

```bash
go get -u github.com/motki/cli/cmd/motki
```

### Command-line options

```
Usage of motki:
  -credentials string
    	Username and password separated by a colon. (ie. "frank:mypass")
  -history-file string
    	Path to the CLI history file. (default ".history")
  -insecure-skip-verify
    	INSECURE: Skip verification of server SSL cert.
  -log-level string
    	Log level. Possible values: debug, info, warn, error. (default "warn")
  -profile string
    	Enable profiling. Writes profiler data to current directory.
    	Possible values: cpu, mem, mutex, block, trace.
  -server string
    	Backend server host and port. (default "motki.org:18443")
  -version
    	Display the application version.
```

## Authenticating

Some functionality in the application requires authenticating with the remote motkid installation (by default, the [Moritake Industries](https://moritakeindustries.com) website).

To authenticate:

1. Ensure you have a valid account with characters linked on the remote motkid installation.
2. Configure `motki` to use your credentials.
   1. Pass them via command-line option:
   ```
   motki -credentials username:password
   ```
       
   2. Pass them via environment variables:
   ```
   MOTKI_USERNAME=username MOTKI_PASSWORD=password motki
   ```


## Building

The recommended way to build this project is using `make`, though it is compatible with `go get` and friends.

Prerequisites:
1. A semi-recent version of `git`.
2. A valid `go` toolchain install and environment configuration.

### Building with `make`

1. Clone or download and extract the source code and place it in the proper place inside your `$GOPATH`.
   ```bash
   mkdir -p $GOPATH/src/github.com/motki
   git clone https://github.com/motki/cli $GOPATH/src/github.com/motki/cli
   cd $GOPATH/src/github.com/motki/cli
   ```

2. Use the included Makefile to build the application.
   ```bash
   cd $GOPATH/src/github.com/motki/cli
   make clean build
   ```
   
3. Run the newly built executable.
   ```bash
   ./build/motki
   ```

#### Cross-compiling the application

Each build target supports specifying `GOOS` and `GOARCH` to facilitate cross-compiling. For example, building the MOTKI CLI for 32-bit ARM linux:

```bash
make build GOOS="linux" GOARCH="arm"
```

> See the [Go language documentation for more information](https://golang.org/doc/install/source#environment) on supported OSes and architectures.

Build the `motki` program for a combination of OSes and architectures with `make matrix`. Use `OSES` and `ARCHES` to configure which platforms to target.

```bash
make matrix OSES="windows linux darwin" ARCHES="amd64 x86"
```

### Building with `go`

Download and install `motki` and its dependencies using `go get`.

```bash
go get -u github.com/motki/cli/...
```

After `go get` exits successfully, you should have a new command in your `$GOBIN` called `motki`. Along with the binary, the application source code will be located within your `$GOPATH`.

Build and run `motki` by switching to the source directory under your `$GOPATH` and using `go build`.

```bash
cd $GOPATH/src/github.com/motki/cli
go build -o ./motki ./cmd/motki/*.go
./motki -h
```

## Profiling the program

Profile the application by passing the `-profile` flag. For example:

```bash
./build/motki -profile cpu
```

The profiler will write to the current directory in a file named after the profile type. The example above results in `cpu.pprof`.

Once you have exited the process, use `go tool pprof` or `go tool trace` to review the output.

```bash
go tool pprof ./build/motki ./cpu.pprof
```

For more information about profiling Go applications, check out the [blog post](https://blog.golang.org/profiling-go-programs) or [documentation](https://golang.org/pkg/runtime/pprof/).
