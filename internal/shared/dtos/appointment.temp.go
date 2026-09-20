package dtos

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

//AppointmentPaymentToken is the token which is to be sent after temporary booking appointment this helps payment module to decide how much amount to deduct and also for emiting event to appointment module
type AppointmentPaymentToken struct {
	AppointmentID  primitive.ObjectID `json:"appointmentID" bson:"appointmentID"`
	AmountToDeduct int                `json:"amountToDeduct" bson:"amountToDeduct"`
	SlotDocumentID primitive.ObjectID `json:"slotDocumentID" bson:"slotDocumentID"`
	UserID         primitive.ObjectID `json:"userID" bson:"userID"`
	ExpiresAt      time.Time          `json:"expiresAt" bson:"expiresAt"`
}
