package event

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

// Provider returns a terraform.ResourceProvider.
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		ResourcesMap: map[string]*schema.Resource{
			"vkg_1on1":      oneOnOne(),
			"vkg_nomikai":   nomikai(),
			"vkg_tsuribori": tsuribori(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	c := &Config{}
	return c, c.loadAndValidate()
}
