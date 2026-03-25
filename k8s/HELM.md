# Lab 10 — Helm Package Manager

### Task 1 — Helm Fundamentals

- Terminal output showing Helm installation and version (should be 4.x)

```
(.venv) arthur@Artur-MacBook-Pro DevOps-Core-Course % helm version
version.BuildInfo{Version:"v4.1.3", GitCommit:"c94d381b03be117e7e57908edbf642104e00eb8f", GitTreeState:"clean", GoVersion:"go1.26.1", KubeClientVersion:"v1.35"}
```

- Output of exploring a public chart (e.g., helm show chart prometheus-community/prometheus)
```
(.venv) arthur@Artur-MacBook-Pro DevOps-Core-Course % helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo update
"prometheus-community" has been added to your repositories
Hang tight while we grab the latest from your chart repositories...
...Successfully got an update from the "prometheus-community" chart repository
Update Complete. ⎈Happy Helming!⎈
(.venv) arthur@Artur-MacBook-Pro DevOps-Core-Course % helm search repo prometheus
NAME                                                    CHART VERSION   APP VERSION     DESCRIPTION                                       
prometheus-community/kube-prometheus-stack              82.14.0         v0.89.0         kube-prometheus-stack collects Kubernetes manif...
prometheus-community/prometheus                         28.14.0         v3.10.0         Prometheus is a monitoring system and time seri...
prometheus-community/prometheus-adapter                 5.3.0           v0.12.0         A Helm chart for k8s prometheus adapter           
prometheus-community/prometheus-blackbox-exporter       11.9.0          v0.28.0         Prometheus Blackbox Exporter                      
prometheus-community/prometheus-cloudwatch-expo...      0.28.1          0.16.0          A Helm chart for prometheus cloudwatch-exporter   
prometheus-community/prometheus-conntrack-stats...      0.5.35          v0.4.42         A Helm chart for conntrack-stats-exporter         
prometheus-community/prometheus-consul-exporter         1.1.1           v0.13.0         A Helm chart for the Prometheus Consul Exporter   
prometheus-community/prometheus-couchdb-exporter        1.0.1           1.0             A Helm chart to export the metrics from couchdb...
prometheus-community/prometheus-druid-exporter          1.2.0           v0.11.0         Druid exporter to monitor druid metrics with Pr...
prometheus-community/prometheus-elasticsearch-e...      7.2.1           v1.10.0         Elasticsearch stats exporter for Prometheus       
prometheus-community/prometheus-fastly-exporter         0.11.0          v10.2.0         A Helm chart for the Prometheus Fastly Exporter   
prometheus-community/prometheus-ipmi-exporter           0.8.0           v1.10.1         This is an IPMI exporter for Prometheus.          
prometheus-community/prometheus-json-exporter           0.19.2          v0.7.0          Install prometheus-json-exporter                  
prometheus-community/prometheus-kafka-exporter          3.0.1           v1.9.0          A Helm chart to export metrics from Kafka in Pr...
prometheus-community/prometheus-memcached-exporter      0.4.5           v0.15.5         Prometheus exporter for Memcached metrics         
prometheus-community/prometheus-modbus-exporter         0.1.4           0.4.1           A Helm chart for prometheus-modbus-exporter       
prometheus-community/prometheus-mongodb-exporter        3.18.0          0.49.0          A Prometheus exporter for MongoDB metrics         
prometheus-community/prometheus-mysql-exporter          2.13.0          v0.19.0         A Helm chart for prometheus mysql exporter with...
prometheus-community/prometheus-nats-exporter           2.21.1          0.18.0          A Helm chart for prometheus-nats-exporter         
prometheus-community/prometheus-nginx-exporter          1.20.8          1.5.1           A Helm chart for NGINX Prometheus Exporter        
prometheus-community/prometheus-node-exporter           4.52.2          1.10.2          A Helm chart for prometheus node-exporter         
prometheus-community/prometheus-opencost-exporter       0.1.2           1.108.0         Prometheus OpenCost Exporter                      
prometheus-community/prometheus-operator                9.3.2           0.38.1          DEPRECATED - This chart will be renamed. See ht...
prometheus-community/prometheus-operator-admiss...      0.37.1          0.90.0          Prometheus Operator Admission Webhook             
prometheus-community/prometheus-operator-crds           28.0.0          v0.90.0         A Helm chart that collects custom resource defi...
prometheus-community/prometheus-pgbouncer-exporter      0.10.0          v0.12.0         A Helm chart for prometheus pgbouncer-exporter    
prometheus-community/prometheus-pingdom-exporter        3.4.2           v0.5.6          A Helm chart for Prometheus Pingdom Exporter      
prometheus-community/prometheus-pingmesh-exporter       0.4.3           v1.2.2          Prometheus Pingmesh Exporter                      
prometheus-community/prometheus-postgres-exporter       7.5.2           v0.19.1         A Helm chart for prometheus postgres-exporter     
prometheus-community/prometheus-pushgateway             3.6.0           v1.11.2         A Helm chart for prometheus pushgateway           
prometheus-community/prometheus-rabbitmq-exporter       2.1.2           1.0.0           Rabbitmq metrics exporter for prometheus          
prometheus-community/prometheus-redis-exporter          6.22.0          v1.82.0         Prometheus exporter for Redis metrics             
prometheus-community/prometheus-smartctl-exporter       0.16.0          v0.14.0         A Helm chart for Kubernetes                       
prometheus-community/prometheus-snmp-exporter           9.13.0          v0.30.1         Prometheus SNMP Exporter                          
prometheus-community/prometheus-sql-exporter            0.5.0           v0.8            Prometheus SQL Exporter                           
prometheus-community/prometheus-stackdriver-exp...      4.12.2          v0.18.0         Stackdriver exporter for Prometheus               
prometheus-community/prometheus-statsd-exporter         1.0.0           v0.28.0         A Helm chart for prometheus stats-exporter        
prometheus-community/prometheus-systemd-exporter        0.5.2           0.7.0           A Helm chart for prometheus systemd-exporter      
prometheus-community/prometheus-to-sd                   0.5.1           v0.9.2          Scrape metrics stored in prometheus format and ...
prometheus-community/prometheus-windows-exporter        0.12.5          0.31.5          A Helm chart for prometheus windows-exporter      
prometheus-community/prometheus-yet-another-clo...      0.42.1          v0.63.0         Yace - Yet Another CloudWatch Exporter            
prometheus-community/alertmanager                       1.34.0          v0.31.1         The Alertmanager handles alerts sent by client ...
prometheus-community/alertmanager-snmp-notifier         2.1.0           v2.1.0          The SNMP Notifier handles alerts coming from Pr...
prometheus-community/jiralert                           1.8.2           v1.3.0          A Helm chart for Kubernetes to install jiralert   
prometheus-community/kube-state-metrics                 7.2.2           2.18.0          Install kube-state-metrics to generate and expo...
prometheus-community/prom-label-proxy                   0.18.0          v0.12.1         A proxy that enforces a given label in a given ...
prometheus-community/yet-another-cloudwatch-exp...      0.39.1          v0.62.1         Yace - Yet Another CloudWatch Exporter            
(.venv) arthur@Artur-MacBook-Pro DevOps-Core-Course % helm show chart prometheus-community/prometheus
annotations:
  artifacthub.io/license: Apache-2.0
  artifacthub.io/links: |
    - name: Chart Source
      url: https://github.com/prometheus-community/helm-charts
    - name: Upstream Project
      url: https://github.com/prometheus/prometheus
apiVersion: v2
appVersion: v3.10.0
dependencies:
- condition: alertmanager.enabled
  name: alertmanager
  repository: https://prometheus-community.github.io/helm-charts
  version: 1.34.*
- condition: kube-state-metrics.enabled
  name: kube-state-metrics
  repository: https://prometheus-community.github.io/helm-charts
  version: 7.2.*
- condition: prometheus-node-exporter.enabled
  name: prometheus-node-exporter
  repository: https://prometheus-community.github.io/helm-charts
  version: 4.52.*
- condition: prometheus-pushgateway.enabled
  name: prometheus-pushgateway
  repository: https://prometheus-community.github.io/helm-charts
  version: 3.6.*
description: Prometheus is a monitoring system and time series database.
home: https://prometheus.io/
icon: https://raw.githubusercontent.com/prometheus/prometheus.github.io/master/assets/prometheus_logo-cb55bb5c346.png
keywords:
- monitoring
- prometheus
kubeVersion: '>=1.19.0-0'
maintainers:
- email: gianrubio@gmail.com
  name: gianrubio
  url: https://github.com/gianrubio
- email: zanhsieh@gmail.com
  name: zanhsieh
  url: https://github.com/zanhsieh
- email: miroslav.hadzhiev@gmail.com
  name: Xtigyro
  url: https://github.com/Xtigyro
- email: naseem@transit.app
  name: naseemkullah
  url: https://github.com/naseemkullah
- email: rootsandtrees@posteo.de
  name: zeritti
  url: https://github.com/zeritti
name: prometheus
sources:
- https://github.com/prometheus/alertmanager
- https://github.com/prometheus/prometheus
- https://github.com/prometheus/pushgateway
- https://github.com/prometheus/node_exporter
- https://github.com/kubernetes/kube-state-metrics
type: application
version: 28.14.0
```

