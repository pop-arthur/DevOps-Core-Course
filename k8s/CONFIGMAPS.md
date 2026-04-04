# CONFIGMAPS & PERSISTENCE IMPLEMENTATION

## 1. Application Changes

### Visits Counter Implementation

The application was extended with a persistent visits counter:

* A file-based counter is stored at `/app/data/visits`
* On each request to `/`, the counter:

  1. Reads current value from file
  2. Increments it
  3. Writes updated value back

This ensures state is preserved across restarts when using persistent storage.

### New Endpoint

A new endpoint was added:

```
GET /visits
```

Returns current visits count:

```json
{"visits": 4}
```

### Local Testing with Docker

* Volume mounted:

```
./data:/app/data
```

* Verified:

  * Counter increases with requests
  * Value persists after container restart

```
(.venv) arthur@192 DevOps-Core-Course %  curl localhost:8000/      

{"service":{"name":"devops-info-service","version":"1.0.0","description":"DevOps course info service","framework":"FastAPI"},"system":{"hostname":"fda9ebf82c54","platform":"Linux","platform_version":"Linux-6.10.14-linuxkit-aarch64-with-glibc2.41","architecture":"aarch64","cpu_count":8,"python_version":"3.13.12"},"runtime":{"uptime_seconds":82,"uptime_human":"0 hours, 1 minutes","current_time":"2026-04-04T11:49:51.372002+00:00","timezone":"UTC"},"request":{"client_ip":"192.168.65.1","user_agent":"curl/8.7.1","method":"GET","path":"/"},"endpoints":[{"path":"/","method":"GET","description":"Service information"},{"path":"/health","method":"GET","description":"Health check"}]}%                       
(.venv) arthur@192 DevOps-Core-Course %  curl localhost:8000/

{"service":{"name":"devops-info-service","version":"1.0.0","description":"DevOps course info service","framework":"FastAPI"},"system":{"hostname":"fda9ebf82c54","platform":"Linux","platform_version":"Linux-6.10.14-linuxkit-aarch64-with-glibc2.41","architecture":"aarch64","cpu_count":8,"python_version":"3.13.12"},"runtime":{"uptime_seconds":38,"uptime_human":"0 hours, 0 minutes","current_time":"2026-04-04T11:52:57.815567+00:00","timezone":"UTC"},"request":{"client_ip":"192.168.65.1","user_agent":"curl/8.7.1","method":"GET","path":"/"},"endpoints":[{"path":"/","method":"GET","description":"Service information"},{"path":"/health","method":"GET","description":"Health check"}]}%                       
(.venv) arthur@192 DevOps-Core-Course %  curl localhost:8000/visits

{"visits":3}%                      
```
---

## 2. ConfigMap Implementation

### Config File

`files/config.json`:

```json
{
  "appName": "devops-info-service",
  "environment": "dev",
  "featureFlags": {
    "enableVisits": true
  }
}
```

---

### ConfigMap Template

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: {{ include "myapp.fullname" . }}-config
data:
  config.json: |-
{{ .Files.Get "files/config.json" | indent 4 }}
```

---

### Mounting ConfigMap as File

Deployment configuration:

```yaml
volumes:
  - name: config-volume
    configMap:
      name: {{ include "myapp.fullname" . }}-config

volumeMounts:
  - name: config-volume
    mountPath: /config
```

---

### Verification (inside pod)

```
(.venv) arthur@192 myapp % kubectl exec -it myapp-myapp-59967b5fbd-6qr5q -- cat /config/config.json 
{ 
    "appName": "devops-info-service", 
    "environment": "dev", 
    "featureFlags": { 
    "enableVisits": true } 
}%
```

---

### ConfigMap as Environment Variables

Second ConfigMap:

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: {{ include "myapp.fullname" . }}-env
data:
  APP_NAME: "devops-info-service"
  APP_ENV: "dev"
  FEATURE_VISITS: "true"
```

Deployment:

```yaml
envFrom:
  - configMapRef:
      name: {{ include "myapp.fullname" . }}-env
```

---

