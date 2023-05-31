# CCE Addon Templates

Addon support configuration input depending on addon type and version. This page contains description of addon
arguments. You can get up to date reference of addon arguments for your cluster using data source
[`sbercloud_cce_addon_template`](https://registry.terraform.io/providers/sbercloud-terraform/sbercloud/latest/docs/data-sources/cce_addon_template)
.

Following addon templates exist in the addon template list:

- [`autoscaler`](#autoscaler)
- [`coredns`](#coredns)
- [`everest`](#everest)
- [`metrics-server`](#metrics-server)
- [`gpu-beta`](#gpu-beta)
- [`nginx-ingress`](#nginx-ingress)

All addons accept `basic` and some can accept `custom`, `flavor` input values.
It is recommended to use `basic_json`, `custom_json` and `flavor_json` for more flexible input.

## Example Usage

### Use basic_json, custom_json and flavor_json

```hcl
variable "cluster_id" {}
variable "tenant_id" {}

data "sbercloud_cce_addon_template" "autoscaler" {
  cluster_id = var.cluster_id
  name       = "autoscaler"
  version    = "1.19.6"
}

resource "sbercloud_cce_addon" "autoscaler" {
  cluster_id    = var.cluster_id
  template_name = "autoscaler"
  version       = "1.19.6"

  values {
    basic_json  = jsonencode(jsondecode(data.sbercloud_cce_addon_template.autoscaler.spec).basic)
    custom_json = jsonencode(merge(
      jsondecode(data.sbercloud_cce_addon_template.autoscaler.spec).parameters.custom,
      {
        cluster_id = var.cluster_id
        tenant_id  = var.tenant_id
      }
    ))
    flavor_json = jsonencode(jsondecode(data.sbercloud_cce_addon_template.autoscaler.spec).parameters.flavor2)
  }
}

```

### Use basic and custom

```hcl
variable "cluster_id" {}
variable "tenant_id" {}

data "sbercloud_cce_addon_template" "autoscaler" {
  cluster_id = var.cluster_id
  name       = "autoscaler"
  version    = "1.19.6"
}

resource "sbercloud_cce_addon" "autoscaler" {
  cluster_id    = var.cluster_id
  template_name = "autoscaler"
  version       = "1.19.6"

  values {
    basic  = jsondecode(data.sbercloud_cce_addon_template.autoscaler.spec).basic
    custom = merge(
      jsondecode(data.sbercloud_cce_addon_template.autoscaler.spec).parameters.custom,
      {
        cluster_id = var.cluster_id
        tenant_id  = var.tenant_id
      }
    )
  }
}

```

## Addon Inputs

### autoscaler

A component that automatically adjusts the size of a Kubernetes Cluster so that all pods have a place to run and there
are no unneeded nodes.
`template_version`: `1.19.1`

#### basic

```json
{
  "cceEndpoint": "https://cce.ru-moscow-1.mysbercloud.com",
  "ecsEndpoint": "https://ecs.ru-moscow-1.mysbercloud.com",
  "image_version": "1.19.6",
  "platform": "linux-amd64",
  "region": "ru-moscow-1",
  "swr_addr": "swr.ru-moscow-1.mysbercloud.com",
  "swr_user": "hwofficial"
}
```

#### custom

```json
{
  "cluster_id": "",
  "coresTotal": 32000,
  "expander": "priority",
  "logLevel": 4,
  "maxEmptyBulkDeleteFlag": 10,
  "maxNodeProvisionTime": 15,
  "maxNodesTotal": 1000,
  "memoryTotal": 128000,
  "scaleDownDelayAfterAdd": 10,
  "scaleDownDelayAfterDelete": 10,
  "scaleDownDelayAfterFailure": 3,
  "scaleDownEnabled": false,
  "scaleDownUnneededTime": 10,
  "scaleDownUtilizationThreshold": 0.5,
  "scaleUpCpuUtilizationThreshold": 1,
  "scaleUpMemUtilizationThreshold": 1,
  "scaleUpUnscheduledPodEnabled": true,
  "scaleUpUtilizationEnabled": true,
  "tenant_id": "",
  "unremovableNodeRecheckTimeout": 5
}
```

### coredns

CoreDNS is a DNS server that chains plugins and provides Kubernetes DNS Services.
`template_version`: `1.17.7`

#### basic

```json
{
  "cluster_ip": "10.247.3.10",
  "image_version": "1.17.7",
  "platform": "linux-amd64",
  "swr_addr": "swr.ru-moscow-1.mysbercloud.com",
  "swr_user": "hwofficial"
}
```

#### custom

```json
{
  "stub_domains": "",
  "upstream_nameservers": ""
}
```

### everest

Everest is a cloud native container storage system based on CSI, used to support cloud storages services for Kubernetes.
`template_version`: `1.2.9`

#### basic

```json
{
  "bms_url": "bms.ru-moscow-1.mysbercloud.com",
  "controller_image_version": "1.2.9",
  "driver_image_version": "1.2.9",
  "ecsEndpoint": "https://ecs.ru-moscow-1.mysbercloud.com",
  "evs_url": "evs.ru-moscow-1.mysbercloud.com",
  "iam_url": "iam.ru-moscow-1.mysbercloud.com",
  "ims_url": "ims.ru-moscow-1.mysbercloud.com",
  "obs_url": "obs.ru-moscow-1.mysbercloud.com",
  "platform": "linux-amd64",
  "sfs_turbo_url": "sfs-turbo.ru-moscow-1.mysbercloud.com",
  "sfs_url": "sfs.ru-moscow-1.mysbercloud.com",
  "supportHcs": false,
  "swr_addr": "swr.ru-moscow-1.mysbercloud.com",
  "swr_user": "hwofficial"
}
```

#### custom

```json
{
  "cluster_id": "",
  "default_vpc_id": "",
  "disable_auto_mount_secret": false,
  "project_id": ""
}
```

### metrics-server

Metrics Server is a cluster-level resource usage data aggregator.
`template_version`: `1.1.2`

#### basic

```json
{
  "image_version": "v0.4.4",
  "swr_addr": "swr.ru-moscow-1.mysbercloud.com",
  "swr_user": "hwofficial"
}
```

#### custom

The custom block is *not supported*.

### gpu-beta

A device plugin for nvidia.com/gpu resource on nvidia driver.
`template_version`: `1.2.2`

#### basic

```json
{
  "device_version": "1.2.2",
  "driver_version": "1.2.2",
  "obs_url": "obs.ru-moscow-1.mysbercloud.com",
  "region": "ru-moscow-1",
  "swr_addr": "swr.ru-moscow-1.mysbercloud.com",
  "swr_user": "hwofficial"
}
```

#### custom

```json
{
  "is_driver_from_nvidia": true,
  "nvidia_driver_download_url": ""
}
```

### nginx-ingress

#### basic

```json
{
  "swr_addr": "swr.ru-moscow-1.mysbercloud.com",
  "swr_user": "hwofficial",
  "tag": "v0.46.0"
}
```

#### custom

```json
{
  "config": {},
  "defaultBackend": {
    "enabled": true
  },
  "defaultBackendService": "",
  "headers": {},
  "service": {
    "annotations": {},
    "loadBalancerIP": ""
  },
  "tcp": {},
  "udp": {}
}
```

#### flavor

Should be passed as `flavor_json`.

```json
{
  "description": "custom resources",
  "name": "custom-resources",
  "replicas": 2,
  "resources": [
    {
      "limitsCpu": "8000m",
      "limitsMem": "4000Mi",
      "name": "nginx-ingress",
      "requestsCpu": "8000m",
      "requestsMem": "4000Mi"
    }
  ]
}
```

Overall structure looks like this:

```json
{
  "basic":
  {
    "swr_addr": "swr.ru-moscow-1.mysbercloud.com",
    "swr_user": "hwofficial",
    "tag": "v0.46.0"
  },
  "parameters":
  {
    "custom":
    {
      "config":
      {},
      "defaultBackend":
      {
        "enabled": true
      },
      "defaultBackendService": "",
      "headers":
      {},
      "service":
      {
        "annotations":
        {},
        "loadBalancerIP": ""
      },
      "tcp":
      {},
      "udp":
      {}
    },
    "flavor1":
    {
      "description": "custom resources",
      "name": "custom-resources",
      "replicas": 2,
      "resources":
      [
        {
          "limitsCpu": "8000m",
          "limitsMem": "4000Mi",
          "name": "nginx-ingress",
          "requestsCpu": "8000m",
          "requestsMem": "4000Mi"
        }
      ]
    }
  }
}
```

Friendly reminder, that `flavor_json` should be `jsonencode`d after `merge`.

```hcl
flavor_json = jsonencode(merge(
  jsondecode(data.sbercloud_cce_addon_template.nginx_ingress_template.spec).parameters.flavor1,
  {
    "resources" = [
      {
        "limitsCpu": "200m",
        "limitsMem": "512Mi",
        "requestsCpu": "200m",
        "requestsMem": "256Mi"
      }
    ]
  }
))
```
