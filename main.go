package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
	"github.com/huaweicloud/terraform-provider-sbercloud/sbercloud"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: sbercloud.Provider})
}