### Verification (env vars)

```
(.venv) arthur@192 myapp % kubectl exec -it myapp-myapp-7cc5b75b95-n5pv6 -- printenv | grep APP_
APP_ENV=dev
APP_NAME=devops-info-service
```

---

## 3. Persistent Volume

### PVC Configuration

```yaml
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: {{ include "myapp.fullname" . }}-data
spec:
  accessModes:
    - ReadWriteOnce
  resources:
    requests:
      storage: 100Mi
```

---

### Explanation

* **ReadWriteOnce (RWO)**:

  * Volume can be mounted by a single node
  * Suitable for this application

* **Storage Class**:

  * Default `standard` (Minikube)
  * Automatically provisions storage

---

### Volume Mount

```yaml
volumes:
  - name: data-volume
    persistentVolumeClaim:
      claimName: {{ include "myapp.fullname" . }}-data

volumeMounts:
  - name: data-volume
    mountPath: /app/data
```

---

### PVC Verification
```
(.venv) arthur@192 myapp % kubectl get pvc 
NAME STATUS VOLUME CAPACITY ACCESS MODES STORAGECLASS VOLUMEATTRIBUTESCLASS AGE 
myapp-myapp-data Bound pvc-80043a74-74d7-40d9-a044-de2a26c89c0b 100Mi RWO standard <unset> 13s
```

---

### Persistence Test

