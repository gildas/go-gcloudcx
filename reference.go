package purecloud

// EntityRef describes an Entity that has an ID
type EntityRef struct {
	ID string `json:"id"`
}

// GetID gets the identifier of this
//   implements Identifiable
func (ref EntityRef) GetID() string {
	return ref.ID
}

// AddressableEntityRef describes an Entity that can be addressed
type AddressableEntityRef struct {
	ID      string `json:"id"`
	SelfURI string `json:"selfUri"`
}

// GetID gets the identifier of this
//   implements Identifiable
func (ref AddressableEntityRef) GetID() string {
	return ref.ID
}

// GetURI gets the URI of this
//   implements Addressable
func (ref AddressableEntityRef) GetURI() string {
	return ref.SelfURI
}

// DomainEntityRef describes a DomainEntity Reference
type DomainEntityRef struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	SelfURI string `json:"self_uri"`
}

// GetID gets the identifier of this
//   implements Identifiable
func (ref DomainEntityRef) GetID() string {
	return ref.ID
}

// GetURI gets the URI of this
//   implements Addressable
func (ref DomainEntityRef) GetURI() string {
	return ref.SelfURI
}
