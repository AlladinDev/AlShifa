package controller

import (
	middleware "AlShifa/Middleware"
	structs "AlShifa/Structs"
	models "AlShifa/Users/Models"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestRegisterUser(t *testing.T) {
	user := `{
		"name":     "John",
		"email":    "Saqlain@gmail.com",
		"address":  "soura",
		"password": "Saqlain@123",
	    "mobile":9797798243,
		"role":"user",
		"age":21,
		"pincode":190011
	}`
	testCases := []struct {
		Name               string
		UseCase            string
		mockService        *MockUserService
		ExpectedStatusCode int
		returnUserDate     string
	}{
		//test case when everything is ok
		{
			Name:    "Successfully Register User ",
			UseCase: "To Check user registration when everything is ok",
			mockService: &MockUserService{
				AddUserFn: func(ctx context.Context, user models.User) *structs.IAppError {
					return nil
				},
			},
			returnUserDate:     user,
			ExpectedStatusCode: http.StatusCreated,
		},
		//test case when controller layer should returns error because service mock returned failure
		{
			Name:    "Register User Failure",
			UseCase: "To Check user registration when service layer throws error but json and validation is ok",
			mockService: &MockUserService{
				AddUserFn: func(ctx context.Context, user models.User) *structs.IAppError {
					return &structs.IAppError{
						Message:    "Error While creating user",
						StatusCode: http.StatusInternalServerError,
						Reason:     "Service layer issue",
						ErrorObj:   nil,
					}
				},
			},
			returnUserDate:     user,
			ExpectedStatusCode: http.StatusInternalServerError,
		},
		//test case when service layer mock is ok but invalid json provided to controller layer
		{
			Name:    "User Registration Failure for invalid json",
			UseCase: "To Check user registration when service layer throws error but json and validation is ok",
			mockService: &MockUserService{
				AddUserFn: func(ctx context.Context, user models.User) *structs.IAppError {
					return nil
				},
			},
			returnUserDate: `{
		name":     "John",
		"email":    "Saqlain@gmail.com",
		"address":  "soura",
		"password": "Saqlain@123",
	    "mobile":9797798243,
		role":"user",
		"age:21,
		"pincode":190011,
	}`,
			ExpectedStatusCode: http.StatusBadRequest,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			controller := ReturnNewController(tc.mockService)

			//test for mockservice1 when service layer doesnt return any error
			req := httptest.NewRequest("POST", "/register", strings.NewReader(tc.returnUserDate))
			req.Header.Set("Content-Type", "application/json")
			res := httptest.NewRecorder()
			controller.RegisterUser(res, req)

			if res.Code != tc.ExpectedStatusCode {
				fmt.Print(res.Body)
				t.Fatalf("expected %d, got %d", tc.ExpectedStatusCode, res.Code)
			}

		})
	}

}

func ReturnLoginDetails() string {
	return `{
	"email":"Saqlain@gmail.com",
	"password":"Saqlain@123"
	}`
}
func TestUserLogin(t *testing.T) {

	testCases := []struct {
		Name               string
		Usecase            string
		ExpectedStatusCode int
		mockService        *MockUserService
		loginDetails       string
	}{
		{
			Name:               "User Login failure",
			Usecase:            "User Login failure for invalid json",
			ExpectedStatusCode: http.StatusBadRequest,
			mockService: &MockUserService{
				LoginUserFn: func(ctx context.Context, email, password string) (string, *structs.IAppError) {
					return "Logged in", nil
				},
			},
			loginDetails: `{
	        email":"Saqlain@gmail.com",
	        "password:"Saqlain@123"
	       }`,
		},
		{
			Name:               "User Login Successful ",
			Usecase:            "user login successfull when everything is ok",
			ExpectedStatusCode: http.StatusOK,
			mockService: &MockUserService{
				LoginUserFn: func(ctx context.Context, email, password string) (string, *structs.IAppError) {
					return "Logged in", nil
				},
			},
			loginDetails: ReturnLoginDetails(),
		},
		{
			Name:               "User Login Failure ",
			Usecase:            "user login failure when mock service returns error",
			ExpectedStatusCode: http.StatusInternalServerError,
			mockService: &MockUserService{
				LoginUserFn: func(ctx context.Context, email, password string) (string, *structs.IAppError) {
					return "", &structs.IAppError{
						Message:    "Error from mock service layer",
						StatusCode: http.StatusInternalServerError,
						ErrorObj:   nil,
						Reason:     "Mock Service  Returned Error",
					}
				},
			},
			loginDetails: ReturnLoginDetails(),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			req := httptest.NewRequest("POST", "/login", strings.NewReader(tc.loginDetails))
			req.Header.Set("Content-Type", "application/json")
			res := httptest.NewRecorder()
			controller := ReturnNewController(tc.mockService)
			controller.LoginUser(res, req)

			if res.Code != tc.ExpectedStatusCode {
				fmt.Print(res.Body)
				t.Fatalf("Expected %d got %d for usecase %s", tc.ExpectedStatusCode, res.Code, tc.Usecase)
			}
		})
	}
}

