package validators

import (
	structs "AlShifa/Users/Structs"
	"testing"
)

func TestValidateLoginDetails(t *testing.T) {
	loginDetails := structs.LoginDetails{
		Email:    "Saqlain@gmaill.com",
		Password: "Saqlain@123",
	}

	errors := ValidateLoginDetails(&loginDetails)
	if errors != nil {
		t.Errorf("Expected no validation errors as email and password are valid, but got: %v", errors)
	}

	//payload with invalid password but valid email
	invalidPasswordLoginDetailsPayload := structs.LoginDetails{
		Email:    "Saqlain@gmail.com",
		Password: "short",
	}

	errorsInvalid := ValidateLoginDetails(&invalidPasswordLoginDetailsPayload)
	if errorsInvalid == nil {
		t.Errorf("Expected validation errors for invalid password, but got none")
	}

	//now check if errorsInvalid contains password error
	if _, exists := errorsInvalid["password"]; !exists {
		t.Errorf("Expected password validation error field, but it was not found")
	}

	//payload with invalid email but valid password
	invalidEmailLoginDetailsPayload := structs.LoginDetails{
		Email:    "invalid-email",
		Password: "Saqlain@gmail.com",
	}

	errorsInvalidEmail := ValidateLoginDetails(&invalidEmailLoginDetailsPayload)
	if errorsInvalidEmail == nil {
		t.Errorf("Expected validation errors for invalid email, but got none")
	}

	//now check if errorsInvalidEmail contains email error
	if _, exists := errorsInvalidEmail["email"]; !exists {
		t.Errorf("Expected email validation error field, but it was not found")
	}

}
