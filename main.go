package main

import (
	"github.com/hashicorp/terraform/plugin"
	"sumo-terra-pro/sumologic"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: sumologic.Provider,
	})
}
