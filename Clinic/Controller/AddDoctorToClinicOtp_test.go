package controller

import (
	"AlShifa/clinic/models"
	middleware "AlShifa/middleware"
	structs "AlShifa/structs"
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

func TestAddDoctorToclinicOTPHandler(t *testing.T) {
	testCases := []struct {
		Name                        string
		Data                        models.AddDoctorToclinic
		ErrorExpected               bool
		MockValidateclinicDetailsFn func(clinicDetails *models.AddDoctorToclinic) map[string]string
		MockService                 *MockService
		ExpectedErr                 *structs.IAppError
		ExpectedStatusCode          int
		OwnerID                     string
	}{
		{
			Name: "Successfull eveything is ok",
			Data: models.AddDoctorToclinic{
				DoctorID:    primitive.NewObjectID(),
				StartTime:   time.Now().Add(1 * time.Hour),
				EndTime:     time.Now().Add(7 * time.Hour),
				WorkingDays: []string{"monday"},
			},
			ErrorExpected:      false,
			ExpectedStatusCode: http.StatusOK,
			ExpectedErr:        nil,
			OwnerID:            "696b86d0f1d4392b0bd68155",
			MockValidateclinicDetailsFn: func(clinicDetails *models.AddDoctorToclinic) map[string]string {
				return nil
			},
			MockService: &MockService{
				SearchclinicFn: func(ctx context.Context, filter bson.M) ([]models.ClinicDoctor, *structs.IAppError) {
					return []models.clinicDoctor{
						{ID: primitive.NewObjectID()},
					}, nil
				},
				AddDoctorToclinicFn: func(ctx context.Context, clinicDetails models.AddDoctorToclinic) *structs.IAppError {
					return nil
				},
			},
		},

		{

			Name: "Failed as ownerid is empty else eveything is ok",
			Data: models.AddDoctorToclinic{
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
			MockValidateclinicDetailsFn: func(clinicDetails *models.AddDoctorToclinic) map[string]string {
				return nil
			},
			MockService: &MockService{
				SearchclinicFn: func(ctx context.Context, filter bson.M) ([]models.clinicDoctor, *structs.IAppError) {
					return []models.clinicDoctor{
						{ID: primitive.NewObjectID()},
					}, nil
				},
				AddDoctorToclinicFn: func(ctx context.Context, clinicDetails models.AddDoctorToclinic) *structs.IAppError {
					return nil
				},
			},
		},
		{

			Name: "Failed as search clinic returned []clinic ",
			Data: models.AddDoctorToclinic{
				DoctorID:    primitive.NewObjectID(),
				StartTime:   time.Now().Add(1 * time.Hour),
				EndTime:     time.Now().Add(7 * time.Hour),
				WorkingDays: []string{"monday"},
			},
			ErrorExpected:      true,
			ExpectedStatusCode: http.StatusNotFound,
			ExpectedErr: &structs.IAppError{
				Message:    "No clinic Found",
				StatusCode: http.StatusNotFound,
				Reason:     errors.New("no clinic found").Error(),
			},
			OwnerID: "696b86d0f1d4392b0bd68155",
			MockValidateclinicDetailsFn: func(clinicDetails *models.AddDoctorToclinic) map[string]string {
				return nil
			},
			MockService: &MockService{
				SearchclinicFn: func(ctx context.Context, filter bson.M) ([]models.clinicDoctor, *structs.IAppError) {
					return nil, &structs.IAppError{
						Message:    "No clinic Found",
						StatusCode: http.StatusNotFound,
						Reason:     errors.New("no clinic found").Error(),
					}
				},
				AddDoctorToclinicFn: func(ctx context.Context, clinicDetails models.AddDoctorToclinic) *structs.IAppError {
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
		controller := NewController(tc.MockService, tc.MockValidateclinicDetailsFn)
		res := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/clinic", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		//inject ownerId into context value
		ctx := context.WithValue(req.Context(), middleware.ContextUserIDKey, tc.OwnerID)
		req = req.WithContext(ctx)

		controller.AddDoctorToclinic(res, req)

		if res.Code != tc.ExpectedStatusCode {
			t.Fatalf("expected status %d, got %d. body=%s",
				http.StatusOK,
				res.Code,
				res.Body.String(),
			)
		}

	}
}
