package gcloudcx

type RecordingChatMessage struct {
	ID                    string `json:"id"`
	To                    string `json:"to"`
	From                  string `json:"from"`
	Body                  string `json:"body"`
	UTC                   string `json:"utc"`
	Chat                  string `json:"chat"`
	Message               string `json:"message"`
	Type                  string `json:"type"`
	BodyType              string `json:"bodyType"`
	SenderCommunicationID string `json:"senderCommunicationId"`
	ParticipantPurpose    string `json:"participantPurpose"`
	User                  *User  `json:"user,omitempty"`
}
