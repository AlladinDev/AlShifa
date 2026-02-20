// Package utils provides utility functions and types for the AlShifa application.
package utils

import (
	structs "AlShifa/structs"
)

func ReturnAppError(
	err any,
	statusCode int,
	message string,
	reason string,
) *structs.IAppError {
	//fmt.Print("err received is", err)
	return &structs.IAppError{
		Message:    message,
		StatusCode: statusCode,
		Reason:     reason,
		ErrorObj:   err,
	}
}
