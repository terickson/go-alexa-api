package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

type SimpleActionResponse struct {
	Status  string `json:"status"`
	Intent  string `json:"intent"`
	Message string `json:"message,omitempty"`
}

// handleSimpleAction returns a GET handler for the /action/* endpoints.
// All parameters are read from URL query params — no JSON body needed.
// Designed for easy integration with IFTTT Webhooks.
func handleSimpleAction(room Room) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		token := os.Getenv("GOOGLE_WEBHOOK_TOKEN")
		if token != "" && r.URL.Query().Get("token") != token {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		q := r.URL.Query()
		intent := strings.ToUpper(q.Get("intent"))
		fmt.Println("Simple action intent: " + intent)

		message := ""

		switch intent {
		case "OFF":
			go executeAction(room.TVActionHost, "PowerOff", "")
			if room.ReceiverHost != "" {
				go updateReceiver(room.ReceiverHost, `{"on": false}`)
			}
		case "MUTE":
			if room.ReceiverHost != "" {
				go updateReceiver(room.ReceiverHost, `{"mute": true}`)
			} else {
				go executeAction(room.TVActionHost, "Mute", "")
			}
		case "UNMUTE":
			if room.ReceiverHost != "" {
				go updateReceiver(room.ReceiverHost, `{"mute": false}`)
			}
		case "VOLUME":
			level := strings.TrimSpace(q.Get("level"))
			if level == "" {
				log.Println("VOLUME: missing level param")
				message = "I'm sorry I could not process your request " + intent + "."
			} else {
				if room.ReceiverHost != "" {
					go updateReceiver(room.ReceiverHost, `{"volume": -`+level+`}`)
				} else {
					go executeAction(room.TVActionHost, "Volume", level)
				}
			}
		case "CHANNEL":
			number := strings.TrimSpace(q.Get("number"))
			if number == "" {
				log.Println("CHANNEL: missing number param")
				message = "I'm sorry I could not process your request " + intent + "."
			} else {
				go executeAction(room.TVActionHost, "Channel", number)
			}
		case "CHANNELUP":
			go executeAction(room.TVActionHost, "ChannelUp", "")
		case "CHANNELDOWN":
			go executeAction(room.TVActionHost, "ChannelDown", "")
		case "INPUT":
			input := strings.TrimSpace(q.Get("input"))
			if input == "" {
				log.Println("INPUT: missing input param")
				message = "I'm sorry I could not process your request " + intent + "."
			} else {
				inputType := strings.ReplaceAll(strings.ToUpper(input), " ", "")
				setInput(room, inputType)
			}
		case "HOME":
			go executeAction(room.RokuActionHost, "home", "")
		case "BACK":
			go executeAction(room.RokuActionHost, "back", "")
		case "UP", "DOWN", "LEFT", "RIGHT":
			spaces := strings.TrimSpace(q.Get("spaces"))
			if spaces == "" {
				spaces = "1"
			}
			go executeAction(room.RokuActionHost, strings.ToLower(intent), spaces)
		case "ENTER":
			go executeAction(room.RokuActionHost, "enter", "")
		case "SELECT":
			go executeAction(room.RokuActionHost, "select", "")
		case "PLAY":
			go executeAction(room.RokuActionHost, "right", "")
		case "FORWARD":
			go executeAction(room.RokuActionHost, "forward", "")
		case "REVERSE":
			go executeAction(room.RokuActionHost, "reverse", "")
		case "SEARCH":
			search := strings.TrimSpace(q.Get("search"))
			if search == "" {
				log.Println("SEARCH: missing search param")
				message = "I'm sorry I could not process your request " + intent + "."
			} else {
				go executeAction(room.RokuActionHost, "search", search)
			}
		default:
			message = "I'm sorry I could not process your request " + intent + "."
		}

		status := "ok"
		if message != "" {
			status = "error"
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(SimpleActionResponse{
			Status:  status,
			Intent:  intent,
			Message: message,
		})
	}
}
