package gcloudcx

import (
	"time"
)

// RoutingStatus describes a Routing Status
type RoutingStatus struct {
	UserID    string    `json:"userId"`
	Status    string    `json:"status"` // OFF_QUEUE, IDLE, INTERACTING, NOT_RESPONDING, COMMUNICATING
	StartTime time.Time `json:"startTime"`
}
