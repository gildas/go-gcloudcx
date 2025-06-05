package gcloudcx

import (
	"strings"

	"github.com/gildas/go-core"
	"github.com/gildas/go-errors"
)

type NormalizedMessageContent interface {
	core.TypeCarrier
}

var normalizedMessageContentRegistry = core.TypeRegistry{}

func UnmarshalOpenMessageContent(payload []byte) (NormalizedMessageContent, error) {
	content, err := normalizedMessageContentRegistry.UnmarshalJSON(payload, "contentType")
	if err == nil {
		return content.(NormalizedMessageContent), nil
	}
	if strings.HasPrefix(err.Error(), "Missing JSON Property") {
		return nil, errors.JSONUnmarshalError.Wrap(errors.ArgumentMissing.With("contentType"))
	}
	if strings.HasPrefix(err.Error(), "Unsupported Type") {
		return nil, errors.JSONUnmarshalError.Wrap(
			errors.InvalidType.With(
				strings.TrimSuffix(strings.TrimPrefix(err.Error(), `Unsupported Type "`), `"`),
				strings.Join(normalizedMessageContentRegistry.SupportedTypes(), ","),
			),
		)
	}
	if errors.Is(err, errors.JSONUnmarshalError) {
		return nil, err
	}
	return nil, errors.JSONUnmarshalError.Wrap(err)
}
