## Task 1 — StatefulSet Concepts

```
arthur@Artur-MacBook-Pro DevOps-Core-Course % kubectl explain statefulset | head -n 20

kubectl explain deployment | head -n 20
GROUP:      apps
KIND:       StatefulSet
VERSION:    v1

DESCRIPTION:
    StatefulSet represents a set of pods with consistent identities. Identities
    are defined as:
      - Network: A single stable DNS and hostname.
      - Storage: As many VolumeClaims as requested.
    
    The StatefulSet guarantees that a given network identity will always map to
    the same storage identity.
    
FIELDS:
  apiVersion	<string>
    APIVersion defines the versioned schema of this representation of an object.
    Servers should convert recognized schemas to the latest internal value, and
    may reject unrecognized values. More info:
    https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources

GROUP:      apps
KIND:       Deployment
VERSION:    v1

DESCRIPTION:
    Deployment enables declarative updates for Pods and ReplicaSets.
    
FIELDS:
  apiVersion	<string>
    APIVersion defines the versioned schema of this representation of an object.
    Servers should convert recognized schemas to the latest internal value, and
    may reject unrecognized values. More info:
    https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources

  kind	<string>
    Kind is a string value representing the REST resource this object
    represents. Servers may infer this from the endpoint the client submits
    requests to. Cannot be updated. In CamelCase. More info:
    https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds
```

### StatefulSet Guarantees

StatefulSets provide the following guarantees:

- **Stable Network Identity**
  Each pod gets a predictable name:
  `pod-0`, `pod-1`, `pod-2`

- **Stable Storage**
  Each pod has its own PersistentVolumeClaim (PVC)
  Storage is preserved even if the pod is restarted

- **Ordered Deployment & Scaling**
  Pods are created and deleted in order:
  `pod-0 → pod-1 → pod-2`

### Deployment vs StatefulSet

| Feature          | Deployment     | StatefulSet               |
|------------------|----------------|---------------------------|
| Pod Names        | Random suffix  | Stable (`pod-0`, `pod-1`) |
| Storage          | Shared or none | Per-pod PVC               |
| Scaling          | Parallel       | Ordered                   |
| Network Identity | Dynamic        | Stable DNS                |

### When to Use

#### Deployment:

- Stateless apps
- Web services
- APIs

#### StatefulSet:

- Databases (PostgreSQL, MongoDB)
- Message brokers (Kafka, RabbitMQ)
- Distributed systems (Elasticsearch)

### Headless Service

A Headless Service is defined with:

```yaml
clusterIP: None
```

How DNS works with StatefulSets:

StatefulSets create a Headless Service with `clusterIP: None`. Each pod gets a
stable DNS name: `<statefulset>-<ordinal>.<headless-service>.<namespace>.svc.cluster.local`.

For example, if the StatefulSet is named `myapp` and the headless service is named `myapp-headless`, the pod named `myapp-0` would have the DNS name `myapp-0.myapp-headless.default.svc.cluster.local`.

This allows each pod to discover the others by their stable DNS name.


## Task 2 — Convert Deployment to StatefulSet

```
arthur@Artur-MacBook-Pro DevOps-Core-Course % kubectl get storageclass
NAME                 PROVISIONER                RECLAIMPOLICY   VOLUMEBINDINGMODE   ALLOWVOLUMEEXPANSION   AGE
standard (default)   k8s.io/minikube-hostpath   Delete          Immediate           false                  7d13h
```

```
arthur@Artur-MacBook-Pro DevOps-Core-Course % helm upgrade --install myapp k8s/myapp
Release "myapp" does not exist. Installing it now.
NAME: myapp
LAST DEPLOYED: Wed Apr 15 10:45:16 2026
NAMESPACE: default
STATUS: deployed
REVISION: 1
DESCRIPTION: Install complete
TEST SUITE: None
arthur@Artur-MacBook-Pro DevOps-Core-Course % kubectl get statefulset
kubectl get pods
kubectl get pvc
NAME          READY   AGE
myapp-myapp   2/2     56s
NAME            READY   STATUS    RESTARTS   AGE
myapp-myapp-0   1/1     Running   0          55s
myapp-myapp-1   1/1     Running   0          41s
NAME                 STATUS   VOLUME                                     CAPACITY   ACCESS MODES   STORAGECLASS   VOLUMEATTRIBUTESCLASS   AGE
data-myapp-myapp-0   Bound    pvc-c4aab6e5-2413-4a61-8b7a-1d4ec82428ab   1Gi        RWO            standard       <unset>                 55s
data-myapp-myapp-1   Bound    pvc-8fdc60ad-bd96-4a6b-a4ac-b2d4b516f66a   1Gi        RWO            standard       <unset>                 41s
arthur@Artur-MacBook-Pro DevOps-Core-Course % 
```

