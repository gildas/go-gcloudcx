package gcloudcx

import (
	"encoding/json"
	"mime"
	"net/url"
	"path"
	"strings"

	"github.com/gildas/go-core"
	"github.com/gildas/go-errors"
	"github.com/gildas/go-request"
	nanoid "github.com/matoous/go-nanoid/v2"
)

type OpenMessageAttachment struct {
	ID       string   `json:"id,omitempty"`
	Type     string   `json:"mediaType"`
	URL      *url.URL `json:"-"`
	Mime     string   `json:"mime,omitempty"`
	Filename string   `json:"filename,omitempty"`
	Text     string   `json:"text,omitempty"`
	Hash     string   `json:"sha256,omitempty"`
}

func (attachment OpenMessageAttachment) WithContent(content *request.Content) *OpenMessageAttachment {
	var attachmentType string
	switch {
	case len(content.Type) == 0:
		attachmentType = "Link"
	case strings.HasPrefix(content.Type, "audio"):
		attachmentType = "Audio"
	case strings.HasPrefix(content.Type, "image"):
		attachmentType = "Image"
	case strings.HasPrefix(content.Type, "video"):
		attachmentType = "Video"
	default:
		attachmentType = "File"
	}

	attachment.Type = attachmentType
	attachment.Mime = content.Type
	attachment.Filename = content.Name
	attachment.URL = content.URL

	if attachmentType != "Link" && len(content.Name) == 0 {
		fileExtension := path.Ext(content.URL.Path)
		if content.Type == "audio/mpeg" {
			fileExtension = ".mp3"
		} else if fileExtensions, err := mime.ExtensionsByType(content.Type); err == nil && len(fileExtensions) > 0 {
			fileExtension = fileExtensions[0]
		}
		fileID, _ := nanoid.New()
		attachment.Filename = strings.ToLower(attachmentType) + "-" + fileID + fileExtension
	}

	return &attachment
}

// MarshalJSON marshals this into JSON
func (attachment OpenMessageAttachment) MarshalJSON() ([]byte, error) {
	type surrogate OpenMessageAttachment
	data, err := json.Marshal(struct {
		surrogate
		U *core.URL `json:"url"`
	}{
		surrogate: surrogate(attachment),
		U:         (*core.URL)(attachment.URL),
	})
	return data, errors.JSONMarshalError.Wrap(err)
}

// UnmarshalJSON unmarshals JSON into this
func (attachment *OpenMessageAttachment) UnmarshalJSON(payload []byte) (err error) {
	type surrogate OpenMessageAttachment
	var inner struct {
		surrogate
		U *core.URL `json:"url"`
	}

	if err = json.Unmarshal(payload, &inner); err != nil {
		return errors.JSONUnmarshalError.Wrap(err)
	}
	*attachment = OpenMessageAttachment(inner.surrogate)
	attachment.URL = (*url.URL)(inner.U)
	return
}
