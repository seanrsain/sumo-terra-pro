package sumologic

import (
	"log"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceSumologicRole() *schema.Resource {
	return &schema.Resource{
		Create: resourceSumologicRoleCreate,
		Read:   resourceSumologicRoleRead,
		Update: resourceSumologicRoleUpdate,
		Delete: resourceSumologicRoleDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
			},
			"filter_predicate": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
			},
			"users": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: false,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"capabilities": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: false,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceSumologicRoleRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	id := d.Id()
	role, err := c.GetRole(id)

	if err != nil {
		return err
	}

	if role == nil {
		log.Printf("[WARN] Role not found, removing from state: %v - %v", id, err)
		d.SetId("")
		return nil
	}

	d.Set("name", role.Name)
	d.Set("description", role.Description)
	d.Set("filterPredicate", role.FilterPredicate)
	d.Set("users", role.Users)
	d.Set("capabilities", role.Capabilities)

	return nil
}

func resourceSumologicRoleDelete(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)
	return c.DeleteRole(d.Id())
}

func resourceSumologicRoleCreate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	if d.Id() == "" {
		role := resourceToRole(d)
		id, err := c.CreateRole(role)

		if err != nil {
			return err
		}

		d.SetId(id)
	}

	return resourceSumologicRoleRead(d, meta)
}

func resourceSumologicRoleUpdate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	role := resourceToRole(d)
	err := c.UpdateRole(role)

	if err != nil {
		return err
	}

	return resourceSumologicRoleRead(d, meta)
}

func resourceSumologicRoleExists(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	_, err := c.GetRole(d.Id())
	if err != nil {
		return err
	}

	return nil
}

func resourceToRole(d *schema.ResourceData) Role {
	rawUsers := d.Get("users").([]interface{})
	users := make([]string, len(rawUsers))
	for i, v := range rawUsers {
		users[i] = v.(string)
	}
	rawCapabilities := d.Get("capabilities").([]interface{})
	capabilities := make([]string, len(rawCapabilities))
	for i, v := range rawCapabilities {
		capabilities[i] = v.(string)
	}

	return Role{
		ID:              d.Id(),
		Name:            d.Get("name").(string),
		Description:     d.Get("description").(string),
		FilterPredicate: d.Get("filter_predicate").(string),
		Users:           users,
		Capabilities:    capabilities,
	}
}
