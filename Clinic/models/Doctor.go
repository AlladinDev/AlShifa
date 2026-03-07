package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Doctor struct {
	RegistrationDate time.Time          `json:"registrationDate,omitempty" bson:"registrationDate"`
	ID               primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name             string             `json:"name,omitempty" bson:"name"`
	Qualifications   string             `json:"qualifications,omitempty" bson:"qualifications"`
	Address          string             `json:"address,omitempty" bson:"address"`
	Experience       string             `json:"experience" bson:"experience"`
	Post             string             `json:"post" bson:"post"`
	WorkingAt        string             `json:"workingAt,omitempty" bson:"workingAt"`
	Role             string             `json:"role,omitempty" bson:"role"`
}
