package cdn

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/cdn/v1/domains"

	cdnv1 "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/cdn/v1"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/cdn/v1/model"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var httpsConfig = schema.Schema{
	Type:     schema.TypeList,
	Optional: true,
	Computed: true,
	MaxItems: 1,
	Elem: &schema.Resource{
		Schema: map[string]*schema.Schema{
			"https_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"certificate_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"certificate_body": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"private_key": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
				Computed:  true,
			},
			"certificate_source": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.IntInSlice([]int{
					0, 1,
				}),
			},
			"http2_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"tls_version": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"https_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"http2_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	},
}

var requestAndResponseHeader = schema.Schema{
	Type:     schema.TypeList,
	Optional: true,
	Computed: true,
	Elem: &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"action": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"set", "delete",
				}, false),
			},
			"value": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	},
}

var authOpts = schema.Schema{
	Type:     schema.TypeList,
	Optional: true,
	Computed: true,
	MaxItems: 1,
	Elem: &schema.Resource{
		Schema: map[string]*schema.Schema{
			"enabled": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"type_a", "type_b", "type_c1", "type_c2",
				}, false),
			},
			"key": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"time_format": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"dec", "hex",
				}, false),
			},
			"expire_time": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	},
}

var forceRedirectAndCompress = schema.Schema{
	Type:     schema.TypeList,
	Optional: true,
	Computed: true,
	MaxItems: 1,
	Elem: &schema.Resource{
		Schema: map[string]*schema.Schema{
			"enabled": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	},
}

var cacheUrlParameterFilter = schema.Schema{
	Type:     schema.TypeList,
	Optional: true,
	Computed: true,
	MaxItems: 1,
	Elem: &schema.Resource{
		Schema: map[string]*schema.Schema{
			"type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"full_url", "ignore_url_params", "del_args", "reserve_args",
				}, false),
			},
			"value": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	},
}

func ResourceCdnDomainV1() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCdnDomainV1Create,
		ReadContext:   resourceCdnDomainV1Read,
		UpdateContext: resourceCdnDomainV1Update,
		DeleteContext: resourceCdnDomainV1Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

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
				ValidateFunc: validation.StringInSlice([]string{
					"web", "download", "video", "wholeSite",
				}, true),
			},
			"sources": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"origin": {
							Type:     schema.TypeString,
							Required: true,
						},
						"origin_type": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"ipaddr", "domain", "obs_bucket",
							}, true),
						},
						"active": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  1,
						},
						"obs_web_hosting_enabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"http_port": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"https_port": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"retrieval_host": {
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
				ForceNew: true,
			},
			"service_area": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"configs": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"origin_protocol": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ValidateFunc: validation.StringInSlice([]string{
								"follow", "http", "https",
							}, false),
						},
						"ipv6_enable": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"range_based_retrieval_enabled": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"https_settings":             &httpsConfig,
						"retrieval_request_header":   &requestAndResponseHeader,
						"http_response_header":       &requestAndResponseHeader,
						"url_signing":                &authOpts,
						"force_redirect":             &forceRedirectAndCompress,
						"compress":                   &forceRedirectAndCompress,
						"cache_url_parameter_filter": &cacheUrlParameterFilter,
					},
				},
			},

			"cache_settings": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"follow_origin": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"rules": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"rule_type": {
										Type:     schema.TypeString,
										Required: true,
										DiffSuppressFunc: func(k, o, n string, d *schema.ResourceData) bool {
											// Convert several original types and change parameter types while ensuring
											// that the original configuration is available.
											// Notes: the state file no longer save the original types.
											return parseCacheRuleType(n) == o
										},
									},
									"content": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"ttl": {
										Type:        schema.TypeInt,
										Optional:    true,
										Computed:    true,
										Description: "schema: Required",
									},
									"ttl_type": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "schema: Required",
										DiffSuppressFunc: func(k, o, n string, d *schema.ResourceData) bool {
											// Convert several original types and change parameter types while ensuring
											// that the original configuration is available.
											// Notes: the state file no longer save the original types.
											return parseCacheTTLUnits(n) == o
										},
									},
									"priority": {
										Type:        schema.TypeInt,
										Optional:    true,
										Computed:    true,
										Description: "schema: Required",
									},
								},
							},
						},
					},
				},
			},
			"tags": common.TagsSchema(),
			"cname": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"domain_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

