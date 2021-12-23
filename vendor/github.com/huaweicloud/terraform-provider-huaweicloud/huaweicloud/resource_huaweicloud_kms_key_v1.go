package huaweicloud

import (
	"time"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/common/tags"
	"github.com/chnsz/golangsdk/openstack/kms/v1/keys"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"
)

const WaitingForEnableState = "1"
const EnabledState = "2"
const DisabledState = "3"
const PendingDeletionState = "4"

func ResourceKmsKeyV1() *schema.Resource {
	return &schema.Resource{
		Create: resourceKmsKeyV1Create,
		Read:   resourceKmsKeyV1Read,
		Update: resourceKmsKeyV1Update,
		Delete: resourceKmsKeyV1Delete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"key_alias": {
				Type:     schema.TypeString,
				Required: true,
			},
			"key_description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"key_algorithm": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"is_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"pending_days": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "7",
			},
			"tags": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"key_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"domain_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"creation_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"scheduled_deletion_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"default_key_flag": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"expiration_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceKmsKeyV1Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	kmsKeyV1Client, err := config.KmsKeyV1Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud kms key client: %s", err)
	}

	createOpts := &keys.CreateOpts{
		KeyAlias:            d.Get("key_alias").(string),
		KeyDescription:      d.Get("key_description").(string),
		KeySpec:             d.Get("key_algorithm").(string),
		EnterpriseProjectID: GetEnterpriseProjectID(d, config),
	}

	logp.Printf("[DEBUG] Create Options: %#v", createOpts)
	v, err := keys.Create(kmsKeyV1Client, createOpts).ExtractKeyInfo()
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud key: %s", err)
	}
	logp.Printf("[INFO] Key ID: %s", v.KeyID)

	// Wait for the key to become enabled.
	logp.Printf("[DEBUG] Waiting for key (%s) to become enabled", v.KeyID)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{WaitingForEnableState, DisabledState},
		Target:     []string{EnabledState},
		Refresh:    KeyV1StateRefreshFunc(kmsKeyV1Client, v.KeyID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmtp.Errorf(
			"Error waiting for key (%s) to become ready: %s",
			v.KeyID, err)
	}

	if !d.Get("is_enabled").(bool) {
		key, err := keys.DisableKey(kmsKeyV1Client, v.KeyID).ExtractKeyInfo()
		if err != nil {
			return fmtp.Errorf("Error disabling key: %s.", err)
		}

		if key.KeyState != DisabledState {
			return fmtp.Errorf("Error disabling key, the key state is: %s", key.KeyState)
		}
	}

	tagRaw := d.Get("tags").(map[string]interface{})
	if len(tagRaw) > 0 {
		taglist := utils.ExpandResourceTags(tagRaw)
		tagErr := tags.Create(kmsKeyV1Client, "kms", v.KeyID, taglist).ExtractErr()
		if tagErr != nil {
			logp.Printf("Error creating tags for kms key(%s): %s", v.KeyID, err)
		}
	}

	// Store the key ID now
	d.SetId(v.KeyID)
	d.Set("key_id", v.KeyID)

	return resourceKmsKeyV1Read(d, meta)
}

func resourceKmsKeyV1Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)

	kmsRegion := GetRegion(d, config)
	kmsKeyV1Client, err := config.KmsKeyV1Client(kmsRegion)
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud kms key client: %s", err)
	}
	v, err := keys.Get(kmsKeyV1Client, d.Id()).ExtractKeyInfo()
	if err != nil {
		return err
	}

	logp.Printf("[DEBUG] Kms key %s: %+v", d.Id(), v)
	if v.KeyState == PendingDeletionState {
		logp.Printf("[WARN] Removing KMS key %s because it's already gone", d.Id())
		d.SetId("")
		return nil
	}

	d.SetId(v.KeyID)
	d.Set("key_id", v.KeyID)
	d.Set("domain_id", v.DomainID)
	d.Set("key_alias", v.KeyAlias)
	d.Set("region", kmsRegion)
	d.Set("key_description", v.KeyDescription)
	d.Set("key_algorithm", v.KeySpec)
	d.Set("creation_date", v.CreationDate)
	d.Set("scheduled_deletion_date", v.ScheduledDeletionDate)
	d.Set("is_enabled", v.KeyState == EnabledState)
	d.Set("default_key_flag", v.DefaultKeyFlag)
	d.Set("expiration_time", v.ExpirationTime)
	d.Set("enterprise_project_id", v.EnterpriseProjectID)

	// Set kms tags
	if resourceTags, err := tags.Get(kmsKeyV1Client, "kms", d.Id()).Extract(); err == nil {
		tagmap := utils.TagsToMap(resourceTags.Tags)
		if err := d.Set("tags", tagmap); err != nil {
			return fmtp.Errorf("Error saving tags to state for kms key(%s): %s", d.Id(), err)
		}
	} else {
		logp.Printf("[WARN] Error fetching tags of kms key(%s): %s", d.Id(), err)
	}

	return nil
}

