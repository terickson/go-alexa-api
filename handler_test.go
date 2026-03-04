package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	alexa "github.com/mikeflynn/go-alexa/skillserver"
)

func newEchoRequest(intentName string, slots map[string]string) *alexa.EchoRequest {
	req := &alexa.EchoRequest{}
	req.Request.Type = "IntentRequest"
	req.Request.Intent.Name = intentName
	if len(slots) > 0 {
		req.Request.Intent.Slots = make(map[string]alexa.EchoSlot)
		for k, v := range slots {
			req.Request.Intent.Slots[k] = alexa.EchoSlot{Name: k, Value: v}
		}
	}
	return req
}

func testRoom(tvURL, rokuURL, receiverURL string) Room {
	return Room{
		Name:           "Test Room",
		TVActionHost:   tvURL,
		RokuActionHost: rokuURL,
		ReceiverHost:   receiverURL,
		DefaultVolume:  -30,
		InputMap: map[string]InputConfig{
			"NETFLIX": {ReceiverInput: "HDMI1", TVInput: "HDMI1", RokuApp: "Netflix"},
			"TV":      {ReceiverInput: "AV1", TVInput: "InputTV"},
		},
	}
}

func TestHandleIntent_OFF_WithReceiver(t *testing.T) {
	var mu sync.Mutex
	var tvCalls, receiverCalls []string

	tvServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		mu.Lock()
		tvCalls = append(tvCalls, string(body))
		mu.Unlock()
	}))
	defer tvServer.Close()

	receiverServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		mu.Lock()
		receiverCalls = append(receiverCalls, string(body))
		mu.Unlock()
	}))
	defer receiverServer.Close()

	room := testRoom(tvServer.URL, "", receiverServer.URL)
	handler := handleIntent(room)

	req := newEchoRequest("OFF", nil)
	resp := alexa.NewEchoResponse()
	handler(req, resp)

	time.Sleep(100 * time.Millisecond)

	mu.Lock()
	defer mu.Unlock()

	if len(tvCalls) != 1 {
		t.Fatalf("expected 1 TV call, got %d", len(tvCalls))
	}
	if tvCalls[0] != `{"command": "PowerOff"}` {
		t.Errorf("unexpected TV call: %s", tvCalls[0])
	}
	if len(receiverCalls) != 1 {
		t.Fatalf("expected 1 receiver call, got %d", len(receiverCalls))
	}
	if receiverCalls[0] != `{"on": false}` {
		t.Errorf("unexpected receiver call: %s", receiverCalls[0])
	}
}

func TestHandleIntent_OFF_NoReceiver(t *testing.T) {
	var mu sync.Mutex
	var tvCalls []string

	tvServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		mu.Lock()
		tvCalls = append(tvCalls, string(body))
		mu.Unlock()
	}))
	defer tvServer.Close()

	room := testRoom(tvServer.URL, "", "")
	handler := handleIntent(room)

	req := newEchoRequest("OFF", nil)
	resp := alexa.NewEchoResponse()
	handler(req, resp)

	time.Sleep(100 * time.Millisecond)

	mu.Lock()
	defer mu.Unlock()

	if len(tvCalls) != 1 {
		t.Fatalf("expected 1 TV call, got %d", len(tvCalls))
	}
}

func TestHandleIntent_MUTE_WithReceiver(t *testing.T) {
	var receiverBody string
	receiverServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		receiverBody = string(body)
	}))
	defer receiverServer.Close()

	tvServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	defer tvServer.Close()

	room := testRoom(tvServer.URL, "", receiverServer.URL)
	handler := handleIntent(room)

	req := newEchoRequest("MUTE", nil)
	resp := alexa.NewEchoResponse()
	handler(req, resp)

	time.Sleep(100 * time.Millisecond)

	if receiverBody != `{"mute": true}` {
		t.Errorf("expected mute payload, got %q", receiverBody)
	}
}

func TestHandleIntent_MUTE_NoReceiver(t *testing.T) {
	var tvBody string
	tvServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		tvBody = string(body)
	}))
	defer tvServer.Close()

	room := testRoom(tvServer.URL, "", "")
	handler := handleIntent(room)

	req := newEchoRequest("MUTE", nil)
	resp := alexa.NewEchoResponse()
	handler(req, resp)

	time.Sleep(100 * time.Millisecond)

	if tvBody != `{"command": "Mute"}` {
		t.Errorf("expected Mute command, got %q", tvBody)
	}
}

func TestHandleIntent_CHANNEL(t *testing.T) {
	var tvBody string
	tvServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		tvBody = string(body)
	}))
	defer tvServer.Close()

	room := testRoom(tvServer.URL, "", "")
	handler := handleIntent(room)

	req := newEchoRequest("Channel", map[string]string{"Number": "42"})
	resp := alexa.NewEchoResponse()
	handler(req, resp)

	time.Sleep(100 * time.Millisecond)

	if tvBody != `{"command": "Channel", "value": "42"}` {
		t.Errorf("unexpected TV body: %s", tvBody)
	}
}

func TestHandleIntent_HOME(t *testing.T) {
	var rokuBody string
	rokuServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		rokuBody = string(body)
	}))
	defer rokuServer.Close()

	tvServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	defer tvServer.Close()

	room := testRoom(tvServer.URL, rokuServer.URL, "")
	handler := handleIntent(room)

	req := newEchoRequest("HOME", nil)
	resp := alexa.NewEchoResponse()
	handler(req, resp)

	time.Sleep(100 * time.Millisecond)

	if rokuBody != `{"command": "home"}` {
		t.Errorf("unexpected Roku body: %s", rokuBody)
	}
}

func TestHandleIntent_Direction_WithSpaces(t *testing.T) {
	var rokuBody string
	rokuServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		rokuBody = string(body)
	}))
	defer rokuServer.Close()

	tvServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	defer tvServer.Close()

	room := testRoom(tvServer.URL, rokuServer.URL, "")
	handler := handleIntent(room)

	req := newEchoRequest("UP", map[string]string{"Spaces": "3"})
	resp := alexa.NewEchoResponse()
	handler(req, resp)

	time.Sleep(100 * time.Millisecond)

	if rokuBody != `{"command": "up", "value": "3"}` {
		t.Errorf("unexpected Roku body: %s", rokuBody)
	}
}

func TestHandleIntent_Default(t *testing.T) {
	tvServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	defer tvServer.Close()

	room := testRoom(tvServer.URL, "", "")
	handler := handleIntent(room)

	req := newEchoRequest("UNKNOWNINTENT", nil)
	resp := alexa.NewEchoResponse()
	handler(req, resp)

	if resp.Response.OutputSpeech == nil {
		t.Fatal("expected output speech")
	}
	if resp.Response.OutputSpeech.Text != "I'm sorry I could not process your request UNKNOWNINTENT." {
		t.Errorf("unexpected output: %s", resp.Response.OutputSpeech.Text)
	}
}

func TestHandleIntent_SlotError(t *testing.T) {
	tvServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	defer tvServer.Close()

	room := testRoom(tvServer.URL, "", "")
	handler := handleIntent(room)

	// CHANNEL without Number slot
	req := newEchoRequest("Channel", nil)
	resp := alexa.NewEchoResponse()
	handler(req, resp)

	if resp.Response.OutputSpeech == nil {
		t.Fatal("expected output speech")
	}
	if resp.Response.OutputSpeech.Text != "I'm sorry I could not process your request CHANNEL." {
		t.Errorf("unexpected output: %s", resp.Response.OutputSpeech.Text)
	}
}