```
arthur@Artur-MacBook-Pro DevOps-Core-Course % kubectl run test-dns --image=busybox:1.28 --rm -it --restart=Never -- nslookup myapp-myapp-0.myapp-myapp-headless
Server:    10.96.0.10
Address 1: 10.96.0.10 kube-dns.kube-system.svc.cluster.local

Name:      myapp-myapp-0.myapp-myapp-headless
Address 1: 10.244.0.129 10-244-0-129.myapp-myapp.default.svc.cluster.local
pod "test-dns" deleted from default namespace
arthur@Artur-MacBook-Pro DevOps-Core-Course % 
```

## Task 3 — Headless Service & Pod Identity

1. **Test DNS Resolution**
   - Exec into a pod
   - Resolve other pods via DNS
   - Document the DNS naming pattern

```
arthur@Artur-MacBook-Pro DevOps-Core-Course % kubectl run dns-test --image=busybox:1.28 --rm -it --restart=Never -- \
nslookup myapp-myapp-1.myapp-myapp-headless
Server:    10.96.0.10
Address 1: 10.96.0.10 kube-dns.kube-system.svc.cluster.local

Name:      myapp-myapp-1.myapp-myapp-headless
Address 1: 10.244.0.131 10-244-0-131.myapp-myapp-preview.default.svc.cluster.local
pod "dns-test" deleted from default namespace
```
DNS Naming Pattern

- StatefulSet pods follow a predictable DNS pattern:

- `<pod-name>.<headless-service-name>.<namespace>.svc.cluster.local`

- In this case:

- `myapp-myapp-1.myapp-myapp-headless.default.svc.cluster.local`


2. **Test Per-Pod Storage**
   - Access your app through each pod
   - Verify each pod maintains its own visit count
   - Demonstrate isolation between pods


