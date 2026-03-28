package cli

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	wsa "github.com/alafeefidev/go-serve/internal/ws"
)

// TODO: change
// TODO: change the js name to be auto without naming it, will see
const (
	TEST_JS_SCRIPT = `new WebSocket("%s").onmessage=e=>e.data==="refresh"&&location.reload();`
	TEST_JS_NAME   = "BIG_GLOB_GOLD.js"
)

func runCMD(args []string) error {
	ServeFile(args[0])
	return nil
}

func ServeFile(file string) error {
	if strings.HasSuffix(file, "/") || strings.HasSuffix(file, "\\") {
		return fmt.Errorf("%s is a directory not a file", file)
	}

	// Go build/run doesn't allow supplying multiple
	//go files in different directories, so user can
	//supply without extension
	if !strings.HasSuffix(file, ".go") {
		file = file + ".go"
	}

	if _, err := os.Stat(file); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("file does not exist: %s\n", file)
		} else {
			return fmt.Errorf("error checking file: %w", err)
		}
	}

	tempDir, err := os.MkdirTemp("", "go-serve-*")
	if err != nil {
		return fmt.Errorf("error making temp dir: %w", err)
	}
	defer os.RemoveAll(tempDir)

	bin, err := buildBinary(file, tempDir)
	if err != nil {
		return err
	}

	if err := os.Chdir(tempDir); err != nil {
		return fmt.Errorf("error changing current directory to temp: %w", err)
	}

	ws := wsa.NewWS()
	go ws.Run()

	wsAddr, err := wsa.StartServerWS(ws)
	if err != nil {
		return err
	}

	//TODO: change name and ws code insertion
	script := fmt.Sprintf(TEST_JS_SCRIPT, wsAddr)
	if err := createScript(TEST_JS_NAME, script, tempDir); err != nil {
		return err
	}

	go func(bin string) {
		if err := runBinary(bin); err != nil {
			panic(err)
		}
	}(bin)

	for {
		<-ws.Connected()
		slog.Info("client connected")

		ticker := time.NewTicker(5 * time.Second)
	loop:
		for {
			select {
			case <-ticker.C:
				ws.Send([]byte("refresh"))
			case <-ws.Disconnected():
				slog.Info("client disconnected")
				ticker.Stop()
				break loop
			}
		}
	}
}

// Takes the go file to build and the temp directory where to build it
//
// Returns the path for the executable with the extension if needed
func buildBinary(file, tempDir string) (string, error) {
	fileNameNoExt := strings.TrimSuffix(filepath.Base(file), ".go")

	var ext string
	if runtime.GOOS == "windows" {
		ext = ".exe"
	}

	resultFile := filepath.Join(tempDir, fileNameNoExt+ext)
	cmd := exec.Command("go", "build", "-o", tempDir, file)
	if o, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed at building go binary: %w\n%s", err, o)
	}
	return resultFile, nil
}

// Takes the script name.ext and content and the temp directory and adds it there
func createScript(name, content, tempDir string) error {
	scriptPath := filepath.Join(tempDir, name)
	sc, err := os.Create(scriptPath)
	if err != nil {
		return fmt.Errorf("error creating script file: %w", err)
	}
	defer sc.Close()

	if _, err := sc.WriteString(content); err != nil {
		return fmt.Errorf("error writing to script file: %w", err)
	}
	return nil
}

func runBinary(bin string) error {
	cmd := exec.Command(bin)
	fmt.Println(cmd)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdin
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error running binary: %w", err)
	}
	return nil
}
