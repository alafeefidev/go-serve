package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	run()
}

const (
	TmpDir     = "temp"
	ScriptName = "ALPHA_SWEET_GOLD.js"
)

const Script = `
alert("%s!");
`

func run() {
	if len(os.Args) < 2 {
		nerr("No file provided")
	}
	farg := os.Args[1]

	if strings.HasSuffix(farg, "/") || strings.HasSuffix(farg, "\\") {
		nerrf("%s is a folder\n", farg)
	}
	// go won't allow it but why not
	if !strings.HasSuffix(farg, ".go") {
		farg = farg + ".go"
	}

	if _, err := os.Stat(farg); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			nerrf("File does not exist: %s\n", farg)
		} else {
			nerrf("error checking file")
		}
	}

}

func nerr(a ...any) {
	fmt.Fprintln(os.Stderr, a...)
	os.Exit(1)
}

func nerrf(format string, a ...any) {
	fmt.Fprintf(os.Stderr, format, a...)
	os.Exit(1)
}

func BuildBinary(file, tempFolder string) (string, error) {
	if err := os.MkdirAll(tempFolder, 0677); err != nil {
		return "", fmt.Errorf("error creating temp folder: %w", err)
	}

	fileName := filepath.Base(file)
	fileName = strings.Split(fileName, ".")[0]
	tempfile := filepath.Join(tempFolder, fileName)

	cmd := exec.Command("go", "build", "-o", tempfile, file)

	return tempfile, nil
}

type ScriptFile struct {
	name string
	dir  string
}

func NewScriptFile(name, dir string) *ScriptFile {
	return &ScriptFile{name, dir}
}

func (s *ScriptFile) Create(c []byte) error {
	path := filepath.Join(s.dir, s.name)
	if _, err := os.Stat(path); err != nil && !errors.Is(err, os.ErrNotExist) {
		if err := s.Remove(); err != nil {
			return fmt.Errorf("error removing old script: %w", err)
		}
	}
	sc, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed creating script file: %w", err)
	}
	defer sc.Close()

	if _, err := sc.Write(c); err != nil {
		return fmt.Errorf("error writing script to file: %w", err)
	}
	return nil
}

// Does not return an error if file doesn't exists
func (s *ScriptFile) Remove() error {
	path := filepath.Join(s.dir, s.name)
	if err := os.Remove(path); err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}
	return nil
}
