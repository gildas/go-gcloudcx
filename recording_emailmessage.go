package gcloudcx

type RecordingEmailMessage struct {
	ID          string            `json:"id"`
	To          []EmailAddress    `json:"to"`
	Cc          []EmailAddress    `json:"cc"`
	Bcc         []EmailAddress    `json:"bcc"`
	From        EmailAddress      `json:"from"`
	Subject     string            `json:"subject"`
	Time        string            `json:"time"`
	HTMLBody    string            `json:"htmlBody"`
	TextBody    string            `json:"textBody"`
	Body        string            `json:"body"`
	Attachments []EmailAttachment `json:"attachments"`
}
