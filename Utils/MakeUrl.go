package utils

import "fmt"

// MakeURL creates a pattern compatible with native Go routing.
// Usage: MakeURL("GET", "/user/{id}") -> "GET /api/v1/user/{id}"
// Usage: MakeURL("", "/public")      -> "/api/v1/public"
func MakeURL(method string, path string) string {
	fullPath := APIVERSION + path
	if method == "" {
		return fullPath
	}
	// Go 1.22+ requires an uppercase method followed by a space
	return fmt.Sprintf("%s %s", method, fullPath)
}
