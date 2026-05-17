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

func newActionRequest(params map[string]string) *http.Request {
	req := httptest.NewRequest(http.MethodGet, "/action/fr", nil)
	q := req.URL.Query()
	for k, v := range params {
		q.Set(k, v)
	}
	req.URL.RawQuery = q.Encode()
	return req
}

func decodeActionResponse(t *testing.T, rr *httptest.ResponseRecorder) SimpleActionResponse {
	t.Helper()
	var resp SimpleActionResponse
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode action response: %v", err)
	}
	return resp
}

func TestSimpleAction_OFF_WithReceiver(t *testing.T) {
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
	handler := handleSimpleAction(room)

	rr := httptest.NewRecorder()
	handler(rr, newActionRequest(map[string]string{"intent": "OFF"}))
	time.Sleep(100 * time.Millisecond)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	resp := decodeActionResponse(t, rr)
	if resp.Status != "ok" {
		t.Errorf("expected status ok, got %q", resp.Status)
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

func TestSimpleAction_OFF_NoReceiver(t *testing.T) {
	var tvBody string
	tvServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		tvBody = string(body)
	}))
	defer tvServer.Close()

	room := testRoom(tvServer.URL, "", "")
	handler := handleSimpleAction(room)

	rr := httptest.NewRecorder()
	handler(rr, newActionRequest(map[string]string{"intent": "OFF"}))
	time.Sleep(100 * time.Millisecond)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	if tvBody != `{"command": "PowerOff"}` {
		t.Errorf("unexpected TV body: %q", tvBody)
	}
}

func TestSimpleAction_MUTE_WithReceiver(t *testing.T) {
	var receiverBody string
	receiverServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		receiverBody = string(body)
	}))
	defer receiverServer.Close()
	tvServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	defer tvServer.Close()

	room := testRoom(tvServer.URL, "", receiverServer.URL)
	handler := handleSimpleAction(room)

	rr := httptest.NewRecorder()
	handler(rr, newActionRequest(map[string]string{"intent": "MUTE"}))
	time.Sleep(100 * time.Millisecond)

	if receiverBody != `{"mute": true}` {
		t.Errorf("unexpected receiver body: %q", receiverBody)
	}
}

func TestSimpleAction_VOLUME(t *testing.T) {
	var receiverBody string
	receiverServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		receiverBody = string(body)
	}))
	defer receiverServer.Close()
	tvServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	defer tvServer.Close()

	room := testRoom(tvServer.URL, "", receiverServer.URL)
	handler := handleSimpleAction(room)

	rr := httptest.NewRecorder()
	handler(rr, newActionRequest(map[string]string{"intent": "VOLUME", "level": "30"}))
	time.Sleep(100 * time.Millisecond)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	if receiverBody != `{"volume": -30}` {
		t.Errorf("unexpected receiver body: %q", receiverBody)
	}
}

func TestSimpleAction_VOLUME_NoParam(t *testing.T) {
	tvServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	defer tvServer.Close()

	room := testRoom(tvServer.URL, "", "")
	handler := handleSimpleAction(room)

	rr := httptest.NewRecorder()
	handler(rr, newActionRequest(map[string]string{"intent": "VOLUME"}))

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	resp := decodeActionResponse(t, rr)
	if resp.Status != "error" {
		t.Errorf("expected status error, got %q", resp.Status)
	}
	if !strings.Contains(resp.Message, "I'm sorry") {
		t.Errorf("expected error message, got %q", resp.Message)
	}
}

func TestSimpleAction_CHANNEL(t *testing.T) {
	var tvBody string
	tvServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		tvBody = string(body)
	}))
	defer tvServer.Close()

	room := testRoom(tvServer.URL, "", "")
	handler := handleSimpleAction(room)

	rr := httptest.NewRecorder()
	handler(rr, newActionRequest(map[string]string{"intent": "CHANNEL", "number": "5"}))
	time.Sleep(100 * time.Millisecond)

	if tvBody != `{"command": "Channel", "value": "5"}` {
		t.Errorf("unexpected TV body: %q", tvBody)
	}
}

