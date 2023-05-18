## 1.10.1 (May 17, 2023)

ENHANCEMENTS:

* Upgrade to terraform-provider-huaweicloud `v1.48.0` ([#171](https://github.com/sbercloud-terraform/terraform-provider-sbercloud/issues/171))
* Upgrade to new golangsdk ([#172](https://github.com/sbercloud-terraform/terraform-provider-sbercloud/issues/172))

## 1.10.0 (October 28, 2022)

FEATURES:

* **New Resource:** `sbercloud_aom_service_discovery_rule` ([#151](https://github.com/sbercloud-terraform/terraform-provider-sbercloud/issues/151)
* **New Resource:** `sbercloud_dli_database` ([#161](https://github.com/sbercloud-terraform/terraform-provider-sbercloud/issues/161)
* **New Resource:** `sbercloud_dli_package` ([#161](https://github.com/sbercloud-terraform/terraform-provider-sbercloud/issues/161)
* **New Resource:** `sbercloud_dli_spark_job` ([#161](https://github.com/sbercloud-terraform/terraform-provider-sbercloud/issues/161)

ENHANCEMENTS:

* Update the documentation for ELB resources ([#159](https://github.com/sbercloud-terraform/terraform-provider-sbercloud/pull/159))
* Update the documentation for DLI resources ([#161](https://github.com/sbercloud-terraform/terraform-provider-sbercloud/issues/161)

## 1.9.0 (August 31, 2022)

FEATURES:

* **New Data Source:** `sbercloud_compute_instance` ([#148](https://github.com/sbercloud-terraform/terraform-provider-sbercloud/issues/148))
* **New Data Source:** `sbercloud_compute_instances` ([#148](https://github.com/sbercloud-terraform/terraform-provider-sbercloud/issues/148))
* **New Data Source:** `sbercloud_cce_clusters` ([#155](https://github.com/sbercloud-terraform/terraform-provider-sbercloud/issues/155)
* **New Data Source:** `sbercloud_cce_nodes` ([#155](https://github.com/sbercloud-terraform/terraform-provider-sbercloud/issues/155)
* **New Resource:** `sbercloud_cce_node_attach` ([#155](https://github.com/sbercloud-terraform/terraform-provider-sbercloud/issues/155)
* **New Resource:** `sbercloud_cce_namespace` ([#155](https://github.com/sbercloud-terraform/terraform-provider-sbercloud/issues/155)
* **New Resource:** `sbercloud_cce_pvc` ([#155](https://github.com/sbercloud-terraform/terraform-provider-sbercloud/issues/155)

ENHANCEMENTS:

* Update the documentation for Network ACL resources ([#141](https://github.com/sbercloud-terraform/terraform-provider-sbercloud/issues/141))
* Update the documentation for CCE objects ([#155](https://github.com/sbercloud-terraform/terraform-provider-sbercloud/pull/155))

## 1.8.1 (June 24, 2022)

ENHANCEMENTS:

* Update documentation for DCS objects ([#146](https://github.com/sbercloud-terraform/terraform-provider-sbercloud/pull/146))

## 1.8.0 (June 6, 2022)

FEATURES:

* **New Data Source:** `sbercloud_cce_addon_template` ([#43](https://github.com/sbercloud-terraform/terraform-provider-sbercloud/issues/43))
* **New Resource:** `sbercloud_cce_addon` ([#43](https://github.com/sbercloud-terraform/terraform-provider-sbercloud/issues/43))
* **New Resource:** `sbercloud_mapreduce_cluster` ([#67](https://github.com/sbercloud-terraform/terraform-provider-sbercloud/issues/67))
* **New Resource:** `sbercloud_mapreduce_job` ([#67](https://github.com/sbercloud-terraform/terraform-provider-sbercloud/issues/67))

ENHANCEMENTS:

* Update documentation for CCE objects ([#145](https://github.com/sbercloud-terraform/terraform-provider-sbercloud/pull/145))

## 1.7.0 (May 10, 2022)

FEATURES:

* **New Resource:** `sbercloud_lts_group` ([#118](https://github.com/sbercloud-terraform/terraform-provider-sbercloud/issues/118))
* **New Resource:** `sbercloud_lts_stream` ([#118](https://github.com/sbercloud-terraform/terraform-provider-sbercloud/issues/118))

## 1.6.3 (April 4, 2022)

BUG FIXES:

* Fix broken creation of ECS resource with SAS disk ([#137](https://github.com/sbercloud-terraform/terraform-provider-sbercloud/issues/137))

## 1.6.2 (April 1, 2022)

ENHANCEMENTS:

* Add `security_token` description to provider documentation ([#133](https://github.com/sbercloud-terraform/terraform-provider-sbercloud/pull/133))

BUG FIXES:

* Fix broken creation of ECS and security group resources ([#132](https://github.com/sbercloud-terraform/terraform-provider-sbercloud/issues/132))

## 1.6.1 (March 30, 2022)

ENHANCEMENTS:

* Add support for `security_token` parameter to authenticate with a temporary security credentials ([#126](https://github.com/sbercloud-terraform/terraform-provider-sbercloud/pull/126))
* Upgrade to terraform-provider-huaweicloud `v1.34.1` ([#128](https://github.com/sbercloud-terraform/terraform-provider-sbercloud/pull/128))
* Update ECS doc examples with required params ([#129](https://github.com/sbercloud-terraform/terraform-provider-sbercloud/pull/129))
* Remove unsupported disk type from doc examples ([#130](https://github.com/sbercloud-terraform/terraform-provider-sbercloud/pull/130))

BUG FIXES:

* Fix an issue when ECS ipv4 gets imported as fixed_ip_v6 ([#125](https://github.com/sbercloud-terraform/terraform-provider-sbercloud/issues/125))

## 1.6.0 (February 25, 2022)

FEATURES:

* **New Data Source:** `sbercloud_cbr_vaults` ([#117](https://github.com/sbercloud-terraform/terraform-provider-sbercloud/issues/117))
* **New Data Source:** `sbercloud_enterprise_project` ([#119](https://github.com/sbercloud-terraform/terraform-provider-sbercloud/issues/119))
* **New Resource:** `sbercloud_cbr_policy` ([#117](https://github.com/sbercloud-terraform/terraform-provider-sbercloud/issues/117))
* **New Resource:** `sbercloud_cbr_vault` ([#117](https://github.com/sbercloud-terraform/terraform-provider-sbercloud/issues/117))
* **New Resource:** `sbercloud_enterprise_project` ([#119](https://github.com/sbercloud-terraform/terraform-provider-sbercloud/issues/119))
* **New Resource:** `sbercloud_identity_project` ([#123](https://github.com/sbercloud-terraform/terraform-provider-sbercloud/pull/123))

BUG FIXES:

* Fix the resource schema version for `sbercloud_vpc_route` resource ([#122](https://github.com/sbercloud-terraform/terraform-provider-sbercloud/pull/122))

## 1.5.1 (February 8, 2022)

ENHANCEMENTS:

* Add more useful examples ([#111](https://github.com/sbercloud-terraform/terraform-provider-sbercloud/pull/111))
* Upgrade to terraform-provider-huaweicloud `v1.32.2` ([#113](https://github.com/sbercloud-terraform/terraform-provider-sbercloud/pull/113))

## 1.5.0 (December 29, 2021)

FEATURES:

* **New Data Source:** `sbercloud_vpcs` ([#103](https://github.com/sbercloud-terraform/terraform-provider-sbercloud/issues/103))
* **New Data Source:** `sbercloud_vpc_eip` ([#103](https://github.com/sbercloud-terraform/terraform-provider-sbercloud/issues/103))
* **New Data Source:** `sbercloud_vpc_ids` ([#103](https://github.com/sbercloud-terraform/terraform-provider-sbercloud/issues/103))
* **New Data Source:** `sbercloud_vpc_peering_connection` ([#103](https://github.com/sbercloud-terraform/terraform-provider-sbercloud/issues/103))
* **New Data Source:** `sbercloud_vpc_route_table` ([#103](https://github.com/sbercloud-terraform/terraform-provider-sbercloud/issues/103))
* **New Data Source:** `sbercloud_vpc_subnets` ([#103](https://github.com/sbercloud-terraform/terraform-provider-sbercloud/issues/103))
* **New Data Source:** `sbercloud_identity_custom_role` ([#105](https://github.com/sbercloud-terraform/terraform-provider-sbercloud/issues/105))
* **New Data Source:** `sbercloud_identity_group` ([#105](https://github.com/sbercloud-terraform/terraform-provider-sbercloud/issues/105))
* **New Resource:** `sbercloud_vpc_peering_connection_accepter` ([#103](https://github.com/sbercloud-terraform/terraform-provider-sbercloud/issues/103))
* **New Resource:** `sbercloud_vpc_route_table` ([#103](https://github.com/sbercloud-terraform/terraform-provider-sbercloud/issues/103))
* **New Resource:** `sbercloud_identity_access_key` ([#105](https://github.com/sbercloud-terraform/terraform-provider-sbercloud/issues/105))
* **New Resource:** `sbercloud_identity_acl` ([#105](https://github.com/sbercloud-terraform/terraform-provider-sbercloud/issues/105))
* **New Resource:** `sbercloud_identity_role` ([#105](https://github.com/sbercloud-terraform/terraform-provider-sbercloud/issues/105))
* **New Resource:** `sbercloud_dms_kafka_instance` ([#107](https://github.com/sbercloud-terraform/terraform-provider-sbercloud/issues/107))
* **New Resource:** `sbercloud_dms_kafka_topic` ([#107](https://github.com/sbercloud-terraform/terraform-provider-sbercloud/issues/107))
* **New Resource:** `sbercloud_dms_rabbitmq_instance` ([#107](https://github.com/sbercloud-terraform/terraform-provider-sbercloud/issues/107))

ENHANCEMENTS:

* Upgrade to terraform-plugin-sdk v2 ([#99](https://github.com/sbercloud-terraform/terraform-provider-sbercloud/issues/99))
* Upgrade to new golangsdk ([#101](https://github.com/sbercloud-terraform/terraform-provider-sbercloud/issues/101))
* Upgrade to terraform-provider-huaweicloud `v1.31.0` ([#104](https://github.com/sbercloud-terraform/terraform-provider-sbercloud/pull/104))

DEPRECATE:

* data/sbercloud_dcs_az
* data/sbercloud_dcs_product

## 1.4.0 (August 02, 2021)

FEATURES:

* **New Data Source:** `sbercloud_dis_partition` ([#68](https://github.com/sbercloud-terraform/terraform-provider-sbercloud/issues/68))
* **New Data Source:** `sbercloud_cdm_flavors` ([#69](https://github.com/sbercloud-terraform/terraform-provider-sbercloud/issues/69))
* **New Resource:** `sbercloud_dis_stream` ([#68](https://github.com/sbercloud-terraform/terraform-provider-sbercloud/issues/68))
* **New Resource:** `sbercloud_cdm_cluster` ([#69](https://github.com/sbercloud-terraform/terraform-provider-sbercloud/issues/69))
* **New Resource:** `sbercloud_dli_queue` ([#72](https://github.com/sbercloud-terraform/terraform-provider-sbercloud/issues/72))
* **New Resource:** `sbercloud_dws_cluster` ([#73](https://github.com/sbercloud-terraform/terraform-provider-sbercloud/issues/73))
* **New Resource:** `sbercloud_ges_graph` ([#75](https://github.com/sbercloud-terraform/terraform-provider-sbercloud/issues/75))
* **New Resource:** `sbercloud_ces_alarmrule` ([#77](https://github.com/sbercloud-terraform/terraform-provider-sbercloud/issues/77))
* **New Resource:** `sbercloud_css_cluster` ([#78](https://github.com/sbercloud-terraform/terraform-provider-sbercloud/issues/78))

ENHANCEMENTS:

* Update GNUmakefile to make log message configurable ([#80](https://github.com/sbercloud-terraform/terraform-provider-sbercloud/pull/80))
* Update go to `1.16` in setup-go action to enable darwin/arm64 builds ([#81](https://github.com/sbercloud-terraform/terraform-provider-sbercloud/pull/81))
* Add support of scale up and class changing for the `sbercloud_dds_instance` resource ([#91](https://github.com/sbercloud-terraform/terraform-provider-sbercloud/pull/91))

BUG FIXES:

* Fix an issue when the `sbercloud_rds_instance` resource cannot be created with empty database port value ([#90](https://github.com/sbercloud-terraform/terraform-provider-sbercloud/pull/90))

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
