package rds

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/bss/v2/orders"
	"github.com/chnsz/golangsdk/openstack/common/tags"
	"github.com/chnsz/golangsdk/openstack/rds/v3/backups"
	"github.com/chnsz/golangsdk/openstack/rds/v3/instances"
	"github.com/chnsz/golangsdk/openstack/rds/v3/securities"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/hashcode"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

type ctxType string

// ResourceRdsInstance is the impl for huaweicloud_rds_instance resource
func ResourceRdsInstance() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRdsInstanceCreate,
		ReadContext:   resourceRdsInstanceRead,
		UpdateContext: resourceRdsInstanceUpdate,
		DeleteContext: resourceRdsInstanceDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create:  schema.DefaultTimeout(30 * time.Minute),
			Update:  schema.DefaultTimeout(30 * time.Minute),
			Delete:  schema.DefaultTimeout(30 * time.Minute),
			Default: schema.DefaultTimeout(15 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"availability_zone": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"flavor": {
				Type:     schema.TypeString,
				Required: true,
			},

			"db": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:             schema.TypeString,
							Required:         true,
							ForceNew:         true,
							DiffSuppressFunc: utils.SuppressCaseDiffs,
						},
						"version": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"password": {
							Type:      schema.TypeString,
							Sensitive: true,
							Optional:  true,
						},
						"port": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"user_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"volume": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"size": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"type": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"disk_encryption_id": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ForceNew: true,
						},
						"limit_size": {
							Type:         schema.TypeInt,
							Optional:     true,
							RequiredWith: []string{"volume.0.trigger_threshold"},
						},
						"trigger_threshold": {
							Type:         schema.TypeInt,
							Optional:     true,
							RequiredWith: []string{"volume.0.limit_size"},
						},
					},
				},
			},

			"restore": {
				Type:          schema.TypeList,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"period"},
				MaxItems:      1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance_id": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"backup_id": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"database_name": {
							Type:     schema.TypeMap,
							Optional: true,
							ForceNew: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},

			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"subnet_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"security_group_id": {
				Type:     schema.TypeString,
				Required: true,
			},

			"backup_strategy": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"start_time": {
							Type:     schema.TypeString,
							Required: true,
						},
						"keep_days": {
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							Description: "schema: Required",
						},
						"period": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},

			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"fixed_ip": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: utils.ValidateIP,
			},

			"ha_replication_mode": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"param_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"collation": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"switch_strategy": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"ssl_enable": {
				Type:     schema.TypeBool,
				Optional: true,
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"dss_pool_id": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"tags": common.TagsSchema(),

			"time_zone": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"parameters": {
				Type: schema.TypeSet,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"value": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
				Set:      parameterToHash,
				Optional: true,
				Computed: true,
			},

			"maintain_begin": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"maintain_end": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				RequiredWith: []string{"maintain_begin"},
			},

			"nodes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"availability_zone": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"role": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"private_ips": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"public_ips": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"created": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"lower_case_table_names": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			// charge info: charging_mode, period_unit, period, auto_renew, auto_pay
			"charging_mode": common.SchemaChargingMode(nil),
			"period_unit":   common.SchemaPeriodUnit(nil),
			"period":        common.SchemaPeriod(nil),
			"auto_renew":    common.SchemaAutoRenewUpdatable(nil),
			"auto_pay":      common.SchemaAutoPay(nil),
		},
	}
}

func buildRdsInstanceDBPort(d *schema.ResourceData) string {
	if v, ok := d.GetOk("db.0.port"); ok {
		return strconv.Itoa(v.(int))
	}
	return ""
}

func isMySQLDatabase(d *schema.ResourceData) bool {
	dbType := d.Get("db.0.type").(string)
	// Database type is not case sensitive.
	return strings.ToLower(dbType) == "mysql"
}

func resourceRdsInstanceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	region := config.GetRegion(d)
	client, err := config.RdsV3Client(region)
	if err != nil {
		return diag.Errorf("error creating RDS client: %s", err)
	}

	if d.Get("ssl_enable").(bool) && !isMySQLDatabase(d) {
		return diag.Errorf("only MySQL database support SSL enable and disable")
	}

	createOpts := instances.CreateOpts{
		Name:                d.Get("name").(string),
		FlavorRef:           d.Get("flavor").(string),
		VpcId:               d.Get("vpc_id").(string),
		SubnetId:            d.Get("subnet_id").(string),
		SecurityGroupId:     d.Get("security_group_id").(string),
		ConfigurationId:     d.Get("param_group_id").(string),
		TimeZone:            d.Get("time_zone").(string),
		FixedIp:             d.Get("fixed_ip").(string),
		DiskEncryptionId:    d.Get("volume.0.disk_encryption_id").(string),
		Collation:           d.Get("collation").(string),
		Port:                buildRdsInstanceDBPort(d),
		EnterpriseProjectId: config.GetEnterpriseProjectID(d),
		Region:              region,
		AvailabilityZone:    buildRdsInstanceAvailabilityZone(d),
		Datastore:           buildRdsInstanceDatastore(d),
		Volume:              buildRdsInstanceVolume(d),
		Ha:                  buildRdsInstanceHaReplicationMode(d),
		UnchangeableParam:   buildRdsInstanceUnchangeableParam(d),
		RestorePoint:        buildRdsInstanceRestorePoint(d),
		DssPoolId:           d.Get("dss_pool_id").(string),
	}

	// PrePaid
	if d.Get("charging_mode") == "prePaid" {
		if err := common.ValidatePrePaidChargeInfo(d); err != nil {
			return diag.FromErr(err)
		}

		chargeInfo := &instances.ChargeInfo{
			ChargeMode: d.Get("charging_mode").(string),
			PeriodType: d.Get("period_unit").(string),
			PeriodNum:  d.Get("period").(int),
		}
		if d.Get("auto_pay").(string) != "false" {
			chargeInfo.IsAutoPay = true
		}
		if d.Get("auto_renew").(string) == "true" {
			chargeInfo.IsAutoRenew = true
		}
		createOpts.ChargeInfo = chargeInfo
	}

	log.Printf("[DEBUG] Create Options: %#v", createOpts)
	// Add password here so it wouldn't go in the above log entry
	createOpts.Password = d.Get("db.0.password").(string)

	res, err := instances.Create(client, createOpts).Extract()
	if err != nil {
		return diag.Errorf("error creating RDS instance: %s", err)
	}
	d.SetId(res.Instance.Id)
	instanceID := d.Id()

	// wait for order success
	if res.OrderId != "" {
		bssClient, err := config.BssV2Client(config.GetRegion(d))
		if err != nil {
			return diag.Errorf("error creating BSS V2 client: %s", err)
		}
		if err := orders.WaitForOrderSuccess(bssClient, int(d.Timeout(schema.TimeoutCreate)/time.Second), res.OrderId); err != nil {
			return diag.Errorf("error waiting for RDS order %s succuss: %s", res.OrderId, err)
		}
	}

	if res.JobId != "" {
		if err := checkRDSInstanceJobFinish(client, res.JobId, d.Timeout(schema.TimeoutCreate)); err != nil {
			return diag.Errorf("error creating instance (%s): %s", instanceID, err)
		}
	}
	// for prePaid charge mode
	stateConf := &resource.StateChangeConf{
		Target:       []string{"ACTIVE", "BACKING UP"},
		Refresh:      rdsInstanceStateRefreshFunc(client, instanceID),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        20 * time.Second,
		PollInterval: 10 * time.Second,
	}
	if _, err = stateConf.WaitForStateContext(ctx); err != nil {
		return diag.Errorf("error waiting for RDS instance (%s) creation completed: %s", instanceID, err)
	}

	if err = updateRdsInstanceDescription(d, client, instanceID); err != nil {
		return diag.FromErr(err)
	}

	if err = updateRdsInstanceSSLConfig(ctx, d, client, instanceID); err != nil {
		return diag.FromErr(err)
	}

	if err = updateRdsInstanceMaintainWindow(d, client, instanceID); err != nil {
		return diag.FromErr(err)
	}

	if v, ok := d.GetOk("switch_strategy"); ok && v.(string) != "reliability" {
		if err = updateRdsInstanceSwitchStrategy(ctx, d, client, instanceID); err != nil {
			return diag.FromErr(err)
		}
	}

	tagRaw := d.Get("tags").(map[string]interface{})
	if len(tagRaw) > 0 {
		taglist := utils.ExpandResourceTags(tagRaw)
		if tagErr := tags.Create(client, "instances", instanceID, taglist).ExtractErr(); tagErr != nil {
			return diag.Errorf("error setting tags of RDS instance (%s): %s", instanceID, tagErr)
		}
	}

	// Set Parameters
	if parameters := d.Get("parameters").(*schema.Set); parameters.Len() > 0 {
		if err = initializeParameters(ctx, d.Timeout(schema.TimeoutCreate), client, instanceID, parameters); err != nil {
			return diag.FromErr(err)
		}
	}

	if size := d.Get("volume.0.limit_size").(int); size > 0 {
		if err = enableVolumeAutoExpand(ctx, d, client, instanceID, size); err != nil {
			return diag.FromErr(err)
		}
	}

	if err := updateRdsInstanceBackupStrategy(d, client, instanceID); err != nil {
		return diag.FromErr(err)
	}

	return resourceRdsInstanceRead(ctx, d, meta)
}

func resourceRdsInstanceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	client, err := config.RdsV3Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating RDS client: %s", err)
	}

	instanceID := d.Id()
	instance, err := GetRdsInstanceByID(client, instanceID)
	if err != nil {
		return diag.Errorf("error getting RDS instance: %s", err)
	}
	if instance.Id == "" {
		d.SetId("")
		return nil
	}

	log.Printf("[DEBUG] Retrieved RDS instance (%s): %#v", instanceID, instance)
	d.Set("region", instance.Region)
	d.Set("name", instance.Name)
	d.Set("description", instance.Alias)
	d.Set("status", instance.Status)
	d.Set("created", instance.Created)
	d.Set("ha_replication_mode", instance.Ha.ReplicationMode)
	d.Set("vpc_id", instance.VpcId)
	d.Set("subnet_id", instance.SubnetId)
	d.Set("security_group_id", instance.SecurityGroupId)
	d.Set("flavor", instance.FlavorRef)
	d.Set("time_zone", instance.TimeZone)
	d.Set("collation", instance.Collation)
	d.Set("enterprise_project_id", instance.EnterpriseProjectId)
	d.Set("switch_strategy", instance.SwitchStrategy)
	d.Set("charging_mode", instance.ChargeInfo.ChargeMode)
	d.Set("tags", utils.TagsToMap(instance.Tags))

	publicIps := make([]interface{}, len(instance.PublicIps))
	for i, v := range instance.PublicIps {
		publicIps[i] = v
	}
	d.Set("public_ips", publicIps)

	privateIps := make([]string, len(instance.PrivateIps))
	for i, v := range instance.PrivateIps {
		privateIps[i] = v
	}
	d.Set("private_ips", privateIps)
	// If the creation of the RDS instance is failed, the length of the private IP list will be zero.
	if len(privateIps) > 0 {
		d.Set("fixed_ip", privateIps[0])
	}

	maintainWindow := strings.Split(instance.MaintenanceWindow, "-")
	if len(maintainWindow) == 2 {
		d.Set("maintain_begin", maintainWindow[0])
		d.Set("maintain_end", maintainWindow[1])
	}

	volume := map[string]interface{}{
		"type":               instance.Volume.Type,
		"size":               instance.Volume.Size,
		"disk_encryption_id": instance.DiskEncryptionId,
	}
	// Only MySQL engines are supported.
	resp, err := instances.GetAutoExpand(client, instanceID)
	if err != nil {
		log.Printf("[ERROR] error query automatic expansion configuration of the instance storage: %s", err)
	}
	if resp.SwitchOption {
		volume["limit_size"] = resp.LimitSize
		volume["trigger_threshold"] = resp.TriggerThreshold
	}
	if err := d.Set("volume", []map[string]interface{}{volume}); err != nil {
		return diag.Errorf("error saving volume to RDS instance (%s): %s", instanceID, err)
	}

	dbList := make([]map[string]interface{}, 1)
	database := map[string]interface{}{
		"type":      instance.DataStore.Type,
		"version":   instance.DataStore.Version,
		"port":      instance.Port,
		"user_name": instance.DbUserName,
	}
	if len(d.Get("db").([]interface{})) > 0 {
		database["password"] = d.Get("db.0.password")
	}
	dbList[0] = database
	if err := d.Set("db", dbList); err != nil {
		return diag.Errorf("error saving data base to RDS instance (%s): %s", instanceID, err)
	}

	backupStrategy, err := backups.Get(client, instanceID).Extract()
	if err != nil {
		return diag.Errorf("error getting RDS backup strategy: %s", err)
	}

	backup := make([]map[string]interface{}, 1)
	backup[0] = map[string]interface{}{
		"start_time": instance.BackupStrategy.StartTime,
		"keep_days":  instance.BackupStrategy.KeepDays,
		"period":     backupStrategy.Period,
	}
	if err := d.Set("backup_strategy", backup); err != nil {
		return diag.Errorf("error saving backup strategy to RDS instance (%s): %s", instanceID, err)
	}

	nodes := make([]map[string]interface{}, len(instance.Nodes))
	for i, v := range instance.Nodes {
		nodes[i] = map[string]interface{}{
			"id":                v.Id,
			"name":              v.Name,
			"role":              v.Role,
			"status":            v.Status,
			"availability_zone": v.AvailabilityZone,
		}
	}
	if err := d.Set("nodes", nodes); err != nil {
		return diag.Errorf("error saving nodes to RDS instance (%s): %s", instanceID, err)
	}

	return setRdsInstanceParameters(ctx, d, client, instanceID)
}