- Brief explanation of Helm's value proposition

Helm is a package manager for Kubernetes that simplifies application deployment and management. It enables reusable configurations through templating, supports versioning and rollbacks, and helps manage complex applications with dependencies. Helm standardizes deployments and makes it easier to configure applications across different environments (e.g., development and production).

### Task 2 — Create Your Helm Chart

```
(.venv) arthur@Artur-MacBook-Pro DevOps-Core-Course % helm create k8s/myapp
Creating k8s/myapp
```

```
(.venv) arthur@Artur-MacBook-Pro DevOps-Core-Course % helm lint k8s/myapp                  

==> Linting k8s/myapp
[INFO] Chart.yaml: icon is recommended

1 chart(s) linted, 0 chart(s) failed
```

```
(.venv) arthur@Artur-MacBook-Pro DevOps-Core-Course % helm template myapp k8s/myapp

---
# Source: myapp/templates/service.yaml
apiVersion: v1
kind: Service
metadata:
  name: myapp-myapp
  labels:
    app.kubernetes.io/name: myapp
    app.kubernetes.io/instance: myapp
spec:
  type: NodePort
  selector:
    app.kubernetes.io/name: myapp
    app.kubernetes.io/instance: myapp
  ports:
    - port: 80
      targetPort: 8000
---
# Source: myapp/templates/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: myapp-myapp
  labels:
    app.kubernetes.io/name: myapp
    app.kubernetes.io/instance: myapp
spec:
  replicas: 2
  selector:
    matchLabels:
      app.kubernetes.io/name: myapp
      app.kubernetes.io/instance: myapp
  template:
    metadata:
      labels:
        app.kubernetes.io/name: myapp
        app.kubernetes.io/instance: myapp
    spec:
      containers:
      - name: myapp
        image: "yourdocker/myapp:1.0"
        imagePullPolicy: IfNotPresent
        ports:
        - containerPort: 8000

        resources:
          limits:
            cpu: 200m
            memory: 256Mi
          requests:
            cpu: 100m
            memory: 128Mi

        livenessProbe:
          httpGet:
            path: /health
            port: 8000
          initialDelaySeconds: 10
          periodSeconds: 5

        readinessProbe:
          httpGet:
            path: /ready
            port: 8000
          initialDelaySeconds: 5
          periodSeconds: 3
```

