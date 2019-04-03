package purecloud

import (
	"fmt"
	"encoding/json"
)

const CredentialsExpired = "credentials.expired"
const AuthenticationRequired = "authentication.required"
const BadCredentials = "bad.credentials"

// APIError represents an error from the PureCloud API
type APIError struct {
	Status            int               `json:"status,omitempty"`
	Code              string            `json:"code,omitempty"`
	Message           string            `json:"message,omitempty"`
	MessageParams     map[string]string `json:messageParams,omitempty"`
	MessageWithParams string            `json:messageWithParams,omitempty"`
	EntityID          string            `json:entityId,omitempty"`
	EntityName        string            `json:entityName,omitempty"`
	ContextID         string            `json:"contextId,omitempty"`
	Details           []APIErrorDetails `json:"details,omitempty"`
	Errors            []APIError        `json:"errors,omitempty"`
}

// APIErrorDetails contains the details of an APIError
type APIErrorDetails struct {
	ErrorCode  string `json:"errorCode,omitempty"`
	FieldName  string `json:"fieldName,omitempty"`
	EntityID   string `json:entityId,omitempty"`
	EntityName string `json:entityName,omitempty"`
}

// String returns a string representation of this
func (e APIError) String() string {
	return e.Message
}

// Error returns a string representation of this error
func (e APIError) Error() string {
	return e.Message
}

// UnmarshalJSON decodes a JSON payload into an APIError
func (e *APIError) UnmarshalJSON(payload []byte) (err error) {
	// Try to get an error from the login API (/oauth/token)
	oauthError := struct{
		Error       string `json:"error,omitempty"`
		Description string `json:"description,omitempty"`
	}{}
	if err = json.Unmarshal(payload, &oauthError); err == nil {
		e.Message = fmt.Sprintf("%s: %s", oauthError.Description, oauthError.Error)
		e.Code    = BadCredentials
	}

	// Get the standard structure
	type surrogate APIError
	inner := surrogate{}
	if err = json.Unmarshal(payload, &inner); err != nil { return err }
	*e = APIError(inner)
	return nil
}