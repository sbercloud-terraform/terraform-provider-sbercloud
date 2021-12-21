package products

import (
	"strings"

	"github.com/chnsz/golangsdk"
)

// endpoint/products
const resourcePath = "products"

// getURL will build the get url of get function
func getURL(client *golangsdk.ServiceClient, engine string) string {
	// remove projectid from endpoint
	return strings.Replace(client.ServiceURL(resourcePath+"?engine="+engine), "/"+client.ProjectID, "", -1)
}
