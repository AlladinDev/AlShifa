// Package constants provides constants
package constants

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

// JwtCustomClaims struct with userID, role, and expiration
type JwtCustomClaims struct {
	UserID     string `json:"userID"`
	Role       string `json:"role"`
	Email      string `json:"email"`
	Mobile     string `json:"mobile"`
	IsVerified bool   `json:"isVerified"`
	jwt.RegisteredClaims
}

const (
	KeyUserID                contextKey = "UserID"
	KeyUserRole              contextKey = "Role"
	KeyEmail                 contextKey = "Email"
	KeyMobile                contextKey = "Mobile"
	StatusAppointmentPending            = "pending"

	NameOTPToken                  = "OTP_TOKEN"
	NameAuthToken                 = "AUTH_TOKEN"
	NameAppStore                  = "AppStore"
	NameDISystem                  = "DISystem"
	NameUserModule                = "UserModule"
	NamePlanModule                = "PlanModule"
	NamePaymentsModule            = "PaymentModule"
	NameClinicModule              = "ClinicModule"
	NameAppointmentModule         = "AppointmentModule"
	NameDoctorModule              = "DoctorModule"
	NameDoctorClinicMappingModule = "DoctorClinicMappingModule"
	NameOwnerModule               = "OwnerModule"
	NameAdminModule               = "AdminModule"

	AppointmentBookingFeesPending      = "APPOINTMENT_BOOKING_FEES-PENDING"
	AppointmentBookingFeesPaid         = "Appointment_Booking_Fees_Paid"
	AppointmentBookingOrderCreated     = "AppointmentBookingOrderCreated"
	NameAppointmentBookingPaymentToken = "AppoinmentBookingAppointmentToken"
	AppointmentInitiated               = "AppointmentInitiated"
	IndianCurrencyCode                 = "INR"

	AppointmentPending   = "Appointment_Pending"
	AppointmentCompleted = "Appointment_Completed"
	AppointmentCancelled = "Appointment_Cancelled"
	AppointmentRejected  = "Appointment_Rejected"
	AppointmentConfirmed = "Appointment_Confirmed"

	//appointmentfeespending or appointmentfeespaid represents clinic or doctor fees like rs 500 paid when patient visits clinic and pays them it goes to clinic
	AppointmentFeesPending = "APPOINTMENT_FEES-PENDING"
	AppointmentFeesPaid    = "Appointment_Fees_Paid"

	//it means appointment was booked or confirmed but patient didnt came so appointment became stale
	NoShowAppointment = "NoShowAppointment"

	EmailSendingURL    = "BREVO_EMAIL_SENDING_URL"
	EmailSendingAPIKey = "BREVO_EMAIL_API_KEY"
	EncryptionKey      = "ENCRYPTION_KEY"

	NameUserService        = "UserService"
	NameClinicService      = "ClinicService"
	NameDoctorService      = "DoctorService"
	NameAppointmentService = "AppointmentService"
	NamePlanService        = "PlanService"

	//clinic plans
	ClinicPlanBasic   = "basic"
	ClinicSilverPlan  = "silver"
	ClinicPlanPremium = "premium"
	ClinicPlanGold    = "gold"

	//appointment plan
	PlanAppointment = "AppointmentPlan"

	FreeWalletBalance = 100

	RequestTimeout = 10 * time.Second
	APIVERSION     = "/v1"
	JwtExpiryTime  = time.Hour * 24 * 7

	//otp expiry time
	OTPExpiry                           = 5 * time.Minute
	DoctorAddToclinicOTPKey             = "DoctorAddToclinicOTPKey"
	CacheTTL                            = OTPExpiry
	AppointmentBookingPaymentExpiryTime = 5 * time.Minute

	//Roles
	RoleUser               = "User"
	RoleAdmin              = "Admin"
	RoleDoctor             = "Doctor"
	RoleclinicOwner        = "ClinicOwner"
	RoleClinicReceptionist = "ClinicReceptionist"

	//errors

	//Plans
	PlanPaid = "Paid"
	PlanFree = "Free"

	//jwtprefix
	JwtPrefix = "BEARER "

	//some length constants for validation
	MaxNameLength                = 20
	MaxAddressLength             = 50
	AadhaarNumberDigits          = 12
	MobileNumberDigits           = 10
	MinAddressLength             = 5
	MaxQualificationsLength      = 100
	MinQualificationsLength      = 4
	MaxAge                       = 80
	MinAge                       = 15
	MaxPasswordLength            = 30
	MinPasswordLength            = 8
	MaxEmailLength               = 40
	MinNameLength                = 3
	MinSpecializationFieldLength = 4
	MaxSpecializationFieldLength = 40
)

var RolesAllowed = []string{RoleclinicOwner, RoleDoctor, RoleUser, RoleAdmin}
