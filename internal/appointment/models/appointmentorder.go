package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AppointmentBookingPaymentOrder struct {
	ID              primitive.ObjectID `json:"_id" bson:"_id"`
	Amount          int                `json:"amount" bson:"amount"`
	Currency        string             `json:"currency" bson:"currency"`
	CreatedAt       time.Time          `json:"createdAt" bson:"createdAt"`
	AppointmentID   primitive.ObjectID `json:"appointmentID" bson:"appointmentID"`
	ProviderOrderID string             `json:"providerOrderID" bson:"providerOrderID"`
	PaymentStatus   string             `json:"paymentStatus" bson:"paymentStatus"`
	Provider        string             `json:"provider" bson:"provider"`
}
