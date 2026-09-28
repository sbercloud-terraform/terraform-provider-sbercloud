---
subcategory: "Cloud Container Instance (CCI)"
layout: "sbercloud"
page_title: "Sbercloud: sbercloud_cci_namespace"
description: |-
  Manages a CCI v2 namespace resource within Sbercloud.
---

# sbercloud_cci_namespace

Manages a CCI v2 namespace resource within Sbercloud.

## Example Usage

```hcl
variable "namespace_name" {}

resource "sbercloud_cci_namespace" "test" {
  name = var.namespace_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the CCI namespace resource.
  If omitted, the provider-level region will be used. Changing this will create a new CCI namespace resource.

* `name` - (Required, String, NonUpdatable) Specifies the unique name of the CCI namespace.
  This parameter can contain a maximum of `63` characters, which may consist of lowercase letters, digits and
  hyphens (-), and must start and end with lowercase letters and digits.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, which equals to the namespace name.

* `annotations` - The annotations of the namespace.

* `labels` - The labels of the namespace.

* `creation_timestamp` - The creation timestamp of the namespace.

* `resource_version` - The resource version of the namespace.

* `uid` - The uid of the namespace.

* `api_version` - The API version of the namespace.

* `kind` - The kind of the namespace.

* `finalizers` - The finalizers of the namespace.

* `status` - The status of the namespace.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 5 minutes.
* `delete` - Default is 3 minutes.

## Import

CCI namespaces can be imported using their `name`, e.g.

```bash
$ terraform import sbercloud_cci_namespace.test <name>
```
