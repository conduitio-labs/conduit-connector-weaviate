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

package destination

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/conduitio-labs/conduit-connector-weaviate/destination/mock"
	"github.com/conduitio-labs/conduit-connector-weaviate/destination/weaviate"
	sdk "github.com/conduitio/conduit-connector-sdk"
	"github.com/google/go-cmp/cmp"
	"github.com/matryer/is"
	"go.uber.org/mock/gomock"
)

type eqMatcher struct {
	want interface{}
}

func newEqMatcher(want interface{}) eqMatcher {
	return eqMatcher{
		want: want,
	}
}

func (eq eqMatcher) Matches(got interface{}) bool {
	return gomock.Eq(eq.want).Matches(got)
}

func (eq eqMatcher) Got(got interface{}) string {
	diff := cmp.Diff(
		got,
		eq.want,
	)
	return fmt.Sprintf(
		"%v (%T)\nDiff (-got +want):\n%s",
		got,
		got,
		strings.TrimSpace(diff),
	)
}

func (eq eqMatcher) String() string {
	return fmt.Sprintf("%v (%T)\n", eq.want, eq.want)
}

func TestDestination_Teardown_NoOpen(t *testing.T) {
	is := is.New(t)
	underTest := New()
	err := underTest.Teardown(context.Background())
	is.NoErr(err)
}

func TestDestination_Open_WCSAuth(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	cfg := map[string]string{
		"endpoint":           "test-endpoint",
		"scheme":             "test-scheme",
		"auth.mechanism":     "wcsCreds",
		"auth.wcs.username":  "conduit-user",
		"auth.wcs.password":  "secret",
		"class":              "test-class",
		"moduleHeader.name":  "X-OpenAI-Api-Key",
		"moduleHeader.value": "test-OpenAI-Api-Key",
		"generateUUID":       "true",
	}

	ctrl := gomock.NewController(t)
	client := mock.NewWeaviateClient(ctrl)
	client.EXPECT().
		Open(gomock.Eq(weaviate.Config{
			WCSAuth: weaviate.WCSAuth{
				Username: cfg["auth.wcs.username"],
				Password: cfg["auth.wcs.password"],
			},
			Endpoint: cfg["endpoint"],
			Scheme:   cfg["scheme"],
			Headers: map[string]string{
				"X-OpenAI-Api-Key": "test-OpenAI-Api-Key",
			},
		}))

	underTest := NewWithClient(client)
	err := underTest.Configure(ctx, cfg)
	is.NoErr(err)

	err = underTest.Open(ctx)
	is.NoErr(err)
}

func TestDestination_Open_OpensClient(t *testing.T) {
	ctx := context.Background()
	cfg := map[string]string{
		"endpoint":           "test-endpoint",
		"scheme":             "test-scheme",
		"auth.mechanism":     "apiKey",
		"auth.apiKey":        "test-api-key",
		"class":              "test-class",
		"moduleHeader.name":  "X-OpenAI-Api-Key",
		"moduleHeader.value": "test-OpenAI-Api-Key",
		"generateUUID":       "true",
	}

	// setupTest is doing the basic checks
	_, _ = setupTest(t, ctx, cfg)
}

