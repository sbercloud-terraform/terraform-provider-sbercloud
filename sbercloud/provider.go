package sbercloud

import (
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/apig"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/cts"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/dew"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/dns"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/ecs"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/kafka"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/lts"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/nat"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/obs"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/rabbitmq"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/sfsturbo"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/swr"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/vpn"
	"github.com/sbercloud-terraform/terraform-provider-sbercloud/sbercloud/services/cbh"
	deprecated_sbc "github.com/sbercloud-terraform/terraform-provider-sbercloud/sbercloud/services/deprecated"
	ges_sbercloud "github.com/sbercloud-terraform/terraform-provider-sbercloud/sbercloud/services/ges"
	"github.com/sbercloud-terraform/terraform-provider-sbercloud/sbercloud/services/rds"
	vpc2 "github.com/sbercloud-terraform/terraform-provider-sbercloud/sbercloud/services/vpc"
	"github.com/sbercloud-terraform/terraform-provider-sbercloud/sbercloud/services/vpcep"

	elb2 "github.com/sbercloud-terraform/terraform-provider-sbercloud/sbercloud/services/elb"
	"sync"

	dcs2 "github.com/sbercloud-terraform/terraform-provider-sbercloud/sbercloud/services/dcs"
	dds2 "github.com/sbercloud-terraform/terraform-provider-sbercloud/sbercloud/services/dds"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/mutexkv"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/aom"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/as"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/cbr"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/cce"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/cdm"
	css_huawei "github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/css"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/dcs"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/deprecated"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/dis"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/dli"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/dms"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/dws"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/eip"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/elb"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/eps"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/evs"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/fgs"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/iam"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/ims"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/lb"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/mrs"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/smn"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/vpc"
	"github.com/sbercloud-terraform/terraform-provider-sbercloud/sbercloud/services/ces"
	"github.com/sbercloud-terraform/terraform-provider-sbercloud/sbercloud/services/css"
	dli_sbercloud "github.com/sbercloud-terraform/terraform-provider-sbercloud/sbercloud/services/dli"
	"github.com/sbercloud-terraform/terraform-provider-sbercloud/sbercloud/services/drs"
)

// This is a global MutexKV for use within this plugin.
var osMutexKV = mutexkv.NewMutexKV()