```
(.venv) arthur@Artur-MacBook-Pro DevOps-Core-Course % helm install --dry-run --debug test-release k8s/myapp

level=WARN msg="--dry-run is deprecated and should be replaced with '--dry-run=client'"
level=DEBUG msg="Original chart version" version=""
level=DEBUG msg="Chart path" path=/Users/arthur/PycharmProjects/DevOps-Core-Course/k8s/myapp
level=DEBUG msg="number of dependencies in the chart" chart=myapp dependencies=0
NAME: test-release
LAST DEPLOYED: Wed Mar 25 10:11:09 2026
NAMESPACE: default
STATUS: pending-install
REVISION: 1
DESCRIPTION: Dry run complete
TEST SUITE: None
USER-SUPPLIED VALUES:
{}

COMPUTED VALUES:
image:
  pullPolicy: IfNotPresent
  repository: yourdocker/myapp
  tag: "1.0"
livenessProbe:
  httpGet:
    path: /health
    port: 8000
  initialDelaySeconds: 10
  periodSeconds: 5
readinessProbe:
  httpGet:
    path: /ready
    port: 8000
  initialDelaySeconds: 5
  periodSeconds: 3
replicaCount: 2
resources:
  limits:
    cpu: 200m
    memory: 256Mi
  requests:
    cpu: 100m
    memory: 128Mi
service:
  port: 80
  targetPort: 8000
  type: NodePort
serviceAccount:
  create: true
  name: ""

HOOKS:
MANIFEST:
---
# Source: myapp/templates/service.yaml
apiVersion: v1
kind: Service
metadata:
  name: test-release-myapp
  labels:
    app.kubernetes.io/name: myapp
    app.kubernetes.io/instance: test-release
spec:
  type: NodePort
  selector:
    app.kubernetes.io/name: myapp
    app.kubernetes.io/instance: test-release
  ports:
    - port: 80
      targetPort: 8000
---
# Source: myapp/templates/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: test-release-myapp
  labels:
    app.kubernetes.io/name: myapp
    app.kubernetes.io/instance: test-release
spec:
  replicas: 2
  selector:
    matchLabels:
      app.kubernetes.io/name: myapp
      app.kubernetes.io/instance: test-release
  template:
    metadata:
      labels:
        app.kubernetes.io/name: myapp
        app.kubernetes.io/instance: test-release
    spec:
      containers:
      - name: myapp
        image: "yourdocker/myapp:1.0"
        imagePullPolicy: IfNotPresent
        ports:
        - containerPort: 8000

        resources:
          limits:
            cpu: 200m
            memory: 256Mi
          requests:
            cpu: 100m
            memory: 128Mi

        livenessProbe:
          httpGet:
            path: /health
            port: 8000
          initialDelaySeconds: 10
          periodSeconds: 5

        readinessProbe:
          httpGet:
            path: /ready
            port: 8000
          initialDelaySeconds: 5
          periodSeconds: 3
```

