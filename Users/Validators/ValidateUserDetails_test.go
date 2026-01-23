package validators

import (
	models "AlShifa/Users/Models"
	utils "AlShifa/Utils"
	"fmt"
	"testing"
)

func ReturnNewUserValidDetails() models.User {
	return models.User{
		Name:     "Saqlain",
		Email:    "protonium789@gmail.com",
		Password: "Saqlain@123",
		Mobile:   9797798243,
		Address:  "Soura srinagar",
		Pincode:  190011,
		Age:      23,
	}
}
func TestValidateUserDetails(t *testing.T) {
	testCases := []struct {
		name               string
		useCase            string
		modifyUser         func(u *models.User)
		expectedField      string
		expectedErrMessage string
	}{
		// ---------- NAME ----------
		{
			name:    "name_short",
			useCase: "Short Name Validation",
			modifyUser: func(u *models.User) {
				u.Name = utils.GenerateRandomString(utils.MinNameLength - 1)
			},
			expectedField:      "name",
			expectedErrMessage: "name length is invalid",
		},
		{
			name:    "name_long",
			useCase: "Long Name Validation",
			modifyUser: func(u *models.User) {
				u.Name = utils.GenerateRandomString(utils.MaxNameLength + 1)
			},
			expectedField:      "name",
			expectedErrMessage: "name length is invalid",
		},

		// ---------- PASSWORD ----------
		{
			name:    "password_empty",
			useCase: "Empty Password Validation",
			modifyUser: func(u *models.User) {
				u.Password = ""
			},
			expectedField:      "password",
			expectedErrMessage: "password is required",
		},
		//short password test case
		{
			name:    "password_short",
			useCase: "Short Password Validation",
			modifyUser: func(u *models.User) {
				u.Password = utils.GenerateRandomString(utils.MinPasswordLength - 1)
			},
			expectedField:      "password",
			expectedErrMessage: "password is too short",
		},
		//long password test case
		{
			name:    "password_long",
			useCase: "Long Password Validation",
			modifyUser: func(u *models.User) {
				u.Password = utils.GenerateRandomString(utils.MaxPasswordLength + 1)
			},
			expectedField:      "password",
			expectedErrMessage: "password is too long",
		},
		{
			name:    "password_missing_uppercase",
			useCase: "Password Missing Uppercase Letter",
			modifyUser: func(u *models.User) {
				u.Password = "lowercase1@"
			},
			expectedField:      "password",
			expectedErrMessage: "password must contain upper, lower, digit and special character",
		},
		{
			name:    "password_missing_lowercase",
			useCase: "Password Missing Lowercase Letter",
			modifyUser: func(u *models.User) {
				u.Password = "UPPERCASE1@"
			},
			expectedField:      "password",
			expectedErrMessage: "password must contain upper, lower, digit and special character",
		},
		{
			name:    "password_missing_digit",
			useCase: "Password Missing Digit",
			modifyUser: func(u *models.User) {
				u.Password = "NoDigit@"
			},
			expectedField:      "password",
			expectedErrMessage: "password must contain upper, lower, digit and special character",
		},
		{
			name:    "password_missing_special",
			useCase: "Password Missing Special Character",
			modifyUser: func(u *models.User) {
				u.Password = "NoSpecial1"
			},
			expectedField:      "password",
			expectedErrMessage: "password must contain upper, lower, digit and special character",
		},

		// ---------- EMAIL ----------
		{
			name:    "email_invalid",
			useCase: "Invalid Email Format",
			modifyUser: func(u *models.User) {
				u.Email = "invalid-email"
			},
			expectedField:      "email",
			expectedErrMessage: "invalid email format",
		},
		{
			name:    "email_empty",
			useCase: "Empty Email Validation",
			modifyUser: func(u *models.User) {
				u.Email = ""
			},
			expectedField:      "email",
			expectedErrMessage: "email is required",
		},

		// ---------- MOBILE ----------
		{
			name:    "mobile_short",
			useCase: "Short Mobile Number",
			modifyUser: func(u *models.User) {
				u.Mobile = 12345
			},
			expectedField:      "mobile",
			expectedErrMessage: "mobile number must be 10 digits",
		},
		{
			name:    "mobile_long",
			useCase: "Long Mobile Number",
			modifyUser: func(u *models.User) {
				u.Mobile = 1234567890123
			},
			expectedField:      "mobile",
			expectedErrMessage: "mobile number must be 10 digits",
		},
		//invalid mobile number test case eg -9797798243 in negative
		{
			name:    "mobile_negative",
			useCase: "Negative Mobile Number",
			modifyUser: func(u *models.User) {
				u.Mobile = -979779824
			},
			expectedField:      "mobile",
			expectedErrMessage: "mobile number must be positive",
		},

		// ---------- PINCODE ----------
		{
			name:    "pincode_invalid",
			useCase: "Invalid Pincode Length",
			modifyUser: func(u *models.User) {
				u.Pincode = 123
			},
			expectedField:      "pincode",
			expectedErrMessage: "pincode must be 6 digits",
		},

		// ---------- AGE ----------
		{
			name:    "age_invalid_short",
			useCase: "Invalid Age",
			modifyUser: func(u *models.User) {
				u.Age = 0
			},
			expectedField:      "age",
			expectedErrMessage: fmt.Sprintf("age must be between  %d and %d", utils.MinAge, utils.MaxAge),
		},
		//max age validation
		{
			name:    "age_invalid_long",
			useCase: "Max  Age Validation",
			modifyUser: func(u *models.User) {
				u.Age = utils.MaxAge + 1
			},
			expectedField:      "age",
			expectedErrMessage: fmt.Sprintf("age must be between  %d and %d", utils.MinAge, utils.MaxAge),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			user := ReturnNewUserValidDetails()
			tc.modifyUser(&user)

			validationErrors := ValidateUser(&user)

			if validationErrors == nil {
				t.Fatalf("[%s] Expected validation errors, got none", tc.useCase)
			}

			errMsg, exists := validationErrors[tc.expectedField]
			if !exists {
				t.Fatalf("[%s] Expected error for field %q, got %v",
					tc.useCase, tc.expectedField, validationErrors)
			}

			if errMsg != tc.expectedErrMessage {
				t.Fatalf("[%s] Expected error message %q, got %q",
					tc.useCase, tc.expectedErrMessage, errMsg)
			}
		})
	}
}
