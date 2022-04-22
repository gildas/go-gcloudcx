package gcloudcx

type AuthorizationGrantPolicy struct {
	EntityName string   `json:"entityName"`
	Domain     string   `json:"domain"`
	Condition  string   `json:"condition"`
	Actions    []string `json:"actions"`
}