---
subcategory: "Cloud Container Instance (CCI)"
layout: "sbercloud"
page_title: "Sbercloud: sbercloud_cci_pvc"
description: |-
  Manages a CCI v2 persistent volume claim resource within Sbercloud.
---

# sbercloud_cci_pvc

Manages a CCI v2 persistent volume claim resource within Sbercloud.

## Example Usage

```hcl
variable "namespace_name" {}
variable "pvc_name" {}

resource "sbercloud_cci_pvc" "test" {
  namespace          = var.namespace_name
  name               = var.pvc_name
  access_modes       = ["ReadWriteMany"]
  storage_class_name = "csi-obs"
  volume_mode        = "Filesystem"

  annotations = {
    "everest.io/obs-volume-type" = "STANDARD"
    "csi.storage.k8s.io/fstype"  = "s3fs"
  }

  resources {
    requests = {
      storage = "1Gi"
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the persistent volume claim.
  If omitted, the provider-level region will be used. Changing this will create a new resource.

* `namespace` - (Required, String, NonUpdatable) Specifies the namespace of the persistent volume claim.

* `name` - (Required, String, NonUpdatable) Specifies the name of the persistent volume claim.

* `annotations` - (Optional, Map, NonUpdatable) Specifies the annotations of the persistent volume claim.

* `labels` - (Optional, Map, NonUpdatable) Specifies the labels of the persistent volume claim.

* `access_modes` - (Optional, List, NonUpdatable) Specifies the access modes of the persistent volume claim.

* `storage_class_name` - (Optional, String, NonUpdatable) Specifies the storage class name of the persistent volume
  claim.

* `volume_mode` - (Optional, String, NonUpdatable) Specifies the volume mode of the persistent volume claim.

* `valume_name` - (Optional, String, NonUpdatable) Specifies the volume name of the persistent volume claim.

* `resources` - (Optional, List) Specifies the resources of the persistent volume claim.
  The [resources](#cci_pvc_resources) structure is documented below.

* `selector` - (Optional, List) Specifies the selector of the persistent volume claim.
  The [selector](#cci_pvc_selector) structure is documented below.

<a name="cci_pvc_resources"></a>
The `resources` block supports:

* `limits` - (Optional, Map) Specifies the resource limits.

* `requests` - (Optional, Map) Specifies the resource requests, e.g. `storage = "1Gi"`.

<a name="cci_pvc_selector"></a>
The `selector` block supports:

* `match_labels` - (Optional, Map) Specifies the labels to match.

* `match_expressions` - (Optional, List) Specifies the label selector requirements.
  The [match_expressions](#cci_pvc_selector_match_expressions) structure is documented below.

<a name="cci_pvc_selector_match_expressions"></a>
The `match_expressions` block supports:

* `key` - (Optional, String) Specifies the label key.

* `operator` - (Optional, String) Specifies the operator, e.g. **In**, **NotIn**, **Exists** and **DoesNotExist**.

* `values` - (Optional, List) Specifies the label values.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in format `<namespace>/<name>`.

* `creation_timestamp` - The creation timestamp of the persistent volume claim.

* `resource_version` - The resource version of the persistent volume claim.

* `uid` - The uid of the persistent volume claim.

* `api_version` - The API version of the persistent volume claim.

* `kind` - The kind of the persistent volume claim.

* `finalizers` - The finalizers of the persistent volume claim.

* `status` - The status of the persistent volume claim.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes.
* `delete` - Default is 10 minutes.

## Import

CCI persistent volume claims can be imported using their `namespace` and `name`, separated by a slash, e.g.

```bash
$ terraform import sbercloud_cci_pvc.test <namespace>/<name>
```
