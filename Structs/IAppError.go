// Package structs contains shared data structures used across the application.
package structs

// swagger:model ErrorResponse

type IAppError struct {
	Message    string `json:"message"`
	StatusCode int    `json:"statusCode"`
	Reason     string `json:"reason"`
	ErrorObj   any    `json:"errorObj"`
}

func (appError *IAppError) Error() string {
	return appError.Message
}
