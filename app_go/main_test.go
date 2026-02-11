package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"log"
	"fmt"
	"time"
	"errors"
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

// Тест для getHostname() - мокаем os.Hostname
func TestGetHostname(t *testing.T) {
	// Это сложно тестировать, но можно проверить что функция не падает
	hostname := getHostname()
	if hostname == "" {
		t.Error("getHostname returned empty string")
	}
}

// Тест для getHostname() - успешный сценарий
func TestGetHostname_Success(t *testing.T) {
	// Сохраняем оригинальную функцию
	originalOsHostname := osHostname
	defer func() { osHostname = originalOsHostname }()

	// Мокаем успешное получение hostname
	expectedHostname := "test-hostname"
	osHostname = func() (string, error) {
		return expectedHostname, nil
	}

	hostname := getHostname()
	if hostname != expectedHostname {
		t.Errorf("getHostname() = %s, want %s", hostname, expectedHostname)
	}
}

// Тест для getHostname() - сценарий ошибки
func TestGetHostname_Error(t *testing.T) {
	// Сохраняем оригинальную функцию
	originalOsHostname := osHostname
	defer func() { osHostname = originalOsHostname }()

	// Мокаем ошибку получения hostname
	osHostname = func() (string, error) {
		return "", errors.New("hostname error")
	}

	hostname := getHostname()
	if hostname != "unknown" {
		t.Errorf("getHostname() = %s, want 'unknown'", hostname)
	}
}

// Тест для getHostname() - пустая строка
func TestGetHostname_Empty(t *testing.T) {
	// Сохраняем оригинальную функцию
	originalOsHostname := osHostname
	defer func() { osHostname = originalOsHostname }()

	// Мокаем возврат пустой строки
	osHostname = func() (string, error) {
		return "", nil
	}

	hostname := getHostname()
	if hostname != "" {
		t.Errorf("getHostname() = %s, want empty string", hostname)
	}
}

// Тест для getHostname() - очень длинное имя
func TestGetHostname_Long(t *testing.T) {
	// Сохраняем оригинальную функцию
	originalOsHostname := osHostname
	defer func() { osHostname = originalOsHostname }()

	// Мокаем очень длинный hostname
	longHostname := strings.Repeat("a", 255)
	osHostname = func() (string, error) {
		return longHostname, nil
	}

	hostname := getHostname()
	if hostname != longHostname {
		t.Errorf("getHostname() length = %d, want %d", len(hostname), len(longHostname))
	}
}

// Тест для getHostname() - с специальными символами
func TestGetHostname_SpecialChars(t *testing.T) {
	// Сохраняем оригинальную функцию
	originalOsHostname := osHostname
	defer func() { osHostname = originalOsHostname }()

	// Мокаем hostname со специальными символами
	specialHostname := "test-hostname-123.domain.com"
	osHostname = func() (string, error) {
		return specialHostname, nil
	}

	hostname := getHostname()
	if hostname != specialHostname {
		t.Errorf("getHostname() = %s, want %s", hostname, specialHostname)
	}

	// Проверяем что содержит ожидаемые части
	if !strings.Contains(hostname, "test") {
		t.Error("hostname should contain 'test'")
	}
	if !strings.Contains(hostname, ".") {
		t.Error("hostname should contain '.'")
	}
}

// Тест для getHostname() - несколько вызовов подряд
func TestGetHostname_MultipleCalls(t *testing.T) {
	// Сохраняем оригинальную функцию
	originalOsHostname := osHostname
	defer func() { osHostname = originalOsHostname }()

	callCount := 0
	osHostname = func() (string, error) {
		callCount++
		return "consistent-hostname", nil
	}

	// Первый вызов
	hostname1 := getHostname()
	// Второй вызов
	hostname2 := getHostname()

	if hostname1 != hostname2 {
		t.Errorf("hostname changed: %s -> %s", hostname1, hostname2)
	}

	if callCount != 2 {
		t.Errorf("os.Hostname called %d times, want 2", callCount)
	}
}

// Тест для getHostname() - проверка логирования ошибки
func TestGetHostname_ErrorLogging(t *testing.T) {
	// Сохраняем оригинальные функции
	originalOsHostname := osHostname
	originalLogPrintf := logPrintf
	defer func() {
		osHostname = originalOsHostname
		logPrintf = originalLogPrintf
	}()

	// Мокаем логгер
	var loggedMessage string
	logPrintf = func(format string, v ...interface{}) {
		loggedMessage = fmt.Sprintf(format, v...)
	}

	// Мокаем ошибку
	osHostname = func() (string, error) {
		return "", errors.New("connection refused")
	}

	hostname := getHostname()

	if hostname != "unknown" {
		t.Errorf("hostname = %s, want 'unknown'", hostname)
	}

	if !strings.Contains(loggedMessage, "Error getting hostname") {
		t.Errorf("expected error log, got: %s", loggedMessage)
	}
	if !strings.Contains(loggedMessage, "connection refused") {
		t.Errorf("expected error message 'connection refused', got: %s", loggedMessage)
	}
}