func TestSimpleAction_INPUT(t *testing.T) {
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
	handler := handleSimpleAction(room)

	rr := httptest.NewRecorder()
	// setInput blocks for 500ms internally
	handler(rr, newActionRequest(map[string]string{"intent": "INPUT", "input": "netflix"}))
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

func TestSimpleAction_Directional_WithSpaces(t *testing.T) {
	var rokuBody string
	rokuServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		rokuBody = string(body)
	}))
	defer rokuServer.Close()
	tvServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	defer tvServer.Close()

	room := testRoom(tvServer.URL, rokuServer.URL, "")
	handler := handleSimpleAction(room)

	rr := httptest.NewRecorder()
	handler(rr, newActionRequest(map[string]string{"intent": "UP", "spaces": "3"}))
	time.Sleep(100 * time.Millisecond)

	if rokuBody != `{"command": "up", "value": "3"}` {
		t.Errorf("unexpected Roku body: %q", rokuBody)
	}
}

func TestSimpleAction_Directional_DefaultSpaces(t *testing.T) {
	var rokuBody string
	rokuServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		rokuBody = string(body)
	}))
	defer rokuServer.Close()
	tvServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	defer tvServer.Close()

	room := testRoom(tvServer.URL, rokuServer.URL, "")
	handler := handleSimpleAction(room)

	rr := httptest.NewRecorder()
	handler(rr, newActionRequest(map[string]string{"intent": "UP"}))
	time.Sleep(100 * time.Millisecond)

	if rokuBody != `{"command": "up", "value": "1"}` {
		t.Errorf("unexpected Roku body: %q", rokuBody)
	}
}

func TestSimpleAction_SEARCH(t *testing.T) {
	var rokuBody string
	rokuServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		rokuBody = string(body)
	}))
	defer rokuServer.Close()
	tvServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	defer tvServer.Close()

	room := testRoom(tvServer.URL, rokuServer.URL, "")
	handler := handleSimpleAction(room)

	rr := httptest.NewRecorder()
	handler(rr, newActionRequest(map[string]string{"intent": "SEARCH", "search": "breaking bad"}))
	time.Sleep(100 * time.Millisecond)

	if rokuBody != `{"command": "search", "value": "breaking bad"}` {
		t.Errorf("unexpected Roku body: %q", rokuBody)
	}
}

func TestSimpleAction_UnknownIntent(t *testing.T) {
	tvServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	defer tvServer.Close()

	room := testRoom(tvServer.URL, "", "")
	handler := handleSimpleAction(room)

	rr := httptest.NewRecorder()
	handler(rr, newActionRequest(map[string]string{"intent": "FOOBAR"}))

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	resp := decodeActionResponse(t, rr)
	if resp.Status != "error" {
		t.Errorf("expected status error, got %q", resp.Status)
	}
}

func TestSimpleAction_TokenAuth_Missing(t *testing.T) {
	t.Setenv("GOOGLE_WEBHOOK_TOKEN", "secret123")

	tvServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	defer tvServer.Close()

	room := testRoom(tvServer.URL, "", "")
	handler := handleSimpleAction(room)

	rr := httptest.NewRecorder()
	handler(rr, newActionRequest(map[string]string{"intent": "OFF"}))

	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", rr.Code)
	}
}

func TestSimpleAction_TokenAuth_Wrong(t *testing.T) {
	t.Setenv("GOOGLE_WEBHOOK_TOKEN", "secret123")

	tvServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	defer tvServer.Close()

	room := testRoom(tvServer.URL, "", "")
	handler := handleSimpleAction(room)

	rr := httptest.NewRecorder()
	handler(rr, newActionRequest(map[string]string{"intent": "OFF", "token": "wrongtoken"}))

	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", rr.Code)
	}
}

func TestSimpleAction_TokenAuth_Valid(t *testing.T) {
	t.Setenv("GOOGLE_WEBHOOK_TOKEN", "secret123")

	tvServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	defer tvServer.Close()

	room := testRoom(tvServer.URL, "", "")
	handler := handleSimpleAction(room)

	rr := httptest.NewRecorder()
	handler(rr, newActionRequest(map[string]string{"intent": "OFF", "token": "secret123"}))

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
}
