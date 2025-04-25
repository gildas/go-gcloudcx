package gcloudcx

type EmailAddress struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type EmailAttachment struct {
	Domain DomainEntityRef    `json:"domain"`
	Route  *EmailInboundRoute `json:"route"`
}
