package utils

import structs "AlShifa/Structs"

func ReturnAppSuccess(statusCode int, message string, data any) *structs.IAppSuccess {

	return &structs.IAppSuccess{
		Message:    message,
		StatusCode: statusCode,
		Data:       data,
	}
}
