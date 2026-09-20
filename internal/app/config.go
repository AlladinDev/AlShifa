// Package app provides internal application dependencies and configuration.
package app

import (
	"net/http"

	emailNotifier "github.com/AlladinDev/AlShifa/internal/shared/emailnotifier"
	"github.com/AlladinDev/AlShifa/internal/shared/interfaces"
	"github.com/go-chi/chi/v5"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
)

type ICache interfaces.Cache[string, []byte]
type App struct {
	*chi.Mux
	DB            *mongo.Database
	Redis         *redis.Client
	Broker        interfaces.IMessageBroker
	httpClint     *http.Client
	ImageUploader interfaces.IImageUploader
	EmailNotifier interfaces.INotifier
	KVStore       ICache
}

func NewApp() *App {
	return &App{}
}

func (app *App) WithRedis(r *redis.Client) *App {
	app.Redis = r
	return app
}

func (app *App) WithEmailNotifier() *App {
	emailService := emailNotifier.NewEmailNotifier()
	app.EmailNotifier = emailService
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

func (app *App) WithImageUploader(imageUploader interfaces.IImageUploader) *App {
	app.ImageUploader = imageUploader
	return app
}

func (app *App) WithHTTPClient() *App {
	app.httpClint = &http.Client{}
	return app
}

func (app *App) WithBroker(broker interfaces.IMessageBroker) *App {
	app.Broker = broker
	return app
}

func (app *App) WithKVStore(cache ICache) *App {
	app.KVStore = cache
	return app
}
