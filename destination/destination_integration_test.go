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

package destination

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	sdk "github.com/conduitio/conduit-connector-sdk"
	"github.com/matryer/is"
	"github.com/weaviate/weaviate-go-client/v4/weaviate"
)

func TestDestination_Integration_Insert(t *testing.T) {
	is := is.New(t)
	openAIKey := os.Getenv("OPENAI_APIKEY")
	is.True(openAIKey != "") // expected OPENAI_APIKEY to be set

	ctx := context.Background()
	class := fmt.Sprintf("products_%v", time.Now().UnixMilli())
	cfg := map[string]string{
		"endpoint":           "localhost:18080",
		"scheme":             "http",
		"class":              class,
		"moduleAPIKey.name":  "X-OpenAI-Api-Key",
		"moduleAPIKey.value": openAIKey,
		"generateUUID":       "false",
	}

	client, err := newWeaviateClient(cfg)
	is.NoErr(err)
	defer func() {
		err = client.Schema().
			ClassDeleter().
			WithClassName(class).
			Do(ctx)
		is.NoErr(err)
	}()

	underTest := New()
	err = underTest.Configure(ctx, cfg)
	is.NoErr(err)

	err = underTest.Open(ctx)
	is.NoErr(err)

	id := "f9a510b3-5865-40e4-9fe8-e7fbab25b8bc"
	payload := map[string]any{
		"product_name": "computer",
		"price":        float64(1000),
		"labels":       []string{"laptop", "navy-blue"},
		"used":         true,
	}
	rec := sdk.Util.Source.NewRecordCreate(
		sdk.Position("test-position"),
		map[string]string{},
		sdk.RawData(id),
		sdk.StructuredData(payload),
	)

	n, err := underTest.Write(ctx, []sdk.Record{rec})
	is.NoErr(err)
	is.Equal(1, n)

	objects, err := client.Data().
		ObjectsGetter().
		WithClassName(cfg["class"]).
		WithID(id).
		Do(ctx)
	is.NoErr(err)
	is.Equal(1, len(objects))

	obj := objects[0]
	objProps, ok := obj.Properties.(map[string]any)
	is.True(ok) // expected object properties to be a map[string]any
	is.Equal(payload, objProps)
}

func newWeaviateClient(cfg map[string]string) (*weaviate.Client, error) {
	wcfg := weaviate.Config{
		Host:   "localhost:18080",
		Scheme: "http",
		Headers: map[string]string{
			cfg["moduleAPIKey.name"]: cfg["moduleAPIKey.value"],
		},
	}

	return weaviate.NewClient(wcfg)
}
