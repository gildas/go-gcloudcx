package purecloud

// Division describes an Authorization Division
type Division struct {
	ID      string  `json:"id"`
	Name    string  `json:"name"`
	SelfURI string  `json:"selfUri"`
}

// GetID gets the identifier of this
//   implements Identifiable
func (division Division) GetID() string {
	return division.ID
}

// String gets a string version
//   implements the fmt.Stringer interface
func (division Division) String() string {
	if len(division.Name) != 0 {
		return division.Name
	}
	return division.ID
}