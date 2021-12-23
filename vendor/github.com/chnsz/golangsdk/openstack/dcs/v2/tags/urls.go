package tags

import "github.com/chnsz/golangsdk"

func resourceURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL("instances", id, "tags")
}

func actionURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL("dcs", id, "tags", "action")
}

func listURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL("dcs", "tags")
}
