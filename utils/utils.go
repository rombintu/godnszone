package utils

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

// Execute any commands from shell
func ExecCommand(command string, params ...string) (string, error) {
	cmd := exec.Command(command, params...)
	var output string
	out, err := cmd.Output()
	if err != nil {
		return output, err
	}

	return string(out), nil
}

// Precheck file path
func ToValidPath(path string) string {
	if filepath.IsAbs(path) {
		return path
	} else {
		workDir, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}
		return filepath.Join(workDir, path)
	}
}
