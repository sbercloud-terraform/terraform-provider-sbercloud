package dli

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/dli/v1/tables"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func ResourceDliTable() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDliTableCreate,
		ReadContext:   resourceDliTableRead,
		DeleteContext: resourceDliTableDelete,
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
			"database_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"data_location": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"columns": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"type": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"description": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"is_partition": {
							Type:     schema.TypeBool,
							Optional: true,
							ForceNew: true,
						},
					},
				},
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"data_format": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{"parquet", "orc", "csv", "json", "carbon", "avro"},
					true),
				Computed: true,
			},
			"bucket_location": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"with_column_header": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"delimiter": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"quote_char": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"escape_char": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"date_format": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"timestamp_format": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
		},

		Timeouts: &schema.ResourceTimeout{
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
	}
}

func resourceDliTableCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.DliV1Client(region)
	if err != nil {
		return diag.Errorf("error creating DLI v1 client, err=%s", err)
	}
	databaseName := d.Get("database_name").(string)
	tableName := d.Get("name").(string)
	opts := tables.CreateTableOpts{
		TableName:       tableName,
		DataLocation:    d.Get("data_location").(string),
		Columns:         buildColumnParam(d),
		Description:     d.Get("description").(string),
		DataType:        d.Get("data_format").(string),
		DataPath:        d.Get("bucket_location").(string),
		Delimiter:       d.Get("delimiter").(string),
		QuoteChar:       d.Get("quote_char").(string),
		EscapeChar:      d.Get("escape_char").(string),
		DateFormat:      d.Get("date_format").(string),
		TimestampFormat: d.Get("timestamp_format").(string),
	}

	if v, ok := d.GetOk("with_column_header"); ok {
		opts.WithColumnHeader = utils.Bool(v.(bool))
	}

	log.Printf("[DEBUG] Creating new DLI table opts: %#v", opts)

	rst, createErr := tables.Create(client, databaseName, opts)
	if createErr != nil {
		return diag.Errorf("error creating DLI table: %s", createErr)
	}

	if rst != nil && !rst.IsSuccess {
		return diag.Errorf("error creating DLI table: %s", rst.Message)
	}

	d.SetId(fmt.Sprintf("%s/%s", databaseName, tableName))
	return resourceDliTableRead(ctx, d, meta)
}

func resourceDliTableRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.DliV1Client(region)
	if err != nil {
		return diag.Errorf("error creating DLI v1 client, err=%s", err)
	}

	databaseName, tableName := ParseTableInfoFromId(d.Id())

	detail, err := tables.Get(client, databaseName, tableName)
	if err != nil {
		return common.CheckDeletedDiag(d, parseDliErrorToError404(err), "DLI table")
	}

	if detail != nil && !detail.IsSuccess {
		return diag.Errorf("error query DLI Table: %s", detail.Message)
	}

	tbList, err := tables.List(client, databaseName, tables.ListOpts{
		Keyword:    tableName,
		WithDetail: utils.Bool(true),
		WithPriv:   utils.Bool(true),
	})
	if err != nil {
		return diag.Errorf("error query DLI Table %q:%s", d.Id(), err)
	}

	if tbList != nil && !tbList.IsSuccess {
		return diag.Errorf("error query DLI Table: %s", tbList.Message)
	}

	tb, err := filterByTableName(tbList.Tables, tableName)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "DLI table")
	}

	mErr := multierror.Append(
		d.Set("database_name", databaseName),
		d.Set("name", tableName),
		d.Set("data_location", tb.DataLocation),
		setColumnsToState(d, detail.Columns),
		d.Set("description", detail.TableComment),
		d.Set("data_format", tb.DataType),
		d.Set("bucket_location", tb.Location),
		setStoragePropertiesToState(d, detail.StorageProperties),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func setColumnsToState(d *schema.ResourceData, columns []tables.Column) error {
	if len(columns) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(columns))
	for _, column := range columns {
		item := map[string]interface{}{
			"name":         column.ColumnName,
			"type":         column.Type,
			"description":  column.Description,
			"is_partition": column.IsPartitionColumn,
		}
		result = append(result, item)
	}

	return d.Set("columns", result)
}

func setStoragePropertiesToState(d *schema.ResourceData, storageProperties []map[string]interface{}) error {
	if len(storageProperties) == 0 {
		return nil
	}
	var mErr *multierror.Error
	for _, properties := range storageProperties {
		switch properties["key"] {
		case "delimiter":
			mErr = multierror.Append(d.Set("delimiter", properties["value"]))
		case "quote":
			mErr = multierror.Append(d.Set("quote_char", properties["value"]))
		case "escape":
			mErr = multierror.Append(d.Set("escape_char", properties["value"]))
		case "dateformat":
			mErr = multierror.Append(d.Set("date_format", properties["value"]))
		case "timestampformat":
			mErr = multierror.Append(d.Set("timestamp_format", properties["value"]))
		case "header":
			mErr = multierror.Append(d.Set("with_column_header", properties["value"].(string) == "true"))
		}
	}
	return mErr.ErrorOrNil()
}

func resourceDliTableDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.DliV1Client(region)
	if err != nil {
		return diag.Errorf("error creating DLI v1 client, err=%s", err)
	}

	databaseName, tableName := ParseTableInfoFromId(d.Id())

	resp, dErr := tables.Delete(client, databaseName, tableName, false)
	if dErr != nil {
		return diag.Errorf("error delete DLI Table %q:%s", d.Id(), dErr)
	}

	if resp != nil && !resp.IsSuccess {
		return diag.Errorf("error delete DLI Table: %s", resp.Message)
	}

	return nil
}

func buildColumnParam(d *schema.ResourceData) []tables.ColumnOpts {
	var rt []tables.ColumnOpts
	columns := d.Get("columns").([]interface{})
	if len(columns) > 0 {
		for _, raw := range columns {
			columnRaw := raw.(map[string]interface{})
			column := tables.ColumnOpts{
				ColumnName:        columnRaw["name"].(string),
				Type:              columnRaw["type"].(string),
				Description:       columnRaw["description"].(string),
				IsPartitionColumn: utils.Bool(columnRaw["is_partition"].(bool)),
			}
			rt = append(rt, column)
		}
	}

	return rt
}

func ParseTableInfoFromId(id string) (databaseName, tableName string) {
	idArrays := strings.Split(id, "/")
	databaseName = idArrays[0]
	tableName = idArrays[1]
	return
}

func filterByTableName(tablesResp []tables.Table4List, tableName string) (*tables.Table4List, error) {
	log.Printf("[DEBUG]The list of table from SDK:%+v", tablesResp)
	for _, v := range tablesResp {
		if v.TableName == tableName {
			return &v, nil
		}
	}
	return &tables.Table4List{}, golangsdk.ErrDefault404{}
}
