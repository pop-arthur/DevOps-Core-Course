terraform {
  required_version = ">= 1.0.0"
  required_providers {
    yandex = {
      source  = "yandex-cloud/yandex"
      version = ">= 0.87.0"
    }
    github = {
      source  = "integrations/github"
      version = ">= 4.0.0"
    }
  }
}

provider "yandex" {
  token     = var.yc_token # Вместо service_account_key_file
  folder_id = var.folder_id
  zone      = var.zone
}

data "yandex_compute_image" "ubuntu" {
  family = "ubuntu-2204-lts"
}

resource "yandex_vpc_network" "lab_network" {
  name = "lab-network-${var.instance_name}"
}

resource "yandex_vpc_subnet" "lab_subnet" {
  name           = "lab-subnet-${var.instance_name}"
  zone           = var.zone
  network_id     = yandex_vpc_network.lab_network.id
  v4_cidr_blocks = ["10.0.1.0/24"]
}

resource "yandex_compute_instance" "vm" {
  name        = var.instance_name
  platform_id = "standard-v2"
  zone        = var.zone

  resources {
    cores         = 2
    core_fraction = 20
    memory        = 1
  }

  boot_disk {
    initialize_params {
      image_id = data.yandex_compute_image.ubuntu.id
      size     = 10
      type     = "network-hdd"
    }
  }

  network_interface {
    subnet_id = yandex_vpc_subnet.lab_subnet.id
    nat       = true
  }

  metadata = {
    ssh-keys = "ubuntu:${var.ssh_public_key}"
  }
}
