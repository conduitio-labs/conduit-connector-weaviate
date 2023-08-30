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
	"errors"
	"testing"

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
				"auth.mechanism": "none",
			},
			wantCfg: destination.Config{
				Config: config.Config{
					Auth: config.Auth{
						Mechanism: "none",
					},
				},
			},
		},
		{
			name: "API key only",
			cfgMap: map[string]string{
				"auth.mechanism": "apiKey",
				"auth.apiKey":    "xyz",
			},
			wantCfg: destination.Config{
				Config: config.Config{
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
				"auth.mechanism":         "wcsCreds",
				"auth.wcsCreds.username": "abc",
				"auth.wcsCreds.password": "xyz",
			},
			wantCfg: destination.Config{
				Config: config.Config{
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
				"auth.mechanism":         "wcsCreds",
				"auth.wcsCreds.username": "abc",
			},
			wantErr: config.ErrUsernamePasswordMissing,
		},
		{
			name: "partial WCS auth (password)",
			cfgMap: map[string]string{
				"auth.mechanism":         "wcsCreds",
				"auth.wcsCreds.password": "xyz",
			},
			wantErr: config.ErrUsernamePasswordMissing,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			is := is.New(t)

			cfg := destination.Config{}
			err := sdk.Util.ParseConfig(tc.cfgMap, &cfg)
			is.NoErr(err)

			err = cfg.Validate()
			if tc.wantErr == nil {
				is.NoErr(err)
				is.Equal(tc.wantCfg, cfg)
			} else {
				is.True(errors.Is(err, tc.wantErr))
			}
		})
	}
}
