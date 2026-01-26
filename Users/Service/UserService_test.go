package service

import (
	structs "AlShifa/Structs"
	models "AlShifa/Users/Models"
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var dummyTime = time.Now()
var dummyObjectID = primitive.NewObjectID()

func ReturnDummyUser() models.User {
	return models.User{
		Name:             "Saqlain",
		Age:              23,
		Password:         "Saqlain@123",
		Email:            "Saqlain@gmail.com",
		Address:          "Soura",
		RegistrationDate: dummyTime,
		ID:               dummyObjectID,
		Role:             "user",
		AppointmentIDS:   nil,
		Appointments:     nil,
	}
}

func TestUserServiceAddUser(t *testing.T) {
	testCases := []struct {
		Name        string
		Data        models.User
		ExpectedErr *structs.IAppError
		mockRepo    *MockUserRepo
	}{
		//test case when everything is ok and service also didnt returned any error
		{Name: "User Registration Successfull When everything is ok",
			Data:        ReturnDummyUser(),
			ExpectedErr: nil,
			mockRepo: &MockUserRepo{
				AddUserFn: func(ctx context.Context, user models.User) error {
					return nil
				},
				SearchUserFn: func(ctx context.Context, filter bson.M) (*models.User, error) {
					return nil, nil
				},
			},
		},

		//test case for registering failure when everything is ok but mock repo fn AddUser returned error
		{
			Name: "User Registration Failed Because Mock Repo function Returned Error else everything is ok",
			Data: ReturnDummyUser(),
			ExpectedErr: &structs.IAppError{
				Message:    "Failed to Register User",
				ErrorObj:   errors.New("error from mocked repo"),
				StatusCode: 500,
			},
			mockRepo: &MockUserRepo{
				AddUserFn: func(ctx context.Context, user models.User) error {
					return errors.New("error from mocked repo")
				},
				SearchUserFn: func(ctx context.Context, filter bson.M) (*models.User, error) {
					return nil, nil
				},
			},
		},

		//test case for registering failure when everything is ok but mock repo fn SearchUser returned error
		{
			Name: "User Registration Failed Because Mock Repo function SearchUser Returned Error else everything is ok",
			Data: ReturnDummyUser(),
			ExpectedErr: &structs.IAppError{
				Message:    "Registration Failed",
				StatusCode: 500,
				Reason:     "Server Error",
				ErrorObj:   errors.New("error from mock repo failed to search user"),
			},
			mockRepo: &MockUserRepo{
				AddUserFn: func(ctx context.Context, user models.User) error {
					return nil
				},
				SearchUserFn: func(ctx context.Context, filter bson.M) (*models.User, error) {
					return nil, errors.New("error from mock repo failed to search user")
				},
			},
		},

		//Test case for duplicate email error
		{
			Name: "User Registration Failed Because Email Already exists",
			Data: ReturnDummyUser(),
			ExpectedErr: &structs.IAppError{
				Message:    "This Email or Mobile Already Exists",
				StatusCode: 400,
				Reason:     "Duplicate Email or Mobile",
				ErrorObj:   errors.New("user already exists"),
			},
			mockRepo: &MockUserRepo{
				SearchUserFn: func(ctx context.Context, filter bson.M) (*models.User, error) {
					user := ReturnDummyUser()
					return &user, nil
				},
				AddUserFn: func(ctx context.Context, user models.User) error {
					return nil
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			service := ReturnNewService(tc.mockRepo)
			err := service.AddUser(context.Background(), tc.Data)
			if !reflect.DeepEqual(err, tc.ExpectedErr) {
				t.Fatalf("Expected %v to be equal to %v ", err, tc.ExpectedErr)
			}

		})
	}
}

func TestSearchUserByID(t *testing.T) {
	testCases := []struct {
		Name         string
		UserID       primitive.ObjectID
		mockRepo     *MockUserRepo
		expectedUser *models.User
		expectedErr  *structs.IAppError
	}{
		//test case when everything is ok and mock repo should return a user
		{
			Name:   "Searched User Successfully when everything is ok",
			UserID: primitive.NewObjectID(),
			mockRepo: &MockUserRepo{
				SearchUserByIDFn: func(ctx context.Context, userID primitive.ObjectID) (*models.User, error) {
					user := ReturnDummyUser()
					return &user, nil
				},
			},
			expectedUser: func() *models.User {
				user := ReturnDummyUser()
				user.Password = ""
				return &user
			}(),
			expectedErr: nil,
		},
		//test case when it search user fails because mocked repo function returns some server error
		{
			Name:   "Searched User failed because mocked repo function returns some server error",
			UserID: primitive.NewObjectID(),
			mockRepo: &MockUserRepo{
				SearchUserByIDFn: func(ctx context.Context, userID primitive.ObjectID) (*models.User, error) {
					return nil, errors.New("server error from mock repo")
				},
			},
			expectedUser: nil,
			expectedErr: &structs.IAppError{
				Message:    "Registration Failed",
				StatusCode: 500,
				ErrorObj:   errors.New("server error from mock repo"),
				Reason:     "Server Error",
			},
		},

		//test case when  search user fails because mocked repo fn returned nill as user means user not found
		{
			Name:   "Searched User failed because repo fn returned nill as user means user not found ",
			UserID: primitive.NewObjectID(),
			mockRepo: &MockUserRepo{
				SearchUserByIDFn: func(ctx context.Context, userID primitive.ObjectID) (*models.User, error) {
					return nil, mongo.ErrNoDocuments
				},
			},
			expectedUser: nil,
			expectedErr: &structs.IAppError{
				Message:    "User Not Found",
				StatusCode: 404,
				ErrorObj:   mongo.ErrNoDocuments,
				Reason:     "User Doesnt Exist",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			service := ReturnNewService(tc.mockRepo)
			user, err := service.SearchUserByID(context.Background(), tc.UserID)

			if !reflect.DeepEqual(tc.expectedErr, err) {
				t.Fatalf("expected errors %v to be equal to %v", tc.expectedErr, err)
			}

			if !reflect.DeepEqual(user, tc.expectedUser) {
				t.Fatalf("expected user %v to be equal to %v", tc.expectedUser, user)
			}
		})
	}
}

func TestSearchUser(t *testing.T) {
	testCases := []struct {
		Name         string
		UserID       primitive.ObjectID
		mockRepo     *MockUserRepo
		expectedUser *models.User
		expectedErr  *structs.IAppError
	}{
		//test case when everything is ok and mock repo should return a user
		{
			Name:   "Searched User Successfully when everything is ok",
			UserID: primitive.NewObjectID(),
			mockRepo: &MockUserRepo{
				SearchUserFn: func(ctx context.Context, filter bson.M) (*models.User, error) {
					user := ReturnDummyUser()
					return &user, nil
				},
			},
			expectedUser: func() *models.User {
				user := ReturnDummyUser()
				return &user
			}(),
			expectedErr: nil,
		},
		//test case when  search user fails because mocked repo function returns some server error
		{
			Name:   "Searched User failed because mocked repo function returns some server error",
			UserID: primitive.NewObjectID(),
			mockRepo: &MockUserRepo{
				SearchUserFn: func(ctx context.Context, filter bson.M) (*models.User, error) {
					return nil, errors.New("server error from mock repo")
				},
			},
			expectedUser: nil,
			expectedErr: &structs.IAppError{
				Message:    "Failed To Search User",
				StatusCode: 500,
				Reason:     "server error from mock repo",
				ErrorObj:   errors.New("server error from mock repo"),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			service := ReturnNewService(tc.mockRepo)
			user, err := service.SearchUser(context.Background(), bson.M{})

			if !reflect.DeepEqual(tc.expectedErr, err) {
				t.Fatalf("expected errors %v to be equal to %v", tc.expectedErr, err)
			}

			if !reflect.DeepEqual(user, tc.expectedUser) {
				t.Fatalf("expected user %v to be equal to %v", tc.expectedUser, user)
			}
		})
	}
}
