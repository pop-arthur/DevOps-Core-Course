# K8s lab

### Task 1 — Local Kubernetes Setup

- Terminal output showing successful cluster setup

```
(.venv) arthur@Artur-MacBook-Pro app_python % minikube version
minikube version: v1.35.0
commit: dd5d320e41b5451cdf3c01891bc4e13d189586ed-dirty

(.venv) arthur@Artur-MacBook-Pro app_python % minikube start
😄  minikube v1.35.0 on Darwin 26.1 (arm64)
🎉  minikube 1.38.1 is available! Download it: https://github.com/kubernetes/minikube/releases/tag/v1.38.1
💡  To disable this notice, run: 'minikube config set WantUpdateNotification false'

✨  Using the docker driver based on existing profile
👍  Starting "minikube" primary control-plane node in "minikube" cluster
🚜  Pulling base image v0.0.46 ...
❗  minikube cannot pull kicbase image from any docker registry, and is trying to download kicbase tarball from github release page via HTTP.
❗  It's very likely that you have an internet issue. Please ensure that you can access the internet at least via HTTP, directly or with proxy. Currently your proxy configure is:

    > kicbase-v0.0.46-arm64.tar:  1.14 GiB / 1.14 GiB  100.00% 3.63 MiB p/s 5m2
🤷  docker "minikube" container is missing, will recreate.
🔥  Creating docker container (CPUs=2, Memory=2048MB) ...

🧯  Docker is nearly out of disk space, which may cause deployments to fail! (96% of capacity). You can pass '--force' to skip this check.
💡  Suggestion: 

    Try one or more of the following to free up space on the device:
    
    1. Run "docker system prune" to remove unused Docker data (optionally with "-a")
    2. Increase the storage allocated to Docker for Desktop by clicking on:
    Docker icon > Preferences > Resources > Disk Image Size
    3. Run "minikube ssh -- docker system prune" if using the Docker container runtime
🍿  Related issue: https://github.com/kubernetes/minikube/issues/9024


🧯  Docker is nearly out of disk space, which may cause deployments to fail! (96% of capacity). You can pass '--force' to skip this check.
💡  Suggestion: 

    Try one or more of the following to free up space on the device:
    
    1. Run "docker system prune" to remove unused Docker data (optionally with "-a")
    2. Increase the storage allocated to Docker for Desktop by clicking on:
    Docker icon > Preferences > Resources > Disk Image Size
    3. Run "minikube ssh -- docker system prune" if using the Docker container runtime
🍿  Related issue: https://github.com/kubernetes/minikube/issues/9024

🐳  Preparing Kubernetes v1.32.0 on Docker 27.4.1 ...
    ▪ Generating certificates and keys ...
    ▪ Booting up control plane ...
    ▪ Configuring RBAC rules ...
🔗  Configuring bridge CNI (Container Networking Interface) ...
🔎  Verifying Kubernetes components...
    ▪ Using image gcr.io/k8s-minikube/storage-provisioner:v5
    ▪ Using image docker.io/kubernetesui/metrics-scraper:v1.0.8
    ▪ Using image docker.io/kubernetesui/dashboard:v2.7.0
💡  Some dashboard features require the metrics-server addon. To enable all features please run:

        minikube addons enable metrics-server

🌟  Enabled addons: storage-provisioner, default-storageclass, dashboard

❗  /opt/homebrew/bin/kubectl is version 1.35.2, which may have incompatibilities with Kubernetes 1.32.0.
    ▪ Want kubectl v1.32.0? Try 'minikube kubectl -- get pods -A'
🏄  Done! kubectl is now configured to use "minikube" cluster and "default" namespace by default

```

- Output of `kubectl cluster-info` and `kubectl get nodes`
```
(.venv) arthur@Artur-MacBook-Pro app_python % kubectl cluster-info

Kubernetes control plane is running at https://127.0.0.1:63953
CoreDNS is running at https://127.0.0.1:63953/api/v1/namespaces/kube-system/services/kube-dns:dns/proxy

To further debug and diagnose cluster problems, use 'kubectl cluster-info dump'.
```

```
(.venv) arthur@Artur-MacBook-Pro app_python % kubectl get nodes   
NAME       STATUS   ROLES           AGE     VERSION
minikube   Ready    control-plane   3m25s   v1.32.0
```