```
arthur@Artur-MacBook-Pro DevOps-Core-Course % curl localhost:8080       

{"service":{"name":"devops-info-service","version":"1.0.0","description":"DevOps course info service","framework":"FastAPI"},"system":{"hostname":"myapp-myapp-0","platform":"Linux","platform_version":"Linux-6.10.14-linuxkit-aarch64-with-glibc2.41","architecture":"aarch64","cpu_count":8,"python_version":"3.13.12"},"runtime":{"uptime_seconds":903,"uptime_human":"0 hours, 15 minutes","current_time":"2026-04-15T08:00:57.891645+00:00","timezone":"UTC"},"request":{"client_ip":"127.0.0.1","user_agent":"curl/8.7.1","method":"GET","path":"/"},"endpoints":[{"path":"/","method":"GET","description":"Service information"},{"path":"/health","method":"GET","description":"Health check"}]}%                      arthur@Artur-MacBook-Pro DevOps-Core-Course % curl localhost:8080

{"service":{"name":"devops-info-service","version":"1.0.0","description":"DevOps course info service","framework":"FastAPI"},"system":{"hostname":"myapp-myapp-0","platform":"Linux","platform_version":"Linux-6.10.14-linuxkit-aarch64-with-glibc2.41","architecture":"aarch64","cpu_count":8,"python_version":"3.13.12"},"runtime":{"uptime_seconds":906,"uptime_human":"0 hours, 15 minutes","current_time":"2026-04-15T08:01:01.129084+00:00","timezone":"UTC"},"request":{"client_ip":"127.0.0.1","user_agent":"curl/8.7.1","method":"GET","path":"/"},"endpoints":[{"path":"/","method":"GET","description":"Service information"},{"path":"/health","method":"GET","description":"Health check"}]}%                      arthur@Artur-MacBook-Pro DevOps-Core-Course % curl localhost:8081
{"service":{"name":"devops-info-service","version":"1.0.0","description":"DevOps course info service","framework":"FastAPI"},"system":{"hostname":"myapp-myapp-1","platform":"Linux","platform_version":"Linux-6.10.14-linuxkit-aarch64-with-glibc2.41","architecture":"aarch64","cpu_count":8,"python_version":"3.13.12"},"runtime":{"uptime_seconds":906,"uptime_human":"0 hours, 15 minutes","current_time":"2026-04-15T08:01:06.185323+00:00","timezone":"UTC"},"request":{"client_ip":"127.0.0.1","user_agent":"curl/8.7.1","method":"GET","path":"/"},"endpoints":[{"path":"/","method":"GET","description":"Service information"},{"path":"/health","method":"GET","description":"Health check"}]}%                      arthur@Artur-MacBook-Pro DevOps-Core-Course % curl localhost:808 
curl: (7) Failed to connect to localhost port 808 after 0 ms: Couldn't connect to server
arthur@Artur-MacBook-Pro DevOps-Core-Course % curl localhost:8081
{"service":{"name":"devops-info-service","version":"1.0.0","description":"DevOps course info service","framework":"FastAPI"},"system":{"hostname":"myapp-myapp-1","platform":"Linux","platform_version":"Linux-6.10.14-linuxkit-aarch64-with-glibc2.41","architecture":"aarch64","cpu_count":8,"python_version":"3.13.12"},"runtime":{"uptime_seconds":910,"uptime_human":"0 hours, 15 minutes","current_time":"2026-04-15T08:01:10.164551+00:00","timezone":"UTC"},"request":{"client_ip":"127.0.0.1","user_agent":"curl/8.7.1","method":"GET","path":"/"},"endpoints":[{"path":"/","method":"GET","description":"Service information"},{"path":"/health","method":"GET","description":"Health check"}]}%                      arthur@Artur-MacBook-Pro DevOps-Core-Course % curl localhost:8081
{"service":{"name":"devops-info-service","version":"1.0.0","description":"DevOps course info service","framework":"FastAPI"},"system":{"hostname":"myapp-myapp-1","platform":"Linux","platform_version":"Linux-6.10.14-linuxkit-aarch64-with-glibc2.41","architecture":"aarch64","cpu_count":8,"python_version":"3.13.12"},"runtime":{"uptime_seconds":912,"uptime_human":"0 hours, 15 minutes","current_time":"2026-04-15T08:01:11.282567+00:00","timezone":"UTC"},"request":{"client_ip":"127.0.0.1","user_agent":"curl/8.7.1","method":"GET","path":"/"},"endpoints":[{"path":"/","method":"GET","description":"Service information"},{"path":"/health","method":"GET","description":"Health check"}]}%                      arthur@Artur-MacBook-Pro DevOps-Core-Course % curl localhost:8081/visits
{"visits":5}%                                                                           arthur@Artur-MacBook-Pro DevOps-Core-Course % curl localhost:8080/visits
{"visits":2}%                                                                           arthur@Artur-MacBook-Pro DevOps-Core-Course %  
```

3. **Test Persistence**
   - Note visit counts for each pod
   - Delete one pod (not the StatefulSet)
   - Verify the visit count is preserved after restart

```
arthur@Artur-MacBook-Pro DevOps-Core-Course % kubectl get pods -w
NAME            READY   STATUS    RESTARTS   AGE
myapp-myapp-0   1/1     Running   0          8s
myapp-myapp-1   1/1     Running   0          11s
^C%                                                                                     

arthur@Artur-MacBook-Pro DevOps-Core-Course % curl localhost:8080       
{"service":{"name":"devops-info-service","version":"1.0.0","description":"DevOps course info service","framework":"FastAPI"},"system":{"hostname":"myapp-myapp-0","platform":"Linux","platform_version":"Linux-6.10.14-linuxkit-aarch64-with-glibc2.41","architecture":"aarch64","cpu_count":8,"python_version":"3.13.12"},"runtime":{"uptime_seconds":27,"uptime_human":"0 hours, 0 minutes","current_time":"2026-04-15T09:01:25.243767+00:00","timezone":"UTC"},"request":{"client_ip":"127.0.0.1","user_agent":"curl/8.7.1","method":"GET","path":"/"},"endpoints":[{"path":"/","method":"GET","description":"Service information"},{"path":"/health","method":"GET","description":"Health check"}]}%                        arthur@Artur-MacBook-Pro DevOps-Core-Course % curl localhost:8080/visits
{"visits":1}%                 

arthur@Artur-MacBook-Pro DevOps-Core-Course % kubectl delete pod myapp-myapp-0
pod "myapp-myapp-0" deleted from default namespace

arthur@Artur-MacBook-Pro DevOps-Core-Course % kubectl get pods -w             
NAME            READY   STATUS    RESTARTS   AGE
myapp-myapp-0   1/1     Running   0          6s
myapp-myapp-1   1/1     Running   0          49s

arthur@Artur-MacBook-Pro DevOps-Core-Course % curl localhost:8080/visits
{"visits":1}%
```
