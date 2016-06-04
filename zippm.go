package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/gernest/kemi"
)

func fileNameFromURL(urls string) string {
	u, _ := url.Parse(urls)
	path := u.Path
	return filepath.Base(path)
}

func getRelDir(suffix string) string {
	joined := filepath.Join(root(), suffix)
	if _, err := os.Stat(joined); os.IsNotExist(err) {
		os.MkdirAll(joined, 0777)

	}
	return joined
}

func root() string {
	e := os.Getenv("ZIPPM_ROOT")
	if len(e) > 0 {
		return e
	}
	return "C:\\pkg"
}

func install(pkg string) {
	fmt.Printf("Installing %s", pkg)
	resp, _ := http.Get(pkg)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	fname := fileNameFromURL(pkg)
	dlpath := filepath.Join(getRelDir("_download"), fname)
	ioutil.WriteFile(dlpath, body, 0644)
	parts := strings.Split(fname, ".")
	archivename := parts[0]
	tgt := getRelDir(archivename)
	err := kemi.Unpack(dlpath, tgt)
	if err != nil {
		panic("Failed to unpack")

	}
}

func main() {
	cmd := os.Args[1]
	pkg := os.Args[2]
	install(pkg)

	fmt.Printf("hello world %s", cmd)
}
