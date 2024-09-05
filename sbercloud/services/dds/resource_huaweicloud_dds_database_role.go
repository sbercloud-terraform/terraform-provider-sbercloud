package dds

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/dds/v3/instances"
	"github.com/chnsz/golangsdk/openstack/dds/v3/roles"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

// @API DDS POST /v3/{project_id}/instances/{instance_id}/db-role
// @API DDS GET /v3/{project_id}/instances
// @API DDS GET /v3/{project_id}/instances/{instance_id}/db-roles
// @API DDS DELETE /v3/{project_id}/instances/{instance_id}/db-role
func ResourceDatabaseRole() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDatabaseRoleCreate,
		ReadContext:   resourceDatabaseRoleRead,
		DeleteContext: resourceDatabaseRoleDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(60 * time.Minute),
		},

		Importer: &schema.ResourceImporter{
			StateContext: resourceDatabaseRoleImportState,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"db_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"roles": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"db_name": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
					},
				},
			},
			"privileges": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     privilegeSchemaResource(),
			},
			"inherited_privileges": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     privilegeSchemaResource(),
			},
		},
	}
}

func privilegeSchemaResource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"resources": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"collection": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"db_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"actions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func buildDatabaseRoles(roleList []interface{}) []roles.Role {
	result := make([]roles.Role, len(roleList))
	for i, val := range roleList {
		role := val.(map[string]interface{})
		result[i] = roles.Role{
			DbName: role["db_name"].(string),
			Name:   role["name"].(string),
		}
	}
	return result
}

func instanceActionsRefreshFunc(client *golangsdk.ServiceClient, instanceId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		opts := instances.ListInstanceOpts{
			Id: instanceId,
		}
		resp, err := instances.List(client, opts).AllPages()
		if err != nil {
			return nil, "ERROR", fmt.Errorf("error getting DDS instance list")
		}
		result, err := instances.ExtractInstances(resp)
		if err != nil {
			return nil, "ERROR", fmt.Errorf("error extracting DDS instance")
		}
		if len(result.Instances) < 1 {
			return resp, "DELETED", nil
		}
		if len(result.Instances[0].Actions) > 0 {
			return resp, "PENDING", nil
		}
		return resp, "ACTIVE", nil
	}
}

func resourceDatabaseRoleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	client, err := conf.DdsV3Client(conf.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating DDS v3 client: %s", err)
	}

	instanceId := d.Get("instance_id").(string)
	dbName := d.Get("db_name").(string)
	roleName := d.Get("name").(string)
	opt := roles.CreateOpts{
		DbName: dbName,
		Name:   roleName,
		Roles:  buildDatabaseRoles(d.Get("roles").([]interface{})),
	}
	retryFunc := func() (interface{}, bool, error) {
		err = roles.Create(client, instanceId, opt)
		retry, err := handleMultiOperationsError(err)
		return nil, retry, err
	}
	_, err = common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     instanceActionsRefreshFunc(client, instanceId),
		WaitTarget:   []string{"ACTIVE"},
		WaitPending:  []string{"PENDING"},
		Timeout:      d.Timeout(schema.TimeoutCreate),
		DelayTimeout: 1 * time.Second,
		PollInterval: 10 * time.Second,
	})
	if err != nil {
		return diag.Errorf("error creating database role: %v", err)
	}

	// Role names are unique within a single database, but duplicate names may exist across multiple databases.
	d.SetId(fmt.Sprintf("%s/%s/%s", instanceId, dbName, roleName))
	return resourceDatabaseRoleRead(ctx, d, meta)
}

func flattenDatabaseRoles(roleList []roles.RoleDetail) []map[string]interface{} {
	if len(roleList) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, len(roleList))
	for i, role := range roleList {
		result[i] = map[string]interface{}{
			"db_name": role.DbName,
			"name":    role.Name,
		}
	}
	return result
}

