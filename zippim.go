package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/gernest/kemi"
	"github.com/vivainio/walker"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	app     = kingpin.New("zippim", "An unzip-based package manager for Windows")
	getCmd  = app.Command("get", "Get a package from url")
	getName = getCmd.Flag("name", "Custom name for the package").String()
	getURL  = getCmd.Arg("url", "Url to retrieve (can be zip/exe)").Required().String()
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

// not used, bad idea
func mklink(link string, tgt string) {
	fmt.Printf("link %s -> %s\n", link, tgt)
	c := exec.Command("cmd", "/C", "mklink", link, tgt)
	err := c.Run()
	if err != nil {
		log.Fatal(err)
	}
}

func writeString(fname string, cont string) {
	f, _ := os.Create(fname)
	f.WriteString(cont)
	f.Close()
}

func makeLauncher(launcherDir string, tgt string) {
	relpath, _ := filepath.Rel(launcherDir, tgt)
	cont := "@%~dp0/" + relpath + " %*"
	basename := filepath.Base(tgt)
	name := strings.TrimSuffix(basename, filepath.Ext(basename))
	launcherName := name + ".cmd"
	writeString(filepath.Join(launcherDir, launcherName), cont)
}

func makeLaunchers(startPath string, linkPath string) {
	visitor := func(pth string, fileinfos []os.FileInfo) bool {
		fmt.Printf("Link: %s\n", pth)
		for _, f := range fileinfos {
			name := f.Name()
			fmt.Printf("  %s\n", name)
			if matched, _ := regexp.MatchString(".*exe$", name); matched {
				tgtpath := filepath.Join(pth, name)
				makeLauncher(linkPath, tgtpath)
			}
		}
		return true
	}

	walker.WalkOne(startPath, visitor)
}

func getPackageNameFromFilename(fname string) string {
	parts := strings.Split(fname, ".")
	archivename := parts[0]
	return archivename
}

func install(pkg string, pkgname string) {
	fmt.Printf("Installing %s", pkg)
	resp, _ := http.Get(pkg)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	fname := fileNameFromURL(pkg)
	dlpath := filepath.Join(getRelDir("_download"), fname)
	ioutil.WriteFile(dlpath, body, 0644)
	tgt := getRelDir(pkgname)
	err := kemi.Unpack(dlpath, tgt)
	if err != nil {
		panic("Failed to unpack")

	}

	bindir := getRelDir("bin")
	makeLaunchers(tgt, bindir)
}

func main() {
	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case getCmd.FullCommand():
		pkgName := ""
		if len(*getName) == 0 {
			pkgName = getPackageNameFromFilename(fileNameFromURL(*getURL))
		} else {
			pkgName = *getName
		}

		install(*getURL, pkgName)
	}
}