```
(.venv) arthur@192 DevOps-Core-Course % curl localhost:8080/
curl localhost:8080/visits
{"service":{"name":"devops-info-service","version":"1.0.0","description":"DevOps course info service","framework":"FastAPI"},"system":{"hostname":"myapp-myapp-795c48874f-r42dn","platform":"Linux","platform_version":"Linux-6.10.14-linuxkit-aarch64-with-glibc2.41","architecture":"aarch64","cpu_count":8,"python_version":"3.13.12"},"runtime":{"uptime_seconds":307,"uptime_human":"0 hours, 5 minutes","current_time":"2026-04-04T15:55:47.388051+00:00","timezone":"UTC"},"request":{"client_ip":"127.0.0.1","user_agent":"curl/8.7.1","method":"GET","path":"/"},"endpoints":[{"path":"/","method":"GET","description":"Service information"},{"path":"/health","method":"GET","description":"Health check"}]}

{"visits":1}%
                                                                                                                                          
(.venv) arthur@192 DevOps-Core-Course % curl localhost:8080/
curl localhost:8080/visits
{"service":{"name":"devops-info-service","version":"1.0.0","description":"DevOps course info service","framework":"FastAPI"},"system":{"hostname":"myapp-myapp-795c48874f-r42dn","platform":"Linux","platform_version":"Linux-6.10.14-linuxkit-aarch64-with-glibc2.41","architecture":"aarch64","cpu_count":8,"python_version":"3.13.12"},"runtime":{"uptime_seconds":310,"uptime_human":"0 hours, 5 minutes","current_time":"2026-04-04T15:55:50.874110+00:00","timezone":"UTC"},"request":{"client_ip":"127.0.0.1","user_agent":"curl/8.7.1","method":"GET","path":"/"},"endpoints":[{"path":"/","method":"GET","description":"Service information"},{"path":"/health","method":"GET","description":"Health check"}]}

{"visits":2}%
                                                                                                                                          
(.venv) arthur@192 DevOps-Core-Course % curl localhost:8080/
curl localhost:8080/visits
{"service":{"name":"devops-info-service","version":"1.0.0","description":"DevOps course info service","framework":"FastAPI"},"system":{"hostname":"myapp-myapp-795c48874f-r42dn","platform":"Linux","platform_version":"Linux-6.10.14-linuxkit-aarch64-with-glibc2.41","architecture":"aarch64","cpu_count":8,"python_version":"3.13.12"},"runtime":{"uptime_seconds":312,"uptime_human":"0 hours, 5 minutes","current_time":"2026-04-04T15:55:52.290116+00:00","timezone":"UTC"},"request":{"client_ip":"127.0.0.1","user_agent":"curl/8.7.1","method":"GET","path":"/"},"endpoints":[{"path":"/","method":"GET","description":"Service information"},{"path":"/health","method":"GET","description":"Health check"}]}

{"visits":3}%
   
(.venv) arthur@192 DevOps-Core-Course % kubectl get pods                                  
NAME                                    READY   STATUS    RESTARTS       AGE
devops-info-service-698cc89c85-4jqc2    1/1     Running   2 (165m ago)   10d
devops-info-service-698cc89c85-gp6kq    1/1     Running   2 (165m ago)   10d
devops-info-service-698cc89c85-zctbl    1/1     Running   3 (165m ago)   10d
myapp-myapp-795c48874f-9zlch            1/1     Running   0              8m59s
myapp-myapp-795c48874f-r42dn            1/1     Running   0              8m59s
vault-0                                 1/1     Running   1 (165m ago)   3d9h
vault-agent-injector-75998c9b76-rjq52   1/1     Running   1 (165m ago)   3d9h
(.venv) arthur@192 DevOps-Core-Course % kubectl delete pod myapp-myapp-795c48874f-9zlch 
pod "myapp-myapp-795c48874f-9zlch" deleted from default namespace
^[[A%                                                                                                                                        (.venv) arthur@192 DevOps-Core-Course % kubectl delete pod myapp-myapp-795c48874f-r42dn     
pod "myapp-myapp-795c48874f-r42dn" deleted from default namespace
(.venv) arthur@192 DevOps-Core-Course % kubectl get pods                                    
NAME                                    READY   STATUS    RESTARTS       AGE
devops-info-service-698cc89c85-4jqc2    1/1     Running   2 (166m ago)   10d
devops-info-service-698cc89c85-gp6kq    1/1     Running   2 (166m ago)   10d
devops-info-service-698cc89c85-zctbl    1/1     Running   3 (166m ago)   10d
myapp-myapp-795c48874f-clzhc            1/1     Running   0              32s
myapp-myapp-795c48874f-vv2dh            1/1     Running   0              62s
vault-0                                 1/1     Running   1 (166m ago)   3d9h
vault-agent-injector-75998c9b76-rjq52   1/1     Running   1 (166m ago)   3d9h
(.venv) arthur@192 DevOps-Core-Course % curl localhost:8080/                                
curl localhost:8080/visits
{"service":{"name":"devops-info-service","version":"1.0.0","description":"DevOps course info service","framework":"FastAPI"},"system":{"hostname":"myapp-myapp-795c48874f-vv2dh","platform":"Linux","platform_version":"Linux-6.10.14-linuxkit-aarch64-with-glibc2.41","architecture":"aarch64","cpu_count":8,"python_version":"3.13.12"},"runtime":{"uptime_seconds":68,"uptime_human":"0 hours, 1 minutes","current_time":"2026-04-04T16:01:00.985548+00:00","timezone":"UTC"},"request":{"client_ip":"127.0.0.1","user_agent":"curl/8.7.1","method":"GET","path":"/"},"endpoints":[{"path":"/","method":"GET","description":"Service information"},{"path":"/health","method":"GET","description":"Health check"}]}{"visits":4}%                                                                                                                                           (.venv) arthur@192 DevOps-Core-Course % 
   
```

### Persistence test

* Before restart: `visits = 3`
* After restart and call: `visits = 4`

---

### Result

- Data persisted after pod deletion
- PVC successfully retains application state

---

## 4. ConfigMap vs Secret

### ConfigMap

Used for:

* Non-sensitive configuration
* Application settings
* Feature flags

Example:

* appName
* environment
* feature toggles

---

### Secret

Used for:

* Sensitive data
* Credentials
* API keys

Example:

* passwords
* tokens
* database credentials

---

### Key Differences

| Feature   | ConfigMap  | Secret                  |
| --------- | ---------- | ----------------------- |
| Data type | Plain text | Base64 encoded          |
| Security  | Not secure | Secure (RBAC protected) |
| Use case  | Config     | Sensitive data          |


---