type WaitDomainStatus struct {
	ID      string
	Penging []string
	Target  []string
	Opts    *domains.ExtensionOpts
}

func getDomainSources(d *schema.ResourceData) []domains.SourcesOpts {
	var sourceRequests []domains.SourcesOpts

	sources := d.Get("sources").([]interface{})
	for i := range sources {
		source := sources[i].(map[string]interface{})
		sourceRequest := domains.SourcesOpts{
			IporDomain:    source["origin"].(string),
			OriginType:    source["origin_type"].(string),
			ActiveStandby: source["active"].(int),
		}
		sourceRequests = append(sourceRequests, sourceRequest)
	}
	return sourceRequests
}

func buildHTTPSOpts(rawHTTPS []interface{}) *model.HttpPutBody {
	if len(rawHTTPS) != 1 {
		return nil
	}

	https := rawHTTPS[0].(map[string]interface{})
	httpsStatus := ""
	if https["https_enabled"].(bool) {
		httpsStatus = "on"
	}
	http2Status := ""
	if https["http2_enabled"].(bool) {
		http2Status = "on"
	}

	httpsOpts := model.HttpPutBody{
		HttpsStatus:       utils.StringIgnoreEmpty(httpsStatus),
		CertificateName:   utils.StringIgnoreEmpty(https["certificate_name"].(string)),
		CertificateValue:  utils.StringIgnoreEmpty(https["certificate_body"].(string)),
		PrivateKey:        utils.StringIgnoreEmpty(https["private_key"].(string)),
		CertificateSource: utils.Int32IgnoreEmpty(int32(https["certificate_source"].(int))),
		Http2Status:       utils.StringIgnoreEmpty(http2Status),
		TlsVersion:        utils.StringIgnoreEmpty(https["tls_version"].(string)),
	}

	return &httpsOpts
}

func buildOriginRequestHeaderOpts(rawOriginRequestHeader []interface{}) *[]model.OriginRequestHeader {
	if len(rawOriginRequestHeader) < 1 {
		return nil
	}

	originRequestHeaderOpts := make([]model.OriginRequestHeader, len(rawOriginRequestHeader))
	for i, v := range rawOriginRequestHeader {
		header := v.(map[string]interface{})
		originRequestHeaderOpts[i] = model.OriginRequestHeader{
			Name:   header["name"].(string),
			Value:  utils.StringIgnoreEmpty(header["value"].(string)),
			Action: header["action"].(string),
		}
	}

	return &originRequestHeaderOpts
}

func buildHttpResponseHeaderOpts(rawHttpResponseHeader []interface{}) *[]model.HttpResponseHeader {
	if len(rawHttpResponseHeader) < 1 {
		return nil
	}

	httpResponseHeaderOpts := make([]model.HttpResponseHeader, len(rawHttpResponseHeader))
	for i, v := range rawHttpResponseHeader {
		header := v.(map[string]interface{})
		httpResponseHeaderOpts[i] = model.HttpResponseHeader{
			Name:   header["name"].(string),
			Value:  utils.StringIgnoreEmpty(header["value"].(string)),
			Action: header["action"].(string),
		}
	}

	return &httpResponseHeaderOpts
}

func parseFunctionEnabledStatus(enabled bool) string {
	if enabled {
		return "on"
	}
	return "off"
}

func buildUrlAuthOpts(rawUrlAuth []interface{}) *model.UrlAuth {
	if len(rawUrlAuth) != 1 {
		return nil
	}

	urlAuth := rawUrlAuth[0].(map[string]interface{})
	urlAuthOpts := model.UrlAuth{
		Status:     parseFunctionEnabledStatus(urlAuth["enabled"].(bool)),
		Type:       utils.StringIgnoreEmpty(urlAuth["type"].(string)),
		Key:        utils.StringIgnoreEmpty(urlAuth["key"].(string)),
		TimeFormat: utils.StringIgnoreEmpty(urlAuth["time_format"].(string)),
		ExpireTime: utils.Int32IgnoreEmpty(int32(urlAuth["expire_time"].(int))),
	}

	return &urlAuthOpts
}