```
(.venv) arthur@Artur-MacBook-Pro DevOps-Core-Course % kubectl get pods -A
NAMESPACE              NAME                                         READY   STATUS    RESTARTS      AGE
kube-system            coredns-668d6bf9bc-9spct                     1/1     Running   1 (21m ago)   28m
kube-system            coredns-668d6bf9bc-pqvnp                     1/1     Running   1 (21m ago)   28m
kube-system            etcd-minikube                                1/1     Running   1 (21m ago)   28m
kube-system            kube-apiserver-minikube                      1/1     Running   1 (21m ago)   28m
kube-system            kube-controller-manager-minikube             1/1     Running   1 (21m ago)   28m
kube-system            kube-proxy-vvft9                             1/1     Running   1 (21m ago)   28m
kube-system            kube-scheduler-minikube                      1/1     Running   1 (21m ago)   28m
kube-system            storage-provisioner                          1/1     Running   3 (21m ago)   28m
kubernetes-dashboard   dashboard-metrics-scraper-5d59dccf9b-j9j94   1/1     Running   1 (21m ago)   28m
kubernetes-dashboard   kubernetes-dashboard-7779f9b69b-9wvqn        1/1     Running   2 (21m ago)   28m
```

- Brief explanation of your chosen tool (minikube/kind) and why

I used Minikube because it provides a simple and full-featured local Kubernetes cluster suitable for development and testing.

### Task 2 — Application Deployment

Test output:
```
(.venv) arthur@Artur-MacBook-Pro DevOps-Core-Course % kubectl apply -f k8s/deployment.yml
deployment.apps/devops-info-service created

```
```
(.venv) arthur@Artur-MacBook-Pro DevOps-Core-Course % kubectl get deployments

NAME                  READY   UP-TO-DATE   AVAILABLE   AGE
devops-info-service   0/3     3            0           24s
```
```
(.venv) arthur@Artur-MacBook-Pro DevOps-Core-Course % kubectl get pods

NAME                                   READY   STATUS    RESTARTS   AGE
devops-info-service-698cc89c85-d2dq6   1/1     Running   0          42s
devops-info-service-698cc89c85-dr9n4   1/1     Running   0          42s
devops-info-service-698cc89c85-hwtrt   1/1     Running   0          42s
```
```
(.venv) arthur@Artur-MacBook-Pro DevOps-Core-Course % kubectl describe deployment devops-info-service
Name:                   devops-info-service
Namespace:              default
CreationTimestamp:      Wed, 18 Mar 2026 12:14:06 +0300
Labels:                 app=devops-info-service
Annotations:            deployment.kubernetes.io/revision: 1
Selector:               app=devops-info-service
Replicas:               3 desired | 3 updated | 3 total | 3 available | 0 unavailable
StrategyType:           RollingUpdate
MinReadySeconds:        0
RollingUpdateStrategy:  0 max unavailable, 1 max surge
Pod Template:
  Labels:  app=devops-info-service
  Containers:
   devops-info-service:
    Image:      poparthur/devops-info-service:latest
    Port:       8000/TCP
    Host Port:  0/TCP
    Limits:
      cpu:     200m
      memory:  256Mi
    Requests:
      cpu:         100m
      memory:      128Mi
    Liveness:      http-get http://:8000/health delay=10s timeout=1s period=10s #success=1 #failure=3
    Readiness:     http-get http://:8000/health delay=5s timeout=1s period=5s #success=1 #failure=3
    Environment:   <none>
    Mounts:        <none>
  Volumes:         <none>
  Node-Selectors:  <none>
  Tolerations:     <none>
Conditions:
  Type           Status  Reason
  ----           ------  ------
  Available      True    MinimumReplicasAvailable
  Progressing    True    NewReplicaSetAvailable
OldReplicaSets:  <none>
NewReplicaSet:   devops-info-service-698cc89c85 (3/3 replicas created)
Events:
  Type    Reason             Age   From                   Message
  ----    ------             ----  ----                   -------
  Normal  ScalingReplicaSet  77s   deployment-controller  Scaled up replica set devops-info-service-698cc89c85 from 0 to 3
```

### Task 3 — Service Configuration  
```
(.venv) arthur@Artur-MacBook-Pro DevOps-Core-Course % minikube service devops-info-service --url
http://127.0.0.1:62821

(.venv) arthur@Artur-MacBook-Pro DevOps-Core-Course % curl http://127.0.0.1:62821/health
{"status":"healthy","timestamp":"2026-03-25T05:50:43.345887+00:00","uptime_seconds":161}

(.venv) arthur@Artur-MacBook-Pro DevOps-Core-Course % curl http://127.0.0.1:62821/      
{"service":{"name":"devops-info-service","version":"1.0.0","description":"DevOps course info service","framework":"FastAPI"},"system":{"hostname":"devops-info-service-698cc89c85-hwtrt","platform":"Linux","platform_version":"Linux-6.10.14-linuxkit-aarch64-with-glibc2.41","architecture":"aarch64","cpu_count":8,"python_version":"3.13.12"},"runtime":{"uptime_seconds":319,"uptime_human":"0 hours, 5 minutes","current_time":"2026-03-25T05:53:21.416830+00:00","timezone":"UTC"},"request":{"client_ip":"10.244.0.1","user_agent":"curl/8.7.1","method":"GET","path":"/"},"endpoints":[{"path":"/","method":"GET","description":"Service information"},{"path":"/health","method":"GET","description":"Health check"}]}% 
```

