package gcloudcx

// IntegrationConfiguration represents the configuration of an integration.
type IntegrationConfiguration struct {
	ID          string                               `json:"id"`
	Name        string                               `json:"name"`
	Version     uint                                 `json:"version"`
	Notes       string                               `json:"notes"`
	Properties  map[string]any                       `json:"properties,omitempty"`
	Advanced    map[string]any                       `json:"advanced"`
	Credentials map[string]IntegrationCredentialInfo `json:"credentials,omitempty"`
}

// IntegrationConfigurationInfo represents the information about an integration configuration
type IntegrationConfigurationInfo struct {
	Current IntegrationConfiguration `json:"current"`
}
