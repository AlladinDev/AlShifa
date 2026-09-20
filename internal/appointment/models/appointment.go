package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IAppointmentFees struct {
	Amount int  `json:"amount" bson:"amount"`
	Paid   bool `json:"paid" bson:"paid"`

	//createdAt represents time at which appointment was booked
	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	//paymentdoneAt representa time at which user went to clinic and paid the actual doctor or clinic fees like rs 500
	PaymentDoneAt time.Time `json:"paymentDoneAt" bson:"paymentDoneAt"`
}

type Appointment struct {
	// 24 bytes each
	AppointmentDate time.Time `json:"appointmentDate" bson:"appointmentDate"`
	CreatedAt       time.Time `json:"registrationDate" bson:"registrationDate"`

	UserID primitive.ObjectID `json:"userID" bson:"userID"`

	// 16 bytes each
	PatientName    string `json:"patientName" bson:"patientName"`
	PatientAddress string `json:"patientAddress" bson:"patientAddress"`

	//it means confirmed,cancelled,completed or noshow
	AppointmentStatus string `json:"status" bson:"status"`

	//booking fee status like rs 5 for booking it doesnt represent  actual appointment charge ,charged by clinic
	AppointmentBookingFeeStatus string `json:"appointmentBookingFeeStatus" bson:"appointmentBookingFeeStatus"`

	//appointment fees details represents the actual amount user pays at clinic for completing appointment like rs 500 mostly clinics charge this
	AppointmentFeesDetails IAppointmentFees `json:"appointmentFeesDetails" bson:"appointmentFeesDetails"`
	DoctorName             string           `json:"doctorName" bson:"doctorName"`
	ClinicName             string           `json:"clinicName" bson:"clinicName"`
	ClinicAddress          string           `json:"clinicAddress" bson:"clinicAddress"`

	// 12 bytes each
	ID       primitive.ObjectID `json:"id" bson:"_id"`
	ClinicID primitive.ObjectID `json:"clinicID" bson:"clinicID"`
	DoctorID primitive.ObjectID `json:"doctorID" bson:"doctorID"`

	Department string `json:"department" bson:"department"`

	// 8 bytes
	Mobile string `json:"mobile" bson:"mobile"`

	ClinicMaxAppointments int `json:"clinicMaxAppointments" bson:"clinicMaxAppointments"`

	// 1 byte (placed last to avoid padding waste)
	Slot      int       `json:"slot" bson:"slot"`
	StartTime time.Time `json:"startTime" bson:"startTime"`
	EndTime   time.Time `json:"endTime" bson:"endTime"`

	AppointmentBookingOrderID       string `json:"appointmentBookingOrderID" bson:"appointmentBookingOrderID"`
	AppointmentBookingTransactionID string `json:"appointmentBookingTransactionID" bson:"appointmentBookingTransactionID"`

	//this is the id of the transaction when patient pays the clinic fees like rs 500 also through our platform
	AppointmentTransactionID string `json:"appointmentTransactionID" bson:"appointmentTransactionID"`
}