func setRdsInstanceParameters(ctx context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient,
	instanceID string) diag.Diagnostics {
	// Set Parameters
	configs, err := instances.GetConfigurations(client, instanceID).Extract()
	if err != nil {
		log.Printf("[WARN] error fetching parameters of instance (%s): %s", instanceID, err)
		return nil
	}

	var restart []string
	var params []map[string]interface{}
	for _, parameter := range d.Get("parameters").(*schema.Set).List() {
		name := parameter.(map[string]interface{})["name"]
		for _, v := range configs.Parameters {
			if v.Name == name {
				p := map[string]interface{}{
					"name":  v.Name,
					"value": v.Value,
				}
				params = append(params, p)
				if v.Restart {
					restart = append(restart, v.Name)
				}
				break
			}
		}
	}

	if len(params) > 0 {
		if err = d.Set("parameters", params); err != nil {
			log.Printf("error saving parameters to RDS instance (%s): %s", instanceID, err)
		}
		if len(restart) > 0 && ctx.Value(ctxType("parametersChanged")) == "true" {
			return diag.Diagnostics{
				diag.Diagnostic{
					Severity: diag.Warning,
					Summary:  "Parameters Changed",
					Detail:   fmt.Sprintf("Parameters %s changed which needs reboot.", restart),
				},
			}
		}
	}
	return nil
}

func resourceRdsInstanceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	client, err := config.RdsV3Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating RDS Client: %s", err)
	}

	instanceID := d.Id()

	if err := updateRdsInstanceName(d, client, instanceID); err != nil {
		return diag.FromErr(err)
	}

	if err = updateRdsInstanceDescription(d, client, instanceID); err != nil {
		return diag.FromErr(err)
	}

	if err := updateRdsInstanceFlavor(ctx, d, config, client, instanceID, true); err != nil {
		return diag.FromErr(err)
	}

	if err := updateRdsInstanceVolumeSize(ctx, d, client, instanceID); err != nil {
		return diag.FromErr(err)
	}

	if err := updateRdsInstanceBackupStrategy(d, client, instanceID); err != nil {
		return diag.FromErr(err)
	}

	if err = updateRdsInstanceMaintainWindow(d, client, instanceID); err != nil {
		return diag.FromErr(err)
	}

	if err = updateRdsInstanceReplicationMode(ctx, d, client, instanceID); err != nil {
		return diag.FromErr(err)
	}

	if err = updateRdsInstanceSwitchStrategy(ctx, d, client, instanceID); err != nil {
		return diag.FromErr(err)
	}

	if err = updateRdsInstanceCollation(ctx, d, client, instanceID); err != nil {
		return diag.FromErr(err)
	}

	if err := updateRdsInstanceDBPort(ctx, d, client, instanceID); err != nil {
		return diag.FromErr(err)
	}

	if err := updateRdsInstanceFixedIp(ctx, d, client, instanceID); err != nil {
		return diag.FromErr(err)
	}

	if err := updateRdsInstanceSecurityGroup(ctx, d, client, instanceID); err != nil {
		return diag.FromErr(err)
	}

	if err := updateRdsInstanceSSLConfig(ctx, d, client, instanceID); err != nil {
		return diag.FromErr(err)
	}

	if err := updateRdsRootPassword(ctx, d, client, instanceID); err != nil {
		return diag.FromErr(err)
	}

	if d.HasChange("tags") {
		tagErr := utils.UpdateResourceTags(client, d, "instances", instanceID)
		if tagErr != nil {
			return diag.Errorf("error updating tags of RDS instance (%s): %s", instanceID, tagErr)
		}
	}

	if d.HasChange("auto_renew") {
		bssClient, err := config.BssV2Client(config.GetRegion(d))
		if err != nil {
			return diag.Errorf("error creating BSS V2 client: %s", err)
		}
		if err = common.UpdateAutoRenew(bssClient, d.Get("auto_renew").(string), d.Id()); err != nil {
			return diag.Errorf("error updating the auto-renew of the instance (%s): %s", d.Id(), err)
		}
	}

	if ctx, err = updateRdsParameters(ctx, d, client, instanceID); err != nil {
		return diag.FromErr(err)
	}

	if err = updateVolumeAutoExpand(ctx, d, client, instanceID); err != nil {
		return diag.FromErr(err)
	}

	return resourceRdsInstanceRead(ctx, d, meta)
}

func resourceRdsInstanceDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	client, err := config.RdsV3Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating rds client: %s ", err)
	}

	id := d.Id()
	log.Printf("[DEBUG] Deleting Instance %s", id)
	if v, ok := d.GetOk("charging_mode"); ok && v.(string) == "prePaid" {
		resourceIds := []string{id}
		// the image of SQL server is come from cloud market, when creating an SQL server instance resource, two order
		// will be created, one is instance order, the other is market image order, so it is needed to unsubscribe the
		// two order when unsubscribe the instance
		if strings.ToLower(d.Get("db.0.type").(string)) == "sqlserver" {
			resourceIds = append(resourceIds, fmt.Sprintf("%s%s", id, ".marketimage"))
		}
		retryFunc := func() (interface{}, bool, error) {
			err = common.UnsubscribePrePaidResource(d, config, resourceIds)
			retry, err := handleDeletionError(err)
			return nil, retry, err
		}
		_, err = common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
			Ctx:          ctx,
			RetryFunc:    retryFunc,
			WaitFunc:     rdsInstanceStateRefreshFunc(client, id),
			WaitTarget:   []string{"ACTIVE"},
			Timeout:      d.Timeout(schema.TimeoutDelete),
			DelayTimeout: 10 * time.Second,
			PollInterval: 10 * time.Second,
		})
		if err != nil {
			return diag.Errorf("error unsubscribe RDS instance: %s", err)
		}
	} else {
		retryFunc := func() (interface{}, bool, error) {
			result := instances.Delete(client, id)
			retry, err := handleDeletionError(result.Err)
			return nil, retry, err
		}
		_, err = common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
			Ctx:          ctx,
			RetryFunc:    retryFunc,
			WaitFunc:     rdsInstanceStateRefreshFunc(client, id),
			WaitTarget:   []string{"ACTIVE"},
			Timeout:      d.Timeout(schema.TimeoutDelete),
			DelayTimeout: 10 * time.Second,
			PollInterval: 10 * time.Second,
		})
		if err != nil {
			return diag.FromErr(err)
		}
	}

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"ACTIVE", "BACKING UP"},
		Target:       []string{"DELETED"},
		Refresh:      rdsInstanceStateRefreshFunc(client, id),
		Timeout:      d.Timeout(schema.TimeoutDelete),
		Delay:        15 * time.Second,
		PollInterval: 5 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf(
			"error waiting for rds instance (%s) to be deleted: %s ",
			id, err)
	}

	log.Printf("[DEBUG] Successfully deleted RDS instance %s", id)
	return nil
}

func GetRdsInstanceByID(client *golangsdk.ServiceClient, instanceID string) (*instances.RdsInstanceResponse, error) {
	listOpts := instances.ListOpts{
		Id: instanceID,
	}
	pages, err := instances.List(client, listOpts).AllPages()
	if err != nil {
		return nil, fmt.Errorf("An error occurred while querying rds instance %s: %s", instanceID, err)
	}

	resp, err := instances.ExtractRdsInstances(pages)
	if err != nil {
		return nil, err
	}

	instanceList := resp.Instances
	if len(instanceList) == 0 {
		// return an empty rds instance
		log.Printf("[WARN] can not find the specified rds instance %s", instanceID)
		instance := new(instances.RdsInstanceResponse)
		return instance, nil
	}

	if len(instanceList) > 1 {
		return nil, fmt.Errorf("retrieving more than one rds instance by %s", instanceID)
	}
	if instanceList[0].Id != instanceID {
		return nil, fmt.Errorf("the id of rds instance was expected %s, but got %s",
			instanceID, instanceList[0].Id)
	}

	return &instanceList[0], nil
}

func buildRdsInstanceAvailabilityZone(d *schema.ResourceData) string {
	azList := utils.ExpandToStringList(d.Get("availability_zone").([]interface{}))
	return strings.Join(azList, ",")
}

func buildRdsInstanceDatastore(d *schema.ResourceData) *instances.Datastore {
	var database *instances.Datastore
	dbRaw := d.Get("db").([]interface{})

	if len(dbRaw) == 1 {
		database = new(instances.Datastore)
		database.Type = dbRaw[0].(map[string]interface{})["type"].(string)
		database.Version = dbRaw[0].(map[string]interface{})["version"].(string)
	}
	return database
}

func buildRdsInstanceVolume(d *schema.ResourceData) *instances.Volume {
	var volume *instances.Volume
	volumeRaw := d.Get("volume").([]interface{})

	if len(volumeRaw) == 1 {
		volume = new(instances.Volume)
		volume.Type = volumeRaw[0].(map[string]interface{})["type"].(string)
		volume.Size = volumeRaw[0].(map[string]interface{})["size"].(int)
	}
	return volume
}

func buildRdsInstanceBackupStrategy(d *schema.ResourceData) *instances.BackupStrategy {
	var backupStrategy *instances.BackupStrategy
	backupRaw := d.Get("backup_strategy").([]interface{})

	if len(backupRaw) == 1 {
		backupStrategy = new(instances.BackupStrategy)
		backupStrategy.StartTime = backupRaw[0].(map[string]interface{})["start_time"].(string)
		backupStrategy.KeepDays = backupRaw[0].(map[string]interface{})["keep_days"].(int)
	}
	return backupStrategy
}

func buildRdsInstanceUnchangeableParam(d *schema.ResourceData) *instances.UnchangeableParam {
	var unchangeableParam *instances.UnchangeableParam
	if v, ok := d.GetOk("lower_case_table_names"); ok {
		unchangeableParam = new(instances.UnchangeableParam)
		unchangeableParam.LowerCaseTableNames = v.(string)
	}
	return unchangeableParam
}

func buildRdsInstanceRestorePoint(d *schema.ResourceData) *instances.RestorePoint {
	if restoreRaw, ok := d.GetOk("restore"); ok {
		if v, ok := restoreRaw.([]interface{})[0].(map[string]interface{}); ok {
			restorePoint := instances.RestorePoint{
				Type:         "backup",
				InstanceId:   v["instance_id"].(string),
				BackupId:     v["backup_id"].(string),
				DatabaseName: utils.ExpandToStringMap(v["database_name"].(map[string]interface{})),
			}
			return &restorePoint
		}
	}
	return nil
}

func buildRdsInstanceHaReplicationMode(d *schema.ResourceData) *instances.Ha {
	var ha *instances.Ha
	if v, ok := d.GetOk("ha_replication_mode"); ok {
		ha = new(instances.Ha)
		ha.Mode = "ha"
		ha.ReplicationMode = v.(string)
	}
	return ha
}

