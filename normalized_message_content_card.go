package gcloudcx

import (
	"encoding/json"
	"net/url"

	"github.com/gildas/go-core"
	"github.com/gildas/go-errors"
)

// NormalizedMessageCardContent describes the content of a Card
type NormalizedMessageCardContent struct {
	Title         string                        `json:"title"`
	Description   string                        `json:"description,omitempty"`
	ImageURL      *url.URL                      `json:"image,omitempty"`
	VideoURL      *url.URL                      `json:"video,omitempty"`
	DefaultAction *NormalizedMessageCardAction  `json:"defaultAction,omitempty"`
	Actions       []NormalizedMessageCardAction `json:"actions,omitempty"`
}

func init() {
	normalizedMessageContentRegistry.Add(NormalizedMessageCardContent{})
}

// GetType tells the type of this OpenMessageContent
//
// implements core.TypeCarrier
func (card NormalizedMessageCardContent) GetType() string {
	return "Card"
}

// MarshalJSON marshals this into JSON
//
// implements json.Marshaler
func (card NormalizedMessageCardContent) MarshalJSON() ([]byte, error) {
	type surrogate NormalizedMessageCardContent
	type Card struct {
		surrogate
		ImageURL *core.URL `json:"image,omitempty"`
		VideoURL *core.URL `json:"video,omitempty"`
	}
	data, err := json.Marshal(struct {
		ContentType string `json:"contentType"`
		Card        Card   `json:"card"`
	}{
		ContentType: card.GetType(),
		Card: Card{
			surrogate: surrogate(card),
			ImageURL:  (*core.URL)(card.ImageURL),
			VideoURL:  (*core.URL)(card.VideoURL),
		},
	})
	return data, errors.JSONMarshalError.Wrap(err)
}

// UnmarshalJSON unmarshals JSON into this
//
// implements json.Unmarshaler
func (card *NormalizedMessageCardContent) UnmarshalJSON(payload []byte) (err error) {
	type surrogate NormalizedMessageCardContent
	type Card struct {
		surrogate
		ImageURL *core.URL `json:"image"`
		VideoURL *core.URL `json:"video"`
	}
	var inner struct {
		Card Card `json:"card"`
	}
	if err = json.Unmarshal(payload, &inner); err != nil {
		return errors.JSONUnmarshalError.Wrap(err)
	}
	*card = NormalizedMessageCardContent(inner.Card.surrogate)
	card.ImageURL = (*url.URL)(inner.Card.ImageURL)
	card.VideoURL = (*url.URL)(inner.Card.VideoURL)
	return
}
