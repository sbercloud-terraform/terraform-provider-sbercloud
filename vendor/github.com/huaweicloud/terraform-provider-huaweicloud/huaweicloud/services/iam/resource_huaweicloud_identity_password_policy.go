package iam

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API IAM PUT /v3.0/OS-SECURITYPOLICY/domains/{domainID}/password-policy
// @API IAM GET /v3.0/OS-SECURITYPOLICY/domains/{domainID}/password-policy
func ResourceV3PasswordPolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceV3PasswordPolicyCreateOrUpdate,
		ReadContext:   resourceV3PasswordPolicyRead,
		UpdateContext: resourceV3PasswordPolicyCreateOrUpdate,
		DeleteContext: resourceV3PasswordPolicyDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"maximum_consecutive_identical_chars": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: `The maximum number of times that a character is allowed to consecutively present in a password.`,
			},
			"minimum_password_age": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: `The minimum period (minutes) after which users are allowed to make a password change.`,
			},
			"minimum_password_length": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     8,
				Description: `The minimum number of characters that a password must contain.`,
			},
			"number_of_recent_passwords_disallowed": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     1,
				Description: `The member of previously used passwords that are not allowed.`,
			},
			"password_not_username_or_invert": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: `Whether the password can be the username or the username spelled backwards.`,
			},
			"password_validity_period": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: `The password validity period (days).`,
			},
			"password_char_combination": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     2,
				Description: `The minimum number of character types that a password must contain.`,
			},
			"maximum_password_length": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The maximum number of characters that a password can contain.`,
			},
		},
	}
}

func buildV3UpdatePasswordPolicyBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"password_policy": map[string]interface{}{
			"maximum_consecutive_identical_chars":   d.Get("maximum_consecutive_identical_chars").(int),
			"minimum_password_age":                  d.Get("minimum_password_age").(int),
			"minimum_password_length":               d.Get("minimum_password_length").(int),
			"number_of_recent_passwords_disallowed": d.Get("number_of_recent_passwords_disallowed").(int),
			"password_not_username_or_invert":       d.Get("password_not_username_or_invert").(bool),
			"password_validity_period":              d.Get("password_validity_period").(int),
			"password_char_combination":             d.Get("password_char_combination").(int),
		},
	}
}

func resourceV3PasswordPolicyCreateOrUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg      = meta.(*config.Config)
		domainId = cfg.DomainID
		httpUrl  = "v3.0/OS-SECURITYPOLICY/domains/{domain_id}/password-policy"
	)

	client, err := cfg.NewServiceClient("iam", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	httpUrl = client.Endpoint + httpUrl
	httpUrl = strings.ReplaceAll(httpUrl, "{domain_id}", domainId)
	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         buildV3UpdatePasswordPolicyBodyParams(d),
	}

	_, err = client.Request("PUT", httpUrl, &updateOpt)
	if err != nil {
		return diag.Errorf("error updating the IAM account password policy: %s", err)
	}

	if d.IsNewResource() {
		d.SetId(domainId)
	}

	return resourceV3PasswordPolicyRead(ctx, d, meta)
}

func GetV3PasswordPolicy(client *golangsdk.ServiceClient, domainId string) (interface{}, error) {
	httpUrl := "v3.0/OS-SECURITYPOLICY/domains/{domain_id}/password-policy"
	httpUrl = strings.ReplaceAll(httpUrl, "{domain_id}", domainId)
	httpUrl = client.Endpoint + httpUrl
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	requestResp, err := client.Request("GET", httpUrl, &getOpt)
	if err != nil {
		return nil, err
	}
	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return nil, err
	}

	if utils.PathSearch("password_policy.maximum_consecutive_identical_chars", respBody, float64(0)).(float64) == 0 &&
		utils.PathSearch("password_policy.minimum_password_age", respBody, float64(0)).(float64) == 0 &&
		utils.PathSearch("password_policy.minimum_password_length", respBody, float64(0)).(float64) == 8 &&
		utils.PathSearch("password_policy.number_of_recent_passwords_disallowed", respBody, float64(0)).(float64) == 1 &&
		utils.PathSearch("password_policy.password_not_username_or_invert", respBody, false).(bool) &&
		utils.PathSearch("password_policy.password_validity_period", respBody, float64(0)).(float64) == 0 &&
		utils.PathSearch("password_policy.password_char_combination", respBody, float64(0)).(float64) == 2 {
		return nil, golangsdk.ErrDefault404{
			ErrUnexpectedResponseCode: golangsdk.ErrUnexpectedResponseCode{
				Method:    "GET",
				URL:       "/v3.0/OS-SECURITYPOLICY/domains/{domain_id}/password-policy",
				RequestId: "NONE",
				Body:      []byte("All configurations of password policy have been restored to the default value"),
			},
		}
	}
	return respBody, nil
}

func resourceV3PasswordPolicyRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg      = meta.(*config.Config)
		domainId = d.Id()
	)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("iam", region)
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	policy, err := GetV3PasswordPolicy(client, domainId)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving password policy")
	}

	mErr := multierror.Append(nil,
		d.Set("maximum_consecutive_identical_chars",
			int(utils.PathSearch("password_policy.maximum_consecutive_identical_chars", policy, float64(0)).(float64))),
		d.Set("minimum_password_age", int(utils.PathSearch("password_policy.minimum_password_age", policy, float64(0)).(float64))),
		d.Set("minimum_password_length", int(utils.PathSearch("password_policy.minimum_password_length", policy, float64(0)).(float64))),
		d.Set("number_of_recent_passwords_disallowed",
			int(utils.PathSearch("password_policy.number_of_recent_passwords_disallowed", policy, float64(0)).(float64))),
		d.Set("password_not_username_or_invert", utils.PathSearch("password_policy.password_not_username_or_invert", policy, false).(bool)),
		d.Set("password_validity_period", int(utils.PathSearch("password_policy.password_validity_period", policy, float64(0)).(float64))),
		d.Set("password_char_combination", int(utils.PathSearch("password_policy.password_char_combination", policy, float64(0)).(float64))),
		d.Set("maximum_password_length", int(utils.PathSearch("password_policy.maximum_password_length", policy, float64(0)).(float64))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceV3PasswordPolicyDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg      = meta.(*config.Config)
		region   = cfg.GetRegion(d)
		domainId = d.Id()
	)

	client, err := cfg.NewServiceClient("iam", region)
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	httpUrl := "v3.0/OS-SECURITYPOLICY/domains/{domain_id}/password-policy"
	restorePath := client.Endpoint + httpUrl
	restorePath = strings.ReplaceAll(restorePath, "{domain_id}", domainId)

	restoreOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody: map[string]interface{}{
			// Default configurations of password policy to be restored
			"password_policy": map[string]interface{}{
				"maximum_consecutive_identical_chars":   0,
				"minimum_password_age":                  0,
				"minimum_password_length":               8,
				"number_of_recent_passwords_disallowed": 1,
				"password_not_username_or_invert":       true,
				"password_validity_period":              0,
				"password_char_combination":             2,
			},
		},
	}

	_, err = client.Request("PUT", restorePath, &restoreOpt)
	if err != nil {
		return diag.Errorf("error restoring the default password policy: %s", err)
	}

	return nil
}
