package sumologic

import (
	"log"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceSumologicUser() *schema.Resource {
	return &schema.Resource{
		Create: resourceSumologicUserCreate,
		Read:   resourceSumologicUserRead,
		Update: resourceSumologicUserUpdate,
		Delete: resourceSumologicUserDelete,

		Schema: map[string]*schema.Schema{
			"first_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},
			"last_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},
			"email": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},
			"role_ids": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: false,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"is_active": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: false,
			},
		},
	}
}

func resourceSumologicUserRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	id := d.Id()
	user, err := c.GetUser(id)

	if err != nil {
		return err
	}

	if user == nil {
		log.Printf("[WARN] User not found, removing from state: %v - %v", id, err)
		d.SetId("")
		return nil
	}

	d.Set("firstName", user.FirstName)
	d.Set("lastName", user.LastName)
	d.Set("email", user.Email)
	d.Set("roleIds", user.RoleIds)
	d.Set("isActive", user.IsActive)

	return nil
}

func resourceSumologicUserDelete(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)
	return c.DeleteUser(d.Id())
}

func resourceSumologicUserCreate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	if d.Id() == "" {
		user := resourceToUser(d)
		id, err := c.CreateUser(user)

		if err != nil {
			return err
		}

		d.SetId(id)
	}

	return resourceSumologicUserRead(d, meta)
}

func resourceSumologicUserUpdate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	var user User

	rawRoleIds := d.Get("role_ids").([]interface{})
	roleIds := make([]string, len(rawRoleIds))
	for i, v := range rawRoleIds {
		roleIds[i] = v.(string)
	}

	user.ID = d.Id()
	user.FirstName = d.Get("first_name").(string)
	user.LastName = d.Get("last_name").(string)
	user.Email = d.Get("email").(string)
	user.RoleIds = roleIds
	user.IsActive = d.Get("is_active").(bool)
	err := c.UpdateUser(user)

	if err != nil {
		return err
	}

	return resourceSumologicUserRead(d, meta)
}

func resourceSumologicUserExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	c := meta.(*Client)

	user, err := c.GetUser(d.Id())
	if err != nil {
		return false, err
	}

	return user != nil, nil
}

func resourceToUser(d *schema.ResourceData) User {
	rawRoleIds := d.Get("role_ids").([]interface{})
	roleIds := make([]string, len(rawRoleIds))
	for i, v := range rawRoleIds {
		roleIds[i] = v.(string)
	}

	return User{
		FirstName: d.Get("first_name").(string),
		LastName:  d.Get("last_name").(string),
		Email:     d.Get("email").(string),
		RoleIds:   roleIds,
		IsActive:  d.Get("is_active").(bool),
	}
}
