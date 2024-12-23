Terraform SberCloud Provider
==============================

- Website: https://www.terraform.io
- [![Gitter chat](https://badges.gitter.im/hashicorp-terraform/Lobby.png)](https://gitter.im/hashicorp-terraform/Lobby)
- Mailing list: [Google Groups](http://groups.google.com/group/terraform-tool)

Requirements
------------

-	[Terraform](https://www.terraform.io/downloads.html) 0.12.x
-	[Go](https://golang.org/doc/install) 1.18 (to build the provider plugin)

Building The Provider
---------------------

Clone repository to: `$GOPATH/src/github.com/sbercloud-terraform/terraform-provider-sbercloud`

```sh
mkdir -p $GOPATH/src/github.com/sbercloud-terraform; cd $GOPATH/src/github.com/sbercloud-terraform
git clone https://github.com/sbercloud-terraform/terraform-provider-sbercloud
```

Enter the provider directory and build the provider

```sh
cd $GOPATH/src/github.com/sbercloud-terraform/terraform-provider-sbercloud
make build
```
After this step, change your .terraformrc file like provided below:
```
provider_installation {
  dev_overrides {
    "mycloud.com/myorg/sbercloud" = "/Users/user/go/bin/"
  }
  direct {}
}
```
In your .tf files add provider initialization like shown below:
```terraform
terraform {
  required_providers {
    sbercloud = {
      source  = "mycloud.com/myorg/sbercloud"
    }
  }
}
provider "sbercloud" {
  auth_url = "https://iam.ru-moscow-1.hc.sbercloud.ru/v3" 
  region   = "ru-moscow-1" 
  access_key = var.access_key
  secret_key = var.secret_key
}
```
Add file variables.tf to same directory, where your .tf file placed.
```terraform
variable "access_key" {
        default = "access-key-id"
}

variable "secret_key" {
        default = "secret-access-key"
}
```
Run `terraform plan` (skip `terraform init`) 

Using the provider
----------------------
Please see the documentation at [Terraform Registry](https://registry.terraform.io/providers/sbercloud-terraform/sbercloud/latest/docs).

Or you can browse the documentation within this repo [here](https://github.com/sbercloud-terraform/terraform-provider-sbercloud/tree/master/docs).

Developing the Provider
---------------------------

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (version 1.11+ is *required*). You'll also need to correctly setup a [GOPATH](http://golang.org/doc/code.html#GOPATH), as well as adding `$GOPATH/bin` to your `$PATH`.

To compile the provider, run `make build`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

```sh
make build
...
$GOPATH/bin/terraform-provider-sbercloud
...
```

In order to test the provider, you can simply run `make test`.

```sh
make test
```

In order to run the full suite of Acceptance tests, run `make testacc`.

*Note:* Acceptance tests create real resources, and often cost money to run.

```sh
make testacc
```

## License

Terraform-Provider-Sbercloud is under the Mozilla Public License 2.0. See the [LICENSE](LICENSE) file for details.

