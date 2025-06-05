package gcloudcx

import (
	"encoding/json"

	"github.com/gildas/go-errors"
)

// NormalizedMessageQuickReplyV2Content describes the content of a QuickReplyV2
type NormalizedMessageQuickReplyV2Content struct {
	Title   string                                `json:"title"`
	Actions []NormalizedMessageQuickReplyV2Action `json:"actions"`
}

func init() {
	normalizedMessageContentRegistry.Add(NormalizedMessageQuickReplyV2Content{})
}

// GetType tells the type of this OpenMessageContent
//
// implements core.TypeCarrier
func (quickReply NormalizedMessageQuickReplyV2Content) GetType() string {
	return "QuickReplyV2"
}

// MarshalJSON marshals this into JSON
//
// implements json.Marshaler
func (quickReply NormalizedMessageQuickReplyV2Content) MarshalJSON() ([]byte, error) {
	type surrogate NormalizedMessageQuickReplyV2Content
	type QuickReplyV2 struct {
		surrogate
	}
	data, err := json.Marshal(struct {
		ContentType  string       `json:"contentType"`
		QuickReplyV2 QuickReplyV2 `json:"quickReplyV2"`
	}{
		ContentType: quickReply.GetType(),
		QuickReplyV2: QuickReplyV2{
			surrogate: surrogate(quickReply),
		},
	})
	return data, errors.JSONMarshalError.Wrap(err)
}

// UnmarshalJSON unmarshals JSON into this
//
// implements json.Unmarshaler
func (quickReply *NormalizedMessageQuickReplyV2Content) UnmarshalJSON(payload []byte) (err error) {
	type surrogate NormalizedMessageQuickReplyV2Content
	type QuickReplyV2 struct {
		surrogate
	}
	var inner struct {
		QuickReplyV2 QuickReplyV2 `json:"quickReplyV2"`
	}
	if err = json.Unmarshal(payload, &inner); err != nil {
		return errors.JSONUnmarshalError.Wrap(err)
	}
	*quickReply = NormalizedMessageQuickReplyV2Content(inner.QuickReplyV2.surrogate)
	return
}
