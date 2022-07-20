package utils

import (
	"os/exec"
)

func ExecCommand(command string, params ...string) (string, error) {
	cmd := exec.Command(command, params...)
	var output string
	out, err := cmd.Output()
	if err != nil {
		return output, err
	}

	return string(out), nil
}
