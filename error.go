package gcloudcx

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gildas/go-errors"
)

var (
	// AuthenticationRequestTimeoutError means the request timed out
	AuthenticationRequestTimeoutError = APIError{Status: 504, Code: "authentication.request.timeout", Message: "Authentication request timeout."}
	// BadRequestError means the request was badly formed
	BadRequestError = APIError{Status: 400, Code: "bad.request", Message: "The request could not be understood by the server due to malformed syntax."}
	// InternalServerError means the server experiences an internal error
	InternalServerError = APIError{Status: 500, Code: "internal.server.error", Message: "The server encountered an unexpected condition which prevented it from fulfilling the request."}
	// InvalidDateError means the given date was invalid
	InvalidDateError = APIError{Status: 400, Code: "invalid.date", Message: "Dates must be specified as ISO-8601 strings. For example: yyyy-MM-ddTHH:mm:ss.SSSZ"}
	// InvalidValueError means the value was invalid
	InvalidValueError = APIError{Status: 400, Code: "invalid.value", Message: "Value [%s] is not valid for field type [%s]. Allowable values are: %s"}
	// MissingAnyPermissionsError means the request was missing some permissions
	MissingAnyPermissionsError = APIError{Status: 403, Code: "missing.any.permissions", Message: "Unable to perform the requested action. You must have at least one of the following permissions assigned: %s"}
	// MissingPermissionsError means the request was missing some permissions
	MissingPermissionsError = APIError{Status: 403, Code: "missing.permissions", Message: "Unable to perform the requested action. You are missing the following permission(s): %s"}
	// NotAuthorizedError means the request was not authorized
	NotAuthorizedError = APIError{Status: 403, Code: "not.authorized", Message: "You are not authorized to perform the requested action."}
	// NotFoundError means the wanted resource was not found
	NotFoundError = APIError{Status: 404, Code: "not.found", Message: "The requested resource was not found."}
	// RequestTimeoutError means the request timed out
	RequestTimeoutError = APIError{Status: 504, Code: "request.timeout", Message: "The request timed out."}
	// ServiceUnavailableError means the service is not available
	ServiceUnavailableError = APIError{Status: 503, Code: "service.unavailable", Message: "Service Unavailable - The server is currently unavailable (because it is overloaded or down for maintenance)."}
	// TooManyRequestsError means the client sent too many requests and should wait before sending more
	TooManyRequestsError = APIError{Status: 429, Code: "too.many.requests", Message: "Rate limit exceeded the maximum [%s] requests within [%s] seconds"}
	// UnsupportedMediaTypeError means the media type is not supported
	UnsupportedMediaTypeError = APIError{Status: 415, Code: "unsupported.media.type", Message: "Unsupported Media Type - Unsupported or incorrect media type, such as an incorrect Content-Type value in the header."}

	// AuthenticationRequiredError means the request should authenticate first
	AuthenticationRequiredError = APIError{Status: 401, Code: "authentication.required", Message: "No authentication bearer token specified in authorization header."}
	// BadCredentialsError means the credentials are invalid
	BadCredentialsError = APIError{Status: 401, Code: "bad.credentials", Message: "Invalid login credentials (%s)."}
	// CredentialsExpiredError means the credentials are expired
	CredentialsExpiredError = APIError{Status: 401, Code: "credentials.expired", Message: "The supplied credentials are expired and cannot be used."}

	// ChatConversationStateError  means the conversation does not permit the request
	ChatConversationStateError = APIError{Status: 400, Code: "chat.error.conversation.state", Message: "The conversation is in a state which does not permit this action."}
	// ChatMemberStateError means the chat member does not permit the request
	ChatMemberStateError = APIError{Status: 400, Code: "chat.error.member.state", Message: "The conversation member is in a state which does not permit this action."}
	// ChatDeploymentBadAuthError means the authentication failed
	ChatDeploymentBadAuthError = APIError{Status: 400, Code: "chat.deployment.bad.auth", Message: "The customer member authentication has failed."}
	// ChatDeploymentDisabledError means the deployment is disabled
	ChatDeploymentDisabledError = APIError{Status: 400, Code: "chat.deployment.disabled", Message: "The web chat deployment is currently disabled."}
	// ChatDeploymentRequireAuth means the deployment requires some authentication
	ChatDeploymentRequireAuth = APIError{Status: 400, Code: "chat.deployment.require.auth", Message: "The deployment requires the customer member to be authenticated."}
	// ChatInvalidQueueError means the queue is not valid
	ChatInvalidQueueError = APIError{Status: 400, Code: "chat.error.invalid.queue", Message: "The specified queue is not valid."}
	// ChatCreateConversationRequestRoutingTargetError means the routing target is not valid
	ChatCreateConversationRequestRoutingTargetError = APIError{Status: 400, Code: "chat.error.createconversationrequest.routingtarget", Message: "The routing target is not valid."}
	// MessageExpired means the message has expired
	MessageExpired = APIError{Status: 400, Code: "message.expired", Message: "%s"}
	// RateLimited means the request was rate limited
	RateLimited = APIError{Status: 429, Code: "rate.limited", Message: "%s"}
	// MessageNotAllowed means the message is not allowed
	MessageNotAllowed = APIError{Status: 400, Code: "message.not.allowed", Message: "%s"}
	// GeneralError means the request failed for a general reason
	GeneralError = APIError{Status: 400, Code: "general.error", Message: "%s"}
	// UnsupportedMessage means the message is not supported
	UnsupportedMessage = APIError{Status: 400, Code: "unsupported.message", Message: "%s"}
	// UnknownMessage means the message is unknown
	UnknownMessage = APIError{Status: 400, Code: "unknown.message", Message: "%s"}
	// InvalidMessageStructure means the message structure is invalid
	InvalidMessageStructure = APIError{Status: 400, Code: "invalid.message.structure", Message: "%s"}
	// InvalidDestination means the destination is invalid
	InvalidDestination = APIError{Status: 400, Code: "invalid.destination", Message: "%s"}
	// ServerError means the server failed
	ServerError = APIError{Status: 500, Code: "server.error", Message: "%s"}
	// MediaTypeNotAllowed means the media type is not allowed
	MediaTypeNotAllowed = APIError{Status: 415, Code: "media.type.not.allowed", Message: "%s"}
	// InvalidMediaContentLength means the media content length is invalid
	InvalidMediaContentLength = APIError{Status: 400, Code: "invalid.media.content.length", Message: "%s"}
	// RecipientOptedOut means the recipient opted out
	RecipientOptedOut = APIError{Status: 400, Code: "recipient.opted.out", Message: "%s"}
	// JSONUnmarshalError means the JSON could not be unmarshaled
	JSONUnmarshalError = APIError{Status: http.StatusNotAcceptable, Code: "json.unmarshal", Message: "Unmarshable Data for type %s"}
)

