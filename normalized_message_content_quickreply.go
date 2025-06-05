package gcloudcx

import (
	"encoding/json"
	"net/url"

	"github.com/gildas/go-core"
	"github.com/gildas/go-errors"
)

type NormalizedMessageQuickReplyContent struct {
	Text     string   `json:"text"`
	Payload  string   `json:"payload"`
	ImageURL *url.URL `json:"image,omitempty"`
	Action   string   `json:"action,omitempty"` // Message
}

func init() {
	normalizedMessageContentRegistry.Add(NormalizedMessageQuickReplyContent{})
}

// GetType tells the type of this OpenMessageContent
//
// implements core.TypeCarrier
func (quickReply NormalizedMessageQuickReplyContent) GetType() string {
	return "QuickReply"
}

// MarshalJSON marshals this into JSON
//
// implements json.Marshaler
func (quickReply NormalizedMessageQuickReplyContent) MarshalJSON() ([]byte, error) {
	type surrogate NormalizedMessageQuickReplyContent
	type QuickReply struct {
		surrogate
		ImageURL *core.URL `json:"image,omitempty"`
	}
	data, err := json.Marshal(struct {
		ContentType string     `json:"contentType"`
		QuickReply  QuickReply `json:"quickReply"`
	}{
		ContentType: quickReply.GetType(),
		QuickReply: QuickReply{
			surrogate: surrogate(quickReply),
			ImageURL:  (*core.URL)(quickReply.ImageURL),
		},
	})
	return data, errors.JSONMarshalError.Wrap(err)
}

// UnmarshalJSON unmarshals JSON into this
//
// implements json.Unmarshaler
func (quickReply *NormalizedMessageQuickReplyContent) UnmarshalJSON(payload []byte) (err error) {
	type surrogate NormalizedMessageQuickReplyContent
	type QuickReply struct {
		surrogate
		ImageURL *core.URL `json:"image"`
	}
	var inner struct {
		QuickReply QuickReply `json:"quickReply"`
	}
	if err = json.Unmarshal(payload, &inner); err != nil {
		return errors.JSONUnmarshalError.Wrap(err)
	}
	*quickReply = NormalizedMessageQuickReplyContent(inner.QuickReply.surrogate)
	quickReply.ImageURL = (*url.URL)(inner.QuickReply.ImageURL)
	return nil
}
