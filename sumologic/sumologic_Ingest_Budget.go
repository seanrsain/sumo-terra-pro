package sumologic

import (
	"encoding/json"
	"fmt"
)

func (s *Client) CreateIngestBudget(ingestBudget IngestBudget) (string, error) {
	data, err := s.Post("v1/ingestBudgets", ingestBudget)
	if err != nil {
		return "", err
	}

	var createdingestBudget IngestBudget
	err = json.Unmarshal(data, &createdingestBudget)
	if err != nil {
		return "", err
	}

	return createdingestBudget.ID, nil
}

func (s *Client) DeleteIngestBudget(id string) error {
	_, err := s.Delete(fmt.Sprintf("v1/ingestBudgets/%s", id))
	return err
}

func (s *Client) GetIngestBudget(id string) (*IngestBudget, error) {
	data, _, err := s.Get(fmt.Sprintf("v1/ingestBudgets/%s", id))
	if err != nil {
		return nil, err
	}
	if data == nil {
		return nil, nil
	}

	var ingestBudget IngestBudget
	err = json.Unmarshal(data, &ingestBudget)
	if err != nil {
		return nil, err
	}
	return &ingestBudget, nil
}

func (s *Client) UpdateIngestBudget(ingestBudget IngestBudget) error {
	url := fmt.Sprintf("v1/ingestBudgets/%s", ingestBudget.ID)

	ingestBudget.ID = ""

	_, err := s.Put(url, ingestBudget)
	return err
}

// models
type IngestBudget struct {
	ID string `json:"id,omitempty"`
	// Display name of the ingest budget.
	Name string `json:"name"`
	// Custom field value that is used to assign Collectors to the ingest budget.
	FieldValue string `json:"fieldValue"`
	// Capacity of the ingest budget, in bytes. It takes a few minutes for Collectors to stop collecting when capacity is reached. We recommend setting a soft limit that is lower than your needed hard limit.
	CapacityBytes int `json:"capacityBytes"`
	// Time zone of the reset time for the ingest budget. Follow the format in the [IANA Time Zone Database](https://en.wikipedia.org/wiki/List_of_tz_database_time_zones#List).
	Timezone string `json:"timezone"`
	// Reset time of the ingest budget in HH:MM format.
	ResetTime string `json:"resetTime"`
	// Description of the ingest budget.
	Description string `json:"description,omitempty"`
	// Action to take when ingest budget's capacity is reached. All actions are audited. Supported values are:   * `stopCollecting`   * `keepCollecting`
	Action string `json:"action"`
	// The threshold as a percentage of when an ingest budget's capacity usage is logged in the Audit Index.
	AuditThreshold int `json:"auditThreshold,omitempty"`
}
