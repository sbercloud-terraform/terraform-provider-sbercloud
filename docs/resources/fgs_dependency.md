---
subcategory: "FunctionGraph"
layout: "sbercloud"
page_title: "SberCloud: sbercloud_fgs_dependency"
description: ""
---

# sbercloud_fgs_dependency

Manages a custom dependency package within SberCloud FunctionGraph.

~> This resource will be deprecated in a future version. Please use `sbercloud_fgs_dependency_version` resource to
replace it. For specific usage instructions, please refer to the corresponding document.

## Example Usage

### Create a custom dependency package using a OBS bucket path where the zip file is located

```hcl
variable "package_name"
variable "package_location"
variable "dependency_name"

resource "sbercloud_obs_bucket" "test" {
  ...
}

resource "sbercloud_obs_bucket_object" "test" {
  bucket = sbercloud_obs_bucket.test.bucket
  key    = format("terraform_dependencies/%s", var.package_name)
  source = var.package_location
}

resource "sbercloud_fgs_dependency" "test" {
  name    = var.dependency_name
  runtime = "Python3.6"
  link    = format("https://%s/%s", sbercloud_obs_bucket.test.bucket_domain_name, sbercloud_obs_bucket_object.test.key)
}
```

## Argument Reference

* `region` - (Optional, String, ForceNew) Specifies the region in which to create a custom dependency package.
  If omitted, the provider-level region will be used. Changing this will create a new resource.

* `runtime` - (Required, String) Specifies the dependency package runtime.
  The valid values are as follows:
  + **Java8**
  + **Java11**
  + **Node.js6.10**
  + **Node.js8.10**
  + **Node.js10.16**
  + **Node.js12.13**
  + **Node.js14.18**
  + **Python2.7**
  + **Python3.6**
  + **Python3.9**
  + **Go1.8**
  + **Go1.x**
  + **C#(.NET Core 2.0)**
  + **C#(.NET Core 2.1)**
  + **C#(.NET Core 3.1)**
  + **PHP7.3**
  + **Custom**
  + **http**

* `name` - (Required, String) Specifies the dependency name.
  The name can contain a maximum of `96` characters and must start with a letter and end with a letter or digit.
  Only letters, digits, underscores (_), periods (.), and hyphens (-) are allowed.

* `link` - (Required, String) Specifies the OBS bucket path where the dependency package is located. The OBS object URL
  must be in zip format, such as `https://obs-terraform.obs.cn-north-4.mysbercloud.com/sbercloudsdkcore.zip`.

-> A link can only be used to create at most one dependency package.

* `description` - (Optional, String) Specifies the dependency description.
  The description can contain a maximum of `512` characters.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The dependency ID in UUID format.

* `owner` - The base64 encoded digest of the dependency after encryption by MD5.

* `etag` - The unique ID of the dependency package.

* `size` - The dependency package size in bytes.

* `version` - The dependency package version.

## Import

Dependencies can be imported using the `id`, e.g.:

```bash
$ terraform import sbercloud_fgs_dependency.test 795e722f-0c23-41b6-a189-dcd56f889cf6
```