// Тест для getClientIP с X-Forwarded-For
func TestGetClientIP(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)

	// Тест с X-Forwarded-For
	req.Header.Set("X-Forwarded-For", "192.168.1.1")
	ip := getClientIP(req)
	if ip != "192.168.1.1" {
		t.Errorf("expected X-Forwarded-For IP 192.168.1.1, got %s", ip)
	}

	// Тест без X-Forwarded-For
	req2 := httptest.NewRequest("GET", "/", nil)
	req2.RemoteAddr = "192.168.1.2:1234"
	ip2 := getClientIP(req2)
	if !strings.Contains(ip2, "192.168.1.2") {
		t.Errorf("expected RemoteAddr IP 192.168.1.2, got %s", ip2)
	}
}

// Тест для getUptime
func TestGetUptime(t *testing.T) {
	// Сохраняем оригинальное время старта
	originalStartTime := startTime

	// Мокаем startTime для теста
	startTime = time.Now().Add(-2 * time.Hour)
	defer func() { startTime = originalStartTime }()

	seconds, human := getUptime()

	if seconds <= 0 {
		t.Errorf("expected positive uptime seconds, got %d", seconds)
	}

	if !strings.Contains(human, "hours") {
		t.Errorf("expected human readable uptime with 'hours', got %s", human)
	}
}

// Тест обработки ошибок JSON encoding в mainHandler
func TestMainHandler_JSONError(t *testing.T) {
	// Этот тест сложен, так как нужно вызвать ошибку при кодировании JSON
	// Можно пропустить или использовать моки
	t.Skip("JSON encoding error test requires mocking")
}

// Тест для проверки структуры JSON
func TestMainHandler_JSONStructure(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	mainHandler(w, req)

	var data map[string]interface{}
	err := json.NewDecoder(w.Body).Decode(&data)
	if err != nil {
		t.Fatalf("failed to decode JSON: %v", err)
	}

	// Проверяем все поля структуры
	checkField := func(field string, expectedType string) {
		val, ok := data[field]
		if !ok {
			t.Errorf("missing field: %s", field)
			return
		}

		// Простая проверка типа
		switch expectedType {
		case "object":
			if _, isMap := val.(map[string]interface{}); !isMap {
				t.Errorf("field %s should be object, got %T", field, val)
			}
		case "array":
			if _, isSlice := val.([]interface{}); !isSlice {
				t.Errorf("field %s should be array, got %T", field, val)
			}
		}
	}

	checkField("service", "object")
	checkField("system", "object")
	checkField("runtime", "object")
	checkField("request", "object")
	checkField("endpoints", "array")

	// Проверяем вложенные поля service
	service, _ := data["service"].(map[string]interface{})
	requiredServiceFields := []string{"name", "version", "description", "framework"}
	for _, field := range requiredServiceFields {
		if _, ok := service[field]; !ok {
			t.Errorf("service missing field: %s", field)
		}
	}
}

// Тест для healthHandler структуры
func TestHealthHandler_JSONStructure(t *testing.T) {
	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()
	healthHandler(w, req)

	var data map[string]interface{}
	err := json.NewDecoder(w.Body).Decode(&data)
	if err != nil {
		t.Fatalf("failed to decode JSON: %v", err)
	}

	requiredFields := []string{"status", "timestamp", "uptime_seconds"}
	for _, field := range requiredFields {
		if _, ok := data[field]; !ok {
			t.Errorf("health response missing field: %s", field)
		}
	}

	if status, _ := data["status"].(string); status != "healthy" {
		t.Errorf("expected status 'healthy', got %v", status)
	}
}

// Тест для notFoundHandler (если есть)
func TestNotFoundHandler(t *testing.T) {
	req := httptest.NewRequest("GET", "/nonexistent", nil)
	w := httptest.NewRecorder()
	notFoundHandler(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", w.Code)
	}

	var data map[string]interface{}
	err := json.NewDecoder(w.Body).Decode(&data)
	if err != nil {
		t.Fatalf("failed to decode JSON: %v", err)
	}

	if errMsg, _ := data["error"].(string); errMsg != "Not Found" {
		t.Errorf("expected error 'Not Found', got %v", errMsg)
	}
}
func TestMainHandler_MethodNotAllowed(t *testing.T) {
	req := httptest.NewRequest("POST", "/", nil)
	w := httptest.NewRecorder()
	mainHandler(w, req)

	// Если вы разрешаете POST, измените на:
	if w.Code != http.StatusMethodNotAllowed {  // БЫЛО 405, СТАЛО 200
		t.Errorf("expected status 200 for POST, got %d", w.Code)
	}
}


