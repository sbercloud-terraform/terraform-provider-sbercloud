package dds

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/common/tags"
	"github.com/chnsz/golangsdk/openstack/dds/v3/instances"
	"github.com/chnsz/golangsdk/openstack/dds/v3/jobs"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func ResourceDdsInstanceV3() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDdsInstanceV3Create,
		ReadContext:   resourceDdsInstanceV3Read,
		UpdateContext: resourceDdsInstanceV3Update,
		DeleteContext: resourceDdsInstanceV3Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Update: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(60 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"datastore": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"version": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"storage_engine": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
					},
				},
			},
			"availability_zone": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
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
			"port": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"password": {
				Type:      schema.TypeString,
				Sensitive: true,
				Required:  true,
			},
			"disk_encryption_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"mode": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"configuration": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"id": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
					},
				},
			},
			"flavor": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"num": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"storage": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"size": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"spec_code": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
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
							Type:     schema.TypeInt,
							Required: true,
						},
					},
				},
			},
			"ssl": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"charging_mode": common.SchemaChargingMode(nil),
			"period_unit":   common.SchemaPeriodUnit(nil),
			"period":        common.SchemaPeriod(nil),
			"auto_renew":    common.SchemaAutoRenew(nil),
			"auto_pay":      common.SchemaAutoPay(nil),
			"tags":          common.TagsSchema(),
			"db_username": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"nodes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"role": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"private_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"public_ip": {
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
		},
	}
}

func resourceDdsDataStore(d *schema.ResourceData) instances.DataStore {
	var dataStore instances.DataStore
	datastoreRaw := d.Get("datastore").([]interface{})
	log.Printf("[DEBUG] datastoreRaw: %+v", datastoreRaw)
	if len(datastoreRaw) == 1 {
		dataStore.Type = datastoreRaw[0].(map[string]interface{})["type"].(string)
		dataStore.Version = datastoreRaw[0].(map[string]interface{})["version"].(string)
		dataStore.StorageEngine = datastoreRaw[0].(map[string]interface{})["storage_engine"].(string)
	}
	log.Printf("[DEBUG] datastore: %+v", dataStore)
	return dataStore
}

func resourceDdsConfiguration(d *schema.ResourceData) []instances.Configuration {
	var configurations []instances.Configuration
	configurationRaw := d.Get("configuration").([]interface{})
	log.Printf("[DEBUG] configurationRaw: %+v", configurationRaw)
	for i := range configurationRaw {
		configuration := configurationRaw[i].(map[string]interface{})
		flavorReq := instances.Configuration{
			Type: configuration["type"].(string),
			Id:   configuration["id"].(string),
		}
		configurations = append(configurations, flavorReq)
	}
	log.Printf("[DEBUG] configurations: %+v", configurations)
	return configurations
}

func resourceDdsFlavors(d *schema.ResourceData) []instances.Flavor {
	var flavors []instances.Flavor
	flavorRaw := d.Get("flavor").([]interface{})
	log.Printf("[DEBUG] flavorRaw: %+v", flavorRaw)
	for i := range flavorRaw {
		flavor := flavorRaw[i].(map[string]interface{})
		flavorReq := instances.Flavor{
			Type:     flavor["type"].(string),
			Num:      flavor["num"].(int),
			Storage:  flavor["storage"].(string),
			Size:     flavor["size"].(int),
			SpecCode: flavor["spec_code"].(string),
		}
		flavors = append(flavors, flavorReq)
	}
	log.Printf("[DEBUG] flavors: %+v", flavors)
	return flavors
}

func resourceDdsBackupStrategy(d *schema.ResourceData) instances.BackupStrategy {
	var backupStrategy instances.BackupStrategy
	backupStrategyRaw := d.Get("backup_strategy").([]interface{})
	log.Printf("[DEBUG] backupStrategyRaw: %+v", backupStrategyRaw)
	startTime := "00:00-01:00"
	keepDays := 7
	if len(backupStrategyRaw) == 1 {
		startTime = backupStrategyRaw[0].(map[string]interface{})["start_time"].(string)
		keepDays = backupStrategyRaw[0].(map[string]interface{})["keep_days"].(int)
	}
	backupStrategy.StartTime = startTime
	backupStrategy.KeepDays = &keepDays
	log.Printf("[DEBUG] backupStrategy: %+v", backupStrategy)
	return backupStrategy
}

