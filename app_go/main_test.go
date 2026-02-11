package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"encoding/json"
)

func TestMainHandler_StatusOK(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	mainHandler(w, req)
	resp := w.Result()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	var data map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		t.Errorf("invalid JSON response: %v", err)
	}

	requiredFields := []string{"service", "system", "runtime", "request", "endpoints"}
	for _, field := range requiredFields {
		if _, ok := data[field]; !ok {
			t.Errorf("missing field in response: %s", field)
		}
	}
}

func TestMainHandler_JSONFields(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	mainHandler(w, req)
	resp := w.Result()

	if ct := resp.Header.Get("Content-Type"); ct != "application/json" {
		t.Errorf("expected application/json, got %s", ct)
	}
}

func TestHealthHandler_StatusOK(t *testing.T) {
	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()
	healthHandler(w, req)
	resp := w.Result()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	var healthResp map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&healthResp); err != nil {
		t.Errorf("invalid JSON response: %v", err)
	}

	if status, ok := healthResp["status"].(string); !ok || status != "healthy" {
		t.Errorf("expected status 'healthy', got %v", healthResp["status"])
	}
}

func TestMainHandler_NotFound(t *testing.T) {
	// Тестируем несуществующий путь
	req := httptest.NewRequest("GET", "/nonexistent", nil)
	w := httptest.NewRecorder()
	mainHandler(w, req)
	resp := w.Result()

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("expected status 404 for non-existent path, got %d", resp.StatusCode)
	}
}

func TestMainHandler_WithUserAgent(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("User-Agent", "TestAgent/1.0")
	w := httptest.NewRecorder()
	mainHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200 with User-Agent, got %d", w.Code)
	}
}

func TestMainHandler_PostAllowed(t *testing.T) {
	req := httptest.NewRequest("POST", "/", nil)
	w := httptest.NewRecorder()
	mainHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("POST request returned status %d, expected 200 if POST is allowed", w.Code)
	}
}

func TestHealthHandler_PostAllowed(t *testing.T) {
	req := httptest.NewRequest("POST", "/health", nil)
	w := httptest.NewRecorder()
	healthHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("POST to health returned status %d", w.Code)
	}
}