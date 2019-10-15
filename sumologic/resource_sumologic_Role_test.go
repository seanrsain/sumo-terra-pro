package sumologic

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccRoleCreate(t *testing.T) {
	var role Role
	testname := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	testdescription := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	testfilterPredicate := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	testusers := []string{strconv.Quote(acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum))}
	testcapabilities := []string{strconv.Quote(acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum))}
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { TestAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRoleDestroy(role),
		Steps: []resource.TestStep{
			{
				Config: testAccSumologicRole(testname, testdescription, testfilterPredicate, testusers, testcapabilities),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRoleExists("sumologic_role.test", &role, t),
					testAccCheckRoleAttributes("sumologic_role.test"),
					resource.TestCheckResourceAttr("sumologic_role.test", "name", testname),
					resource.TestCheckResourceAttr("sumologic_role.test", "description", testdescription),
					resource.TestCheckResourceAttr("sumologic_role.test", "filter_predicate", testfilterPredicate),
					resource.TestCheckResourceAttr("sumologic_role.test", "users.#", "0"),
					resource.TestCheckResourceAttr("sumologic_role.test", "capabilities.#", "0"),
				),
			},
		},
	})
}

func testAccCheckRoleDestroy(role Role) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*Client)
		for _, r := range s.RootModule().Resources {
			id := r.Primary.ID
			u, err := client.GetRole(id)
			if err != nil {
				return fmt.Errorf("Encountered an error: " + err.Error())
			}
			if u != nil {
				return fmt.Errorf("Role still exists")
			}
		}
		return nil
	}
}

func testAccCheckRoleExists(name string, role *Role, t *testing.T) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Role not found: %s", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Role ID is not set")
		}

		id := rs.Primary.ID
		c := testAccProvider.Meta().(*Client)
		newRole, err := c.GetRole(id)
		if err != nil {
			return fmt.Errorf("Role %s not found", id)
		}
		role = newRole
		return nil
	}
}

func TestAccRoleUpdate(t *testing.T) {
	var role Role
	testname := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	testdescription := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	testfilterPredicate := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	testusers := []string{strconv.Quote(acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum))}
	testcapabilities := []string{strconv.Quote(acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum))}

	testUpdatedname := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	testUpdateddescription := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	testUpdatedfilterPredicate := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	testUpdatedusers := []string{strconv.Quote(acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum))}
	testUpdatedcapabilities := []string{strconv.Quote(acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum))}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { TestAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRoleDestroy(role),
		Steps: []resource.TestStep{
			{
				Config: testAccSumologicRole(testname, testdescription, testfilterPredicate, testusers, testcapabilities),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRoleExists("sumologic_role.test", &role, t),
					testAccCheckRoleAttributes("sumologic_role.test"),
					resource.TestCheckResourceAttr("sumologic_role.test", "name", testname),
					resource.TestCheckResourceAttr("sumologic_role.test", "description", testdescription),
					resource.TestCheckResourceAttr("sumologic_role.test", "filter_predicate", testfilterPredicate),
					resource.TestCheckResourceAttr("sumologic_role.test", "users.#", "0"),
					resource.TestCheckResourceAttr("sumologic_role.test", "capabilities.#", "0"),
				),
			},
			{
				Config: testAccSumologicRoleUpdate(testUpdatedname, testUpdateddescription, testUpdatedfilterPredicate, testUpdatedusers, testUpdatedcapabilities),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRoleExists("sumologic_role.test", &role, t),
					testAccCheckRoleAttributes("sumologic_role.test"),
					resource.TestCheckResourceAttr("sumologic_role.test", "name", testUpdatedname),
					resource.TestCheckResourceAttr("sumologic_role.test", "description", testUpdateddescription),
					resource.TestCheckResourceAttr("sumologic_role.test", "filter_predicate", testUpdatedfilterPredicate),
					resource.TestCheckResourceAttr("sumologic_role.test", "users.#", "0"),
					resource.TestCheckResourceAttr("sumologic_role.test", "capabilities.#", "0"),
				),
			},
		},
	})
}

func testAccSumologicRole(name string, description string, filterPredicate string, users []string, capabilities []string) string {
	return fmt.Sprintf(`
resource "sumologic_role" "test" {
    name = "%s"

    description = "%s"

    filter_predicate = "%s"

    users = []

    capabilities = []

}
`, name, description, filterPredicate)
}

func testAccSumologicRoleUpdate(name string, description string, filterPredicate string, users []string, capabilities []string) string {
	return fmt.Sprintf(`
resource "sumologic_role" "test" {
      name = "%s"

      description = "%s"

      filter_predicate = "%s"

      users = []

      capabilities = []

}
`, name, description, filterPredicate)
}

func testAccCheckRoleAttributes(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		f := resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttrSet(name, "name"),
			resource.TestCheckResourceAttrSet(name, "description"),
			resource.TestCheckResourceAttrSet(name, "filter_predicate"),
		)
		return f(s)
	}
}
