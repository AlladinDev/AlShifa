// Package utils provides utility functions and types for the AlShifa application.
package utils

import structs "AlShifa/Structs"

func ReturnAppError(
	err any,
	statusCode int,
	message string,
	reason string,
) *structs.IAppError {
	
	return &structs.IAppError{
		Message:    message,
		StatusCode: statusCode,
		Reason:     reason,
		ErrorObj:   err,
	}
}
