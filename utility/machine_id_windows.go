


package main

import (
	"os/exec"
	"strings"
)

func getMachineID() string {
	cmd := exec.Command("wmic", "baseboard", "get", "serialnumber")
	output, err := cmd.Output()
	if err != nil {
		return "unknown-machine-id"
	}
	lines := strings.Split(string(output), "\n")
	if len(lines) > 1 {
		return strings.TrimSpace(lines[1])
	}
	return "unknown-machine-id"
}