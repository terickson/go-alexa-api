package main

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"
)

func newGoogleRequest(intent string, params map[string]interface{}) *http.Request {
	body := DialogflowRequest{
		QueryResult: QueryResult{
			Intent:     DialogflowIntent{DisplayName: intent},
			Parameters: params,
		},
	}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/google/fr", strings.NewReader(string(b)))
	req.Header.Set("Content-Type", "application/json")
	return req
}

func decodeGoogleResponse(t *testing.T, rr *httptest.ResponseRecorder) DialogflowResponse {
	t.Helper()
	var resp DialogflowResponse
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode google response: %v", err)
	}
	return resp
}

func TestGoogleOff_WithReceiver(t *testing.T) {
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
	handler := handleGoogleIntent(room)

	rr := httptest.NewRecorder()
	handler(rr, newGoogleRequest("OFF", nil))
	time.Sleep(100 * time.Millisecond)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	mu.Lock()
	defer mu.Unlock()
	if len(tvCalls) != 1 || tvCalls[0] != `{"command": "PowerOff"}` {
		t.Errorf("unexpected TV calls: %v", tvCalls)
	}
	if len(receiverCalls) != 1 || receiverCalls[0] != `{"on": false}` {
		t.Errorf("unexpected receiver calls: %v", receiverCalls)
	}
}

func TestGoogleOff_NoReceiver(t *testing.T) {
	var tvBody string
	tvServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		tvBody = string(body)
	}))
	defer tvServer.Close()

	room := testRoom(tvServer.URL, "", "")
	handler := handleGoogleIntent(room)

	rr := httptest.NewRecorder()
	handler(rr, newGoogleRequest("OFF", nil))
	time.Sleep(100 * time.Millisecond)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	if tvBody != `{"command": "PowerOff"}` {
		t.Errorf("unexpected TV body: %q", tvBody)
	}
}

func TestGoogleMute_WithReceiver(t *testing.T) {
	var receiverBody string
	receiverServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		receiverBody = string(body)
	}))
	defer receiverServer.Close()
	tvServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	defer tvServer.Close()

	room := testRoom(tvServer.URL, "", receiverServer.URL)
	handler := handleGoogleIntent(room)

	rr := httptest.NewRecorder()
	handler(rr, newGoogleRequest("MUTE", nil))
	time.Sleep(100 * time.Millisecond)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	if receiverBody != `{"mute": true}` {
		t.Errorf("unexpected receiver body: %q", receiverBody)
	}
}

func TestGoogleMute_NoReceiver(t *testing.T) {
	var tvBody string
	tvServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		tvBody = string(body)
	}))
	defer tvServer.Close()

	room := testRoom(tvServer.URL, "", "")
	handler := handleGoogleIntent(room)

	rr := httptest.NewRecorder()
	handler(rr, newGoogleRequest("MUTE", nil))
	time.Sleep(100 * time.Millisecond)

	if tvBody != `{"command": "Mute"}` {
		t.Errorf("unexpected TV body: %q", tvBody)
	}
}

func TestGoogleVolume_WithReceiver(t *testing.T) {
	var receiverBody string
	receiverServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		receiverBody = string(body)
	}))
	defer receiverServer.Close()
	tvServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	defer tvServer.Close()

	room := testRoom(tvServer.URL, "", receiverServer.URL)
	handler := handleGoogleIntent(room)

	rr := httptest.NewRecorder()
	handler(rr, newGoogleRequest("VOLUME", map[string]interface{}{"Level": float64(30)}))
	time.Sleep(100 * time.Millisecond)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	if receiverBody != `{"volume": -30}` {
		t.Errorf("unexpected receiver body: %q", receiverBody)
	}
}

func TestGoogleVolume_NoParam(t *testing.T) {
	tvServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	defer tvServer.Close()

	room := testRoom(tvServer.URL, "", "")
	handler := handleGoogleIntent(room)

	rr := httptest.NewRecorder()
	handler(rr, newGoogleRequest("VOLUME", nil))

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	resp := decodeGoogleResponse(t, rr)
	if !strings.Contains(resp.FulfillmentText, "I'm sorry") {
		t.Errorf("expected error text, got %q", resp.FulfillmentText)
	}
}

func TestGoogleChannel(t *testing.T) {
	var tvBody string
	tvServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		tvBody = string(body)
	}))
	defer tvServer.Close()

	room := testRoom(tvServer.URL, "", "")
	handler := handleGoogleIntent(room)

	rr := httptest.NewRecorder()
	handler(rr, newGoogleRequest("CHANNEL", map[string]interface{}{"Number": float64(5)}))
	time.Sleep(100 * time.Millisecond)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	if tvBody != `{"command": "Channel", "value": "5"}` {
		t.Errorf("unexpected TV body: %q", tvBody)
	}
}

