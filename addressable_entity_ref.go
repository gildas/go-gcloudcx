package purecloud

// AddressableEntityRef describes an Entity that can be addressed
type AddressableEntityRef struct {
	ID      string `json:"id"`
	SelfURI string `json:"selfUri"`
}