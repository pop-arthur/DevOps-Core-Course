# Lab 7
## Documentation

## Architecture 

```
┌─────────────────────────────────────────────────────┐
│                    DOCKER HOST                       │
│                                                       │
│  ┌──────────────┐      ┌──────────────┐             │
│  │   Python App  │      │   Promtail   │             │
│  │   port:8000   │─────▶│  Discovers   │             │
│  │label:logging= │      │  containers  │             │
│  │   promtail    │      │ via Docker   │             │
│  └──────────────┘      │    socket    │             │
│         │               └──────┬───────┘             │
│         │                      │                      │
│         │                      │ http://loki:3100     │
│         ▼                      ▼                      │
│  ┌─────────────────────────────────────────────┐     │
│  │                   Loki                        │     │
│  │         Log storage with TSDB index           │     │
│  │          Filesystem storage, 7-day retention  │     │
│  └─────────────────────┬─────────────────────────┘     │
│                        │                                │
│                        │ http://loki:3100 (query)       │
│                        ▼                                │
│  ┌─────────────────────────────────────────────┐     │
│  │                  Grafana                      │     │
│  │         Port: 3002, Login required            │     │
│  │         Dashboards with LogQL queries         │     │
│  └─────────────────────────────────────────────┘     │
│                                                       │
│  Volumes: loki-data, grafana-data, promtail-positions│
│  Network: logging (all containers communicate)        │
└─────────────────────────────────────────────────────┘
```

## Setup Guide

### 1. Create Project Structure
```bash
mkdir -p monitoring/{loki,promtail,grafana/datasources}
cd monitoring
```

### 2. Create .env File (Don't commit this!)
```bash
cat > .env << EOF
GRAFANA_ADMIN_USER=admin
GRAFANA_ADMIN_PASSWORD=YourSecurePassword123!
EOF
chmod 600 .env
echo ".env" >> .gitignore
```

### 3. Deploy
```bash
docker compose up -d
docker compose ps
```

### 4. Add Loki Data Source in Grafana
- Open `http://localhost:3002`
- Login with credentials from `.env`
- Go to **Connections** → **Data sources** → **Add data source** → **Loki**
- URL: `http://loki:3100`
- Click **Save & Test**


## Config

### Loki Config (`loki/config.yml`)
```yaml
auth_enabled: false
server:
  http_listen_port: 3100
common:
  path_prefix: /loki
  storage:
    filesystem:
      chunks_directory: /loki/chunks
      rules_directory: /loki/rules
schema_config:
  configs:
    - from: 2024-01-01
      store: tsdb          # Faster queries
      object_store: filesystem
      schema: v13
      index:
        prefix: index_
        period: 24h
limits_config:
  retention_period: 168h   # 7 days
compactor:
  retention_enabled: true
```

### Promtail Config (`promtail/config.yml`)
```yaml
server:
  http_listen_port: 9080
positions:
  filename: /tmp/positions.yaml
clients:
  - url: http://loki:3100/loki/api/v1/push
scrape_configs:
  - job_name: docker
    docker_sd_configs:
      - host: unix:///var/run/docker.sock
        refresh_interval: 5s
        filters:
          - name: label
            values: ["logging=promtail"]  # Only collect from labeled containers
    relabel_configs:
      - source_labels: ['__meta_docker_container_name']
        target_label: container
        regex: '/(.*)'
      - source_labels: ['__meta_docker_container_label_app']
        target_label: app
```

