package main

import (
	"archive/tar"
	"fmt"
	"github.com/sagernet/cronet-go/utils"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/klauspost/compress/gzip"
)

var logger = log.New(os.Stderr, "prebuild", log.LstdFlags)

const (
	clangVersion = "llvmorg-15-init-11722-g3f3a235a-2"
)

var (
	goos   string
	goarch string
)

func init() {
	goos = os.Getenv("GOOS")
	if goos == "" {
		goos = runtime.GOOS
	}
	goarch = os.Getenv("GOARCH")
	if goarch == "" {
		goarch = runtime.GOARCH
	}
}

func main() {
	if !utils.FileExists("build/llvm/bin/clang") {
		os.RemoveAll("build/llvm")
		os.MkdirAll("build/llvm", 0o755)
		clangDownload := os.ExpandEnv("https://commondatastorage.googleapis.com/chromium-browser-clang/" + clangOsString() + "/clang-" + clangVersion + ".tgz")
		logger.Println(">> ", clangDownload)
		clangDownloadResponse, err := http.Get(clangDownload)
		if err != nil {
			logger.Fatal(err)
		}
		gzReader, err := gzip.NewReader(clangDownloadResponse.Body)
		if err != nil {
			logger.Fatal(err)
		}
		tarReader := tar.NewReader(gzReader)
		linkName := make(map[string]string)
		for {
			header, err := tarReader.Next()
			if err != nil {
				if err == io.EOF {
					break
				}
				logger.Fatal(err)
			}
			path := filepath.Join("build/llvm", header.Name)
			if header.FileInfo().IsDir() {
				continue
			}
			logger.Println(">> ", path)
			if header.Linkname != "" {
				linkName[path] = filepath.Join(filepath.Dir(path), header.Linkname)
				linkName[path], _ = filepath.Abs(linkName[path])
				continue
			}
			err = os.MkdirAll(filepath.Dir(path), 0o755)
			if err != nil {
				return
			}
			file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.FileMode(header.Mode))
			if err != nil {
				logger.Fatal(err)
			}

			_, err = io.CopyN(file, tarReader, header.Size)
			if err != nil {
				logger.Fatal(err)
			}
			file.Close()
		}
		clangDownloadResponse.Body.Close()
		var notExists, leftNotExists int
		for {
			for dst, src := range linkName {
				if !utils.FileExists(src) {
					notExists++
					continue
				}
				logger.Println(">> ", src, " => ", dst)
				os.MkdirAll(filepath.Dir(dst), 0o755)
				err = os.Symlink(src, dst)
				if err != nil {
					logger.Fatal(err)
				}
				delete(linkName, dst)
			}
			if notExists == 0 {
				break
			}
			if notExists == leftNotExists {
				logger.Fatal("untar: link failed")
			}
			leftNotExists = notExists
			notExists = 0
		}
	}

	output := filepath.Join("build", goos, goarch)
	p := filepath.Join(output, "libcronet.so")
	if !utils.FileExists(p) {
		logger.Fatal(fmt.Sprintf("libcronet.so not found in '%s'", p))
	}
}

func clangOsString() string {
	clangOs := strings.ToUpper(runtime.GOOS[:1]) + runtime.GOOS[1:]
	clangArch := runtime.GOARCH
	switch clangArch {
	case "amd64":
		clangArch = "x64"
	case "386":
		clangArch = "x86"
	case "mipsle":
		clangArch = "mipsel"
	case "mips64le":
		clangArch = "mips64el"
	}
	return clangOs + "_" + clangArch
}
