package gcloudcx

import "fmt"

// BotConnectorErrorInfo is used to provide error information in the Message Requests/Responses
type BotConnectorErrorInfo struct {
	ErrorCode    string `json:"errorCode"`    // The error code, e.g. "BotNotFound"
	ErrorMessage string `json:"errorMessage"` // A human-readable error message
}

// Error displays the error information in a human-readable format
//
// implements the error interface
func (e *BotConnectorErrorInfo) Error() string {
	return fmt.Sprintf("BotConnectorErrorInfo: %s - %s", e.ErrorCode, e.ErrorMessage)
}
