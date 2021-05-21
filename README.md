Terraform SberCloud Provider
==============================

- Website: https://www.terraform.io
- [![Gitter chat](https://badges.gitter.im/hashicorp-terraform/Lobby.png)](https://gitter.im/hashicorp-terraform/Lobby)
- Mailing list: [Google Groups](http://groups.google.com/group/terraform-tool)

<img src="https://cdn.rawgit.com/hashicorp/terraform-website/master/content/source/assets/images/logo-hashicorp.svg" width="600px">

Requirements
------------

-	[Terraform](https://www.terraform.io/downloads.html) 0.12.x
-	[Go](https://golang.org/doc/install) 1.11 (to build the provider plugin)

Building The Provider
---------------------

Clone repository to: `$GOPATH/src/github.com/sbercloud-terraform/terraform-provider-sbercloud`

```sh
$ mkdir -p $GOPATH/src/github.com/sbercloud-terraform; cd $GOPATH/src/github.com/sbercloud-terraform
$ git clone https://github.com/sbercloud-terraform/terraform-provider-sbercloud
```

Enter the provider directory and build the provider

```sh
$ cd $GOPATH/src/github.com/sbercloud-terraform/terraform-provider-sbercloud
$ make build
```

Using the provider
----------------------
Please see the documentation at [Terraform Registry](https://registry.terraform.io/providers/sbercloud-terraform/sbercloud/latest/docs).

Or you can browse the documentation within this repo [here](https://github.com/sbercloud-terraform/terraform-provider-sbercloud/tree/master/docs).

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

Terraform-Provider-Sbercloud is under the Mozilla Public License 2.0. See the [LICENSE](LICENSE) file for details.

