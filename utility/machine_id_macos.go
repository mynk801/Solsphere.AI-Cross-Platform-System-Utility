// machine_id_macos.go


package main

import (
	"os/exec"
	"regexp"
	"strings"
)

func getMachineID() string {
	cmd := exec.Command("system_profiler", "SPHardwareDataType")
	output, err := cmd.Output()
	if err != nil {
		return "unknown-machine-id"
	}

	re := regexp.MustCompile(`Hardware UUID:\s*(.+)`)
	matches := re.FindStringSubmatch(string(output))

	if len(matches) > 1 {
		return strings.TrimSpace(matches[1])
	}

	return "unknown-machine-id"
}