package gcloudcx

import (
	"context"

	"github.com/gildas/go-errors"
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
	client                     *Client         `json:"-"`
	logger                     *logger.Logger  `json:"-"`
}

// Fetch fetches an Organization
//
// implements Fetchable
func (organization *Organization) Fetch(ctx context.Context, client *Client, parameters ...interface{}) error {
	id, name, selfURI, log := client.ParseParameters(ctx, organization, parameters...)

	if id != uuid.Nil {
		if err := client.Get(ctx, NewURI("/organizations/%s", id), &organization); err != nil {
			return err
		}
		organization.logger = log
	} else if len(selfURI) > 0 {
		if err := client.Get(ctx, selfURI, &organization); err != nil {
			return err
		}
		organization.logger = log.Record("id", organization.ID)
	} else if len(name) > 0 {
		return errors.NotImplemented.WithStack()
	} else {
		if err := client.Get(ctx, NewURI("/organizations/me"), &organization); err != nil {
			return err
		}
		organization.logger = log.Record("id", organization.ID)
	}
	organization.client = client
	return nil
}

// GetMyOrganization retrives the current Organization
func (client *Client) GetMyOrganization(context context.Context) (*Organization, error) {
	organization := &Organization{}
	if err := client.Get(context, "/organizations/me", &organization); err != nil {
		return nil, err
	}
	organization.client = client
	organization.logger = client.Logger.Child("organization", "organization", "id", organization.ID)
	return organization, nil
}

// GetID gets the identifier of this
//
// implements Identifiable
func (organization Organization) GetID() uuid.UUID {
	return organization.ID
}

// GetURI gets the URI of this
//
// implements Addressable
func (organization Organization) GetURI(ids ...uuid.UUID) URI {
	if len(ids) > 0 {
		return NewURI("/api/v2/organizations/%s", ids[0])
	}
	if organization.ID != uuid.Nil {
		return NewURI("/api/v2/organizations/%s", organization.ID)
	}
	return URI("/api/v2/organizations/")
}

// String gets a string version
//
// implements the fmt.Stringer interface
func (organization Organization) String() string {
	if len(organization.Name) > 0 {
		return organization.Name
	}
	return organization.ID.String()
}