func buildRdsInstanceParameters(params *schema.Set) instances.ModifyConfigurationOpts {
	var configOpts instances.ModifyConfigurationOpts

	values := make(map[string]string)
	for _, v := range params.List() {
		key := v.(map[string]interface{})["name"].(string)
		value := v.(map[string]interface{})["value"].(string)
		values[key] = value
	}
	configOpts.Values = values
	return configOpts
}

func initializeParameters(ctx context.Context, timeout time.Duration, client *golangsdk.ServiceClient,
	instanceID string, parametersRaw *schema.Set) error {
	configOpts := buildRdsInstanceParameters(parametersRaw)
	retryFunc := func() (interface{}, bool, error) {
		_, err := instances.ModifyConfiguration(client, instanceID, configOpts).Extract()
		retry, err := handleMultiOperationsError(err)
		return nil, retry, err
	}
	_, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     rdsInstanceStateRefreshFunc(client, instanceID),
		WaitTarget:   []string{"ACTIVE"},
		Timeout:      timeout,
		DelayTimeout: 10 * time.Second,
		PollInterval: 10 * time.Second,
	})
	if err != nil {
		return fmt.Errorf("error modifying parameters for RDS instance (%s): %s", instanceID, err)
	}

	// Check if we need to restart
	restart, err := checkRdsInstanceRestart(client, instanceID, parametersRaw.List())
	if err != nil {
		return err
	}

	if restart {
		return restartRdsInstance(ctx, timeout, client, instanceID)
	}
	return nil
}

func checkRdsInstanceRestart(client *golangsdk.ServiceClient, instanceID string, parameters []interface{}) (bool, error) {
	configs, err := instances.GetConfigurations(client, instanceID).Extract()
	if err != nil {
		return false, fmt.Errorf("error fetching the instance parameters (%s): %s", instanceID, err)
	}

	for _, parameter := range parameters {
		name := parameter.(map[string]interface{})["name"]
		for _, v := range configs.Parameters {
			if v.Name == name && v.Restart {
				return true, nil
			}
		}
	}
	return false, nil
}

func restartRdsInstance(ctx context.Context, timeout time.Duration, client *golangsdk.ServiceClient,
	instanceID string) error {
	// If parameters which requires restart changed, reboot the instance.
	retryFunc := func() (interface{}, bool, error) {
		_, err := instances.RebootInstance(client, instanceID).Extract()
		retry, err := handleMultiOperationsError(err)
		return nil, retry, err
	}
	_, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     rdsInstanceStateRefreshFunc(client, instanceID),
		WaitTarget:   []string{"ACTIVE"},
		Timeout:      timeout,
		DelayTimeout: 10 * time.Second,
		PollInterval: 10 * time.Second,
	})
	if err != nil {
		return fmt.Errorf("error rebooting for RDS instance (%s): %s", instanceID, err)
	}

	// wait for the instance state to be 'ACTIVE'.
	stateConf := &resource.StateChangeConf{
		Target:       []string{"ACTIVE"},
		Refresh:      rdsInstanceStateRefreshFunc(client, instanceID),
		Timeout:      timeout,
		Delay:        5 * time.Second,
		PollInterval: 5 * time.Second,
	}
	if _, err = stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("error waiting for RDS instance (%s) become active status: %s", instanceID, err)
	}
	return nil
}

func updateRdsInstanceName(d *schema.ResourceData, client *golangsdk.ServiceClient, instanceID string) error {
	if !d.HasChange("name") {
		return nil
	}

	renameOpts := instances.RenameInstanceOpts{
		Name: d.Get("name").(string),
	}
	r := instances.Rename(client, renameOpts, instanceID)
	if r.Result.Err != nil {
		return fmt.Errorf("error renaming RDS instance (%s): %s", instanceID, r.Err)
	}

	return nil
}

func updateRdsInstanceDescription(d *schema.ResourceData, client *golangsdk.ServiceClient, instanceID string) error {
	if !d.HasChange("description") {
		return nil
	}

	modifyAliasOpts := instances.ModifyAliasOpts{
		Alias: d.Get("description").(string),
	}
	log.Printf("[DEBUG] Modify RDS instance description opts: %+v", modifyAliasOpts)
	r := instances.ModifyAlias(client, modifyAliasOpts, instanceID)
	if r.Err != nil {
		return fmt.Errorf("error modify RDS instance (%s) description: %s", instanceID, r.Err)
	}

	return nil
}

func updateRdsInstanceFlavor(ctx context.Context, d *schema.ResourceData, cfg *config.Config,
	client *golangsdk.ServiceClient, instanceID string, isSupportAutoPay bool) error {
	if !d.HasChange("flavor") {
		return nil
	}

	resizeFlavor := instances.SpecCode{
		Speccode:  d.Get("flavor").(string),
		IsAutoPay: true,
	}
	if isSupportAutoPay && d.Get("auto_pay").(string) == "false" {
		resizeFlavor.IsAutoPay = false
	}
	var resizeFlavorOpts instances.ResizeFlavorOpts
	resizeFlavorOpts.ResizeFlavor = &resizeFlavor

	retryFunc := func() (interface{}, bool, error) {
		res, err := instances.Resize(client, resizeFlavorOpts, instanceID).Extract()
		retry, err := handleMultiOperationsError(err)
		return res, retry, err
	}
	r, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     rdsInstanceStateRefreshFunc(client, instanceID),
		WaitTarget:   []string{"ACTIVE"},
		Timeout:      d.Timeout(schema.TimeoutUpdate),
		DelayTimeout: 10 * time.Second,
		PollInterval: 10 * time.Second,
	})
	if err != nil {
		return fmt.Errorf("error updating instance Flavor from result: %s ", err)
	}

	res := r.(*instances.ResizeFlavor)
	// wait for order success
	if res.OrderId != "" {
		bssClient, err := cfg.BssV2Client(cfg.GetRegion(d))
		if err != nil {
			return fmt.Errorf("error creating BSS V2 client: %s", err)
		}
		if err := orders.WaitForOrderSuccess(bssClient, int(d.Timeout(schema.TimeoutUpdate)/time.Second), res.OrderId); err != nil {
			return fmt.Errorf("error waiting for RDS order %s succuss: %s", res.OrderId, err)
		}
	}

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"MODIFYING"},
		Target:       []string{"ACTIVE"},
		Refresh:      rdsInstanceStateRefreshFunc(client, instanceID),
		Timeout:      d.Timeout(schema.TimeoutUpdate),
		Delay:        15 * time.Second,
		PollInterval: 15 * time.Second,
	}
	if _, err = stateConf.WaitForState(); err != nil {
		return fmt.Errorf("error waiting for instance (%s) flavor to be Updated: %s ", instanceID, err)
	}
	return nil
}

