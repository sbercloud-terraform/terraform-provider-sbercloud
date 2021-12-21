package huaweicloud

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/opengauss/v3/backups"
	"github.com/chnsz/golangsdk/openstack/opengauss/v3/instances"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"
)

func resourceOpenGaussInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceOpenGaussInstanceCreate,
		Read:   resourceOpenGaussInstanceRead,
		Delete: resourceOpenGaussInstanceDelete,
		Update: resourceOpenGaussInstanceUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(120 * time.Minute),
			Update: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},
		CustomizeDiff: func(_ context.Context, d *schema.ResourceDiff, v interface{}) error {
			if d.HasChange("coordinator_num") {
				d.SetNewComputed("private_ips")
			}
			return nil
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
				ForceNew: true,
			},
			"availability_zone": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
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
			"port": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"configuration_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"sharding_num": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      3,
				ValidateFunc: validation.IntBetween(1, 9),
			},
			"coordinator_num": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      3,
				ValidateFunc: validation.IntBetween(1, 9),
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"time_zone": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "UTC+08:00",
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
						},
						"version": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ForceNew: true,
						},
					},
				},
			},
			"ha": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"mode": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"replication_mode": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"consistency": {
							Type:     schema.TypeString,
							Optional: true,
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
			"volume": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"size": {
							Type:     schema.TypeInt,
							Required: true,
						},
					},
				},
			},
			"force_import": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
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
			"endpoints": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"db_user_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"switch_strategy": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"maintenance_window": {
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
						"role": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
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
		},
	}
}

func resourceOpenGaussDataStore(d *schema.ResourceData) instances.DataStoreOpt {
	var db instances.DataStoreOpt

	datastoreRaw := d.Get("datastore").([]interface{})
	if len(datastoreRaw) == 1 {
		datastore := datastoreRaw[0].(map[string]interface{})
		db.Type = datastore["engine"].(string)
		db.Version = datastore["version"].(string)
	} else {
		db.Type = "GaussDB(for openGauss)"
	}
	return db
}

func resourceOpenGaussBackupStrategy(d *schema.ResourceData) *instances.BackupStrategyOpt {
	var backupOpt instances.BackupStrategyOpt

	backupStrategyRaw := d.Get("backup_strategy").([]interface{})
	if len(backupStrategyRaw) == 1 {
		strategy := backupStrategyRaw[0].(map[string]interface{})
		backupOpt.StartTime = strategy["start_time"].(string)
		backupOpt.KeepDays = strategy["keep_days"].(int)
		return &backupOpt
	}

	return nil
}

