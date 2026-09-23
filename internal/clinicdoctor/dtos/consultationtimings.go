package dtos

import "time"

type ConsultationTimings struct {
	Day       string    `json:"day" bson:"day"`
	StartTime time.Time `json:"startTime" bson:"startTime"`
	EndTime   time.Time `json:"endTime" bson:"endTime"`
}