func ReturnDummyUser() models.User {
	return models.User{
		Name:             "Saqlain mushtaq",
		Age:              23,
		Address:          "Soura",
		Password:         "",
		Email:            "Saqlain@gmail.com",
		Mobile:           9797798243,
		Pincode:          190011,
		RegistrationDate: time.Now(),
		ID:               primitive.NewObjectID(),
		AppointmentIDS:   nil,
		Appointments:     nil,
		Role:             "User",
	}

}
func TestSearchUser(t *testing.T) {
	testCases := []struct {
		Name               string
		UseCase            string
		ExpectedStatusCode int
		userID             string
		mockService        *MockUserService
	}{
		{
			Name:               "Search User Returns OK ",
			UseCase:            "To Check User is successfully fetched when everything is ok",
			ExpectedStatusCode: http.StatusOK,
			userID:             primitive.NewObjectID().Hex(),
			mockService: &MockUserService{
				SearchUserByIDFn: func(ctx context.Context, userID primitive.ObjectID) (*models.User, *structs.IAppError) {
					user := ReturnDummyUser()
					return &user, nil
				},
			},
		},
		{

			Name:               "Search User failed",
			UseCase:            "To Check User search fails because mock service returns error rest details are fine",
			ExpectedStatusCode: http.StatusInternalServerError,
			userID:             primitive.NewObjectID().Hex(),
			mockService: &MockUserService{
				SearchUserByIDFn: func(ctx context.Context, userID primitive.ObjectID) (*models.User, *structs.IAppError) {
					return nil, &structs.IAppError{
						Message:    "error from mock service",
						StatusCode: http.StatusInternalServerError,
						Reason:     "Mock Server returned error",
						ErrorObj:   nil,
					}
				},
			},
		},

		{
			Name:               "Search User failed",
			UseCase:            "To Check User search fails because userId is nill objectid",
			ExpectedStatusCode: http.StatusBadRequest,
			userID:             primitive.NilObjectID.Hex(),
			mockService: &MockUserService{
				SearchUserByIDFn: func(ctx context.Context, userID primitive.ObjectID) (*models.User, *structs.IAppError) {
					user := ReturnDummyUser()
					return &user, nil
				},
			},
		},
		{
			Name:               "Search User failed",
			UseCase:            "To Check User search fails because userId is not a valid objectid",
			ExpectedStatusCode: http.StatusBadRequest,
			userID:             "Invalid userid",
			mockService: &MockUserService{
				SearchUserByIDFn: func(ctx context.Context, userID primitive.ObjectID) (*models.User, *structs.IAppError) {
					user := ReturnDummyUser()
					return &user, nil
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			ctx := context.WithValue(context.Background(), middleware.ContextUserIDKey, tc.userID)
			req := httptest.NewRequestWithContext(ctx, "Get", "/user", strings.NewReader(""))
			res := httptest.NewRecorder()

			controller := ReturnNewController(tc.mockService)
			controller.SearchUser(res, req)

			if res.Code != tc.ExpectedStatusCode {
				fmt.Print(res.Body)
				t.Fatalf("Expected %d got %d for usecase %s", tc.ExpectedStatusCode, res.Code, tc.UseCase)
			}

		})
	}
}
