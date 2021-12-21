package volumes

import "github.com/chnsz/golangsdk"

func createURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL("cloudvolumes")
}

func getURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL("os-vendor-volumes", id)
}
