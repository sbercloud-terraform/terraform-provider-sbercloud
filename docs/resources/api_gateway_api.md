---
subcategory: "API Gateway (APIG)"
---

# sbercloud\_api\_gateway\_api

Provides an API gateway API resource.

## Example Usage

```hcl
resource "sbercloud_api_gateway_group" "tf_apigw_group" {
  name        = "tf_apigw_group"
  description = "your descpiption"
}

resource "sbercloud_api_gateway_api" "tf_apigw_api" {
  group_id                 = sbercloud_api_gateway_group.tf_apigw_group.id
  name                     = "tf_apigw_api"
  description              = "your descpiption"
  tags                     = ["tag1", "tag2"]
  visibility               = 2
  auth_type                = "IAM"
  backend_type             = "HTTP"
  request_protocol         = "HTTPS"
  request_method           = "GET"
  request_uri              = "/test/path1"
  example_success_response = "example response"

  http_backend {
    protocol   = "HTTPS"
    method     = "GET"
    uri        = "/web/openapi"
    url_domain = "hc.sbercloud.ru"
    timeout    = 10000
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the API resource. If omitted, the provider-level
  region will be used. Changing this creates a new API resource.

* `name` - (Required, String) Specifies the name of the API. An API name consists of 3–64 characters, starting with a
  letter. Only letters, digits, and underscores (_) are allowed.

* `group_id` - (Required, String, ForceNew) Specifies the ID of the API group. Changing this creates a new resource.

* `description` - (Optional, String) Specifies the description of the API. The description cannot exceed 255 characters.

* `visibility` - (Optional, Int) Specifies whether the API is available to the public. The value can be 1 (public) and
  2 (private). Defaults to 2.

* `auth_type` - (Required, String) Specifies the security authentication mode. The value can be 'APP', 'IAM', and '
  NONE'.

* `request_protocol` - (Optional, String) Specifies the request protocol. The value can be 'HTTP', 'HTTPS', and 'BOTH'
  which means the API can be accessed through both 'HTTP' and 'HTTPS'. Defaults to 'HTTPS'.

* `request_method` - (Required, String) Specifies the request method, including 'GET','POST','PUT' and etc..

* `request_uri` - (Required, String) Specifies the request path of the API. The value must comply with URI
  specifications.

* `backend_type` - (Required, String) Specifies the service backend type. The value can be:
    + 'HTTP': the web service backend
    + 'FUNCTION': the FunctionGraph service backend
    + 'MOCK': the Mock service backend

* `http_backend` - (Optional, List) Specifies the configuration when backend_type selected 'HTTP' (documented below).

* `function_backend` - (Optional, List) Specifies the configuration when backend_type selected 'FUNCTION' (documented
  below).

* `mock_backend` - (Optional, List) Specifies the configuration when backend_type selected 'MOCK' (documented below).

* `request_parameter` - (Optional, List) the request parameter list (documented below).

* `backend_parameter` - (Optional, List) the backend parameter list (documented below).

* `tags` - (Optional, List) the tags of API in format of string list.

* `version` - (Optional, String) Specifies the version of the API. A maximum of 16 characters are allowed.

* `cors` - (Optional, Bool) Specifies whether CORS is supported or not.

* `example_success_response` - (Required, String) Specifies the example response for a successful request. The length
  cannot exceed 20,480 characters.

* `example_failure_response` - (Optional, String) Specifies the example response for a failed request The length cannot
  exceed 20,480 characters.

The `http_backend` object supports the following:

* `protocol` - (Required, String) Specifies the backend request protocol. The value can be 'HTTP' and 'HTTPS'.

* `method` - (Required, String) Specifies the backend request method, including 'GET','POST','PUT' and etc..

* `uri` - (Required, String) Specifies the backend request path. The value must comply with URI specifications.

* `vpc_channel` - (Optional, String) Specifies the VPC channel ID. This parameter and `url_domain` are alternative.

* `url_domain` - (Optional, String) Specifies the backend service address. An endpoint URL is in the format of
  "domain name (or IP address):port number", with up to 255 characters. This parameter and `vpc_channel` are
  alternative.

* `timeout` - (Optional, Int) Timeout duration (in ms) for API Gateway to request for the backend service. Defaults to
    50000.

The `function_backend` object supports the following:

* `function_urn` - (Required, String) Specifies the function URN.

* `invocation_type` - (Required, String) Specifies the invocation mode, which can be 'async' or 'sync'.

* `version` - (Required, String) Specifies the function version.

* `timeout` - (Optional, Int) Timeout duration (in ms) for API Gateway to request for FunctionGraph. Defaults to 50000.

The `mock_backend` object supports the following:

* `result_content` - (Optional, String) Specifies the return result.

* `version` - (Optional, String) Specifies the version of the Mock backend.

* `description` - (Optional, String) Specifies the description of the Mock backend. The description cannot exceed 255
  characters.

The `request_parameter` object supports the following:

* `name` - (Required, String) Specifies the input parameter name. A parameter name consists of 1–32 characters, starting
  with a letter. Only letters, digits, periods (.), hyphens (-), and underscores (_) are allowed.

* `location` - (Required, String) Specifies the input parameter location, which can be 'PATH', 'QUERY' or 'HEADER'.

* `type` - (Required, String) Specifies the input parameter type, which can be 'STRING' or 'NUMBER'.

* `required` - (Required, Bool) Specifies whether the parameter is mandatory or not.

* `default` - (Optional, String) Specifies the default value when the parameter is optional.

* `description` - (Optional, String) Specifies the description of the parameter. The description cannot exceed 255
  characters.

The `backend_parameter` object supports the following:

* `name` - (Required, String) Specifies the parameter name. A parameter name consists of 1–32 characters, starting with
  a letter. Only letters, digits, periods (.), hyphens (-), and underscores (_) are allowed.

* `location` - (Required, String) Specifies the parameter location, which can be 'PATH', 'QUERY' or 'HEADER'.

* `value` - (Required, String) Specifies the parameter value, which is a string of not more than 255 characters. The
  value varies depending on the parameter type:
    + 'REQUEST': parameter name in `request_parameter`
    + 'CONSTANT': real value of the parameter
    + 'SYSTEM': gateway parameter name

* `type` - (Optional, String) Specifies the parameter type, which can be 'REQUEST', 'CONSTANT', or 'SYSTEM'.

* `description` - (Optional, String) Specifies the description of the parameter. The description cannot exceed 255
  characters.

Вот параметры, которых нет в вашей документации, в формате Markdown:

---

### Additional Arguments Reference

* `instance_id`* - (Required, String, ForceNew) The ID of the instance to which the API belongs.

* `type`* - (Required, String) The API type.

* `security_authentication`* - (Optional, String) The security authentication mode of the API request.

* `simple_authentication`* - (Optional, Bool, Computed) Whether the authentication of the application code is enabled.

* `authorizer_id`* - (Optional, String) The ID of the authorizer to which the API request used.

* `content_type`* - (Optional, String, Computed) The content type of the request body.

* `is_send_fg_body_base64`* - (Optional, Bool) Whether to perform Base64 encoding on the body for interaction with FunctionGraph.

* `body_description`* - (Optional, String) The description of the API request body, which can be an example request body, media type or parameters.

* `matching`* - (Optional, String) The matching mode of the API.

* `response_id`* - (Optional, String) The ID of the custom response that API used.

* `success_response`* - (Optional, String) The example response for a successful request.

* `failure_response`* - (Optional, String) The example response for a failure request.

* `mock`* - (Optional, List, Computed, ForceNew, MaxItems: 1) The mock backend details.  
Each `mock` block supports:

  * `status_code`* - (Optional, Int, Computed) The custom status code of the mock response.

  * `response`* - (Optional, String) The response content of the mock.
  
  * `authorizer_id`* - (Optional, String) The ID of the backend custom authorization.

* `func_graph`* - (Optional, List, Computed, ForceNew, MaxItems: 1) The FunctionGraph backend details.  
  Each `func_graph` block supports:
  * `function_urn`* - (Required, String) The URN of the FunctionGraph function.
  * `version`* - (Optional, String) The version of the FunctionGraph function.
  * `function_alias_urn`* - (Optional, String) The alias URN of the FunctionGraph function.
  * `network_type`* - (Optional, String) The network architecture (framework) type.
  * `request_protocol`* - (Optional, String) The request protocol of the FunctionGraph function.
  * `timeout`* - (Optional, Int) The timeout for API requests to backend service.
  * `invocation_type`* - (Optional, String) The invocation type.
  * `authorizer_id`* - (Optional, String) The ID of the backend custom authorization.

* `web`* - (Optional, List, Computed, ForceNew, MaxItems: 1) The web backend details.  
  Each `web` block supports:
  * `path`* - (Required, String) The backend request path.
  * `host_header`* - (Optional, String) The proxy host header.
  * `vpc_channel_id`* - (Optional, String) The VPC channel ID.
  * `backend_address`* - (Optional, String) The backend service address, which consists of a domain name or IP address, and a port number.
  * `request_method`* - (Optional, String) The backend request method of the API.
  * `request_protocol`* - (Optional, String) The web protocol type of the API request.
  * `timeout`* - (Optional, Int) The timeout for API requests to backend service.
  * `retry_count`* - (Optional, Int) The number of retry attempts to request the backend service.
  * `ssl_enable`* - (Optional, Bool) Whether to enable two-way authentication.
  * `authorizer_id`* - (Optional, String) The ID of the backend custom authorization.

* `mock_policy`* - (Optional, Set, MaxItems: 5) The mock policy backends.  
  Each `mock_policy` block supports:
  * `name`* - (Required, String) The backend policy name.
  * `conditions`* - (Required, Set, MaxItems: 5) The policy conditions.
  * `status_code`* - (Optional, Int, Computed) The custom status code of the mock response.
  * `response`* - (Optional, String) The response content of the mock.
  * `effective_mode`* - (Optional, String) The effective mode of the backend policy.
  * `backend_params`* - (Optional, Set) The configuration list of backend parameters.
  * `authorizer_id`* - (Optional, String) The ID of the backend custom authorization.

* `func_graph_policy`* - (Optional, Set, MaxItems: 5) The policy backends of the FunctionGraph function.  
  Each `func_graph_policy` block supports:
  * `name`* - (Required, String) The name of the backend policy.
  * `function_urn`* - (Required, String) The URN of the FunctionGraph function.
  * `version`* - (Optional, String) The version of the FunctionGraph function.
  * `function_alias_urn`* - (Optional, String) The alias URN of the FunctionGraph function.
  * `network_type`* - (Optional, String) The network (framework) type of the FunctionGraph function.
  * `request_protocol`* - (Optional, String) The request protocol of the FunctionGraph function.
  * `conditions`* - (Required, Set, MaxItems: 5) The policy conditions.
  * `invocation_type`* - (Optional, String) The invocation mode of the FunctionGraph function.
  * `effective_mode`* - (Optional, String) The effective mode of the backend policy.
  * `timeout`* - (Optional, Int) The timeout for API requests to backend service.
  * `backend_params`* - (Optional, Set) The configuration list of the backend parameters.
  * `authorizer_id`* - (Optional, String) The ID of the backend custom authorization.
  * `invocation_mode`* - (Optional, String, Deprecated) The invocation mode of the FunctionGraph function.  
    _Deprecated: Use `invocation_type` instead._

* `web_policy`* - (Optional, Set, MaxItems: 5) The web policy backends.  
  Each `web_policy` block supports:
  * `name`* - (Required, String) The name of the web policy.
  * `path`* - (Required, String) The backend request address.
  * `request_method`* - (Required, String) The backend request method of the API.
  * `conditions`* - (Required, Set, MaxItems: 5) The policy conditions.
  * `host_header`* - (Optional, String) The proxy host header.
  * `vpc_channel_id`* - (Optional, String) The VPC channel ID.
  * `backend_address`* - (Optional, String) The backend service address.
  * `request_protocol`* - (Optional, String) The backend request protocol.
  * `effective_mode`* - (Optional, String) The effective mode of the backend policy.
  * `timeout`* - (Optional, Int) The timeout for API requests to backend service.
  * `retry_count`* - (Optional, Int) The number of retry attempts to request the backend service.
  * `backend_params`* - (Optional, Set) The configuration list of the backend parameters.
  * `authorizer_id`* - (Optional, String) The ID of the backend custom authorization.

* `registered_at`* - (Computed, String) The registered time of the API.

* `updated_at`* - (Computed, String) The latest update time of the API.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the API.
* `group_name` - The name of the API group to which the API belongs.

## Timeouts
This resource provides the following timeouts configuration options:
- `create` - Default is 10 minute.
- `delete` - Default is 10 minute.

## Import

API can be imported using the `id`, e.g.

```
$ terraform import sbercloud_api_gateway_api.api "774438a28a574ac8a496325d1bf51807"
```
