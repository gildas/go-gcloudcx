package gcloudcx

import (
	"time"

	"github.com/google/uuid"
)

// EvaluationForm describes an Evaluation Form
type EvaluationForm struct {
	ID      uuid.UUID `json:"id"`
	Name    string    `json:"name"`
	SelfURI URI       `json:"selfUri"`

	ModifiedDate time.Time `json:"modifiedDate"`
	Published    bool      `json:"published"`
	ContextID    string    `json:"contextId"`

	QuestionGroups    []EvaluationQuestionGroup         `json:"questionGroups"`
	PublishedVersions DomainEntityListingEvaluationForm `json:"publishedVersions"`
}

// GetID gets the identifier of this
//   implements Identifiable
func (form EvaluationForm) GetID() uuid.UUID {
	return form.ID
}

// GetURI gets the URI of this
//   implements Addressable
func (form EvaluationForm) GetURI() URI {
	return form.SelfURI
}

// String gets a string version
//   implements the fmt.Stringer interface
func (form EvaluationForm) String() string {
	if len(form.Name) != 0 {
		return form.Name
	}
	return form.ID.String()
}
