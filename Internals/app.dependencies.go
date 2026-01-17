// Package internals contains config rrelated things
package internals

import (
	"net/http"

	"go.mongodb.org/mongo-driver/mongo"
)

type App struct {
	Server *http.ServeMux
	DB     *mongo.Database
}