func buildForceRedirectOpts(rawForceRedirect []interface{}) *model.ForceRedirectConfig {
	if len(rawForceRedirect) != 1 {
		return nil
	}

	forceRedirect := rawForceRedirect[0].(map[string]interface{})
	forceRedirectOpts := model.ForceRedirectConfig{
		Status: parseFunctionEnabledStatus(forceRedirect["enabled"].(bool)),
		Type:   utils.StringIgnoreEmpty(forceRedirect["type"].(string)),
	}

	return &forceRedirectOpts
}

func buildCompressOpts(rawCompress []interface{}) *model.Compress {
	if len(rawCompress) != 1 {
		return nil
	}

	compress := rawCompress[0].(map[string]interface{})
	compressOpts := model.Compress{
		Status: parseFunctionEnabledStatus(compress["enabled"].(bool)),
		Type:   utils.StringIgnoreEmpty(compress["type"].(string)),
	}

	return &compressOpts
}

func buildCacheUrlParameterFilterOpts(rawCacheUrlParameterFilter []interface{}) *model.CacheUrlParameterFilter {
	if len(rawCacheUrlParameterFilter) != 1 {
		return nil
	}

	cacheUrlParameterFilter := rawCacheUrlParameterFilter[0].(map[string]interface{})
	cacheUrlParameterFilterOpts := model.CacheUrlParameterFilter{
		Value: utils.StringIgnoreEmpty(cacheUrlParameterFilter["value"].(string)),
		Type:  utils.StringIgnoreEmpty(cacheUrlParameterFilter["type"].(string)),
	}

	return &cacheUrlParameterFilterOpts
}

func buildSourcesOpts(rawSources []interface{}) *[]model.SourcesConfig {
	if len(rawSources) < 1 {
		return nil
	}
	sourcesOpts := make([]model.SourcesConfig, len(rawSources))
	for i, v := range rawSources {
		source := v.(map[string]interface{})
		var priority int32
		if source["active"].(int) == 1 {
			priority = 70
		} else {
			priority = 30
		}
		sourcesOpts[i] = model.SourcesConfig{
			OriginAddr:          source["origin"].(string),
			OriginType:          source["origin_type"].(string),
			Priority:            priority,
			ObsWebHostingStatus: utils.String(parseFunctionEnabledStatus(source["obs_web_hosting_enabled"].(bool))),
			HttpPort:            utils.Int32IgnoreEmpty(int32(source["http_port"].(int))),
			HttpsPort:           utils.Int32IgnoreEmpty(int32(source["https_port"].(int))),
			HostName:            utils.StringIgnoreEmpty(source["retrieval_host"].(string)),
		}
	}
	return &sourcesOpts
}

func parseCacheRuleType(ruleType string) string {
	var cacheRuleTypes = map[string]string{
		"0": "all",
		"1": "file_extension",
		"2": "catalog",
		"3": "full_path",
		"5": "home_page",
	}
	if val, ok := cacheRuleTypes[ruleType]; ok {
		return val
	}
	return ruleType
}

func parseCacheTTLUnits(ttlUnit string) string {
	var cacheTTLUnits = map[string]string{
		"1": "s",
		"2": "m",
		"3": "h",
		"4": "d",
	}
	if val, ok := cacheTTLUnits[ttlUnit]; ok {
		return val
	}
	return ttlUnit
}

func buildCacheRules(followOrigin bool, rules []interface{}) *[]model.CacheRules {
	result := make([]model.CacheRules, len(rules))
	for i, val := range rules {
		rule := val.(map[string]interface{})
		result[i] = model.CacheRules{
			FollowOrigin: parseFunctionEnabledStatus(followOrigin),
			MatchType:    parseCacheRuleType(rule["rule_type"].(string)),
			MatchValue:   utils.StringIgnoreEmpty(rule["content"].(string)),
			Ttl:          int32(rule["ttl"].(int)),
			TtlUnit:      parseCacheTTLUnits(rule["ttl_type"].(string)),
			Priority:     int32(rule["priority"].(int)),
		}
	}
	return &result
}

