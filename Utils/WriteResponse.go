package utils

import (
	"encoding/json"
	"net/http"

	structs "github.com/AlladinDev/AlShifa/structs"
)

func WriteResponse(res http.ResponseWriter, status int, response any) error {
	res.WriteHeader(int(status))
	return json.NewEncoder(res).Encode(response)
}
func InvalidMethodResponse(methodName string, res http.ResponseWriter) error {

	return WriteResponse(res, http.StatusMethodNotAllowed, structs.IAppError{
		Message:    "Only " + methodName + " Requests Are Allowed",
		StatusCode: http.StatusMethodNotAllowed,
		ErrorObj:   nil,
	})
}