// APIError represents an error from the Gcloud API
type APIError struct {
	Status            int               `json:"status,omitempty"`
	Code              string            `json:"code,omitempty"`
	Message           string            `json:"message,omitempty"`
	MessageParams     map[string]string `json:"messageParams,omitempty"`
	MessageWithParams string            `json:"messageWithParams,omitempty"`
	EntityID          string            `json:"entityId,omitempty"`
	EntityName        string            `json:"entityName,omitempty"`
	ContextID         string            `json:"contextId,omitempty"`
	CorrelationID     string            `json:"correlationId,omitempty"`
	Details           []APIErrorDetails `json:"details,omitempty"`
	Errors            []APIError        `json:"errors,omitempty"`
	Stack             errors.StackTrace `json:"-"`
}

// APIErrorDetails contains the details of an APIError
type APIErrorDetails struct {
	ErrorCode  string `json:"errorCode,omitempty"`
	FieldName  string `json:"fieldName,omitempty"`
	EntityID   string `json:"entityId,omitempty"`
	EntityName string `json:"entityName,omitempty"`
}

// Clone creates an exact copy of this Error
func (e APIError) Clone() *APIError {
	final := e
	return &final
}

// Error returns a string representation of this error
func (e APIError) Error() string {
	var buffer strings.Builder
	if len(e.Code) > 0 {
		buffer.WriteString(fmt.Sprintf("GCloudCX Code %s: ", e.Code))
	}
	if len(e.MessageWithParams) > 0 {
		buffer.WriteString(e.MessageWithParams)
	}
	if len(e.Message) > 0 {
		buffer.WriteString(e.Message)
	}
	if len(e.CorrelationID) > 0 {
		buffer.WriteString(fmt.Sprintf(" (Correlation: %s)", e.CorrelationID))
	}
	return buffer.String()
}

