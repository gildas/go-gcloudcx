package gcloudcx

import (
	"github.com/google/uuid"
)

// EvaluationQuestion describe an Evaluation Question
type EvaluationQuestion struct {
	ID                  uuid.UUID            `json:"id"`
	Type                string               `json:"type"`
	Text                string               `json:"text"`
	HelpText            string               `json:"helpText"`
	NAEnabled           bool                 `json:"naEnabled"`
	CommentsRequired    bool                 `json:"commentsRequired"`
	IsKill              bool                 `json:"isKill"`
	IsCritical          bool                 `json:"isCritical"`
	VisibilityCondition *VisibilityCondition `json:"visibilityCondition"`
	AnswerOptions       []*AnswerOption      `json:"answerOptions"`
}

// GetID gets the identifier of this
//
//	implements Identifiable
func (question EvaluationQuestion) GetID() uuid.UUID {
	return question.ID
}

// String gets a string version
//
//	implements the fmt.Stringer interface
func (question EvaluationQuestion) String() string {
	return question.ID.String()
}