func ddsInstanceStateRefreshFunc(client *golangsdk.ServiceClient, instanceID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		opts := instances.ListInstanceOpts{
			Id: instanceID,
		}
		allPages, err := instances.List(client, &opts).AllPages()
		if err != nil {
			return nil, "", err
		}
		instancesList, err := instances.ExtractInstances(allPages)
		if err != nil {
			return nil, "", err
		}

		if instancesList.TotalCount == 0 {
			var instance instances.InstanceResponse
			return instance, "deleted", nil
		}
		insts := instancesList.Instances

		status := insts[0].Status
		// wait for updating
		if status == "normal" && len(insts[0].Actions) > 0 {
			status = "updating"
		}
		return insts[0], status, nil
	}
}

func buildChargeInfoParams(d *schema.ResourceData) instances.ChargeInfo {
	chargeInfo := instances.ChargeInfo{
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
	return chargeInfo
}

func resourceDdsInstanceV3Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	client, err := conf.DdsV3Client(conf.GetRegion(d))
	if err != nil {
		return diag.Errorf("Error creating DDS client: %s ", err)
	}

	createOpts := instances.CreateOpts{
		Name:                d.Get("name").(string),
		DataStore:           resourceDdsDataStore(d),
		Region:              conf.GetRegion(d),
		AvailabilityZone:    d.Get("availability_zone").(string),
		VpcId:               d.Get("vpc_id").(string),
		SubnetId:            d.Get("subnet_id").(string),
		SecurityGroupId:     d.Get("security_group_id").(string),
		DiskEncryptionId:    d.Get("disk_encryption_id").(string),
		Mode:                d.Get("mode").(string),
		Configuration:       resourceDdsConfiguration(d),
		Flavor:              resourceDdsFlavors(d),
		BackupStrategy:      resourceDdsBackupStrategy(d),
		EnterpriseProjectID: conf.GetEnterpriseProjectID(d),
	}
	if d.Get("ssl").(bool) {
		createOpts.Ssl = "1"
	} else {
		createOpts.Ssl = "0"
	}
	if d.Get("charging_mode").(string) == "prePaid" {
		chargeInfo := buildChargeInfoParams(d)
		createOpts.ChargeInfo = &chargeInfo
	}
	log.Printf("[DEBUG] Create Options: %#v", createOpts)
	// Add password here so it wouldn't go in the above log entry
	createOpts.Password = d.Get("password").(string)

	if val, ok := d.GetOk("port"); ok {
		createOpts.Port = strconv.Itoa(val.(int))
	}

	instance, err := instances.Create(client, createOpts).Extract()
	if err != nil {
		return diag.Errorf("Error getting instance from result: %s ", err)
	}
	log.Printf("[DEBUG] Create : instance %s: %#v", instance.Id, instance)

	if instance.OrderId != "" {
		bssClient, err := conf.BssV2Client(conf.GetRegion(d))
		if err != nil {
			return diag.Errorf("error creating BSS v2 client: %s", err)
		}
		err = common.WaitOrderComplete(ctx, bssClient, instance.OrderId, d.Timeout(schema.TimeoutCreate))
		if err != nil {
			return diag.FromErr(err)
		}
		resourceId, err := common.WaitOrderResourceComplete(ctx, bssClient, instance.OrderId,
			d.Timeout(schema.TimeoutCreate))
		if err != nil {
			return diag.FromErr(err)
		}

		d.SetId(resourceId)
	} else {
		d.SetId(instance.Id)
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"creating", "updating"},
		Target:     []string{"normal"},
		Refresh:    ddsInstanceStateRefreshFunc(client, instance.Id),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      120 * time.Second,
		MinTimeout: 20 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf(
			"Error waiting for instance (%s) to become ready: %s ",
			instance.Id, err)
	}

	// set tags
	tagRaw := d.Get("tags").(map[string]interface{})
	if len(tagRaw) > 0 {
		taglist := utils.ExpandResourceTags(tagRaw)
		if tagErr := tags.Create(client, "instances", instance.Id, taglist).ExtractErr(); tagErr != nil {
			return diag.Errorf("Error setting tags of DDS instance %s: %s", instance.Id, tagErr)
		}
	}

	return resourceDdsInstanceV3Read(ctx, d, meta)
}

