package utils

import (
	structs "AlShifa/Structs"
	"encoding/json"
	"net/http"
)

func WriteResponse(res http.ResponseWriter, response any) error {
	return json.NewEncoder(res).Encode(response)
}
func InvalidMethodResponse(methodName string, res http.ResponseWriter) error {

	return WriteResponse(res, structs.IAppError{
		Message:    "Only " + methodName + " Requests Are Allowed",
		StatusCode: 400,
		ErrorObj:   nil,
	})
}