func updateRdsInstanceVolumeSize(ctx context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient,
	instanceID string) error {
	if !d.HasChange("volume.0.size") {
		return nil
	}

	volumeRaw := d.Get("volume").([]interface{})
	volumeItem := volumeRaw[0].(map[string]interface{})
	enlargeOpts := instances.EnlargeVolumeOpts{
		EnlargeVolume: &instances.EnlargeVolumeSize{
			Size: volumeItem["size"].(int),
		},
	}

	log.Printf("[DEBUG] Enlarge Volume opts: %+v", enlargeOpts)

	retryFunc := func() (interface{}, bool, error) {
		instance, err := instances.EnlargeVolume(client, enlargeOpts, instanceID).Extract()
		retry, err := handleMultiOperationsError(err)
		return instance, retry, err
	}
	r, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     rdsInstanceStateRefreshFunc(client, instanceID),
		WaitTarget:   []string{"ACTIVE"},
		Timeout:      d.Timeout(schema.TimeoutUpdate),
		DelayTimeout: 10 * time.Second,
		PollInterval: 10 * time.Second,
	})
	if err != nil {
		return fmt.Errorf("error updating instance volume from result: %s ", err)
	}

	instance := r.(*instances.EnlargeVolumeResp)
	if err := checkRDSInstanceJobFinish(client, instance.JobId, d.Timeout(schema.TimeoutUpdate)); err != nil {
		return fmt.Errorf("error updating instance (%s): %s", instanceID, err)
	}

	return nil
}

func updateRdsInstanceBackupStrategy(d *schema.ResourceData, client *golangsdk.ServiceClient, instanceID string) error {
	if !d.HasChange("backup_strategy") {
		return nil
	}

	backupRaw := d.Get("backup_strategy").([]interface{})
	rawMap := backupRaw[0].(map[string]interface{})
	keepDays := rawMap["keep_days"].(int)
	period := rawMap["period"].(string)
	if period == "" {
		period = "1,2,3,4,5,6,7"
	}
	updateOpts := backups.UpdateOpts{
		KeepDays:  &keepDays,
		StartTime: rawMap["start_time"].(string),
		Period:    period,
	}

	log.Printf("[DEBUG] updateOpts: %#v", updateOpts)
	err := backups.Update(client, instanceID, updateOpts).ExtractErr()
	if err != nil {
		return fmt.Errorf("error updating RDS instance backup strategy (%s): %s", instanceID, err)
	}

	return nil
}

func updateRdsInstanceDBPort(ctx context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient,
	instanceID string) error {
	if !d.HasChange("db.0.port") {
		return nil
	}

	updateOpts := securities.PortOpts{
		Port: d.Get("db.0.port").(int),
	}
	log.Printf("[DEBUG] Update opts of Database port: %+v", updateOpts)

	retryFunc := func() (interface{}, bool, error) {
		_, err := securities.UpdatePort(client, instanceID, updateOpts).Extract()
		retry, err := handleMultiOperationsError(err)
		return nil, retry, err
	}
	_, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     rdsInstanceStateRefreshFunc(client, instanceID),
		WaitTarget:   []string{"ACTIVE"},
		Timeout:      d.Timeout(schema.TimeoutUpdate),
		DelayTimeout: 10 * time.Second,
		PollInterval: 10 * time.Second,
	})
	if err != nil {
		return fmt.Errorf("error updating instance database port: %s ", err)
	}

	// for prePaid charge mode
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"MODIFYING DATABASE PORT"},
		Target:       []string{"ACTIVE"},
		Refresh:      rdsInstanceStateRefreshFunc(client, instanceID),
		Timeout:      d.Timeout(schema.TimeoutUpdate),
		Delay:        5 * time.Second,
		PollInterval: 3 * time.Second,
	}
	if _, err = stateConf.WaitForState(); err != nil {
		return fmt.Errorf("error waiting for RDS instance (%s) creation completed: %s", instanceID, err)
	}

	return nil
}

func updateRdsInstanceFixedIp(ctx context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient,
	instanceID string) error {
	if !d.HasChange("fixed_ip") {
		return nil
	}

	updateOpts := securities.DataIpOpts{
		NewIp: d.Get("fixed_ip").(string),
	}
	log.Printf("[DEBUG] Update opts of RDS database fixed IP: %+v", updateOpts)

	retryFunc := func() (interface{}, bool, error) {
		res, err := securities.UpdateDataIp(client, instanceID, updateOpts).Extract()
		retry, err := handleMultiOperationsError(err)
		return res, retry, err
	}
	res, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     rdsInstanceStateRefreshFunc(client, instanceID),
		WaitTarget:   []string{"ACTIVE"},
		Timeout:      d.Timeout(schema.TimeoutUpdate),
		DelayTimeout: 1 * time.Second,
		PollInterval: 10 * time.Second,
	})
	if err != nil {
		return fmt.Errorf("error updating instance database fixed IP: %s ", err)
	}
	job := res.(*securities.WorkFlow)

	if err := checkRDSInstanceJobFinish(client, job.WorkflowId, d.Timeout(schema.TimeoutUpdate)); err != nil {
		return fmt.Errorf("error waiting for RDS instance (%s) update fixed IP completed: %s", instanceID, err)
	}

	return nil
}

func updateRdsInstanceSecurityGroup(ctx context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient,
	instanceID string) error {
	if !d.HasChange("security_group_id") {
		return nil
	}

	updateOpts := securities.SecGroupOpts{
		SecurityGroupId: d.Get("security_group_id").(string),
	}
	log.Printf("[DEBUG] Update opts of security group: %+v", updateOpts)

	retryFunc := func() (interface{}, bool, error) {
		_, err := securities.UpdateSecGroup(client, instanceID, updateOpts).Extract()
		retry, err := handleMultiOperationsError(err)
		return nil, retry, err
	}
	_, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     rdsInstanceStateRefreshFunc(client, instanceID),
		WaitTarget:   []string{"ACTIVE"},
		Timeout:      d.Timeout(schema.TimeoutUpdate),
		DelayTimeout: 10 * time.Second,
		PollInterval: 10 * time.Second,
	})
	if err != nil {
		return fmt.Errorf("error updating instance security group: %s ", err)
	}

	return nil
}

