package groups

import "github.com/chnsz/golangsdk"

func createURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL("api-groups")
}

func groupURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL("api-groups", id)
}
