package gcloudcx

import (
	"encoding/json"
	"net/url"
	"time"

	"github.com/gildas/go-core"
	"github.com/gildas/go-errors"
)

// NormalizedMessageDatePickerContent describes the content of a DatePicker
type NormalizedMessageDatePickerContent struct {
	Title          string                     `json:"title,omitempty"`
	Subtitle       string                     `json:"subtitle,omitempty"`
	ImageURL       *url.URL                   `json:"imageUrl,omitempty"`
	MinDate        *time.Time                 `json:"dateMinimum,omitempty"`
	MaxDate        *time.Time                 `json:"dateMaximum,omitempty"`
	AvailableTimes []OpenMessageAvailableTime `json:"availableTimes"`
}

// OpenMessageAvailableTime describes the content of an available time within a DatePicker
type OpenMessageAvailableTime struct {
	Time     time.Time     `json:"dateTime"`
	Duration time.Duration `json:"duration"`
}

func init() {
	normalizedMessageContentRegistry.Add(NormalizedMessageDatePickerContent{})
}

// GetType tells the type of this OpenMessageContent
//
// implements core.TypeCarrier
func (datePicker NormalizedMessageDatePickerContent) GetType() string {
	return "DatePicker"
}

// MarshalJSON marshals this into JSON
//
// implements json.marshaler
func (availableTime OpenMessageAvailableTime) MarshalJSON() ([]byte, error) {
	type surrogate OpenMessageAvailableTime
	data, err := json.Marshal(struct {
		surrogate
		Time     core.Time `json:"dateTime"`
		Duration uint64    `json:"duration"`
	}{
		surrogate: surrogate(availableTime),
		Time:      core.Time(availableTime.Time),
		Duration:  uint64(availableTime.Duration.Seconds()),
	})
	return data, errors.JSONMarshalError.Wrap(err)
}

// MarshalJSON marshals this into JSON
//
// implements json.Marshaler
func (datePicker NormalizedMessageDatePickerContent) MarshalJSON() ([]byte, error) {
	type surrogate NormalizedMessageDatePickerContent
	type DatePicker struct {
		surrogate
		ImageURL *core.URL  `json:"imageUrl,omitempty"`
		MinDate  *core.Time `json:"dateMinimum,omitempty"`
		MaxDate  *core.Time `json:"dateMaximum,omitempty"`
	}
	data, err := json.Marshal(struct {
		ContentType string     `json:"contentType"`
		DatePicker  DatePicker `json:"datePicker"`
	}{
		ContentType: datePicker.GetType(),
		DatePicker: DatePicker{
			surrogate: surrogate(datePicker),
			ImageURL:  (*core.URL)(datePicker.ImageURL),
			MinDate:   (*core.Time)(datePicker.MinDate),
			MaxDate:   (*core.Time)(datePicker.MaxDate),
		},
	})
	return data, errors.JSONMarshalError.Wrap(err)
}

// UnmarshalJSON unmarshals this from JSON
//
// implements json.Unmarshaler
func (availableTime *OpenMessageAvailableTime) UnmarshalJSON(payload []byte) (err error) {
	type surrogate OpenMessageAvailableTime
	var inner struct {
		surrogate
		Time     core.Time `json:"dateTime"`
		Duration uint64    `json:"duration"`
	}
	if err = json.Unmarshal(payload, &inner); err != nil {
		return errors.JSONUnmarshalError.Wrap(err)
	}
	*availableTime = OpenMessageAvailableTime(inner.surrogate)
	availableTime.Time = time.Time(inner.Time)
	availableTime.Duration = time.Duration(inner.Duration) * time.Second
	return nil
}

// UnmarshalJSON unmarshals this from JSON
//
// implements json.Unmarshaler
func (datePicker *NormalizedMessageDatePickerContent) UnmarshalJSON(payload []byte) (err error) {
	type surrogate NormalizedMessageDatePickerContent
	type DatePicker struct {
		surrogate
		ImageURL *core.URL  `json:"imageUrl,omitempty"`
		MinDate  *core.Time `json:"dateMinimum,omitempty"`
		MaxDate  *core.Time `json:"dateMaximum,omitempty"`
	}
	var inner struct {
		Type       string     `json:"contentType"`
		DatePicker DatePicker `json:"datePicker"`
	}
	if err = json.Unmarshal(payload, &inner); err != nil {
		return errors.JSONUnmarshalError.WrapIfNotMe(err)
	}
	if inner.Type != datePicker.GetType() {
		return errors.InvalidType.With(inner.Type, datePicker.GetType())
	}
	*datePicker = NormalizedMessageDatePickerContent(inner.DatePicker.surrogate)
	datePicker.ImageURL = (*url.URL)(inner.DatePicker.ImageURL)
	datePicker.MinDate = (*time.Time)(inner.DatePicker.MinDate)
	datePicker.MaxDate = (*time.Time)(inner.DatePicker.MaxDate)
	return nil
}