func updateDomainFullConfigs(client *cdnv1.CdnClient, cfg *config.Config, d *schema.ResourceData) error {
	rawConfigs := d.Get("configs").([]interface{})
	if len(rawConfigs) < 1 || rawConfigs[0] == nil {
		return nil
	}
	configs := rawConfigs[0].(map[string]interface{})

	ipv6Accelerate := 0
	if configs["ipv6_enable"].(bool) {
		ipv6Accelerate = 1
	}
	configsOpts := model.Configs{
		Sources:           buildSourcesOpts(d.Get("sources").([]interface{})),
		Ipv6Accelerate:    utils.Int32(int32(ipv6Accelerate)),
		OriginRangeStatus: utils.String(parseFunctionEnabledStatus(configs["range_based_retrieval_enabled"].(bool))),
	}
	if d.HasChange("configs.0.https_settings") {
		configsOpts.Https = buildHTTPSOpts(configs["https_settings"].([]interface{}))
	}
	if d.HasChange("configs.0.retrieval_request_header") {
		configsOpts.OriginRequestHeader = buildOriginRequestHeaderOpts(configs["retrieval_request_header"].([]interface{}))
	}
	if d.HasChange("configs.0.http_response_header") {
		configsOpts.HttpResponseHeader = buildHttpResponseHeaderOpts(configs["http_response_header"].([]interface{}))
	}
	if d.HasChange("configs.0.url_signing") {
		configsOpts.UrlAuth = buildUrlAuthOpts(configs["url_signing"].([]interface{}))
	}
	if d.HasChange("configs.0.origin_protocol") {
		configsOpts.OriginProtocol = utils.StringIgnoreEmpty(configs["origin_protocol"].(string))
	}
	if d.HasChange("configs.0.force_redirect") {
		configsOpts.ForceRedirect = buildForceRedirectOpts(configs["force_redirect"].([]interface{}))
	}
	if d.HasChange("configs.0.compress") {
		configsOpts.Compress = buildCompressOpts(configs["compress"].([]interface{}))
	}
	if d.HasChange("configs.0.cache_url_parameter_filter") {
		configsOpts.CacheUrlParameterFilter = buildCacheUrlParameterFilterOpts(configs["cache_url_parameter_filter"].([]interface{}))
	}

	if d.HasChange("cache_settings") {
		cacheSettings := d.Get("cache_settings").([]interface{})
		if len(cacheSettings) > 0 {
			cacheSetting := cacheSettings[0].(map[string]interface{})
			configsOpts.CacheRules = buildCacheRules(cacheSetting["follow_origin"].(bool), cacheSetting["rules"].([]interface{}))
		}
	}

	req := model.UpdateDomainFullConfigRequest{
		DomainName:          d.Get("name").(string),
		EnterpriseProjectId: utils.StringIgnoreEmpty(cfg.GetEnterpriseProjectID(d)),
		Body: &model.ModifyDomainConfigRequestBody{
			Configs: &configsOpts,
		},
	}

	_, err := client.UpdateDomainFullConfig(&req)
	if err != nil {
		return err
	}
	return nil
}

func resourceCdnDomainV1Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	cdnClient, err := cfg.CdnV1Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating CDN v1 client: %s", err)
	}

	createOpts := &domains.CreateOpts{
		DomainName:          d.Get("name").(string),
		BusinessType:        d.Get("type").(string),
		Sources:             getDomainSources(d),
		ServiceArea:         d.Get("service_area").(string),
		EnterpriseProjectId: cfg.GetEnterpriseProjectID(d),
	}

	log.Printf("[DEBUG] Create Options: %#v", createOpts)
	v, err := domains.Create(cdnClient, createOpts).Extract()
	if err != nil {
		return diag.Errorf("error creating CDN Domain: %s", err)
	}
	d.SetId(v.ID)

	// Wait for CDN domain to become active again before continuing
	opts := getResourceExtensionOpts(d, cfg)
	timeout := d.Timeout(schema.TimeoutCreate)
	log.Printf("[INFO] Waiting for CDN domain %s to become online.", v.ID)
	err = waitDomainOnline(ctx, cdnClient, v.ID, opts, timeout)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceCdnDomainV1Update(ctx, d, meta)
}

