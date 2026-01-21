package utils

import "time"

const RequestTimeout = 2 * time.Second
const APIVERSION = "/v1"
const JwtExpiryTime = time.Hour * 24 * 7
const RoleUser = "User"
const RoleAdmin = "Admin"
const RoleDoctor = "Doctor"
const RoleClinicOwner = "ClinicOwner"

const (
	// Name
	MinNameLength = 2
	MaxNameLength = 50

	//jwtprefix
	JwtPrefix = "BEARER "

	// Password
	MinPasswordLength = 8
	MaxPasswordLength = 30

	// Age
	MinAge = 1
	MaxAge = 50

	// Address
	MinAddressLength = 5

	// Mobile & Pincode
	MobileLength   = 10
	PincodeLength  = 6
	MaxEmailLength = 50

	// Regex
	EmailRegex = `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`
)
