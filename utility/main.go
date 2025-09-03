// main.go

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"runtime"
	"time"
)

const backendURL = "http://127.0.0.1:5000/api/report"
const checkInterval = 1 * time.Minute

var lastSentPayload ReportPayload

type SystemHealth struct {
	DiskEncrypted   bool `json:"diskEncrypted"`
	OSUpToDate      bool `json:"osUpToDate"`
	AntivirusActive bool `json:"antivirusActive"`
	SleepSettingsOK bool `json:"sleepSettingsOk"`
}

type ReportPayload struct {
	MachineID    string       `json:"machineId"`
	OS           string       `json:"os"`
	CheckResults SystemHealth `json:"checkResults"`
}

func main() {
	fmt.Println("System utility starting... Watching for changes.")
	ticker := time.NewTicker(checkInterval)
	defer ticker.Stop()

	performChecksAndReport()

	for range ticker.C {
		performChecksAndReport()
	}
}

func performChecksAndReport() {
	fmt.Printf("\nRunning checks at %s\n", time.Now().Format(time.RFC1123))

	currentPayload := ReportPayload{
		MachineID: getMachineID(),
		OS:        runtime.GOOS,
		CheckResults: SystemHealth{
			DiskEncrypted:   checkDiskEncryption(),
			OSUpToDate:      checkOSUpdates(),
			AntivirusActive: checkAntivirus(),
			SleepSettingsOK: checkSleepSettings(),
		},
	}

	if !reflect.DeepEqual(currentPayload, lastSentPayload) {
		fmt.Println("Change detected, sending new data to the backend...")
		sendDataToBackend(&currentPayload)
		lastSentPayload = currentPayload
	} else {
		fmt.Println("No changes detected. Skipping report.")
	}
}

func sendDataToBackend(payload *ReportPayload) {
	jsonData, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("Error creating JSON:", err)
		return
	}
	req, err := http.NewRequest("POST", backendURL, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", apiKey)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		fmt.Println("Data sent successfully to the backend.")
	} else {
		fmt.Println("Backend returned an error. Status:", resp.Status)
	}
}