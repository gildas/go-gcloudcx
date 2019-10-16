package purecloud

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
}

// GetMyOrganization retrives the current Organization
func (client *Client) GetMyOrganization() (*Organization, error) {
	organization := &Organization{}
	if err := client.Get("/organizations/me", &organization); err != nil {
		return nil, err
	}
	organization.Client = client
	return organization, nil
}

// String gets a string version
//   implements the fmt.Stringer interface
func (organization Organization) String() string {
	if len(organization.Name) > 0 {
		return organization.Name
	}
	return organization.ID
}