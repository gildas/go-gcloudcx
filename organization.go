package purecloud

import (
	"context"
	"github.com/gildas/go-core"
	"net/http"
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
}

// GetMyOrganization retrives the current Organization
func (client *Client) GetMyOrganization() (*Organization, error) {
	organization := &Organization{}
	/*
	if err := client.Get("organizations/me", nil, &organization); err != nil {
		return nil, err
	}
	*/
	url, err := client.parseURL("organizations/me")
	if err != nil {
		return nil, APIError{ Code: "url.parse", Message: err.Error() }
	}
	res, err := core.SendRequest(context.Background(), &core.RequestOptions{
		Method:     http.MethodPost,
		URL:        url,
		Proxy:      client.Proxy,
		UserAgent:  APP + " " + VERSION,
		Headers:    map[string]string {
			"Authorization": client.Authorization.TokenType + " " + client.Authorization.Token,
		},
		Logger: client.Logger,
	}, &organization)

	if err != nil {
		client.Logger.Record("err", err).Errorf("Core SendRequest error", err)
		if res != nil {
			client.Logger.Infof("Reading error from res")
			apiError := APIError{}
			err = res.UnmarshalContentJSON(&apiError)
			if err != nil { return nil, err }
			return nil, apiError
		}
		return nil, err
	}
	return organization, nil
}