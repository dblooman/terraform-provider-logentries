package logentries

import (
	"errors"

	logentries "github.com/depop/logentries"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceLogentriesLogSet() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceLogentriesLogSetRead,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func dataSourceLogentriesLogSetRead(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*logentries.Client)
	res, err := client.LogSets.Read(&logentries.LogSetsReadRequest{})
	if err != nil {
		panic(err)
	}
	if err != nil {
		return err
	}

	for _, logset := range res.LogSets {
		if logset.Name == d.Get("name").(string) {
			d.SetId(logset.ID)
		}
	}

	if d.Get("Id") == "" {
		return errors.New("logset does not exist")
	}

	return nil

}
