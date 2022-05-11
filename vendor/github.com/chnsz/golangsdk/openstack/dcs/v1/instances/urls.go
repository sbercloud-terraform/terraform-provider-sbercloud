package instances

import "github.com/chnsz/golangsdk"

// endpoint/instances
const resourcePath = "instances"
const passwordPath = "password"
const extendPath = "extend"

// createURL will build the rest query url of creation
func createURL(client *golangsdk.ServiceClient) string {
	return client.ServiceURL(client.ProjectID, resourcePath)
}

// deleteURL will build the url of deletion
func deleteURL(client *golangsdk.ServiceClient, id string) string {
	return client.ServiceURL(client.ProjectID, resourcePath, id)
}

// updateURL will build the url of update
func updateURL(client *golangsdk.ServiceClient, id string) string {
	return client.ServiceURL(client.ProjectID, resourcePath, id)
}

// getURL will build the get url of get function
func getURL(client *golangsdk.ServiceClient, id string) string {
	return client.ServiceURL(client.ProjectID, resourcePath, id)
}

// passwordURL will build the password update function
func passwordURL(client *golangsdk.ServiceClient, id string) string {
	return client.ServiceURL(client.ProjectID, resourcePath, id, passwordPath)
}

// extendURL will build the extend update function
func extendURL(client *golangsdk.ServiceClient, id string) string {
	return client.ServiceURL(client.ProjectID, resourcePath, id, extendPath)
}

func listURL(client *golangsdk.ServiceClient) string {
	return client.ServiceURL(client.ProjectID, resourcePath)
}