```
(.venv) arthur@Artur-MacBook-Pro DevOps-Core-Course % helm install myrelease k8s/myapp
NAME: myrelease
LAST DEPLOYED: Wed Mar 25 10:11:38 2026
NAMESPACE: default
STATUS: deployed
REVISION: 1
DESCRIPTION: Install complete
TEST SUITE: None
```

### Task 3 — Multi-Environment Support

```
(.venv) arthur@Artur-MacBook-Pro DevOps-Core-Course % helm install myapp-dev k8s/myapp -f k8s/myapp/values-dev.yaml
NAME: myapp-dev
LAST DEPLOYED: Wed Mar 25 10:15:27 2026
NAMESPACE: default
STATUS: deployed
REVISION: 1
DESCRIPTION: Install complete
TEST SUITE: None
(.venv) arthur@Artur-MacBook-Pro DevOps-Core-Course % kubectl get pods
kubectl get svc
NAME                                   READY   STATUS              RESTARTS   AGE
devops-info-service-698cc89c85-4jqc2   1/1     Running             0          65m
devops-info-service-698cc89c85-gp6kq   1/1     Running             0          66m
devops-info-service-698cc89c85-zctbl   1/1     Running             0          66m
myapp-dev-myapp-6c6dccf6c4-867p6       0/1     ContainerCreating   0          4s
myrelease-myapp-68fb47bf89-lhmtt       0/1     ImagePullBackOff    0          3m52s
myrelease-myapp-68fb47bf89-pznv6       0/1     ErrImagePull        0          3m52s
NAME                  TYPE        CLUSTER-IP       EXTERNAL-IP   PORT(S)        AGE
devops-info-service   NodePort    10.105.95.35     <none>        80:30080/TCP   6d21h
kubernetes            ClusterIP   10.96.0.1        <none>        443/TCP        6d22h
myapp-dev-myapp       NodePort    10.104.127.109   <none>        80:31895/TCP   4s
myrelease-myapp       NodePort    10.97.121.210    <none>        80:31761/TCP   3m53s
(.venv) arthur@Artur-MacBook-Pro DevOps-Core-Course % helm upgrade myapp-dev k8s/myapp -f k8s/myapp/values-prod.yaml
Release "myapp-dev" has been upgraded. Happy Helming!
NAME: myapp-dev
LAST DEPLOYED: Wed Mar 25 10:15:48 2026
NAMESPACE: default
STATUS: deployed
REVISION: 2
DESCRIPTION: Upgrade complete
TEST SUITE: None
(.venv) arthur@Artur-MacBook-Pro DevOps-Core-Course % kubectl get pods
kubectl get svc
NAME                                   READY   STATUS              RESTARTS   AGE
devops-info-service-698cc89c85-4jqc2   1/1     Running             0          66m
devops-info-service-698cc89c85-gp6kq   1/1     Running             0          66m
devops-info-service-698cc89c85-zctbl   1/1     Running             0          66m
myapp-dev-myapp-6c6dccf6c4-867p6       0/1     ImagePullBackOff    0          31s
myapp-dev-myapp-6c6dccf6c4-cshnn       0/1     ContainerCreating   0          10s
myapp-dev-myapp-6c6dccf6c4-fnxpq       0/1     ContainerCreating   0          10s
myapp-dev-myapp-6c6dccf6c4-ns4h9       0/1     Terminating         0          10s
myapp-dev-myapp-c549cb59f-59gd9        0/1     ErrImagePull        0          10s
myapp-dev-myapp-c549cb59f-ccbd9        0/1     ContainerCreating   0          9s
myrelease-myapp-68fb47bf89-lhmtt       0/1     ImagePullBackOff    0          4m19s
myrelease-myapp-68fb47bf89-pznv6       0/1     ImagePullBackOff    0          4m19s
NAME                  TYPE           CLUSTER-IP       EXTERNAL-IP   PORT(S)        AGE
devops-info-service   NodePort       10.105.95.35     <none>        80:30080/TCP   6d21h
kubernetes            ClusterIP      10.96.0.1        <none>        443/TCP        6d22h
myapp-dev-myapp       LoadBalancer   10.104.127.109   <pending>     80:31895/TCP   31s
myrelease-myapp       NodePort       10.97.121.210    <none>        80:31761/TCP   4m20s
```

