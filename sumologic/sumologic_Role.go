package sumologic

import (
	"encoding/json"
	"fmt"
)

func (s *Client) CreateRole(role Role) (string, error) {
	data, err := s.Post("v1/roles", role)
	if err != nil {
		return "", err
	}

	var createdrole Role
	err = json.Unmarshal(data, &createdrole)
	if err != nil {
		return "", err
	}

	return createdrole.ID, nil
}

func (s *Client) DeleteRole(id string) error {
	_, err := s.Delete(fmt.Sprintf("v1/roles/%s", id))
	return err
}

func (s *Client) GetRole(id string) (*Role, error) {
	data, _, err := s.Get(fmt.Sprintf("v1/roles/%s", id))
	if err != nil {
		return nil, err
	}
	if data == nil {
		return nil, nil
	}

	var role Role
	err = json.Unmarshal(data, &role)
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (s *Client) UpdateRole(role Role) error {
	url := fmt.Sprintf("v1/roles/%s", role.ID)

	role.ID = ""

	_, err := s.Put(url, role)
	return err
}

// models
type Role struct {
	ID string `json:"id,omitempty"`
	// Name of the role.
	Name string `json:"name"`
	// Description of the role.
	Description string `json:"description,omitempty"`
	// A search filter to restrict access to specific logs. The filter is silently added to the beginning of each query a user runs. For example, using '!_sourceCategory=billing' as a filter predicate will prevent users assigned to the role from viewing logs from the source category named 'billing'.
	FilterPredicate string `json:"filterPredicate"`
	// List of user identifiers to assign the role to.
	Users []string `json:"users,omitempty"`
	// List of [capabilities](https://help.sumologic.com/Manage/Users-and-Roles/Manage-Roles/Role-Capabilities) associated with this role. Valid values are   ### Connections   - manageConnections   ### Collectors   - manageCollectors   - viewCollectors   ### Dashboards   - shareDashboardWhitelist   - shareDashboardWorld   ### Data Management   - manageContent   - manageDataVolumeFeed   - manageFieldExtractionRules   - manageIndexes   - manageS3DataForwarding   ### Metrics   - manageMonitors   - metricsExtraction   ### Security   - ipWhitelisting   - manageAccessKeys   - manageAuditDataFeed   - managePasswordPolicy   - manageSaml   - manageSupportAccountAccess   - manageUsersAndRoles   - shareDashboardOutsideOrg
	Capabilities []string `json:"capabilities,omitempty"`
}
