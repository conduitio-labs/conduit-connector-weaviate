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

package destination_test

import (
	"context"
	"errors"
	"fmt"
	"os"
	"testing"
	"time"

	weaviateConn "github.com/conduitio-labs/conduit-connector-weaviate"
	"github.com/conduitio-labs/conduit-connector-weaviate/destination"
	"github.com/conduitio/conduit-commons/opencdc"
	sdk "github.com/conduitio/conduit-connector-sdk"
	"github.com/google/uuid"
	"github.com/matryer/is"
	"github.com/weaviate/weaviate-go-client/v4/weaviate"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/fault"
)

func TestDestination_Integration_Insert(t *testing.T) {
	integrationTest(t)

	is := is.New(t)
	openAIKey := os.Getenv("OPENAI_APIKEY")
	is.True(openAIKey != "") // expected OPENAI_APIKEY to be set

	testCases := []struct {
		name   string
		vector []float32
	}{
		{
			name: "no vector (auto-generated)",
		},
		{
			name:   "with vector",
			vector: []float32{123.45, 678.9},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			integrationTestInsert(t, openAIKey, tc.vector)
		})
	}
}

func integrationTestInsert(t *testing.T, openAIKey string, wantVector []float32) {
	is := is.New(t)
	ctx := context.Background()
	class := fmt.Sprintf("products_%v", time.Now().UnixMilli())
	cfg := integrationTestCfg(class, openAIKey)

	client, err := newWeaviateClient(cfg)
	is.NoErr(err)
	defer func() {
		err = client.Schema().
			ClassDeleter().
			WithClassName(class).
			Do(ctx)
		is.NoErr(err)
	}()

	underTest := destination.New()
	err = sdk.Util.ParseConfig(ctx, cfg, underTest.Config(), weaviateConn.Connector.NewSpecification().DestinationParams)
	is.NoErr(err)

	err = underTest.Open(ctx)
	is.NoErr(err)

	id := uuid.NewString()
	wantProperties := map[string]any{
		"product_name": "computer",
		"price":        220.15,
		"labels":       []any{"laptop", "navy-blue"},
		"used":         true,
	}
	rec := sdk.Util.Source.NewRecordCreate(
		opencdc.Position("test-position"),
		map[string]string{},
		opencdc.RawData(id),
		opencdc.StructuredData(wantProperties),
	)
	if wantVector != nil {
		rec.Metadata[destination.MetadataVector] = toString(wantVector)
	}

	n, err := underTest.Write(ctx, []opencdc.Record{rec})
	is.NoErr(err)
	is.Equal(1, n)

	wID := uuid.NewMD5(uuid.NameSpaceOID, []byte(id)).String()
	objects, err := client.Data().
		ObjectsGetter().
		WithClassName(cfg["class"]).
		WithID(wID).
		WithVector().
		Do(ctx)
	is.NoErr(err)
	is.Equal(1, len(objects))

	obj := objects[0]
	gotProperties, ok := obj.Properties.(map[string]any)
	is.True(ok) // expected object properties to be a map[string]any
	is.Equal(wantProperties, gotProperties)
	is.True(obj.Vector != nil)
	if wantVector != nil {
		is.Equal(wantVector, []float32(obj.Vector))
	}
}

func toString(floats []float32) string {
	s := ""
	for i, f := range floats {
		s += fmt.Sprint(f)
		if i < len(floats)-1 {
			s += ","
		}
	}

	return s
}

