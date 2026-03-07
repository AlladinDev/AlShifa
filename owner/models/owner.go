package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Owner struct {
	CreatedAt time.Time          `json:"createdAt,omitzero" bson:"createdAt"` // 24 bytes
	Name      string             `json:"name,omitempty" bson:"name"`          // 16 bytes
	Address   string             `json:"address,omitempty" bson:"address"`    // 16 bytes
	Gender    string             `json:"gender,omitempty" bson:"gender"`      // 16 bytes
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Role      string             `json:"role,omitempty" bson:"role"` // 12 bytes (placed near end)
	Mobile    int64              `json:"mobile,omitempty" bson:"mobile"`
}
