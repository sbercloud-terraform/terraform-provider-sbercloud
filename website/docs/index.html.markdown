---
layout: "sbercloud"
page_title: "Provider: SberCloud"
sidebar_current: "docs-sbercloud-index"
description: |-
  The SberCloud provider is used to interact with the many resources supported by SberCloud. The provider needs to be configured with the proper credentials before it can be used.
---

# SberCloud Provider

The SberCloud provider is used to interact with the
many resources supported by SberCloud. The provider needs to be configured
with the proper credentials before it can be used.

Use the navigation to the left to read about the available resources.

## Example Usage

```hcl
# Configure the SberCloud Provider
provider "sbercloud" {
  region       = "ru-moscow-1"
  account_name = "my-account-name"
  user_name    = "my-username"
  password     = "my-password"
}

# Create a user
resource "sbercloud_identity_user_v3" "example" {
  name = "terraform"
  password = "password123@!"
  enabled = true
  description = "created by terraform"
}
```

## Authentication

The Sber Cloud provider offers a flexible means of providing credentials for
authentication. The following methods are supported, in this order, and
explained below:

- Static credentials
- Environment variables

### Static credentials ###

!> **Warning:** Hard-coding credentials into any Terraform configuration is not
recommended, and risks secret leakage should this file ever be committed to a
public version control system.

Static credentials can be provided by adding an `user_name` and `password`
in-line in the provider block:

Usage:

```hcl
provider "sbercloud" {
  region       = "ru-moscow-1"
  account_name = "my-account-name"
  user_name    = "my-username"
  password     = "my-password"
}
```

### Environment variables

You can provide your credentials via the `SBC_USERNAME` and
`SBC_PASSWORD`, environment variables, representing your Sber
Cloud Username and Password, respectively.

```hcl
provider "sbercloud" {}
```

Usage:

```sh
$ export SBC_USERNAME="user-name"
$ export SBC_PASSWORD="password"
$ export SBC_REGION_NAME="ru-moscow-1"
$ export SBC_ACCOUNT_NAME="account-name"
$ terraform plan
```


## Configuration Reference

The following arguments are supported:

* `region` - (Required) This is the Sber Cloud region. It must be provided,
  but it can also be sourced from the `SBC_REGION_NAME` environment variables.

* `account_name` - (Optional, Required for IAM resources) The
  of IAM to scope to. If omitted, the `SBC_ACCOUNT_NAME` environment variable is used.

* `access_key` - (Optional) The access key of the SberCloud to use.
  If omitted, the `SBC_ACCESS_KEY` environment variable is used.

* `secret_key` - (Optional) The secret key of the SberCloud to use.
  If omitted, the `SBC_SECRET_KEY` environment variable is used.

* `user_name` - (Optional) The Username to login with. If omitted, the
  `SBC_USERNAME` environment variable is used.

* `password` - (Optional) The Password to login with. If omitted, the
  `SBC_PASSWORD` environment variable is used.

* `project_name` - (Optional) The Name of the Project to login with.
  If omitted, the `SBC_PROJECT_NAME` environment variable are used.

* `auth_url` - (Optional) The Identity authentication URL. If omitted, the
  `SBC_AUTH_URL` environment variable is used.

* `insecure` - (Optional) Trust self-signed SSL certificates. If omitted, the
  `SBC_INSECURE` environment variable is used.


## Testing and Development

In order to run the Acceptance Tests for development, the following environment
variables must also be set:

* `SBC_REGION_NAME` - The region in which to create resources.

* `SBC_USERNAME` - The username to login with.

* `SBC_PASSWORD` - The password to login with.

* `SBC_ACCOUNT_NAME` - The IAM account name.

You should be able to use any SberCloud environment to develop on as long as the
above environment variables are set.
