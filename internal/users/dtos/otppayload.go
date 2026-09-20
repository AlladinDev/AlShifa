// Package dtos provides dtos for user module
package dtos

import "time"

type OTPPayload struct {
	OTPHash      string    `json:"otpHash"`
	ExpiresAt    time.Time `json:"expiresAt"`
	MobileNumber string    `json:"mobileNumber"`
}
