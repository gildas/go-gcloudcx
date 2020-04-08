package purecloud

// MediaParticipantRequest describes a request Media Participant
type MediaParticipantRequest struct {
	Wrapup        *Wrapup `json:"wrapup,omitempty"`
	State         string  `json:"state,omitempty"`     // alerting, dialing, contacting, offering, connected, disconnected, terminated, converting, uploading, transmitting, none
	Recording     bool    `json:"recording,omitempty"`
	Muted         bool    `json:"muted,omitempty"`
	Confined      bool    `json:"confined,omitempty"`
	Held          bool    `json:"held,omitempty"`
	WrapupSkipped bool    `json:"wrapupSkipped,omitempty"`
}