func resourceDdsInstanceV3Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	client, err := conf.DdsV3Client(conf.GetRegion(d))
	if err != nil {
		return diag.Errorf("Error creating DDS client: %s", err)
	}

	instanceID := d.Id()
	opts := instances.ListInstanceOpts{
		Id: instanceID,
	}
	allPages, err := instances.List(client, &opts).AllPages()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving DdsInstance")
	}
	instanceList, err := instances.ExtractInstances(allPages)
	if err != nil {
		return diag.Errorf("Error extracting DDS instance: %s", err)
	}
	if instanceList.TotalCount == 0 {
		log.Printf("[WARN] DDS instance (%s) was not found", instanceID)
		d.SetId("")
		return nil
	}
	insts := instanceList.Instances
	instanceObj := insts[0]

	log.Printf("[DEBUG] Retrieved instance %s: %#v", instanceID, instanceObj)

	mErr := multierror.Append(
		d.Set("region", instanceObj.Region),
		d.Set("name", instanceObj.Name),
		d.Set("vpc_id", instanceObj.VpcId),
		d.Set("subnet_id", instanceObj.SubnetId),
		d.Set("security_group_id", instanceObj.SecurityGroupId),
		d.Set("disk_encryption_id", instanceObj.DiskEncryptionId),
		d.Set("mode", instanceObj.Mode),
		d.Set("db_username", instanceObj.DbUserName),
		d.Set("status", instanceObj.Status),
		d.Set("enterprise_project_id", instanceObj.EnterpriseProjectID),
		d.Set("nodes", flattenDdsInstanceV3Nodes(instanceObj)),
	)

	port, err := strconv.Atoi(instanceObj.Port)
	if err != nil {
		log.Printf("[WARNING] Port %s invalid, Type conversion error: %s", instanceObj.Port, err)
	}
	mErr = multierror.Append(mErr, d.Set("port", port))

	sslEnable := true
	if instanceObj.Ssl == 0 {
		sslEnable = false
	}
	mErr = multierror.Append(mErr, d.Set("ssl", sslEnable))

	datastoreList := make([]map[string]interface{}, 0, 1)
	datastore := map[string]interface{}{
		"type":           instanceObj.DataStore.Type,
		"version":        instanceObj.DataStore.Version,
		"storage_engine": instanceObj.Engine,
	}
	datastoreList = append(datastoreList, datastore)
	mErr = multierror.Append(mErr, d.Set("datastore", datastoreList))

	backupStrategyList := make([]map[string]interface{}, 0, 1)
	backupStrategy := map[string]interface{}{
		"start_time": instanceObj.BackupStrategy.StartTime,
		"keep_days":  instanceObj.BackupStrategy.KeepDays,
	}
	backupStrategyList = append(backupStrategyList, backupStrategy)
	mErr = multierror.Append(mErr, d.Set("backup_strategy", backupStrategyList))

	// save tags
	if resourceTags, err := tags.Get(client, "instances", d.Id()).Extract(); err == nil {
		tagmap := utils.TagsToMap(resourceTags.Tags)
		mErr = multierror.Append(mErr, d.Set("tags", tagmap))
	} else {
		log.Printf("[WARN] Error fetching tags of DDS instance (%s): %s", d.Id(), err)
	}

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("Error setting dds instance fields: %s", err)
	}

	return nil
}

func JobStateRefreshFunc(client *golangsdk.ServiceClient, jobId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := jobs.Get(client, jobId)
		if err != nil {
			return nil, "", err
		}

		return resp, resp.Status, nil
	}
}

func waitForInstanceReady(ctx context.Context, client *golangsdk.ServiceClient, instanceId string, timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"updating"},
		Target:     []string{"normal"},
		Refresh:    ddsInstanceStateRefreshFunc(client, instanceId),
		Timeout:    timeout,
		Delay:      15 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf(
			"Error waiting for instance (%s) to become ready: %s ",
			instanceId, err)
	}

	return nil
}

