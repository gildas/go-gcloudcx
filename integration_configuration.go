package gcloudcx

// IntegrationConfiguration represents the configuration of an integration.
type IntegrationConfiguration struct {
	ID          string                             `json:"id"`
	Name        string                             `json:"name"`
	Version     string                             `json:"version"`
	Notes       string                             `json:"notes"`
	Properties  any                                `json:"properties"`
	Advanced    any                                `json:"advanced"`
	Credentials map[string]CredentialSpecification `json:"credentials,omitempty"`
}

// IntegrationConfigurationInfo represents the information about an integration configuration
type IntegrationConfigurationInfo struct {
	Current IntegrationConfiguration `json:"current"`
}
