package gcloudcx

type OpenMessageTo struct {
	ID   string `json:"id"`
	Type string `json:"idType,omitempty"`
}

type OpenMessageFrom struct {
	ID        string `json:"id"`
	Type      string `json:"idType,omitempty"`
	Firstname string `json:"firstName,omitempty"`
	Lastname  string `json:"lastName,omitempty"`
	Nickname  string `json:"nickname,omitempty"`
}

// Redact redacts sensitive data
//
// implements logger.Redactable
func (from OpenMessageFrom) Redact() interface{} {
	redacted := from
	if len(from.Firstname) > 0 {
		redacted.Firstname = "REDACTED"
	}
	if len(from.Lastname) > 0 {
		redacted.Lastname = "REDACTED"
	}
	if len(from.Nickname) > 0 {
		redacted.Nickname = "REDACTED"
	}
	return &redacted
}
