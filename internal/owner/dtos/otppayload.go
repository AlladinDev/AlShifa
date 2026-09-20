package dtos

import (
	"time"
)

type OTPPayload struct {
	ExpiresAt time.Time `json:"expiresAt"`
	Payload   string    `json:"payload"`
	OTPHash   string    `json:"otpHash"`
}