### Task 4 — Scaling and Updates

Scaling

```
(.venv) arthur@Artur-MacBook-Pro DevOps-Core-Course % kubectl scale deployment/devops-info-service --replicas=5
deployment.apps/devops-info-service scaled

(.venv) arthur@Artur-MacBook-Pro DevOps-Core-Course % kubectl get pods
NAME                                   READY   STATUS    RESTARTS        AGE
devops-info-service-698cc89c85-d2dq6   1/1     Running   1 (5d15h ago)   6d20h
devops-info-service-698cc89c85-dr9n4   1/1     Running   1 (5d15h ago)   6d20h
devops-info-service-698cc89c85-hwtrt   1/1     Running   1 (5d15h ago)   6d20h
devops-info-service-698cc89c85-rk7jq   1/1     Running   0               65s
devops-info-service-698cc89c85-tlpss   1/1     Running   0               65s
```

Rolling Updates
```
(.venv) arthur@Artur-MacBook-Pro DevOps-Core-Course % kubectl rollout status deployment/devops-info-service
deployment "devops-info-service" successfully rolled out

(.venv) arthur@Artur-MacBook-Pro DevOps-Core-Course % kubectl apply -f k8s/deployment.yml
deployment.apps/devops-info-service configured

(.venv) arthur@Artur-MacBook-Pro DevOps-Core-Course % kubectl get pods -w
NAME                                   READY   STATUS              RESTARTS        AGE
devops-info-service-698cc89c85-d2dq6   1/1     Running             1 (5d15h ago)   6d20h
devops-info-service-698cc89c85-dr9n4   1/1     Running             1 (5d15h ago)   6d20h
devops-info-service-698cc89c85-rk7jq   1/1     Terminating         0               10m
devops-info-service-df6cfbdc-n4r4h     0/1     ContainerCreating   0               1s
devops-info-service-df6cfbdc-q488d     1/1     Running             0               18s
devops-info-service-698cc89c85-rk7jq   0/1     Completed           0               10m
devops-info-service-698cc89c85-rk7jq   0/1     Completed           0               10m
devops-info-service-698cc89c85-rk7jq   0/1     Completed           0               10m
devops-info-service-df6cfbdc-n4r4h     0/1     Running             0               5s
devops-info-service-df6cfbdc-n4r4h     1/1     Running             0               13s
devops-info-service-698cc89c85-d2dq6   1/1     Terminating         1 (5d15h ago)   6d20h
devops-info-service-df6cfbdc-ck2fj     0/1     Pending             0               0s
devops-info-service-df6cfbdc-ck2fj     0/1     Pending             0               0s
devops-info-service-df6cfbdc-ck2fj     0/1     ContainerCreating   0               0s
devops-info-service-698cc89c85-d2dq6   0/1     Completed           1 (5d15h ago)   6d20h
devops-info-service-698cc89c85-d2dq6   0/1     Completed           1 (5d15h ago)   6d20h
devops-info-service-698cc89c85-d2dq6   0/1     Completed           1 (5d15h ago)   6d20h
devops-info-service-df6cfbdc-ck2fj     0/1     Running             0               5s
devops-info-service-df6cfbdc-ck2fj     1/1     Running             0               11s
devops-info-service-698cc89c85-dr9n4   1/1     Terminating         1 (5d15h ago)   6d20h
devops-info-service-698cc89c85-dr9n4   0/1     Completed           1 (5d15h ago)   6d20h
devops-info-service-698cc89c85-dr9n4   0/1     Completed           1 (5d15h ago)   6d20h
devops-info-service-698cc89c85-dr9n4   0/1     Completed           1 (5d15h ago)   6d20h

(.venv) arthur@Artur-MacBook-Pro DevOps-Core-Course % kubectl get pods
NAME                                 READY   STATUS    RESTARTS   AGE
devops-info-service-df6cfbdc-ck2fj   1/1     Running   0          92s
devops-info-service-df6cfbdc-n4r4h   1/1     Running   0          105s
devops-info-service-df6cfbdc-q488d   1/1     Running   0          2m2s

(.venv) arthur@Artur-MacBook-Pro DevOps-Core-Course % kubectl rollout history deployment/devops-info-service
deployment.apps/devops-info-service 
REVISION  CHANGE-CAUSE
1         <none>
2         <none>
```

Rollback

