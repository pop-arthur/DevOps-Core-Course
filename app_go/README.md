# DevOps Info Service (Go Version)

A Go-based web service that provides system information, runtime metrics, and health monitoring. Built as part of the
DevOps course.

## Overview

This web service offers endpoints that provide detailed information about the service and the system it runs on,
including:

- Runtime metrics
- Request details
- Health status
- System information

**Key Advantage:** Single binary deployment - no dependencies, no Python/Java runtime required.

## Prerequisites

- Go 1.21+ (or just download the binary)

## Installation & Running

```bash
# Run directly
go run main.go

# Or build and run
go build -o devops-service
./devops-service
```

## Running the Application

**Default (0.0.0.0:8000):**

```bash
./devops-service
```

**Custom configuration:**

```bash
PORT=8080 ./devops-service
HOST=127.0.0.1 PORT=3000 ./devops-service
```

## API Endpoints

### `GET /`

Returns comprehensive service and system information.

```bash
curl http://localhost:8000/
```

**Response includes:**

- Service metadata (name, version, framework)
- System info (hostname, platform, CPU count, Go version)
- Runtime metrics (uptime, current time, timezone)
- Request details (client IP, user agent)
- Available endpoints

### `GET /health`

Health check endpoint for monitoring tools.

```bash
curl http://localhost:8000/health
```

**Response:**

```json
{
  "status": "healthy",
  "timestamp": "2026-01-26T07:32:24.854178Z",
  "uptime_seconds": 42
}
```

## Configuration

Environment variables:

| Variable | Default   | Description         |
|----------|-----------|---------------------|
| `HOST`   | `0.0.0.0` | Server bind address |
| `PORT`   | `8000`    | Server port number  |

## Testing

```bash
# Test endpoints
curl http://localhost:8000/
curl http://localhost:8000/health

# Pretty-print JSON
curl -s http://localhost:8000/ | python3 -m json.tool
```

## Key Benefits Over Python Version

**No dependencies** - Single binary, no Python/pip installation needed  
**Fast startup** - No virtual environment setup  
**Easy deployment** - Copy binary to any server and run  
**Better performance** - Compiled language, smaller memory footprint  
**Cross-compile** - Build for Linux/Windows/Mac from any machine

## Troubleshooting

**Port already in use:**

```bash
PORT=8080 ./devops-service
```

**Permission denied:**

```bash
chmod +x devops-service
```

## Project Structure

```
app-go/
├── main.go              # Main application (single file!)
├── README.md           # This file
├── go.mod              # Go module definition
└── docs/               # Documentation
    ├── LAB01.md       # Lab report
    ├── GO.md       # Language justification
    └── screenshots/   # Evidence
```
