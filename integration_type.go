package gcloudcx

import (
	"encoding/json"
	"net/url"

	"github.com/gildas/go-core"
	"github.com/gildas/go-errors"
)

// IntegrationType represents the type of integration.
type IntegrationType struct {
	ID                        string                                        `json:"id"`
	Name                      string                                        `json:"name"`
	Description               string                                        `json:"description,omitempty"`
	Provider                  string                                        `json:"provider"`
	Category                  string                                        `json:"category"`
	ConfigPropertiesSchemaURI *url.URL                                      `json:"configPropertiesSchemaUri,omitempty"`
	ConfigAdvancedSchemaURI   *url.URL                                      `json:"configAdvancedSchemaUri,omitempty"`
	HelpURI                   *url.URL                                      `json:"helpUri,omitempty"`
	TermsOfServiceURI         *url.URL                                      `json:"termsOfServiceUri,omitempty"`
	VendorName                string                                        `json:"vendorName,omitempty"`
	VendorWebsiteURI          *url.URL                                      `json:"vendorWebsiteUri,omitempty"`
	MarketplaceURI            *url.URL                                      `json:"marketplaceUri,omitempty"`
	FAQURI                    *url.URL                                      `json:"faqUri,omitempty"`
	PrivacyPolicyURI          *url.URL                                      `json:"privacyPolicyUri,omitempty"`
	SupportContactURI         *url.URL                                      `json:"supportContactUri,omitempty"`
	SalesContactURI           *url.URL                                      `json:"salesContactUri,omitempty"`
	HelpLinks                 []HelpLink                                    `json:"helpLinks,omitempty"`
	NonInstallable            bool                                          `json:"nonInstallable"`
	MaxInstances              int                                           `json:"maxInstances"`
	Credentials               map[string]IntegrationCredentialSpecification `json:"credentials,omitempty"`
	UserPermissions           []string                                      `json:"userPermissions,omitempty"`
	VendorOAUTHClientIDs      []string                                      `json:"vendorOauthClientIds,omitempty"`
	Images                    []Image                                       `json:"images,omitempty"`
}

const (
	IntegrationTypeGenesysDigitalBotConnector = "genesys-digital-bot-connector"
)

// GetID gets the ID of this integration type
//
// implements core.FetchableByNamedID
func (integrationType IntegrationType) GetID() string {
	return integrationType.ID
}

// GetURI gets the URI of this
//
//	implements AddressableByStringID
func (integrationType IntegrationType) GetURI(ids ...string) URI {
	if len(ids) > 0 {
		return NewURI("/api/v2/integrations/types/%s", ids[0])
	}
	if integrationType.ID != "" {
		return NewURI("/api/v2/integrations/types/%s", integrationType.ID)
	}
	return URI("/api/v2/integrations/types/")
}

// Initialize initializes the object
//
// accepted parameters: *gcloudcx.Client, *logger.Logger
//
// implements Initializable
func (integrationType *IntegrationType) Initialize(parameters ...interface{}) {
	for _, raw := range parameters {
		switch parameter := raw.(type) {
		case string:
			integrationType.ID = parameter
		}
	}
}