func waitforCDNV1DomainStatus(ctx context.Context, c *golangsdk.ServiceClient,
	waitstatus *WaitDomainStatus, timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending:    waitstatus.Penging,
		Target:     waitstatus.Target,
		Refresh:    resourceCDNV1DomainRefreshFunc(c, waitstatus.ID, waitstatus.Opts),
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("error waiting for CDN domain %s to become %s: %s",
			waitstatus.ID, waitstatus.Target, err)
	}
	return nil
}

func resourceCDNV1DomainRefreshFunc(c *golangsdk.ServiceClient, id string, opts *domains.ExtensionOpts) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		domain, err := domains.Get(c, id, opts).Extract()
		if err != nil {
			return nil, "", err
		}

		// return DomainStatus attribute of CDN domain resource
		return domain, domain.DomainStatus, nil
	}
}

func analyseFunctionEnabledStatus(enabledStatus string) bool {
	return enabledStatus == "on"
}

func analyseFunctionEnabledStatusPtr(enabledStatus *string) bool {
	return enabledStatus != nil && *enabledStatus == "on"
}

func flattenHTTPSAttrs(https *model.HttpGetBody, privateKey string) []map[string]interface{} {
	if https == nil {
		return nil
	}
	httpsAttrs := map[string]interface{}{
		"https_status":       https.HttpsStatus,
		"certificate_name":   https.CertificateName,
		"certificate_body":   https.CertificateValue,
		"private_key":        privateKey,
		"certificate_source": https.CertificateSource,
		"http2_status":       https.Http2Status,
		"tls_version":        https.TlsVersion,
		"https_enabled":      analyseFunctionEnabledStatusPtr(https.HttpsStatus),
		"http2_enabled":      analyseFunctionEnabledStatusPtr(https.Http2Status),
	}

	return []map[string]interface{}{httpsAttrs}
}

func flattenOriginRequestHeaderAttrs(originRequestHeader *[]model.OriginRequestHeader) []map[string]interface{} {
	if originRequestHeader == nil || len(*originRequestHeader) == 0 {
		return nil
	}

	originRequestHeaderAttrs := make([]map[string]interface{}, len(*originRequestHeader))
	for i, v := range *originRequestHeader {
		originRequestHeaderAttrs[i] = map[string]interface{}{
			"name":   v.Name,
			"value":  v.Value,
			"action": v.Action,
		}
	}

	return originRequestHeaderAttrs
}

func flattenHttpResponseHeaderAttrs(httpResponseHeader *[]model.HttpResponseHeader) []map[string]interface{} {
	if httpResponseHeader == nil || len(*httpResponseHeader) == 0 {
		return nil
	}

	httpResponseHeaderAttrs := make([]map[string]interface{}, len(*httpResponseHeader))
	for i, v := range *httpResponseHeader {
		httpResponseHeaderAttrs[i] = map[string]interface{}{
			"name":   v.Name,
			"value":  v.Value,
			"action": v.Action,
		}
	}

	return httpResponseHeaderAttrs
}

func flattenUrlAuthAttrs(urlAuth *model.UrlAuthGetBody, urlAuthKey string) []map[string]interface{} {
	if urlAuth == nil {
		return nil
	}

	urlAuthAttrs := map[string]interface{}{
		"enabled":     analyseFunctionEnabledStatus(urlAuth.Status),
		"status":      urlAuth.Status,
		"type":        urlAuth.Type,
		"key":         urlAuthKey,
		"time_format": urlAuth.TimeFormat,
		"expire_time": urlAuth.ExpireTime,
	}

	return []map[string]interface{}{urlAuthAttrs}
}

func flattenForceRedirectAttrs(forceRedirect *model.ForceRedirectConfig) []map[string]interface{} {
	if forceRedirect == nil {
		return nil
	}

	forceRedirectAttrs := map[string]interface{}{
		"status":  forceRedirect.Status,
		"type":    forceRedirect.Type,
		"enabled": analyseFunctionEnabledStatus(forceRedirect.Status),
	}

	return []map[string]interface{}{forceRedirectAttrs}
}

