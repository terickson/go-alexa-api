package main

import (
	"fmt"
	"log"
	"strings"

	alexa "github.com/mikeflynn/go-alexa/skillserver"
)

// handleIntent returns an Alexa intent handler for the given room.
func handleIntent(room Room) func(*alexa.EchoRequest, *alexa.EchoResponse) {
	return func(echoReq *alexa.EchoRequest, echoResp *alexa.EchoResponse) {
		intent := strings.ToUpper(echoReq.GetIntentName())
		fmt.Println("Intent passed: " + intent)
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
			slotLevel, err := echoReq.GetSlotValue("Level")
			if err != nil {
				log.Println(err)
				output = "I'm sorry I could not process your request " + intent + "."
			} else {
				if room.ReceiverHost != "" {
					go updateReceiver(room.ReceiverHost, `{"volume": -`+slotLevel+`}`)
				} else {
					go executeAction(room.TVActionHost, "Volume", slotLevel)
				}
			}
		case "CHANNEL":
			slotNumber, err := echoReq.GetSlotValue("Number")
			if err != nil {
				log.Println(err)
				output = "I'm sorry I could not process your request " + intent + "."
			} else {
				go executeAction(room.TVActionHost, "Channel", slotNumber)
			}
		case "CHANNELUP":
			go executeAction(room.TVActionHost, "ChannelUp", "")
		case "CHANNELDOWN":
			go executeAction(room.TVActionHost, "ChannelDown", "")
		case "INPUT":
			slotInputType, err := echoReq.GetSlotValue("InputType")
			if err != nil {
				log.Println(err)
				output = "I'm sorry I could not process your request " + intent + "."
			} else {
				inputType := strings.ReplaceAll(strings.ToUpper(slotInputType), " ", "")
				setInput(room, inputType)
			}
		case "HOME":
			go executeAction(room.RokuActionHost, "home", "")
		case "BACK":
			go executeAction(room.RokuActionHost, "back", "")
		case "UP", "DOWN", "LEFT", "RIGHT":
			spaces := "1"
			slotSpaces, err := echoReq.GetSlotValue("Spaces")
			if err == nil && len(slotSpaces) > 0 {
				spaces = slotSpaces
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
			slotSearchType, err := echoReq.GetSlotValue("SearchType")
			if err != nil {
				log.Println(err)
				output = "I'm sorry I could not process your request " + intent + "."
			} else {
				go executeAction(room.RokuActionHost, "search", slotSearchType)
			}
		default:
			output = "I'm sorry I could not process your request " + intent + "."
		}

		echoResp.OutputSpeech(output)
	}
}