// Provider returns a schema.Provider for SberCloud.
func Provider() *schema.Provider {
	provider := &schema.Provider{
		Schema: map[string]*schema.Schema{
			"access_key": {
				Type:         schema.TypeString,
				Optional:     true,
				DefaultFunc:  schema.EnvDefaultFunc("SBC_ACCESS_KEY", nil),
				Description:  descriptions["access_key"],
				RequiredWith: []string{"secret_key"},
			},

			"secret_key": {
				Type:         schema.TypeString,
				Optional:     true,
				DefaultFunc:  schema.EnvDefaultFunc("SBC_SECRET_KEY", nil),
				Description:  descriptions["secret_key"],
				RequiredWith: []string{"access_key"},
			},

			"security_token": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  descriptions["security_token"],
				RequiredWith: []string{"access_key"},
				DefaultFunc:  schema.EnvDefaultFunc("SBC_SECURITY_TOKEN", nil),
			},

			"auth_url": {
				Type:     schema.TypeString,
				Optional: true,
				DefaultFunc: schema.EnvDefaultFunc(
					"SBC_AUTH_URL", "https://iam.ru-moscow-1.hc.sbercloud.ru/v3"),
				Description: descriptions["auth_url"],
			},

			"region": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  descriptions["region"],
				DefaultFunc:  schema.EnvDefaultFunc("SBC_REGION_NAME", nil),
				InputDefault: "ru-moscow-1",
			},

			"user_name": {
				Type:         schema.TypeString,
				Optional:     true,
				DefaultFunc:  schema.EnvDefaultFunc("SBC_USERNAME", ""),
				Description:  descriptions["user_name"],
				RequiredWith: []string{"password", "account_name"},
			},

			"project_name": {
				Type:     schema.TypeString,
				Optional: true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"SBC_PROJECT_NAME",
				}, ""),
				Description: descriptions["project_name"],
			},

			"password": {
				Type:         schema.TypeString,
				Optional:     true,
				Sensitive:    true,
				DefaultFunc:  schema.EnvDefaultFunc("SBC_PASSWORD", ""),
				Description:  descriptions["password"],
				RequiredWith: []string{"user_name", "account_name"},
			},

			"assume_role": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"agency_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: descriptions["assume_role_agency_name"],
							DefaultFunc: schema.EnvDefaultFunc("SBC_ASSUME_ROLE_AGENCY_NAME", nil),
						},
						"domain_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: descriptions["assume_role_domain_name"],
							DefaultFunc: schema.EnvDefaultFunc("SBC_ASSUME_ROLE_DOMAIN_NAME", nil),
						},
					},
				},
			},

			"account_name": {
				Type:     schema.TypeString,
				Optional: true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"SBC_ACCOUNT_NAME",
				}, ""),
				Description:  descriptions["account_name"],
				RequiredWith: []string{"password", "user_name"},
			},

			"insecure": {
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("SBC_INSECURE", false),
				Description: descriptions["insecure"],
			},

			"enterprise_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: descriptions["enterprise_project_id"],
				DefaultFunc: schema.EnvDefaultFunc("SBC_ENTERPRISE_PROJECT_ID", ""),
			},

			"max_retries": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: descriptions["max_retries"],
				DefaultFunc: schema.EnvDefaultFunc("SBC_MAX_RETRIES", 5),
			},
			"domain_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: descriptions["domain_id"],
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"SBC_DOMAIN_ID",
					"OS_DOMAIN_ID",
					"OS_USER_DOMAIN_ID",
					"OS_PROJECT_DOMAIN_ID",
				}, ""),
			},
			"domain_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: descriptions["domain_name"],
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"SBC_DOMAIN_NAME",
					"OS_DOMAIN_NAME",
					"OS_USER_DOMAIN_NAME",
					"OS_PROJECT_DOMAIN_NAME",
				}, ""),
			},
		},

		DataSourcesMap: map[string]*schema.Resource{
			"sbercloud_availability_zones": huaweicloud.DataSourceAvailabilityZones(),

			"sbercloud_apig_acl_policies":                       apig.DataSourceAclPolicies(),
			"sbercloud_apig_api_associated_acl_policies":        apig.DataSourceApiAssociatedAclPolicies(),
			"sbercloud_apig_api_associated_applications":        apig.DataSourceApiAssociatedApplications(),
			"sbercloud_apig_api_associated_plugins":             apig.DataSourceApiAssociatedPlugins(),
			"sbercloud_apig_api_associated_signatures":          apig.DataSourceApiAssociatedSignatures(),
			"sbercloud_apig_api_associated_throttling_policies": apig.DataSourceApiAssociatedThrottlingPolicies(),
			"sbercloud_apig_api_basic_configurations":           apig.DataSourceApiBasicConfigurations(),
			"sbercloud_apig_api":                                apig.DataSourceApi(),
			"sbercloud_apig_appcodes":                           apig.DataSourceAppcodes(),
			"sbercloud_apig_applications":                       apig.DataSourceApplications(),
			"sbercloud_apig_application_acl":                    apig.DataSourceApplicationAcl(),
			"sbercloud_apig_application_quotas":                 apig.DataSourceApigApplicationQuotas(),
			"sbercloud_apig_channels":                           apig.DataSourceChannels(),
			"sbercloud_apig_custom_authorizers":                 apig.DataSourceCustomAuthorizers(),
			"sbercloud_apig_endpoint_connections":               apig.DataSourceApigEndpointConnections(),
			"sbercloud_apig_environment_variables":              apig.DataSourceApigEnvironmentVariables(),
			"sbercloud_apig_environments":                       apig.DataSourceEnvironments(),
			"sbercloud_apig_groups":                             apig.DataSourceGroups(),
			"sbercloud_apig_instance_features":                  apig.DataSourceInstanceFeatures(),
			"sbercloud_apig_instance_supported_features":        apig.DataSourceInstanceSupportedFeatures(),
			"sbercloud_apig_instances":                          apig.DataSourceInstances(),
			"sbercloud_apig_signatures":                         apig.DataSourceSignatures(),
			"sbercloud_apig_throttling_policies":                apig.DataSourceThrottlingPolicies(),

			"sbercloud_css_flavors": css_huawei.DataSourceCssFlavors(),

			"sbercloud_cbh_instances":          cbh.DataSourceCbhInstances(),
			"sbercloud_cbh_flavors":            cbh.DataSourceCbhFlavors(),
			"sbercloud_cbh_availability_zones": cbh.DataSourceAvailabilityZones(),

			"sbercloud_cbr_backup":              cbr.DataSourceBackup(),
			"sbercloud_cbr_vaults":              cbr.DataSourceVaults(),
			"sbercloud_cbr_policies":            cbr.DataSourcePolicies(),
			"sbercloud_cce_addon_template":      cce.DataSourceAddonTemplate(),
			"sbercloud_cce_cluster":             cce.DataSourceCCEClusterV3(),
			"sbercloud_cce_clusters":            cce.DataSourceCCEClusters(),
			"sbercloud_cce_node":                cce.DataSourceNode(),
			"sbercloud_cce_nodes":               cce.DataSourceNodes(),
			"sbercloud_cce_node_pool":           cce.DataSourceCCENodePoolV3(),
			"sbercloud_cce_cluster_certificate": cce.DataSourceCCEClusterCertificate(),
			"sbercloud_cdm_flavors":             cdm.DataSourceCdmFlavors(),
			"sbercloud_compute_flavors":         ecs.DataSourceEcsFlavors(),
			"sbercloud_compute_instance":        ecs.DataSourceComputeInstance(),
			"sbercloud_compute_instances":       ecs.DataSourceComputeInstances(),
			"sbercloud_compute_servergroups":    ecs.DataSourceComputeServerGroups(),
			"sbercloud_dcs_flavors":             dcs.DataSourceDcsFlavorsV2(),
			"sbercloud_dcs_accounts":            dcs.DataSourceDcsAccounts(),
			"sbercloud_dcs_az":                  deprecated.DataSourceDcsAZV1(),
			"sbercloud_dcs_maintainwindow":      dcs.DataSourceDcsMaintainWindow(),
			"sbercloud_dcs_product":             deprecated.DataSourceDcsProductV1(),
			"sbercloud_dds_flavors":             dds2.DataSourceDDSFlavorV3(),
			"sbercloud_dms_az":                  deprecated.DataSourceDmsAZ(),

			"sbercloud_kps_failed_tasks":  dew.DataSourceDewKpsFailedTasks(),
			"sbercloud_kps_running_tasks": dew.DataSourceDewKpsRunningTasks(),
			"sbercloud_kps_keypairs":      dew.DataSourceKeypairs(),

			"sbercloud_dms_product":               dms.DataSourceDmsProduct(),
			"sbercloud_dms_maintainwindow":        dms.DataSourceDmsMaintainWindow(),
			"sbercloud_dms_kafka_instances":       kafka.DataSourceDmsKafkaInstances(),
			"sbercloud_dms_kafka_flavors":         kafka.DataSourceKafkaFlavors(),
			"sbercloud_dms_kafka_users":           kafka.DataSourceDmsKafkaUsers(),
			"sbercloud_dms_kafka_messages":        kafka.DataSourceDmsKafkaMessages(),
			"sbercloud_dms_kafka_consumer_groups": kafka.DataSourceDmsKafkaConsumerGroups(),

			"sbercloud_dms_rabbitmq_flavors": rabbitmq.DataSourceRabbitMQFlavors(),

			"sbercloud_dws_flavors":          dws.DataSourceDwsFlavors(),
			"sbercloud_elb_certificate":      elb.DataSourceELBCertificateV3(),
			"sbercloud_elb_flavors":          elb.DataSourceElbFlavorsV3(),
			"sbercloud_elb_pools":            elb.DataSourcePools(),
			"sbercloud_enterprise_project":   eps.DataSourceEnterpriseProject(),
			"sbercloud_evs_volumes":          evs.DataSourceEvsVolumesV2(),
			"sbercloud_identity_role":        iam.DataSourceIdentityRole(),
			"sbercloud_identity_custom_role": iam.DataSourceIdentityCustomRole(),
			"sbercloud_identity_group":       iam.DataSourceIdentityGroup(),
			"sbercloud_identity_projects":    iam.DataSourceIdentityProjects(),
			"sbercloud_identity_users":       iam.DataSourceIdentityUsers(),
			"sbercloud_images_image":         ims.DataSourceImagesImageV2(),
			"sbercloud_images_images":        ims.DataSourceImagesImages(),
			"sbercloud_kms_key":              dew.DataSourceKmsKey(),
			"sbercloud_kms_data_key":         dew.DataSourceKmsDataKeyV1(),
			"sbercloud_lb_listeners":         lb.DataSourceListeners(),
			"sbercloud_lb_loadbalancer":      lb.DataSourceELBV2Loadbalancer(),
			"sbercloud_lb_certificate":       lb.DataSourceLBCertificateV2(),
			"sbercloud_lb_pools":             lb.DataSourcePools(),
			"sbercloud_nat_gateway":          nat.DataSourcePublicGateway(),
			"sbercloud_networking_port":      vpc.DataSourceNetworkingPortV2(),
			"sbercloud_networking_secgroup":  vpc.DataSourceNetworkingSecGroup(),
			"sbercloud_networking_secgroups": vpc.DataSourceNetworkingSecGroups(),
			"sbercloud_obs_buckets":          obs.DataSourceObsBuckets(),
			"sbercloud_obs_bucket_object":    obs.DataSourceObsBucketObject(),

			"sbercloud_rds_pg_plugins":                      rds.DataSourcePgPlugins(),
			"sbercloud_rds_pg_accounts":                     rds.DataSourcePgAccounts(),
			"sbercloud_rds_pg_roles":                        rds.DataSourceRdsPgRoles(),
			"sbercloud_rds_pg_databases":                    rds.DataSourcePgDatabases(),
			"sbercloud_rds_pg_sql_limits":                   rds.DataSourceRdsPgSqlLimits(),
			"sbercloud_rds_pg_plugin_parameter_value_range": rds.DataSourceRdsPgPluginParameterValueRange(),
			"sbercloud_rds_pg_plugin_parameter_values":      rds.DataSourceRdsPgPluginParameterValues(),

			"sbercloud_rds_flavors":         rds.DataSourceRdsFlavor(),
			"sbercloud_rds_backups":         rds.DataSourceRdsBackups(),
			"sbercloud_rds_engine_versions": rds.DataSourceRdsEngineVersionsV3(),
			"sbercloud_rds_instances":       rds.DataSourceRdsInstances(),
			"sbercloud_rds_storage_types":   rds.DataSourceStoragetype(),
			//"sbercloud_sfs_file_system":                sfs.DataSourceSFSFileSystemV2(),
			//"sbercloud_sfs_turbos":                     sfs.DataSourceTurbos(),
			"sbercloud_sfs_turbos":            sfsturbo.DataSourceTurbos(),
			"sbercloud_sfs_turbo_data_tasks":  sfsturbo.DataSourceSfsTurboDataTasks(),
			"sbercloud_sfs_turbo_du_tasks":    sfsturbo.DataSourceSfsTurboDuTasks(),
			"sbercloud_sfs_turbo_obs_targets": sfsturbo.DataSourceSfsTurboObsTargets(),
			"sbercloud_sfs_turbo_perm_rules":  sfsturbo.DataSourceSfsTurboPermRules(),
			"sbercloud_sfs_file_system":       deprecated.DataSourceSFSFileSystemV2(),
			"sbercloud_sfs_file_system_v2":    deprecated.DataSourceSFSFileSystemV2(),

			"sbercloud_vpc":                            vpc.DataSourceVpcV1(),
			"sbercloud_vpcs":                           vpc.DataSourceVpcs(),
			"sbercloud_vpc_address_groups":             vpc.DataSourceVpcAddressGroups(),
			"sbercloud_vpc_bandwidth":                  eip.DataSourceBandWidth(),
			"sbercloud_vpc_eip":                        eip.DataSourceVpcEip(),
			"sbercloud_vpc_eips":                       eip.DataSourceVpcEips(),
			"sbercloud_vpc_ids":                        vpc.DataSourceVpcIdsV1(),
			"sbercloud_vpc_peering_connection":         vpc.DataSourceVpcPeeringConnectionV2(),
			"sbercloud_vpc_routes":                     vpc.DataSourceVpcRoutes(),
			"sbercloud_vpc_route":                      vpc.DataSourceVpcRouteV2(),
			"sbercloud_vpc_route_table":                vpc.DataSourceVPCRouteTable(),
			"sbercloud_vpc_subnet":                     vpc.DataSourceVpcSubnetV1(),
			"sbercloud_vpc_subnets":                    vpc.DataSourceVpcSubnets(),
			"sbercloud_vpc_subnet_ids":                 vpc.DataSourceVpcSubnetIdsV1(),
			"sbercloud_vpcep_public_services":          vpcep.DataSourceVPCEPPublicServices(),
			"sbercloud_vpn_gateway_availability_zones": vpn.DataSourceVpnGatewayAZs(),
			"sbercloud_vpn_gateways":                   vpn.DataSourceGateways(),
			"sbercloud_vpn_customer_gateways":          vpn.DataSourceVpnCustomerGateways(),
			"sbercloud_vpn_connections":                vpn.DataSourceVpnConnections(),
			"sbercloud_vpn_connection_health_checks":   vpn.DataSourceVpnConnectionHealthChecks(),

			// Legacy
			"sbercloud_identity_role_v3": iam.DataSourceIdentityRole(),
		},

		ResourcesMap: map[string]*schema.Resource{
			"sbercloud_aom_service_discovery_rule": aom.ResourceServiceDiscoveryRule(),
			"sbercloud_api_gateway_api":            apig.ResourceApigAPIV2(),
			"sbercloud_api_gateway_group":          apig.ResourceApigGroupV2(),

			"sbercloud_apig_acl_policy":                     apig.ResourceAclPolicy(),
			"sbercloud_apig_acl_policy_associate":           apig.ResourceAclPolicyAssociate(),
			"sbercloud_apig_api":                            apig.ResourceApigAPIV2(),
			"sbercloud_apig_api_publishment":                apig.ResourceApigApiPublishment(),
			"sbercloud_apig_appcode":                        apig.ResourceAppcode(),
			"sbercloud_apig_application":                    apig.ResourceApigApplicationV2(),
			"sbercloud_apig_application_acl":                apig.ResourceApplicationAcl(),
			"sbercloud_apig_application_authorization":      apig.ResourceAppAuth(),
			"sbercloud_apig_application_quota":              apig.ResourceApplicationQuota(),
			"sbercloud_apig_application_quota_associate":    apig.ResourceApplicationQuotaAssociate(),
			"sbercloud_apig_certificate":                    apig.ResourceCertificate(),
			"sbercloud_apig_channel":                        apig.ResourceChannel(),
			"sbercloud_apig_custom_authorizer":              apig.ResourceApigCustomAuthorizerV2(),
			"sbercloud_apig_endpoint_connection_management": apig.ResourceEndpointConnectionManagement(),
			"sbercloud_apig_environment":                    apig.ResourceApigEnvironmentV2(),
			"sbercloud_apig_environment_variable":           apig.ResourceEnvironmentVariable(),
			"sbercloud_apig_group":                          apig.ResourceApigGroupV2(),
			"sbercloud_apig_instance_feature":               apig.ResourceInstanceFeature(),
			"sbercloud_apig_instance_routes":                apig.ResourceInstanceRoutes(),
			"sbercloud_apig_instance":                       apig.ResourceApigInstanceV2(),
			"sbercloud_apig_plugin_associate":               apig.ResourcePluginAssociate(),
			"sbercloud_apig_plugin":                         apig.ResourcePlugin(),
			"sbercloud_apig_response":                       apig.ResourceApigResponseV2(),
			"sbercloud_apig_signature_associate":            apig.ResourceSignatureAssociate(),
			"sbercloud_apig_signature":                      apig.ResourceSignature(),
			"sbercloud_apig_throttling_policy_associate":    apig.ResourceThrottlingPolicyAssociate(),
			"sbercloud_apig_throttling_policy":              apig.ResourceApigThrottlingPolicyV2(),
			"sbercloud_apig_endpoint_whitelist":             apig.ResourceEndpointWhiteList(),

			"sbercloud_as_configuration":    as.ResourceASConfiguration(),
			"sbercloud_as_group":            as.ResourceASGroup(),
			"sbercloud_as_policy":           as.ResourceASPolicy(),
			"sbercloud_as_bandwidth_policy": as.ResourceASBandWidthPolicy(),

			"sbercloud_cbr_backup_share_accepter": cbr.ResourceBackupShareAccepter(),
			"sbercloud_cbr_backup_share":          cbr.ResourceBackupShare(),
			"sbercloud_cbr_checkpoint":            cbr.ResourceCheckpoint(),

			"sbercloud_cbh_instance":                   cbh.ResourceCBHInstance(),
			"sbercloud_cbh_ha_instance":                cbh.ResourceCBHHAInstance(),
			"sbercloud_cbh_asset_agency_authorization": cbh.ResourceAssetAgencyAuthorization(),

			"sbercloud_cbr_policy": cbr.ResourcePolicy(),
			"sbercloud_cbr_vault":  cbr.ResourceVault(),

			"sbercloud_css_cluster":       css.ResourceCssCluster(),
			"sbercloud_css_configuration": css_huawei.ResourceCssConfiguration(),

			"sbercloud_cce_addon":       cce.ResourceAddon(),
			"sbercloud_cce_cluster":     cce.ResourceCluster(),
			"sbercloud_cce_namespace":   cce.ResourceCCENamespaceV1(),
			"sbercloud_cce_node":        cce.ResourceNode(),
			"sbercloud_cce_node_attach": cce.ResourceNodeAttach(),
			"sbercloud_cce_node_pool":   cce.ResourceNodePool(),
			"sbercloud_cce_pvc":         cce.ResourceCcePersistentVolumeClaimsV1(),

			"sbercloud_cdm_cluster": cdm.ResourceCdmCluster(),

			"sbercloud_compute_instance":         ecs.ResourceComputeInstance(),
			"sbercloud_compute_interface_attach": ecs.ResourceComputeInterfaceAttach(),
			"sbercloud_compute_servergroup":      ecs.ResourceComputeServerGroup(),
			"sbercloud_compute_eip_associate":    ecs.ResourceComputeEIPAssociate(),
			"sbercloud_compute_volume_attach":    ecs.ResourceComputeVolumeAttach(),

			"sbercloud_compute_keypair": huaweicloud.ResourceComputeKeypairV2(),

			"sbercloud_ces_alarmrule": ces.ResourceAlarmRule(),

			"sbercloud_cts_tracker":      cts.ResourceCTSTracker(),
			"sbercloud_cts_data_tracker": cts.ResourceCTSDataTracker(),
			"sbercloud_cts_notification": cts.ResourceCTSNotification(),

			"sbercloud_dcs_instance":   dcs2.ResourceDcsInstance(),
			"sbercloud_dcs_backup":     dcs.ResourceDcsBackup(),
			"sbercloud_dcs_restore":    dcs2.ResourceDcsRestore(),
			"sbercloud_dcs_parameters": deprecated_sbc.ResourceDcsParameters(),
			"sbercloud_dcs_account":    dcs.ResourceDcsAccount(),
			"sbercloud_dds_instance":   dds2.ResourceDdsInstanceV3(),

			"sbercloud_dis_stream": dis.ResourceDisStream(),

			"sbercloud_dli_database":  dli.ResourceDliSqlDatabaseV1(),
			"sbercloud_dli_package":   dli.ResourceDliPackageV2(),
			"sbercloud_dli_queue":     dli.ResourceDliQueue(),
			"sbercloud_dli_spark_job": dli_sbercloud.ResourceDliSparkJobV2(),

			"sbercloud_dms_instance":              deprecated.ResourceDmsInstancesV1(),
			"sbercloud_dms_kafka_instance":        kafka.ResourceDmsKafkaInstance(),
			"sbercloud_dms_kafka_topic":           kafka.ResourceDmsKafkaTopic(),
			"sbercloud_dms_kafka_permissions":     kafka.ResourceDmsKafkaPermissions(),
			"sbercloud_dms_kafka_user":            kafka.ResourceDmsKafkaUser(),
			"sbercloud_dms_kafka_message_produce": kafka.ResourceDmsKafkaMessageProduce(),
			"sbercloud_dms_kafka_consumer_group":  kafka.ResourceDmsKafkaConsumerGroup(),

			"sbercloud_dms_rabbitmq_instance": rabbitmq.ResourceDmsRabbitmqInstance(),

			"sbercloud_dns_recordset": dns.ResourceDNSRecordSetV2(),
			"sbercloud_dns_zone":      dns.ResourceDNSZone(),

			"sbercloud_drs_job": drs.ResourceDrsJob(),

			"sbercloud_dws_cluster": dws.ResourceDwsCluster(),

			"sbercloud_elb_certificate":     elb.ResourceCertificateV3(),
			"sbercloud_elb_l7policy":        elb.ResourceL7PolicyV3(),
			"sbercloud_elb_l7rule":          elb.ResourceL7RuleV3(),
			"sbercloud_elb_listener":        elb.ResourceListenerV3(),
			"sbercloud_elb_loadbalancer":    elb.ResourceLoadBalancerV3(),
			"sbercloud_elb_monitor":         elb2.ResourceMonitorV3(),
			"sbercloud_elb_ipgroup":         elb.ResourceIpGroupV3(),
			"sbercloud_elb_pool":            elb.ResourcePoolV3(),
			"sbercloud_elb_member":          elb.ResourceMemberV3(),
			"sbercloud_elb_security_policy": elb.ResourceSecurityPolicy(),

			"sbercloud_enterprise_project": eps.ResourceEnterpriseProject(),

			"sbercloud_evs_snapshot": evs.ResourceEvsSnapshotV2(),
			"sbercloud_evs_volume":   evs.ResourceEvsVolume(),

			"sbercloud_fgs_function": fgs.ResourceFgsFunction(),

			"sbercloud_ges_graph": ges_sbercloud.ResourceGesGraph(),

			"sbercloud_identity_access_key":            iam.ResourceIdentityKey(),
			"sbercloud_identity_acl":                   iam.ResourceIdentityACL(),
			"sbercloud_identity_agency":                iam.ResourceIAMAgencyV3(),
			"sbercloud_identity_group":                 iam.ResourceIdentityGroup(),
			"sbercloud_identity_group_membership":      iam.ResourceIdentityGroupMembership(),
			"sbercloud_identity_group_role_assignment": iam.ResourceIdentityGroupRoleAssignment(),
			"sbercloud_identity_project":               iam.ResourceIdentityProject(),
			"sbercloud_identity_provider":              iam.ResourceIdentityProvider(),
			"sbercloud_identity_provider_conversion":   iam.ResourceIAMProviderConversion(),
			"sbercloud_identity_role":                  iam.ResourceIdentityRole(),
			"sbercloud_identity_role_assignment":       iam.ResourceIdentityGroupRoleAssignment(),
			"sbercloud_identity_user":                  iam.ResourceIdentityUser(),

			"sbercloud_images_image": deprecated.ResourceImsImage(),

			"sbercloud_kms_key":     dew.ResourceKmsKey(),
			"sbercloud_kps_keypair": dew.ResourceKeypair(),

			"sbercloud_lb_certificate":  lb.ResourceCertificateV2(),
			"sbercloud_lb_l7policy":     lb.ResourceL7PolicyV2(),
			"sbercloud_lb_l7rule":       lb.ResourceL7RuleV2(),
			"sbercloud_lb_listener":     lb.ResourceListener(),
			"sbercloud_lb_loadbalancer": lb.ResourceLoadBalancer(),
			"sbercloud_lb_member":       lb.ResourceMemberV2(),
			"sbercloud_lb_monitor":      lb.ResourceMonitorV2(),
			"sbercloud_lb_pool":         lb.ResourcePoolV2(),
			"sbercloud_lb_whitelist":    lb.ResourceWhitelistV2(),

			"sbercloud_lts_group":  lts.ResourceLTSGroup(),
			"sbercloud_lts_stream": lts.ResourceLTSStream(),

			"sbercloud_mapreduce_cluster": mrs.ResourceMRSClusterV2(),
			"sbercloud_mapreduce_job":     mrs.ResourceMRSJobV2(),

			"sbercloud_nat_dnat_rule": nat.ResourcePublicDnatRule(),
			"sbercloud_nat_gateway":   nat.ResourcePublicGateway(),
			"sbercloud_nat_snat_rule": nat.ResourcePublicSnatRule(),

			"sbercloud_network_acl":      deprecated.ResourceNetworkACL(),
			"sbercloud_network_acl_rule": deprecated.ResourceNetworkACLRule(),

			"sbercloud_networking_eip_associate": eip.ResourceEIPAssociate(),

			"sbercloud_networking_secgroup":      vpc.ResourceNetworkingSecGroup(),
			"sbercloud_networking_secgroup_rule": vpc.ResourceNetworkingSecGroupRule(),
			"sbercloud_networking_vip":           vpc.ResourceNetworkingVip(),
			"sbercloud_networking_vip_associate": vpc.ResourceNetworkingVIPAssociateV2(),

			"sbercloud_obs_bucket":        obs.ResourceObsBucket(),
			"sbercloud_obs_bucket_object": obs.ResourceObsBucketObject(),
			"sbercloud_obs_bucket_policy": obs.ResourceObsBucketPolicy(),
			"sbercloud_obs_bucket_acl":    obs.ResourceOBSBucketAcl(),

			"sbercloud_rds_instance":              rds.ResourceRdsInstance(),
			"sbercloud_rds_parametergroup":        rds.ResourceRdsConfiguration(),
			"sbercloud_rds_backup":                rds.ResourceBackup(),
			"sbercloud_rds_read_replica_instance": rds.ResourceRdsReadReplicaInstance(),
			"sbercloud_rds_pg_database":           rds.ResourcePgDatabase(),
			"sbercloud_rds_pg_account_roles":      rds.ResourcePgAccountRoles(),
			"sbercloud_rds_pg_plugin":             rds.ResourceRdsPgPlugin(),
			"sbercloud_rds_pg_plugin_update":      rds.ResourceRdsPgPluginUpdate(),
			"sbercloud_rds_pg_hba":                rds.ResourcePgHba(),
			"sbercloud_rds_pg_sql_limit":          rds.ResourcePgSqlLimit(),

			"sbercloud_rds_pg_account":          rds.ResourcePgAccount(),
			"sbercloud_rds_pg_plugin_parameter": rds.ResourcePgPluginParameter(),

			"sbercloud_rds_mysql_account":                rds.ResourceMysqlAccount(),
			"sbercloud_rds_mysql_binlog":                 rds.ResourceMysqlBinlog(),
			"sbercloud_rds_mysql_database":               rds.ResourceMysqlDatabase(),
			"sbercloud_rds_mysql_database_privilege":     rds.ResourceMysqlDatabasePrivilege(),
			"sbercloud_rds_mysql_database_table_restore": rds.ResourceMysqlDatabaseTableRestore(),
			// "sbercloud_rds_mysql_proxy":                  rds.ResourceMysqlProxy(),
			// "sbercloud_rds_mysql_proxy_restart":          rds.ResourceMysqlProxyRestart(),
			"sbercloud_rds_sqlserver_account":            rds.ResourceSQLServerAccount(),
			"sbercloud_rds_sqlserver_database":           rds.ResourceSQLServerDatabase(),
			"sbercloud_rds_sqlserver_database_privilege": rds.ResourceSQLServerDatabasePrivilege(),
			"sbercloud_rds_sql_audit":                    rds.ResourceSQLAudit(),

			"sbercloud_sfs_turbo":            sfsturbo.ResourceSFSTurbo(),
			"sbercloud_sfs_turbo_dir":        sfsturbo.ResourceSfsTurboDir(),
			"sbercloud_sfs_turbo_dir_quota":  sfsturbo.ResourceSfsTurboDirQuota(),
			"sbercloud_sfs_turbo_data_task":  sfsturbo.ResourceDataTask(),
			"sbercloud_sfs_turbo_du_task":    sfsturbo.ResourceDuTask(),
			"sbercloud_sfs_turbo_obs_target": sfsturbo.ResourceOBSTarget(),
			"sbercloud_sfs_turbo_perm_rule":  sfsturbo.ResourceSFSTurboPermRule(),

			"sbercloud_sfs_access_rule": deprecated.ResourceSFSAccessRuleV2(),
			"sbercloud_sfs_file_system": deprecated.ResourceSFSFileSystemV2(),

			"sbercloud_smn_subscription": smn.ResourceSubscription(),
			"sbercloud_smn_topic":        smn.ResourceTopic(),

			"sbercloud_swr_organization":             swr.ResourceSWROrganization(),
			"sbercloud_swr_organization_permissions": swr.ResourceSWROrganizationPermissions(),
			"sbercloud_swr_repository":               swr.ResourceSWRRepository(),

			"sbercloud_vpc":                             vpc.ResourceVirtualPrivateCloudV1(),
			"sbercloud_vpc_peering_connection":          vpc.ResourceVpcPeeringConnectionV2(),
			"sbercloud_vpc_peering_connection_accepter": vpc.ResourceVpcPeeringConnectionAccepterV2(),
			"sbercloud_vpc_route":                       vpc.ResourceVPCRouteTableRoute(),
			"sbercloud_vpc_route_table":                 vpc.ResourceVPCRouteTable(),
			"sbercloud_vpc_subnet":                      vpc2.ResourceVpcSubnetV1(),
			"sbercloud_vpc_address_group":               vpc.ResourceVpcAddressGroup(),

			"sbercloud_vpc_bandwidth": eip.ResourceVpcBandWidthV2(),
			"sbercloud_vpc_eip":       eip.ResourceVpcEIPV1(),

			"sbercloud_vpcep_endpoint": vpcep.ResourceVPCEndpoint(),
			"sbercloud_vpcep_service":  vpcep.ResourceVPCEndpointService(),

			"sbercloud_vpn_gateway":                 vpn.ResourceGateway(),
			"sbercloud_vpn_customer_gateway":        vpn.ResourceCustomerGateway(),
			"sbercloud_vpn_connection":              vpn.ResourceConnection(),
			"sbercloud_vpn_connection_health_check": vpn.ResourceConnectionHealthCheck(),
			// Legacy
			"sbercloud_identity_role_assignment_v3":  iam.ResourceIdentityGroupRoleAssignment(),
			"sbercloud_identity_user_v3":             iam.ResourceIdentityUser(),
			"sbercloud_identity_group_v3":            iam.ResourceIdentityGroup(),
			"sbercloud_identity_group_membership_v3": iam.ResourceIdentityGroupMembership(),
		},
	}

	provider.ConfigureFunc = func(d *schema.ResourceData) (interface{}, error) {
		terraformVersion := provider.TerraformVersion
		if terraformVersion == "" {
			// Terraform 0.12 introduced this field to the protocol
			// We can therefore assume that if it's missing it's 0.10 or 0.11
			terraformVersion = "0.11+compatible"
		}
		return configureProvider(d, terraformVersion)
	}

	return provider
}