func flattenCompressAttrs(compress *model.Compress) []map[string]interface{} {
	if compress == nil {
		return nil
	}

	compressAttrs := map[string]interface{}{
		"status":  compress.Status,
		"type":    compress.Type,
		"enabled": analyseFunctionEnabledStatus(compress.Status),
	}

	return []map[string]interface{}{compressAttrs}
}

func flattenCacheUrlParameterFilterAttrs(cacheUrlParameterFilter *model.CacheUrlParameterFilter) []map[string]interface{} {
	if cacheUrlParameterFilter == nil {
		return nil
	}

	cacheUrlParameterFilterAttrs := map[string]interface{}{
		"value": cacheUrlParameterFilter.Value,
		"type":  cacheUrlParameterFilter.Type,
	}

	return []map[string]interface{}{cacheUrlParameterFilterAttrs}
}

func flattenSourcesAttrs(sources *[]model.SourcesConfig) []map[string]interface{} {
	if sources == nil || len(*sources) == 0 {
		return nil
	}

	sourcesAttrs := make([]map[string]interface{}, len(*sources))
	for i, v := range *sources {
		var active int
		if v.Priority == 70 {
			active = 1
		}
		sourcesAttrs[i] = map[string]interface{}{
			"origin":                  v.OriginAddr,
			"origin_type":             v.OriginType,
			"active":                  active,
			"obs_web_hosting_enabled": analyseFunctionEnabledStatusPtr(v.ObsWebHostingStatus),
			"http_port":               v.HttpPort,
			"https_port":              v.HttpsPort,
			"retrieval_host":          v.HostName,
		}
	}

	return sourcesAttrs
}

func flattenCacheRulesAttrs(cacheRulesPtr *[]model.CacheRules) []map[string]interface{} {
	if cacheRulesPtr == nil || len(*cacheRulesPtr) == 0 {
		return nil
	}

	cacheRules := *cacheRulesPtr
	sourcesAttrs := make([]map[string]interface{}, len(cacheRules))
	for i, v := range cacheRules {
		sourcesAttrs[i] = map[string]interface{}{
			"rule_type": v.MatchType,
			"content":   v.MatchValue,
			"ttl":       v.Ttl,
			"ttl_type":  v.TtlUnit,
			"priority":  v.Priority,
		}
	}

	return []map[string]interface{}{
		{
			"follow_origin": analyseFunctionEnabledStatus(cacheRules[0].FollowOrigin),
			"rules":         sourcesAttrs,
		},
	}
}

func getConfigsAttrs(hcCdnClient *cdnv1.CdnClient, domainName, epsId, privateKey, urlAuthKey string) (sources, configs,
	cacheRules []map[string]interface{}, err error) {
	req := model.ShowDomainFullConfigRequest{
		DomainName:          domainName,
		EnterpriseProjectId: utils.StringIgnoreEmpty(epsId),
	}
	resp, err := hcCdnClient.ShowDomainFullConfig(&req)
	if err != nil {
		return
	}

	if resp.Configs == nil {
		err = fmt.Errorf("unbale to find the configs of domain: %s", domainName)
		return
	}

	configsResp := resp.Configs
	configsAttrs := map[string]interface{}{
		"https_settings":                flattenHTTPSAttrs(configsResp.Https, privateKey),
		"retrieval_request_header":      flattenOriginRequestHeaderAttrs(configsResp.OriginRequestHeader),
		"http_response_header":          flattenHttpResponseHeaderAttrs(configsResp.HttpResponseHeader),
		"url_signing":                   flattenUrlAuthAttrs(configsResp.UrlAuth, urlAuthKey),
		"origin_protocol":               configsResp.OriginProtocol,
		"force_redirect":                flattenForceRedirectAttrs(configsResp.ForceRedirect),
		"compress":                      flattenCompressAttrs(configsResp.Compress),
		"cache_url_parameter_filter":    flattenCacheUrlParameterFilterAttrs(configsResp.CacheUrlParameterFilter),
		"ipv6_enable":                   configsResp.Ipv6Accelerate != nil && *configsResp.Ipv6Accelerate == 1,
		"range_based_retrieval_enabled": analyseFunctionEnabledStatusPtr(configsResp.OriginRangeStatus),
	}

	sources = flattenSourcesAttrs(configsResp.Sources)
	configs = []map[string]interface{}{configsAttrs}
	cacheRules = flattenCacheRulesAttrs(configsResp.CacheRules)
	return
}

