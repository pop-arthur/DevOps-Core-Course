# Lab 01: Implementing in Go

Bonus task: Creating a compiled version of our DevOps service in Go.

## Why Go Instead of Python?

I chose Go for this implementation to see how a compiled language compares to our Python version. The main benefit? **One file that runs anywhere**.

With Go, I compile everything into a single binary. No need to install Python, set up virtual environments, or manage pip packages. Just copy the file to any server and run it.

## No Framework Needed

- Service is simple (just 2 endpoints)
- Zero dependencies means easier maintenance

## What I Did Better This Time

**Cleaner Structure:**
```go
type ServiceInfo struct {
    Service   Service   `json:"service"`
    System    System    `json:"system"`
    // ...
}
```

**Production-Ready Logging:**
```go
log.SetFlags(log.LstdFlags | log.Lshortfile)
log.Printf("Request from %s", ip)
```
Added file and line numbers to logs for easier debugging.

## The Endpoints

Both endpoints work exactly like the Python version:

**Main endpoint (`GET /`):**
```bash
curl http://localhost:8000/
```
Returns all service info, now showing `"framework": "Go net/http"`.

**Health check (`GET /health`):**
```bash
curl http://localhost:8000/health
```
Same JSON response: status, timestamp, and uptime.

## Building and Running

**Development mode (like Python):**
```bash
go run main.go
PORT=8080 go run main.go
```

**Production build:**
```bash
# Creates a single executable
go build -o devops-service main.go
./devops-service
```

## The Size Difference

Here's what surprised me:

```
Go binary:      5.2 MB  (as shown on screenshot)
Python package: ~25 MB  (.venv + file)
```
