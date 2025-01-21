package gcloudcx

import (
	"encoding/json"
	"net/url"

	"github.com/gildas/go-core"
	"github.com/gildas/go-errors"
)

// OpenMessageCardContent describes the content of a Card
type OpenMessageCardContent struct {
	Title         string                  `json:"title"`
	Description   string                  `json:"description,omitempty"`
	ImageURL      *url.URL                `json:"image,omitempty"`
	VideoURL      *url.URL                `json:"video,omitempty"`
	DefaultAction *OpenMessageCardAction  `json:"defaultAction,omitempty"`
	Actions       []OpenMessageCardAction `json:"actions,omitempty"`
}

func init() {
	openMessageContentRegistry.Add(OpenMessageCardContent{})
}

// GetType tells the type of this OpenMessageContent
//
// implements core.TypeCarrier
func (card OpenMessageCardContent) GetType() string {
	return "Card"
}

// MarshalJSON marshals this into JSON
//
// implements json.Marshaler
func (card OpenMessageCardContent) MarshalJSON() ([]byte, error) {
	type surrogate OpenMessageCardContent
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
func (card *OpenMessageCardContent) UnmarshalJSON(payload []byte) (err error) {
	type surrogate OpenMessageCardContent
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
	*card = OpenMessageCardContent(inner.Card.surrogate)
	card.ImageURL = (*url.URL)(inner.Card.ImageURL)
	card.VideoURL = (*url.URL)(inner.Card.VideoURL)
	return
}
