package sumologic

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccSumoLogicUser(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { TestAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckUserDestroy,
		Steps: []resource.TestStep{
			// Create a User
			{
				Config: newUserConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckUserExists("sumologic_user.foo"),
					resource.TestCheckResourceAttr(
						"sumologic_user.foo", "first_name", "Sean"),
					resource.TestCheckResourceAttr(
						"sumologic_user.foo", "last_name", "Terraform"),
					resource.TestCheckResourceAttr(
						"sumologic_user.foo", "email", "ssain+terraform@demo.com"),
					resource.TestCheckResourceAttr(
						"sumologic_user.foo", "role_ids.#", "1"),
					resource.TestCheckResourceAttr(
						"sumologic_user.foo", "role_ids.0", "000000000000022D"),
					resource.TestCheckResourceAttr(
						"sumologic_user.foo", "is_active", "true"),
				),
			},
			// Update a User
			{
				Config: updateUserConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckUserExists("sumologic_user.foo"),
					resource.TestCheckResourceAttr(
						"sumologic_user.foo", "first_name", "Sean"),
					resource.TestCheckResourceAttr(
						"sumologic_user.foo", "last_name", "Terraform Updated"),
					resource.TestCheckResourceAttr(
						"sumologic_user.foo", "email", "ssain+terraform@demo.com"),
					resource.TestCheckResourceAttr(
						"sumologic_user.foo", "role_ids.#", "1"),
					resource.TestCheckResourceAttr(
						"sumologic_user.foo", "role_ids.0", "000000000000022D"),
					resource.TestCheckResourceAttr(
						"sumologic_user.foo", "is_active", "true"),
				),
			},
		},
	})
}

func testAccCheckUserDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*Client)
	for _, r := range s.RootModule().Resources {
		id := r.Primary.ID
		u, err := client.GetUser(id)

		if err != nil {
			return fmt.Errorf("Encountered an error: " + err.Error())
		}

		if u != nil {
			return fmt.Errorf("User still exists")
		}
	}
	return nil
}

func testAccCheckUserExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*Client)
		for _, r := range s.RootModule().Resources {
			id := r.Primary.ID
			if _, err := client.GetUser(id); err != nil {
				return fmt.Errorf("Received an error retrieving user %s", err)
			}
		}
		return nil
	}
}

const newUserConfig = `
resource "sumologic_user" "foo" {
	first_name     = "Sean"
	last_name      = "Terraform"
	email          = "ssain+terraform@demo.com"
	role_ids       = ["000000000000022D"]
	is_active      = true
}
`

const updateUserConfig = `
resource "sumologic_user" "foo" {
	first_name     = "Sean"
	last_name      = "Terraform Updated"
	email          = "ssain+terraform@demo.com"
	role_ids       = ["000000000000022D"]
	is_active      = true
}
`
