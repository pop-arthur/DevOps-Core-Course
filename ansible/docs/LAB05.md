# Lab 05
---

# 1. Ansible
## 1.1 Ansible Version & Environment

```
(venv) arthur@Artur-MacBook-Pro ansible % ansible --version
ansible [core 2.20.2]
  config file = /Users/arthur/PycharmProjects/DevOps-Core-Course/ansible/ansible.cfg
  configured module search path = ['/Users/arthur/.ansible/plugins/modules', '/usr/share/ansible/plugins/modules']
  ansible python module location = /opt/homebrew/Cellar/ansible/13.3.0/libexec/lib/python3.14/site-packages/ansible
  ansible collection location = /Users/arthur/.ansible/collections:/usr/share/ansible/collections
  executable location = /opt/homebrew/bin/ansible
  python version = 3.14.3 (main, Feb  3 2026, 15:32:20) [Clang 17.0.0 (clang-1700.6.3.2)] (/opt/homebrew/Cellar/ansible/13.3.0/libexec/bin/python)
  jinja version = 3.1.6
  pyyaml version = 6.0.3 (with libyaml v0.2.5)
```

## 1.2 The Target VM

```
(venv) arthur@Artur-MacBook-Pro DevOps-Core-Course % ssh ubuntu@89.169.147.65
Welcome to Ubuntu 22.04.4 LTS (GNU/Linux 5.15.0-117-generic x86_64)
```

## 1.3 How I Organized Everything

```
ansible/
├── roles/                              # All my reusable components live here
│   ├── docker/                         # This one handles server setup
│   │   ├── tasks/
│   │   │   └── main.yml                # Actually installs Docker
│   │   ├── handlers/
│   │   │   └── main.yml                # Restarts Docker when needed
│   │   └── defaults/
│   │       └── main.yml                # Default settings (if any)
│   │
│   └── app_deploy/                      # This one deploys the actual app
│       ├── tasks/
│       │   └── main.yml                  # Logs into Docker Hub, pulls image, runs container
│       └── defaults/
│           └── main.yml                   # Default port, app name, etc.
│
├── group_vars/
│   └── all/
│       ├── main.yml                        # Regular variables (app name, port)
│       └── vault.yml 🔒                     # Encrypted secrets (DHub credentials)
│
├── playbooks/
│   └── deploy.yml                           # The main playbook tying it all together
│
└── hosts.ini                                # Tells Ansible about our VM
```

## 1.4 Why Bother with Roles?

**Without roles:**

- One massive file where Docker installation and app deployment are all mixed together
- Want to use this for another project? Better copy-paste everything
- Something breaks? Good luck finding where in those 500 lines
- Working with a teammate? Hope you don't both edit the same file

**With roles:**

- Docker stuff stays in the `docker` folder, app stuff in `app_deploy` - nice and separate
- I can reuse the `docker` role on any Ubuntu server, takes 2 seconds
- If the app fails to deploy, I know exactly where to look
- Someone else can work on improving the Docker role while I work on the app role

## 1.5 Ping
```
$ ansible all -m ping

lab-vm | SUCCESS => {
    "changed": false,
    "ping": "pong"
}
```
---

# 2. Roles 

## 2.1 Role: `common`
**System preparation and baseline configuration**

### Purpose
Prepares the VM with essential packages and system settings before Docker or any application is installed. Handles apt lock issues and ensures the system is ready for further provisioning.

### Variables
| Variable | Default | Description |
|----------|---------|-------------|
| `common_packages` | `['curl', 'wget', 'git', 'htop', 'ufw', 'fail2ban']` | List of common utilities to install |
| `common_timezone` | `UTC` | System timezone (set to your region if needed) |

### Handlers
*None defined in this role*

### Dependencies
*None - this is the foundation role*


## 2.2 Role: `docker`
**Docker engine installation and configuration**

### Purpose
Installs Docker CE from official repositories, configures the service, and prepares the system to run containers. Handles multi-architecture support (AMD64/ARM64).

### Variables
| Variable | Default | Description |
|----------|---------|-------------|
| `docker_packages` | `['docker-ce', 'docker-ce-cli', 'containerd.io', 'docker-buildx-plugin', 'docker-compose-plugin']` | Docker packages to install |
| `docker_user` | `ubuntu` | User to add to docker group (for running containers without sudo) |

### Handlers
| Name | Trigger | Action |
|------|---------|--------|
| `restart docker` | Docker package installation | Restarts Docker daemon to apply changes |

### Dependencies
- Depends on: `common` (for apt prerequisites)


## 2.3 Role: `app_deploy`
**Application container deployment and health verification**

### Purpose
Logs into Docker Hub, pulls the application image, runs the container with proper configuration, and performs health checks to verify successful deployment.

