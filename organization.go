package gcloudcx

import (
	"context"

	"github.com/gildas/go-logger"
	"github.com/google/uuid"
)

// Organization describes a GCloud Organization
type Organization struct {
	ID                         uuid.UUID       `json:"id"`
	Name                       string          `json:"name"`
	DefaultLanguage            string          `json:"defaultLanguage"`
	ThirdPartyOrganizationName string          `json:"thirdPartyOrgName"`
	ThirdPartyURI              string          `json:"thirdPartyURI"`
	Domain                     string          `json:"domain"`
	State                      string          `json:"state"`
	DefaultSiteID              string          `json:"defaultSiteId"`
	SupportURI                 string          `json:"supportURI"`
	VoicemailEnabled           bool            `json:"voicemailEnabled"`
	SelfURI                    URI             `json:"selfURI"`
	Features                   map[string]bool `json:"features"`
	Version                    uint32          `json:"version"`
	Client                     *Client         `json:"-"`
	Logger                     *logger.Logger  `json:"-"`
}

// Initialize initializes this from the given Client
//   implements Initializable
//   If the organzation ID is not given, /organizations/me is fetched
func (organization *Organization) Initialize(parameters ...interface{}) error {
	context, client, logger, id, err := parseParameters(organization, parameters...)
	if err != nil {
		return err
	}
	if id != uuid.Nil {
		if err := client.Get(context, NewURI("/organizations/%s", id), &organization); err != nil {
			return err
		}
	} else {
		if err := client.Get(context, NewURI("/organizations/me"), &organization); err != nil {
			return err
		}
	}
	organization.Client = client
	organization.Logger = logger.Child("organization", "organization", "organization", organization.ID)
	return nil
}

// GetMyOrganization retrives the current Organization
func (client *Client) GetMyOrganization(context context.Context) (*Organization, error) {
	organization := &Organization{}
	if err := client.Get(context, "/organizations/me", &organization); err != nil {
		return nil, err
	}
	organization.Client = client
	organization.Logger = client.Logger.Child("organization", "organization", "organization", organization.ID)
	return organization, nil
}

// GetID gets the identifier of this
//   implements Identifiable
func (organization Organization) GetID() uuid.UUID {
	return organization.ID
}

// GetURI gets the URI of this
//   implements Addressable
func (organization Organization) GetURI() URI {
	return organization.SelfURI
}

// String gets a string version
//   implements the fmt.Stringer interface
func (organization Organization) String() string {
	if len(organization.Name) > 0 {
		return organization.Name
	}
	return organization.ID.String()
}
