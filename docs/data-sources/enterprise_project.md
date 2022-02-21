---
subcategory: "Enterprise Project Management Service (EPS)"
---

# sbercloud_enterprise_project

Use this data source to get an enterprise project from SberCloud

## Example Usage

```hcl
data "sbercloud_enterprise_project" "test" {
  name = "test"
}
```

## Resources Supported Currently

<!-- markdownlint-disable MD033 -->
Service Name | Resource Name | Sub Resource Name
---- | --- | ---
AS  | sbercloud_as_group |
CBR | sbercloud_cbr_vault |
CCE | sbercloud_cce_cluster | sbercloud_cce_node<br>sbercloud_cce_node_pool
CDM | sbercloud_cdm_cluster |
CES | sbercloud_ces_alarmrule |
DCS | sbercloud_dcs_instance |
DDS | sbercloud_dds_instance |
DMS | sbercloud_dms_kafka_instance<br>sbercloud_dms_rabbitmq_instance |
DNS | sbercloud_dns_ptrrecord<br>sbercloud_dns_zone |
ECS | sbercloud_compute_instance |
EIP | sbercloud_vpc_eip<br>sbercloud_vpc_bandwidth |
ELB | sbercloud_lb_loadbalancer |
EVS | sbercloud_evs_volume |
FGS | sbercloud_fgs_function |
IMS | sbercloud_images_image |
NAT | sbercloud_nat_gateway | sbercloud_nat_snat_rule<br>sbercloud_nat_dnat_rule
OBS | sbercloud_obs_bucket | sbercloud_obs_bucket_object<br>sbercloud_obs_bucket_policy
RDS | sbercloud_rds_instance<br>sbercloud_rds_read_replica_instance |
SFS | sbercloud_sfs_file_system<br>sbercloud_sfs_turbo | sbercloud_sfs_access_rule
VPC | sbercloud_vpc<br>sbercloud_networking_secgroup | sbercloud_vpc_subnet<br>sbercloud_vpc_route<br>sbercloud_networking_secgroup_rule
<!-- markdownlint-enable MD033 -->

## Argument Reference

* `name` - (Optional, String) Specifies the enterprise project name. Fuzzy search is supported.

* `id` - (Optional, String) Specifies the ID of an enterprise project. The value 0 indicates enterprise project default.

* `status` - (Optional, Int) Specifies the status of an enterprise project.
    + 1 indicates Enabled.
    + 2 indicates Disabled.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `description` - Provides supplementary information about the enterprise project.

* `created_at` - Specifies the time (UTC) when the enterprise project was created. Example: 2018-05-18T06:49:06Z

* `updated_at` - Specifies the time (UTC) when the enterprise project was modified. Example: 2018-05-28T02:21:36Z
