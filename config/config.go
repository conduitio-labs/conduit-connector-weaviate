// Copyright © 2022 Meroxa, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package config

import (
	"errors"
	"fmt"
)

var ErrUsernamePasswordMissing = errors.New("username or password missing")

type Config struct {
	// Host of the Weaviate instance.
	Endpoint string `json:"endpoint" validate:"required"`

	// Scheme of the Weaviate instance.
	Scheme string `json:"scheme" default:"https" validate:"inclusion=http|https"`

	Auth Auth `json:"auth"`

	// The class name as defined in the schema.
	// A record will be saved under this class unless
	// it has the `weaviate.class` metadata field.
	Class string `json:"class" validate:"required"`
}

func (c *Config) Validate() error {
	return c.Auth.Validate()
}

type Auth struct {
	// Mechanism specifies in which way the connector will authenticate to Weaviate.
	Mechanism string `json:"mechanism" validate:"inclusion=none|apiKey|wcsCreds" default:"none"`

	// A Weaviate API key.
	APIKey string `json:"apiKey"`

	// Weaviate Cloud Services (WCS) credentials.
	WCSCredentials WCSCredentials `json:"wcsCreds"`
}

func (a Auth) Validate() error {
	if a.Mechanism == "none" {
		return nil
	}

	if a.Mechanism == "apiKey" {
		if a.APIKey == "" {
			return errors.New("authMechanism set to 'apiKey', but apiKey not specified")
		}

		return nil
	}

	if a.Mechanism == "wcsCreds" {
		return a.WCSCredentials.Validate()
	}

	return fmt.Errorf("unknown auth mechanism %v", a.Mechanism)
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
