// Copyright © 2023 Meroxa, Inc.
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

package config_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	weaviate "github.com/conduitio-labs/conduit-connector-weaviate"
	"github.com/conduitio-labs/conduit-connector-weaviate/config"
	"github.com/conduitio-labs/conduit-connector-weaviate/destination"
	sdk "github.com/conduitio/conduit-connector-sdk"
	"github.com/matryer/is"
)

func TestConfig_Auth(t *testing.T) {
	testCases := []struct {
		name    string
		cfgMap  map[string]string
		wantCfg destination.Config
		wantErr error
	}{
		{
			name: "no auth is possible",
			cfgMap: map[string]string{
				"endpoint":       "test-endpoint",
				"scheme":         "https",
				"class":          "test-class",
				"auth.mechanism": "none",
			},
			wantCfg: destination.Config{
				Config: config.Config{
					Endpoint: "test-endpoint",
					Scheme:   "https",
					Class:    "test-class",
					Auth: config.Auth{
						Mechanism: "none",
					},
				},
			},
		},
		{
			name: "API key only",
			cfgMap: map[string]string{
				"endpoint":       "test-endpoint",
				"scheme":         "https",
				"class":          "test-class",
				"auth.mechanism": "apiKey",
				"auth.apiKey":    "xyz",
			},
			wantCfg: destination.Config{
				Config: config.Config{
					Endpoint: "test-endpoint",
					Scheme:   "https",
					Class:    "test-class",
					Auth: config.Auth{
						Mechanism: "apiKey",
						APIKey:    "xyz",
					},
				},
			},
		},
		{
			name: "WCS username and password",
			cfgMap: map[string]string{
				"endpoint":               "test-endpoint",
				"scheme":                 "https",
				"class":                  "test-class",
				"auth.mechanism":         "wcsCreds",
				"auth.wcsCreds.username": "abc",
				"auth.wcsCreds.password": "xyz",
			},
			wantCfg: destination.Config{
				Config: config.Config{
					Endpoint: "test-endpoint",
					Scheme:   "https",
					Class:    "test-class",
					Auth: config.Auth{
						Mechanism: "wcsCreds",
						WCSCredentials: config.WCSCredentials{
							Username: "abc",
							Password: "xyz",
						},
					},
				},
			},
		},
		{
			name: "partial WCS auth (username)",
			cfgMap: map[string]string{
				"endpoint":               "test-endpoint",
				"scheme":                 "https",
				"class":                  "test-class",
				"auth.mechanism":         "wcsCreds",
				"auth.wcsCreds.username": "abc",
			},
			wantErr: config.ErrUsernamePasswordMissing,
		},
		{
			name: "partial WCS auth (password)",
			cfgMap: map[string]string{
				"endpoint":               "test-endpoint",
				"scheme":                 "https",
				"class":                  "test-class",
				"auth.mechanism":         "wcsCreds",
				"auth.wcsCreds.password": "xyz",
			},
			wantErr: config.ErrUsernamePasswordMissing,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			is := is.New(t)
			ctx := context.Background()

			cfg := destination.Config{}
			err := sdk.Util.ParseConfig(ctx, tc.cfgMap, &cfg, weaviate.Connector.NewSpecification().DestinationParams)
			fmt.Println("--------------- ", weaviate.Connector.NewSpecification().DestinationParams)
			if tc.wantErr == nil {
				is.NoErr(err)
			} else {
				is.True(errors.Is(err, tc.wantErr))
			}

			err = cfg.Validate(ctx)
			if tc.wantErr == nil {
				is.NoErr(err)
			} else {
				is.True(errors.Is(err, tc.wantErr))
			}
		})
	}
}