### Variables
| Variable | Default | Description |
|----------|---------|-------------|
| `dockerhub_username` | *required* | Docker Hub username |
| `dockerhub_password` | *required* | Docker Hub password/token |
| `docker_image` | `"{{ dockerhub_username }}/devops-app"` | 
| `docker_image_tag` | `latest` | Image tag to deploy |
| `app_container_name` | `devops-app` | Container name on host |
| `app_port` | `5000` | Port to expose |
| `app_env` | `{}` | Environment variables dict |
| `app_restart_policy` | `always` | Container restart policy |

### Handlers
| Name | Trigger | Action |
|------|---------|--------|
| `restart app container` | Configuration change | Restarts the application container |

### Dependencies
- Depends on: `docker` (for Docker daemon and Python SDK)
- Depends on: `common` (indirectly through docker)

---

# 3. Idempotency Demonstration
## 3.1 First run
```
(venv) arthur@Artur-MacBook-Pro ansible % ansible-playbook playbooks/provision.yml

PLAY [Provision web servers] ****************************************************************************************************************

TASK [Gathering Facts] **********************************************************************************************************************
ok: [lab-vm]

TASK [common : Wait for apt lock to be released] ********************************************************************************************
ok: [lab-vm]

TASK [common : Update apt cache] ************************************************************************************************************
changed: [lab-vm]

TASK [common : Install common packages] *****************************************************************************************************
changed: [lab-vm]

TASK [common : Set timezone (idempotent)] ***************************************************************************************************
ok: [lab-vm]

TASK [docker : Install prerequisites for Docker repo] ***************************************************************************************
ok: [lab-vm]

TASK [docker : Create keyrings directory] ***************************************************************************************************
ok: [lab-vm]

TASK [docker : Add Docker GPG key] **********************************************************************************************************
changed: [lab-vm]

TASK [docker : Add Docker repository] *******************************************************************************************************
[WARNING]: Deprecation warnings can be disabled by setting `deprecation_warnings=False` in ansible.cfg.
[DEPRECATION WARNING]: INJECT_FACTS_AS_VARS default to `True` is deprecated, top-level facts will not be auto injected after the change. This feature will be removed from ansible-core version 2.24.
Origin: /Users/arthur/PycharmProjects/DevOps-Core-Course/ansible/roles/docker/tasks/main.yml:25:11

23 - name: Add Docker repository
24   ansible.builtin.apt_repository:
25     repo: >-
             ^ column 11

Use `ansible_facts["fact_name"]` (no `ansible_` prefix) instead.

changed: [lab-vm]

TASK [docker : Update apt cache after adding Docker repo] ***********************************************************************************
ok: [lab-vm]

TASK [docker : Install Docker] **************************************************************************************************************
changed: [lab-vm]

TASK [docker : Ensure Docker service is enabled and running] ********************************************************************************
ok: [lab-vm]

TASK [docker : Add user to docker group] ****************************************************************************************************
changed: [lab-vm]

TASK [docker : Install Python Docker SDK for Ansible modules] *******************************************************************************
changed: [lab-vm]

RUNNING HANDLER [docker : restart docker] ***************************************************************************************************
changed: [lab-vm]

PLAY RECAP **********************************************************************************************************************************
lab-vm                     : ok=15   changed=8    unreachable=0    failed=0    skipped=0    rescued=0    ignored=0
```

## 3.2 Second run
```
(venv) arthur@Artur-MacBook-Pro ansible % ansible-playbook playbooks/provision.yml

PLAY [Provision web servers] ****************************************************************************************************************

TASK [Gathering Facts] **********************************************************************************************************************
ok: [lab-vm]

TASK [common : Wait for apt lock to be released] ********************************************************************************************
ok: [lab-vm]

TASK [common : Update apt cache] ************************************************************************************************************
ok: [lab-vm]

TASK [common : Install common packages] *****************************************************************************************************
ok: [lab-vm]

TASK [common : Set timezone (idempotent)] ***************************************************************************************************
ok: [lab-vm]

TASK [docker : Install prerequisites for Docker repo] ***************************************************************************************
ok: [lab-vm]

TASK [docker : Create keyrings directory] ***************************************************************************************************
ok: [lab-vm]

TASK [docker : Add Docker GPG key] **********************************************************************************************************
ok: [lab-vm]

TASK [docker : Add Docker repository] *******************************************************************************************************
[WARNING]: Deprecation warnings can be disabled by setting `deprecation_warnings=False` in ansible.cfg.
[DEPRECATION WARNING]: INJECT_FACTS_AS_VARS default to `True` is deprecated, top-level facts will not be auto injected after the change. This feature will be removed from ansible-core version 2.24.
Origin: /Users/arthur/PycharmProjects/DevOps-Core-Course/ansible/roles/docker/tasks/main.yml:25:11

23 - name: Add Docker repository
24   ansible.builtin.apt_repository:
25     repo: >-
             ^ column 11

Use `ansible_facts["fact_name"]` (no `ansible_` prefix) instead.

ok: [lab-vm]

TASK [docker : Update apt cache after adding Docker repo] ***********************************************************************************
ok: [lab-vm]

TASK [docker : Install Docker] **************************************************************************************************************
ok: [lab-vm]

TASK [docker : Ensure Docker service is enabled and running] ********************************************************************************
ok: [lab-vm]

TASK [docker : Add user to docker group] ****************************************************************************************************
ok: [lab-vm]

TASK [docker : Install Python Docker SDK for Ansible modules] *******************************************************************************
ok: [lab-vm]

PLAY RECAP **********************************************************************************************************************************
lab-vm                     : ok=14   changed=0    unreachable=0    failed=0    skipped=0    rescued=0    ignored=0

```


