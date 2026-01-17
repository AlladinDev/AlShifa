package structs

type IAppSuccess struct {
	Message    string `json:"message"`
	Data       any    `json:"data"`
	StatusCode int    `json:"statusCode"`
}