func TestDestination_Integration_Update(t *testing.T) {
	integrationTest(t)

	is := is.New(t)
	openAIKey := os.Getenv("OPENAI_APIKEY")
	is.True(openAIKey != "") // expected OPENAI_APIKEY to be set

	ctx := context.Background()
	class := fmt.Sprintf("Products_%v", time.Now().UnixMilli())
	cfg := integrationTestCfg(class, openAIKey)

	client, err := newWeaviateClient(cfg)
	is.NoErr(err)
	defer func() {
		err = client.Schema().
			ClassDeleter().
			WithClassName(class).
			Do(ctx)
		is.NoErr(err)
	}()

	underTest := destination.New()
	err = sdk.Util.ParseConfig(ctx, cfg, underTest.Config(), weaviateConn.Connector.NewSpecification().DestinationParams)
	is.NoErr(err)

	err = underTest.Open(ctx)
	is.NoErr(err)

	// Insert record
	id := "test-id"
	recInsert := sdk.Util.Source.NewRecordCreate(
		opencdc.Position("test-position"),
		map[string]string{},
		opencdc.RawData(id),
		opencdc.StructuredData(map[string]any{
			"product_name": "computer",
			"price":        220.15,
			"labels":       []any{"laptop", "navy-blue"},
			"used":         true,
		}),
	)

	n, err := underTest.Write(ctx, []opencdc.Record{recInsert})
	is.NoErr(err)
	is.Equal(1, n)

	// Update record
	payloadUpdate := map[string]any{
		"product_name": "computer",
		"price":        330.75,
		"labels":       []any{"laptop", "pink"},
		"used":         true,
	}
	recUpdate := sdk.Util.Source.NewRecordUpdate(
		opencdc.Position("test-position"),
		map[string]string{},
		opencdc.RawData(id),
		nil,
		opencdc.StructuredData(payloadUpdate),
	)

	n, err = underTest.Write(ctx, []opencdc.Record{recUpdate})
	is.NoErr(err)
	is.Equal(1, n)

	// Verify update
	wID := uuid.NewMD5(uuid.NameSpaceOID, []byte(id)).String()
	objects, err := client.Data().
		ObjectsGetter().
		WithClassName(cfg["class"]).
		WithID(wID).
		Do(ctx)
	is.NoErr(err)
	is.Equal(1, len(objects))

	obj := objects[0]
	got, ok := obj.Properties.(map[string]any)
	is.True(ok) // expected object properties to be a map[string]any
	is.Equal(payloadUpdate, got)
}

func TestDestination_Integration_Delete(t *testing.T) {
	integrationTest(t)

	is := is.New(t)
	openAIKey := os.Getenv("OPENAI_APIKEY")
	is.True(openAIKey != "") // expected OPENAI_APIKEY to be set

	ctx := context.Background()
	class := fmt.Sprintf("products_%v", time.Now().UnixMilli())
	cfg := integrationTestCfg(class, openAIKey)

	client, err := newWeaviateClient(cfg)
	is.NoErr(err)
	defer func() {
		err = client.Schema().
			ClassDeleter().
			WithClassName(class).
			Do(ctx)
		is.NoErr(err)
	}()

	underTest := destination.New()
	err = sdk.Util.ParseConfig(ctx, cfg, underTest.Config(), weaviateConn.Connector.NewSpecification().DestinationParams)
	is.NoErr(err)

	err = underTest.Open(ctx)
	is.NoErr(err)

	// Write record
	id := uuid.NewString()
	want := map[string]any{
		"product_name": "computer",
		"price":        220.15,
		"labels":       []any{"laptop", "navy-blue"},
		"used":         true,
	}
	rec := sdk.Util.Source.NewRecordCreate(
		opencdc.Position("test-position"),
		map[string]string{},
		opencdc.RawData(id),
		opencdc.StructuredData(want),
	)

	n, err := underTest.Write(ctx, []opencdc.Record{rec})
	is.NoErr(err)
	is.Equal(1, n)

	recDelete := sdk.Util.Source.NewRecordDelete(nil, nil, opencdc.RawData(id), nil)
	n, err = underTest.Write(ctx, []opencdc.Record{recDelete})
	is.NoErr(err)
	is.Equal(1, n)

	_, err = client.Data().
		ObjectsGetter().
		WithClassName(cfg["class"]).
		WithID(id).
		Do(ctx)
	is.True(err != nil)
	wErr := &fault.WeaviateClientError{}
	is.True(errors.As(err, &wErr))
	is.Equal(404, wErr.StatusCode)
}

func newWeaviateClient(cfg map[string]string) (*weaviate.Client, error) {
	wcfg := weaviate.Config{
		Host:   "localhost:18080",
		Scheme: "http",
		Headers: map[string]string{
			cfg["moduleHeader.name"]: cfg["moduleHeader.value"],
		},
	}

	return weaviate.NewClient(wcfg)
}

func integrationTest(t *testing.T) {
	t.Helper()
	if os.Getenv("RUN_INTEGRATION_TESTS") == "" {
		t.Skip("skipping integration tests, set environment variable RUN_INTEGRATION_TESTS")
	}
}

func integrationTestCfg(class, openAIKey string) map[string]string {
	return map[string]string{
		"endpoint":           "localhost:18080",
		"scheme":             "http",
		"auth.mechanism":     "none",
		"class":              class,
		"moduleHeader.name":  "X-OpenAI-Api-Key",
		"moduleHeader.value": openAIKey,
		"generateUUID":       "true",
	}
}