func resourceKmsKeyV1Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	kmsKeyV1Client, err := config.KmsKeyV1Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud kms key client: %s", err)
	}

	if d.HasChange("key_alias") {
		updateAliasOpts := keys.UpdateAliasOpts{
			KeyID:    d.Id(),
			KeyAlias: d.Get("key_alias").(string),
		}
		_, err = keys.UpdateAlias(kmsKeyV1Client, updateAliasOpts).ExtractKeyInfo()
		if err != nil {
			return fmtp.Errorf("Error updating HuaweiCloud key: %s", err)
		}
	}

	if d.HasChange("key_description") {
		updateDesOpts := keys.UpdateDesOpts{
			KeyID:          d.Id(),
			KeyDescription: d.Get("key_description").(string),
		}
		_, err = keys.UpdateDes(kmsKeyV1Client, updateDesOpts).ExtractKeyInfo()
		if err != nil {
			return fmtp.Errorf("Error updating HuaweiCloud key: %s", err)
		}
	}

	if d.HasChange("is_enabled") {
		v, err := keys.Get(kmsKeyV1Client, d.Id()).ExtractKeyInfo()
		if err != nil {
			return fmtp.Errorf("DescribeKey got an error: %s.", err)
		}

		if d.Get("is_enabled").(bool) && v.KeyState == DisabledState {
			key, err := keys.EnableKey(kmsKeyV1Client, d.Id()).ExtractKeyInfo()
			if err != nil {
				return fmtp.Errorf("Error enabling key: %s.", err)
			}
			if key.KeyState != EnabledState {
				return fmtp.Errorf("Error enabling key, the key state is: %s", key.KeyState)
			}
		}

		if !d.Get("is_enabled").(bool) && v.KeyState == EnabledState {
			key, err := keys.DisableKey(kmsKeyV1Client, d.Id()).ExtractKeyInfo()
			if err != nil {
				return fmtp.Errorf("Error disabling key: %s.", err)
			}
			if key.KeyState != DisabledState {
				return fmtp.Errorf("Error disabling key, the key state is: %s", key.KeyState)
			}
		}
	}

	if d.HasChange("tags") {
		tagErr := utils.UpdateResourceTags(kmsKeyV1Client, d, "kms", d.Id())
		if tagErr != nil {
			return fmtp.Errorf("Error updating tags of kms:%s, err:%s", d.Id(), err)
		}
	}

	return resourceKmsKeyV1Read(d, meta)
}

func resourceKmsKeyV1Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	kmsKeyV1Client, err := config.KmsKeyV1Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud kms key client: %s", err)
	}

	v, err := keys.Get(kmsKeyV1Client, d.Id()).ExtractKeyInfo()
	if err != nil {
		return CheckDeleted(d, err, "key")
	}

	deleteOpts := &keys.DeleteOpts{
		KeyID: d.Id(),
	}
	if v, ok := d.GetOk("pending_days"); ok {
		deleteOpts.PendingDays = v.(string)
	}

	// It's possible that this key was used as a boot device and is currently
	// in a pending deletion state from when the instance was terminated.
	// If this is true, just move on. It'll eventually delete.
	if v.KeyState != PendingDeletionState {
		v, err = keys.Delete(kmsKeyV1Client, deleteOpts).Extract()
		if err != nil {
			return err
		}

		if v.KeyState != PendingDeletionState {
			return fmtp.Errorf("failed to delete key")
		}
	}

	logp.Printf("[DEBUG] KMS Key %s deactivated.", d.Id())
	d.SetId("")
	return nil
}

func KeyV1StateRefreshFunc(client *golangsdk.ServiceClient, keyID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		v, err := keys.Get(client, keyID).ExtractKeyInfo()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return v, PendingDeletionState, nil
			}
			return nil, "", err
		}

		return v, v.KeyState, nil
	}
}