func TestGoogleDirectional_WithSpaces(t *testing.T) {
	var rokuBody string
	rokuServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		rokuBody = string(body)
	}))
	defer rokuServer.Close()
	tvServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	defer tvServer.Close()

	room := testRoom(tvServer.URL, rokuServer.URL, "")
	handler := handleGoogleIntent(room)

	rr := httptest.NewRecorder()
	handler(rr, newGoogleRequest("UP", map[string]interface{}{"Spaces": float64(3)}))
	time.Sleep(100 * time.Millisecond)

	if rokuBody != `{"command": "up", "value": "3"}` {
		t.Errorf("unexpected Roku body: %q", rokuBody)
	}
}

func TestGoogleDirectional_DefaultSpaces(t *testing.T) {
	var rokuBody string
	rokuServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		rokuBody = string(body)
	}))
	defer rokuServer.Close()
	tvServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	defer tvServer.Close()

	room := testRoom(tvServer.URL, rokuServer.URL, "")
	handler := handleGoogleIntent(room)

	rr := httptest.NewRecorder()
	handler(rr, newGoogleRequest("UP", nil))
	time.Sleep(100 * time.Millisecond)

	if rokuBody != `{"command": "up", "value": "1"}` {
		t.Errorf("unexpected Roku body: %q", rokuBody)
	}
}

func TestGoogleSearch(t *testing.T) {
	var rokuBody string
	rokuServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		rokuBody = string(body)
	}))
	defer rokuServer.Close()
	tvServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	defer tvServer.Close()

	room := testRoom(tvServer.URL, rokuServer.URL, "")
	handler := handleGoogleIntent(room)

	rr := httptest.NewRecorder()
	handler(rr, newGoogleRequest("SEARCH", map[string]interface{}{"SearchType": "breaking bad"}))
	time.Sleep(100 * time.Millisecond)

	if rokuBody != `{"command": "search", "value": "breaking bad"}` {
		t.Errorf("unexpected Roku body: %q", rokuBody)
	}
}

func TestGoogleUnknownIntent(t *testing.T) {
	tvServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	defer tvServer.Close()

	room := testRoom(tvServer.URL, "", "")
	handler := handleGoogleIntent(room)

	rr := httptest.NewRecorder()
	handler(rr, newGoogleRequest("FOOBAR", nil))

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	resp := decodeGoogleResponse(t, rr)
	if resp.FulfillmentText != "I'm sorry I could not process your request FOOBAR." {
		t.Errorf("unexpected response: %q", resp.FulfillmentText)
	}
}

func TestGoogleTokenAuth_Missing(t *testing.T) {
	t.Setenv("GOOGLE_WEBHOOK_TOKEN", "secret123")

	tvServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	defer tvServer.Close()

	room := testRoom(tvServer.URL, "", "")
	handler := handleGoogleIntent(room)

	rr := httptest.NewRecorder()
	handler(rr, newGoogleRequest("OFF", nil))

	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", rr.Code)
	}
}

func TestGoogleTokenAuth_Wrong(t *testing.T) {
	t.Setenv("GOOGLE_WEBHOOK_TOKEN", "secret123")

	tvServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	defer tvServer.Close()

	room := testRoom(tvServer.URL, "", "")
	handler := handleGoogleIntent(room)

	rr := httptest.NewRecorder()
	req := newGoogleRequest("OFF", nil)
	req.URL.RawQuery = "token=wrongtoken"
	handler(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", rr.Code)
	}
}

func TestGoogleTokenAuth_Valid(t *testing.T) {
	t.Setenv("GOOGLE_WEBHOOK_TOKEN", "secret123")

	tvServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	defer tvServer.Close()

	room := testRoom(tvServer.URL, "", "")
	handler := handleGoogleIntent(room)

	rr := httptest.NewRecorder()
	req := newGoogleRequest("OFF", nil)
	req.URL.RawQuery = "token=secret123"
	handler(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
}

func TestGoogleInput_Netflix(t *testing.T) {
	var mu sync.Mutex
	var tvCalls, rokuCalls []string

	tvServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		mu.Lock()
		tvCalls = append(tvCalls, string(body))
		mu.Unlock()
	}))
	defer tvServer.Close()

	rokuServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		mu.Lock()
		rokuCalls = append(rokuCalls, string(body))
		mu.Unlock()
	}))
	defer rokuServer.Close()

	room := testRoom(tvServer.URL, rokuServer.URL, "")
	handler := handleGoogleIntent(room)

	rr := httptest.NewRecorder()
	// setInput blocks for 500ms internally; handler returns after setInput completes
	handler(rr, newGoogleRequest("INPUT", map[string]interface{}{"InputType": "netflix"}))
	time.Sleep(200 * time.Millisecond)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	mu.Lock()
	defer mu.Unlock()
	if len(tvCalls) == 0 {
		t.Error("expected TV calls for INPUT, got none")
	}
	if len(rokuCalls) == 0 {
		t.Error("expected Roku app launch for netflix INPUT, got none")
	}
}
