package emailnotifier

import (
	utils "AlShifa/Utils"
	"testing"
)

func TestRealEmailSending(t *testing.T) {
	//os.Setenv("GODEBUG", "netdns=go")
	//compulsory load envs first
	// os.Clearenv()
	utils.LoadEnvs()

	//appEmail := os.Getenv("APP_EMAIL")

	EmailService, err := NewEmailService()
	if err != nil {
		t.Fatalf("Expected no error but got %v when getting  real email service", err)
	}

	emailSendingErr := EmailService.SendNotification("shizuka@gmail.com", "hey buddy testing")
	if emailSendingErr != nil {
		t.Fatalf("Expected no error when sending real email but got %v", emailSendingErr)
	}
}
