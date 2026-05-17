package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type DialogflowRequest struct {
	ResponseID  string      `json:"responseId"`
	QueryResult QueryResult `json:"queryResult"`
}

type QueryResult struct {
	Intent     DialogflowIntent       `json:"intent"`
	Parameters map[string]interface{} `json:"parameters"`
}

type DialogflowIntent struct {
	DisplayName string `json:"displayName"`
}

type DialogflowResponse struct {
	FulfillmentText string `json:"fulfillmentText"`
}

// getParam retrieves a parameter from a Dialogflow parameters map, normalizing
// @sys.number (float64) and @sys.any (string) types to plain strings.
func getParam(params map[string]interface{}, key string) (string, bool) {
	val, ok := params[key]
	if !ok || val == nil {
		return "", false
	}
	switch v := val.(type) {
	case string:
		s := strings.TrimSpace(v)
		return s, s != ""
	case float64:
		return strconv.Itoa(int(v)), true
	default:
		s := fmt.Sprintf("%v", v)
		return s, s != ""
	}
}

// handleGoogleIntent returns a Dialogflow webhook handler for the given room.
func handleGoogleIntent(room Room) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		token := os.Getenv("GOOGLE_WEBHOOK_TOKEN")
		if token != "" && r.URL.Query().Get("token") != token {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		var dfReq DialogflowRequest
		if err := json.NewDecoder(r.Body).Decode(&dfReq); err != nil {
			log.Println("failed to decode dialogflow request:", err)
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		intent := strings.ToUpper(dfReq.QueryResult.Intent.DisplayName)
		params := dfReq.QueryResult.Parameters
		fmt.Println("Google intent passed: " + intent)
		output := "Processing Request."

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
			level, ok := getParam(params, "Level")
			if !ok {
				output = "I'm sorry I could not process your request " + intent + "."
			} else {
				if room.ReceiverHost != "" {
					go updateReceiver(room.ReceiverHost, `{"volume": -`+level+`}`)
				} else {
					go executeAction(room.TVActionHost, "Volume", level)
				}
			}
		case "CHANNEL":
			number, ok := getParam(params, "Number")
			if !ok {
				output = "I'm sorry I could not process your request " + intent + "."
			} else {
				go executeAction(room.TVActionHost, "Channel", number)
			}
		case "CHANNELUP":
			go executeAction(room.TVActionHost, "ChannelUp", "")
		case "CHANNELDOWN":
			go executeAction(room.TVActionHost, "ChannelDown", "")
		case "INPUT":
			inputType, ok := getParam(params, "InputType")
			if !ok {
				output = "I'm sorry I could not process your request " + intent + "."
			} else {
				inputType = strings.ReplaceAll(strings.ToUpper(inputType), " ", "")
				setInput(room, inputType)
			}
		case "HOME":
			go executeAction(room.RokuActionHost, "home", "")
		case "BACK":
			go executeAction(room.RokuActionHost, "back", "")
		case "UP", "DOWN", "LEFT", "RIGHT":
			spaces := "1"
			if s, ok := getParam(params, "Spaces"); ok {
				spaces = s
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
			searchType, ok := getParam(params, "SearchType")
			if !ok {
				output = "I'm sorry I could not process your request " + intent + "."
			} else {
				go executeAction(room.RokuActionHost, "search", searchType)
			}
		default:
			output = "I'm sorry I could not process your request " + intent + "."
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(DialogflowResponse{FulfillmentText: output})
	}
}
