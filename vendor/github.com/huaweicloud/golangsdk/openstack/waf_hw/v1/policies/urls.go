package policies

import "github.com/huaweicloud/golangsdk"

func rootURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL("policy")
}

func resourceURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL("policy", id)
}
