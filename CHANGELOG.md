## 1.3.0 (June 22, 2021)

FEATURES:

* **New Data Source:** `sbercloud_dcs_az` ([#54](https://github.com/sbercloud-terraform/terraform-provider-sbercloud/issues/54))
* **New Data Source:** `sbercloud_dcs_maintainwindow` ([#54](https://github.com/sbercloud-terraform/terraform-provider-sbercloud/issues/54))
* **New Data Source:** `sbercloud_dcs_product` ([#54](https://github.com/sbercloud-terraform/terraform-provider-sbercloud/issues/54))
* **New Data Source:** `sbercloud_kms_data_key` ([#56](https://github.com/sbercloud-terraform/terraform-provider-sbercloud/issues/56))
* **New Data Source:** `sbercloud_kms_key` ([#56](https://github.com/sbercloud-terraform/terraform-provider-sbercloud/issues/56))
* **New Data Source:** `sbercloud_dms_az` ([#57](https://github.com/sbercloud-terraform/terraform-provider-sbercloud/issues/57))
* **New Data Source:** `sbercloud_dms_product` ([#57](https://github.com/sbercloud-terraform/terraform-provider-sbercloud/issues/57))
* **New Data Source:** `sbercloud_dms_maintainwindow` ([#57](https://github.com/sbercloud-terraform/terraform-provider-sbercloud/issues/57))
* **New Data Source:** `sbercloud_dds_flavors` ([#63](https://github.com/sbercloud-terraform/terraform-provider-sbercloud/issues/63))
* **New Resource:** `sbercloud_api_gateway_api` ([#52](https://github.com/sbercloud-terraform/terraform-provider-sbercloud/issues/52))
* **New Resource:** `sbercloud_api_gateway_group` ([#52](https://github.com/sbercloud-terraform/terraform-provider-sbercloud/issues/52))
* **New Resource:** `sbercloud_dcs_instance` ([#54](https://github.com/sbercloud-terraform/terraform-provider-sbercloud/issues/54))
* **New Resource:** `sbercloud_kms_key` ([#56](https://github.com/sbercloud-terraform/terraform-provider-sbercloud/issues/56))
* **New Resource:** `sbercloud_smn_subscription` ([#58](https://github.com/sbercloud-terraform/terraform-provider-sbercloud/issues/58))
* **New Resource:** `sbercloud_smn_topic` ([#58](https://github.com/sbercloud-terraform/terraform-provider-sbercloud/issues/58))
* **New Resource:** `sbercloud_dms_instance` ([#57](https://github.com/sbercloud-terraform/terraform-provider-sbercloud/issues/57))
* **New Resource:** `sbercloud_fgs_function` ([#60](https://github.com/sbercloud-terraform/terraform-provider-sbercloud/issues/60))
* **New Resource:** `sbercloud_dds_instance` ([#63](https://github.com/sbercloud-terraform/terraform-provider-sbercloud/issues/63))

## 1.2.0 (June 04, 2021)

FEATURES:

* **New Data Source:** `sbercloud_cce_node_pool` ([#43](https://github.com/sbercloud-terraform/terraform-provider-sbercloud/issues/43))
* **New Data Source:** `sbercloud_rds_flavors` ([#45](https://github.com/sbercloud-terraform/terraform-provider-sbercloud/issues/45))
* **New Resource:** `sbercloud_network_acl` ([#38](https://github.com/sbercloud-terraform/terraform-provider-sbercloud/issues/38))
* **New Resource:** `sbercloud_network_acl_rule` ([#38](https://github.com/sbercloud-terraform/terraform-provider-sbercloud/issues/38))
* **New Resource:** `sbercloud_identity_agency` ([#39](https://github.com/sbercloud-terraform/terraform-provider-sbercloud/issues/39))
* **New Resource:** `sbercloud_cce_node_pool` ([#43](https://github.com/sbercloud-terraform/terraform-provider-sbercloud/issues/43))
* **New Resource:** `sbercloud_rds_instance` ([#45](https://github.com/sbercloud-terraform/terraform-provider-sbercloud/issues/45))
* **New Resource:** `sbercloud_rds_parametergroup` ([#45](https://github.com/sbercloud-terraform/terraform-provider-sbercloud/issues/45))
* **New Resource:** `sbercloud_rds_read_replica_instance` ([#45](https://github.com/sbercloud-terraform/terraform-provider-sbercloud/issues/45))

## 1.1.0 (May 17, 2021)

FEATURES:

* **New Data Source:** `sbercloud_nat_gateway` ([#27](https://github.com/sbercloud-terraform/terraform-provider-sbercloud/pull/27))
* **New Data Source:** `sbercloud_sfs_file_system` ([#33](https://github.com/sbercloud-terraform/terraform-provider-sbercloud/issues/33))
* **New Resource:** `sbercloud_nat_gateway` ([#27](https://github.com/sbercloud-terraform/terraform-provider-sbercloud/pull/27))
* **New Resource:** `sbercloud_nat_dnat_rule` ([#27](https://github.com/sbercloud-terraform/terraform-provider-sbercloud/pull/27))
* **New Resource:** `sbercloud_nat_snat_rule` ([#27](https://github.com/sbercloud-terraform/terraform-provider-sbercloud/pull/27))
* **New Resource:** `sbercloud_sfs_turbo` ([#29](https://github.com/sbercloud-terraform/terraform-provider-sbercloud/pull/29))
* **New Resource:** `sbercloud_sfs_access_rule` ([#33](https://github.com/sbercloud-terraform/terraform-provider-sbercloud/issues/33))
* **New Resource:** `sbercloud_sfs_file_system` ([#33](https://github.com/sbercloud-terraform/terraform-provider-sbercloud/issues/33))
* **New Resource:** `sbercloud_networking_eip_associate` ([#34](https://github.com/sbercloud-terraform/terraform-provider-sbercloud/issues/34))

## 1.0.0 (December 29, 2020)

FEATURES:

* **New Data Source:** `sbercloud_availability_zones`
* **New Data Source:** `sbercloud_cce_cluster`
* **New Data Source:** `sbercloud_cce_node`
* **New Data Source:** `sbercloud_compute_flavors`
* **New Data Source:** `sbercloud_identity_role`
* **New Data Source:** `sbercloud_images_image`
* **New Data Source:** `sbercloud_networking_port`
* **New Data Source:** `sbercloud_networking_secgroup`
* **New Data Source:** `sbercloud_obs_bucket_object`
* **New Data Source:** `sbercloud_vpc`
* **New Data Source:** `sbercloud_vpc_bandwidth`
* **New Data Source:** `sbercloud_vpc_route`
* **New Data Source:** `sbercloud_vpc_subnet`
* **New Data Source:** `sbercloud_vpc_subnet_ids`
* **New Resource:** `sbercloud_as_configuration`
* **New Resource:** `sbercloud_as_group`
* **New Resource:** `sbercloud_as_policy`
* **New Resource:** `sbercloud_cce_cluster`
* **New Resource:** `sbercloud_cce_node`
* **New Resource:** `sbercloud_dns_recordset`
* **New Resource:** `sbercloud_dns_zone`
* **New Resource:** `sbercloud_identity_role_assignment`
* **New Resource:** `sbercloud_identity_user`
* **New Resource:** `sbercloud_identity_group`
* **New Resource:** `sbercloud_identity_group_membership`
* **New Resource:** `sbercloud_images_image`
* **New Resource:** `sbercloud_compute_instance`
* **New Resource:** `sbercloud_compute_interface_attach`
* **New Resource:** `sbercloud_compute_keypair`
* **New Resource:** `sbercloud_compute_servergroup`
* **New Resource:** `sbercloud_compute_eip_associate`
* **New Resource:** `sbercloud_compute_volume_attach`
* **New Resource:** `sbercloud_evs_snapshot`
* **New Resource:** `sbercloud_evs_volume`
* **New Resource:** `sbercloud_lb_certificate`
* **New Resource:** `sbercloud_lb_l7policy`
* **New Resource:** `sbercloud_lb_l7rule`
* **New Resource:** `sbercloud_lb_listener`
* **New Resource:** `sbercloud_lb_loadbalancer`
* **New Resource:** `sbercloud_lb_member`
* **New Resource:** `sbercloud_lb_monitor`
* **New Resource:** `sbercloud_lb_pool`
* **New Resource:** `sbercloud_lb_whitelist`
* **New Resource:** `sbercloud_obs_bucket`
* **New Resource:** `sbercloud_obs_bucket_object`
* **New Resource:** `sbercloud_obs_bucket_policy`
* **New Resource:** `sbercloud_networking_secgroup`
* **New Resource:** `sbercloud_networking_secgroup_rule`
* **New Resource:** `sbercloud_vpc`
* **New Resource:** `sbercloud_vpc_eip`
* **New Resource:** `sbercloud_vpc_subnet`
* **New Resource:** `sbercloud_vpc_route`
* **New Resource:** `sbercloud_vpc_peering_connection`
