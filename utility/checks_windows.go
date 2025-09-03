// checks_windows.go


package main

import (
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

func checkDiskEncryption() bool {
	cmd := exec.Command("cmd", "/C", "manage-bde", "-status")
	output, err := cmd.Output()
	if err != nil {
		return false
	}
	return strings.Contains(string(output), "Protection On")
}

func checkOSUpdates() bool {
	psCommand := `(New-Object -ComObject "Microsoft.Update.Session").CreateUpdateSearcher().Search("IsInstalled=0").Updates.Count`
	cmd := exec.Command("powershell", "-Command", psCommand)
	output, err := cmd.Output()
	if err != nil {
		return false
	}
	countStr := strings.TrimSpace(string(output))
	count, err := strconv.Atoi(countStr)
	if err != nil {
		return false
	}
	return count == 0
}

func checkAntivirus() bool {
	psCommand := `Get-CimInstance -Namespace root/SecurityCenter2 -ClassName AntiVirusProduct | Select-Object -ExpandProperty productState`
	cmd := exec.Command("powershell", "-Command", psCommand)

	output, err := cmd.Output()
	if err != nil {
		return false
	}

	stateStr := strings.TrimSpace(string(output))
	state, err := strconv.ParseInt(stateStr, 10, 64)
	if err != nil {
		return false
	}
	isEnabled := (state & 4096) != 0

	return isEnabled
}

func checkSleepSettings() bool {
	cmd := exec.Command("powercfg", "/Q")
	output, err := cmd.Output()
	if err != nil {
		return false
	}
	lines := strings.Split(string(output), "\n")
	hexRe := regexp.MustCompile(`(0x[0-9a-fA-F]+)`)
	var inCorrectSection bool = false
	for _, line := range lines {
		line = strings.TrimSpace(line)

		if strings.Contains(line, "3c0bc021-c8a8-4e07-a973-6b14cbcb2b7e") {
			inCorrectSection = true
			continue
		}
		
		if inCorrectSection {
			if strings.HasPrefix(line, "Current AC Power Setting Index:") {
				matches := hexRe.FindStringSubmatch(line)
				if len(matches) >= 2 {
					seconds, err := strconv.ParseInt(strings.TrimPrefix(matches[1], "0x"), 16, 64)
					if err != nil {
						return false 
					}
					return seconds > 0 && seconds <= 600
				}
			}
		}
	}

	return false
}