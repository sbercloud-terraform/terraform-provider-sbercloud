package sbercloud

import (
	"sync"

	"github.com/hashicorp/terraform-plugin-sdk/helper/mutexkv"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud"
)

// This is a global MutexKV for use within this plugin.
var osMutexKV = mutexkv.NewMutexKV()

// Provider returns a schema.Provider for SberCloud.
func Provider() terraform.ResourceProvider {
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
			"sbercloud_availability_zones":  huaweicloud.DataSourceAvailabilityZones(),
			"sbercloud_cce_cluster":         huaweicloud.DataSourceCCEClusterV3(),
			"sbercloud_cce_node":            huaweicloud.DataSourceCCENodeV3(),
			"sbercloud_compute_flavors":     huaweicloud.DataSourceEcsFlavors(),
			"sbercloud_identity_role":       huaweicloud.DataSourceIdentityRoleV3(),
			"sbercloud_images_image":        huaweicloud.DataSourceImagesImageV2(),
			"sbercloud_networking_port":     huaweicloud.DataSourceNetworkingPortV2(),
			"sbercloud_networking_secgroup": huaweicloud.DataSourceNetworkingSecGroupV2(),
			"sbercloud_obs_bucket_object":   huaweicloud.DataSourceObsBucketObject(),
			"sbercloud_vpc":                 huaweicloud.DataSourceVirtualPrivateCloudVpcV1(),
			"sbercloud_vpc_bandwidth":       huaweicloud.DataSourceBandWidth(),
			"sbercloud_vpc_subnet":          huaweicloud.DataSourceVpcSubnetV1(),
			"sbercloud_vpc_subnet_ids":      huaweicloud.DataSourceVpcSubnetIdsV1(),
			"sbercloud_vpc_route":           huaweicloud.DataSourceVPCRouteV2(),
			// Legacy
			"sbercloud_identity_role_v3": huaweicloud.DataSourceIdentityRoleV3(),
		},

		ResourcesMap: map[string]*schema.Resource{
			"sbercloud_as_configuration":          huaweicloud.ResourceASConfiguration(),
			"sbercloud_as_group":                  huaweicloud.ResourceASGroup(),
			"sbercloud_as_policy":                 huaweicloud.ResourceASPolicy(),
			"sbercloud_cce_cluster":               huaweicloud.ResourceCCEClusterV3(),
			"sbercloud_cce_node":                  huaweicloud.ResourceCCENodeV3(),
			"sbercloud_dns_recordset":             huaweicloud.ResourceDNSRecordSetV2(),
			"sbercloud_dns_zone":                  huaweicloud.ResourceDNSZoneV2(),
			"sbercloud_identity_role_assignment":  huaweicloud.ResourceIdentityRoleAssignmentV3(),
			"sbercloud_identity_user":             huaweicloud.ResourceIdentityUserV3(),
			"sbercloud_identity_group":            huaweicloud.ResourceIdentityGroupV3(),
			"sbercloud_identity_group_membership": huaweicloud.ResourceIdentityGroupMembershipV3(),
			"sbercloud_images_image":              huaweicloud.ResourceImsImage(),
			"sbercloud_compute_instance":          huaweicloud.ResourceComputeInstanceV2(),
			"sbercloud_compute_interface_attach":  huaweicloud.ResourceComputeInterfaceAttachV2(),
			"sbercloud_compute_keypair":           huaweicloud.ResourceComputeKeypairV2(),
			"sbercloud_compute_servergroup":       huaweicloud.ResourceComputeServerGroupV2(),
			"sbercloud_compute_eip_associate":     huaweicloud.ResourceComputeFloatingIPAssociateV2(),
			"sbercloud_compute_volume_attach":     huaweicloud.ResourceComputeVolumeAttachV2(),
			"sbercloud_evs_snapshot":              huaweicloud.ResourceEvsSnapshotV2(),
			"sbercloud_evs_volume":                huaweicloud.ResourceEvsStorageVolumeV3(),
			"sbercloud_lb_certificate":            huaweicloud.ResourceCertificateV2(),
			"sbercloud_lb_l7policy":               huaweicloud.ResourceL7PolicyV2(),
			"sbercloud_lb_l7rule":                 huaweicloud.ResourceL7RuleV2(),
			"sbercloud_lb_listener":               huaweicloud.ResourceListenerV2(),
			"sbercloud_lb_loadbalancer":           huaweicloud.ResourceLoadBalancerV2(),
			"sbercloud_lb_member":                 huaweicloud.ResourceMemberV2(),
			"sbercloud_lb_monitor":                huaweicloud.ResourceMonitorV2(),
			"sbercloud_lb_pool":                   huaweicloud.ResourcePoolV2(),
			"sbercloud_lb_whitelist":              huaweicloud.ResourceWhitelistV2(),
			"sbercloud_obs_bucket":                huaweicloud.ResourceObsBucket(),
			"sbercloud_obs_bucket_object":         huaweicloud.ResourceObsBucketObject(),
			"sbercloud_obs_bucket_policy":         huaweicloud.ResourceObsBucketPolicy(),
			"sbercloud_vpc":                       huaweicloud.ResourceVirtualPrivateCloudV1(),
			"sbercloud_vpc_bandwidth":             huaweicloud.ResourceVpcBandWidthV2(),
			"sbercloud_vpc_eip":                   huaweicloud.ResourceVpcEIPV1(),
			"sbercloud_vpc_route":                 huaweicloud.ResourceVPCRouteV2(),
			"sbercloud_vpc_peering_connection":    huaweicloud.ResourceVpcPeeringConnectionV2(),
			"sbercloud_vpc_subnet":                huaweicloud.ResourceVpcSubnetV1(),
			"sbercloud_networking_secgroup":       huaweicloud.ResourceNetworkingSecGroupV2(),
			"sbercloud_networking_secgroup_rule":  huaweicloud.ResourceNetworkingSecGroupRuleV2(),
			// Legacy
			"sbercloud_identity_role_assignment_v3":  huaweicloud.ResourceIdentityRoleAssignmentV3(),
			"sbercloud_identity_user_v3":             huaweicloud.ResourceIdentityUserV3(),
			"sbercloud_identity_group_v3":            huaweicloud.ResourceIdentityGroupV3(),
			"sbercloud_identity_group_membership_v3": huaweicloud.ResourceIdentityGroupMembershipV3(),
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

	config := huaweicloud.Config{
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
