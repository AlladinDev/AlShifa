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

	NameUserService   = "UserService"
	NameClinicService = "ClinicService"
	NameDoctorService = "DoctorService"

	RequestTimeout = 2 * time.Second
	APIVERSION     = "/v1"
	JwtExpiryTime  = time.Hour * 24 * 7

	//otp expiry time
	OTPExpiry               = 5 * time.Minute
	DoctorAddToclinicOTPKey = "DoctorAddToclinicOTPKey"
	CacheTTL                = OTPExpiry

	//Roles
	RoleUser        = "User"
	RoleAdmin       = "Admin"
	RoleDoctor      = "Doctor"
	RoleclinicOwner = "clinicOwner"

	//Plans
	PlanPaid = "Paid"
	PlanFree = "Free"

	//jwtprefix
	JwtPrefix = "BEARER "
)

var RolesAllowed = []string{RoleclinicOwner, RoleDoctor, RoleUser, RoleAdmin}