func resourceDdsInstanceV3Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	client, err := conf.DdsV3Client(conf.GetRegion(d))
	if err != nil {
		return diag.Errorf("Error creating DDS client: %s ", err)
	}

	var opts []instances.UpdateOpt
	if d.HasChange("name") {
		opt := instances.UpdateOpt{
			Param:  "new_instance_name",
			Value:  d.Get("name").(string),
			Action: "modify-name",
			Method: "put",
		}
		opts = append(opts, opt)
	}

	if d.HasChange("password") {
		opt := instances.UpdateOpt{
			Param:  "user_pwd",
			Value:  d.Get("password").(string),
			Action: "reset-password",
			Method: "put",
		}
		opts = append(opts, opt)
	}

	if d.HasChange("security_group_id") {
		opt := instances.UpdateOpt{
			Param:  "security_group_id",
			Value:  d.Get("security_group_id").(string),
			Action: "modify-security-group",
			Method: "post",
		}
		opts = append(opts, opt)
	}

	if d.HasChange("backup_strategy") {
		backupStrategy := resourceDdsBackupStrategy(d)
		backupStrategy.Period = "1,2,3,4,5,6,7"
		opt := instances.UpdateOpt{
			Param:  "backup_policy",
			Value:  backupStrategy,
			Action: "backups/policy",
			Method: "put",
		}
		opts = append(opts, opt)
	}

	if d.HasChange("ssl") {
		opt := instances.UpdateOpt{
			Param:  "ssl_option",
			Action: "switch-ssl",
			Method: "post",
		}
		if d.Get("ssl").(bool) {
			opt.Value = "1"
		} else {
			opt.Value = "0"
		}
		opts = append(opts, opt)
	}

	if len(opts) > 0 {
		retryFunc := func() (interface{}, bool, error) {
			resp, err := instances.Update(client, d.Id(), opts).Extract()
			retry, err := handleMultiOperationsError(err)
			return resp, retry, err
		}
		r, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
			Ctx:          ctx,
			RetryFunc:    retryFunc,
			WaitFunc:     ddsInstanceStateRefreshFunc(client, d.Id()),
			WaitTarget:   []string{"normal"},
			Timeout:      d.Timeout(schema.TimeoutUpdate),
			DelayTimeout: 1 * time.Second,
			PollInterval: 10 * time.Second,
		})

		if err != nil {
			return diag.Errorf("Error updating instance from result: %s ", err)
		}
		resp := r.(*instances.UpdateResp)
		if resp.OrderId != "" {
			bssClient, err := conf.BssV2Client(conf.GetRegion(d))
			if err != nil {
				return diag.Errorf("error creating BSS v2 client: %s", err)
			}
			err = common.WaitOrderComplete(ctx, bssClient, resp.OrderId, d.Timeout(schema.TimeoutCreate))
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	if d.HasChange("port") {
		retryFunc := func() (interface{}, bool, error) {
			resp, err := instances.UpdatePort(client, d.Id(), d.Get("port").(int))
			retry, err := handleMultiOperationsError(err)
			return resp, retry, err
		}
		r, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
			Ctx:          ctx,
			RetryFunc:    retryFunc,
			WaitFunc:     ddsInstanceStateRefreshFunc(client, d.Id()),
			WaitTarget:   []string{"normal"},
			Timeout:      d.Timeout(schema.TimeoutUpdate),
			DelayTimeout: 1 * time.Second,
			PollInterval: 10 * time.Second,
		})
		if err != nil {
			return diag.Errorf("error updating database access port: %s", err)
		}
		resp := r.(*instances.PortUpdateResp)
		stateConf := &resource.StateChangeConf{
			Pending:      []string{"Running"},
			Target:       []string{"Completed"},
			Refresh:      JobStateRefreshFunc(client, resp.JobId),
			Timeout:      d.Timeout(schema.TimeoutUpdate),
			PollInterval: 10 * time.Second,
		}

		_, err = stateConf.WaitForStateContext(ctx)
		if err != nil {
			return diag.Errorf("error waiting for the job (%s) completed: %s ", resp.JobId, err)
		}
	}

	if d.HasChange("tags") {
		tagErr := utils.UpdateResourceTags(client, d, "instances", d.Id())
		if tagErr != nil {
			return diag.Errorf("Error updating tags of DDS instance:%s, err:%s", d.Id(), tagErr)
		}
	}

	// update flavor
	if d.HasChange("flavor") {
		for i := range d.Get("flavor").([]interface{}) {
			numIndex := fmt.Sprintf("flavor.%d.num", i)
			volumeSizeIndex := fmt.Sprintf("flavor.%d.size", i)
			specCodeIndex := fmt.Sprintf("flavor.%d.spec_code", i)

			// The update operation of the volume size must ahead of the update operation of the number. Because the
			// size and number are updated at the same time and the number is increased, and the request will fail.
			// For example, when the number is increased from 2 to 3, and the size of all nodes is increased from 20 to
			// 30, the newly added node will prompt that the storage update failed and cannot be updated from 30 to 30.
			if d.HasChange(volumeSizeIndex) {
				err := flavorSizeUpdate(ctx, conf, client, d, i)
				if err != nil {
					return diag.FromErr(err)
				}
			}
			if d.HasChange(numIndex) {
				err := flavorNumUpdate(ctx, conf, client, d, i)
				if err != nil {
					return diag.FromErr(err)
				}
			}
			if d.HasChange(specCodeIndex) {
				err := flavorSpecCodeUpdate(ctx, conf, client, d, i)
				if err != nil {
					return diag.FromErr(err)
				}
			}
		}
	}

	return resourceDdsInstanceV3Read(ctx, d, meta)
}

