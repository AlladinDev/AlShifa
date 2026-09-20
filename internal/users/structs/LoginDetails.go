// Package structs contains data structures used across the user module.
package structs

type LoginDetails struct {
	Email    string `json:"email" bson:"email"`
	Password string `json:"password" bson:"password"`
}
