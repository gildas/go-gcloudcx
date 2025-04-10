package gcloudcx

import (
	"encoding/json"
	"net/url"

	"github.com/gildas/go-core"
	"github.com/gildas/go-errors"
)

// OpenMessageCarouselContent describes the content of a Carousel
type OpenMessageCarouselContent struct {
	Cards []OpenMessageCarouselCard `json:"-"`
}

// OpenMessageCarouselCard describes the content of a Card in a Carousel
type OpenMessageCarouselCard struct {
	Title         string                  `json:"title"`
	Description   string                  `json:"description,omitempty"`
	ImageURL      *url.URL                `json:"image,omitempty"`
	VideoURL      *url.URL                `json:"video,omitempty"`
	DefaultAction *OpenMessageCardAction  `json:"defaultAction,omitempty"`
	Actions       []OpenMessageCardAction `json:"actions,omitempty"`
}

func init() {
	openMessageContentRegistry.Add(OpenMessageCarouselContent{})
}

// GetType tells the type of this OpenMessageContent
//
// implements core.TypeCarrier
func (carousel OpenMessageCarouselContent) GetType() string {
	return "Carousel"
}

// MarshalJSON marshals this into JSON
//
// implements json.Marshaler
func (carousel OpenMessageCarouselContent) MarshalJSON() ([]byte, error) {
	type surrogate OpenMessageCarouselContent
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
func (carousel *OpenMessageCarouselContent) UnmarshalJSON(payload []byte) (err error) {
	type surrogate OpenMessageCarouselContent
	type Carousel struct {
		Cards []OpenMessageCarouselCard `json:"cards"`
	}
	var inner struct {
		ContentType string   `json:"contentType"`
		Carousel    Carousel `json:"carousel"`
		surrogate
	}
	if err = json.Unmarshal(payload, &inner); err != nil {
		return errors.JSONUnmarshalError.Wrap(err)
	}
	if inner.ContentType != carousel.GetType() {
		return errors.JSONUnmarshalError.Wrap(errors.InvalidType.With("contentType", carousel.GetType()))
	}
	*carousel = OpenMessageCarouselContent(inner.surrogate)
	carousel.Cards = append(carousel.Cards, inner.Carousel.Cards...)
	return nil
}

// UnmarshalJSON unmarshals JSON into this
//
// implements json.Unmarshaler
func (card *OpenMessageCarouselCard) UnmarshalJSON(payload []byte) (err error) {
	type surrogate OpenMessageCarouselCard
	var inner struct {
		surrogate
		ImageURL *core.URL `json:"image,omitempty"`
		VideoURL *core.URL `json:"video,omitempty"`
	}
	if err = json.Unmarshal(payload, &inner); err != nil {
		return errors.JSONUnmarshalError.Wrap(err)
	}
	*card = OpenMessageCarouselCard(inner.surrogate)
	card.ImageURL = (*url.URL)(inner.ImageURL)
	card.VideoURL = (*url.URL)(inner.VideoURL)
	return
}