## 3.3 Analysis

## 📊 What Changed: First Run vs Second Run

| **Task** | **First Run** | **Second Run** | **Why?** |
|----------|---------------|----------------|----------|
| Update apt cache |  **changed** |  ok | First run needed updates |
| Install common packages |  **changed** |  ok | Packages installed first time |
| Add Docker GPG key |  **changed** |  ok | Key added first time |
| Add Docker repository |  **changed** |  ok | Repo added first time |
| Install Docker |  **changed** |  ok | Docker installed first time |
| Add user to docker group |  **changed** |  ok | User added first time |
| Install Python Docker SDK |  **changed** |  ok | SDK installed first time |
| restart docker handler |  **changed** | (not triggered) | No changes = no restart |
| **TOTAL** | **8 changes** | **0 changes** | **Perfect idempotency!** 🎯 |

## 🧠 What Makes Our Roles Idempotent?

### 1. **Built-in Ansible Modules** (They're smart by design)

```yaml
# apt module - checks if package exists before installing
- name: Install common packages
  ansible.builtin.apt:
    name: "{{ common_packages }}"
    state: present  # "present" = install only if missing
```

```yaml
# user module - checks if user already in group
- name: Add user to docker group
  ansible.builtin.user:
    name: "{{ docker_user }}"
    groups: docker
    append: true  # "append: true" = add only if not already there
```

### 2. **Handlers Only Trigger on Actual Changes**

```yaml
# Handler only runs if Docker was *really* installed/changed
- name: Install Docker
  ansible.builtin.apt:
    name: "{{ docker_packages }}"
    state: present
  notify: restart docker  # Notifies BUT only if task changed
```

### 3. **Clever Task Design**

```yaml
# Wait for apt lock - always runs but never "changes"
- name: Wait for apt lock
  ansible.builtin.shell: while fuser...  # changed_when: false
  changed_when: false  # Tells Ansible: "This never changes state"
```

```yaml
# Timezone - only runs if timezone differs
- name: Set timezone
  ansible.builtin.command: "timedatectl set-timezone {{ common_timezone }}"
  changed_when: false  # Idempotent by design - command does nothing if already set
```

### 4. **Cache Validation Prevents Unnecessary Work**

```yaml
# Only updates apt cache if older than 1 hour
- name: Update apt cache
  ansible.builtin.apt:
    update_cache: true
    cache_valid_time: 3600  # 1 hour = no unnecessary updates
```

### 5. **Repository Addition is Idempotent**

```yaml
# apt_repository module checks if repo exists before adding
- name: Add Docker repository
  ansible.builtin.apt_repository:
    repo: "..."
    state: present  # "present" = add only if missing
```

---

# 4. Ansible Vault


### How I store credentials securely

I store all secrets in `group_vars/all.yml` and encrypt them with Ansible Vault:

```bash
ansible-vault create group_vars/all.yml
ansible-vault edit group_vars/all.yml
```

### Vault password management strategy

**Local development:**
```bash
ansible-playbook playbooks/deploy.yml --ask-vault-pass
```

**What's in vault vs regular vars:**

`group_vars/all.yml` (public):
```yaml
dockerhub_username: \***
dockerhub_password: \***

app_name: devops-app
app_port: 5000
docker_image_tag: latest
```

### Example showing it's encrypted

