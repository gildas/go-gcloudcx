package gcloudcx

import (
	"github.com/google/uuid"
)

// EvaluationQuestionGroup describes a Group of Evaluation Questions
type EvaluationQuestionGroup struct {
	ID                      uuid.UUID             `json:"id"`
	Name                    string                `json:"name"`
	Type                    string                `json:"type"`
	DefaultAnswersToHighest bool                  `json:"defaultAnswersToHighest"`
	DefaultAnswersToNA      bool                  `json:"defaultAnswersToNA"`
	NAEnabled               bool                  `json:"naEnabled"`
	Weight                  float64               `json:"weight"`
	ManualWeight            bool                  `json:"manualWeight"`
	Questions               []*EvaluationQuestion `json:"questions"`
	VisibilityCondition     *VisibilityCondition  `json:"visibilityCondition"`
}

// GetID gets the identifier of this
//
//	implements Identifiable
func (group EvaluationQuestionGroup) GetID() uuid.UUID {
	return group.ID
}

// String gets a string version
//
//	implements the fmt.Stringer interface
func (group EvaluationQuestionGroup) String() string {
	if len(group.Name) != 0 {
		return group.Name
	}
	return group.ID.String()
}
