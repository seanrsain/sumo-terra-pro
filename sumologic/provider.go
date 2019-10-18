package sumologic

import (
	"fmt"

	"github.com/go-errors/errors"

	"github.com/hashicorp/terraform/helper/mutexkv"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

const DefaultEnvironment = "us2"

func Provider() terraform.ResourceProvider {

	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"access_id": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("SUMOLOGIC_ACCESSID", nil),
			},
			"access_key": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("SUMOLOGIC_ACCESSKEY", nil),
			},
			"environment": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("SUMOLOGIC_ENVIRONMENT", DefaultEnvironment),
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			// "sumologic_collector":     resourceSumologicCollector(),
			"sumologic_user":          resourceSumologicUser(),
			"sumologic_role":          resourceSumologicRole(),
			"sumologic_ingest_budget": resourceSumologicIngestBudget(),
			"sumologic_folder":        resourceSumologicFolder(),
			"sumologic_content":       resourceSumologicContent(),
		},

		DataSourcesMap: map[string]*schema.Resource{
			// "sumologic_collector": dataSourceSumologicCollector(),
			"sumologic_personal_folder": dataSourceSumologicPersonalFolder(),
		},

		ConfigureFunc: providerConfigure,
	}
}

var SumoMutexKV = mutexkv.NewMutexKV()

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	accessId := d.Get("access_id").(string)
	accessKey := d.Get("access_key").(string)
	environment := d.Get("environment").(string)

	msg := ""
	if accessId == "" {
		msg = "sumologic provider: access_id should be set;"
	}
	if accessKey == "" {
		msg = fmt.Sprintf("%s access_key should be set; ", msg)
	}
	if msg != "" {
		if environment == DefaultEnvironment {
			msg = fmt.Sprintf("%s make sure environment is set or that the default (%s) is appropriate", msg, DefaultEnvironment)
		}
		return nil, errors.New(msg)
	}

	return NewClient(
		accessId,
		accessKey,
		environment,
	)
}
