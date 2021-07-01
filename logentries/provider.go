package logentries

import (
	"github.com/depop/logentries"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

// Provider returns a terraform.ResourceProvider.
func Provider() terraform.ResourceProvider {

	// The actual provider
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"account_key": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("LOGENTRIES_TOKEN", nil),
				Description: descriptions["account_key"],
			},
			"proxy_url": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: descriptions["proxy_url"],
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"logentries_log":    resourceLogentriesLog(),
			"logentries_logset": resourceLogentriesLogSet(),
		},

		DataSourcesMap: map[string]*schema.Resource{
			"logentries_logset": dataSourceLogentriesLogSet(),
		},

		ConfigureFunc: providerConfigure,
	}
}

var descriptions map[string]string

func init() {
	descriptions = map[string]string{
		"account_key": "The Log Entries account key.",
		"proxy_url":   "The proxy url you use as cache",
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	return logentries.New(d.Get("account_key").(string), d.Get("proxy_url").(string)), nil
}
