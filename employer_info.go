package gcloudcx

// EmployerInfo  describes Employer Information
type EmployerInfo struct {
	OfficialName string `json:"officialName"`
	EmployeeID   string `json:"employeeId"`
	EmployeeType string `json:"employeeType"`
	HiredSince   string `json:"dateHire"`
}

// String gets a string version
//   implements the fmt.Stringer interface
func (info EmployerInfo) String() string {
	return info.OfficialName
}
