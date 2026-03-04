package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

func executeAction(host string, command string, value string) {
	log.Println("host:", host)
	bodyStr := `{"command": "` + command + `"}`
	if len(value) > 0 {
		bodyStr = `{"command": "` + command + `", "value": "` + value + `"}`
	}
	log.Println("body: ", bodyStr)

	req, err := http.NewRequest("POST", host, bytes.NewBuffer([]byte(bodyStr)))
	if err != nil {
		log.Println(err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()

	log.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := io.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
}

// updateReceiver sends a PUT request to update receiver state.
// Body string should look like: {"on": true, "volume": "string", "input": "string", "mute": true}
// but should only include the properties that need to be updated.
func updateReceiver(host string, bodyStr string) {
	log.Println("host:", host)
	log.Println("body: ", bodyStr)

	req, err := http.NewRequest(http.MethodPut, host, bytes.NewBuffer([]byte(bodyStr)))
	if err != nil {
		log.Println(err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()

	log.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := io.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
}
