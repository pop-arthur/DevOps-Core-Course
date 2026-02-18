provider "github" {
  token = var.github_token
}

resource "github_repository" "course_repo" {
  count         = var.create_github ? 1 : 0
  name          = "DevOps-Core-Course"
  description   = "🚀Production-grade DevOps course: 18 hands-on labs covering Docker, Kubernetes, Helm, Terraform, Ansible, CI/CD, GitOps (ArgoCD), monitoring (Prometheus/Grafana), and more. Build real-world skills with progressive delivery, secrets management, and cloud-native deployments. (Hi from terraform)"
  visibility    = "public"
  has_issues    = true
  has_wiki      = true
  has_projects  = true
}

variable "github_token" {
  description = "GitHub personal access token (set via env or terraform.tfvars)"
  type        = string
  sensitive   = true
  default     = ""
}

variable "create_github" {
  description = "Whether to create GitHub repository via Terraform"
  type        = bool
  default     = true
}
