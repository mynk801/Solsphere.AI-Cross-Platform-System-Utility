// checks_linux.go


package main

import (
	"os"
	"os/exec"
	"strings"
)

func checkDiskEncryption() bool {
	cmd := exec.Command("lsblk", "-o", "TYPE")
	output, err := cmd.Output()
	if err != nil {
		return false
	}
	return strings.Contains(string(output), "crypt")
}

func checkOSUpdates() bool {
	if _, err := os.Stat("/var/run/reboot-required"); err == nil {
		return false 
	}
	return true 
}

func checkAntivirus() bool {
	cmd := exec.Command("ps", "-e")
	output, err := cmd.Output()
	if err != nil {
		return false
	}
	return strings.Contains(string(output), "clamd")
}

func checkSleepSettings() bool {
	return true
}