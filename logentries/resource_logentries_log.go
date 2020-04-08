package logentries

import (
	"fmt"
	"github.com/depop/logentries"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceLogentriesLog() *schema.Resource {

	return &schema.Resource{
		Create: resourceLogentriesLogCreate,
		Read:   resourceLogentriesLogRead,
		Update: resourceLogentriesLogUpdate,
		Delete: resourceLogentriesLogDelete,
		Exists: resourceLogentriesLogExists,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"token": {
				Type:      schema.TypeString,
				Computed:  true,
				ForceNew:  true,
				Sensitive: true,
			},
			"logset_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"filename": {
				Type:     schema.TypeString,
				Default:  "logentries",
				Optional: true,
			},
			"source": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "token",
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					value := v.(string)
					allowedValues := []string{"token", "syslog", "agent", "api"}
					if !sliceContains(value, allowedValues) {
						errors = append(errors, fmt.Errorf("Invalid log source option: %s (must be one of: %s)", value, allowedValues))
					}
					return
				},
			},
		},
	}
}

func resourceLogentriesLogCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*logentries.Client)

	res, err := client.Log.Create(&logentries.LogCreateRequest{
		Log: logentries.LogCreateRequestFields{
			Name:       d.Get("name").(string),
			SourceType: d.Get("source").(string),
			UserData: logentries.LogUserData{
				LeAgentFilename: d.Get("filename").(string),
			},
			LogsetsInfo: []logentries.LogsetsInfo{
				{
					ID: d.Get("logset_id").(string),
				},
			},
		},
	})

	if err != nil {
		return err
	}

	if d.Get("source").(string) == "token" {
		d.Set("token", res.Tokens[0])
	}

	d.SetId(res.ID)

	return resourceLogentriesLogRead(d, meta)
}

func resourceLogentriesLogExists(d *schema.ResourceData, meta interface{}) (bool, error) {
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

func resourceLogentriesLogRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*logentries.Client)
	res, err := client.Log.Read(&logentries.LogReadRequest{
		ID: d.Id(),
	})

	if err != nil {
		return err
	}

	d.SetId(res.ID)
	d.Set("name", res.Name)
	d.Set("filename", res.UserData.LeAgentFilename)
	d.Set("source", res.SourceType)

	if len(res.Tokens) > 0 {
		d.Set("token", res.Tokens[0])
	}

	if len(res.LogsetsInfo) > 0 {
		d.Set("logset_id", res.LogsetsInfo[0].ID)
	}

	return nil
}

func resourceLogentriesLogUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*logentries.Client)
	_, err := client.Log.Update(&logentries.LogUpdateRequest{
		ID: d.Id(),
		Log: logentries.LogUpdateRequestFields{
			Name: d.Get("name").(string),
			UserData: logentries.LogUserData{
				LeAgentFilename: d.Get("filename").(string),
			},
			LogsetsInfo: []logentries.LogsetsInfo{
				{
					ID: d.Get("logset_id").(string),
				},
			},
		},
	})
	if err != nil {
		return err
	}

	return resourceLogentriesLogRead(d, meta)
}

func resourceLogentriesLogDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*logentries.Client)
	_, err := client.Log.Delete(&logentries.LogDeleteRequest{
		ID: d.Id(),
	})
	return err
}

func sliceContains(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
