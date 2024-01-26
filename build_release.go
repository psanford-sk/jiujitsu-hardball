//go:build ignore
// +build ignore

package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	flag.Parse()
	os.MkdirAll("release", 0777)

	for _, t := range targets {
		cmd := exec.Command("go", "build", "-trimpath", "-o", filepath.Join("release", t.binaryName()))
		env := []string{"GOOS=" + t.goos, "GOARCH=" + t.garch, "GO111MODULE=on"}
		if t.goarm != "" {
			env = append(env, "GOARM="+t.goarm)
		}
		cmd.Env = append(os.Environ(), env...)

		fmt.Printf("run: %s %s %s\n", strings.Join(env, " "), cmd.Path, strings.Join(cmd.Args[1:], " "))

		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s %s err: %s, out: %s\n", t.goos, t.garch, err, out)
			os.Exit(1)
		}
	}
}

type target struct {
	goos  string
	garch string
	goarm string
}

func (t *target) binaryName() string {
	ext := ""
	if t.goos == "windows" {
		ext = ".exe"
	}

	tmpl := "jiujitsu-hardball-%s-%s%s%s"
	return fmt.Sprintf(tmpl, t.goos, t.garch, t.goarm, ext)
}

var targets = []target{
	{"linux", "amd64", ""},
	{"linux", "arm64", ""},
	{"darwin", "amd64", ""},
	{"darwin", "arm64", ""},
	{"windows", "386", ""},
}
