package purecloud

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

// EvaluationForm describes an Evaluation Form
type EvaluationForm struct {
	ID      uuid.UUID `json:"id"`
	Name    string    `json:"name"`
	SelfURI string    `json:"selfUri"`

	ModifiedDate time.Time `json:"modifiedDate"`
	Published    bool      `json:"published"`
	ContextID    string    `json:"contextId"`

	QuestionGroups    []EvaluationQuestionGroup         `json:"questionGroups"`
	PublishedVersions DomainEntityListingEvaluationForm `json:"publishedVersions"`
}

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

// VisibilityCondition  describes visibility conditions
type VisibilityCondition struct {
	CombiningOperation string        `json:"combiningOperation"`
	Predicates         []interface{} `json:"predicates"`
}

// AnswerOption describes an Answer Option
type AnswerOption struct {
	ID    uuid.UUID `json:"id"`
	Text  string `json:"text"`
	Value int    `json:"value"`
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
	SelfUri         string          `json:"selfUri"`
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
//   implements Identifiable
func (evaluation Evaluation) GetID() uuid.UUID {
	return evaluation.ID
}

// String gets a string version
//   implements the fmt.Stringer interface
func (evaluation Evaluation) String() string {
	if len(evaluation.Name) != 0 {
		return evaluation.Name
	}
	return evaluation.ID.String()
}
