package config

type Config struct {
	Endpoint string `json:"endpoint" validate:"required"`
	Scheme   string `json:"scheme" validate:"required"`
	//TODO: support additional auth schemes __sL__
	APIKey string `json:"api_key" validate:"required"`
	Class  string `json:"class" validate:"required"`
}