var descriptions map[string]string

func init() {
	descriptions = map[string]string{
		"auth_url": "The Identity authentication URL.",

		"region": "The SberCloud region to connect to.",

		"user_name": "Username to login with.",

		"access_key": "The access key of the SberCloud to use.",

		"secret_key": "The secret key of the SberCloud to use.",

		"security_token": "The security token to authenticate with a temporary security credential.",

		"project_name": "The name of the Project to login with.",

		"password": "Password to login with.",

		"account_name": "The name of the Account to login with.",

		"insecure": "Trust self-signed certificates.",
	}
}

func configureProvider(d *schema.ResourceData, terraformVersion string) (interface{}, error) {
	var project_name string

	// Use region as project_name if it's not set
	if v, ok := d.GetOk("project_name"); ok && v.(string) != "" {
		project_name = v.(string)
	} else {
		project_name = d.Get("region").(string)
	}

	config := config.Config{
		AccessKey:           d.Get("access_key").(string),
		SecretKey:           d.Get("secret_key").(string),
		SecurityToken:       d.Get("security_token").(string),
		DomainName:          d.Get("account_name").(string),
		IdentityEndpoint:    d.Get("auth_url").(string),
		Insecure:            d.Get("insecure").(bool),
		Password:            d.Get("password").(string),
		Region:              d.Get("region").(string),
		TenantName:          project_name,
		Username:            d.Get("user_name").(string),
		TerraformVersion:    terraformVersion,
		Cloud:               "hc.sbercloud.ru",
		MaxRetries:          d.Get("max_retries").(int),
		EnterpriseProjectID: d.Get("enterprise_project_id").(string),
		RegionClient:        true,
		RegionProjectIDMap:  make(map[string]string),
		RPLock:              new(sync.Mutex),
	}

	if err := config.LoadAndValidate(); err != nil {
		return nil, err
	}

	if config.HwClient != nil && config.HwClient.ProjectID != "" {
		config.RegionProjectIDMap[config.Region] = config.HwClient.ProjectID
	}

	return &config, nil
}