func resourceDdsInstanceV3Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	client, err := conf.DdsV3Client(conf.GetRegion(d))
	if err != nil {
		return diag.Errorf("Error creating DDS client: %s ", err)
	}

	instanceId := d.Id()
	// for prePaid mode, we should unsubscribe the resource
	if d.Get("charging_mode").(string) == "prePaid" {
		retryFunc := func() (interface{}, bool, error) {
			err = common.UnsubscribePrePaidResource(d, conf, []string{instanceId})
			retry, err := handleDeletionError(err)
			return nil, retry, err
		}
		_, err = common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
			Ctx:          ctx,
			RetryFunc:    retryFunc,
			WaitFunc:     ddsInstanceStateRefreshFunc(client, d.Id()),
			WaitTarget:   []string{"normal"},
			Timeout:      d.Timeout(schema.TimeoutDelete),
			DelayTimeout: 1 * time.Second,
			PollInterval: 10 * time.Second,
		})
		if err != nil {
			return diag.Errorf("error unsubscribing DDS instance : %s", err)
		}
	} else {
		retryFunc := func() (interface{}, bool, error) {
			result := instances.Delete(client, instanceId)
			retry, err := handleDeletionError(result.Err)
			return nil, retry, err
		}
		_, err = common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
			Ctx:          ctx,
			RetryFunc:    retryFunc,
			WaitFunc:     ddsInstanceStateRefreshFunc(client, d.Id()),
			WaitTarget:   []string{"normal"},
			Timeout:      d.Timeout(schema.TimeoutDelete),
			DelayTimeout: 1 * time.Second,
			PollInterval: 10 * time.Second,
		})
		if err != nil {
			return diag.FromErr(err)
		}
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"normal", "abnormal", "frozen", "createfail", "enlargefail", "data_disk_full"},
		Target:     []string{"deleted"},
		Refresh:    ddsInstanceStateRefreshFunc(client, instanceId),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      15 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf(
			"Error waiting for instance (%s) to be deleted: %s ",
			instanceId, err)
	}
	log.Printf("[DEBUG] Successfully deleted instance %s", instanceId)
	return nil
}

func flattenDdsInstanceV3Nodes(dds instances.InstanceResponse) interface{} {
	nodesList := make([]map[string]interface{}, 0)
	for _, group := range dds.Groups {
		groupType := group.Type
		for _, Node := range group.Nodes {
			node := map[string]interface{}{
				"type":       groupType,
				"id":         Node.Id,
				"name":       Node.Name,
				"role":       Node.Role,
				"status":     Node.Status,
				"private_ip": Node.PrivateIP,
				"public_ip":  Node.PublicIP,
			}
			nodesList = append(nodesList, node)
		}
	}
	return nodesList
}

