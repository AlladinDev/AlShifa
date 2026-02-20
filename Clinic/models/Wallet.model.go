package models

type WalletDetails struct {
	AvailableBalance int64 `json:"availableBalance" bson:"availableBalance"`
}
