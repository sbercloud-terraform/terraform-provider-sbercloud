package huaweicloud

import (
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/openstack/bss/v2/orders"
	"github.com/huaweicloud/golangsdk/openstack/taurusdb/v3/backups"
	"github.com/huaweicloud/golangsdk/openstack/taurusdb/v3/instances"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"
)

func resourceGaussDBInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceGaussDBInstanceCreate,
		Update: resourceGaussDBInstanceUpdate,
		Read:   resourceGaussDBInstanceRead,
		Delete: resourceGaussDBInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		CustomizeDiff: func(d *schema.ResourceDiff, v interface{}) error {
			if d.HasChange("proxy_node_num") {
				d.SetNewComputed("proxy_address")
				d.SetNewComputed("proxy_port")
			}
			return nil
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Update: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"flavor": {
				Type:     schema.TypeString,
				Required: true,
			},
			"password": {
				Type:      schema.TypeString,
				Sensitive: true,
				Required:  true,
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
				Optional: true,
				ForceNew: true,
			},
			"configuration_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"table_name_case_sensitivity": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"read_replicas": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  1,
			},
			"volume_size": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"time_zone": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "UTC+08:00",
			},
			"availability_zone_mode": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "single",
				ValidateFunc: validation.StringInSlice([]string{
					"single", "multi",
				}, true),
			},
			"master_availability_zone": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"datastore": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"engine": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
							ValidateFunc: validation.StringInSlice([]string{
								"gaussdb-mysql",
							}, true),
						},
						"version": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
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
							Optional: true,
						},
					},
				},
			},
			"proxy_flavor": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"proxy_node_num": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"force_import": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"proxy_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"proxy_port": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"port": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"mode": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"private_write_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"db_user_name": {
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
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"private_read_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"availability_zone": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			// charge info: charging_mode, period_unit, period, auto_renew
			// make ForceNew false here but do nothing in update method!
			"charging_mode": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"prePaid", "postPaid",
				}, false),
			},
			"period_unit": {
				Type:         schema.TypeString,
				Optional:     true,
				RequiredWith: []string{"period"},
				ValidateFunc: validation.StringInSlice([]string{
					"month", "year",
				}, false),
			},
			"period": {
				Type:         schema.TypeInt,
				Optional:     true,
				RequiredWith: []string{"period_unit"},
				ValidateFunc: validation.IntBetween(1, 9),
			},
			"auto_renew": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"true", "false",
				}, false),
			},
		},
	}
}

func resourceGaussDBDataStore(d *schema.ResourceData) instances.DataStoreOpt {
	var db instances.DataStoreOpt

	datastoreRaw := d.Get("datastore").([]interface{})
	if len(datastoreRaw) == 1 {
		datastore := datastoreRaw[0].(map[string]interface{})
		db.Type = datastore["engine"].(string)
		db.Version = datastore["version"].(string)
	} else {
		db.Type = "gaussdb-mysql"
		db.Version = "8.0"
	}
	return db
}

func GaussDBInstanceStateRefreshFunc(client *golangsdk.ServiceClient, instanceID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		v, err := instances.Get(client, instanceID).Extract()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return v, "DELETED", nil
			}
			return nil, "", err
		}

		if v.Id == "" {
			return v, "DELETED", nil
		}
		return v, v.Status, nil
	}
}

func resourceGaussDBInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	client, err := config.GaussdbV3Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud GaussDB client: %s ", err)
	}

	// If force_import set, try to import it instead of creating
	if hasFilledOpt(d, "force_import") {
		logp.Printf("[DEBUG] Gaussdb mysql instance force_import is set, try to import it instead of creating")
		listOpts := instances.ListTaurusDBInstanceOpts{
			Name: d.Get("name").(string),
		}
		pages, err := instances.List(client, listOpts).AllPages()
		if err != nil {
			return err
		}

		allInstances, err := instances.ExtractTaurusDBInstances(pages)
		if err != nil {
			return fmtp.Errorf("Unable to retrieve instances: %s ", err)
		}
		if allInstances.TotalCount > 0 {
			instance := allInstances.Instances[0]
			logp.Printf("[DEBUG] Found existing mysql instance %s with name %s", instance.Id, instance.Name)
			d.SetId(instance.Id)
			return resourceGaussDBInstanceRead(d, meta)
		}
	}

	createOpts := instances.CreateTaurusDBOpts{
		Name:                d.Get("name").(string),
		Flavor:              d.Get("flavor").(string),
		Region:              GetRegion(d, config),
		VpcId:               d.Get("vpc_id").(string),
		SubnetId:            d.Get("subnet_id").(string),
		SecurityGroupId:     d.Get("security_group_id").(string),
		ConfigurationId:     d.Get("configuration_id").(string),
		EnterpriseProjectId: GetEnterpriseProjectID(d, config),
		TimeZone:            d.Get("time_zone").(string),
		SlaveCount:          d.Get("read_replicas").(int),
		Mode:                "Cluster",
		DataStore:           resourceGaussDBDataStore(d),
	}

	if d.Get("table_name_case_sensitivity").(bool) {
		lowerCaseTableNames := 0
		createOpts.LowerCaseTableNames = &lowerCaseTableNames
	}

	azMode := d.Get("availability_zone_mode").(string)
	createOpts.AZMode = azMode
	if azMode == "multi" {
		v, exist := d.GetOk("master_availability_zone")
		if !exist {
			return fmtp.Errorf("missing master_availability_zone in a multi availability zone mode")
		}
		createOpts.MasterAZ = v.(string)
	}

	if hasFilledOpt(d, "volume_size") {
		volume := &instances.VolumeOpt{
			Size: d.Get("volume_size").(int),
		}
		createOpts.Volume = volume
	}

	// PrePaid
	if d.Get("charging_mode") == "prePaid" {
		if err := validatePrePaidChargeInfo(d); err != nil {
			return err
		}

		chargeInfo := &instances.ChargeInfoOpt{
			ChargingMode: d.Get("charging_mode").(string),
			PeriodType:   d.Get("period_unit").(string),
			PeriodNum:    d.Get("period").(int),
			IsAutoPay:    "true",
			IsAutoRenew:  d.Get("auto_renew").(string),
		}
		createOpts.ChargeInfo = chargeInfo
	}

	logp.Printf("[DEBUG] Create Options: %#v", createOpts)
	// Add password here so it wouldn't go in the above log entry
	createOpts.Password = d.Get("password").(string)

	instance, err := instances.Create(client, createOpts).Extract()
	if err != nil {
		return fmtp.Errorf("Error creating GaussDB instance : %s", err)
	}

	id := instance.Instance.Id
	d.SetId(id)

	// waiting for the instance to become ready
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"BUILD", "BACKING UP"},
		Target:       []string{"ACTIVE"},
		Refresh:      GaussDBInstanceStateRefreshFunc(client, id),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        180 * time.Second,
		PollInterval: 20 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmtp.Errorf(
			"Error waiting for instance (%s) to become ready: %s",
			id, err)
	}

	if hasFilledOpt(d, "backup_strategy") {
		var updateOpts backups.UpdateOpts
		backupRaw := d.Get("backup_strategy").([]interface{})
		rawMap := backupRaw[0].(map[string]interface{})
		keep_days := rawMap["keep_days"].(int)
		updateOpts.KeepDays = &keep_days
		updateOpts.StartTime = rawMap["start_time"].(string)
		// Fixed to "1,2,3,4,5,6,7"
		updateOpts.Period = "1,2,3,4,5,6,7"
		logp.Printf("[DEBUG] Update backup_strategy: %#v", updateOpts)

		err = backups.Update(client, id, updateOpts).ExtractErr()
		if err != nil {
			return fmtp.Errorf("Error updating backup_strategy: %s", err)
		}
	}

	if hasFilledOpt(d, "proxy_flavor") {
		proxyOpts := instances.ProxyOpts{
			Flavor:  d.Get("proxy_flavor").(string),
			NodeNum: d.Get("proxy_node_num").(int),
		}
		logp.Printf("[DEBUG] Enable proxy: %#v", proxyOpts)

		n, err := instances.EnableProxy(client, id, proxyOpts).ExtractJobResponse()
		if err != nil {
			return fmtp.Errorf("Error enabling proxy: %s", err)
		}

		if err := instances.WaitForJobSuccess(client, int(d.Timeout(schema.TimeoutCreate)/time.Second), n.JobID); err != nil {
			return err
		}
	}

	// This is a workaround to avoid db connection issue
	time.Sleep(360 * time.Second) //lintignore:R018

	return resourceGaussDBInstanceRead(d, meta)
}

func resourceGaussDBInstanceRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	region := GetRegion(d, config)
	client, err := config.GaussdbV3Client(region)
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud GaussDB client: %s", err)
	}

	instanceID := d.Id()
	instance, err := instances.Get(client, instanceID).Extract()
	if err != nil {
		return CheckDeleted(d, err, "GaussDB instance")
	}
	if instance.Id == "" {
		d.SetId("")
		return nil
	}

	logp.Printf("[DEBUG] Retrieved instance %s: %#v", instanceID, instance)

	d.Set("region", region)
	d.Set("name", instance.Name)
	d.Set("status", instance.Status)
	d.Set("mode", instance.Type)
	d.Set("vpc_id", instance.VpcId)
	d.Set("subnet_id", instance.SubnetId)
	d.Set("security_group_id", instance.SecurityGroupId)
	d.Set("configuration_id", instance.ConfigurationId)
	d.Set("db_user_name", instance.DbUserName)
	d.Set("time_zone", instance.TimeZone)
	d.Set("availability_zone_mode", instance.AZMode)
	d.Set("master_availability_zone", instance.MasterAZ)

	if dbPort, err := strconv.Atoi(instance.Port); err == nil {
		d.Set("port", dbPort)
	}
	if len(instance.PrivateIps) > 0 {
		d.Set("private_write_ip", instance.PrivateIps[0])
	}

	// set data store
	dbList := make([]map[string]interface{}, 1)
	db := map[string]interface{}{
		"version": instance.DataStore.Version,
	}
	// normalize engine
	engine := instance.DataStore.Type
	if engine == "GaussDB(for MySQL)" {
		engine = "gaussdb-mysql"
	}
	db["engine"] = engine
	dbList[0] = db
	d.Set("datastore", dbList)

	// set nodes
	flavor := ""
	slave_count := 0
	volume_size := 0
	nodesList := make([]map[string]interface{}, 0, 1)
	for _, raw := range instance.Nodes {
		node := map[string]interface{}{
			"id":                raw.Id,
			"name":              raw.Name,
			"status":            raw.Status,
			"type":              raw.Type,
			"availability_zone": raw.AvailabilityZone,
		}
		if len(raw.PrivateIps) > 0 {
			node["private_read_ip"] = raw.PrivateIps[0]
		}
		if raw.Volume.Size > 0 {
			volume_size = raw.Volume.Size
		}
		nodesList = append(nodesList, node)
		if raw.Type == "slave" && raw.Status == "ACTIVE" {
			slave_count += 1
		}
		if flavor == "" {
			flavor = raw.Flavor
		}
	}
	d.Set("nodes", nodesList)
	d.Set("read_replicas", slave_count)
	d.Set("volume_size", volume_size)
	if flavor != "" {
		logp.Printf("[DEBUG] Node Flavor: %s", flavor)
		d.Set("flavor", flavor)
	}

	// set backup_strategy
	backupStrategyList := make([]map[string]interface{}, 1)
	backupStrategy := map[string]interface{}{
		"start_time": instance.BackupStrategy.StartTime,
	}
	if days, err := strconv.Atoi(instance.BackupStrategy.KeepDays); err == nil {
		backupStrategy["keep_days"] = days
	}
	backupStrategyList[0] = backupStrategy
	d.Set("backup_strategy", backupStrategyList)

	// set proxy
	proxy, err := instances.GetProxy(client, instanceID).Extract()
	if err != nil {
		logp.Printf("[DEBUG] Instance %s Proxy not enabled: %s", instanceID, err)
	} else {
		d.Set("proxy_flavor", proxy.Flavor)
		d.Set("proxy_node_num", proxy.NodeNum)
		d.Set("proxy_address", proxy.Address)
		d.Set("proxy_port", proxy.Port)
	}

	return nil
}

func resourceGaussDBInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	client, err := config.GaussdbV3Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud GaussDB client: %s ", err)
	}
	bssClient, err := config.BssV2Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud bss V2 client: %s", err)
	}
	instanceId := d.Id()

	if d.HasChange("name") {
		newName := d.Get("name").(string)
		updateNameOpts := instances.UpdateNameOpts{
			Name: newName,
		}
		logp.Printf("[DEBUG] Update Name Options: %+v", updateNameOpts)

		n, err := instances.UpdateName(client, instanceId, updateNameOpts).ExtractJobResponse()
		if err != nil {
			return fmtp.Errorf("Error updating name for instance %s: %s ", instanceId, err)
		}

		if err := instances.WaitForJobSuccess(client, int(d.Timeout(schema.TimeoutUpdate)/time.Second), n.JobID); err != nil {
			return err
		}
		logp.Printf("[DEBUG] Updated Name to %s for instance %s", newName, instanceId)
	}

	if d.HasChange("password") {
		newPass := d.Get("password").(string)
		updatePassOpts := instances.UpdatePassOpts{
			Password: newPass,
		}
		logp.Printf("[DEBUG] Update Password Options: %+v", updatePassOpts)

		_, err := instances.UpdatePass(client, instanceId, updatePassOpts).ExtractJobResponse()
		if err != nil {
			return fmtp.Errorf("Error updating password for instance %s: %s ", instanceId, err)
		}
		logp.Printf("[DEBUG] Updated Password for instance %s", instanceId)
	}

	if d.HasChange("flavor") {
		newFlavor := d.Get("flavor").(string)
		resizeOpts := instances.ResizeOpts{
			Resize: instances.ResizeOpt{
				Spec: newFlavor,
			},
		}
		if d.Get("charging_mode") == "prePaid" {
			resizeOpts.IsAutoPay = "true"
		}
		logp.Printf("[DEBUG] Update Flavor Options: %+v", resizeOpts)

		n, err := instances.Resize(client, instanceId, resizeOpts).ExtractJobResponse()
		if err != nil {
			return fmtp.Errorf("Error updating flavor for instance %s: %s ", instanceId, err)
		}

		// wait for job success
		if n.JobID != "" {
			if err := instances.WaitForJobSuccess(client, int(d.Timeout(schema.TimeoutUpdate)/time.Second), n.JobID); err != nil {
				return err
			}
		}
		// wait for order success
		if n.OrderID != "" {
			if err := orders.WaitForOrderSuccess(bssClient, int(d.Timeout(schema.TimeoutUpdate)/time.Second), n.OrderID); err != nil {
				return err
			}
			// check whether the order take effect
			instance, err := instances.Get(client, instanceId).Extract()
			if err != nil {
				return err
			}
			currFlavor := ""
			for _, raw := range instance.Nodes {
				if currFlavor == "" {
					currFlavor = raw.Flavor
					break
				}
			}
			if currFlavor != newFlavor {
				return fmtp.Errorf("Error updating flavor for instance %s: order failed", instanceId)
			}
		}
		logp.Printf("[DEBUG] Updated Flavor for instance %s", instanceId)
	}

	if d.HasChange("read_replicas") {
		old, newnum := d.GetChange("read_replicas")
		if newnum.(int) > old.(int) {
			expand_size := newnum.(int) - old.(int)
			priorities := []int{}
			for i := 0; i < expand_size; i++ {
				priorities = append(priorities, 1)
			}
			createReplicaOpts := instances.CreateReplicaOpts{
				Priorities: priorities,
			}
			if d.Get("charging_mode") == "prePaid" {
				createReplicaOpts.IsAutoPay = "true"
			}
			logp.Printf("[DEBUG] Create Replica Options: %+v", createReplicaOpts)

			n, err := instances.CreateReplica(client, instanceId, createReplicaOpts).ExtractJobResponse()
			if err != nil {
				return fmtp.Errorf("Error creating read replicas for instance %s: %s ", instanceId, err)
			}

			// wait for job success
			if n.JobID != "" {
				job_list := strings.Split(n.JobID, ",")
				logp.Printf("[DEBUG] Create Replica Jobs: %#v", job_list)
				for i := 0; i < len(job_list); i++ {
					job_id := job_list[i]
					logp.Printf("[DEBUG] Waiting for job: %s", job_id)
					if err := instances.WaitForJobSuccess(client, int(d.Timeout(schema.TimeoutUpdate)/time.Second), job_id); err != nil {
						return err
					}
				}
			}
			// wait for order success
			if n.OrderID != "" {
				if err := orders.WaitForOrderSuccess(bssClient, int(d.Timeout(schema.TimeoutUpdate)/time.Second), n.OrderID); err != nil {
					return err
				}
				// check whether the order take effect
				instance, err := instances.Get(client, instanceId).Extract()
				if err != nil {
					return err
				}
				slave_count := 0
				for _, raw := range instance.Nodes {
					if raw.Type == "slave" && raw.Status == "ACTIVE" {
						slave_count += 1
					}
				}
				if newnum.(int) != slave_count {
					return fmtp.Errorf("Error updating read_replicas for instance %s: order failed", instanceId)
				}
			}
		}
		if newnum.(int) < old.(int) {
			shrink_size := old.(int) - newnum.(int)

			slave_nodes := []string{}
			nodes := d.Get("nodes").([]interface{})
			for _, nodeRaw := range nodes {
				node := nodeRaw.(map[string]interface{})
				if node["type"].(string) == "slave" && node["status"] == "ACTIVE" {
					slave_nodes = append(slave_nodes, node["id"].(string))
				}
			}
			logp.Printf("[DEBUG] Slave Nodes: %+v", slave_nodes)
			if len(slave_nodes) <= shrink_size {
				return fmtp.Errorf("Error deleting read replicas for instance %s: Shrink Size is bigger than active slave nodes", instanceId)
			}
			for i := 0; i < shrink_size; i++ {
				n, err := instances.DeleteReplica(client, instanceId, slave_nodes[i]).ExtractJobResponse()
				if err != nil {
					return fmtp.Errorf("Error creating read replica %s for instance %s: %s ", slave_nodes[i], instanceId, err)
				}

				if err := instances.WaitForJobSuccess(client, int(d.Timeout(schema.TimeoutUpdate)/time.Second), n.JobID); err != nil {
					return err
				}
				logp.Printf("[DEBUG] Deleted Read Replica: %s", slave_nodes[i])
			}
		}
	}

	if d.HasChange("volume_size") {
		extendOpts := instances.ExtendVolumeOpts{
			Size:      d.Get("volume_size").(int),
			IsAutoPay: "true",
		}
		logp.Printf("[DEBUG] Extending Volume: %#v", extendOpts)

		n, err := instances.ExtendVolume(client, d.Id(), extendOpts).ExtractJobResponse()
		if err != nil {
			return fmtp.Errorf("Error extending volume: %s", err)
		}

		// wait for order success
		if n.OrderID != "" {
			if err := orders.WaitForOrderSuccess(bssClient, int(d.Timeout(schema.TimeoutUpdate)/time.Second), n.OrderID); err != nil {
				return err
			}
			// check whether the order take effect
			instance, err := instances.Get(client, instanceId).Extract()
			if err != nil {
				return err
			}
			volume_size := 0
			for _, raw := range instance.Nodes {
				if raw.Volume.Size > 0 {
					volume_size = raw.Volume.Size
					break
				}
			}
			if volume_size != d.Get("volume_size").(int) {
				return fmtp.Errorf("Error updating volume for instance %s: order failed", instanceId)
			}
		}
	}

	if d.HasChange("backup_strategy") {
		var updateOpts backups.UpdateOpts
		backupRaw := d.Get("backup_strategy").([]interface{})
		rawMap := backupRaw[0].(map[string]interface{})
		keep_days := rawMap["keep_days"].(int)
		updateOpts.KeepDays = &keep_days
		updateOpts.StartTime = rawMap["start_time"].(string)
		// Fixed to "1,2,3,4,5,6,7"
		updateOpts.Period = "1,2,3,4,5,6,7"
		logp.Printf("[DEBUG] Update backup_strategy: %#v", updateOpts)

		err = backups.Update(client, d.Id(), updateOpts).ExtractErr()
		if err != nil {
			return fmtp.Errorf("Error updating backup_strategy: %s", err)
		}
	}

	if d.HasChanges("proxy_flavor", "proxy_node_num") {
		if hasFilledOpt(d, "proxy_flavor") {
			proxyOpts := instances.ProxyOpts{
				Flavor:  d.Get("proxy_flavor").(string),
				NodeNum: d.Get("proxy_node_num").(int),
			}
			logp.Printf("[DEBUG] Enable proxy: %#v", proxyOpts)

			ep, err := instances.EnableProxy(client, d.Id(), proxyOpts).ExtractJobResponse()
			if err != nil {
				return fmtp.Errorf("Error enabling proxy: %s", err)
			}

			if err = instances.WaitForJobSuccess(client, int(d.Timeout(schema.TimeoutUpdate)/time.Second), ep.JobID); err != nil {
				return err
			}
		} else {
			dp, err := instances.DeleteProxy(client, d.Id()).ExtractJobResponse()
			if err != nil {
				return fmtp.Errorf("Error disabling proxy: %s", err)
			}

			if err = instances.WaitForJobSuccess(client, int(d.Timeout(schema.TimeoutUpdate)/time.Second), dp.JobID); err != nil {
				return err
			}
		}
	}

	return resourceGaussDBInstanceRead(d, meta)
}

func resourceGaussDBInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	client, err := config.GaussdbV3Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud GaussDB client: %s ", err)
	}

	instanceId := d.Id()
	if d.Get("charging_mode") == "prePaid" {
		if err := UnsubscribePrePaidResource(d, config, []string{instanceId}); err != nil {
			// try to delete the instance directly if unsubscribing failed
			res := instances.Delete(client, instanceId)
			if res.Err != nil {
				return CheckDeleted(d, res.Err, "GaussDB instance")
			}
		}
	} else {
		result := instances.Delete(client, instanceId)
		if result.Err != nil {
			return CheckDeleted(d, result.Err, "GaussDB instance")
		}
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"ACTIVE", "BACKING UP", "FAILED"},
		Target:     []string{"DELETED"},
		Refresh:    GaussDBInstanceStateRefreshFunc(client, instanceId),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmtp.Errorf(
			"Error waiting for instance (%s) to be deleted: %s ",
			instanceId, err)
	}
	logp.Printf("[DEBUG] Successfully deleted instance %s", instanceId)
	return nil
}