func getDdsInstanceV3ShardGroupID(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]string, error) {
	groupIDs := make([]string, 0)

	instanceID := d.Id()
	opts := instances.ListInstanceOpts{
		Id: instanceID,
	}
	allPages, err := instances.List(client, &opts).AllPages()
	if err != nil {
		return groupIDs, fmt.Errorf("Error fetching DDS instance: %s", err)
	}
	instanceList, err := instances.ExtractInstances(allPages)
	if err != nil {
		return groupIDs, fmt.Errorf("Error extracting DDS instance: %s", err)
	}
	if instanceList.TotalCount == 0 {
		log.Printf("[WARN] DDS instance (%s) was not found", instanceID)
		return groupIDs, nil
	}
	insts := instanceList.Instances
	instanceObj := insts[0]

	log.Printf("[DEBUG] Retrieved instance %s: %#v", instanceID, instanceObj)

	for _, group := range instanceObj.Groups {
		if group.Type == "shard" {
			groupIDs = append(groupIDs, group.Id)
		}
	}

	return groupIDs, nil
}

func getDdsInstanceV3MongosNodeID(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]string, error) {
	nodeIDs := make([]string, 0)

	instanceID := d.Id()
	opts := instances.ListInstanceOpts{
		Id: instanceID,
	}
	allPages, err := instances.List(client, &opts).AllPages()
	if err != nil {
		return nodeIDs, fmt.Errorf("Error fetching DDS instance: %s", err)
	}
	instanceList, err := instances.ExtractInstances(allPages)
	if err != nil {
		return nodeIDs, fmt.Errorf("Error extracting DDS instance: %s", err)
	}
	if instanceList.TotalCount == 0 {
		log.Printf("[WARN] DDS instance (%s) was not found", instanceID)
		return nodeIDs, nil
	}
	insts := instanceList.Instances
	instanceObj := insts[0]

	log.Printf("[DEBUG] Retrieved instance %s: %#v", instanceID, instanceObj)

	for _, group := range instanceObj.Groups {
		if group.Type == "mongos" {
			for _, node := range group.Nodes {
				nodeIDs = append(nodeIDs, node.Id)
			}
		}
	}

	return nodeIDs, nil
}

func flavorUpdate(ctx context.Context, conf *config.Config, client *golangsdk.ServiceClient, d *schema.ResourceData,
	opts []instances.UpdateOpt) error {
	retryFunc := func() (interface{}, bool, error) {
		resp, err := instances.Update(client, d.Id(), opts).Extract()
		retry, err := handleMultiOperationsError(err)
		return resp, retry, err
	}
	r, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     ddsInstanceStateRefreshFunc(client, d.Id()),
		WaitTarget:   []string{"normal"},
		Timeout:      d.Timeout(schema.TimeoutUpdate),
		DelayTimeout: 1 * time.Second,
		PollInterval: 10 * time.Second,
	})
	if err != nil {
		return fmt.Errorf("Error updating instance from result: %s ", err)
	}
	resp := r.(*instances.UpdateResp)
	if resp.OrderId != "" {
		bssClient, err := conf.BssV2Client(conf.GetRegion(d))
		if err != nil {
			return fmt.Errorf("error creating BSS v2 client: %s", err)
		}
		err = common.WaitOrderComplete(ctx, bssClient, resp.OrderId, d.Timeout(schema.TimeoutCreate))
		if err != nil {
			return err
		}
	}

	err = waitForInstanceReady(ctx, client, d.Id(), d.Timeout(schema.TimeoutUpdate))
	if err != nil {
		return err
	}

	return nil
}

