package config

import (
	"errors"
	"fmt"
)

var (
	ErrUsernamePasswordMissing = errors.New("username or password missing")
)

type Config struct {
	// Host of the Weaviate instance.
	Endpoint string `json:"endpoint" validate:"required"`

	// Scheme of the Weaviate instance.
	Scheme string `json:"scheme" default:"https" validate:"inclusion=http|https"`

	// AuthMechanism specifies in which way the connector will authenticate to Weaviate.
	AuthMechanism string `json:"auth.mechanism" validate:"inclusion=none|apiKey|wcsCredentials" default:"none"`

	// Weaviate API key.
	APIKey string `json:"apiKey"`

	// Weaviate Cloud Services (WCS) credentials.
	WCSCredentials WCSCredentials `json:"wcs"`

	// The class name as defined in the schema.
	// A record will be saved under this class unless
	// it has the `weaviate.class` metadata field.
	Class string `json:"class" validate:"required"`
}

func (c *Config) Validate() error {
	// Validate authentication configuration
	if c.AuthMechanism == "none" {
		return nil
	}

	if c.AuthMechanism == "apiKey" {
		if c.APIKey == "" {
			return errors.New("authMechanism set to 'apiKey', but apiKey not specified")
		}

		return nil
	}

	if c.AuthMechanism == "wcsCredentials" {
		return c.WCSCredentials.Validate()
	}

	return fmt.Errorf("unknown authMechanism %v", c.AuthMechanism)
}

type WCSCredentials struct {
	// WCS username
	Username string `json:"username"`
	// WCS password
	Password string `json:"password"`
}

func (a *WCSCredentials) Validate() error {
	if a.Username == "" || a.Password == "" {
		return ErrUsernamePasswordMissing
	}

	return nil
}
