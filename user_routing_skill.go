package gcloudcx

import "github.com/google/uuid"

// UserRoutingSkill describe a Routing Skill for a User
type UserRoutingSkill struct {
	ID          uuid.UUID `json:"id"`
	SelfURI     URI       `json:"selfUri"`
	Name        string    `json:"name"`
	SkillURI    string    `json:"skillUri"`
	State       string    `json:"state"`
	Proficiency float64   `json:"proficiency"`
}

// GetID gets the identifier of this
//   implements Identifiable
func (userRoutingSkill UserRoutingSkill) GetID() uuid.UUID {
	return userRoutingSkill.ID
}

// GetURI gets the URI of this
//   implements Addressable
func (userRoutingSkill UserRoutingSkill) GetURI() URI {
	return userRoutingSkill.SelfURI
}
