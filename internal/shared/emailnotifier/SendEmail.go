// Package emailnotifier provides functions for sending email through email
package emailnotifier

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	constants "github.com/AlladinDev/AlShifa/internal/shared/constants"
	"github.com/AlladinDev/AlShifa/internal/shared/interfaces"
)

type Sender struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}
type Recipient struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}
type EmailRequest struct {
	Subject    string         `json:"subject"`
	Sender     Sender         `json:"sender"`
	To         []Recipient    `json:"to"`
	TemplateID int            `json:"templateId"`
	Params     map[string]any `json:"params"`
}

type EmailSender struct {
}

func (e *EmailSender) SendNotification(channel string, message string, title string, info string) error {
	if channel == "" || title == "" || info == "" {

		panic("missing params in sending notification")
	}

	emailServiceURL := os.Getenv(constants.EmailSendingURL)

	if emailServiceURL == "" {
		return errors.New("email service url to hit missing in env file")
	}

	emailSecretKey := os.Getenv(constants.EmailSendingAPIKey)

	if emailSecretKey == "" {
		return errors.New("emailSecretKey missing in env file")
	}

	payload := EmailRequest{
		Sender: Sender{
			Name:  "AlShifa Platform",
			Email: "bhatsaqlain982@gmail.com",
		},
		Subject: title,
		To: []Recipient{
			{
				Email: channel,
				Name:  "Dear User",
			},
		},
		TemplateID: 2, // Replace with your actual template ID
		Params: map[string]any{
			"MESSAGE": message,
			"TITLE":   title,
			"INFO":    info,
		},
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("error marshalling payload: %w", err)
	}
	req, err := http.NewRequest("POST", emailServiceURL, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}
	req.Header.Set("accept", "application/json")
	req.Header.Set("content-type", "application/json")
	req.Header.Set("api-key", emailSecretKey)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response: %w", err)
	}
	fmt.Printf("Status: %d\n", resp.StatusCode)
	fmt.Printf("Response: %s\n", string(respBody))
	return nil
}

func NewEmailNotifier() *EmailSender {
	return &EmailSender{}
}

var _ interfaces.INotifier = (*EmailSender)(nil)
