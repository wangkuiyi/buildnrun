package buildnrun

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"strings"
)

func build(pkg string) error {
	goCompiler, e := exec.LookPath("go")
	if e != nil {
		return fmt.Errorf("Cannot find go command: %v", e)
	}

	if e := exec.Command(goCompiler, "install", pkg).Run(); e != nil {
		return fmt.Errorf("Failed building %s: %v", pkg, e)
	}

	return nil
}

func run(pkg string, args ...string) (string, string, error) {
	c := exec.Command(path.Join(gopath(), "bin", path.Base(pkg)), args...)
	var out, err bytes.Buffer
	op, _ := c.StdoutPipe()
	ep, _ := c.StderrPipe()
	go func() { io.Copy(&out, op) }()
	go func() { io.Copy(&err, ep) }()
	if e := c.Run(); e != nil {
		return "", "", fmt.Errorf("%s failed: %v", path.Base(pkg), e)
	}

	return out.String(), err.String(), nil
}

func Run(pkg string, args ...string) (string, string, error) {
	if e := build(pkg); e != nil {
		return "", "", e
	}
	out, err, e := run(pkg, args...)
	return out, err, e
}

// PkgDir returns the package directory prefixed by $GOPATH.  In case of error,
// it returns "".
func Pkg(pkg string) string {
	return path.Join(gopath(), "src", pkg)
}

func gopath() string {
	return strings.Split(os.Getenv("GOPATH"), ":")[0]
}
