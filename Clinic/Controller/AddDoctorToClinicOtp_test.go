package controller

import (
	"AlShifa/Clinic/models"
	middleware "AlShifa/Middleware"
	structs "AlShifa/Structs"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestAddDoctorToClinicOTPHandler(t *testing.T) {
	testCases := []struct {
		Name                        string
		Data                        models.AddDoctorToClinic
		ErrorExpected               bool
		MockValidateClinicDetailsFn func(clinicDetails *models.AddDoctorToClinic) map[string]string
		MockService                 *MockService
		ExpectedErr                 *structs.IAppError
		ExpectedStatusCode          int
		OwnerID                     string
	}{
		{
			Name: "Successfull eveything is ok",
			Data: models.AddDoctorToClinic{
				DoctorID:    primitive.NewObjectID(),
				StartTime:   time.Now().Add(1 * time.Hour),
				EndTime:     time.Now().Add(7 * time.Hour),
				WorkingDays: []string{"monday"},
			},
			ErrorExpected:      false,
			ExpectedStatusCode: http.StatusOK,
			ExpectedErr:        nil,
			OwnerID:            "696b86d0f1d4392b0bd68155",
			MockValidateClinicDetailsFn: func(clinicDetails *models.AddDoctorToClinic) map[string]string {
				return nil
			},
			MockService: &MockService{
				SearchClinicFn: func(ctx context.Context, filter bson.M) ([]models.Clinic, *structs.IAppError) {
					return []models.Clinic{
						{ID: primitive.NewObjectID()},
					}, nil
				},
				AddDoctorToClinicFn: func(ctx context.Context, clinicDetails models.AddDoctorToClinic) *structs.IAppError {
					return nil
				},
			},
		},

		{

			Name: "Failed as ownerid is empty else eveything is ok",
			Data: models.AddDoctorToClinic{
				DoctorID:    primitive.NewObjectID(),
				StartTime:   time.Now().Add(1 * time.Hour),
				EndTime:     time.Now().Add(7 * time.Hour),
				WorkingDays: []string{"monday"},
			},
			ErrorExpected:      true,
			ExpectedStatusCode: http.StatusBadRequest,
			ExpectedErr: &structs.IAppError{
				Message:    "OwnerID missing",
				Reason:     "OwnerId is missing for authentication",
				StatusCode: 400,
				ErrorObj:   errors.New("ownerid is missing"),
			},
			OwnerID: "",
			MockValidateClinicDetailsFn: func(clinicDetails *models.AddDoctorToClinic) map[string]string {
				return nil
			},
			MockService: &MockService{
				SearchClinicFn: func(ctx context.Context, filter bson.M) ([]models.Clinic, *structs.IAppError) {
					return []models.Clinic{
						{ID: primitive.NewObjectID()},
					}, nil
				},
				AddDoctorToClinicFn: func(ctx context.Context, clinicDetails models.AddDoctorToClinic) *structs.IAppError {
					return nil
				},
			},
		},
		{

			Name: "Failed as search clinic returned []clinic ",
			Data: models.AddDoctorToClinic{
				DoctorID:    primitive.NewObjectID(),
				StartTime:   time.Now().Add(1 * time.Hour),
				EndTime:     time.Now().Add(7 * time.Hour),
				WorkingDays: []string{"monday"},
			},
			ErrorExpected:      true,
			ExpectedStatusCode: http.StatusNotFound,
			ExpectedErr: &structs.IAppError{
				Message:    "No Clinic Found",
				StatusCode: http.StatusNotFound,
				Reason:     errors.New("no clinic found").Error(),
			},
			OwnerID: "696b86d0f1d4392b0bd68155",
			MockValidateClinicDetailsFn: func(clinicDetails *models.AddDoctorToClinic) map[string]string {
				return nil
			},
			MockService: &MockService{
				SearchClinicFn: func(ctx context.Context, filter bson.M) ([]models.Clinic, *structs.IAppError) {
					return nil, &structs.IAppError{
						Message:    "No Clinic Found",
						StatusCode: http.StatusNotFound,
						Reason:     errors.New("no clinic found").Error(),
					}
				},
				AddDoctorToClinicFn: func(ctx context.Context, clinicDetails models.AddDoctorToClinic) *structs.IAppError {
					return nil
				},
			},
		},
	}

	for _, tc := range testCases {
		jsonData, err := json.Marshal(tc.Data)
		if err != nil {
			t.Fatal("Failed to convert data into json")
		}
		controller := NewController(tc.MockService, tc.MockValidateClinicDetailsFn)
		res := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/clinic", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		//inject ownerId into context value
		ctx := context.WithValue(req.Context(), middleware.ContextUserIDKey, tc.OwnerID)
		req = req.WithContext(ctx)

		controller.AddDoctorToClinic(res, req)

		if res.Code != tc.ExpectedStatusCode {
			t.Fatalf("expected status %d, got %d. body=%s",
				http.StatusOK,
				res.Code,
				res.Body.String(),
			)
		}

	}
}
