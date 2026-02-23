package utils

import (
	"fmt"

	"net/url"
	"regexp"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func sanitizeRegexPattern(pattern string) string {
	return regexp.QuoteMeta(pattern)
}
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
		for k, v := range queryValues {
			//dont add regex for date fields
			if strings.Contains(k, "Date") || strings.Contains(k, "date") {
				if t, ok := v.(time.Time); ok {
					dest[k] = t
					continue
				}
			}
			value := fmt.Sprintf("%v", v)
			dest[k] = primitive.Regex{Pattern: sanitizeRegexPattern(value), Options: "i"}
		}
	}
	//also return them
	return queryValues
}
