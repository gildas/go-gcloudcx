package gcloudcx

import (
	"encoding/json"
	"net/url"

	"github.com/gildas/go-core"
	"github.com/gildas/go-errors"
)

type OpenMessageQuickReplyContent struct {
	Text     string   `json:"text"`
	Payload  string   `json:"payload"`
	ImageURL *url.URL `json:"image,omitempty"`
	Action   string   `json:"action,omitempty"` // Message
}

func init() {
	openMessageContentRegistry.Add(OpenMessageQuickReplyContent{})
}

// GetType tells the type of this OpenMessageContent
//
// implements core.TypeCarrier
func (quickReply OpenMessageQuickReplyContent) GetType() string {
	return "QuickReply"
}

// MarshalJSON marshals this into JSON
//
// implements json.Marshaler
func (quickReply OpenMessageQuickReplyContent) MarshalJSON() ([]byte, error) {
	type surrogate OpenMessageQuickReplyContent
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
func (quickReply *OpenMessageQuickReplyContent) UnmarshalJSON(payload []byte) (err error) {
	type surrogate OpenMessageQuickReplyContent
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
	*quickReply = OpenMessageQuickReplyContent(inner.QuickReply.surrogate)
	quickReply.ImageURL = (*url.URL)(inner.QuickReply.ImageURL)
	return nil
}
