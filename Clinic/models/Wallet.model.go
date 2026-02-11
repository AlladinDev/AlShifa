package models

import (
	"time"
)

type WalletDetails struct {
	UpdatedAt        time.Time `json:"updatedAt" bson:"updatedAt"`
	AvailableBalance int64     `json:"availableBalance" bson:"availableBalance"`
}
