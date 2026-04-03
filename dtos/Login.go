// Package dtos provide resuable dtos for alshifa
package dtos

type LoginEmailPasswordDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
