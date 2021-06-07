#
# Declare all required input variables
#
variable "access_key" {
  description = "Access Key to access SberCloud"
  sensitive   = true
}

variable "secret_key" {
  description = "Secret Key to access SberCloud"
  sensitive   = true
}

variable "iam_project_name" {
  description = "IAM project where to deploy infrastructure"
}

#
# Configure the SberCloud Provider
#
terraform {
  required_providers {
    sbercloud = {
      source = "sbercloud-terraform/sbercloud"
    }
  }
}

provider "sbercloud" {
  auth_url = "https://iam.ru-moscow-1.hc.sbercloud.ru/v3"
  region   = "ru-moscow-1"

  access_key   = var.access_key
  secret_key   = var.secret_key
  project_name = var.iam_project_name
}