// Is tells if this error matches the target.
//
// implements errors.Is interface (package "errors").
//
// To check if an error is an errors.Error, simply write:
//
//	if errors.Is(err, gcloudcx.APIError{}) {
//	  // do something with err
//	}
func (e APIError) Is(target error) bool {
	if actual, ok := target.(APIError); ok {
		if len(actual.Code) == 0 {
			return true // no ID means any error is a match
		}
		return e.Code == actual.Code
	}
	return false
}

// As attempts to convert the given error into the given target
//
// As returns true if the conversion was successful and the target is now populated.
//
// Example:
//
//	target := errors.ArgumentInvalid.Clone()
//	if errors.As(err, &target) {
//	  // do something with target
//	}
func (e APIError) As(target interface{}) bool {
	if actual, ok := target.(**APIError); ok {
		if *actual != nil && (*actual).Code != e.Code {
			return false
		}
		copy := e
		*actual = &copy
		return true
	}
	return false
}

// SetCorrelationID creates a new APIError from a given Error and records its correlation ID.
func (e APIError) SetCorrelationID(correlationID string) APIError {
	final := e
	final.CorrelationID = correlationID
	return final
}

// With creates a new Error from a given sentinel telling "what" is wrong and eventually their value.
//
// With also records the stack trace at the point it was called.
func (e APIError) With(what string, values ...interface{}) error {
	final := e
	// insert what in the values
	final.MessageParams = map[string]string{"what": what}
	for i, value := range values {
		final.MessageParams[fmt.Sprintf("value%d", i+1)] = fmt.Sprintf("%v", value)
	}
	values = append([]interface{}{what}, values...)
	final.MessageWithParams = fmt.Sprintf(final.Message, values...)
	final.Stack.Initialize()
	return final
}

// WithParams creates a new Error from a given sentinel telling "what" is wrong and eventually their value.
func (e APIError) WithParams(what string, values map[string]string) error {
	final := e
	final.MessageWithParams = fmt.Sprintf(final.Message, what)
	// TODO: Inject the values in the messageWithParams. Since it is a map, we need to find the keys in the messageWithParams and replace them with the values
	final.MessageParams = values
	final.Stack.Initialize()
	return final
}

// WithStack creates a new error from a given Error and records its stack.
func (e APIError) WithStack() error {
	final := e
	final.Stack.Initialize()
	return final
}

// UnmarshalJSON decodes a JSON payload into an APIError
func (e *APIError) UnmarshalJSON(payload []byte) (err error) {
	// Try to get an error from the login API (/oauth/token)
	oauthError := struct {
		Error       string `json:"error"`
		Description string `json:"description"`
	}{}
	err = json.Unmarshal(payload, &oauthError)
	if err == nil && len(oauthError.Error) > 0 && len(oauthError.Description) > 0 {
		switch oauthError.Error {
		case "invalid_client":
			*e = BadCredentialsError
			e.Message = fmt.Sprintf(e.Message, oauthError.Description)
			e.MessageParams = map[string]string{
				"reason":      oauthError.Error,
				"description": oauthError.Description,
			}
		default:
			*e = APIError{
				Status:  BadCredentialsError.Status,
				Code:    BadCredentialsError.Code,
				Message: fmt.Sprintf("%s: %s", oauthError.Description, oauthError.Error),
			}
		}
		return nil
	}

	// Get the standard structure
	type surrogate APIError
	inner := surrogate{}
	if err = json.Unmarshal(payload, &inner); err != nil {
		return err
	}
	*e = APIError(inner)
	return nil
}
