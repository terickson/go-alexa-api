package main

import (
	"fmt"
	"log"
	"time"
)

// InputConfig defines the device actions for an input type.
type InputConfig struct {
	ReceiverInput string // receiver input to set (e.g., "AV1", "HDMI4"); only used if room has a receiver
	TVInput       string // TV command to switch input (e.g., "HDMI1", "InputTV")
	RokuApp       string // if non-empty, launch this Roku app
}

// addAliases maps multiple aliases to the same InputConfig.
func addAliases(m map[string]InputConfig, cfg InputConfig, aliases ...string) {
	for _, alias := range aliases {
		m[alias] = cfg
	}
}

// setInput switches the input for a room based on the input type.
func setInput(room Room, inputType string) {
	log.Println("Input passed: ", inputType)

	cfg, ok := room.InputMap[inputType]
	if !ok {
		cfg = InputConfig{
			ReceiverInput: "HDMI1",
			TVInput:       "HDMI1",
			RokuApp:       inputType,
		}
	}

	if room.ReceiverHost != "" {
		receiverInput := cfg.ReceiverInput
		if receiverInput == "" {
			receiverInput = "HDMI1"
		}
		payload := fmt.Sprintf(`{"on":true, "volume": %d, "input": "%s"}`, room.DefaultVolume, receiverInput)
		go updateReceiver(room.ReceiverHost, payload)
	}

	if cfg.RokuApp != "" {
		go executeAction(room.RokuActionHost, "input", cfg.RokuApp)
	}

	go executeAction(room.TVActionHost, "PowerOn", "")
	time.Sleep(500 * time.Millisecond)
	go executeAction(room.TVActionHost, cfg.TVInput, "")
}
