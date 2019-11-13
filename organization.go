package purecloud

import (
	"github.com/gildas/go-logger"
)

// Organization describes a PureCloud Organization
type Organization struct {
	ID                         string          `json:"id"`
	Name                       string          `json:"name"`
	DefaultLanguage            string          `json:"defaultLanguage"`
	ThirdPartyOrganizationName string          `json:"thirdPartyOrgName"`
	ThirdPartyURI              string          `json:"thirdPartyURI"`
	Domain                     string          `json:"domain"`
	State                      string          `json:"state"`
	DefaultSiteID              string          `json:"defaultSiteId"`
	SupportURI                 string          `json:"supportURI"`
	VoicemailEnabled           bool            `json:"voicemailEnabled"`
	SelfURI                    string          `json:"selfURI"`
	Features                   map[string]bool `json:"features"`
	Version                    uint32          `json:"version"`
	Client                     *Client         `json:"-"`
	Logger                     *logger.Logger  `json:"-"`
}

// Initialize initializes this from the given Client
//   implements Initializable
//   If the organzation ID is not given, /organizations/me is fetched
func (organization *Organization) Initialize(parameters ...interface{}) error {
	client, logger, err := ExtractClientAndLogger(parameters...)
	if err != nil {
		return err
	}
	if len(organization.ID) > 0 {
		if err := client.Get("/organizations/" + organization.ID, &organization); err != nil {
			return err
		}
	} else {
		if err := client.Get("/organizations/me", &organization); err != nil {
			return err
		}
	}
	organization.Client = client
	organization.Logger = logger.Child("organization", "organization", "organization", organization.ID)
	return nil
}

// GetMyOrganization retrives the current Organization
func (client *Client) GetMyOrganization() (*Organization, error) {
	organization := &Organization{}
	if err := client.Get("/organizations/me", &organization); err != nil {
		return nil, err
	}
	organization.Client = client
	organization.Logger = client.Logger.Child("organization", "organization", "organization", organization.ID)
	return organization, nil
}

// GetID gets the identifier of this
//   implements Identifiable
func (organization Organization) GetID() string {
	return organization.ID
}

// String gets a string version
//   implements the fmt.Stringer interface
func (organization Organization) String() string {
	if len(organization.Name) > 0 {
		return organization.Name
	}
	return organization.ID
}