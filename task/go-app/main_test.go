package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestTimeHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/time", nil)
	rr := httptest.NewRecorder()

	timeHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}

	var resp TimeResponse
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.Time == "" {
		t.Fatal("time field is empty")
	}

	if _, err := time.Parse(time.RFC3339, resp.Time); err != nil {
		t.Fatalf("time is not RFC3339: %v", err)
	}
}

func TestHealthHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rr := httptest.NewRecorder()

	healthHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}

	var resp HealthResponse
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.Status != "ok" {
		t.Fatalf("expected status ok, got %s", resp.Status)
	}
}

func TestNewServer(t *testing.T) {
	srv := newServer()
	if srv.Addr != ":8080" {
		t.Fatalf("expected addr :8080, got %s", srv.Addr)
	}
	if srv.Handler == nil {
		t.Fatal("handler is nil")
	}
}

func TestTimeHandlerContentType(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/time", nil)
	rr := httptest.NewRecorder()

	timeHandler(rr, req)

	ct := rr.Header().Get("Content-Type")
	if ct != "application/json" {
		t.Fatalf("expected application/json, got %s", ct)
	}
}

func TestHealthHandlerContentType(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rr := httptest.NewRecorder()

	healthHandler(rr, req)

	ct := rr.Header().Get("Content-Type")
	if ct != "application/json" {
		t.Fatalf("expected application/json, got %s", ct)
	}
}
