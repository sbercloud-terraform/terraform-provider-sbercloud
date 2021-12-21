package policies

import "github.com/chnsz/golangsdk"

const (
	rootPath = "OS-ROLE"
	rolePath = "roles"
)

func rootURL(client *golangsdk.ServiceClient) string {
	return client.ServiceURL(rootPath, rolePath)
}

func resourceURL(client *golangsdk.ServiceClient, roleID string) string {
	return client.ServiceURL(rootPath, rolePath, roleID)
}

func listURL(client *golangsdk.ServiceClient) string {
	return client.ServiceURL(rootPath, rolePath)
}
