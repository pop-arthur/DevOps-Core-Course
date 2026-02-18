variable "folder_id" {
  description = "Yandex Cloud folder id"
  type        = string
}

variable "zone" {
  description = "Yandex Cloud zone"
  type        = string
  default     = "ru-central1-a"
}

variable "instance_name" {
  description = "Name for the compute instance"
  type        = string
  default     = "lab4"
}

variable "ssh_public_key" {
  description = "Your SSH public key content"
  type        = string
  sensitive   = true
}

variable "yc_token" {
  description = "Yandex Cloud OAuth token"
  type        = string
  sensitive   = true
}