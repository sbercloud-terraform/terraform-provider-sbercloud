package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
	"github.com/sbercloud-terraform/terraform-provider-sbercloud/sbercloud"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: sbercloud.Provider})
}
