package main

import (
	"github.com/depop/terraform-provider-logentries/logentries"
	"github.com/hashicorp/terraform/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: logentries.Provider})
}
