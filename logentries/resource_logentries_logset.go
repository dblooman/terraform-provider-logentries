package logentries

import (
	"github.com/depop/logentries"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceLogentriesLogSet() *schema.Resource {

	return &schema.Resource{
		Create: resourceLogentriesLogSetCreate,
		Read:   resourceLogentriesLogSetRead,
		Update: resourceLogentriesLogSetUpdate,
		Delete: resourceLogentriesLogSetDelete,
		Exists: resourceLogentriesLogSetExists,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceLogentriesLogSetExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	client := meta.(*logentries.Client)
	_, err := client.Log.Read(&logentries.LogReadRequest{
		ID: d.Id(),
	})
	if err != nil {
		if err == logentries.ErrNotFound {
			return false, nil
		}
		return false, err
	}

	return true, err
}

func resourceLogentriesLogSetCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*logentries.Client)
	res, err := client.LogSet.Create(&logentries.LogSetCreateRequest{
		LogSet: logentries.LogSetFields{
			Name: d.Get("name").(string),
		},
	})

	if err != nil {
		return err
	}

	d.SetId(res.ID)

	return resourceLogentriesLogSetRead(d, meta)
}

func resourceLogentriesLogSetRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*logentries.Client)
	res, err := client.LogSet.Read(&logentries.LogSetReadRequest{
		ID: d.Id(),
	})

	if err != nil {
		return err
	}

	d.SetId(res.ID)
	d.Set("name", res.Name)
	d.Set("description", res.Description)

	return nil
}

func resourceLogentriesLogSetUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*logentries.Client)
	_, err := client.LogSet.Update(&logentries.LogSetUpdateRequest{
		ID: d.Id(),
		LogSet: logentries.LogSetFields{
			Name:        d.Get("name").(string),
			Description: d.Get("description").(string),
		},
	})
	if err != nil {
		return err
	}

	return resourceLogentriesLogRead(d, meta)
}

func resourceLogentriesLogSetDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*logentries.Client)
	_, err := client.LogSet.Delete(&logentries.LogSetDeleteRequest{
		ID: d.Id(),
	})
	return err
}