func TestDestination_SingleWrite(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	cfg := map[string]string{
		"endpoint":           "test-endpoint",
		"scheme":             "test-scheme",
		"auth.mechanism":     "apiKey",
		"auth.apiKey":        "test-api-key",
		"class":              "test-class",
		"moduleHeader.name":  "X-OpenAI-Api-Key",
		"moduleHeader.value": "test-OpenAI-Api-Key",
		"generateUUID":       "false",
	}

	testCases := []struct {
		name   string
		record sdk.Record
		want   *weaviate.Object
	}{
		{
			name: "raw payload, raw key, use class from config",
			record: sdk.Util.Source.NewRecordCreate(
				sdk.Position("test-position"),
				map[string]string{},
				sdk.RawData("f9a510b3-5865-40e4-9fe8-e7fbab25b8bc"),
				sdk.RawData(`{
					"product_name": "computer",
					"price":        1000,
					"labels":       ["laptop", "navy-blue"],
					"used":			true
				}`),
			),
			want: &weaviate.Object{
				ID:    "f9a510b3-5865-40e4-9fe8-e7fbab25b8bc",
				Class: cfg["class"],
				Properties: map[string]interface{}{
					"product_name": "computer",
					"price":        float64(1000),
					"labels":       []any{"laptop", "navy-blue"},
					"used":         true,
				},
			},
		},
		{
			name: "structured payload, raw key, use class from config",
			record: sdk.Util.Source.NewRecordCreate(
				sdk.Position("test-position"),
				map[string]string{},
				sdk.RawData("f9a510b3-5865-40e4-9fe8-e7fbab25b8bc"),
				sdk.StructuredData{
					"product_name": "computer",
					"price":        1000,
					"labels":       []string{"laptop", "navy-blue"},
					"used":         true,
				},
			),
			want: &weaviate.Object{
				ID:    "f9a510b3-5865-40e4-9fe8-e7fbab25b8bc",
				Class: cfg["class"],
				Properties: map[string]interface{}{
					"product_name": "computer",
					"price":        float64(1000),
					"labels":       []any{"laptop", "navy-blue"},
					"used":         true,
				},
			},
		},
		{
			name: "structured payload, raw key, use class from metadata",
			record: sdk.Util.Source.NewRecordCreate(
				sdk.Position("test-position"),
				map[string]string{
					metadataClass: "top-secret-class",
				},
				sdk.RawData("f9a510b3-5865-40e4-9fe8-e7fbab25b8bc"),
				sdk.StructuredData{
					"product_name": "computer",
					"price":        1000,
					"labels":       []string{"laptop", "navy-blue"},
					"used":         true,
				},
			),
			want: &weaviate.Object{
				ID:    "f9a510b3-5865-40e4-9fe8-e7fbab25b8bc",
				Class: "top-secret-class",
				Properties: map[string]interface{}{
					"product_name": "computer",
					"price":        float64(1000),
					"labels":       []any{"laptop", "navy-blue"},
					"used":         true,
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			underTest, wClient := setupTest(t, ctx, cfg)
			wClient.EXPECT().Insert(ctx, newEqMatcher(tc.want))

			n, err := underTest.Write(ctx, []sdk.Record{tc.record})
			is.NoErr(err)
			is.Equal(1, n)
		})
	}
}

func TestDestination_RecordWithVector(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	cfg := map[string]string{
		"endpoint":           "test-endpoint",
		"scheme":             "test-scheme",
		"class":              "test-class",
		"auth.mechanism":     "apiKey",
		"auth.apiKey":        "test-api-key",
		"moduleHeader.name":  "X-OpenAI-Api-Key",
		"moduleHeader.value": "test-OpenAI-Api-Key",
		"generateUUID":       "false",
	}

	testCases := []struct {
		name    string
		input   string
		want    []float32
		wantErr error
	}{
		{
			name:  "valid vector",
			input: "110.1,220",
			want:  []float32{110.1, 220},
		},
		{
			name:  "no vector",
			input: "",
			want:  nil,
		},
		{
			name:  "empty element",
			input: "111.1,   ,222.2",
			want:  nil,
			wantErr: errors.New(
				"error routing create: error creating Weaviate object: " +
					"failed parsing vector from metadata, input: 111.1,   ,222.2, error: cannot parse    : " +
					"strconv.ParseFloat: parsing \"   \": invalid syntax",
			),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			underTest, wClient := setupTest(t, ctx, cfg)
			inputRec := sdk.Util.Source.NewRecordCreate(
				sdk.Position("test-position"),
				map[string]string{
					metadataClass:  "top-secret-class",
					metadataVector: tc.input,
				},
				sdk.RawData("f9a510b3-5865-40e4-9fe8-e7fbab25b8bc"),
				sdk.StructuredData{
					"product_name": "computer",
					"price":        1000,
					"labels":       []string{"laptop", "navy-blue"},
					"used":         true,
				},
			)
			wantObj := &weaviate.Object{
				ID:    "f9a510b3-5865-40e4-9fe8-e7fbab25b8bc",
				Class: "top-secret-class",
				Properties: map[string]interface{}{
					"product_name": "computer",
					"price":        float64(1000),
					"labels":       []any{"laptop", "navy-blue"},
					"used":         true,
				},
				Vector: tc.want,
			}

			if tc.wantErr == nil {
				wClient.EXPECT().Insert(ctx, newEqMatcher(wantObj))
			}

			n, err := underTest.Write(ctx, []sdk.Record{inputRec})

			if tc.wantErr == nil {
				is.NoErr(err)
				is.Equal(1, n)
			} else {
				is.True(err != nil)
				is.Equal(tc.wantErr.Error(), err.Error())
				is.Equal(0, n)
			}
		})
	}
}

func setupTest(t *testing.T, ctx context.Context, cfg map[string]string) (sdk.Destination, *mock.WeaviateClient) {
	is := is.New(t)
	ctrl := gomock.NewController(t)
	client := mock.NewWeaviateClient(ctrl)
	client.EXPECT().
		Open(gomock.Eq(weaviate.Config{
			APIKey:   cfg["auth.apiKey"],
			Endpoint: cfg["endpoint"],
			Scheme:   cfg["scheme"],
			Headers: map[string]string{
				"X-OpenAI-Api-Key": "test-OpenAI-Api-Key",
			},
		}))

	underTest := NewWithClient(client)
	err := underTest.Configure(ctx, cfg)
	is.NoErr(err)

	err = underTest.Open(ctx)
	is.NoErr(err)

	return underTest, client
}
