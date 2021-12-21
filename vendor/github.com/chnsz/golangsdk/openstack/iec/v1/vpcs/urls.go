package vpcs

import "github.com/chnsz/golangsdk"

func rootURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL("vpcs")
}

func DeleteURL(c *golangsdk.ServiceClient, vpcId string) string {
	return c.ServiceURL("vpcs", vpcId)
}

func GetURL(c *golangsdk.ServiceClient, vpcId string) string {
	return c.ServiceURL("vpcs", vpcId)
}

func UpdateURL(c *golangsdk.ServiceClient, vpcId string) string {
	return c.ServiceURL("vpcs", vpcId)
}
