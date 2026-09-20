// package doctordtos provides dtos for doctor module
package doctordtos

import "time"

type TOTPPayload struct {
	OTPHash   string    `json:"otpHash" bson:"otpHash"`
	Payload   string    `json:"payload" bson:"payload"`
	ExpiresAt time.Time `json:"expiresAt" bson:"expiresAt"`
}
