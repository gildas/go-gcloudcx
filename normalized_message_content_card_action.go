package gcloudcx

import (
	"fmt"
	"strings"

	"github.com/gildas/go-core"
	"github.com/gildas/go-errors"
)

// NormalizedMessageCardAction describes an action in a Card
type NormalizedMessageCardAction interface {
	core.TypeCarrier
	fmt.Stringer
}

var cardActionTypeRegistry = core.TypeRegistry{}

// UnmarshalMessageCardAction unmarshals a JSON payload into a NormalizedMessageCardAction
func UnmarshalMessageCardAction(payload []byte) (NormalizedMessageCardAction, error) {
	action, err := cardActionTypeRegistry.UnmarshalJSON(payload, "type")
	if err == nil {
		return action.(NormalizedMessageCardAction), nil
	}
	if strings.HasPrefix(err.Error(), "Missing JSON Property") {
		return nil, errors.JSONUnmarshalError.Wrap(errors.ArgumentMissing.With("type"))
	}
	if strings.HasPrefix(err.Error(), "Unsupported Type") {
		return nil, errors.JSONUnmarshalError.Wrap(
			errors.InvalidType.With(
				strings.TrimSuffix(strings.TrimPrefix(err.Error(), `Unsupported Type "`), `"`),
				strings.Join(cardActionTypeRegistry.SupportedTypes(), ","),
			),
		)
	}
	if errors.Is(err, errors.JSONUnmarshalError) {
		return nil, err
	}
	return nil, errors.JSONUnmarshalError.Wrap(err)
}
