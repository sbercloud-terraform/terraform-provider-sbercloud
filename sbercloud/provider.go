package sbercloud

import (
	"sync"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/mutexkv"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/cdm"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/css"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/dcs"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/dds"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/deprecated"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/dis"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/dli"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/dms"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/dws"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/eip"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/eps"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/evs"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/fgs"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/iam"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/rds"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/vpc"
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
		},

		DataSourcesMap: map[string]*schema.Resource{
			"sbercloud_availability_zones":     huaweicloud.DataSourceAvailabilityZones(),
			"sbercloud_cce_cluster":            huaweicloud.DataSourceCCEClusterV3(),
			"sbercloud_cce_node":               huaweicloud.DataSourceCCENodeV3(),
			"sbercloud_cce_node_pool":          huaweicloud.DataSourceCCENodePoolV3(),
			"sbercloud_cdm_flavors":            huaweicloud.DataSourceCdmFlavorV1(),
			"sbercloud_compute_flavors":        huaweicloud.DataSourceEcsFlavors(),
			"sbercloud_dcs_az":                 deprecated.DataSourceDcsAZV1(),
			"sbercloud_dcs_maintainwindow":     dcs.DataSourceDcsMaintainWindow(),
			"sbercloud_dcs_product":            deprecated.DataSourceDcsProductV1(),
			"sbercloud_dds_flavors":            dds.DataSourceDDSFlavorV3(),
			"sbercloud_dms_az":                 deprecated.DataSourceDmsAZ(),
			"sbercloud_dms_product":            dms.DataSourceDmsProduct(),
			"sbercloud_dms_maintainwindow":     dms.DataSourceDmsMaintainWindow(),
			"sbercloud_enterprise_project":     eps.DataSourceEnterpriseProject(),
			"sbercloud_identity_role":          iam.DataSourceIdentityRoleV3(),
			"sbercloud_identity_custom_role":   iam.DataSourceIdentityCustomRole(),
			"sbercloud_identity_group":         iam.DataSourceIdentityGroup(),
			"sbercloud_images_image":           huaweicloud.DataSourceImagesImageV2(),
			"sbercloud_kms_key":                huaweicloud.DataSourceKmsKeyV1(),
			"sbercloud_kms_data_key":           huaweicloud.DataSourceKmsDataKeyV1(),
			"sbercloud_nat_gateway":            huaweicloud.DataSourceNatGatewayV2(),
			"sbercloud_networking_port":        huaweicloud.DataSourceNetworkingPortV2(),
			"sbercloud_networking_secgroup":    huaweicloud.DataSourceNetworkingSecGroupV2(),
			"sbercloud_obs_bucket_object":      huaweicloud.DataSourceObsBucketObject(),
			"sbercloud_rds_flavors":            rds.DataSourceRdsFlavor(),
			"sbercloud_sfs_file_system":        huaweicloud.DataSourceSFSFileSystemV2(),
			"sbercloud_vpc":                    vpc.DataSourceVpcV1(),
			"sbercloud_vpcs":                   vpc.DataSourceVpcs(),
			"sbercloud_vpc_bandwidth":          eip.DataSourceBandWidth(),
			"sbercloud_vpc_eip":                eip.DataSourceVpcEip(),
			"sbercloud_vpc_ids":                vpc.DataSourceVpcIdsV1(),
			"sbercloud_vpc_peering_connection": vpc.DataSourceVpcPeeringConnectionV2(),
			"sbercloud_vpc_route":              vpc.DataSourceVpcRouteV2(),
			"sbercloud_vpc_route_table":        vpc.DataSourceVPCRouteTable(),
			"sbercloud_vpc_subnet":             vpc.DataSourceVpcSubnetV1(),
			"sbercloud_vpc_subnets":            vpc.DataSourceVpcSubnets(),
			"sbercloud_vpc_subnet_ids":         vpc.DataSourceVpcSubnetIdsV1(),
			// Legacy
			"sbercloud_identity_role_v3": iam.DataSourceIdentityRoleV3(),
		},

		ResourcesMap: map[string]*schema.Resource{
			"sbercloud_api_gateway_api":                 huaweicloud.ResourceAPIGatewayAPI(),
			"sbercloud_api_gateway_group":               huaweicloud.ResourceAPIGatewayGroup(),
			"sbercloud_as_configuration":                huaweicloud.ResourceASConfiguration(),
			"sbercloud_as_group":                        huaweicloud.ResourceASGroup(),
			"sbercloud_as_policy":                       huaweicloud.ResourceASPolicy(),
			"sbercloud_css_cluster":                     css.ResourceCssCluster(),
			"sbercloud_cce_cluster":                     huaweicloud.ResourceCCEClusterV3(),
			"sbercloud_cce_node":                        huaweicloud.ResourceCCENodeV3(),
			"sbercloud_cce_node_pool":                   huaweicloud.ResourceCCENodePool(),
			"sbercloud_cdm_cluster":                     cdm.ResourceCdmCluster(),
			"sbercloud_compute_instance":                huaweicloud.ResourceComputeInstanceV2(),
			"sbercloud_compute_interface_attach":        huaweicloud.ResourceComputeInterfaceAttachV2(),
			"sbercloud_compute_keypair":                 huaweicloud.ResourceComputeKeypairV2(),
			"sbercloud_compute_servergroup":             huaweicloud.ResourceComputeServerGroupV2(),
			"sbercloud_compute_eip_associate":           huaweicloud.ResourceComputeFloatingIPAssociateV2(),
			"sbercloud_compute_volume_attach":           huaweicloud.ResourceComputeVolumeAttachV2(),
			"sbercloud_ces_alarmrule":                   huaweicloud.ResourceAlarmRule(),
			"sbercloud_dcs_instance":                    dcs.ResourceDcsInstance(),
			"sbercloud_dds_instance":                    dds.ResourceDdsInstanceV3(),
			"sbercloud_dis_stream":                      dis.ResourceDisStream(),
			"sbercloud_dli_queue":                       dli.ResourceDliQueue(),
			"sbercloud_dms_instance":                    ResourceDmsInstancesV1(),
			"sbercloud_dms_kafka_instance":              dms.ResourceDmsKafkaInstance(),
			"sbercloud_dms_kafka_topic":                 dms.ResourceDmsKafkaTopic(),
			"sbercloud_dms_rabbitmq_instance":           dms.ResourceDmsRabbitmqInstance(),
			"sbercloud_dns_recordset":                   huaweicloud.ResourceDNSRecordSetV2(),
			"sbercloud_dns_zone":                        huaweicloud.ResourceDNSZoneV2(),
			"sbercloud_dws_cluster":                     dws.ResourceDwsCluster(),
			"sbercloud_enterprise_project":              eps.ResourceEnterpriseProject(),
			"sbercloud_evs_snapshot":                    huaweicloud.ResourceEvsSnapshotV2(),
			"sbercloud_evs_volume":                      evs.ResourceEvsVolume(),
			"sbercloud_fgs_function":                    fgs.ResourceFgsFunctionV2(),
			"sbercloud_ges_graph":                       huaweicloud.ResourceGesGraphV1(),
			"sbercloud_identity_access_key":             iam.ResourceIdentityKey(),
			"sbercloud_identity_acl":                    iam.ResourceIdentityACL(),
			"sbercloud_identity_agency":                 iam.ResourceIAMAgencyV3(),
			"sbercloud_identity_group":                  iam.ResourceIdentityGroupV3(),
			"sbercloud_identity_group_membership":       iam.ResourceIdentityGroupMembershipV3(),
			"sbercloud_identity_role":                   iam.ResourceIdentityRole(),
			"sbercloud_identity_role_assignment":        iam.ResourceIdentityRoleAssignmentV3(),
			"sbercloud_identity_user":                   iam.ResourceIdentityUserV3(),
			"sbercloud_images_image":                    huaweicloud.ResourceImsImage(),
			"sbercloud_kms_key":                         huaweicloud.ResourceKmsKeyV1(),
			"sbercloud_lb_certificate":                  huaweicloud.ResourceCertificateV2(),
			"sbercloud_lb_l7policy":                     huaweicloud.ResourceL7PolicyV2(),
			"sbercloud_lb_l7rule":                       huaweicloud.ResourceL7RuleV2(),
			"sbercloud_lb_listener":                     huaweicloud.ResourceListenerV2(),
			"sbercloud_lb_loadbalancer":                 huaweicloud.ResourceLoadBalancerV2(),
			"sbercloud_lb_member":                       huaweicloud.ResourceMemberV2(),
			"sbercloud_lb_monitor":                      huaweicloud.ResourceMonitorV2(),
			"sbercloud_lb_pool":                         huaweicloud.ResourcePoolV2(),
			"sbercloud_lb_whitelist":                    huaweicloud.ResourceWhitelistV2(),
			"sbercloud_nat_dnat_rule":                   huaweicloud.ResourceNatDnatRuleV2(),
			"sbercloud_nat_gateway":                     huaweicloud.ResourceNatGatewayV2(),
			"sbercloud_nat_snat_rule":                   huaweicloud.ResourceNatSnatRuleV2(),
			"sbercloud_network_acl":                     huaweicloud.ResourceNetworkACL(),
			"sbercloud_network_acl_rule":                huaweicloud.ResourceNetworkACLRule(),
			"sbercloud_networking_eip_associate":        eip.ResourceEIPAssociate(),
			"sbercloud_networking_secgroup":             huaweicloud.ResourceNetworkingSecGroupV2(),
			"sbercloud_networking_secgroup_rule":        huaweicloud.ResourceNetworkingSecGroupRuleV2(),
			"sbercloud_obs_bucket":                      huaweicloud.ResourceObsBucket(),
			"sbercloud_obs_bucket_object":               huaweicloud.ResourceObsBucketObject(),
			"sbercloud_obs_bucket_policy":               huaweicloud.ResourceObsBucketPolicy(),
			"sbercloud_rds_instance":                    huaweicloud.ResourceRdsInstanceV3(),
			"sbercloud_rds_parametergroup":              huaweicloud.ResourceRdsConfigurationV3(),
			"sbercloud_rds_read_replica_instance":       huaweicloud.ResourceRdsReadReplicaInstance(),
			"sbercloud_sfs_access_rule":                 huaweicloud.ResourceSFSAccessRuleV2(),
			"sbercloud_sfs_file_system":                 huaweicloud.ResourceSFSFileSystemV2(),
			"sbercloud_sfs_turbo":                       huaweicloud.ResourceSFSTurbo(),
			"sbercloud_smn_subscription":                huaweicloud.ResourceSubscription(),
			"sbercloud_smn_topic":                       huaweicloud.ResourceTopic(),
			"sbercloud_vpc":                             vpc.ResourceVirtualPrivateCloudV1(),
			"sbercloud_vpc_bandwidth":                   eip.ResourceVpcBandWidthV2(),
			"sbercloud_vpc_eip":                         eip.ResourceVpcEIPV1(),
			"sbercloud_vpc_peering_connection":          vpc.ResourceVpcPeeringConnectionV2(),
			"sbercloud_vpc_peering_connection_accepter": vpc.ResourceVpcPeeringConnectionAccepterV2(),
			"sbercloud_vpc_route":                       vpc.ResourceVPCRouteV2(),
			"sbercloud_vpc_route_table":                 vpc.ResourceVPCRouteTable(),
			"sbercloud_vpc_subnet":                      vpc.ResourceVpcSubnetV1(),
			// Legacy
			"sbercloud_identity_role_assignment_v3":  iam.ResourceIdentityRoleAssignmentV3(),
			"sbercloud_identity_user_v3":             iam.ResourceIdentityUserV3(),
			"sbercloud_identity_group_v3":            iam.ResourceIdentityGroupV3(),
			"sbercloud_identity_group_membership_v3": iam.ResourceIdentityGroupMembershipV3(),
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