func flattenDatabasePrivilegeResource(roleResource roles.Resource) []map[string]interface{} {
	if (roleResource == roles.Resource{}) {
		return nil
	}

	return []map[string]interface{}{
		{
			"collection": roleResource.Collection,
			"db_name":    roleResource.DbName,
		},
	}
}

func flattenDatabasePrivileges(privileges []roles.Privilege) (result []map[string]interface{}) {
	if len(privileges) < 1 {
		return nil
	}

	for _, privilege := range privileges {
		result = append(result, map[string]interface{}{
			"resources": flattenDatabasePrivilegeResource(privilege.Resource),
			"actions":   privilege.Actions,
		})
	}
	return
}

func resourceDatabaseRoleRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)
	client, err := conf.DdsV3Client(region)
	if err != nil {
		return diag.Errorf("error creating DDS v3 client: %s", err)
	}

	instanceId := d.Get("instance_id").(string)
	name := d.Get("name").(string)
	opts := roles.ListOpts{
		DbName: d.Get("db_name").(string),
		Name:   name,
	}
	resp, err := roles.List(client, instanceId, opts)
	if err != nil {
		return common.CheckDeletedDiag(d, err, fmt.Sprintf("error getting database role (%s) from DDS instance (%s)",
			name, instanceId))
	}
	if len(resp) < 1 {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "")
	}
	role := resp[0]
	log.Printf("[DEBUG] The role response is: %#v", role)

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("name", role.Name),
		d.Set("db_name", role.DbName),
		d.Set("roles", flattenDatabaseRoles(role.Roles)),
		d.Set("privileges", flattenDatabasePrivileges(role.Privileges)),
		d.Set("inherited_privileges", flattenDatabasePrivileges(role.InheritedPrivileges)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func databaseRoleRefreshFunc(client *golangsdk.ServiceClient, instanceId, dbName,
	roleName string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		opts := roles.ListOpts{
			DbName: dbName,
			Name:   roleName,
		}
		resp, err := roles.List(client, instanceId, opts)
		if err != nil {
			return nil, "ERROR", fmt.Errorf("error getting database role (%s) from DDS instance (%s)", roleName,
				instanceId)
		}
		if len(resp) < 1 {
			return resp, "DELETED", nil
		}
		return resp, "ACTIVE", nil
	}
}

func resourceDatabaseRoleDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	client, err := conf.DdsV3Client(conf.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating DDS v3 client: %s", err)
	}

	instanceId := d.Get("instance_id").(string)
	name := d.Get("name").(string)
	dbName := d.Get("db_name").(string)
	opts := roles.DeleteOpts{
		DbName: dbName,
		Name:   name,
	}
	retryFunc := func() (interface{}, bool, error) {
		err = roles.Delete(client, instanceId, opts)
		retry, err := handleMultiOperationsError(err)
		return nil, retry, err
	}
	_, err = common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     instanceActionsRefreshFunc(client, instanceId),
		WaitTarget:   []string{"ACTIVE"},
		WaitPending:  []string{"PENDING"},
		Timeout:      d.Timeout(schema.TimeoutCreate),
		DelayTimeout: 1 * time.Second,
		PollInterval: 10 * time.Second,
	})
	if err != nil {
		return diag.Errorf("error deleting database role (%s) from DDS instance (%s): %v", name, instanceId, err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"ACTIVE"},
		Target:       []string{"DELETED"},
		Refresh:      databaseRoleRefreshFunc(client, instanceId, dbName, name),
		Timeout:      d.Timeout(schema.TimeoutDelete),
		Delay:        2 * time.Second,
		PollInterval: 10 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for database role deletion to complete: %s", err)
	}

	return nil
}

func resourceDatabaseRoleImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid format specified for import ID, must be <instance_id>/<db_name>/<name>")
	}

	d.Set("instance_id", parts[0])
	d.Set("db_name", parts[1])
	d.Set("name", parts[2])
	return []*schema.ResourceData{d}, nil
}
