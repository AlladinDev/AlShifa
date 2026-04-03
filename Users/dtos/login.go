// Package dtos provides dtos for user module
package dtos

type LoginDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
