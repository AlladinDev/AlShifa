package validators

import (
	"AlShifa/clinic/models"
	middleware "AlShifa/middleware"
	"context"
	"fmt"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestValidateAppointmentDetails(t *testing.T) {
	// Helper: valid ObjectID as string
	validID := primitive.NewObjectID().Hex()
	//invalidID := "12345" // not a valid MongoID

	tests := []struct {
		name        string
		ctxUserID   interface{}
		appointment models.Appointment
		wantErrors  map[string]string
	}{
		{
			name:      "all valid",
			ctxUserID: validID,
			appointment: models.Appointment{
				AppointmentDate: time.Now().Add(24 * time.Hour),
				PatientName:     "John Doe",
				PatientAddress:  "123 Street",
				PatientMobile:   9876543210,
				Clinic:          primitive.NewObjectID(),
				Doctor:          primitive.NewObjectID(),
			},
			wantErrors: nil,
		},
		{
			name:      "missing user in context",
			ctxUserID: nil,
			appointment: models.Appointment{
				AppointmentDate: time.Now().Add(24 * time.Hour),
				PatientName:     "John Doe",
				PatientAddress:  "123 Street",
				PatientMobile:   9876543210,
				Clinic:          primitive.NewObjectID(),
				Doctor:          primitive.NewObjectID(),
			},
			wantErrors: map[string]string{"userId": "User ID is missing from context"},
		},
		{
			name:      "invalid user in context",
			ctxUserID: 123,
			appointment: models.Appointment{
				AppointmentDate: time.Now().Add(24 * time.Hour),
				PatientName:     "John Doe",
				PatientAddress:  "123 Street",
				PatientMobile:   9876543210,
				Clinic:          primitive.NewObjectID(),
				Doctor:          primitive.NewObjectID(),
			},
			wantErrors: map[string]string{"userId": "User ID is invalid"},
		},
		{
			name:      "appointment in past",
			ctxUserID: validID,
			appointment: models.Appointment{
				AppointmentDate: time.Now().Add(-24 * time.Hour),
				PatientName:     "John Doe",
				PatientAddress:  "123 Street",
				PatientMobile:   9876543210,
				Clinic:          primitive.NewObjectID(),
				Doctor:          primitive.NewObjectID(),
			},
			wantErrors: map[string]string{"appointmentDate": "Appointment Date cannot be in past"},
		},
		{
			name:      "empty patient name and address",
			ctxUserID: validID,
			appointment: models.Appointment{
				AppointmentDate: time.Now().Add(24 * time.Hour),
				PatientName:     "",
				PatientAddress:  "",
				PatientMobile:   9876543210,
				Clinic:          primitive.NewObjectID(),
				Doctor:          primitive.NewObjectID(),
			},
			wantErrors: map[string]string{
				"patientName":    "Patient name is required",
				"patientAddress": "Patient address is required",
			},
		},
		{
			name:      "invalid patient mobile",
			ctxUserID: validID,
			appointment: models.Appointment{
				AppointmentDate: time.Now().Add(24 * time.Hour),
				PatientName:     "John Doe",
				PatientAddress:  "123 Street",
				PatientMobile:   12345,
				Clinic:          primitive.NewObjectID(),
				Doctor:          primitive.NewObjectID(),
			},
			wantErrors: map[string]string{"patientMobile": "Patient mobile must be a valid 10-digit number"},
		},
		{
			name:      "invalid Clinic and doctor IDs",
			ctxUserID: validID,
			appointment: models.Appointment{
				AppointmentDate: time.Now().Add(24 * time.Hour),
				PatientName:     "John Doe",
				PatientAddress:  "123 Street",
				PatientMobile:   9876543210,
				Clinic:          primitive.NilObjectID, // simulate invalid
				Doctor:          primitive.NilObjectID,
			},
			wantErrors: map[string]string{
				"Clinic": "Clinic ID is not a valid MongoDB ObjectID",
				"doctor": "Doctor ID is not a valid MongoDB ObjectID",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			if tt.ctxUserID != nil {
				ctx = context.WithValue(ctx, middleware.ContextUserIDKey, tt.ctxUserID)
			}

			errors := ValidateAppointmentDetails(&tt.appointment, ctx)

			if len(errors) != len(tt.wantErrors) {
				fmt.Print(errors)

				t.Errorf("expected errors: %v, got: %v", tt.wantErrors, errors)
			}

			for key, wantMsg := range tt.wantErrors {
				gotMsg, exists := errors[key]
				if !exists {
					t.Errorf("expected error for key %s but not found", key)
				} else if gotMsg != wantMsg {
					t.Errorf("expected error %q for key %s, got %q", wantMsg, key, gotMsg)
				}
			}
		})
	}
}
