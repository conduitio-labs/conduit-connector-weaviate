// Copyright © 2024 Meroxa, Inc.
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

package destination

import (
	"context"
	"errors"
	"fmt"

	"github.com/conduitio-labs/conduit-connector-weaviate/config"
	sdk "github.com/conduitio/conduit-connector-sdk"
)

type Config struct {
	sdk.DefaultDestinationMiddleware
	config.Config
	// TODO: better naming for this value __sL__
	// Vectorizers which can be configured client side
	// mostly require an API key only.
	// However, OpenAI can also be configured with an organization
	// via the X-OpenAI-Organization header.

	ModuleHeader ModuleHeader `json:"moduleHeader"`
	// Whether a UUID for records should be automatically generated.
	// The generated UUIDs are MD5 sums of record keys.
	GenerateUUID bool `json:"generateUUID"`
}

type ModuleHeader struct {
	// Name of the header configuring a module (e.g. `X-OpenAI-Api-Key`)
	Name string `json:"name"`
	// Value for header given in `moduleHeader.name`.
	Value string `json:"value"`
}

func (m ModuleHeader) IsValid() bool {
	return (m.Name == "" && m.Value == "") ||
		(m.Name != "" && m.Value != "")
}

func (c *Config) Validate(ctx context.Context) error {
	err := c.DefaultDestinationMiddleware.Validate(ctx)
	if err != nil {
		return err
	}

	if !c.ModuleHeader.IsValid() {
		return errors.New("invalid module configuration")
	}

	err = c.Config.Validate()
	if err != nil {
		return fmt.Errorf("invalid configuration: %w", err)
	}

	return nil
}
