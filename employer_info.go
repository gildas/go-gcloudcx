package purecloud

// EmployerInfo  describes Employer Information
type EmployerInfo struct {
	OfficialName string `json:"officialName"`
	EmployeeID   string `json:"employeeId"`
	EmployeeType string `json:"employeeType"`
	HiredSince   string `json:"dateHire"`
}