package utils

import (
	"maps"
	"net/url"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TransformParamIDS(params url.Values, dest map[string]any) map[string]any {
	queryValues := make(map[string]any)

	for key := range params {
		value := params.Get(key)
		//if param can be converted into mongodb format convert it if not just continue
		mongoDBID, err := primitive.ObjectIDFromHex(value)
		if err != nil {
			queryValues[key] = value
			continue
		}
		queryValues[key] = mongoDBID
	}
	//if user passes dest then copy them to that dest for convenince
	if dest != nil {
		maps.Copy(dest, queryValues)
	}
	//also return them
	return queryValues
}
