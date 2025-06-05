package gcloudcx

import (
	"encoding/json"
	"net/url"

	"github.com/gildas/go-core"
	"github.com/gildas/go-errors"
)

type NormalizedMessageStoryContent struct {
	StoryType string   `json:"type"` // "Mention", "Reply"
	URL       *url.URL `json:"url,omitempty"`
	ReplyToID string   `json:"replyToId,omitempty"`
}

func init() {
	normalizedMessageContentRegistry.Add(NormalizedMessageStoryContent{})
}

// GetType tells the type of this OpenMessageContent
//
// implements core.TypeCarrier
func (story NormalizedMessageStoryContent) GetType() string {
	return "Story"
}

// MarshalJSON marshals this into JSON
//
// implements json.Marshaler
func (story NormalizedMessageStoryContent) MarshalJSON() ([]byte, error) {
	type surrogate NormalizedMessageStoryContent
	type Story struct {
		surrogate
		URL *core.URL `json:"url,omitempty"`
	}
	data, err := json.Marshal(struct {
		ContentType string `json:"contentType"`
		Story       Story  `json:"story"`
	}{
		ContentType: story.GetType(),
		Story: Story{
			surrogate: surrogate(story),
			URL:       (*core.URL)(story.URL),
		},
	})
	return data, errors.JSONMarshalError.Wrap(err)
}

// UnmarshalJSON unmarshals JSON into this
//
// implements json.Unmarshaler
func (story *NormalizedMessageStoryContent) UnmarshalJSON(payload []byte) (err error) {
	type surrogate NormalizedMessageStoryContent
	type Story struct {
		surrogate
		URL *core.URL `json:"url"`
	}
	var inner struct {
		Story Story `json:"story"`
	}
	if err = json.Unmarshal(payload, &inner); err != nil {
		return errors.JSONUnmarshalError.Wrap(err)
	}
	*story = NormalizedMessageStoryContent(inner.Story.surrogate)
	story.URL = (*url.URL)(inner.Story.URL)
	return nil
}