### Docker Compose (`docker-compose.yml`)
```yaml
services:
  loki:
    image: grafana/loki:3.0.0
    ports: ["3100:3100"]
    volumes:
      - ./loki/config.yml:/etc/loki/config.yml:ro
      - loki-data:/loki
    networks: [logging]
    deploy:
      resources:
        limits:
          cpus: "1.0"
          memory: 1024M
    healthcheck:
      test: ["CMD-SHELL", "wget -qO- http://localhost:3100/ready || exit 1"]

  promtail:
    image: grafana/promtail:3.0.0
    volumes:
      - ./promtail/config.yml:/etc/promtail/config.yml:ro
      - /var/lib/docker/containers:/var/lib/docker/containers:ro
      - /var/run/docker.sock:/var/run/docker.sock:ro
    networks: [logging]
    depends_on: [loki]

  grafana:
    image: grafana/grafana:12.3.1
    ports: ["3002:3002"]
    environment:
      - GF_SERVER_HTTP_PORT=3002
      - GF_AUTH_ANONYMOUS_ENABLED=false
      - GF_SECURITY_ADMIN_USER=${GRAFANA_ADMIN_USER}
      - GF_SECURITY_ADMIN_PASSWORD=${GRAFANA_ADMIN_PASSWORD}
    networks: [logging]
    depends_on: [loki]
    env_file:
      - .env
    deploy:
      resources:
        limits:
          cpus: "0.5"
          memory: 512M

  devops-info-service:
    image: poparthur/devops-info-service:latest
    ports: ["8000:5000"]
    labels:
      logging: "promtail"    # Required for Promtail to discover
      app: "devops-python"
    networks: [logging]

volumes:
  loki-data:
  grafana-data:
networks:
  logging:
```

## JSON Logging in Python App

Add to your Python app:

```python
import logging
from pythonjsonlogger import jsonlogger

# Configure JSON logging
logger = logging.getLogger()
logHandler = logging.StreamHandler()
formatter = jsonlogger.JsonFormatter('%(timestamp)s %(level)s %(message)s')
logHandler.setFormatter(formatter)
logger.addHandler(logHandler)
logger.setLevel(logging.INFO)

# Log example
logger.info('Request completed', extra={
    'method': 'GET',
    'path': '/',
    'status': 200
})
```

Example log output:
```json
{"timestamp": "2026-03-09 15:30:45", "level": "INFO", "message": "Request completed", "method": "GET", "status": 200}
```

## Dashboard Panels

![6 dashboards.png](screenshots/6%20dashboards.png)

| Panel | Query | Visualization |
|-------|-------|---------------|
| Recent Logs | `{app="devops-python"}` | Logs table |
| Request Rate | `rate({app="devops-python"}[1m])` | Time series graph |
| Errors Only | `{app="devops-python"} |= "ERROR"` | Logs with highlighting |
| Log Levels | `sum by (level) (count_over_time({app="devops-python"} \| json [1h]))` | Pie chart |

## Production

| Feature | How We Did It |
|---------|---------------|
| **No anonymous access** | `GF_AUTH_ANONYMOUS_ENABLED=false` |
| **Secure passwords** | Credentials in `.env` file (not in code) |
| **Resource limits** | CPU and memory limits on all services |
| **Health checks** | Each service checks if it's working |
| **7-day retention** | `retention_period: 168h` in Loki config |
| **Persistence** | Named volumes for Loki and Grafana data |


## Testing Commandss

```bash
```bash
# Generate traffic
for i in {1..20}; do curl http://localhost:8000/; done
for i in {1..20}; do curl http://localhost:8000/health; done
for i in {1..5}; do curl http://localhost:8000/error; done

# Check services
curl http://localhost:3100/ready                    # Should return "ready"
curl http://localhost:9080/targets | jq .           # Should show your app
curl -u admin:YourPassword! http://localhost:3002/api/health  # Test auth
```


## Common Problems & Solutions

| Problem | Solution |
|---------|----------|
| No logs in Grafana | Check container has label `logging=promtail` |
| Port 3000 already in use | Changed Grafana to port 3002 |

## Evidence

![1 job docker.png](screenshots/1%20job%20docker.png)
![2 query.png](screenshots/2%20query.png)
![3 errors.png](screenshots/3%20errors.png)
![4 query.png](screenshots/4%20query.png)
![5 logs from app.png](screenshots/5%20logs%20from%20app.png)
![6 dashboards.png](screenshots/6%20dashboards.png)
![7 docker ps.png](screenshots/7%20docker%20ps.png)
![8 grafana auth.png](screenshots/8%20grafana%20auth.png)
![9 grafana auth.png](screenshots/9%20grafana%20auth.png)
