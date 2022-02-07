package apis

import "github.com/chnsz/golangsdk"

func createURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL("apis")
}

func groupURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL("apis", id)
}
