package controller

import (
	structs "AlShifa/Structs"
	interfaces "AlShifa/Users/Interfaces"
	models "AlShifa/Users/Models"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MockUserService struct {
	AddUserFn        func(ctx context.Context, user models.User) *structs.IAppError
	LoginUserFn      func(ctx context.Context, email, password string) (string, *structs.IAppError)
	SearchUserFn     func(ctx context.Context, filter bson.M) (*models.User, *structs.IAppError)
	SearchUserByIDFn func(ctx context.Context, userID primitive.ObjectID) (*models.User, *structs.IAppError)
}

var _ interfaces.IService = (*MockUserService)(nil)

func (m *MockUserService) AddUser(ctx context.Context, user models.User) *structs.IAppError {
	if m.AddUserFn != nil {
		return m.AddUserFn(ctx, user)
	}
	return nil
}

func (m *MockUserService) LoginUser(ctx context.Context, email, password string) (string, *structs.IAppError) {
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
