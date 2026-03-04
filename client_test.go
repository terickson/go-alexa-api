package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestExecuteAction_NoValue(t *testing.T) {
	var receivedMethod, receivedBody, receivedContentType string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedMethod = r.Method
		receivedContentType = r.Header.Get("Content-Type")
		body, _ := io.ReadAll(r.Body)
		receivedBody = string(body)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	executeAction(server.URL, "PowerOff", "")

	if receivedMethod != "POST" {
		t.Errorf("expected POST, got %s", receivedMethod)
	}
	if receivedContentType != "application/json" {
		t.Errorf("expected application/json, got %s", receivedContentType)
	}
	expected := `{"command": "PowerOff"}`
	if receivedBody != expected {
		t.Errorf("expected body %q, got %q", expected, receivedBody)
	}
}

func TestExecuteAction_WithValue(t *testing.T) {
	var receivedBody string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		receivedBody = string(body)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	executeAction(server.URL, "Channel", "42")

	expected := `{"command": "Channel", "value": "42"}`
	if receivedBody != expected {
		t.Errorf("expected body %q, got %q", expected, receivedBody)
	}
}

func TestExecuteAction_ServerDown(t *testing.T) {
	// Should log error but not panic/fatal
	executeAction("http://127.0.0.1:1", "PowerOff", "")
}

func TestUpdateReceiver(t *testing.T) {
	var receivedMethod, receivedBody, receivedContentType string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedMethod = r.Method
		receivedContentType = r.Header.Get("Content-Type")
		body, _ := io.ReadAll(r.Body)
		receivedBody = string(body)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	payload := `{"on": false}`
	updateReceiver(server.URL, payload)

	if receivedMethod != http.MethodPut {
		t.Errorf("expected PUT, got %s", receivedMethod)
	}
	if receivedContentType != "application/json" {
		t.Errorf("expected application/json, got %s", receivedContentType)
	}
	if receivedBody != payload {
		t.Errorf("expected body %q, got %q", payload, receivedBody)
	}
}

func TestUpdateReceiver_ServerDown(t *testing.T) {
	// Should log error but not panic/fatal
	updateReceiver("http://127.0.0.1:1", `{"on": false}`)
}

func TestUpdateReceiver_MutePayload(t *testing.T) {
	var receivedBody string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		receivedBody = string(body)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	updateReceiver(server.URL, `{"mute": true}`)

	if !strings.Contains(receivedBody, `"mute": true`) {
		t.Errorf("expected mute payload, got %q", receivedBody)
	}
}
