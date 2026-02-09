// Package structs contains data structures used across the user module.
package structs

type LoginDetails struct {
	Mobile   int    `json:"mobile" bson:"mobile"`
	Password string `json:"password" bson:"password"`
}
