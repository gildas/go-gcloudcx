package gcloudcx

// ProcessAutomationCriteria is a Process Automation Criteria
type ProcessAutomationCriteria struct {
	JSONPath string   `json:"jsonPath"`
	Operator string   `json:"operator"` // GreaterThanOrEqual, LessThanOrEqual, Equal, NotEqual, LessThan, GreaterThan, NotIn, In, Contains, All, Exists, Size
	Value    string   `json:"value,omitempty"`
	Values   []string `json:"values,omitempty"`
}