func updateRdsInstanceSSLConfig(ctx context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient, instanceID string) error {
	if !d.HasChange("ssl_enable") {
		return nil
	}
	if !isMySQLDatabase(d) {
		return fmt.Errorf("only MySQL database support SSL enable and disable")
	}
	return configRdsInstanceSSL(ctx, d, client, instanceID)
}

func updateRdsParameters(ctx context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient,
	instanceID string) (context.Context, error) {
	values := make(map[string]string)
	if !d.HasChange("parameters") {
		return ctx, nil
	}

	o, n := d.GetChange("parameters")
	os, ns := o.(*schema.Set), n.(*schema.Set)
	change := ns.Difference(os).List()
	if len(change) > 0 {
		for _, v := range change {
			key := v.(map[string]interface{})["name"].(string)
			value := v.(map[string]interface{})["value"].(string)
			values[key] = value
		}

		configOpts := instances.ModifyConfigurationOpts{
			Values: values,
		}
		retryFunc := func() (interface{}, bool, error) {
			_, err := instances.ModifyConfiguration(client, instanceID, configOpts).Extract()
			retry, err := handleMultiOperationsError(err)
			return nil, retry, err
		}
		_, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
			Ctx:          ctx,
			RetryFunc:    retryFunc,
			WaitFunc:     rdsInstanceStateRefreshFunc(client, instanceID),
			WaitTarget:   []string{"ACTIVE"},
			Timeout:      d.Timeout(schema.TimeoutUpdate),
			DelayTimeout: 10 * time.Second,
			PollInterval: 10 * time.Second,
		})
		if err != nil {
			return ctx, fmt.Errorf("error modifying parameters for RDS instance (%s): %s", instanceID, err)
		}
	}

	// Sending parametersChanged to Read to warn users the instance needs a reboot.
	ctx = context.WithValue(ctx, ctxType("parametersChanged"), "true")

	return ctx, nil
}

func updateVolumeAutoExpand(ctx context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient,
	instanceID string) error {
	if !d.HasChanges("volume.0.limit_size", "volume.0.trigger_threshold") {
		return nil
	}

	limitSize := d.Get("volume.0.limit_size").(int)
	if limitSize > 0 {
		if err := enableVolumeAutoExpand(ctx, d, client, instanceID, limitSize); err != nil {
			return err
		}
	} else {
		if err := disableVolumeAutoExpand(ctx, d.Timeout(schema.TimeoutUpdate), client, instanceID); err != nil {
			return err
		}
	}
	return nil
}

func enableVolumeAutoExpand(ctx context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient,
	instanceID string, limitSize int) error {
	opts := instances.EnableAutoExpandOpts{
		InstanceId:       instanceID,
		LimitSize:        limitSize,
		TriggerThreshold: d.Get("volume.0.trigger_threshold").(int),
	}
	retryFunc := func() (interface{}, bool, error) {
		err := instances.EnableAutoExpand(client, opts)
		retry, err := handleMultiOperationsError(err)
		return nil, retry, err
	}
	_, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     rdsInstanceStateRefreshFunc(client, instanceID),
		WaitTarget:   []string{"ACTIVE"},
		Timeout:      d.Timeout(schema.TimeoutUpdate),
		DelayTimeout: 10 * time.Second,
		PollInterval: 10 * time.Second,
	})
	if err != nil {
		return fmt.Errorf("an error occurred while enable automatic expansion of instance storage: %v", err)
	}
	return nil
}

func disableVolumeAutoExpand(ctx context.Context, timeout time.Duration, client *golangsdk.ServiceClient,
	instanceID string) error {
	retryFunc := func() (interface{}, bool, error) {
		err := instances.DisableAutoExpand(client, instanceID)
		retry, err := handleMultiOperationsError(err)
		return nil, retry, err
	}
	_, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     rdsInstanceStateRefreshFunc(client, instanceID),
		WaitTarget:   []string{"ACTIVE"},
		Timeout:      timeout,
		DelayTimeout: 10 * time.Second,
		PollInterval: 10 * time.Second,
	})
	if err != nil {
		return fmt.Errorf("an error occurred while disable automatic expansion of instance storage: %v", err)
	}
	return nil
}

func configRdsInstanceSSL(ctx context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient,
	instanceID string) error {
	sslEnable := d.Get("ssl_enable").(bool)
	updateOpts := securities.SSLOpts{
		SSLEnable: &sslEnable,
	}
	log.Printf("[DEBUG] Update opts of SSL configuration: %+v", updateOpts)

	retryFunc := func() (interface{}, bool, error) {
		err := securities.UpdateSSL(client, instanceID, updateOpts).ExtractErr()
		retry, err := handleMultiOperationsError(err)
		return nil, retry, err
	}
	_, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     rdsInstanceStateRefreshFunc(client, instanceID),
		WaitTarget:   []string{"ACTIVE"},
		Timeout:      d.Timeout(schema.TimeoutUpdate),
		DelayTimeout: 10 * time.Second,
		PollInterval: 10 * time.Second,
	})
	if err != nil {
		return fmt.Errorf("error updating instance SSL configuration: %s ", err)
	}
	return nil
}

func checkRDSInstanceJobFinish(client *golangsdk.ServiceClient, jobID string, timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"Running"},
		Target:       []string{"Completed"},
		Refresh:      rdsInstanceJobRefreshFunc(client, jobID),
		Timeout:      timeout,
		Delay:        20 * time.Second,
		PollInterval: 10 * time.Second,
	}
	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("error waiting for RDS instance (%s) job to be completed: %s ", jobID, err)
	}
	return nil
}

func rdsInstanceJobRefreshFunc(client *golangsdk.ServiceClient, jobID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		jobOpts := instances.RDSJobOpts{
			JobID: jobID,
		}
		jobList, err := instances.GetRDSJob(client, jobOpts).Extract()
		if err != nil {
			return nil, "FOUND ERROR", err
		}

		return jobList.Job, jobList.Job.Status, nil
	}
}