### Task 4 — Chart Hooks

Added: 
- [pre-install-job.yaml](myapp/templates/hooks/pre-install-job.yaml)
- [post-install-job.yaml](myapp/templates/hooks/post-install-job.yaml)


```
(.venv) arthur@Artur-MacBook-Pro DevOps-Core-Course % helm lint k8s/myapp
==> Linting k8s/myapp
[INFO] Chart.yaml: icon is recommended

1 chart(s) linted, 0 chart(s) failed
(.venv) arthur@Artur-MacBook-Pro DevOps-Core-Course % helm install --dry-run --debug test-hooks k8s/myapp
level=WARN msg="--dry-run is deprecated and should be replaced with '--dry-run=client'"
level=DEBUG msg="Original chart version" version=""
level=DEBUG msg="Chart path" path=/Users/arthur/PycharmProjects/DevOps-Core-Course/k8s/myapp
level=DEBUG msg="number of dependencies in the chart" chart=myapp dependencies=0
NAME: test-hooks
LAST DEPLOYED: Wed Mar 25 12:05:23 2026
NAMESPACE: default
STATUS: pending-install
REVISION: 1
DESCRIPTION: Dry run complete
TEST SUITE: None
USER-SUPPLIED VALUES:
{}

COMPUTED VALUES:
image:
  pullPolicy: IfNotPresent
  repository: poparthur/devops-info-service:latest
  tag: "1.0"
livenessProbe:
  httpGet:
    path: /health
    port: 8000
  initialDelaySeconds: 10
  periodSeconds: 5
readinessProbe:
  httpGet:
    path: /ready
    port: 8000
  initialDelaySeconds: 5
  periodSeconds: 3
replicaCount: 2
resources:
  limits:
    cpu: 200m
    memory: 256Mi
  requests:
    cpu: 100m
    memory: 128Mi
service:
  port: 80
  targetPort: 8000
  type: NodePort
serviceAccount:
  create: true
  name: ""

HOOKS:
---
# Source: myapp/templates/hooks/post-install-job.yaml
apiVersion: batch/v1
kind: Job
metadata:
  name: "test-hooks-post-install"
  annotations:
    "helm.sh/hook": post-install
    "helm.sh/hook-weight": "5"
    "helm.sh/hook-delete-policy": hook-succeeded
spec:
  template:
    spec:
      restartPolicy: Never
      containers:
        - name: post-install
          image: busybox
          command: ["sh", "-c", "echo Post-install validation && sleep 5"]
---
# Source: myapp/templates/hooks/pre-install-job.yaml
apiVersion: batch/v1
kind: Job
metadata:
  name: "test-hooks-pre-install"
  annotations:
    "helm.sh/hook": pre-install
    "helm.sh/hook-weight": "-5"
    "helm.sh/hook-delete-policy": hook-succeeded
spec:
  template:
    spec:
      restartPolicy: Never
      containers:
        - name: pre-install
          image: busybox
          command: ["sh", "-c", "echo Pre-install check && sleep 5"]
MANIFEST:
---
# Source: myapp/templates/service.yaml
apiVersion: v1
kind: Service
metadata:
  name: test-hooks-myapp
  labels:
    app.kubernetes.io/name: myapp
    app.kubernetes.io/instance: test-hooks
spec:
  type: NodePort
  selector:
    app.kubernetes.io/name: myapp
    app.kubernetes.io/instance: test-hooks
  ports:
    - port: 80
      targetPort: 8000
---
# Source: myapp/templates/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: test-hooks-myapp
  labels:
    app.kubernetes.io/name: myapp
    app.kubernetes.io/instance: test-hooks
spec:
  replicas: 2
  selector:
    matchLabels:
      app.kubernetes.io/name: myapp
      app.kubernetes.io/instance: test-hooks
  template:
    metadata:
      labels:
        app.kubernetes.io/name: myapp
        app.kubernetes.io/instance: test-hooks
    spec:
      containers:
      - name: myapp
        image: "poparthur/devops-info-service:latest:1.0"
        imagePullPolicy: IfNotPresent
        ports:
        - containerPort: 8000

        resources:
          limits:
            cpu: 200m
            memory: 256Mi
          requests:
            cpu: 100m
            memory: 128Mi

        livenessProbe:
          httpGet:
            path: /health
            port: 8000
          initialDelaySeconds: 10
          periodSeconds: 5

        readinessProbe:
          httpGet:
            path: /ready
            port: 8000
          initialDelaySeconds: 5
          periodSeconds: 3

(.venv) arthur@Artur-MacBook-Pro DevOps-Core-Course % helm install hooks-test k8s/myapp
NAME: hooks-test
LAST DEPLOYED: Wed Mar 25 12:06:07 2026
NAMESPACE: default
STATUS: deployed
REVISION: 1
DESCRIPTION: Install complete
TEST SUITE: None
```