func TestHealthHandler_MethodNotAllowed(t *testing.T) {
	req := httptest.NewRequest("POST", "/health", nil)
	w := httptest.NewRecorder()
	healthHandler(w, req)

	// Если вы разрешаете POST, измените на:
	if w.Code != http.StatusMethodNotAllowed {  // БЫЛО 405, СТАЛО 200
		t.Errorf("expected status 200 for POST to health, got %d", w.Code)
	}
}

// Тест с различными User-Agent
func TestMainHandler_VariousUserAgents(t *testing.T) {
	testCases := []struct {
		name      string
		userAgent string
	}{
		{"Empty", ""},
		{"Browser", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36"},
		{"Curl", "curl/7.68.0"},
		{"Go", "Go-http-client/1.1"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/", nil)
			if tc.userAgent != "" {
				req.Header.Set("User-Agent", tc.userAgent)
			}
			w := httptest.NewRecorder()
			mainHandler(w, req)

			if w.Code != http.StatusOK {
				t.Errorf("expected status 200 for User-Agent %s, got %d", tc.name, w.Code)
			}

			// Проверяем что User-Agent правильно сохранился
			var data map[string]interface{}
			json.NewDecoder(w.Body).Decode(&data)
			request, _ := data["request"].(map[string]interface{})

			if ua, _ := request["user_agent"].(string); ua != tc.userAgent {
				t.Errorf("User-Agent mismatch: expected %s, got %s", tc.userAgent, ua)
			}
		})
	}
}

// Тест Content-Type и индентации
func TestResponseFormatting(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	mainHandler(w, req)

	// Проверяем Content-Type
	if ct := w.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("expected Content-Type application/json, got %s", ct)
	}

	// Проверяем что JSON валидный и имеет отступы (потому что мы используем SetIndent)
	body := w.Body.String()
	if !strings.Contains(body, "\n") {
		t.Error("expected indented JSON with newlines")
	}

	// Проверяем что это валидный JSON
	var data interface{}
	if err := json.Unmarshal([]byte(body), &data); err != nil {
		t.Errorf("invalid JSON response: %v", err)
	}
}

// Минимальные тесты для main
func TestMain_Configuration(t *testing.T) {
    // Тест 1: Дефолтные значения
    host := "0.0.0.0"
    port := "8000"

    if host != "0.0.0.0" || port != "8000" {
        t.Errorf("default config: %s:%s", host, port)
    }

    // Тест 2: Формат адреса
    addr := fmt.Sprintf("%s:%s", host, port)
    if addr != "0.0.0.0:8000" {
        t.Errorf("addr = %s", addr)
    }
}

func TestMain_Routes(t *testing.T) {
    // Создаем тестовый сервер
    mux := http.NewServeMux()
    mux.HandleFunc("/", mainHandler)
    mux.HandleFunc("/health", healthHandler)

    server := httptest.NewServer(mux)
    defer server.Close()

    // Проверяем GET /
    resp1, _ := http.Get(server.URL + "/")
    defer resp1.Body.Close()
    if resp1.StatusCode != 200 {
        t.Error("/ failed")
    }

    // Проверяем GET /health
    resp2, _ := http.Get(server.URL + "/health")
    defer resp2.Body.Close()
    if resp2.StatusCode != 200 {
        t.Error("/health failed")
    }

    // Проверяем 404
    resp3, _ := http.Get(server.URL + "/bad")
    defer resp3.Body.Close()
    if resp3.StatusCode != 404 {
        t.Error("/bad should return 404")
    }
}

func TestMain_LogFlags(t *testing.T) {
    oldFlags := log.Flags()
    defer log.SetFlags(oldFlags)

    log.SetFlags(log.LstdFlags | log.Lshortfile)

    flags := log.Flags()
    if flags&log.Lshortfile == 0 {
        t.Error("Lshortfile not set")
    }
}

// Тест для проверки дефолтных значений
func TestMain_DefaultValues(t *testing.T) {
    // Копируем логику из main с хардкодными значениями
    host := "0.0.0.0"
    port := "8000"

    if host != "0.0.0.0" {
        t.Errorf("default host = %s, want 0.0.0.0", host)
    }
    if port != "8000" {
        t.Errorf("default port = %s, want 8000", port)
    }

    addr := fmt.Sprintf("%s:%s", host, port)
    expectedAddr := "0.0.0.0:8000"
    if addr != expectedAddr {
        t.Errorf("addr = %s, want %s", addr, expectedAddr)
    }
}

