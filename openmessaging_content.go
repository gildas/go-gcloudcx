package gcloudcx

import (
	"strings"

	"github.com/gildas/go-core"
	"github.com/gildas/go-errors"
)

type OpenMessageContent interface {
	core.TypeCarrier
}

var openMessageContentRegistry = core.TypeRegistry{}

func UnmarshalOpenMessageContent(payload []byte) (OpenMessageContent, error) {
	content, err := openMessageContentRegistry.UnmarshalJSON(payload, "contentType")
	if err == nil {
		return content.(OpenMessageContent), nil
	}
	if strings.HasPrefix(err.Error(), "Missing JSON Property") {
		return nil, errors.JSONUnmarshalError.Wrap(errors.ArgumentMissing.With("contentType"))
	}
	// if !Contains([]string{"Attachment", "Location", "QuickReply", "ButtonResponse", "Notification", "GenericTemplate", "ListTemplate", "Postback", "Reactions", "Mention"}, content.Type) {
	if strings.HasPrefix(err.Error(), "Unsupported Type") {
		supportedTypes := make([]string, 0, len(openMessageContentRegistry))
		for key := range openMessageContentRegistry {
			supportedTypes = append(supportedTypes, key)
		}
		return nil, errors.JSONUnmarshalError.Wrap(
			errors.InvalidType.With(
				strings.TrimSuffix(strings.TrimPrefix(err.Error(), `Unsupported Type "`), `"`),
				strings.Join(supportedTypes, ","),
			),
		)
	}
	if errors.Is(err, errors.JSONUnmarshalError) {
		return nil, err
	}
	return nil, errors.JSONUnmarshalError.Wrap(err)
}
