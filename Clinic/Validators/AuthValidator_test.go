package validators

import (
	"AlShifa/Clinic/models"
	utils "AlShifa/Utils"
	"testing"
)

func NewOwnerDetails() models.Owner {
	return models.Owner{
		Name:     "Saqlain",
		Email:    "Saqlain@gmail.com",
		Password: "Saqlain@123",
		Mobile:   9797798243,
		Address:  "Soura Srinagar",
		Gender:   "Male",
	}
}
func TestValidateOwnerDetails(t *testing.T) {
	testCases := []struct {
		Name               string
		UseCase            string
		ExpectedField      string
		ModifyOwnerDetails func(o *models.Owner)
		ExpectedErrorMsg   string
	}{

		//------owner name validation test cases------
		//owner name is missing
		{
			Name:    "OwnerName_Missing",
			UseCase: "Owner name is required",
			ModifyOwnerDetails: func(o *models.Owner) {
				o.Name = ""
			},
			ExpectedField:    "name",
			ExpectedErrorMsg: utils.NameMissingErrMsg,
		},

		//owner name is too short
		{
			Name:    "OwnerName_Short",
			UseCase: "Owner name is too short",
			ModifyOwnerDetails: func(o *models.Owner) {
				o.Name = utils.GenerateRandomString(utils.MinNameLength - 1)
			},
			ExpectedField:    "name",
			ExpectedErrorMsg: utils.ShortNameErrMsg,
		},
		//owner name is too long
		{
			Name:    "OwnerName_long",
			UseCase: "Owner name is too Long",
			ModifyOwnerDetails: func(o *models.Owner) {
				o.Name = utils.GenerateRandomString(utils.MaxNameLength + 1)
			},
			ExpectedField:    "name",
			ExpectedErrorMsg: utils.LongNameErrMsg,
		},

		//-----email validations test cases-----
		// email missing
		{
			Name:    "OwnerEmail_Missing",
			UseCase: "Owner email is required",
			ModifyOwnerDetails: func(o *models.Owner) {
				o.Email = ""
			},
			ExpectedField:    "email",
			ExpectedErrorMsg: utils.EmailMissingErrMsg,
		},
		// invalid email format
		{
			Name:    "OwnerEmail_Invalid",
			UseCase: "Owner email format is invalid",
			ModifyOwnerDetails: func(o *models.Owner) {
				o.Email = "invalidemail.com"
			},
			ExpectedField:    "email",
			ExpectedErrorMsg: utils.InvalidEmailFormatMsg,
		},
		//Email is too long
		{
			Name:    "OwnerEmail_TooLong",
			UseCase: "Owner email is too long",
			ModifyOwnerDetails: func(o *models.Owner) {
				o.Email = utils.GenerateRandomString(utils.MaxEmailLength + 1)
			},
			ExpectedField:    "email",
			ExpectedErrorMsg: utils.LongEmailErrMsg,
		},
		//Email is too small
		{
			Name:    "OwnerEmail_TooSmall",
			UseCase: "Owner email is too small",
			ModifyOwnerDetails: func(o *models.Owner) {
				o.Email = utils.GenerateRandomString(utils.MinEmailLength - 1)
			},
			ExpectedField:    "email",
			ExpectedErrorMsg: utils.ShortEmailErrMsg,
		},

		//-----password validations test cases-----
		//password missing
		{
			Name:    "OwnerPassword_Missing",
			UseCase: "Owner password is required",
			ModifyOwnerDetails: func(o *models.Owner) {
				o.Password = ""
			},
			ExpectedField:    "password",
			ExpectedErrorMsg: utils.PasswordMissingErrMsg,
		},

		//Weak Password validations
		//password doesnt contain lowercase
		{
			Name:    "Owner password_weak",
			UseCase: "Owner password doesnt contain small case",
			ModifyOwnerDetails: func(o *models.Owner) {
				o.Password = "SAQLAIN@23"
			},
			ExpectedField:    "password",
			ExpectedErrorMsg: utils.PasswordWeakErrMsg,
		},
		//password doesnt contain uppercase
		{
			Name:    "Owner password_weak",
			UseCase: "Owner password doesnt contain upper case",
			ModifyOwnerDetails: func(o *models.Owner) {
				o.Password = "saqlain@23"
			},
			ExpectedField:    "password",
			ExpectedErrorMsg: utils.PasswordWeakErrMsg,
		},
		//password doesnt contain numbers
		{
			Name:    "Owner password_weak",
			UseCase: "Owner password doesnt contain numbers",
			ModifyOwnerDetails: func(o *models.Owner) {
				o.Password = "Saqlain@"
			},
			ExpectedField:    "password",
			ExpectedErrorMsg: utils.PasswordWeakErrMsg,
		},
		//password doesnt contain special characters
		{
			Name:    "Owner password_weak",
			UseCase: "Owner password doesnt contain special characters",
			ModifyOwnerDetails: func(o *models.Owner) {
				o.Password = "Saqlain23"
			},
			ExpectedField:    "password",
			ExpectedErrorMsg: utils.PasswordWeakErrMsg,
		},
		//password is too long
		{
			Name:    "OwnerPassword_TooLong",
			UseCase: "Owner password is too long",
			ModifyOwnerDetails: func(o *models.Owner) {
				o.Password = utils.GenerateRandomString(utils.MaxPasswordLength + 1)
			},
			ExpectedField:    "password",
			ExpectedErrorMsg: utils.LongPasswordErrMsg,
		},
		//password is too short
		{
			Name:    "OwnerPassword_TooShort",
			UseCase: "Owner password is too short",
			ModifyOwnerDetails: func(o *models.Owner) {
				o.Password = utils.GenerateRandomString(utils.MinPasswordLength - 1)
			},
			ExpectedField:    "password",
			ExpectedErrorMsg: utils.ShortPasswordErrMsg,
		},

		//---mobile number validation test cases---
		//owner mobile is missing
		{
			Name:    "OwnerMobile_Missing",
			UseCase: "Owner mobile number is required",
			ModifyOwnerDetails: func(o *models.Owner) {
				o.Mobile = 0
			},

			ExpectedField:    "mobile",
			ExpectedErrorMsg: utils.InvalidMobileNumberMsg,
		},
		//mobile is too long
		{
			Name:    "OwnerMobile_TooLong",
			UseCase: "Owner mobile number is too long",
			ModifyOwnerDetails: func(o *models.Owner) {
				o.Mobile = 97977982431
			},

			ExpectedField:    "mobile",
			ExpectedErrorMsg: utils.InvalidMobileNumberMsg,
		},
		//mobile is too short
		{
			Name:    "OwnerMobile_TooShort",
			UseCase: "Owner mobile number is too Short",
			ModifyOwnerDetails: func(o *models.Owner) {
				o.Mobile = 97
			},

			ExpectedField:    "mobile",
			ExpectedErrorMsg: utils.InvalidMobileNumberMsg,
		},
		//----Address validation test cases-----
		//address is missing
		{
			Name:    "OwnerAddress_Missing",
			UseCase: "Owner address is required",
			ModifyOwnerDetails: func(o *models.Owner) {
				o.Address = ""
			},
			ExpectedField:    "address",
			ExpectedErrorMsg: utils.AddressMissingErrMsg,
		},
		//address is too long
		{
			Name:    "OwnerAddress_TooLong",
			UseCase: "Owner address is too long",
			ModifyOwnerDetails: func(o *models.Owner) {
				o.Address = utils.GenerateRandomString(utils.MaxAddressLength + 1)
			},
			ExpectedField:    "address",
			ExpectedErrorMsg: utils.LongAddressErrMsg,
		},
		//address is too short
		{
			Name:    "OwnerAddress_TooShort",
			UseCase: "Owner address is too short",
			ModifyOwnerDetails: func(o *models.Owner) {
				o.Address = utils.GenerateRandomString(utils.MinAddressLength - 1)
			},
			ExpectedField:    "address",
			ExpectedErrorMsg: utils.ShortAddressErrMsg,
		},

		//----gender validation test cases
		//gender is missing
		{
			Name:    "OwnerGender_Missing",
			UseCase: "Owner gender Missing validation",
			ModifyOwnerDetails: func(o *models.Owner) {
				o.Gender = ""
			},
			ExpectedField:    "gender",
			ExpectedErrorMsg: utils.GenderMissingErrMsg,
		},
		//gender is invalid not withing male female or other
		{
			Name:    "OwnerGender_Invalid",
			UseCase: "Owner gender invalid validation",
			ModifyOwnerDetails: func(o *models.Owner) {
				o.Gender = "asdasjdhs"
			},
			ExpectedField:    "gender",
			ExpectedErrorMsg: utils.InvalidGenderErrMsg,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			ownerDetails := NewOwnerDetails()

			//introduce some edge cases which should cause some errors
			tc.ModifyOwnerDetails(&ownerDetails)

			validationErrors := ValidateOwnerDetails(&ownerDetails)

			//all cases have errors so this statement shouldnt run if it runs it means validation function is allowing some edge case
			if validationErrors == nil {
				t.Fatalf("[%s] expected validation error but got none", tc.UseCase)
			}

			//now check if error contains expected field or not
			_, expectedErrFieldExists := validationErrors[tc.ExpectedField]
			if !expectedErrFieldExists {
				t.Fatalf("[%s] expected error field %s but not found", tc.UseCase, tc.ExpectedField)
			}

			//now check if error message is present or not it  shoudlnt be empty
			errMsg := validationErrors[tc.ExpectedField]
			if errMsg == "" {
				t.Fatalf("[%s] expected error Msg %s but it is empty", tc.UseCase, tc.ExpectedField)
			}

			//now check if error message is as intended in test case or not
			if errMsg != tc.ExpectedErrorMsg {
				t.Fatalf("[%s]  error Msg %s doesnt not match test case expectedErrorMsg %s", tc.UseCase, errMsg, tc.ExpectedErrorMsg)
			}
		})
	}
}
