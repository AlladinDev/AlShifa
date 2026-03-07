package workers

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func CleanInCompleteProfileCredientials(db *mongo.Database, expiryTime time.Duration, interval time.Duration) {
	log.Println("Clean worker for incomplete profile credientials started")
	ticker := time.NewTicker(interval)
	for range ticker.C {
		func() {
			ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
			defer cancel()
			expiryTime := time.Now().Add(-expiryTime)
			filter := bson.M{
				"profileCompleted": false,
				"createdAt": bson.M{
					"$lte": expiryTime,
				},
			}
			if _, err := db.Collection("Auth").DeleteMany(ctx, filter); err != nil {
				log.Println("Failed to clean incomplete profiles", err)
			}
		}()
	}
}
