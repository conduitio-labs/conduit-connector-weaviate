package config

type Config struct {
	Endpoint string `json:"endpoint" validate:"required"`
	Scheme   string `json:"scheme" validate:"required"`
	ApiKey   string `json:"api_key" validate:omitempty"`
	Class    string `json:"class" validate:"required"`
}