```
3m2s        Normal    Scheduled           pod/hooks-test-pre-install-8wg8z           Successfully assigned default/hooks-test-pre-install-8wg8z to minikube
3m3s        Normal    SuccessfulCreate    job/hooks-test-pre-install                 Created pod: hooks-test-pre-install-8wg8z
3m1s        Normal    Pulling             pod/hooks-test-pre-install-8wg8z           Pulling image "busybox"
2m30s       Normal    Pulled              pod/hooks-test-pre-install-8wg8z           Successfully pulled image "busybox" in 30.911s (30.912s including waiting). Image size: 4105254 bytes.
2m30s       Normal    Created             pod/hooks-test-pre-install-8wg8z           Created container: pre-install
2m30s       Normal    Started             pod/hooks-test-pre-install-8wg8z           Started container pre-install
2m22s       Normal    SuccessfulCreate    replicaset/hooks-test-myapp-596497567f     Created pod: hooks-test-myapp-596497567f-nvbcl
2m22s       Normal    Scheduled           pod/hooks-test-myapp-596497567f-nvbcl      Successfully assigned default/hooks-test-myapp-596497567f-nvbcl to minikube
2m22s       Normal    SuccessfulCreate    job/hooks-test-post-install                Created pod: hooks-test-post-install-q9drb
2m22s       Normal    Scheduled           pod/hooks-test-myapp-596497567f-kkl7z      Successfully assigned default/hooks-test-myapp-596497567f-kkl7z to minikube
2m22s       Normal    Completed           job/hooks-test-pre-install                 Job completed
2m22s       Normal    SuccessfulCreate    replicaset/hooks-test-myapp-596497567f     Created pod: hooks-test-myapp-596497567f-kkl7z
2m22s       Normal    ScalingReplicaSet   deployment/hooks-test-myapp                Scaled up replica set hooks-test-myapp-596497567f from 0 to 2
2m22s       Normal    Scheduled           pod/hooks-test-post-install-q9drb          Successfully assigned default/hooks-test-post-install-q9drb to minikube
2m20s       Normal    Pulling             pod/hooks-test-post-install-q9drb          Pulling image "busybox"
25s         Warning   Failed              pod/hooks-test-myapp-596497567f-nvbcl      Error: InvalidImageName
14s         Warning   InspectFailed       pod/hooks-test-myapp-596497567f-nvbcl      Failed to apply default image tag "poparthur/devops-info-service:latest:1.0": couldn't parse image name "poparthur/devops-info-service:latest:1.0": invalid reference format
25s         Warning   Failed              pod/hooks-test-myapp-596497567f-kkl7z      Error: InvalidImageName
10s         Warning   InspectFailed       pod/hooks-test-myapp-596497567f-kkl7z      Failed to apply default image tag "poparthur/devops-info-service:latest:1.0": couldn't parse image name "poparthur/devops-info-service:latest:1.0": invalid reference format
2m13s       Normal    Pulled              pod/hooks-test-post-install-q9drb          Successfully pulled image "busybox" in 5.127s (7.704s including waiting). Image size: 4105254 bytes.
2m13s       Normal    Created             pod/hooks-test-post-install-q9drb          Created container: post-install
2m13s       Normal    Started             pod/hooks-test-post-install-q9drb          Started container post-install
2m5s        Normal    Completed           job/hooks-test-post-install                Job completed
```


