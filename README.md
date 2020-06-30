Terraform SberCloud Provider
==============================

- Website: https://www.terraform.io
- [![Gitter chat](https://badges.gitter.im/hashicorp-terraform/Lobby.png)](https://gitter.im/hashicorp-terraform/Lobby)
- Mailing list: [Google Groups](http://groups.google.com/group/terraform-tool)

<img src="https://cdn.rawgit.com/hashicorp/terraform-website/master/content/source/assets/images/logo-hashicorp.svg" width="600px">

Requirements
------------

-	[Terraform](https://www.terraform.io/downloads.html) 0.10+
-	[Go](https://golang.org/doc/install) 1.11 (to build the provider plugin)

Building The Provider
---------------------

Clone repository to: `$GOPATH/src/github.com/huaweicloud/terraform-provider-sbercloud`

```sh
$ mkdir -p $GOPATH/src/github.com/huaweicloud; cd $GOPATH/src/github.com/huaweicloud
$ git clone https://github.com/huaweicloud/terraform-provider-sbercloud
```

Enter the provider directory and build the provider

```sh
$ cd $GOPATH/src/github.com/huaweicloud/terraform-provider-sbercloud
$ make build
```

## Exact steps on clean Ubuntu 16.04

```sh
# prerequisites are sudo privileges, unzip, make, wget and git.  Use apt install if missing.
$ wget https://storage.googleapis.com/golang/go1.11.1.linux-amd64.tar.gz
$ sudo tar -C /usr/local -xzf go1.11.1.linux-amd64.tar.gz
$ export PATH=$PATH:/usr/local/go/bin # You should put in your .profile or .bashrc
$ go version # to verify it runs and version #
$ go get github.com/huaweicloud/terraform-provider-sbercloud
$ cd ~/go/src/github.com/huaweicloud/terraform-provider-sbercloud
$ make build
$ export PATH=$PATH:~/go/bin # You should put in your .profile or .bashrc
# compatible with terraform 0.12
$ wget https://releases.hashicorp.com/terraform/0.12.0/terraform_0.12.0_linux_amd64.zip
$ unzip terraform_0.12.0_linux_amd64.zip
$ mv terraform ~/go/bin/
$ terraform version # to verify it runs and version #
$ vi test.tf # paste in Quick Start contents, fix authentication information
$ terraform init
$ terraform plan
$ terraform apply # Should all work if everything is correct.

```

## Quick Start


```hcl
# Configure the SberCloud Provider with Username/Password
# This will work with a single defined/default network, otherwise you need to specify network
# to fix errors about multiple networks found.
provider "sbercloud" {
  user_name   = "user"
  region      = "cn-north-1"
  domain_name = "domain"
  password    = "pwd"
}

# Create a web server
resource "huaweicloud_compute_instance_v2" "test-server" {
  name		        = "test-server"
  image_name        = "Standard_CentOS_7_latest"
  flavor_name       = "s1.medium"
  availability_zone = "cn-north-1a"
}
```

Using the provider
----------------------
Please see the documentation at [provider usage](website/docs/index.html.markdown).

Or you can browse the documentation within this repo [here](https://github.com/huaweicloud/terraform-provider-sbercloud/tree/master/website/docs).

Developing the Provider
---------------------------

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (version 1.11+ is *required*). You'll also need to correctly setup a [GOPATH](http://golang.org/doc/code.html#GOPATH), as well as adding `$GOPATH/bin` to your `$PATH`.

To compile the provider, run `make build`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

```sh
$ make build
...
$ $GOPATH/bin/terraform-provider-sbercloud
...
```

In order to test the provider, you can simply run `make test`.

```sh
$ make test
```

In order to run the full suite of Acceptance tests, run `make testacc`.

*Note:* Acceptance tests create real resources, and often cost money to run.

```sh
$ make testacc
```

## License

Terraform-Provider-Huaweicloud is under the Mozilla Public License 2.0. See the [LICENSE](LICENSE) file for details.

