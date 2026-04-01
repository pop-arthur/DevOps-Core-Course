## Lab 11 — Kubernetes Secrets & HashiCorp Vault

--- 

### Task 1 — Kubernetes Secrets Fundamentals 

```
(.venv) arthur@Artur-MacBook-Pro DevOps-Core-Course % kubectl create secret generic app-credentials \
  --from-literal=username=admin \
  --from-literal=password=supersecret
secret/app-credentials created
```

```
(.venv) arthur@Artur-MacBook-Pro DevOps-Core-Course % kubectl get secret app-credentials -o yaml
apiVersion: v1
data:
  password: c3VwZXJzZWNyZXQ=
  username: YWRtaW4=
kind: Secret
metadata:
  creationTimestamp: "2026-04-01T06:01:52Z"
  name: app-credentials
  namespace: default
  resourceVersion: "55281"
  uid: 1752ec71-2963-453c-871b-4c6a8164e6ec
type: Opaque
```

```
(.venv) arthur@Artur-MacBook-Pro DevOps-Core-Course % echo "c3VwZXJzZWNyZXQ=" | base64 -d
supersecret%                           
(.venv) arthur@Artur-MacBook-Pro DevOps-Core-Course % echo "YWRtaW4=" | base64 -d
admin%          
```

Base64 encoding is not encryption. It is a reversible encoding mechanism used for data representation, not protection. Kubernetes Secrets use base64 encoding by default, which means anyone with access to the cluster can decode and read sensitive data. Proper security requires encryption at rest or external secret management solutions like HashiCorp Vault.

--- 

### Task 2 — Helm-Managed Secrets

```
(.venv) arthur@Artur-MacBook-Pro myapp % helm upgrade --install myapp .
Release "myapp" does not exist. Installing it now.
NAME: myapp
LAST DEPLOYED: Wed Apr  1 09:27:34 2026
NAMESPACE: default
STATUS: deployed
REVISION: 1
DESCRIPTION: Install complete
TEST SUITE: None
```

```
(.venv) arthur@Artur-MacBook-Pro myapp % kubectl get pods
NAME                                   READY   STATUS    RESTARTS      AGE
devops-info-service-698cc89c85-4jqc2   1/1     Running   1 (52m ago)   7d
devops-info-service-698cc89c85-gp6kq   1/1     Running   1 (52m ago)   7d
devops-info-service-698cc89c85-zctbl   1/1     Running   2 (52m ago)   7d
myapp-myapp-78c896d65c-9wfwm           1/1     Running   0             25s
myapp-myapp-78c896d65c-qwtzx           1/1     Running   0             12s
```

```
(.venv) arthur@Artur-MacBook-Pro myapp % kubectl exec -it myapp-myapp-78c896d65c-9wfwm -- /bin/sh
$ env | grep username                                     
username=admin
$ env | grep password
password=supersecret
```

--- 

### Task 3 — HashiCorp Vault Integration

added vault
```
(.venv) arthur@Artur-MacBook-Pro DevOps-Core-Course % kubectl get pods -w
NAME                                    READY   STATUS    RESTARTS      AGE
devops-info-service-698cc89c85-4jqc2    1/1     Running   1 (56m ago)   7d
devops-info-service-698cc89c85-gp6kq    1/1     Running   1 (56m ago)   7d
devops-info-service-698cc89c85-zctbl    1/1     Running   2 (56m ago)   7d
myapp-myapp-78c896d65c-9wfwm            1/1     Running   0             4m39s
myapp-myapp-78c896d65c-qwtzx            1/1     Running   0             4m26s
vault-0                                 1/1     Running   0             51s
vault-agent-injector-75998c9b76-rjq52   1/1     Running   0             52s
```

created secret
```
/ $ vault kv put secret/myapp/config \
>   username="admin" \
>   password="supersecret"
====== Secret Path ======
secret/data/myapp/config

======= Metadata =======
Key                Value
---                -----
created_time       2026-04-01T07:11:30.006891676Z
custom_metadata    <nil>
deletion_time      n/a
destroyed          false
version            1
```

enabled auth in k8s and configs
```
/ $ vault auth enable kubernetes
Success! Enabled kubernetes auth method at: kubernetes/

/ $ vault write auth/kubernetes/config \
>   kubernetes_host="https://$KUBERNETES_PORT_443_TCP_ADDR:443"
Success! Data written to: auth/kubernetes/config
/ $ vault policy write myapp-policy - <<EOF
> path "secret/data/myapp/config" {
>   capabilities = ["read"]
> }
> EOF
Success! Uploaded policy: myapp-policy
/ $ vault write auth/kubernetes/role/myapp-role \
>   bound_service_account_names=default \
>   bound_service_account_namespaces=default \
>   policies=myapp-policy \
>   ttl=1h

```

