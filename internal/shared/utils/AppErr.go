// Package utils provides utility functions and types for the github.com/AlladinDev/AlShifa application.
package utils

import (
	response "github.com/AlladinDev/AlShifa/internal/shared/response"
)

func ReturnAppError(
	err any,
	statusCode int,
	message string,
	reason string,
) *response.IAppError {
	//fmt.Print("err received is", err)
	return &response.IAppError{
		Message:    message,
		StatusCode: statusCode,
		Reason:     reason,
		ErrorObj:   err,
	}
}
