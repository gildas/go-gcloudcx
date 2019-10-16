package purecloud

import (
	"time"
)

// OutOfOffice describes the Out Of Office status
type OutOfOffice struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	User       *User     `json:"user"`
	Active     bool      `json:"active"`
	Indefinite bool      `json:"indefinite"`
	StartDate  time.Time `json:"startDate"`
	EndDate    time.Time `json:"endDate"`
	SelfURI    string    `json:"selfUri"`
}