func getCacheAttrs(hcCdnClient *cdnv1.CdnClient, domainId, epsId string) ([]map[string]interface{}, error) {
	req := model.ShowCacheRulesRequest{
		DomainId:            domainId,
		EnterpriseProjectId: utils.StringIgnoreEmpty(epsId),
	}
	resp, err := hcCdnClient.ShowCacheRules(&req)
	if err != nil {
		return nil, err
	}

	if resp.CacheConfig == nil {
		return nil, fmt.Errorf("unbale to find the cache config of domain: %s", domainId)
	}

	cacheConfig := resp.CacheConfig
	cacheAttrs := map[string]interface{}{
		"follow_origin": cacheConfig.FollowOrigin,
	}

	if cacheConfig.Rules == nil {
		return nil, fmt.Errorf("unbale to find the cache config rules of domain: %s", domainId)
	}
	rules := make([]map[string]interface{}, len(*cacheConfig.Rules))
	for i, v := range *cacheConfig.Rules {
		rules[i] = map[string]interface{}{
			"rule_type": v.RuleType,
			"content":   v.Content,
			"ttl":       v.Ttl,
			"ttl_type":  v.TtlType,
			"priority":  v.Priority,
		}
	}

	cacheAttrs["rules"] = rules

	return []map[string]interface{}{cacheAttrs}, nil
}

func resourceCdnDomainV1Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	cdnClient, err := cfg.CdnV1Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating CDN v1 client: %s", err)
	}

	hcCdnClient, err := cfg.HcCdnV1Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating CDN v1 client: %s", err)
	}

	id := d.Id()
	epsId := cfg.GetEnterpriseProjectID(d)

	opts := getResourceExtensionOpts(d, cfg)
	v, err := domains.Get(cdnClient, id, opts).Extract()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error reading CDN Domain")
	}

	log.Printf("[DEBUG] Retrieved CDN domain %s: %+v", id, v)

	privateKey := d.Get("configs.0.https_settings.0.private_key").(string)
	urlAuthKey := d.Get("configs.0.url_signing.0.key").(string)
	sources, configAttrs, cacheRules, err := getConfigsAttrs(hcCdnClient, v.DomainName, epsId, privateKey, urlAuthKey)
	if err != nil {
		return diag.Errorf("error reading CDN Domain configs settings: %s", err)
	}

	mErr := multierror.Append(nil,
		d.Set("name", v.DomainName),
		d.Set("type", v.BusinessType),
		d.Set("cname", v.CName),
		d.Set("domain_status", v.DomainStatus),
		d.Set("service_area", v.ServiceArea),
		d.Set("sources", sources),
		d.Set("configs", configAttrs),
		d.Set("cache_settings", cacheRules),
	)

	// Set domain tags
	tags, err := hcCdnClient.ShowTags(&model.ShowTagsRequest{ResourceId: id})
	if err != nil {
		return diag.Errorf("error reading CDN Domain tags: %s", err)
	}
	if tags.Tags != nil {
		tagsToSet := make(map[string]interface{}, len(*tags.Tags))
		for _, tag := range *tags.Tags {
			if tag.Value != nil {
				tagsToSet[tag.Key] = *tag.Value
			} else {
				tagsToSet[tag.Key] = ""
			}
		}

		mErr = multierror.Append(mErr, d.Set("tags", tagsToSet))
	}

	if mErr.ErrorOrNil() != nil {
		return diag.FromErr(mErr)
	}
	return nil
}

