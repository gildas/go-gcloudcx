package gcloudcx

// UserRoutingSkill describe a Routing Skill for a User
type UserRoutingSkill struct {
	ID          string  `json:"id"`
	SelfURI     string  `json:"selfUri"`
	Name        string  `json:"name"`
	SkillURI    string  `json:"skillUri"`
	State       string  `json:"state"`
	Proficiency float64 `json:"proficiency"`
}
