package service

import (
	interfaces "AlShifa/Users/Interfaces"
	models "AlShifa/Users/Models"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MockUserRepo struct {
	AddUserFn        func(ctx context.Context, user models.User) error
	LoginUserFn      func(ctx context.Context, email string, password string) (string, error)
	SearchUserFn     func(ctx context.Context, filter bson.M) (*models.User, error)
	SearchUserByIDFn func(ctx context.Context, userID primitive.ObjectID) (*models.User, error)
}

var _ interfaces.IRepository = (*MockUserRepo)(nil)

func (m *MockUserRepo) RegisterUser(ctx context.Context, user models.User) error {
	if m.AddUserFn == nil {
		panic("AddUserFn not implemented inside mock")
	}
	return m.AddUserFn(ctx, user)
}

func (m *MockUserRepo) LoginUser(ctx context.Context, email string, password string) (string, error) {
	if m.LoginUserFn == nil {
		panic("LoginUser not implemented inside mock")
	}
	return m.LoginUserFn(ctx, email, password)
}

func (m *MockUserRepo) SearchUser(ctx context.Context, filter bson.M) (*models.User, error) {
	if m.SearchUserFn == nil {
		panic("SearchUser not implemented inside mock")
	}
	return m.SearchUserFn(ctx, filter)
}

func (m *MockUserRepo) SearchUserByID(ctx context.Context, userID primitive.ObjectID) (*models.User, error) {
	if m.SearchUserByIDFn == nil {
		panic("SearchUserByID not implemented inside mock")
	}
	return m.SearchUserByIDFn(ctx, userID)
}
