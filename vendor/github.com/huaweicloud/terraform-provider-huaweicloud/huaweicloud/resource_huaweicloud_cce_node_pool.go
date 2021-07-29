package huaweicloud

import (
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/openstack/cce/v3/nodepools"
	"github.com/huaweicloud/golangsdk/openstack/cce/v3/nodes"
	"github.com/huaweicloud/golangsdk/openstack/common/tags"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"
)

func ResourceCCENodePool() *schema.Resource {
	return &schema.Resource{
		Create: resourceCCENodePoolCreate,
		Read:   resourceCCENodePoolRead,
		Update: resourceCCENodePoolUpdate,
		Delete: resourceCCENodePoolDelete,
		Importer: &schema.ResourceImporter{
			State: resourceCCENodePoolV3Import,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
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
			"initial_node_count": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"cluster_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"flavor_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"labels": { //(k8s_tags)
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"root_volume": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"size": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"volumetype": {
							Type:     schema.TypeString,
							Required: true,
						},
						"hw_passthrough": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"extend_param": {
							Type:       schema.TypeString,
							Optional:   true,
							Deprecated: "use extend_params instead",
						},
						"extend_params": {
							Type:     schema.TypeMap,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					}},
			},
			"data_volumes": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"size": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"volumetype": {
							Type:     schema.TypeString,
							Required: true,
						},
						"hw_passthrough": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"extend_param": {
							Type:       schema.TypeString,
							Optional:   true,
							Deprecated: "use extend_params instead",
						},
						"extend_params": {
							Type:     schema.TypeMap,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					}},
			},
			"availability_zone": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "random",
			},
			"os": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"key_pair": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ExactlyOneOf: []string{"password", "key_pair"},
			},
			"password": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Sensitive:    true,
				ExactlyOneOf: []string{"password", "key_pair"},
			},
			"taints": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:     schema.TypeString,
							Required: true,
						},
						"value": {
							Type:     schema.TypeString,
							Required: true,
						},
						"effect": {
							Type:     schema.TypeString,
							Required: true,
						},
					}},
			},
			"tags": tagsSchema(),
			"billing_mode": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"max_pods": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"preinstall": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				StateFunc: func(v interface{}) string {
					switch v.(type) {
					case string:
						return installScriptHashSum(v.(string))
					default:
						return ""
					}
				},
			},
			"postinstall": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				StateFunc: func(v interface{}) string {
					switch v.(type) {
					case string:
						return installScriptHashSum(v.(string))
					default:
						return ""
					}
				},
			},
			"runtime": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"docker", "containerd",
				}, false),
			},
			"extend_param": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"subnet_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"scall_enable": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"min_node_count": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"max_node_count": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"scale_down_cooldown_time": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"priority": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceCCENodePoolTags(d *schema.ResourceData) []tags.ResourceTag {
	tagRaw := d.Get("tags").(map[string]interface{})
	return utils.ExpandResourceTags(tagRaw)
}

func resourceCCENodePoolCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	nodePoolClient, err := config.CceV3Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud CCE Node Pool client: %s", err)
	}

	// wait for the cce cluster to become available
	clusterid := d.Get("cluster_id").(string)
	stateCluster := &resource.StateChangeConf{
		Target:     []string{"Available"},
		Refresh:    waitForClusterAvailable(nodePoolClient, clusterid),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      5 * time.Second,
		MinTimeout: 5 * time.Second,
	}
	_, err = stateCluster.WaitForState()

	initialNodeCount := d.Get("initial_node_count").(int)
	createOpts := nodepools.CreateOpts{
		Kind:       "NodePool",
		ApiVersion: "v3",
		Metadata: nodepools.CreateMetaData{
			Name: d.Get("name").(string),
		},
		Spec: nodepools.CreateSpec{
			Type: d.Get("type").(string),
			NodeTemplate: nodes.Spec{
				Flavor:      d.Get("flavor_id").(string),
				Az:          d.Get("availability_zone").(string),
				Os:          d.Get("os").(string),
				RootVolume:  resourceCCERootVolume(d),
				DataVolumes: resourceCCEDataVolume(d),
				K8sTags:     resourceCCENodeK8sTags(d),
				BillingMode: 0,
				Count:       1,
				NodeNicSpec: nodes.NodeNicSpec{
					PrimaryNic: nodes.PrimaryNic{
						SubnetId: d.Get("subnet_id").(string),
					},
				},
				ExtendParam: resourceCCEExtendParam(d),
				Taints:      resourceCCETaint(d),
				UserTags:    resourceCCENodePoolTags(d),
			},
			Autoscaling: nodepools.AutoscalingSpec{
				Enable:                d.Get("scall_enable").(bool),
				MinNodeCount:          d.Get("min_node_count").(int),
				MaxNodeCount:          d.Get("max_node_count").(int),
				ScaleDownCooldownTime: d.Get("scale_down_cooldown_time").(int),
				Priority:              d.Get("priority").(int),
			},
			InitialNodeCount: &initialNodeCount,
		},
	}

	if v, ok := d.GetOk("runtime"); ok {
		createOpts.Spec.NodeTemplate.RunTime = &nodes.RunTimeSpec{
			Name: v.(string),
		}
	}

	logp.Printf("[DEBUG] Create Options: %#v", createOpts)
	// Add loginSpec here so it wouldn't go in the above log entry
	var loginSpec nodes.LoginSpec
	if hasFilledOpt(d, "key_pair") {
		loginSpec = nodes.LoginSpec{
			SshKey: d.Get("key_pair").(string),
		}
	} else if hasFilledOpt(d, "password") {
		loginSpec = nodes.LoginSpec{
			UserPassword: nodes.UserPassword{
				Username: "root",
				Password: d.Get("password").(string),
			},
		}
	}
	createOpts.Spec.NodeTemplate.Login = loginSpec

	s, err := nodepools.Create(nodePoolClient, clusterid, createOpts).Extract()
	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault403); ok {
			retryNode, err := recursiveNodePoolCreate(nodePoolClient, createOpts, clusterid, 403)
			if err == "fail" {
				return fmtp.Errorf("Error creating HuaweiCloud Node Pool")
			}
			s = retryNode
		} else {
			return fmtp.Errorf("Error creating HuaweiCloud Node Pool: %s", err)
		}
	}

	if len(s.Metadata.Id) == 0 {
		return fmtp.Errorf("Error fetching CreateNodePool id")
	}

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"Synchronizing"},
		Target:       []string{""},
		Refresh:      waitForCceNodePoolActive(nodePoolClient, clusterid, s.Metadata.Id),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        120 * time.Second,
		PollInterval: 20 * time.Second,
	}
	_, err = stateConf.WaitForState()
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud CCE Node Pool: %s", err)
	}

	logp.Printf("[DEBUG] Create node pool: %v", s)

	d.SetId(s.Metadata.Id)
	return resourceCCENodePoolRead(d, meta)
}

func resourceCCENodePoolRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	nodePoolClient, err := config.CceV3Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud CCE Node Pool client: %s", err)
	}
	clusterid := d.Get("cluster_id").(string)
	s, err := nodepools.Get(nodePoolClient, clusterid, d.Id()).Extract()

	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault404); ok {
			d.SetId("")
			return nil
		}

		return fmtp.Errorf("Error retrieving HuaweiCloud Node Pool: %s", err)
	}

	d.Set("name", s.Metadata.Name)
	d.Set("flavor_id", s.Spec.NodeTemplate.Flavor)
	d.Set("availability_zone", s.Spec.NodeTemplate.Az)
	d.Set("os", s.Spec.NodeTemplate.Os)
	d.Set("billing_mode", s.Spec.NodeTemplate.BillingMode)
	d.Set("key_pair", s.Spec.NodeTemplate.Login.SshKey)
	d.Set("initial_node_count", s.Spec.InitialNodeCount)
	d.Set("scall_enable", s.Spec.Autoscaling.Enable)
	d.Set("min_node_count", s.Spec.Autoscaling.MinNodeCount)
	d.Set("max_node_count", s.Spec.Autoscaling.MaxNodeCount)
	d.Set("scale_down_cooldown_time", s.Spec.Autoscaling.ScaleDownCooldownTime)
	d.Set("priority", s.Spec.Autoscaling.Priority)
	d.Set("type", s.Spec.Type)

	// set extend_param
	var extend_param = s.Spec.NodeTemplate.ExtendParam
	d.Set("max_pods", extend_param["maxPods"].(float64))
	delete(extend_param, "maxPods")
	d.Set("extend_param", extend_param)

	if s.Spec.NodeTemplate.RunTime != nil {
		d.Set("runtime", s.Spec.NodeTemplate.RunTime.Name)
	}

	labels := map[string]string{}
	for key, val := range s.Spec.NodeTemplate.K8sTags {
		if strings.Contains(key, "cce.cloud.com") {
			continue
		}
		labels[key] = val
	}
	d.Set("labels", labels)

	var volumes []map[string]interface{}
	for _, pairObject := range s.Spec.NodeTemplate.DataVolumes {
		volume := make(map[string]interface{})
		volume["size"] = pairObject.Size
		volume["volumetype"] = pairObject.VolumeType
		volume["hw_passthrough"] = pairObject.HwPassthrough
		volume["extend_params"] = pairObject.ExtendParam
		volume["extend_param"] = ""
		volumes = append(volumes, volume)
	}
	if err := d.Set("data_volumes", volumes); err != nil {
		return fmtp.Errorf("[DEBUG] Error saving dataVolumes to state for HuaweiCloud Node Pool (%s): %s", d.Id(), err)
	}

	rootVolume := []map[string]interface{}{
		{
			"size":           s.Spec.NodeTemplate.RootVolume.Size,
			"volumetype":     s.Spec.NodeTemplate.RootVolume.VolumeType,
			"hw_passthrough": s.Spec.NodeTemplate.RootVolume.HwPassthrough,
			"extend_params":  s.Spec.NodeTemplate.RootVolume.ExtendParam,
			"extend_param":   "",
		},
	}
	if err := d.Set("root_volume", rootVolume); err != nil {
		return fmtp.Errorf("[DEBUG] Error saving root Volume to state for HuaweiCloud Node Pool (%s): %s", d.Id(), err)
	}

	tagmap := utils.TagsToMap(s.Spec.NodeTemplate.UserTags)
	// ignore "CCE-Dynamic-Provisioning-Node"
	delete(tagmap, "CCE-Dynamic-Provisioning-Node")
	if err := d.Set("tags", tagmap); err != nil {
		return fmtp.Errorf("Error saving tags to state for CCE Node Pool(%s): %s", d.Id(), err)
	}

	d.Set("status", s.Status.Phase)

	return nil
}

func resourceCCENodePoolUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	nodePoolClient, err := config.CceV3Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud CCE client: %s", err)
	}

	initialNodeCount := d.Get("initial_node_count").(int)
	var loginSpec nodes.LoginSpec
	if hasFilledOpt(d, "key_pair") {
		loginSpec = nodes.LoginSpec{SshKey: d.Get("key_pair").(string)}
	} else if hasFilledOpt(d, "password") {
		loginSpec = nodes.LoginSpec{
			UserPassword: nodes.UserPassword{
				Username: "root",
				Password: d.Get("password").(string),
			},
		}
	}

	updateOpts := nodepools.UpdateOpts{
		Kind:       "NodePool",
		ApiVersion: "v3",
		Metadata: nodepools.UpdateMetaData{
			Name: d.Get("name").(string),
		},
		Spec: nodepools.UpdateSpec{
			InitialNodeCount: &initialNodeCount,
			Autoscaling: nodepools.AutoscalingSpec{
				Enable:                d.Get("scall_enable").(bool),
				MinNodeCount:          d.Get("min_node_count").(int),
				MaxNodeCount:          d.Get("max_node_count").(int),
				ScaleDownCooldownTime: d.Get("scale_down_cooldown_time").(int),
				Priority:              d.Get("priority").(int),
			},
			NodeTemplate: nodes.Spec{
				Flavor:      d.Get("flavor_id").(string),
				Az:          d.Get("availability_zone").(string),
				Login:       loginSpec,
				RootVolume:  resourceCCERootVolume(d),
				DataVolumes: resourceCCEDataVolume(d),
				Count:       1,
				UserTags:    resourceCCENodePoolTags(d),
			},
			Type: d.Get("type").(string),
		},
	}

	clusterid := d.Get("cluster_id").(string)
	_, err = nodepools.Update(nodePoolClient, clusterid, d.Id(), updateOpts).Extract()
	if err != nil {
		return fmtp.Errorf("Error updating HuaweiCloud Node Node Pool: %s", err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"Synchronizing"},
		Target:       []string{""},
		Refresh:      waitForCceNodePoolActive(nodePoolClient, clusterid, d.Id()),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        15 * time.Second,
		PollInterval: 10 * time.Second,
	}
	_, err = stateConf.WaitForState()
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud CCE Node Pool: %s", err)
	}

	return resourceCCENodePoolRead(d, meta)
}

func resourceCCENodePoolDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	nodePoolClient, err := config.CceV3Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud CCE client: %s", err)
	}
	clusterid := d.Get("cluster_id").(string)
	err = nodepools.Delete(nodePoolClient, clusterid, d.Id()).ExtractErr()
	if err != nil {
		return fmtp.Errorf("Error deleting HuaweiCloud CCE Node Pool: %s", err)
	}
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"Deleting"},
		Target:       []string{"Deleted"},
		Refresh:      waitForCceNodePoolDelete(nodePoolClient, clusterid, d.Id()),
		Timeout:      d.Timeout(schema.TimeoutDelete),
		Delay:        60 * time.Second,
		PollInterval: 20 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmtp.Errorf("Error deleting HuaweiCloud CCE Node Pool: %s", err)
	}

	d.SetId("")
	return nil
}

func waitForCceNodePoolActive(cceClient *golangsdk.ServiceClient, clusterId, nodePoolId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		n, err := nodepools.Get(cceClient, clusterId, nodePoolId).Extract()
		if err != nil {
			return nil, "", err
		}
		return n, n.Status.Phase, nil
	}
}

func waitForCceNodePoolDelete(cceClient *golangsdk.ServiceClient, clusterId, nodePoolId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		logp.Printf("[DEBUG] Attempting to delete HuaweiCloud CCE Node Pool %s.\n", nodePoolId)

		r, err := nodepools.Get(cceClient, clusterId, nodePoolId).Extract()

		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				logp.Printf("[DEBUG] Successfully deleted HuaweiCloud CCE Node Pool %s", nodePoolId)
				return r, "Deleted", nil
			}
			return r, "Deleting", err
		}

		logp.Printf("[DEBUG] HuaweiCloud CCE Node Pool %s still available.\n", nodePoolId)
		return r, r.Status.Phase, nil
	}
}

func recursiveNodePoolCreate(cceClient *golangsdk.ServiceClient, opts nodepools.CreateOptsBuilder, ClusterID string, errCode int) (*nodepools.NodePool, string) {
	if errCode == 403 {
		stateCluster := &resource.StateChangeConf{
			Target:       []string{"Available"},
			Refresh:      waitForClusterAvailable(cceClient, ClusterID),
			Timeout:      15 * time.Minute,
			Delay:        15 * time.Second,
			PollInterval: 10 * time.Second,
		}
		_, stateErr := stateCluster.WaitForState()
		if stateErr != nil {
			logp.Printf("[INFO] Cluster Unavailable %s.\n", stateErr)
		}
		s, err := nodepools.Create(cceClient, ClusterID, opts).Extract()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault403); ok {
				return recursiveNodePoolCreate(cceClient, opts, ClusterID, 403)
			} else {
				return s, "fail"
			}
		} else {
			return s, "success"
		}
	}
	return nil, "fail"
}

func resourceCCENodePoolV3Import(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	parts := strings.SplitN(d.Id(), "/", 2)
	if len(parts) != 2 {
		err := fmtp.Errorf("Invalid format specified for CCE Node Pool. Format must be <cluster id>/<node pool id>")
		return nil, err
	}

	clusterID := parts[0]
	nodePoolID := parts[1]

	d.SetId(nodePoolID)
	d.Set("cluster_id", clusterID)

	return []*schema.ResourceData{d}, nil
}
