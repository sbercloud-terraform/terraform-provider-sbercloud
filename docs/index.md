# SberCloud Provider

The SberCloud provider is used to interact with the
many resources supported by SberCloud. The provider needs to be configured
with the proper credentials before it can be used.

Use the navigation to the left to read about the available resources.

## Example Usage

```hcl
# Configure the SberCloud Provider
provider "sbercloud" {
  region     = "ru-moscow-1"
  access_key = "my-access-key"
  secret_key = "my-secret-key"
}

# Create a VPC
resource "sbercloud_vpc" "example" {
  name = "my_vpc"
  cidr = "192.168.0.0/16"
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

Static credentials can be provided by adding an `access_key` and `secret_key`
in-line in the provider block:

Usage:

```hcl
provider "sbercloud" {
  region     = "ru-moscow-1"
  access_key = "my-access-key"
  secret_key = "my-secret-key"
}
```

### Environment variables

You can provide your credentials via the `SBC_ACCESS_KEY` and
`SBC_SECRET_KEY`, environment variables, representing your Sber
Cloud Username and Password, respectively.

```hcl
provider "sbercloud" {}
```

Usage:

```sh
$ export SBC_ACCESS_KEY="user-name"
$ export SBC_SECRET_KEY="password"
$ export SBC_REGION_NAME="ru-moscow-1"
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

* `security_token` - (Optional) The security token to authenticate with a
  [temporary security credential](https://support.hc.sbercloud.ru/en-us/api/iam/en-us_topic_0097949518.html).
  If omitted, the `SBC_SECURITY_TOKEN` environment variable is used.

* `project_name` - (Optional) The Name of the Project to login with.
  If omitted, the `SBC_PROJECT_NAME` environment variable are used.

* `auth_url` - (Optional) The Identity authentication URL. If omitted, the
  `SBC_AUTH_URL` environment variable is used.

* `insecure` - (Optional) Trust self-signed SSL certificates. If omitted, the
  `SBC_INSECURE` environment variable is used.

* `max_retries` - (Optional) This is the maximum number of times an API
  call is retried, in the case where requests are being throttled or
  experiencing transient failures. The delay between the subsequent API
  calls increases exponentially. The default value is `5`.
  If omitted, the `SBC_MAX_RETRIES` environment variable is used.

* `enterprise_project_id` - (Optional) Default Enterprise Project ID for supported resources.
  If omitted, the `SBC_ENTERPRISE_PROJECT_ID` environment variable is used.


## Testing and Development

In order to run the Acceptance Tests for development, the following environment
variables must also be set:

* `SBC_REGION_NAME` - The region in which to create resources.

* `SBC_ACCESS_KEY` - The username to login with.

* `SBC_SECRET_KEY` - The password to login with.


You should be able to use any SberCloud environment to develop on as long as the
above environment variables are set.
