package apitxt

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/AlladinDev/AlShifa/internal/shared/constants"
	"github.com/AlladinDev/AlShifa/internal/shared/interfaces"
)

type Messanger struct {
	apiClient *http.Client
}

func (m *Messanger) SendNotification(channel string, message string, title string, info string) error {
	apiTxtUrl := os.Getenv(constants.ApiTxtURL)
	apiTxtSecret := os.Getenv(constants.ApiTxtAuthKey)

	if apiTxtSecret == "" || apiTxtUrl == "" {
		panic("apitext secret or apitext url missing from env")
	}

	log.Println(message)

	// Build the full URL with query parameters
	fullURL := fmt.Sprintf("%s?authkey=%s&mobile=%s&otp=%s", apiTxtUrl, apiTxtSecret, channel, message)

	// 1. Construct the request object
	request, err := http.NewRequest(http.MethodGet, fullURL, nil)
	if err != nil {
		return err
	}

	// 2. Set headers (Note: fixed 'req' to 'request')
	request.Header.Set("Content-Type", "application/json")

	// 3. EXECUTE THE REQUEST (Crucial step!)
	resp, err := m.apiClient.Do(request)
	if err != nil {
		log.Println(err)
		return errors.New("failed to send otp for verification")
	}

	// 4. Always close the response body to reuse the TCP connection
	defer resp.Body.Close()

	// Optional: Check if the API returned an error status code
	if resp.Status != "200 OK" {
		return fmt.Errorf("api returned bad status: %s", resp.Status)
	}

	return nil
}

func SetupApiTxtModule() (interfaces.INotifier, error) {
	// Create one shared client to be reused across all notification requests
	client := &http.Client{}

	return &Messanger{
		apiClient: client,
	}, nil
}