func flavorNumUpdate(ctx context.Context, conf *config.Config, client *golangsdk.ServiceClient, d *schema.ResourceData, i int) error {
	groupTypeIndex := fmt.Sprintf("flavor.%d.type", i)
	groupType := d.Get(groupTypeIndex).(string)
	if groupType != "mongos" && groupType != "shard" {
		return fmt.Errorf("Error updating instance: %s does not support adding nodes", groupType)
	}
	specCodeIndex := fmt.Sprintf("flavor.%d.spec_code", i)
	volumeSizeIndex := fmt.Sprintf("flavor.%d.size", i)
	volumeSize := d.Get(volumeSizeIndex).(int)
	numIndex := fmt.Sprintf("flavor.%d.num", i)
	oldNumRaw, newNumRaw := d.GetChange(numIndex)
	oldNum := oldNumRaw.(int)
	newNum := newNumRaw.(int)
	if newNum < oldNum {
		return fmt.Errorf("Error updating instance: the new num(%d) must be greater than the old num(%d)", newNum, oldNum)
	}

	var numUpdateOpts []instances.UpdateOpt

	if groupType == "mongos" {
		updateNodeNumOpts := instances.UpdateNodeNumOpts{
			Type:     groupType,
			SpecCode: d.Get(specCodeIndex).(string),
			Num:      newNum - oldNum,
		}
		if d.Get("charging_mode").(string) == "prePaid" && d.Get("auto_pay").(string) != "false" {
			updateNodeNumOpts.IsAutoPay = true
		}
		opt := instances.UpdateOpt{
			Param:  "",
			Value:  updateNodeNumOpts,
			Action: "enlarge",
			Method: "post",
		}

		numUpdateOpts = append(numUpdateOpts, opt)
	} else {
		volume := instances.VolumeOpts{
			Size: &volumeSize,
		}
		updateNodeNumOpts := instances.UpdateNodeNumOpts{
			Type:     groupType,
			SpecCode: d.Get(specCodeIndex).(string),
			Num:      newNum - oldNum,
			Volume:   &volume,
		}
		if d.Get("charging_mode").(string) == "prePaid" && d.Get("auto_pay").(string) != "false" {
			updateNodeNumOpts.IsAutoPay = true
		}
		opt := instances.UpdateOpt{
			Param:  "",
			Value:  updateNodeNumOpts,
			Action: "enlarge",
			Method: "post",
		}
		numUpdateOpts = append(numUpdateOpts, opt)
	}
	err := flavorUpdate(ctx, conf, client, d, numUpdateOpts)
	if err != nil {
		return err
	}
	return nil
}

func flavorSizeUpdate(ctx context.Context, conf *config.Config, client *golangsdk.ServiceClient, d *schema.ResourceData, i int) error {
	volumeSizeIndex := fmt.Sprintf("flavor.%d.size", i)
	oldSizeRaw, newSizeRaw := d.GetChange(volumeSizeIndex)
	oldSize := oldSizeRaw.(int)
	newSize := newSizeRaw.(int)
	if newSize < oldSize {
		return fmt.Errorf("Error updating instance: the new size(%d) must be greater than the old size(%d)", newSize, oldSize)
	}
	groupTypeIndex := fmt.Sprintf("flavor.%d.type", i)
	groupType := d.Get(groupTypeIndex).(string)
	if groupType != "replica" && groupType != "single" && groupType != "shard" {
		return fmt.Errorf("Error updating instance: %s does not support scaling up storage space", groupType)
	}

	if groupType == "shard" {
		groupIDs, err := getDdsInstanceV3ShardGroupID(client, d)
		if err != nil {
			return err
		}

		for _, groupID := range groupIDs {
			var sizeUpdateOpts []instances.UpdateOpt
			updateVolumeOpts := instances.UpdateVolumeOpts{
				Volume: instances.VolumeOpts{
					GroupID: groupID,
					Size:    &newSize,
				},
			}
			if d.Get("charging_mode").(string) == "prePaid" && d.Get("auto_pay").(string) != "false" {
				updateVolumeOpts.IsAutoPay = true
			}
			opt := instances.UpdateOpt{
				Param:  "",
				Value:  updateVolumeOpts,
				Action: "enlarge-volume",
				Method: "post",
			}
			sizeUpdateOpts = append(sizeUpdateOpts, opt)
			err := flavorUpdate(ctx, conf, client, d, sizeUpdateOpts)
			if err != nil {
				return err
			}
		}
	} else {
		var sizeUpdateOpts []instances.UpdateOpt
		updateVolumeOpts := instances.UpdateVolumeOpts{
			Volume: instances.VolumeOpts{
				Size: &newSize,
			},
		}
		if d.Get("charging_mode").(string) == "prePaid" && d.Get("auto_pay").(string) != "false" {
			updateVolumeOpts.IsAutoPay = true
		}
		opt := instances.UpdateOpt{
			Param:  "volume",
			Value:  updateVolumeOpts,
			Action: "enlarge-volume",
			Method: "post",
		}
		sizeUpdateOpts = append(sizeUpdateOpts, opt)
		err := flavorUpdate(ctx, conf, client, d, sizeUpdateOpts)
		if err != nil {
			return err
		}
	}
	return nil
}

