package customerrors

import (
	"errors"
)

var ErrClinicWalletNotFound = errors.New("clinic doesnt have any wallet yet")
var ErrClinicPlanNotFound = errors.New("clinic doesnt have any plan yet")
var ErrClinicInsufficientWalletBalance = errors.New("clinic wallet has insufficient balance")
var ErrClinicNotFound = errors.New("clinic not found")
var ErrUnknownErr = errors.New("unknown error occurred")
var ErrOnboardingNotAllowed = errors.New("onboarding not allowed")
var ErrNoDocumentModified = errors.New("no document modified")
var ErrAppointmentDateNotPossible = errors.New("appointment date not possible plz choose some other date")
var ErrDoctorClinicMappingNotFound = errors.New("doctor is not onboarded to the specified clinic")