// MarshalJSON marshals this into JSON
//
// implements json.Marshaler
func (integrationType IntegrationType) MarshalJSON() ([]byte, error) {
	type surrogate IntegrationType
	data, err := json.Marshal(struct {
		surrogate
		ConfigPropertiesSchemaURI *core.URL `json:"configPropertiesSchemaUri,omitempty"`
		ConfigAdvancedSchemaURI   *core.URL `json:"configAdvancedSchemaUri,omitempty"`
		HelpURI                   *core.URL `json:"helpUri,omitempty"`
		TermsOfServiceURI         *core.URL `json:"termsOfServiceUri,omitempty"`
		VendorWebsiteURI          *core.URL `json:"vendorWebsiteUri,omitempty"`
		MarketplaceURI            *core.URL `json:"marketplaceUri,omitempty"`
		FAQURI                    *core.URL `json:"faqUri,omitempty"`
		PrivacyPolicyURI          *core.URL `json:"privacyPolicyUri,omitempty"`
		SupportContactURI         *core.URL `json:"supportContactUri,omitempty"`
		SalesContactURI           *core.URL `json:"salesContactUri,omitempty"`
	}{
		surrogate:                 surrogate(integrationType),
		ConfigPropertiesSchemaURI: (*core.URL)(integrationType.ConfigPropertiesSchemaURI),
		ConfigAdvancedSchemaURI:   (*core.URL)(integrationType.ConfigAdvancedSchemaURI),
		HelpURI:                   (*core.URL)(integrationType.HelpURI),
		TermsOfServiceURI:         (*core.URL)(integrationType.TermsOfServiceURI),
		VendorWebsiteURI:          (*core.URL)(integrationType.VendorWebsiteURI),
		MarketplaceURI:            (*core.URL)(integrationType.MarketplaceURI),
		FAQURI:                    (*core.URL)(integrationType.FAQURI),
		PrivacyPolicyURI:          (*core.URL)(integrationType.PrivacyPolicyURI),
		SupportContactURI:         (*core.URL)(integrationType.SupportContactURI),
		SalesContactURI:           (*core.URL)(integrationType.SalesContactURI),
	})
	if err != nil {
		return nil, errors.JSONMarshalError.Wrap(err)
	}
	return data, nil
}

// UnmarshalJSON unmarshals JSON into this
//
// implements json.Unmarshaler
func (integrationType *IntegrationType) UnmarshalJSON(payload []byte) error {
	type surrogate IntegrationType
	var inner struct {
		surrogate
		ConfigPropertiesSchemaURI *core.URL `json:"configPropertiesSchemaUri,omitempty"`
		ConfigAdvancedSchemaURI   *core.URL `json:"configAdvancedSchemaUri,omitempty"`
		HelpURI                   *core.URL `json:"helpUri,omitempty"`
		TermsOfServiceURI         *core.URL `json:"termsOfServiceUri,omitempty"`
		VendorWebsiteURI          *core.URL `json:"vendorWebsiteUri,omitempty"`
		MarketplaceURI            *core.URL `json:"marketplaceUri,omitempty"`
		FAQURI                    *core.URL `json:"faqUri,omitempty"`
		PrivacyPolicyURI          *core.URL `json:"privacyPolicyUri,omitempty"`
		SupportContactURI         *core.URL `json:"supportContactUri,omitempty"`
		SalesContactURI           *core.URL `json:"salesContactUri,omitempty"`
	}
	if err := json.Unmarshal(payload, &inner); err != nil {
		return errors.JSONUnmarshalError.WrapIfNotMe(err)
	}
	*integrationType = IntegrationType(inner.surrogate)
	integrationType.ConfigPropertiesSchemaURI = (*url.URL)(inner.ConfigPropertiesSchemaURI)
	integrationType.ConfigAdvancedSchemaURI = (*url.URL)(inner.ConfigAdvancedSchemaURI)
	integrationType.HelpURI = (*url.URL)(inner.HelpURI)
	integrationType.TermsOfServiceURI = (*url.URL)(inner.TermsOfServiceURI)
	integrationType.VendorWebsiteURI = (*url.URL)(inner.VendorWebsiteURI)
	integrationType.MarketplaceURI = (*url.URL)(inner.MarketplaceURI)
	integrationType.FAQURI = (*url.URL)(inner.FAQURI)
	integrationType.PrivacyPolicyURI = (*url.URL)(inner.PrivacyPolicyURI)
	integrationType.SupportContactURI = (*url.URL)(inner.SupportContactURI)
	integrationType.SalesContactURI = (*url.URL)(inner.SalesContactURI)
	return nil
}
