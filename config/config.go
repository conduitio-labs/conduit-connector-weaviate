package config

import (
	"errors"
	"fmt"
)

var (
	ErrMultipleAuth   = errors.New("only one auth. option can be used (API key or WCS)")
	ErrNoAuth         = errors.New("no authentication provided")
	ErrIncompleteAuth = errors.New("authentication info incomplete")
)

type Config struct {
	// Host of the Weaviate instance.
	Endpoint string `json:"endpoint" validate:"required"`

	// Scheme of the Weaviate instance.
	Scheme string `json:"scheme" default:"https" validate:"inclusion=http|https"`

	//TODO: support additional auth schemes __sL__
	// https://github.com/conduitio-labs/conduit-connector-weaviate/issues/4

	// A Weaviate API key
	APIKey string `json:"apiKey"`

	WCS WCSAuth `json:"wcs"`

	// The class name as defined in the schema.
	// A record will be saved under this class unless
	// it has the `weaviate.class` metadata field.
	Class string `json:"class" validate:"required"`
}

func (c *Config) Validate() error {
	if c.APIKey == "" {
		if c.WCS.Valid() {
			return nil
		}
		if c.WCS.isSet() {
			return fmt.Errorf("WCS: %w", ErrIncompleteAuth)
		}
		return ErrNoAuth
	}

	// API key is set
	if c.WCS.isSet() {
		return ErrMultipleAuth
	}

	return nil
}

type WCSAuth struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (a *WCSAuth) Valid() bool {
	return a.Username != "" && a.Password != ""
}

func (a *WCSAuth) isSet() bool {
	return a.Username != "" || a.Password != ""
}
