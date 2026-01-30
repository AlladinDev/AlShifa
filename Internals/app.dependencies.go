// Package internals contains config rrelated things
package internals

import (
	"net/http"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
)

type App struct {
	Server *http.ServeMux
	DB     *mongo.Database
	Redis  *redis.Client
}
