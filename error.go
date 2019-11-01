package purecloud

import (
	"fmt"
	"encoding/json"
)

var (
	AuthenticationRequestTimeoutError = APIError{ Status: 504, Code: "authentication.request.timeout", Message: "Authentication request timeout." }
	BadRequestError = APIError{ Status: 400, Code: "bad.request", Message: "The request could not be understood by the server due to malformed syntax." }
	InternalServerError = APIError{ Status: 500, Code: "internal.server.error", Message: "The server encountered an unexpected condition which prevented it from fulfilling the request." }
	InvalidDateError = APIError{ Status: 400, Code: "invalid.date", Message: "Dates must be specified as ISO-8601 strings. For example: yyyy-MM-ddTHH:mm:ss.SSSZ" }
	InvalidValueError = APIError{ Status: 400, Code: "invalid.value", Message: "Value [%s] is not valid for field type [%s]. Allowable values are: %s" }
	MissingAnyPermissionsError = APIError{ Status: 403, Code: "missing.any.permissions", Message: "Unable to perform the requested action. You must have at least one of the following permissions assigned: %s" }
	MissingPermissionsError = APIError{ Status: 403, Code: "missing.permissions", Message: "Unable to perform the requested action. You are missing the following permission(s): %s" }
	NotAuthorizedError = APIError{ Status: 403, Code: "not.authorized", Message: "You are not authorized to perform the requested action." }
	NotFoundError = APIError{ Status: 404, Code: "not.found", Message: "The requested resource was not found." }
	RequestTimeoutError = APIError{ Status: 504, Code: "request.timeout", Message: "The request timed out." }
	ServiceUnavailableError = APIError{ Status: 503, Code: "service.unavailable", Message: "Service Unavailable - The server is currently unavailable (because it is overloaded or down for maintenance)." }
	TooManyRequestsError = APIError{ Status: 429, Code: "too.many.requests", Message: "Rate limit exceeded the maximum [%s] requests within [%s] seconds" }
	UnsupportedMediaTypeError = APIError{ Status: 415, Code: "unsupported.media.type", Message: "Unsupported Media Type - Unsupported or incorrect media type, such as an incorrect Content-Type value in the header." }

	AuthenticationRequiredError = APIError{ Status: 401, Code: "authentication.required" , Message: "No authentication bearer token specified in authorization header." }
	BadCredentialsError = APIError{ Status: 401, Code: "bad.credentials" , Message: "Invalid login credentials." }
	CredentialsExpiredError = APIError{ Status: 401, Code: "credentials.expired" , Message: "The supplied credentials are expired and cannot be used." }

	ChatConversationStateError = APIError{ Status: 400, Code: "chat.error.conversation.state", Message: "The conversation is in a state which does not permit this action." }
	ChatMemberStateError = APIError{ Status: 400, Code: "chat.error.member.state", Message: "The conversation member is in a state which does not permit this action." }
	ChatDeploymentBadAuthError = APIError{ Status: 400, Code: "chat.deployment.bad.auth", Message: "The customer member authentication has failed." }
	ChatDeploymentDisabledError = APIError{ Status: 400, Code: "chat.deployment.disabled", Message: "The web chat deployment is currently disabled." }
	ChatDeploymentRequireAuth = APIError{ Status: 400, Code: "chat.deployment.require.auth", Message: "The deployment requires the customer member to be authenticated." }
	ChatInvalidQueueError  = APIError{ Status: 400, Code: "chat.error.invalid.queue", Message: "The specified queue is not valid." }
	ChatCreateConversationRequestRoutingTargetError = APIError{ Status: 400, Code: "chat.error.createconversationrequest.routingtarget", Message: "The routing target is not valid." }
)

// APIError represents an error from the PureCloud API
type APIError struct {
	Status            int               `json:"status,omitempty"`
	Code              string            `json:"code,omitempty"`
	Message           string            `json:"message,omitempty"`
	MessageParams     map[string]string `json:"messageParams,omitempty"`
	MessageWithParams string            `json:"messageWithParams,omitempty"`
	EntityID          string            `json:"entityId,omitempty"`
	EntityName        string            `json:"entityName,omitempty"`
	ContextID         string            `json:"contextId,omitempty"`
	Details           []APIErrorDetails `json:"details,omitempty"`
	Errors            []APIError        `json:"errors,omitempty"`
}

// APIErrorDetails contains the details of an APIError
type APIErrorDetails struct {
	ErrorCode  string `json:"errorCode,omitempty"`
	FieldName  string `json:"fieldName,omitempty"`
	EntityID   string `json:"entityId,omitempty"`
	EntityName string `json:"entityName,omitempty"`
}

// Error returns a string representation of this error
func (e APIError) Error() string {
	if len(e.MessageWithParams) > 0 { return e.MessageWithParams }
	if len(e.Message) > 0           { return e.Message }
	return e.Code
}

// UnmarshalJSON decodes a JSON payload into an APIError
func (e *APIError) UnmarshalJSON(payload []byte) (err error) {
	// Try to get an error from the login API (/oauth/token)
	oauthError := struct{
		Error       string `json:"error"`
		Description string `json:"description"`
	}{}
	err = json.Unmarshal(payload, &oauthError)
	if err == nil && len(oauthError.Error) > 0 && len(oauthError.Description) > 0 {
		*e = APIError{
			Code:    BadCredentialsError.Code,
			Message: fmt.Sprintf("%s: %s", oauthError.Description, oauthError.Error),
		}
		return nil
	}

	// Get the standard structure
	type surrogate APIError
	inner := surrogate{}
	if err = json.Unmarshal(payload, &inner); err != nil { return err }
	*e = APIError(inner)
	return nil
}