func rdsInstanceStateRefreshFunc(client *golangsdk.ServiceClient, instanceID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		instance, err := GetRdsInstanceByID(client, instanceID)
		if err != nil {
			return nil, "FOUND ERROR", err
		}
		if instance.Id == "" {
			return instance, "DELETED", nil
		}
		if instance.Status == "FAILED" {
			return nil, instance.Status, fmt.Errorf("the instance status is: %s", instance.Status)
		}
		return instance, instance.Status, nil
	}
}

func updateRdsRootPassword(ctx context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient,
	instanceID string) error {
	if !d.HasChange("db.0.password") {
		return nil
	}

	updateOpts := instances.RestRootPasswordOpts{
		DbUserPwd: d.Get("db.0.password").(string),
	}

	retryFunc := func() (interface{}, bool, error) {
		_, err := instances.RestRootPassword(client, instanceID, updateOpts)
		retry, err := handleMultiOperationsError(err)
		return nil, retry, err
	}
	_, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     rdsInstanceStateRefreshFunc(client, instanceID),
		WaitTarget:   []string{"ACTIVE"},
		Timeout:      d.Timeout(schema.TimeoutUpdate),
		DelayTimeout: 10 * time.Second,
		PollInterval: 10 * time.Second,
	})
	if err != nil {
		return fmt.Errorf("error resetting the root password: %s", err)
	}
	return nil
}

func updateRdsInstanceMaintainWindow(d *schema.ResourceData, client *golangsdk.ServiceClient, instanceID string) error {
	if !d.HasChanges("maintain_begin", "maintain_end") {
		return nil
	}

	modifyMaintainWindowOpts := instances.ModifyMaintainWindowOpts{
		StartTime: d.Get("maintain_begin").(string),
		EndTime:   d.Get("maintain_end").(string),
	}

	log.Printf("[DEBUG] Modify RDS instance maintain window opts: %+v", modifyMaintainWindowOpts)
	r := instances.ModifyMaintainWindow(client, modifyMaintainWindowOpts, instanceID)
	if r.Err != nil {
		return fmt.Errorf("error modify RDS instance (%s) maintain window: %s", instanceID, r.Err)
	}
	return nil
}

func updateRdsInstanceReplicationMode(ctx context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient,
	instanceID string) error {
	if !d.HasChanges("ha_replication_mode") {
		return nil
	}

	modifyReplicationModeOpts := instances.ModifyReplicationModeOpts{
		Mode: d.Get("ha_replication_mode").(string),
	}

	log.Printf("[DEBUG] Modify RDS instance replication mode opts: %+v", modifyReplicationModeOpts)
	retryFunc := func() (interface{}, bool, error) {
		res, err := instances.ModifyReplicationMode(client, modifyReplicationModeOpts, instanceID).Extract()
		retry, err := handleMultiOperationsError(err)
		return res, retry, err
	}
	res, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     rdsInstanceStateRefreshFunc(client, instanceID),
		WaitTarget:   []string{"ACTIVE"},
		Timeout:      d.Timeout(schema.TimeoutUpdate),
		DelayTimeout: 1 * time.Second,
		PollInterval: 10 * time.Second,
	})
	if err != nil {
		return fmt.Errorf("error modify RDS instance (%s) replication mode: %s", instanceID, err)
	}
	job := res.(*instances.ReplicationMode)

	if err = checkRDSInstanceJobFinish(client, job.WorkflowId, d.Timeout(schema.TimeoutUpdate)); err != nil {
		return fmt.Errorf("error waiting for RDS instance (%s) update replication mode completed: %s", instanceID, err)
	}
	return nil
}

func updateRdsInstanceSwitchStrategy(ctx context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient,
	instanceID string) error {
	if !d.HasChanges("switch_strategy") {
		return nil
	}

	modifySwitchStrategyOpts := instances.ModifySwitchStrategyOpts{
		RepairStrategy: d.Get("switch_strategy").(string),
	}

	log.Printf("[DEBUG] Modify RDS instance switch strategy opts: %+v", modifySwitchStrategyOpts)
	retryFunc := func() (interface{}, bool, error) {
		res := instances.ModifySwitchStrategy(client, modifySwitchStrategyOpts, instanceID)
		retry, err := handleMultiOperationsError(res.Err)
		return res, retry, err
	}
	_, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     rdsInstanceStateRefreshFunc(client, instanceID),
		WaitTarget:   []string{"ACTIVE"},
		Timeout:      d.Timeout(schema.TimeoutUpdate),
		DelayTimeout: 1 * time.Second,
		PollInterval: 10 * time.Second,
	})
	if err != nil {
		return fmt.Errorf("error modify RDS instance (%s) switch strategy: %s", instanceID, err)
	}
	return nil
}

func updateRdsInstanceCollation(ctx context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient,
	instanceID string) error {
	if !d.HasChanges("collation") {
		return nil
	}

	modifyCollationOpts := instances.ModifyCollationOpts{
		Collation: d.Get("collation").(string),
	}

	log.Printf("[DEBUG] Modify RDS instance collation opts: %+v", modifyCollationOpts)
	retryFunc := func() (interface{}, bool, error) {
		res, err := instances.ModifyCollation(client, modifyCollationOpts, instanceID).Extract()
		retry, err := handleMultiOperationsError(err)
		return res, retry, err
	}
	res, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     rdsInstanceStateRefreshFunc(client, instanceID),
		WaitTarget:   []string{"ACTIVE"},
		Timeout:      d.Timeout(schema.TimeoutUpdate),
		DelayTimeout: 1 * time.Second,
		PollInterval: 10 * time.Second,
	})
	if err != nil {
		return fmt.Errorf("error modify RDS instance (%s) collation: %s", instanceID, err)
	}
	job := res.(*instances.Collation)

	if err = checkRDSInstanceJobFinish(client, job.JobId, d.Timeout(schema.TimeoutUpdate)); err != nil {
		return fmt.Errorf("error waiting for RDS instance (%s) update collation completed: %s", instanceID, err)
	}
	return nil
}

func parameterToHash(v interface{}) int {
	m := v.(map[string]interface{})
	return hashcode.String(m["name"].(string) + m["value"].(string))
}
