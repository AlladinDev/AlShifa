package utils

import "fmt"

func FormatBearerToken(token string) string {
	return fmt.Sprintf("BEARER %s", token)
}
