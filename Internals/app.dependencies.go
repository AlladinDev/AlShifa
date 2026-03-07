// Package internals provides internal application dependencies and configuration.
package internals

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
)

type App struct {
	*chi.Mux
	DB        *mongo.Database
	Redis     *redis.Client
	DI        *DI
	httpClint *http.Client
}

func NewApp() *App {
	return &App{}
}

func (app *App) WithRedis(r *redis.Client) *App {
	app.Redis = r
	return app
}

func (app *App) WithServer(mux *chi.Mux) *App {
	app.Mux = mux
	return app
}

func (app *App) WithDB(db *mongo.Database) *App {
	app.DB = db
	return app
}

func (app *App) WithDI() *App {
	app.DI = NewDI()
	return app
}

func (app *App) WithHTTPClient() *App {
	app.httpClint = &http.Client{}
	return app
}
