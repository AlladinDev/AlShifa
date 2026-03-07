// Package dtos provides dtos for auth module
package dtos

type LoginDetails struct {
	Email    string `json:"email" bson:"email"`
	Password string `json:"password" bson:"password"`
}
