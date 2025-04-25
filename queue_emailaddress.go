package gcloudcx

type QueueEmailAddress struct {
	Domain DomainEntityRef    `json:"domain"`
	Route  *EmailInboundRoute `json:"route,omitempty"`
}
