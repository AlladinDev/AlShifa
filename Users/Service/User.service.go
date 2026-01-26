// Package service provides service layer functions for user module
package service

import (
	structs "AlShifa/Structs"
	interfaces "AlShifa/Users/Interfaces"
	models "AlShifa/Users/Models"
	utils "AlShifa/Utils"
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Service struct {
	repo interfaces.IRepository
}

var _ interfaces.IService = (*Service)(nil)

func ReturnNewService(repo interfaces.IRepository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) AddUser(ctx context.Context, user models.User) *structs.IAppError {
	//first check if this user exists by this email or mobile number or not

	userExists, userExistsDBErr := s.repo.SearchUser(ctx, bson.M{"email": user.Email, "mobile": user.Mobile})
	//if userExistsDBErr is not nill it means either some db error or no document found error ,handle db related error only if no document found error it is ok
	if userExistsDBErr != nil {
		if userExistsDBErr != mongo.ErrNoDocuments {
			return &structs.IAppError{
				Message:    "Registration Failed",
				StatusCode: 500,
				Reason:     "Server Error",
				ErrorObj:   userExistsDBErr,
			}
		}
	}

	if userExists != nil {
		return &structs.IAppError{
			Message:    "This Email or Mobile Already Exists",
			StatusCode: 400,
			Reason:     "Duplicate Email or Mobile",
			ErrorObj:   errors.New("user already exists"),
		}
	}

	//add default things like userId and registrationDate and null values like appointments
	user.RegistrationDate = time.Now()
	user.ID = primitive.NewObjectID()
	user.AppointmentIDS = nil
	user.Appointments = nil
	user.Role = utils.RoleUser

	//hash password before storing
	hashedPassword, hashErr := utils.HashPasswordArgon2id(user.Password)
	if hashErr != nil {
		return &structs.IAppError{
			Message:    "Failed to Register User",
			ErrorObj:   hashErr,
			StatusCode: 500,
			Reason:     "Server Error",
		}
	}

	user.Password = hashedPassword

	if err := s.repo.RegisterUser(ctx, user); err != nil {
		return &structs.IAppError{
			Message:    "Failed to Register User",
			ErrorObj:   err,
			StatusCode: 500,
		}
	}

	return nil
}

func (s *Service) SearchUserByID(ctx context.Context, userID primitive.ObjectID) (*models.User, *structs.IAppError) {

	user, err := s.repo.SearchUserByID(ctx, userID)
	if err != nil {
		if err != mongo.ErrNoDocuments {
			return nil, &structs.IAppError{
				Message:    "Registration Failed",
				StatusCode: 500,
				ErrorObj:   err,
				Reason:     "Server Error",
			}
		}
		return nil, &structs.IAppError{
			Message:    "User Not Found",
			StatusCode: 404,
			ErrorObj:   err,
			Reason:     "User Doesnt Exist",
		}

	}

	//null password of user before sending it to user
	user.Password = ""
	return user, nil
}

func (s *Service) SearchUser(ctx context.Context, filter bson.M) (*models.User, *structs.IAppError) {
	user, err := s.repo.SearchUser(ctx, filter)
	if err != nil {
		return nil, &structs.IAppError{
			Message:    "Failed To Search User",
			StatusCode: 500,
			Reason:     err.Error(),
			ErrorObj:   err,
		}
	}

	return user, nil
}

func (s *Service) LoginUser(ctx context.Context, email string, password string) (string, *structs.IAppError) {
	user, err := s.repo.SearchUser(ctx, bson.M{"email": email})
	if err != nil {
		fmt.Print(err)
		if err == mongo.ErrNoDocuments {
			return "", &structs.IAppError{
				Message:    "User Doesnt Exist With This email",
				Reason:     "No User Found",
				ErrorObj:   mongo.ErrNoDocuments,
				StatusCode: http.StatusNotFound,
			}
		}
		return "", &structs.IAppError{
			Message:    "Login Failed",
			Reason:     "Server Error",
			ErrorObj:   err,
			StatusCode: http.StatusInternalServerError,
		}
	}

	//now check if password matches or not
	passwordMatches, passwordMatchErr := utils.VerifyPasswordArgon2id(password, user.Password)
	if passwordMatchErr != nil {
		fmt.Print(passwordMatchErr)
		return "", &structs.IAppError{
			Message:    "Login Failed plz try again",
			StatusCode: http.StatusInternalServerError,
			ErrorObj:   errors.New("server error"),
			Reason:     "Server Error",
		}
	}
	if !passwordMatches {
		return "", &structs.IAppError{
			Message:    "Invalid Email or Password",
			StatusCode: http.StatusBadRequest,
			ErrorObj:   errors.New("login failed invalid email or password"),
			Reason:     "invalid details",
		}
	}

	//now generate token
	jwtToken, tokenErr := utils.GenerateJWT(user.ID.Hex(), user.Role)
	if tokenErr != nil {
		return "", &structs.IAppError{
			Message:    "Login Failed",
			Reason:     "Server Error",
			ErrorObj:   errors.New("failed to login user : server error"),
			StatusCode: http.StatusInternalServerError,
		}
	}

	return jwtToken, nil
}
