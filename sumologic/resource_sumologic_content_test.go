package sumologic

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

//Testing create functionality for Content resources
func TestAccContentCreate(t *testing.T) {
	var content Content
	personalContentId := os.Getenv("SUMOLOGIC_PF")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { TestAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContentDestroy(content),
		Steps: []resource.TestStep{
			{
				Config: testAccSumologicContent(personalContentId, configJson),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckContentExists("sumologic_content.test", &content, t),
					testAccCheckContentAttributes("sumologic_content.test"),
					//					testAccCheckContentConfig(&content),
					resource.TestCheckResourceAttr("sumologic_content.test", "parent_id", personalContentId),
				),
			},
		},
	})
}

func TestAccContentUpdate(t *testing.T) {
	var content Content
	personalContentId := os.Getenv("SUMOLOGIC_PF")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { TestAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContentDestroy(content),
		Steps: []resource.TestStep{
			{
				Config: testAccSumologicContent(personalContentId, configJson),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckContentExists("sumologic_content.test", &content, t),
					testAccCheckContentAttributes("sumologic_content.test"),
					//                                      testAccCheckContentConfig(&content),
					resource.TestCheckResourceAttr("sumologic_content.test", "parent_id", personalContentId),
				),
			}, {
				Config: testAccSumologicContent(personalContentId, updateConfigJson),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckContentExists("sumologic_content.test", &content, t),
					testAccCheckContentAttributes("sumologic_content.test"),
					//                                      testAccCheckContentConfig(&content),
					resource.TestCheckResourceAttr("sumologic_content.test", "parent_id", personalContentId),
				),
			},
		},
	})
}

func testAccCheckContentExists(name string, content *Content, t *testing.T) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Content not found: %s", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Content ID is not set")
		}

		id := rs.Primary.ID
		c := testAccProvider.Meta().(*Client)
		newContent, err := c.GetContent(id)
		if err != nil {
			return fmt.Errorf("Content %s not found", id)
		}
		content = newContent
		return nil
	}
}

func testAccCheckContentAttributes(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		f := resource.ComposeTestCheckFunc(
			//			resource.TestCheckResourceAttrSet(name, "config"),
			resource.TestCheckResourceAttrSet(name, "parent_id"),
		)
		return f(s)
	}
}

func testAccCheckContentConfig(content *Content) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		var expectedContent Content
		//unmarshal the expected config for comparison. Ignore the error here, configuration is known
		_ = json.Unmarshal([]byte(configJson), &expectedContent)

		//if the configuration structs are not equal
		if !reflect.DeepEqual(expectedContent, content) {
			return fmt.Errorf("Configuration is not equal, %v does not match expected %v", content, expectedContent)
		}
		return nil
	}
}

func testAccCheckContentDestroy(content Content) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*Client)
		_, err := client.GetContent(content.ID)
		if err == nil {
			return fmt.Errorf("Content still exists")
		}
		return nil
	}
}

var updateConfigJson = `{
    "type": "SavedSearchWithScheduleSyncDefinition",
    "name": "test-333",
    "search": {
        "queryText": "\"warn\"",
        "defaultTimeRange": "-15m",
        "byReceiptTime": false,
        "viewName": "",
        "viewStartTime": "1970-01-01T00:00:00Z",
        "queryParameters": []
    },
    "searchSchedule": {
        "cronExpression": "0 0 * * * ? *",
        "displayableTimeRange": "-10m",
        "parseableTimeRange": {
            "type": "BeginBoundedTimeRange",
            "from": {
                "type": "RelativeTimeRangeBoundary",
                "relativeTime": "-50m"
            },
            "to": null
        },
        "timeZone": "America/Los_Angeles",
        "threshold": null,
        "notification": {
            "taskType": "EmailSearchNotificationSyncDefinition",
            "toList": [
                "ops@acme.org"
            ],
            "subjectTemplate": "Search Results: {{SearchName}}",
            "includeQuery": true,
            "includeResultSet": true,
            "includeHistogram": false,
            "includeCsvAttachment": false
        },
        "scheduleType": "1Hour",
        "muteErrorEmails": false,
        "parameters": []
    },
    "description": "Runs every hour with timerange of 15m and sends email notifications"
}
`

var configJson = `{
    "type": "SavedSearchWithScheduleSyncDefinition",
    "name": "test-121",
    "search": {
        "queryText": "\"error\"",
        "defaultTimeRange": "-15m",
        "byReceiptTime": false,
        "viewName": "",
        "viewStartTime": "1970-01-01T00:00:00Z",
        "queryParameters": []
    },
    "searchSchedule": {
        "cronExpression": "0 0 * * * ? *",
        "displayableTimeRange": "-10m",
        "parseableTimeRange": {
            "type": "BeginBoundedTimeRange",
            "from": {
                "type": "RelativeTimeRangeBoundary",
                "relativeTime": "-50m"
            },
            "to": null
        },
        "timeZone": "America/Los_Angeles",
        "threshold": null,
        "notification": {
            "taskType": "EmailSearchNotificationSyncDefinition",
            "toList": [
                "ops@acme.org"
            ],
            "subjectTemplate": "Search Results: {{SearchName}}",
            "includeQuery": true,
            "includeResultSet": true,
            "includeHistogram": false,
            "includeCsvAttachment": false
        },
        "scheduleType": "1Hour",
        "muteErrorEmails": false,
        "parameters": []
    },
    "description": "Runs every hour with timerange of 15m and sends email notifications"
}
`

func testAccSumologicContent(parentId string, configJson string) string {
	return fmt.Sprintf(`
resource "sumologic_content" "test" {
  parent_id = "%s"
  config = <<JSON
%s
JSON
}
`, parentId, configJson)
}
