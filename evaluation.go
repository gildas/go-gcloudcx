package gcloudcx

import (
	"time"

	"github.com/google/uuid"
)

// Evaluation describes an Evaluation (like belonging to Participant)
type Evaluation struct {
	ID             uuid.UUID             `json:"id"`
	Name           string                `json:"name"`
	Status         string                `json:"status"`
	Queue          *Queue                `json:"queue"`
	Conversation   *Conversation         `json:"conversation"`
	EvaluationForm *EvaluationForm       `json:"evaluationForm"`
	Evaluator      *User                 `json:"evaluator"`
	Agent          *User                 `json:"agent"`
	Calibration    *Calibration          `json:"calibration"`
	Answers        *EvaluationScoringSet `json:"answers"`
	AgentHasRead   bool                  `json:"agentHasRead"`
	ReleaseDate    time.Time             `json:"releaseDate"`
	AssignedDate   time.Time             `json:"assignedDate"`
	ChangedDate    time.Time             `json:"changedDate"`
}

// VisibilityCondition  describes visibility conditions
type VisibilityCondition struct {
	CombiningOperation string        `json:"combiningOperation"`
	Predicates         []interface{} `json:"predicates"`
}

// AnswerOption describes an Answer Option
type AnswerOption struct {
	ID    uuid.UUID `json:"id"`
	Text  string    `json:"text"`
	Value int       `json:"value"`
}

// DomainEntityListingEvaluationForm describes ...
type DomainEntityListingEvaluationForm struct {
	Entities    []*EvaluationForm `json:"entities"`
	Total       int               `json:"total"`
	PageSize    int               `json:"pageSize"`
	PageNumber  int               `json:"pageNumber"`
	PageCount   int               `json:"pageCount"`
	FirstUri    string            `json:"firstUri"`
	SelfUri     string            `json:"selfUri"`
	PreviousUri string            `json:"previousUri"`
	NextUri     string            `json:"nextUri"`
	LastUri     string            `json:"lastUri"`
}

// Calibration  describe a Calibration
type Calibration struct {
	ID              uuid.UUID       `json:"id"`
	Name            string          `json:"name"`
	SelfURI         URI             `json:"selfUri"`
	Calibrator      *User           `json:"calibrator"`
	Agent           *User           `json:"agent"`
	Conversation    *Conversation   `json:"conversation"`
	EvaluationForm  *EvaluationForm `json:"evaluationForm"`
	ContextID       string          `json:"contextId"`
	AverageScore    int             `json:"averageScore"`
	HighScore       int             `json:"highScore"`
	LowScore        int             `json:"lowScore"`
	CreatedDate     time.Time       `json:"createdDate"`
	Evaluations     []*Evaluation   `json:"evaluations"`
	Evaluators      []*User         `json:"evaluators"`
	ScoringIndex    *Evaluation     `json:"scoringIndex"`
	ExpertEvaluator *User           `json:"expertEvaluator"`
}

// EvaluationScoringSet  describes an Evaluation Scoring Set
type EvaluationScoringSet struct {
	TotalScore         float64 `json:"totalScore"`
	TotalCriticalScore float64 `json:"totalCriticalScore"`
	//...
}

// GetID gets the identifier of this
//
//	implements Identifiable
func (evaluation Evaluation) GetID() uuid.UUID {
	return evaluation.ID
}

// String gets a string version
//
//	implements the fmt.Stringer interface
func (evaluation Evaluation) String() string {
	if len(evaluation.Name) != 0 {
		return evaluation.Name
	}
	return evaluation.ID.String()
}