func OpenGaussInstanceStateRefreshFunc(client *golangsdk.ServiceClient, instanceID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		v, err := instances.GetInstanceByID(client, instanceID)
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

func resourceOpenGaussInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	client, err := config.OpenGaussV3Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud GaussDB client: %s ", err)
	}

	// If force_import set, try to import it instead of creating
	if hasFilledOpt(d, "force_import") {
		logp.Printf("[DEBUG] Gaussdb opengauss instance force_import is set, try to import it instead of creating")
		listOpts := instances.ListGaussDBInstanceOpts{
			Name: d.Get("name").(string),
		}
		pages, err := instances.List(client, listOpts).AllPages()
		if err != nil {
			return err
		}

		allInstances, err := instances.ExtractGaussDBInstances(pages)
		if err != nil {
			return fmtp.Errorf("Unable to retrieve instances: %s ", err)
		}
		if allInstances.TotalCount > 0 {
			instance := allInstances.Instances[0]
			logp.Printf("[DEBUG] Found existing opengauss instance %s with name %s", instance.Id, instance.Name)
			d.SetId(instance.Id)
			return resourceOpenGaussInstanceRead(d, meta)
		}
	}

	createOpts := instances.CreateGaussDBOpts{
		Name:                d.Get("name").(string),
		Flavor:              d.Get("flavor").(string),
		Region:              GetRegion(d, config),
		VpcId:               d.Get("vpc_id").(string),
		SubnetId:            d.Get("subnet_id").(string),
		SecurityGroupId:     d.Get("security_group_id").(string),
		Port:                d.Get("port").(string),
		EnterpriseProjectId: GetEnterpriseProjectID(d, config),
		TimeZone:            d.Get("time_zone").(string),
		AvailabilityZone:    d.Get("availability_zone").(string),
		ConfigurationId:     d.Get("configuration_id").(string),
		ShardingNum:         d.Get("sharding_num").(int),
		CoordinatorNum:      d.Get("coordinator_num").(int),
		DataStore:           resourceOpenGaussDataStore(d),
		BackupStrategy:      resourceOpenGaussBackupStrategy(d),
	}

	haRaw := d.Get("ha").([]interface{})
	if len(haRaw) > 0 {
		logp.Printf("[DEBUG] ha: %+v", haRaw)
		ha := haRaw[0].(map[string]interface{})
		createOpts.Ha = &instances.HaOpt{
			Mode:            ha["mode"].(string),
			ReplicationMode: ha["replication_mode"].(string),
			Consistency:     ha["consistency"].(string),
		}
	}

	dn_num := d.Get("sharding_num").(int)
	volumeRaw := d.Get("volume").([]interface{})
	if len(volumeRaw) > 0 {
		logp.Printf("[DEBUG] volume: %+v", volumeRaw)
		volume := volumeRaw[0].(map[string]interface{})
		dn_size := volume["size"].(int)
		volume_size := dn_size * dn_num
		createOpts.Volume = instances.VolumeOpt{
			Type: volume["type"].(string),
			Size: volume_size,
		}
	}

	logp.Printf("[DEBUG] Create Options: %#v", createOpts)
	// Add password here so it wouldn't go in the above log entry
	createOpts.Password = d.Get("password").(string)

	instance, err := instances.Create(client, createOpts).Extract()
	if err != nil {
		return fmtp.Errorf("Error creating OpenGauss instance : %s", err)
	}

	id := instance.Instance.Id
	d.SetId(id)

	// waiting for the instance to become ready
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"BUILD", "BACKING UP"},
		Target:       []string{"ACTIVE"},
		Refresh:      OpenGaussInstanceStateRefreshFunc(client, id),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        180 * time.Second,
		PollInterval: 30 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmtp.Errorf(
			"Error waiting for instance (%s) to become ready: %s",
			id, err)
	}

	// This is a workaround to avoid db connection issue
	time.Sleep(360 * time.Second) //lintignore:R018

	return resourceOpenGaussInstanceRead(d, meta)
}

func resourceOpenGaussInstanceRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	region := GetRegion(d, config)
	client, err := config.OpenGaussV3Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud GaussDB client: %s", err)
	}

	instanceID := d.Id()
	instance, err := instances.GetInstanceByID(client, instanceID)
	if err != nil {
		return CheckDeleted(d, err, "OpenGauss instance")
	}
	if instance.Id == "" {
		d.SetId("")
		return nil
	}

	logp.Printf("[DEBUG] Retrieved instance %s: %#v", instanceID, instance)

	d.Set("region", region)
	d.Set("name", instance.Name)
	d.Set("status", instance.Status)
	d.Set("type", instance.Type)
	d.Set("vpc_id", instance.VpcId)
	d.Set("subnet_id", instance.SubnetId)
	d.Set("security_group_id", instance.SecurityGroupId)
	d.Set("db_user_name", instance.DbUserName)
	d.Set("time_zone", instance.TimeZone)
	d.Set("flavor", instance.FlavorRef)
	d.Set("port", strconv.Itoa(instance.Port))
	d.Set("switch_strategy", instance.SwitchStrategy)
	d.Set("maintenance_window", instance.MaintenanceWindow)
	d.Set("public_ips", instance.PublicIps)

	if len(instance.PrivateIps) > 0 {
		private_ips := instance.PrivateIps[0]
		ip_list := strings.Split(private_ips, "/")
		endpoints := []string{}
		for i := 0; i < len(ip_list); i++ {
			ip_list[i] = strings.Trim(ip_list[i], " ")
			endpoint := fmt.Sprintf("%s:%d", ip_list[i], instance.Port)
			endpoints = append(endpoints, endpoint)
		}
		d.Set("private_ips", ip_list)
		d.Set("endpoints", endpoints)
	}

	// set data store
	dbList := make([]map[string]interface{}, 1)
	db := map[string]interface{}{
		"version": instance.DataStore.Version,
		"engine":  instance.DataStore.Type,
	}
	dbList[0] = db
	d.Set("datastore", dbList)

	// set nodes
	sharding_num := 0
	coordinator_num := 0
	nodesList := make([]map[string]interface{}, 0, 1)
	for _, raw := range instance.Nodes {
		node := map[string]interface{}{
			"id":                raw.Id,
			"name":              raw.Name,
			"status":            raw.Status,
			"role":              raw.Role,
			"availability_zone": raw.AvailabilityZone,
		}
		nodesList = append(nodesList, node)
		if strings.Contains(raw.Name, "_gaussdbv5cn") {
			coordinator_num += 1
		} else if strings.Contains(raw.Name, "_gaussdbv5dn") {
			sharding_num += 1
		}
	}
	d.Set("nodes", nodesList)
	d.Set("coordinator_num", coordinator_num)

	dn_num := sharding_num / 3
	d.Set("sharding_num", dn_num)

	// set backup_strategy
	backupStrategyList := make([]map[string]interface{}, 1)
	backupStrategy := map[string]interface{}{
		"start_time": instance.BackupStrategy.StartTime,
		"keep_days":  instance.BackupStrategy.KeepDays,
	}
	backupStrategyList[0] = backupStrategy
	d.Set("backup_strategy", backupStrategyList)

	// set volume
	volume_size := instance.Volume.Size
	dn_size := volume_size / dn_num
	volumeList := make([]map[string]interface{}, 1)
	volume := map[string]interface{}{
		"type": instance.Volume.Type,
		"size": dn_size,
	}
	volumeList[0] = volume
	d.Set("volume", volumeList)

	return nil
}

func resourceOpenGaussInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	client, err := config.OpenGaussV3Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud GaussDB client: %s ", err)
	}

	instanceId := d.Id()
	result := instances.Delete(client, instanceId)
	if result.Err != nil {
		return CheckDeleted(d, result.Err, "OpenGauss instance")
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"ACTIVE", "BACKING UP", "FAILED"},
		Target:     []string{"DELETED"},
		Refresh:    OpenGaussInstanceStateRefreshFunc(client, instanceId),
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

func resourceOpenGaussInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	client, err := config.OpenGaussV3Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud GaussDB client: %s ", err)
	}

	logp.Printf("[DEBUG] Updating OpenGaussDB instances %s", d.Id())
	instanceId := d.Id()

	dn_num := d.Get("sharding_num").(int)
	if d.HasChange("sharding_num") {
		old, newnum := d.GetChange("sharding_num")
		if newnum.(int) < old.(int) {
			return fmtp.Errorf(
				"Error expanding shard for instance (%s): new num must be larger than the old one.",
				instanceId)
		}
		dn_num = newnum.(int)
		expand_size := newnum.(int) - old.(int)
		updateClusterOpts := instances.UpdateClusterOpts{
			Shard: &instances.Shard{
				Count: expand_size,
			},
		}
		logp.Printf("[DEBUG] Expand Shard Options: %+v", updateClusterOpts)
		result := instances.UpdateCluster(client, updateClusterOpts, instanceId)
		if result.Err != nil {
			return fmtp.Errorf("Error expanding shard for instance %s: %s ", instanceId, result.Err)
		}

		stateConf := &resource.StateChangeConf{
			Pending:    []string{"MODIFYING", "EXPANDING", "BACKING UP"},
			Target:     []string{"ACTIVE"},
			Refresh:    OpenGaussInstanceStateRefreshFunc(client, instanceId),
			Timeout:    d.Timeout(schema.TimeoutUpdate),
			Delay:      60 * time.Second,
			MinTimeout: 30 * time.Second,
		}

		_, err = stateConf.WaitForState()
		if err != nil {
			return fmtp.Errorf(
				"Error waiting for instance (%s) shard to be Updated: %s ",
				instanceId, err)
		}
	}

	if d.HasChange("coordinator_num") {
		old, newnum := d.GetChange("coordinator_num")
		if newnum.(int) < old.(int) {
			return fmtp.Errorf(
				"Error expanding coordinator for instance (%s): new num must be larger than the old one.",
				instanceId)
		}
		expand_size := newnum.(int) - old.(int)

		var coordinators []instances.Coordinator
		azs := d.Get("availability_zone").(string)
		az_list := strings.Split(azs, ",")
		for i := 0; i < expand_size; i++ {
			coordinator := instances.Coordinator{
				AzCode: az_list[0],
			}
			coordinators = append(coordinators, coordinator)
		}
		updateClusterOpts := instances.UpdateClusterOpts{
			Coordinators: coordinators,
		}
		logp.Printf("[DEBUG] Expand Coordinator Options: %+v", updateClusterOpts)
		result := instances.UpdateCluster(client, updateClusterOpts, instanceId)
		if result.Err != nil {
			return fmtp.Errorf("Error expanding coordinator for instance %s: %s ", instanceId, result.Err)
		}

		stateConf := &resource.StateChangeConf{
			Pending:    []string{"MODIFYING", "EXPANDING", "BACKING UP"},
			Target:     []string{"ACTIVE"},
			Refresh:    OpenGaussInstanceStateRefreshFunc(client, instanceId),
			Timeout:    d.Timeout(schema.TimeoutUpdate),
			Delay:      60 * time.Second,
			MinTimeout: 30 * time.Second,
		}

		_, err = stateConf.WaitForState()
		if err != nil {
			return fmtp.Errorf(
				"Error waiting for instance (%s) coordinator to be Updated: %s ",
				instanceId, err)
		}
	}

	if d.HasChange("volume") {
		volumeRaw := d.Get("volume").([]interface{})
		dn_size := volumeRaw[0].(map[string]interface{})["size"].(int)
		volume_size := dn_size * dn_num
		updateVolumeOpts := instances.UpdateVolumeOpts{
			Size: volume_size,
		}
		logp.Printf("[DEBUG] Update Volume Options: %+v", updateVolumeOpts)
		result := instances.UpdateVolume(client, updateVolumeOpts, instanceId)
		if result.Err != nil {
			return fmtp.Errorf("Error updating instance %s: %s ", instanceId, result.Err)
		}

		stateConf := &resource.StateChangeConf{
			Pending:    []string{"MODIFYING", "EXPANDING", "BACKING UP"},
			Target:     []string{"ACTIVE"},
			Refresh:    OpenGaussInstanceStateRefreshFunc(client, instanceId),
			Timeout:    d.Timeout(schema.TimeoutUpdate),
			Delay:      40 * time.Second,
			MinTimeout: 20 * time.Second,
		}

		_, err = stateConf.WaitForState()
		if err != nil {
			return fmtp.Errorf(
				"Error waiting for instance (%s) volume to be Updated: %s ",
				instanceId, err)
		}
	}
	logp.Printf("[DEBUG] Successfully updated instance %s", instanceId)

	if d.HasChange("backup_strategy") {
		backupRaw := d.Get("backup_strategy").([]interface{})
		rawMap := backupRaw[0].(map[string]interface{})
		keep_days := rawMap["keep_days"].(int)

		updateOpts := backups.UpdateOpts{
			KeepDays:  &keep_days,
			StartTime: rawMap["start_time"].(string),
			// Fixed to "1,2,3,4,5,6,7"
			Period: "1,2,3,4,5,6,7",
			// Fixed to "30"
			DifferentialPeriod: "30",
		}

		logp.Printf("[DEBUG] Update backup_strategy: %#v", updateOpts)

		err = backups.Update(client, d.Id(), updateOpts).ExtractErr()
		if err != nil {
			return fmtp.Errorf("Error updating backup_strategy: %s", err)
		}
	}

	if d.HasChange("name") {
		renameOpts := instances.RenameOpts{
			Name: d.Get("name").(string),
		}
		_, err = instances.Rename(client, renameOpts, instanceId).Extract()
		if err != nil {
			return fmtp.Errorf("Error updating name for instance (%s): %s ", instanceId, err)
		}
	}

	if d.HasChange("password") {
		restorePasswordOpts := instances.RestorePasswordOpts{
			Password: d.Get("password").(string),
		}
		r := golangsdk.ErrResult{}
		r.Result = instances.RestorePassword(client, restorePasswordOpts, instanceId)
		if r.ExtractErr() != nil {
			return fmtp.Errorf("Error updating password for instance (%s): %s ", instanceId, r.Err)
		}
	}

	return resourceOpenGaussInstanceRead(d, meta)
}
