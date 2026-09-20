package utils

import response "github.com/AlladinDev/AlShifa/internal/shared/response"

func ReturnAppSuccess(statusCode int, message string, data any) *response.IAppSuccess {

	return &response.IAppSuccess{
		Message:    message,
		StatusCode: statusCode,
		Data:       data,
	}
}
