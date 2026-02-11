package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMainHandler_StatusOK(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	mainHandler(w, req)  // Исправлено: mainHandler вместо indexHandler
	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}
}

func TestMainHandler_JSONFields(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	mainHandler(w, req)  // Исправлено: mainHandler вместо indexHandler
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
}

func TestHealthHandler_JSONFields(t *testing.T) {
	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()
	healthHandler(w, req)
	resp := w.Result()
	if ct := resp.Header.Get("Content-Type"); ct != "application/json" {
		t.Errorf("expected application/json, got %s", ct)
	}
}

func TestMainHandler_MethodNotAllowed(t *testing.T) {
	req := httptest.NewRequest("POST", "/", nil)
	w := httptest.NewRecorder()
	mainHandler(w, req)  // Исправлено: mainHandler вместо indexHandler
	resp := w.Result()
	if resp.StatusCode == http.StatusOK {
		t.Errorf("expected not 200 for POST, got %d", resp.StatusCode)
	}
}

func TestHealthHandler_MethodNotAllowed(t *testing.T) {
	req := httptest.NewRequest("POST", "/health", nil)
	w := httptest.NewRecorder()
	healthHandler(w, req)
	resp := w.Result()
	if resp.StatusCode == http.StatusOK {
		t.Errorf("expected not 200 for POST, got %d", resp.StatusCode)
	}
}

// Добавьте дополнительные тесты для покрытия

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

func TestHealthHandler_Content(t *testing.T) {
	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()
	healthHandler(w, req)

	if w.Body.Len() == 0 {
		t.Error("health handler returned empty body")
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
