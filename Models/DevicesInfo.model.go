// Package models contains the data models for the application.
package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DevicesInfo struct {
	ID        primitive.ObjectID // 12 bytes
	UserID    primitive.ObjectID // 12 bytes
	CreatedAt time.Time          // 24 bytes
	DeviceID  string             // 16 bytes (header only)
}
