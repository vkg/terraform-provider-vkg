package main

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/vkg/terraform-provider-vkg/event"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: event.Provider,
	})
}