```bash
$ANSIBLE_VAULT;1.1;AES256
36306361633239663463666464336163356161336664663831373364363734343234613663323138
6130623735363732316135316131303838666235353461360a313162663132633661333835653135
63363865623630613332316238376662623065646663393561343830393665336461626334393835
3662393230313838620a383838323431336564356138323736636364376665643161613562356239
34333236326662373331323633663437393463356264613035636336636237396264633762353132
39623664303163353334313633303639653438666466333335363037613037613166343031393934
61303961343738303937386464633230383933316665346665383138396131633339323330636461
62383663653833326231303233396339353233383635613236316466343662623138326631303966
38653532336564633133626665376439386562303265316261353062616335663839363131633661
33653034383965373130333462386262343434346434336562366164363638376237353633303337
37373063313864376136653733396461616635363031633563646436363434393933623136343361
66323764386138643465353235303536366165313735616565393966666265363462346463626630
30313231333834336231623362346166323432306433666339333933623932346661393936653033
65363163373362623733336162326661633439393030396335383730653161396630326239333039
35396162633131316165643364666231396231613963383231393062623865306537643838616439
63373061376538646632336364343331336133613937643864333130613035353936383962663133
66653261396234643836363331393263633965613731346165343365306136373961616130623330
3739613632356239336565616166653462623833636131373261

```

### Why Ansible Vault is important

**1. Prevents credential leaks**
- Without Vault: "Oops, I committed my password to GitHub" 😱
- With Vault: Encrypted file in repo is useless without the password

**2. Separates secrets from configuration**
- Public vars go in `main.yml` (app port, image tag)
- Secrets go in `vault.yml` (passwords, tokens, API keys)
- Team members know what's sensitive vs what's safe

**3. Enables secure collaboration**
- Encrypted file lives in git with the code
- Team shares vault password separately (LastPass, 1Password, etc.)
- New devs get code + password = everything just works

**4. Compliance and auditing**
- You know exactly what's encrypted
- Easy to rotate credentials (just edit vault file)
- No "find and replace" hunting for hardcoded passwords


# 5. Deployment Verification

## 5.1 Terminal output from deploy.yml run
(venv) arthur@Artur-MacBook-Pro ansible % ansible-playbook playbooks/deploy.yml --ask-vault-pass

Vault password:

PLAY [Deploy application] *******************************************************************************************************************

TASK [Gathering Facts] **********************************************************************************************************************
ok: [lab-vm]

TASK [app_deploy : Ensure community.docker collection is available (control node hint)] *****************************************************
ok: [lab-vm] => {
    "msg": "If docker_* modules fail: run 'ansible-galaxy collection install community.docker' on your Mac."
}

TASK [app_deploy : Login to Docker Hub] *****************************************************************************************************
ok: [lab-vm]

TASK [app_deploy : Pull image] **************************************************************************************************************
ok: [lab-vm]

TASK [app_deploy : Remove existing container if any] ****************************************************************************************
changed: [lab-vm]

TASK [app_deploy : Run container] ***********************************************************************************************************
changed: [lab-vm]

TASK [app_deploy : Wait for app port] *******************************************************************************************************
ok: [lab-vm]

TASK [app_deploy : Health check] ************************************************************************************************************
ok: [lab-vm]

TASK [app_deploy : Show health response] ****************************************************************************************************
ok: [lab-vm] => {
    "healthcheck.content": "{\"status\":\"healthy\",\"timestamp\":\"2026-02-25T03:59:59.169862+00:00\",\"uptime_seconds\":13}"
}

RUNNING HANDLER [app_deploy : restart app container] ****************************************************************************************
changed: [lab-vm]

PLAY RECAP **********************************************************************************************************************************
lab-vm                     : ok=10   changed=3    unreachable=0    failed=0    skipped=0    rescued=0    ignored=0

## 5.2 Container status: docker ps output
ubuntu@fhmqqgnfe74gev2b72hj:~$ docker ps
CONTAINER ID   IMAGE                                  COMMAND           CREATED              STATUS                        PORTS                    NAMES
ee2afb843c45   poparthur/devops-info-service:latest   "python app.py"   About a minute ago   Up About a minute (healthy)   0.0.0.0:5000->5000/tcp   devops-info-service


## 5.3 Health check verification: curl outputs
(venv) arthur@Artur-MacBook-Pro ansible % curl http://89.169.147.65:5000/health
{"status":"healthy","timestamp":"2026-02-25T04:02:09.299939+00:00","uptime_seconds":120}%

---
# 6. Key Decisions

**Why use roles instead of plain playbooks?**
- Roles organize code into reusable components
- Clear separation of concerns instead of one many-line playbook

**How do roles improve reusability?**
- `docker` role can provision any Ubuntu server, not just this project's VM
- if next project needs Docker, we can drop in the same role—no copy-pasting, no rewriting

**What makes a task idempotent?**
A task is idempotent when running it 100 times produces the same result as running it once

**How do handlers improve efficiency?**
- Handlers only run when notified AND only if a task actually made changes
- Docker only restarts when the package was *really* installed, not every time the playbook runs—saving time and avoiding unnecessary service interruptions.

**Why is Ansible Vault necessary?**
- Without Vault, secrets live in plain text—one accidental `git push` away from disaster
- => Vault encrypts credentials so they can live safely in version control
