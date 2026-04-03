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

	NameUserService        = "UserService"
	NameClinicService      = "ClinicService"
	NameDoctorService      = "DoctorService"
	NameAppointmentService = "AppointmentService"

	//clinic plans
	ClinicPlanBasic  = "basic"
	ClinicSilverPlan = "silver"

	RequestTimeout = 2 * time.Second
	APIVERSION     = "/v1"
	JwtExpiryTime  = time.Hour * 24 * 7

	//otp expiry time
	OTPExpiry               = 5 * time.Minute
	DoctorAddToclinicOTPKey = "DoctorAddToclinicOTPKey"
	CacheTTL                = OTPExpiry

	//Roles
	RoleUser               = "User"
	RoleAdmin              = "Admin"
	RoleDoctor             = "Doctor"
	RoleclinicOwner        = "ClinicOwner"
	RoleClinicReceptionist = "ClinicReceptionist"

	//Plans
	PlanPaid = "Paid"
	PlanFree = "Free"

	//jwtprefix
	JwtPrefix = "BEARER "

	//some length constants for validation
	MaxNameLength           = 20
	MaxAddressLength        = 50
	MaxQualificationsLength = 100
	MaxAge                  = 80
	MinAge                  = 15
	MaxPasswordLength       = 30
	MinPasswordLength       = 8
	MaxEmailLength          = 20
)

var RolesAllowed = []string{RoleclinicOwner, RoleDoctor, RoleUser, RoleAdmin}