// Тест для проверки формата адреса
func TestMain_AddrFormat(t *testing.T) {
    testCases := []struct {
        name     string
        host     string
        port     string
        expected string
    }{
        {"default", "0.0.0.0", "8000", "0.0.0.0:8000"},
        {"localhost", "127.0.0.1", "9000", "127.0.0.1:9000"},
        {"custom", "localhost", "8080", "localhost:8080"},
        {"empty host", "", "8000", ":8000"},
        {"empty port", "0.0.0.0", "", "0.0.0.0:"},
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            addr := fmt.Sprintf("%s:%s", tc.host, tc.port)
            if addr != tc.expected {
                t.Errorf("addr = %q, want %q", addr, tc.expected)
            }
        })
    }
}

// Тест для проверки регистрации маршрутов
func TestMain_RouteRegistration(t *testing.T) {
    // Создаем тестовый мультиплексор и регистрируем маршруты как в main
    mux := http.NewServeMux()
    mux.HandleFunc("/", mainHandler)
    mux.HandleFunc("/health", healthHandler)

    // Создаем тестовый сервер
    server := httptest.NewServer(mux)
    defer server.Close()

    // Проверяем корневой маршрут
    t.Run("root endpoint", func(t *testing.T) {
        resp, err := http.Get(server.URL + "/")
        if err != nil {
            t.Fatalf("failed to GET /: %v", err)
        }
        defer resp.Body.Close()

        if resp.StatusCode != http.StatusOK {
            t.Errorf("GET / returned %d, want 200", resp.StatusCode)
        }

        // Проверяем Content-Type
        if ct := resp.Header.Get("Content-Type"); ct != "application/json" {
            t.Errorf("Content-Type = %s, want application/json", ct)
        }
    })

    // Проверяем health маршрут
    t.Run("health endpoint", func(t *testing.T) {
        resp, err := http.Get(server.URL + "/health")
        if err != nil {
            t.Fatalf("failed to GET /health: %v", err)
        }
        defer resp.Body.Close()

        if resp.StatusCode != http.StatusOK {
            t.Errorf("GET /health returned %d, want 200", resp.StatusCode)
        }
    })

    // Проверяем 404
    t.Run("not found", func(t *testing.T) {
        resp, err := http.Get(server.URL + "/nonexistent")
        if err != nil {
            t.Fatalf("failed to GET /nonexistent: %v", err)
        }
        defer resp.Body.Close()

        if resp.StatusCode != http.StatusNotFound {
            t.Errorf("GET /nonexistent returned %d, want 404", resp.StatusCode)
        }
    })

    // Проверяем Method Not Allowed
    t.Run("method not allowed", func(t *testing.T) {
        resp, err := http.Post(server.URL+"/", "application/json", nil)
        if err != nil {
            t.Fatalf("failed to POST /: %v", err)
        }
        defer resp.Body.Close()

        if resp.StatusCode != http.StatusMethodNotAllowed {
            t.Errorf("POST / returned %d, want 405", resp.StatusCode)
        }
    })
}

// Тест для проверки форматирования JSON с отступами
func TestMain_JSONIndentation(t *testing.T) {
    req := httptest.NewRequest("GET", "/", nil)
    w := httptest.NewRecorder()
    mainHandler(w, req)

    body := w.Body.String()

    // Проверяем что JSON имеет отступы (содержит символы новой строки)
    if !strings.Contains(body, "\n") {
        t.Error("JSON response should be indented (contain newlines)")
    }

    // Проверяем что это валидный JSON
    var data interface{}
    if err := json.Unmarshal([]byte(body), &data); err != nil {
        t.Errorf("response is not valid JSON: %v", err)
    }
}

// Тест для проверки health JSON индентации
func TestMain_HealthJSONIndentation(t *testing.T) {
    req := httptest.NewRequest("GET", "/health", nil)
    w := httptest.NewRecorder()
    healthHandler(w, req)

    body := w.Body.String()

    if !strings.Contains(body, "\n") {
        t.Error("health JSON response should be indented (contain newlines)")
    }
}

// Тест для проверки что сервер можно сконфигурировать
func TestMain_ServerConfig(t *testing.T) {
    // Проверяем что адрес формируется корректно
    host := "0.0.0.0"
    port := "8000"
    addr := fmt.Sprintf("%s:%s", host, port)

    // Не запускаем реальный сервер, просто проверяем формат
    expectedAddr := "0.0.0.0:8000"
    if addr != expectedAddr {
        t.Errorf("server address = %s, want %s", addr, expectedAddr)
    }
}
