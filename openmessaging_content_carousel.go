package gcloudcx

import (
	"encoding/json"

	"github.com/gildas/go-errors"
)

// OpenMessageCarouselContent describes the content of a Carousel
type OpenMessageCarouselContent struct {
	Cards []OpenMessageCardContent `json:"cards"`
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
