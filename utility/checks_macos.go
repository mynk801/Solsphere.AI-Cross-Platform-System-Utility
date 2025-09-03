
package main

import (
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

func checkDiskEncryption() bool {
	cmd := exec.Command("fdesetup", "status")
	output, err := cmd.Output()
	if err != nil {
		return false
	}
	return strings.Contains(string(output), "FileVault is On.")
}

func checkOSUpdates() bool {
	cmd := exec.Command("softwareupdate", "-l")
	output, err := cmd.Output()
	if err != nil {
		return false
	}
	return strings.Contains(string(output), "No new software available.")
}

func checkAntivirus() bool {
	avProcesses := []string{
		"Sophos",    
		"falcond",   
		"CbDefense",
	}
	cmd := exec.Command("ps", "aux")
	output, err := cmd.Output()
	if err != nil {
		return false
	}
	runningProcesses := string(output)
	for _, process := range avProcesses {
		if strings.Contains(runningProcesses, process) {
			return true
		}
	}
	return false
}

func checkSleepSettings() bool {
	cmd := exec.Command("pmset", "-g")
	output, err := cmd.Output()
	if err != nil {
		return false
	}
	re := regexp.MustCompile(`\s*displaysleep\s*(\d+)`)
	matches := re.FindStringSubmatch(string(output))
	if len(matches) > 1 {
		minutes, err := strconv.Atoi(matches[1])
		if err != nil {
			return false
		}
		return minutes <= 10
	}
	return false
}