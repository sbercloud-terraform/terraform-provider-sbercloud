package policies

import "github.com/chnsz/golangsdk"

const rootPath = "policies"

func rootURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL(rootPath)
}

func resourceURL(c *golangsdk.ServiceClient, policyid string) string {
	return c.ServiceURL(rootPath, policyid)
}
