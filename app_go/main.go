package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"time"
)

// ==================== ЗАМЕНЯЕМЫЕ ПЕРЕМЕННЫЕ ДЛЯ ТЕСТИРОВАНИЯ ====================
var (
	osHostname         = os.Hostname
	logPrintf          = log.Printf
	logFatalf          = log.Fatalf
	httpListenAndServe = http.ListenAndServe
	osGetenv           = os.Getenv
	timeNow            = time.Now
	timeSince          = time.Since
)

// ==================== СТРУКТУРЫ ДАННЫХ ====================
type Service struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	Description string `json:"description"`
	Framework   string `json:"framework"`
}

type System struct {
	Hostname        string `json:"hostname"`
	Platform        string `json:"platform"`
	PlatformVersion string `json:"platform_version"`
	Architecture    string `json:"architecture"`
	CPUCount        int    `json:"cpu_count"`
	GoVersion       string `json:"go_version"`
}

type HealthResp struct {
	Status        string `json:"status"`
	Timestamp     string `json:"timestamp"`
	UptimeSeconds int    `json:"uptime_seconds"`
}

type Runtime struct {
	UptimeSeconds int    `json:"uptime_seconds"`
	UptimeHuman   string `json:"uptime_human"`
	CurrentTime   string `json:"current_time"`
	Timezone      string `json:"timezone"`
}

type Request struct {
	ClientIP  string `json:"client_ip"`
	UserAgent string `json:"user_agent"`
	Method    string `json:"method"`
	Path      string `json:"path"`
}

type Endpoint struct {
	Path        string `json:"path"`
	Method      string `json:"method"`
	Description string `json:"description"`
}

type ServiceInfo struct {
	Service   Service    `json:"service"`
	System    System     `json:"system"`
	Runtime   Runtime    `json:"runtime"`
	Request   Request    `json:"request"`
	Endpoints []Endpoint `json:"endpoints"`
}

// ==================== HELPER FUNCTIONS ====================
func getHostname() string {
	hostname, err := osHostname()
	if err != nil {
		logPrintf("Error getting hostname: %v", err)
		return "unknown"
	}
	return hostname
}

func getClientIP(r *http.Request) string {
	forwarded := r.Header.Get("X-Forwarded-For")
	if forwarded != "" {
		return forwarded
	}
	return r.RemoteAddr
}

// Application start time
var startTime = timeNow()

func getUptime() (int, string) {
	elapsed := timeSince(startTime)
	seconds := int(elapsed.Seconds())
	hours := seconds / 3600
	minutes := (seconds % 3600) / 60
	return seconds, fmt.Sprintf("%d hours, %d minutes", hours, minutes)
}

// ==================== HANDLERS ====================
func mainHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{
			"error":   "Method Not Allowed",
			"message": "Only GET method is allowed for this endpoint",
		})
		return
	}

	logPrintf("Request: %s %s from %s", r.Method, r.URL.Path, getClientIP(r))

	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	uptimeSeconds, uptimeHuman := getUptime()
	location, _ := timeNow().Local().Zone()

	info := ServiceInfo{
		Service: Service{
			Name:        "devops-info-service",
			Version:     "1.0.0",
			Description: "DevOps course info service",
			Framework:   "Go net/http",
		},
		System: System{
			Hostname:        getHostname(),
			Platform:        runtime.GOOS,
			PlatformVersion: runtime.Version(),
			Architecture:    runtime.GOARCH,
			CPUCount:        runtime.NumCPU(),
			GoVersion:       runtime.Version(),
		},
		Runtime: Runtime{
			UptimeSeconds: uptimeSeconds,
			UptimeHuman:   uptimeHuman,
			CurrentTime:   timeNow().Format(time.RFC3339),
			Timezone:      location,
		},
		Request: Request{
			ClientIP:  getClientIP(r),
			UserAgent: r.UserAgent(),
			Method:    r.Method,
			Path:      r.URL.Path,
		},
		Endpoints: []Endpoint{
			{Path: "/", Method: "GET", Description: "Service information"},
			{Path: "/health", Method: "GET", Description: "Health check"},
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(info); err != nil {
		logPrintf("Error encoding JSON: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{
			"error":   "Method Not Allowed",
			"message": "Only GET method is allowed for this endpoint",
		})
		return
	}

	logPrintf("Health check from %s", getClientIP(r))

	uptimeSeconds, _ := getUptime()

	health := HealthResp{
		Status:        "healthy",
		Timestamp:     timeNow().UTC().Format(time.RFC3339),
		UptimeSeconds: uptimeSeconds,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(health); err != nil {
		logPrintf("Error encoding JSON: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	logPrintf("404 Not Found: %s %s", r.Method, r.URL.Path)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)

	response := map[string]string{
		"error":   "Not Found",
		"message": "Endpoint does not exist",
	}

	json.NewEncoder(w).Encode(response)
}

// ==================== SERVER ====================
func run() error {
	host := osGetenv("HOST")
	if host == "" {
		host = "0.0.0.0"
	}

	port := osGetenv("PORT")
	if port == "" {
		port = "8000"
	}

	log.SetFlags(log.LstdFlags | log.Lshortfile)
	logPrintf("Starting DevOps Info Service (Go) on %s:%s", host, port)

	http.HandleFunc("/", mainHandler)
	http.HandleFunc("/health", healthHandler)

	addr := fmt.Sprintf("%s:%s", host, port)
	logPrintf("Server is running on http://%s", addr)
	logPrintf("Press Ctrl+C to stop")

	return httpListenAndServe(addr, nil)
}

func main() {
	if err := run(); err != nil {
		logFatalf("Server failed to start: %v", err)
	}
}