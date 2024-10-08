package main

import (
	"github.com/sagernet/cronet-go/utils"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"syscall"

	"github.com/sagernet/sing/common"
)

var cflags map[string]string
var ldflags map[string]string

func cgoCFLAGS(goos string, goarch string, flags string) {
	if cflags == nil {
		cflags = make(map[string]string)
	}
	cflags[goos+goarch] = flags
}

func cgoLDFLAGS(goos string, goarch string, flags string) {
	if ldflags == nil {
		ldflags = make(map[string]string)
	}
	ldflags[goos+goarch] = flags
}

var logger = log.New(os.Stderr, "build", log.LstdFlags)

func appendEnv(key string, value string) {
	common.Must(os.Setenv(key, strings.TrimSpace(os.ExpandEnv("$"+key+" "+value))))
}

func main() {
	var args []string

	userArgs := os.Args[1:]

	if len(userArgs) > 0 && userArgs[0] == "test" {
		args = append(args, "test")
		userArgs = userArgs[1:]
	} else {
		args = append(args, "build")
	}
	args = append(args, "-gcflags", "-c "+strconv.Itoa(runtime.NumCPU()))

	goos := os.Getenv("GOOS")
	if goos == "" {
		goos = runtime.GOOS
	}
	goarch := os.Getenv("GOARCH")
	if goarch == "" {
		goarch = runtime.GOARCH
	}
	//openwrt := os.Getenv("OPENWRT")

	var disablePie bool
	if goarch == "mipsle" || goarch == "mips64le" {
		disablePie = true
	} else if goos == "windows" && goarch == "arm64" {
		disablePie = true
	} else if goarch == "186" && goos != "android" {
		disablePie = true
	}

	if !disablePie {
		args = append(args, "-buildmode", "pie")
		appendEnv("CGO_LDFLAGS", "-pie")
	} else {
		appendEnv("CGO_LDFLAGS", "-nopie")
	}

	switch goos {
	case "windows":
		os.Setenv("MSYS", "winsymlinks:nativestrict")
	case "linux":
		sysRoot := os.ExpandEnv("$PWD/build/" + goos + "/" + goarch + "/sysroot")
		if utils.FileExists(sysRoot) {
			appendEnv("CGO_CFLAGS", "--sysroot="+sysRoot)
			appendEnv("CGO_LDFLAGS", "--sysroot="+sysRoot)
		}
		appendEnv("CGO_CFLAGS", cflags[goos+goarch])
		appendEnv("CGO_LDFLAGS", ldflags[goos+goarch])
	}

	os.Setenv("PATH", os.ExpandEnv("$PWD/build/llvm/bin:$PATH"))
	os.Setenv("CC", "clang")
	os.Setenv("CGO_ENABLED", "1")
	os.Setenv("CGO_LDFLAGS_ALLOW", ".*")

	args = append(args, userArgs...)

	err := execve("go", args...)
	if err != nil {
		logger.Fatal(err)
	}
}

func execve(name string, args ...string) error {
	path, _ := exec.LookPath(name)
	args = append([]string{path}, args...)
	env := os.Environ()
	return syscall.Exec(path, args, env)
}
