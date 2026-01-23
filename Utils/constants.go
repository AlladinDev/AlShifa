package utils

import "time"

const (
	RequestTimeout = 2 * time.Second
	APIVERSION     = "/v1"
	JwtExpiryTime  = time.Hour * 24 * 7

	//Roles
	RoleUser        = "User"
	RoleAdmin       = "Admin"
	RoleDoctor      = "Doctor"
	RoleClinicOwner = "ClinicOwner"

	// Name
	MinNameLength = 2
	MaxNameLength = 50

	//jwtprefix
	JwtPrefix = "BEARER "

	// Password
	MinPasswordLength        = 8
	MaxPasswordLength        = 30
	PasswordRegex     string = `^(?=.*[a-z])(?=.*[A-Z])(?=.*[!@#$%^&*(),.?":{}|<>]).{8,}$`

	// Age
	MinAge = 1
	MaxAge = 100

	// Address
	MinAddressLength = 5
	MaxAddressLength = 100

	// Mobile & Pincode
	MobileLength   = 10
	PincodeLength  = 6
	MaxEmailLength = 50

	// Regex
	MinEmailLength = 5
	EmailRegex     = `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`

	// Name validation
	NameMissingErrMsg = "Name is required"
	ShortNameErrMsg   = "Name must be at least 2 characters long"
	LongNameErrMsg    = "Name cannot exceed 50 characters"
	InvalidNameErrMsg = "Name contains invalid characters"

	// Email validation
	EmailMissingErrMsg    = "Email is required"
	InvalidEmailFormatMsg = "Email format is invalid"
	LongEmailErrMsg       = "Email cannot exceed 100 characters"
	ShortEmailErrMsg      = "Email is too short"

	// Password validation
	PasswordMissingErrMsg = "Password is required"
	PasswordWeakErrMsg    = "Password is too weak, must include uppercase, lowercase, number, and special character"
	ShortPasswordErrMsg   = "Password must be at least 8 characters long"
	LongPasswordErrMsg    = "Password cannot exceed 128 characters"

	// Age validation
	AgeMissingErrMsg  = "Age is required"
	InvalidAgeErrMsg  = "Age must be a valid number"
	AgeTooYoungErrMsg = "Age must be at least 18"
	AgeTooOldErrMsg   = "Age cannot exceed 120"

	// Gender validation
	GenderMissingErrMsg = "Gender is required"
	InvalidGenderErrMsg = "Gender must be 'Male', 'Female', or 'Other'"

	// Mobile / Phone validation
	MobileMissingErrMsg    = "Mobile number is required"
	InvalidMobileNumberMsg = "10 Digit Mobile Number Required"
	ShortMobileErrMsg      = "Mobile number is too short"
	LongMobileErrMsg       = "Mobile number is too long"

	// Address validation
	AddressMissingErrMsg = "Address is required"
	ShortAddressErrMsg   = "Address is too short, provide more details"
	LongAddressErrMsg    = "Address is too long, please shorten it"

	// Username / User ID validation
	UsernameMissingErrMsg = "Username is required"
	ShortUsernameErrMsg   = "Username must be at least 3 characters long"
	LongUsernameErrMsg    = "Username cannot exceed 30 characters"
	InvalidUsernameErrMsg = "Username contains invalid characters"

	// Optional additional fields
	CityMissingErrMsg       = "City is required"
	StateMissingErrMsg      = "State is required"
	PostalCodeMissingErrMsg = "Postal code is required"
	InvalidPostalCodeErrMsg = "Postal code format is invalid"
	CountryMissingErrMsg    = "Country is required"
)