added annotations, restart, check
```
(.venv) arthur@Artur-MacBook-Pro myapp % helm upgrade --install myapp .
Release "myapp" has been upgraded. Happy Helming!
NAME: myapp
LAST DEPLOYED: Wed Apr  1 10:23:00 2026
NAMESPACE: default
STATUS: deployed
REVISION: 5
DESCRIPTION: Upgrade complete
TEST SUITE: None

(.venv) arthur@Artur-MacBook-Pro myapp % kubectl get pods
NAME                                    READY   STATUS    RESTARTS      AGE
devops-info-service-698cc89c85-4jqc2    1/1     Running   1 (82m ago)   7d1h
devops-info-service-698cc89c85-gp6kq    1/1     Running   1 (82m ago)   7d1h
devops-info-service-698cc89c85-zctbl    1/1     Running   2 (82m ago)   7d1h
myapp-myapp-78c896d65c-qwtzx            1/1     Running   0             30m
myapp-myapp-7bd9c57847-l8ssj            1/2     Running   0             4s
myapp-myapp-7bd9c57847-wdkx4            2/2     Running   0             17s
vault-0                                 1/1     Running   0             27m
vault-agent-injector-75998c9b76-rjq52   1/1     Running   0             27m

(.venv) arthur@Artur-MacBook-Pro DevOps-Core-Course % kubectl exec -it myapp-myapp-7bd9c57847-wdkx4 -- /bin/sh
Defaulted container "myapp" out of: myapp, vault-agent, vault-agent-init (init)
$ ls /vault/secrets/
config
$ cat /vault/secrets/config
data: map[password:supersecret username:admin]
metadata: map[created_time:2026-04-01T07:11:30.006891676Z custom_metadata:<nil> deletion_time: destroyed:false version:1]
```

--- 

### Task 4 — Documentation

## 1. Kubernetes Secrets

### Output of creating and viewing your secret

Secret was created using `kubectl create secret generic` and verified via `kubectl get secret -o yaml`, showing base64-encoded values for username and password.

```
(.venv) arthur@Artur-MacBook-Pro DevOps-Core-Course % kubectl create secret generic app-credentials \
  --from-literal=username=admin \
  --from-literal=password=supersecret
secret/app-credentials created
```



### Decoded secret values demonstration

```
(.venv) arthur@Artur-MacBook-Pro DevOps-Core-Course % echo "c3VwZXJzZWNyZXQ=" | base64 -d
supersecret%                           
(.venv) arthur@Artur-MacBook-Pro DevOps-Core-Course % echo "YWRtaW4=" | base64 -d
admin%          
```

### Explanation of base64 encoding vs encryption

Base64 is a reversible encoding mechanism, not encryption. It does not provide security, as anyone with access can decode the data. Kubernetes Secrets are only base64-encoded by default, meaning sensitive data is not protected unless additional mechanisms (e.g., etcd encryption or external tools) are used.

---

## 2. Helm Secret Integration

### Chart structure showing secrets.yaml

A `templates/secrets.yaml` file was added to the Helm chart. It defines a Kubernetes Secret using values from `values.yaml`.

### How secrets are consumed in deployment

Secrets are injected into the container using:

```yaml
envFrom:
  - secretRef:
      name: {{ .Release.Name }}-secret
```

This makes all secret keys available as environment variables.

### Verification output (env vars in pod)

above

---

## Resource Management

### Resource limits configuration

```yaml
resources:
  requests:
    memory: "64Mi"
    cpu: "100m"
  limits:
    memory: "128Mi"
    cpu: "200m"
```

### Explanation of requests vs limits

* **Requests**: minimum resources guaranteed for the container
* **Limits**: maximum resources the container can use

### How to choose appropriate values

Values should be based on application profiling:

* start with small defaults
* monitor usage (CPU/memory)
* adjust to avoid throttling or OOM kills

---

## Vault Integration

### Vault installation verification
above

### Policy and role configuration (sanitized)

* **Policy** allows read access:

```hcl
path "secret/data/myapp/config" {
  capabilities = ["read"]
}
```

* **Role** binds Kubernetes ServiceAccount to policy:

```bash
bound_service_account_names=default
bound_service_account_namespaces=default
policies=myapp-policy
```

---

### Proof of secret injection
above

---

### Explanation of the sidecar injection pattern

Vault uses a **sidecar container (`vault-agent`)** that:

* authenticates with Vault using Kubernetes ServiceAccount
* fetches secrets
* writes them to a shared volume (`/vault/secrets`)

This results in:

* 2 containers in the pod (app + Vault agent)
* secrets dynamically injected at runtime

---

## Security Analysis

### Comparison: K8s Secrets vs Vault

| Feature         | Kubernetes Secrets | Vault                   |
| --------------- | ------------------ | ----------------------- |
| Storage         | etcd (base64)      | External secure storage |
| Encryption      | Optional           | Built-in                |
| Access control  | RBAC               | Policies + roles        |
| Secret delivery | Static             | Dynamic                 |
| Rotation        | Manual             | Automatic               |

---

### When to use each approach

* **Kubernetes Secrets**:

  * simple setups
  * non-critical data
  * development environments

* **Vault**:

  * production systems
  * sensitive credentials
  * dynamic secrets and rotation

---

### Production recommendations

* Enable **etcd encryption at rest**
* Restrict access via **RBAC**
* Avoid storing secrets in Git
* Use **external secret managers (Vault)** for critical systems
* Implement **secret rotation and short TTLs**
