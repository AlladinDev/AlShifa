// Package doctormodels provides models for doctor module
package doctormodels

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TSpecializationCertificates struct {
	URL string `json:"url" bson:"url"`
	ID  string `json:"id" bson:"id" `
}

type Doctor struct {
	CreatedAt                   time.Time                     `json:"createdAt" bson:"createdAt" Form:"createdAt"`
	ID                          primitive.ObjectID            `json:"_id,omitempty" bson:"_id,omitempty"`
	Name                        string                        `json:"name,omitempty" bson:"name" Form:"name"`
	Mobile                      int64                         `json:"mobile" bson:"mobile" Form:"mobile"`
	Field                       string                        `json:"field" bson:"field" Form:"field"`
	Qualifications              string                        `json:"qualifications,omitempty" bson:"qualifications" Form:"qualifications"`
	Address                     string                        `json:"address,omitempty" bson:"address" Form:"address"`
	Age                         int8                          `json:"age" bson:"age" Form:"age"`
	Experience                  int8                          `json:"experience" bson:"experience" Form:"experience"`
	Post                        string                        `json:"post" bson:"post" Form:"post"`
	ProfilePhoto                string                        `json:"profilePhoto" bson:"profilePhoto" Form:"profilePhoto"`
	ProfilePhotoID              string                        `json:"profilePhotoID" bson:"profilePhotoID" Form:"profilePhotoID"`
	SpecializationCertificates  []TSpecializationCertificates `json:"specializationCertificates" bson:"specializationCertificates" Form:"specializationCertificates"`
	NmrNumber                   string                        `json:"nmrNumber" bson:"nmrNumber" Form:"nmrNumber"`
	ProfessionalDegreesVerified bool                          `json:"professionalDegreesVerified" bson:"professionalDegreesVerified" Form:"professionalDegreesVerified"`
	WorkingAt                   string                        `json:"workingAt,omitempty" bson:"workingAt" Form:"workingAt"`
	Role                        string                        `json:"role,omitempty" bson:"role" Form:"role"`
	Password                    string                        `json:"password" bson:"password" Form:"password"`
	AadhaarVerified             bool                          `json:"aadhaarVerified" bson:"aadhaarVerified" Form:"aadhaarVerified"`
	EmailVerified               bool                          `json:"emailVerified" bson:"emailVerified" Form:"emailVerified"`
	MobileVerified              bool                          `json:"mobileVerified" bson:"mobileVerified" Form:"mobileVerified"`
	AadhaarNumber               int64                         `json:"aadhaarNumber" bson:"aadhaarNumber" Form:"aadhaarNumber"`
	Email                       string                        `json:"email" bson:"email" Form:"email"`
	AccountOpeningAmountPaid    bool                          `json:"accountOpeningAmountPaid" bson:"accountOpeningAmountPaid"`
	ClinicsAllowedToOnboard     []primitive.ObjectID          `json:"clinicsAllowedToOnboard" bson:"clinicsAllowedToOnboard"`
}
