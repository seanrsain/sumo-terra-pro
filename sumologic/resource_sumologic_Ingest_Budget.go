package sumologic

import (
	"log"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceSumologicIngestBudget() *schema.Resource {
	return &schema.Resource{
		Create: resourceSumologicIngestBudgetCreate,
		Read:   resourceSumologicIngestBudgetRead,
		Update: resourceSumologicIngestBudgetUpdate,
		Delete: resourceSumologicIngestBudgetDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},
			"field_value": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},
			"capacity_bytes": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: false,
			},
			"timezone": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},
			"reset_time": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
			},
			"action": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},
			"audit_threshold": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: false,
			},
		},
	}
}

func resourceSumologicIngestBudgetRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	id := d.Id()
	ingestbudget, err := c.GetIngestBudget(id)

	if err != nil {
		return err
	}

	if ingestbudget == nil {
		log.Printf("[WARN] IngestBudget not found, removing from state: %v - %v", id, err)
		d.SetId("")
		return nil
	}

	d.Set("name", ingestbudget.Name)
	d.Set("field_value", ingestbudget.FieldValue)
	d.Set("capacity_bytes", ingestbudget.CapacityBytes)
	d.Set("timezone", ingestbudget.Timezone)
	d.Set("reset_time", ingestbudget.ResetTime)
	d.Set("description", ingestbudget.Description)
	d.Set("action", ingestbudget.Action)
	d.Set("audit_threshold", ingestbudget.AuditThreshold)

	return nil
}

func resourceSumologicIngestBudgetDelete(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)
	return c.DeleteIngestBudget(d.Id())
}

func resourceSumologicIngestBudgetCreate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	if d.Id() == "" {
		ingestbudget := resourceToIngestBudget(d)
		id, err := c.CreateIngestBudget(ingestbudget)

		if err != nil {
			return err
		}

		d.SetId(id)
	}

	return resourceSumologicIngestBudgetRead(d, meta)
}

func resourceSumologicIngestBudgetUpdate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	ingestbudget := resourceToIngestBudget(d)
	err := c.UpdateIngestBudget(ingestbudget)

	if err != nil {
		return err
	}

	return resourceSumologicIngestBudgetRead(d, meta)
}

func resourceSumologicIngestBudgetExists(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	_, err := c.GetIngestBudget(d.Id())
	if err != nil {
		return err
	}

	return nil
}

func resourceToIngestBudget(d *schema.ResourceData) IngestBudget {

	return IngestBudget{
		ID:             d.Id(),
		Name:           d.Get("name").(string),
		FieldValue:     d.Get("field_value").(string),
		CapacityBytes:  d.Get("capacity_bytes").(int),
		Timezone:       d.Get("timezone").(string),
		ResetTime:      d.Get("reset_time").(string),
		Description:    d.Get("description").(string),
		Action:         d.Get("action").(string),
		AuditThreshold: d.Get("audit_threshold").(int),
	}
}
