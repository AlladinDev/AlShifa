package utils

import (
	"fmt"
	"maps"
	"net/url"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TransformParamIDS(params url.Values, dest map[string]any) map[string]any {
	queryValues := make(map[string]any)

	for key := range params {

		value := params.Get(key)
		fmt.Println("parsed key,valye is", key, value)
		//if param can be converted into mongodb format convert it if not just continue
		mongoDBID, err := primitive.ObjectIDFromHex(value)
		if err == nil {
			queryValues[key] = mongoDBID
			continue
		}

		// 2️⃣ Try parsing RFC3339 date (ISO format)
		fmt.Println("parsed key,valye is", key, value)
		if parsedTime, err := time.Parse(time.RFC3339, value); err == nil {
			fmt.Println("parsed key,valye is", key, value)
			queryValues[key] = parsedTime
			continue
		}

		queryValues[key] = value
	}
	//if user passes dest then copy them to that dest for convenince
	if dest != nil {
		maps.Copy(dest, queryValues)
	}
	//also return them
	return queryValues
}
