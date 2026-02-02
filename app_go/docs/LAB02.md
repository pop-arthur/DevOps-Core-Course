# Multi-Stage Build Strategy Explanation

## **Why Multi-Stage Builds?**

For compiled languages like Go, we need two different environments:
1. **Builder environment:** Has Go compiler, build tools, dependencies (large)
2. **Runtime environment:** Only needs the compiled binary (tiny)

## **Two-Stage Strategy**

### **Stage 1: Builder**
```dockerfile
FROM golang:1.22-alpine AS builder
WORKDIR /build
COPY go.mod .
RUN go mod download
COPY main.go .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o devops-info-service .
```

**Purpose:** Create the executable binary
- Uses `golang:1.22-alpine` (300MB) for compilation
- Downloads dependencies separately (better caching)
- Compiles with optimizations:
  - `CGO_ENABLED=0`: Pure Go, no C dependencies
  - `-ldflags="-w -s"`: Strips debug symbols, reduces size
- Produces a self-contained Linux binary

### **Stage 2: Runtime**
```dockerfile
FROM gcr.io/distroless/static:nonroot
WORKDIR /app
COPY --from=builder /build/devops-info-service .
EXPOSE 8000
USER nonroot
ENTRYPOINT ["/app/devops-info-service"]
```

**Purpose:** Run the compiled application
- Uses `distroless/static:nonroot` (2MB) - minimal secure base
- Copies ONLY the compiled binary from builder
- Runs as non-root user for security
- Contains nothing but the binary and essential runtime libraries

## **Size Comparison & Analysis**

### **Builder Stage Size:** ~300MB
Contains:
- Go compiler and tools
- Alpine Linux with build utilities
- Source code and dependencies
- Intermediate build artifacts

### **Final Image Size:** ~7MB  
Contains:
- Compiled Go binary (~6.8MB)
- Distroless base image (~200KB)

### **Why this matters:**
- **Faster deployment:** Smaller images transfer quicker
- **Reduced storage:** 7MB vs 300MB per instance
- **Lower costs:** Less storage and bandwidth usage

## **Terminal Output**
```bash
docker build -t devops-info-service-go .       
[+] Building 1.1s (15/15) FINISHED                                                                                      docker:desktop-linux
 => [internal] load build definition from Dockerfile                                                                                    0.0s
 => => transferring dockerfile: 472B                                                                                                    0.0s
 => [internal] load metadata for gcr.io/distroless/static:nonroot                                                                       0.7s
 => [internal] load metadata for docker.io/library/golang:1.22-alpine                                                                   1.0s
 => [internal] load .dockerignore                                                                                                       0.0s
 => => transferring context: 170B                                                                                                       0.0s
 => [builder 1/6] FROM docker.io/library/golang:1.22-alpine@sha256:1699c10032ca2582ec89a24a1312d986a3f094aed3d5c1147b19880afe40e052     0.0s
 => [stage-1 1/3] FROM gcr.io/distroless/static:nonroot@sha256:cba10d7abd3e203428e86f5b2d7fd5eb7d8987c387864ae4996cf97191b33764         0.0s
 => [internal] load build context                                                                                                       0.0s
 => => transferring context: 91B                                                                                                        0.0s
 => CACHED [stage-1 2/3] WORKDIR /app                                                                                                   0.0s
 => CACHED [builder 2/6] WORKDIR /build                                                                                                 0.0s
 => CACHED [builder 3/6] COPY go.mod .                                                                                                  0.0s
 => CACHED [builder 4/6] RUN go mod download                                                                                            0.0s
 => CACHED [builder 5/6] COPY main.go .                                                                                                 0.0s
 => CACHED [builder 6/6] RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o devops-info-service .                   0.0s
 => CACHED [stage-1 3/3] COPY --from=builder /build/devops-info-service .                                                               0.0s
 => exporting to image                                                                                                                  0.0s
 => => exporting layers                                                                                                                 0.0s
 => => writing image sha256:70a6dc5efe1e690d4e10058ed7eb171751c53029d6a7763f55d1aab3aed12f3d                                            0.0s
 => => naming to docker.io/library/devops-info-service-go                                                                               0.0s

arthur@Artur-MacBook-Pro app_go % docker run -p 8000:8000 devops-info-service-go

2026/02/02 11:18:39 main.go:205: Starting DevOps Info Service (Go) on 0.0.0.0:8000
2026/02/02 11:18:39 main.go:213: Server is running on http://0.0.0.0:8000
2026/02/02 11:18:39 main.go:214: Press Ctrl+C to stop
2026/02/02 11:19:12 main.go:98: Request: GET /healt from 192.168.65.1:54644
2026/02/02 11:19:14 main.go:155: Health check from 192.168.65.1:54988


arthur@Artur-MacBook-Pro DevOps-Core-Course % curl http://localhost:8000/health

{
  "status": "healthy",
  "timestamp": "2026-02-02T11:19:14Z",
  "uptime_seconds": 35
}
```

## **Why Multi-Stage Matters for Compiled Languages**

### **1. Security Isolation**
- Builder stage can have vulnerabilities in compilers
- Final stage has only the minimal runtime
- Attackers can't exploit build tools that don't exist

### **2. Optimized Caching**
```dockerfile
COPY go.mod .
RUN go mod download    # Cached unless go.mod changes
COPY main.go .
RUN go build          # Cached unless code changes
```
Dependencies download only when `go.mod` changes, not on every code change.

### **3. Platform Flexibility**
- Build once, run anywhere architecture
- Builder can target different OS/arch
- Final image stays platform-agnostic

### **4. Production vs Development**
**Development needs:**
- Debuggers
- Compilers
- Source code
- Build tools

**Production needs:**
- Executable binary
- Runtime libraries
- Configuration

Multi-stage builds give us both in one Dockerfile.

## **Technical Breakdown of Each Stage**

### **Builder Stage Details:**
```dockerfile
FROM golang:1.22-alpine AS builder   # Lightest Go image available
WORKDIR /build                       # Isolated build directory
COPY go.mod .                        # Copy dependency manifest
RUN go mod download                  # Cache dependencies layer
COPY main.go .                       # Copy source code
RUN CGO_ENABLED=0 \                  # Disable C bindings
    GOOS=linux \                     # Target Linux
    GOARCH=amd64 \                   # Target x86-64
    go build \                       # Compile
    -ldflags="-w -s" \               # Strip debug info
    -o devops-info-service .         # Output name
```

**Key optimization:** `CGO_ENABLED=0` creates a static binary with no external dependencies.

### **Runtime Stage Details:**
```dockerfile
FROM gcr.io/distroless/static:nonroot  # Google's minimal secure image
WORKDIR /app                           # Application directory
COPY --from=builder /build/devops-info-service .  # Only binary
EXPOSE 8000                             # Document port
USER nonroot                            # Non-root execution
ENTRYPOINT ["/app/devops-info-service"] # Binary as entrypoint
```

**Benefits:**
- No shell (`/bin/bash`), no package manager
- No unnecessary libraries
- Non-root user by default
- Regularly security-scanned
