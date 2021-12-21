package availablezones

import (
	"strings"

	"github.com/chnsz/golangsdk"
)

// endpoint/availablezones
const resourcePath = "availableZones"

// getURL will build the get url of get function
func getURL(client *golangsdk.ServiceClient) string {
	// remove projectid from endpoint
	return strings.Replace(client.ServiceURL(resourcePath), "/"+client.ProjectID, "", -1)
}
