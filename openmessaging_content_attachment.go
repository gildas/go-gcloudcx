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

// OpenMessageAttachmentContent describes an Attachment Content for an OpenMessage
type OpenMessageAttachmentContent struct {
	ID        string   `json:"id,omitempty"`
	MediaType string   `json:"mediaType"` // Audio, File, Image, Link, Video
	URL       *url.URL `json:"-"`
	Mime      string   `json:"mime,omitempty"`
	Filename  string   `json:"filename,omitempty"`
	Length    uint64   `json:"contentSizeBytes,omitempty"`
	Text      string   `json:"text,omitempty"`
	Hash      string   `json:"sha256,omitempty"`
}

func init() {
	openMessageContentRegistry.Add(OpenMessageAttachmentContent{})
}

// GetType tells the type of this OpenMessageContent
//
// implements core.TypeCarrier
func (attachment OpenMessageAttachmentContent) GetType() string {
	return "Attachment"
}

// WithContent sets the content of this Attachment from a request.Content
func (attachment OpenMessageAttachmentContent) WithContent(content *request.Content) *OpenMessageAttachmentContent {
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

	attachment.MediaType = attachmentType
	attachment.Mime = content.Type
	attachment.Filename = content.Name
	attachment.URL = content.URL
	attachment.Length = content.Length

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
func (attachment OpenMessageAttachmentContent) MarshalJSON() ([]byte, error) {
	type surrogate OpenMessageAttachmentContent
	type Attachment struct {
		surrogate
		URL *core.URL `json:"url"`
	}
	data, err := json.Marshal(struct {
		ContentType string     `json:"contentType"`
		Attachment  Attachment `json:"attachment"`
	}{
		ContentType: attachment.GetType(),
		Attachment: Attachment{
			surrogate: surrogate(attachment),
			URL:       (*core.URL)(attachment.URL),
		},
	})
	return data, errors.JSONMarshalError.Wrap(err)
}

// UnmarshalJSON unmarshals JSON into this
func (attachment *OpenMessageAttachmentContent) UnmarshalJSON(payload []byte) (err error) {
	type surrogate OpenMessageAttachmentContent
	var inner struct {
		ContentType string `json:"contentType"`
		Attachment  struct {
			surrogate
			URL *core.URL `json:"url"`
		} `json:"attachment"`
	}

	if err = json.Unmarshal(payload, &inner); err != nil {
		return errors.JSONUnmarshalError.Wrap(err)
	}
	*attachment = OpenMessageAttachmentContent(inner.Attachment.surrogate)
	attachment.URL = (*url.URL)(inner.Attachment.URL)
	validMediaTypes := []string{"Audio", "File", "Image", "Link", "Video"}
	if !core.Contains(validMediaTypes, attachment.MediaType) {
		return errors.JSONUnmarshalError.Wrap(errors.InvalidType.With(attachment.MediaType, strings.Join(validMediaTypes, ", ")))
	}
	return
}