```
(.venv) arthur@Artur-MacBook-Pro DevOps-Core-Course % kubectl rollout undo deployment/devops-info-service
deployment.apps/devops-info-service rolled back
(.venv) arthur@Artur-MacBook-Pro DevOps-Core-Course % kubectl get pods -w
NAME                                   READY   STATUS              RESTARTS   AGE
devops-info-service-698cc89c85-gp6kq   0/1     ContainerCreating   0          4s
devops-info-service-df6cfbdc-ck2fj     1/1     Running             0          2m40s
devops-info-service-df6cfbdc-n4r4h     1/1     Running             0          2m53s
devops-info-service-df6cfbdc-q488d     1/1     Running             0          3m10s
devops-info-service-698cc89c85-gp6kq   0/1     Running             0          9s
devops-info-service-698cc89c85-gp6kq   1/1     Running             0          16s
devops-info-service-df6cfbdc-ck2fj     1/1     Terminating         0          2m52s
devops-info-service-698cc89c85-zctbl   0/1     Pending             0          0s
devops-info-service-698cc89c85-zctbl   0/1     Pending             0          0s
devops-info-service-698cc89c85-zctbl   0/1     ContainerCreating   0          0s
devops-info-service-df6cfbdc-ck2fj     0/1     Completed           0          2m53s
devops-info-service-df6cfbdc-ck2fj     0/1     Completed           0          2m53s
devops-info-service-df6cfbdc-ck2fj     0/1     Completed           0          2m53s
```

### Task 5 - Documentation

## 1. Architecture Overview

The application is deployed in a local Kubernetes cluster using Minikube. A Deployment manages multiple identical Pods running the devops-info-service container on port 8000.

Initially, the system runs 3 Pods, later scaled to 5, ensuring high availability and load distribution. Each Pod includes liveness and readiness probes (/health) and defined resource limits.

A NodePort Service exposes the application:
```
Client → NodePort (30080) → Service → Pod → Container (8000)
```
The Service uses label selectors to route traffic and provides basic load balancing across Pods.

Resource allocation is configured with requests (100m CPU, 128Mi memory) and limits (200m CPU, 256Mi memory) to ensure stable performance and prevent resource overuse.

## 2. Manifest Files

Two main Kubernetes manifests were created: deployment.yml and service.yml.

Deployment (deployment.yml) defines the application Pods and their lifecycle. It uses 3 replicas (scaled to 5 for testing) to ensure high availability. Resource requests (100m CPU, 128Mi memory) and limits (200m CPU, 256Mi) were set to balance performance and cluster stability. Liveness and readiness probes (/health) were added to enable self-healing and proper traffic routing. A rolling update strategy with zero downtime (maxUnavailable: 0) was configured.

Service (service.yml) exposes the application using a NodePort. It maps port 80 → 8000 and uses label selectors to route traffic to the Pods. NodePort was chosen to allow external access in a local Minikube environment.

## 3. Deployment Evidence

Above

## 4. Operations Performed

Above

## 5.  Production Considerations

**Health Checks**

Liveness and readiness probes were implemented using the /health endpoint.
	•	Liveness probe ensures containers are restarted if the application becomes unresponsive
	•	Readiness probe ensures traffic is only sent to healthy Pods
This improves reliability and prevents failed instances from receiving requests

**Resource Limits Rationale**

Requests (100m CPU, 128Mi memory) guarantee minimum resources for stable startup, while limits (200m CPU, 256Mi) prevent a single container from exhausting node resources. This balance ensures predictable performance and cluster stability.

**Improvements for Production**

- Use Horizontal Pod Autoscaler (HPA) for automatic scaling
- Replace NodePort with Ingress or LoadBalancer
- Use separate environments (dev/staging/prod)
- Store configuration in ConfigMaps and Secrets
- Use versioned images instead of latest

**Monitoring & Observability**

- Collect metrics using Prometheus
- Visualize with Grafana
- Use centralized logging (e.g., ELK stack)
- Monitor health, resource usage, and request rates

## 6. Challenges & Solutions

**Issues Encountered**

- Initial issues with cluster connectivity (connection refused) when Minikube was not running
- Limited Docker disk space caused warnings and potential deployment failures
- Difficulty editing resources using kubectl edit (editor usability issues)

**Debugging Approach**

- Used kubectl get pods to check Pod status
- Used kubectl describe pod to inspect events and configuration
- Used kubectl logs to analyze application behavior
- Monitored rollout process with `kubectl get pods -w`

**What I Learned**

- Kubernetes follows a declarative model, maintaining desired state automatically
- Deployments manage scaling, updates, and self-healing of Pods
- Services provide abstraction and load balancing between Pods
- Rolling updates and rollbacks allow safe, zero-downtime changes
