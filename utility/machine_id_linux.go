// machine_id_linux.go


package main

import (
	"os"
	"strings"
)

func getMachineID() string {
	idBytes, err := os.ReadFile("/sys/class/dmi/id/product_uuid")
	if err != nil {
		return "unknown-machine-id"
	}
	return strings.TrimSpace(string(idBytes))
}