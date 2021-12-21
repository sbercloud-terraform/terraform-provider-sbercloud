package services

import "github.com/chnsz/golangsdk"

const (
	rootPath     = "vpn"
	resourcePath = "vpnservices"
)

func rootURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL(rootPath, resourcePath)
}

func resourceURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL(rootPath, resourcePath, id)
}