func resourceCdnDomainV1Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	cdnClient, err := cfg.CdnV1Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating CDN v1 client: %s", err)
	}

	hcCdnClient, err := cfg.HcCdnV1Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating CDN v1 client: %s", err)
	}

	id := d.Id()
	opts := getResourceExtensionOpts(d, cfg)
	timeout := d.Timeout(schema.TimeoutCreate)

	if d.HasChanges("sources", "configs", "cache_settings") || d.IsNewResource() {
		err = updateDomainFullConfigs(hcCdnClient, cfg, d)
		if err != nil {
			return diag.Errorf("error updating CDN Domain configs settings: %s", err)
		}

		// Wait for CDN domain to become active again before continuing
		log.Printf("[INFO] Waiting for CDN domain %s to become online.", id)
		err = waitDomainOnline(ctx, cdnClient, id, opts, timeout)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("tags") {
		oTagsRaw, nTagsRaw := d.GetChange("tags")
		oTagsMap := oTagsRaw.(map[string]interface{})
		nTagsMap := nTagsRaw.(map[string]interface{})

		// remove old tags
		if len(oTagsMap) > 0 {
			var taglist []string
			for k := range oTagsMap {
				taglist = append(taglist, k)
			}
			deleteTagsReq := model.BatchDeleteTagsRequest{
				Body: &model.DeleteTagsRequestBody{
					ResourceId: id,
					Tags:       taglist,
				},
			}
			_, err := hcCdnClient.BatchDeleteTags(&deleteTagsReq)
			if err != nil {
				return diag.Errorf("error deleting CDN Domain tags: %s", err)
			}
		}

		// set new tags
		if len(nTagsMap) > 0 {
			taglist := make([]model.Map, 0, len(nTagsMap))
			for k, v := range nTagsMap {
				tag := model.Map{
					Key:   k,
					Value: utils.String(v.(string)),
				}
				taglist = append(taglist, tag)
			}
			createTagsReq := model.CreateTagsRequest{
				Body: &model.CreateTagsRequestBody{
					ResourceId: id,
					Tags:       taglist,
				},
			}
			_, err := hcCdnClient.CreateTags(&createTagsReq)
			if err != nil {
				return diag.Errorf("error creating CDN Domain tags: %s", err)
			}
		}
	}

	return resourceCdnDomainV1Read(ctx, d, meta)
}

func waitDomainOnline(ctx context.Context, cdnClient *golangsdk.ServiceClient,
	id string, opts *domains.ExtensionOpts, timeout time.Duration) error {
	wait := &WaitDomainStatus{
		ID:      id,
		Penging: []string{"configuring"},
		Target:  []string{"online"},
		Opts:    opts,
	}
	err := waitforCDNV1DomainStatus(ctx, cdnClient, wait, timeout)
	if err != nil {
		return fmt.Errorf("error waiting for CDN domain %s to become online: %s", id, err)
	}

	return nil
}

func resourceCdnDomainV1Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	cdnClient, err := cfg.CdnV1Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating CDN v1 client: %s", err)
	}

	id := d.Id()
	opts := getResourceExtensionOpts(d, cfg)
	timeout := d.Timeout(schema.TimeoutCreate)

	if d.Get("domain_status").(string) == "online" {
		// make sure the status has changed to offline
		log.Printf("[INFO] Disable CDN domain %s.", id)
		if err = domains.Disable(cdnClient, id, opts).Err; err != nil {
			return diag.Errorf("error disable CDN Domain %s: %s", id, err)
		}

		log.Printf("[INFO] Waiting for disabling CDN domain %s.", id)
		wait := &WaitDomainStatus{
			ID:      id,
			Penging: []string{"configuring", "online"},
			Target:  []string{"offline"},
			Opts:    opts,
		}

		err = waitforCDNV1DomainStatus(ctx, cdnClient, wait, timeout)
		if err != nil {
			return diag.Errorf("error waiting for CDN domain %s to become offline: %s", id, err)
		}
	}

	log.Printf("[INFO] Waiting for deleting CDN domain %s.", id)
	_, err = domains.Delete(cdnClient, id, opts).Extract()
	if err != nil {
		return diag.Errorf("error deleting CDN Domain %s: %s", id, err)
	}

	// an API issue will be raised in ForceNew scene, so wait for a while
	time.Sleep(3 * time.Second) // lintignore:R018

	d.SetId("")
	return nil
}

func getResourceExtensionOpts(d *schema.ResourceData, cfg *config.Config) *domains.ExtensionOpts {
	epsID := cfg.GetEnterpriseProjectID(d)
	if epsID != "" {
		return &domains.ExtensionOpts{
			EnterpriseProjectId: epsID,
		}
	}

	return nil
}