func flavorSpecCodeUpdate(ctx context.Context, conf *config.Config, client *golangsdk.ServiceClient, d *schema.ResourceData, i int) error {
	specCodeIndex := fmt.Sprintf("flavor.%d.spec_code", i)
	groupTypeIndex := fmt.Sprintf("flavor.%d.type", i)
	groupType := d.Get(groupTypeIndex).(string)
	if groupType == "config" {
		return fmt.Errorf("Error updating instance: %s does not support updating spec_code", groupType)
	}
	switch groupType {
	case "mongos":
		nodeIDs, err := getDdsInstanceV3MongosNodeID(client, d)
		if err != nil {
			return err
		}
		for _, ID := range nodeIDs {
			var specUpdateOpts []instances.UpdateOpt
			updateSpecOpts := instances.UpdateSpecOpts{
				Resize: instances.SpecOpts{
					TargetType:     "mongos",
					TargetID:       ID,
					TargetSpecCode: d.Get(specCodeIndex).(string),
				},
			}
			if d.Get("charging_mode").(string) == "prePaid" && d.Get("auto_pay").(string) != "false" {
				updateSpecOpts.IsAutoPay = true
			}
			opt := instances.UpdateOpt{
				Param:  "",
				Value:  updateSpecOpts,
				Action: "resize",
				Method: "post",
			}
			specUpdateOpts = append(specUpdateOpts, opt)
			err := flavorUpdate(ctx, conf, client, d, specUpdateOpts)
			if err != nil {
				return err
			}
		}
	case "shard":
		groupIDs, err := getDdsInstanceV3ShardGroupID(client, d)
		if err != nil {
			return err
		}

		for _, ID := range groupIDs {
			var specUpdateOpts []instances.UpdateOpt
			updateSpecOpts := instances.UpdateSpecOpts{
				Resize: instances.SpecOpts{
					TargetType:     "shard",
					TargetID:       ID,
					TargetSpecCode: d.Get(specCodeIndex).(string),
				},
			}
			if d.Get("charging_mode").(string) == "prePaid" && d.Get("auto_pay").(string) != "false" {
				updateSpecOpts.IsAutoPay = true
			}
			opt := instances.UpdateOpt{
				Param:  "resize",
				Value:  updateSpecOpts,
				Action: "resize",
				Method: "post",
			}
			specUpdateOpts = append(specUpdateOpts, opt)
			err := flavorUpdate(ctx, conf, client, d, specUpdateOpts)
			if err != nil {
				return err
			}
		}
	default:
		var specUpdateOpts []instances.UpdateOpt
		updateSpecOpts := instances.UpdateSpecOpts{
			Resize: instances.SpecOpts{
				TargetID:       d.Id(),
				TargetSpecCode: d.Get(specCodeIndex).(string),
			},
		}
		if d.Get("charging_mode").(string) == "prePaid" && d.Get("auto_pay").(string) != "false" {
			updateSpecOpts.IsAutoPay = true
		}
		opt := instances.UpdateOpt{
			Param:  "resize",
			Value:  updateSpecOpts,
			Action: "resize",
			Method: "post",
		}
		specUpdateOpts = append(specUpdateOpts, opt)
		err := flavorUpdate(ctx, conf, client, d, specUpdateOpts)
		if err != nil {
			return err
		}
	}

	return nil
}
