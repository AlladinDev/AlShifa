package utils_test

import (
	"net/url"
	"testing"

	utils "github.com/AlladinDev/AlShifa/internal/shared/utils"
)

type Doctor struct {
	Name           string   `form:"name"`
	Email          string   `form:"email"`
	Age            int      `form:"age"`
	Salary         float64  `form:"salary"`
	IsVerified     bool     `form:"verified"`
	Qualifications []string `form:"qualifications"`

	Experience *int `form:"experience"`

	// No tag, should fallback to struct field name
	Address string

	// Should never be set
	Password string `form:"-"`
}

func TestUnmarshalFormValues(t *testing.T) {

	form := url.Values{
		"name":           {"Saqlain"},
		"email":          {"saqlain@gmail.com"},
		"age":            {"25"},
		"salary":         {"55000.75"},
		"verified":       {"true"},
		"experience":     {"5"},
		"qualifications": {"MBBS", "MD"},
		"Address":        {"Srinagar"},
		"Password":       {"123456"},
	}

	var doctor Doctor

	utils.UnmarshalFormValues(form, &doctor)

	if doctor.Name != "Saqlain" {
		t.Fatalf("expected Name=Saqlain got %v", doctor.Name)
	}

	if doctor.Email != "saqlain@gmail.com" {
		t.Fatalf("expected Email")
	}

	if doctor.Age != 25 {
		t.Fatalf("expected Age=25 got %v", doctor.Age)
	}

	if doctor.Salary != 55000.75 {
		t.Fatalf("expected Salary=55000.75 got %v", doctor.Salary)
	}

	if doctor.IsVerified != true {
		t.Fatalf("expected Verified=true")
	}

	if doctor.Experience == nil {
		t.Fatal("expected Experience pointer to be initialized")
	}

	if *doctor.Experience != 5 {
		t.Fatalf("expected Experience=5 got %v", *doctor.Experience)
	}

	if len(doctor.Qualifications) != 2 {
		t.Fatalf("expected 2 qualifications got %d", len(doctor.Qualifications))
	}

	if doctor.Qualifications[0] != "MBBS" {
		t.Fatal("qualification 1 incorrect")
	}

	if doctor.Qualifications[1] != "MD" {
		t.Fatal("qualification 2 incorrect")
	}

	if doctor.Address != "Srinagar" {
		t.Fatalf("expected Address=Srinagar got %v", doctor.Address)
	}

	// Should remain empty because of form:"-"
	if doctor.Password != "" {
		t.Fatal("Password should not be populated")
	}
}
