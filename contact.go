package gcloudcx

import "github.com/gildas/go-logger"

// Contact describes something that can be contacted
type Contact struct {
	Type      string `json:"type"`      // PRIMARY, WORK, WORK2, WORK3, WORK4, HOME, MOBILE, MAIN
	MediaType string `json:"mediaType"` // PHONE, EMAIL, SMS
	Display   string `json:"display,omitempty"`
	Address   string `json:"address,omitempty"`   // If present, there is no Extension
	Extension string `json:"extension,omitempty"` // If present, there is no Address
}

// String gets a string version
//
//	implements the fmt.Stringer interface
func (contact Contact) String() string {
	return contact.Display
}

// Redact redacts sensitive data
//
//	implements logger.Redactable
func (contact Contact) Redact() interface{} {
	redacted := contact
	if len(contact.Display) > 0 {
		redacted.Display = logger.RedactWithHash(contact.Display)
	}
	if len(contact.Address) > 0 {
		redacted.Address = logger.RedactWithHash(contact.Address)
	}
	return redacted
}
