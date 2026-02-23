package controller

import (
	structs "AlShifa/structs"
	interfaces "AlShifa/users/interfaces"
	models "AlShifa/users/models"
	"context"

	sharedModels "AlShifa/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MockUserService struct {
	AddUserFn           func(ctx context.Context, user models.User) *structs.IAppError
	LoginUserFn         func(ctx context.Context, email string, password string) (string, *structs.IAppError)
	SearchUserFn        func(ctx context.Context, filter bson.M) (*models.User, *structs.IAppError)
	SearchUserByIDFn    func(ctx context.Context, userID primitive.ObjectID) (*models.User, *structs.IAppError)
	FetchAppointmentsFn func(ctx context.Context, groupingID string, filter bson.M) ([]sharedModels.Appointment, *structs.IAppError)
}

var _ interfaces.IService = (*MockUserService)(nil)

func (m *MockUserService) AddUser(ctx context.Context, user models.User) *structs.IAppError {
	if m.AddUserFn != nil {
		return m.AddUserFn(ctx, user)
	}
	return nil
}

func (m *MockUserService) FetchAppointments(ctx context.Context, groupingID string, filter bson.M) ([]sharedModels.Appointment, *structs.IAppError) {
	if m.FetchAppointmentsFn != nil {
		return m.FetchAppointmentsFn(ctx, groupingID, filter)
	}
	return nil, nil
}
func (m *MockUserService) LoginUser(ctx context.Context, email string, password string) (string, *structs.IAppError) {
	if m.LoginUserFn != nil {
		return m.LoginUserFn(ctx, email, password)
	}
	return "", nil
}

func (m *MockUserService) SearchUser(ctx context.Context, filter bson.M) (*models.User, *structs.IAppError) {
	if m.SearchUserFn != nil {
		return m.SearchUserFn(ctx, filter)
	}
	return nil, nil
}

func (m *MockUserService) SearchUserByID(ctx context.Context, userID primitive.ObjectID) (*models.User, *structs.IAppError) {
	if m.SearchUserByIDFn != nil {
		return m.SearchUserByIDFn(ctx, userID)
	}
	return nil, nil
}