---

# Helm Chart Documentation

## 1. Chart Overview

The Helm chart packages a Kubernetes application using reusable templates and configurable values.

**Structure:**

```
k8s/myapp/
├── Chart.yaml
├── values.yaml
├── values-dev.yaml
├── values-prod.yaml
└── templates/
    ├── deployment.yaml
    ├── service.yaml
    ├── _helpers.tpl
    └── hooks/
        ├── pre-install-job.yaml
        └── post-install-job.yaml
```

**Key Templates:**

* `deployment.yaml` — defines application pods and container configuration
* `service.yaml` — exposes the application via Kubernetes Service
* `_helpers.tpl` — reusable templates for naming and labels
* `hooks/` — lifecycle jobs executed before and after deployment

**Values Strategy:**
Configuration is centralized in `values.yaml`, with overrides for dev and prod environments.

---

## 2. Configuration Guide

**Important Values:**

* `replicaCount` — number of pod replicas
* `image.repository` / `image.tag` — container image
* `service.type` — NodePort or LoadBalancer
* `resources` — CPU and memory limits
* `livenessProbe` / `readinessProbe` — health checks

**Environment Customization:**

* `values-dev.yaml`: 1 replica, NodePort, lower resources
* `values-prod.yaml`: multiple replicas, LoadBalancer, higher resources

**Example Usage:**

```bash
helm install myapp-dev k8s/myapp -f values-dev.yaml
helm upgrade myapp-dev k8s/myapp -f values-prod.yaml
```

---

## 3. Hook Implementation

Two hooks were implemented:

* **Pre-install hook** — runs before deployment to validate environment
* **Post-install hook** — runs after deployment to verify application

Execution Order:

* Pre-install hook (weight: -5)
* Application deployment
* Post-install hook (weight: 5)

**Deletion Policy:**

* `hook-succeeded` — automatically removes jobs after successful execution

---

## 4. Installation Evidence

**Helm Releases:**

```bash
helm list
```

**Kubernetes Resources:**

```bash
kubectl get all
```

**Hook Execution:**

```bash
kubectl get events
kubectl describe job <job-name>
```

**Environment Deployment:**

* Dev: 1 replica, NodePort
* Prod: multiple replicas, LoadBalancer

---

## 5. Operations

**Install:**

```bash
helm install myrelease k8s/myapp
```

**Upgrade:**

```bash
helm upgrade myrelease k8s/myapp -f values-prod.yaml
```

**Rollback:**

```bash
helm rollback myrelease 1
```

**Uninstall:**

```bash
helm uninstall myrelease
```

## 6. Testing & Validation

All above
