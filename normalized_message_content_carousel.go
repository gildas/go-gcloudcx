package gcloudcx

import (
	"encoding/json"
	"net/url"

	"github.com/gildas/go-core"
	"github.com/gildas/go-errors"
)

// NormalizedMessageCarouselContent describes the content of a Carousel
type NormalizedMessageCarouselContent struct {
	Cards []NormalizedMessageCarouselCard `json:"cards"`
}

// NormalizedMessageCarouselCard describes the content of a Card in a Carousel
type NormalizedMessageCarouselCard struct {
	Title         string                        `json:"title"`
	Description   string                        `json:"description,omitempty"`
	ImageURL      *url.URL                      `json:"image,omitempty"`
	VideoURL      *url.URL                      `json:"video,omitempty"`
	DefaultAction *NormalizedMessageCardAction  `json:"defaultAction,omitempty"`
	Actions       []NormalizedMessageCardAction `json:"actions,omitempty"`
}

func init() {
	normalizedMessageContentRegistry.Add(NormalizedMessageCarouselContent{})
}

// GetType tells the type of this OpenMessageContent
//
// implements core.TypeCarrier
func (carousel NormalizedMessageCarouselContent) GetType() string {
	return "Carousel"
}

// MarshalJSON marshals this into JSON
//
// implements json.Marshaler
func (carousel NormalizedMessageCarouselContent) MarshalJSON() ([]byte, error) {
	type surrogate NormalizedMessageCarouselContent
	type Carousel struct {
		surrogate
	}
	data, err := json.Marshal(struct {
		ContentType string   `json:"contentType"`
		Carousel    Carousel `json:"carousel"`
	}{
		ContentType: carousel.GetType(),
		Carousel: Carousel{
			surrogate: surrogate(carousel),
		},
	})
	return data, errors.JSONMarshalError.Wrap(err)
}

// UnmarshalJSON unmarshals JSON into this
//
// implements json.Unmarshaler
func (carousel *NormalizedMessageCarouselContent) UnmarshalJSON(payload []byte) (err error) {
	type surrogate NormalizedMessageCarouselContent
	type Carousel struct {
		surrogate
	}
	var inner struct {
		ContentType string   `json:"contentType"`
		Carousel    Carousel `json:"carousel"`
	}
	if err = json.Unmarshal(payload, &inner); err != nil {
		return errors.JSONUnmarshalError.Wrap(err)
	}
	if inner.ContentType != carousel.GetType() {
		return errors.JSONUnmarshalError.Wrap(errors.InvalidType.With("contentType", carousel.GetType()))
	}
	*carousel = NormalizedMessageCarouselContent(inner.Carousel.surrogate)
	return nil
}

// / MarshalJSON marshals this into JSON
//
// implements json.Marshaler
func (card NormalizedMessageCarouselCard) MarshalJSON() ([]byte, error) {
	type surrogate NormalizedMessageCarouselCard
	data, err := json.Marshal(struct {
		surrogate
		ImageURL *core.URL `json:"image,omitempty"`
		VideoURL *core.URL `json:"video,omitempty"`
	}{
		surrogate: surrogate(card),
		ImageURL:  (*core.URL)(card.ImageURL),
		VideoURL:  (*core.URL)(card.VideoURL),
	})
	return data, errors.JSONMarshalError.Wrap(err)
}

// UnmarshalJSON unmarshals JSON into this
//
// implements json.Unmarshaler
func (card *NormalizedMessageCarouselCard) UnmarshalJSON(payload []byte) (err error) {
	type surrogate NormalizedMessageCarouselCard
	var inner struct {
		surrogate
		ImageURL *core.URL `json:"image"`
		VideoURL *core.URL `json:"video"`
	}
	if err = json.Unmarshal(payload, &inner); err != nil {
		return errors.JSONUnmarshalError.Wrap(err)
	}
	*card = NormalizedMessageCarouselCard(inner.surrogate)
	card.ImageURL = (*url.URL)(inner.ImageURL)
	card.VideoURL = (*url.URL)(inner.VideoURL)
	return
}
