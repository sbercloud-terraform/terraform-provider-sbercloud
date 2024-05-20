package dli

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DLI POST /v3/{project_id}/datasource/auth-infos
// @API DLI GET /v3/{project_id}/datasource/auth-infos
// @API DLI PUT /v3/{project_id}/datasource/auth-infos
// @API DLI DELETE /v3/{project_id}/datasource/auth-infos/{auth_info_name}
func ResourceDatasourceAuth() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDatasourceAuthCreate,
		UpdateContext: resourceDatasourceAuthUpdate,
		ReadContext:   resourceDatasourceAuthRead,
		DeleteContext: resourceDatasourceAuthDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The name of a datasource authentication.`,
				ValidateFunc: validation.All(
					validation.StringMatch(regexp.MustCompile(`^[A-Za-z0-9][\w]*$`),
						"Only letters, digits and underscores (_) are allowed."),
					validation.StringDoesNotMatch(regexp.MustCompile(`^[0-9]*$`), "The name cannot be all digits."),
				),
			},
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Data source type.`,
			},
			"username": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Username for accessing the security cluster or datasource.`,
				ConflictsWith: []string{
					"truststore_location",
				},
			},
			"password": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Sensitive:   true,
				Description: `The password for accessing the security cluster or datasource.`,
				RequiredWith: []string{
					"username",
				},
			},
			"certificate_location": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The OBS path of the security cluster certificate.`,
				ConflictsWith: []string{
					"truststore_location", "krb5_conf",
				},
			},
			"truststore_location": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The OBS path of the **truststore** configuration file.`,
				RequiredWith: []string{
					"truststore_password", "keystore_location", "keystore_password", "key_password",
				},
			},
			"truststore_password": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Computed:    true,
				Description: `The password of the **truststore** configuration file.`,
			},
			"keystore_location": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The OBS path of the **keystore** configuration file.`,
			},
			"keystore_password": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Computed:    true,
				Description: `The password of the **keystore ** configuration file.`,
			},
			"key_password": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				ForceNew:    true,
				Computed:    true,
				Description: `The key password.`,
			},
			"krb5_conf": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The OBS path of the **krb5** configuration file.`,
				RequiredWith: []string{
					"keytab",
				},
			},
			"keytab": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The OBS path of the **keytab** configuration file.`,
			},
			"owner": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The user name of owner.`,
			},
		},
	}
}

func resourceDatasourceAuthCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// createDatasourceAuth: create a DLI datasource authentication.
	var (
		createDatasourceAuthHttpUrl = "v3/{project_id}/datasource/auth-infos"
		createDatasourceAuthProduct = "dli"
	)
	createDatasourceAuthClient, err := cfg.NewServiceClient(createDatasourceAuthProduct, region)
	if err != nil {
		return diag.Errorf("error creating DLI Client: %s", err)
	}

	createDatasourceAuthPath := createDatasourceAuthClient.Endpoint + createDatasourceAuthHttpUrl
	createDatasourceAuthPath = strings.ReplaceAll(createDatasourceAuthPath, "{project_id}", createDatasourceAuthClient.ProjectID)

	createDatasourceAuthOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			201,
		},
	}
	createDatasourceAuthOpt.JSONBody = utils.RemoveNil(buildCreateDatasourceAuthBodyParams(d, cfg))
	_, err = createDatasourceAuthClient.Request("POST", createDatasourceAuthPath, &createDatasourceAuthOpt)
	if err != nil {
		return diag.Errorf("error creating DatasourceAuth: %s", err)
	}

	d.SetId(d.Get("name").(string))

	return resourceDatasourceAuthRead(ctx, d, meta)
}

func buildCreateDatasourceAuthBodyParams(d *schema.ResourceData, _ *config.Config) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"auth_info_name":       utils.ValueIngoreEmpty(d.Get("name")),
		"datasource_type":      utils.ValueIngoreEmpty(d.Get("type")),
		"username":             utils.ValueIngoreEmpty(d.Get("username")),
		"password":             utils.ValueIngoreEmpty(d.Get("password")),
		"certificate_location": utils.ValueIngoreEmpty(d.Get("certificate_location")),
		"truststore_location":  utils.ValueIngoreEmpty(d.Get("truststore_location")),
		"truststore_password":  utils.ValueIngoreEmpty(d.Get("truststore_password")),
		"keystore_location":    utils.ValueIngoreEmpty(d.Get("keystore_location")),
		"keystore_password":    utils.ValueIngoreEmpty(d.Get("keystore_password")),
		"key_password":         utils.ValueIngoreEmpty(d.Get("key_password")),
		"krb5_conf":            utils.ValueIngoreEmpty(d.Get("krb5_conf")),
		"keytab":               utils.ValueIngoreEmpty(d.Get("keytab")),
	}
	return bodyParams
}

func resourceDatasourceAuthRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// getDatasourceAuth: Query the DLI datasource authentication.
	var (
		getDatasourceAuthHttpUrl = "v3/{project_id}/datasource/auth-infos"
		getDatasourceAuthProduct = "dli"
	)
	getDatasourceAuthClient, err := cfg.NewServiceClient(getDatasourceAuthProduct, region)
	if err != nil {
		return diag.Errorf("error creating DLI Client: %s", err)
	}

	getDatasourceAuthPath := getDatasourceAuthClient.Endpoint + getDatasourceAuthHttpUrl
	getDatasourceAuthPath = strings.ReplaceAll(getDatasourceAuthPath, "{project_id}", getDatasourceAuthClient.ProjectID)

	getDatasourceAuthqueryParams := buildGetDatasourceAuthQueryParams(d)
	getDatasourceAuthPath += getDatasourceAuthqueryParams

	getDatasourceAuthOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	getDatasourceAuthResp, err := getDatasourceAuthClient.Request("GET", getDatasourceAuthPath, &getDatasourceAuthOpt)

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving DatasourceAuth")
	}

	getDatasourceAuthRespBody, err := utils.FlattenResponse(getDatasourceAuthResp)
	if err != nil {
		return diag.FromErr(err)
	}

	v := utils.PathSearch("auth_infos[0]", getDatasourceAuthRespBody, nil)
	if v == nil {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "error retrieving DatasourceAuth")
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("name", utils.PathSearch("auth_infos[0].auth_info_name", getDatasourceAuthRespBody, nil)),
		d.Set("type", utils.PathSearch("auth_infos[0].datasource_type", getDatasourceAuthRespBody, nil)),
		d.Set("username", utils.PathSearch("auth_infos[0].user_name", getDatasourceAuthRespBody, nil)),
		d.Set("certificate_location", utils.PathSearch("auth_infos[0].certificate_location", getDatasourceAuthRespBody, nil)),
		d.Set("truststore_location", utils.PathSearch("auth_infos[0].truststore_location", getDatasourceAuthRespBody, nil)),
		d.Set("keystore_location", utils.PathSearch("auth_infos[0].keystore_location", getDatasourceAuthRespBody, nil)),
		d.Set("krb5_conf", utils.PathSearch("auth_infos[0].krb5_conf", getDatasourceAuthRespBody, nil)),
		d.Set("keytab", utils.PathSearch("auth_infos[0].keytab", getDatasourceAuthRespBody, nil)),
		d.Set("owner", utils.PathSearch("auth_infos[0].owner", getDatasourceAuthRespBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildGetDatasourceAuthQueryParams(d *schema.ResourceData) string {
	res := ""
	res = fmt.Sprintf("%s&auth_info_name=%v", res, d.Id())

	if res != "" {
		res = "?" + res[1:]
	}
	return res
}

func resourceDatasourceAuthUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	updateDatasourceAuthChanges := []string{
		"name",
		"username",
		"password",
		"truststore_location",
		"truststore_password",
		"keystore_location",
		"keystore_password",
		"krb5_conf",
		"keytab",
	}

	if d.HasChanges(updateDatasourceAuthChanges...) {
		var (
			updateDatasourceAuthHttpUrl = "v3/{project_id}/datasource/auth-infos"
			updateDatasourceAuthProduct = "dli"
		)
		updateDatasourceAuthClient, err := cfg.NewServiceClient(updateDatasourceAuthProduct, region)
		if err != nil {
			return diag.Errorf("error creating DLI Client: %s", err)
		}

		updateDatasourceAuthPath := updateDatasourceAuthClient.Endpoint + updateDatasourceAuthHttpUrl
		updateDatasourceAuthPath = strings.ReplaceAll(updateDatasourceAuthPath, "{project_id}", updateDatasourceAuthClient.ProjectID)

		updateDatasourceAuthOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			OkCodes: []int{
				200,
			},
		}
		updateDatasourceAuthOpt.JSONBody = utils.RemoveNil(buildUpdateDatasourceAuthBodyParams(d, cfg))
		_, err = updateDatasourceAuthClient.Request("PUT", updateDatasourceAuthPath, &updateDatasourceAuthOpt)
		if err != nil {
			return diag.Errorf("error updating DatasourceAuth: %s", err)
		}
	}
	return resourceDatasourceAuthRead(ctx, d, meta)
}

func buildUpdateDatasourceAuthBodyParams(d *schema.ResourceData, _ *config.Config) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"auth_info_name":      utils.ValueIngoreEmpty(d.Get("name")),
		"username":            utils.ValueIngoreEmpty(d.Get("username")),
		"password":            utils.ValueIngoreEmpty(d.Get("password")),
		"truststore_location": utils.ValueIngoreEmpty(d.Get("truststore_location")),
		"truststore_password": utils.ValueIngoreEmpty(d.Get("truststore_password")),
		"keystore_location":   utils.ValueIngoreEmpty(d.Get("keystore_location")),
		"keystore_password":   utils.ValueIngoreEmpty(d.Get("keystore_password")),
		"krb5_conf":           utils.ValueIngoreEmpty(d.Get("krb5_conf")),
		"keytab":              utils.ValueIngoreEmpty(d.Get("keytab")),
	}
	return bodyParams
}

func resourceDatasourceAuthDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// deleteDatasourceAuth: missing operation notes
	var (
		deleteDatasourceAuthHttpUrl = "v3/{project_id}/datasource/auth-infos/{id}"
		deleteDatasourceAuthProduct = "dli"
	)
	deleteDatasourceAuthClient, err := cfg.NewServiceClient(deleteDatasourceAuthProduct, region)
	if err != nil {
		return diag.Errorf("error creating DLI Client: %s", err)
	}

	deleteDatasourceAuthPath := deleteDatasourceAuthClient.Endpoint + deleteDatasourceAuthHttpUrl
	deleteDatasourceAuthPath = strings.ReplaceAll(deleteDatasourceAuthPath, "{project_id}", deleteDatasourceAuthClient.ProjectID)
	deleteDatasourceAuthPath = strings.ReplaceAll(deleteDatasourceAuthPath, "{id}", d.Id())

	deleteDatasourceAuthOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	_, err = deleteDatasourceAuthClient.Request("DELETE", deleteDatasourceAuthPath, &deleteDatasourceAuthOpt)
	if err != nil {
		return diag.Errorf("error deleting DatasourceAuth: %s", err)
	}

	return nil
}
