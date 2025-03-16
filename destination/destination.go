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

//go:generate mockgen -source=destination.go -package=mock -destination=mock/client_mock.go -mock_names=weaviateClient=WeaviateClient . weaviateClient

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/conduitio/conduit-commons/opencdc"

	"github.com/conduitio-labs/conduit-connector-weaviate/destination/weaviate"
	sdk "github.com/conduitio/conduit-connector-sdk"
	"github.com/google/uuid"
)

var (
	MetadataClass  = "weaviate.class"
	MetadataVector = "weaviate.vector"
)

type weaviateClient interface {
	Open(weaviate.Config) error

	Insert(context.Context, *weaviate.Object) error
	Update(context.Context, *weaviate.Object) error
	Delete(context.Context, *weaviate.Object) error
}

type Destination struct {
	sdk.UnimplementedDestination

	config Config
	client weaviateClient
}

func New() sdk.Destination {
	return NewWithClient(&weaviate.Client{})
}

func NewWithClient(client weaviateClient) sdk.Destination {
	return sdk.DestinationWithMiddleware(
		&Destination{client: client},
	)
}

func (d *Destination) Config() sdk.DestinationConfig {
	return &d.config
}

func (d *Destination) Open(context.Context) error {
	err := d.client.Open(d.weaviateConfig())
	if err != nil {
		return fmt.Errorf("error creating client: %w}", err)
	}

	return nil
}

func (d *Destination) Write(ctx context.Context, records []opencdc.Record) (int, error) {
	for i, record := range records {
		err := sdk.Util.Destination.Route(
			ctx,
			record,
			d.insert,
			d.update,
			d.delete,
			d.insert,
		)
		if err != nil {
			return i, fmt.Errorf("error routing %v: %w", record.Operation, err)
		}
	}

	return len(records), nil
}

func (d *Destination) Teardown(context.Context) error {
	// Teardown signals to the plugin that all records were written and there
	// will be no more calls to any other function. After Teardown returns, the
	// plugin should be ready for a graceful shutdown.
	return nil
}

func (d *Destination) insert(ctx context.Context, record opencdc.Record) error {
	obj, err := d.toWeaviateObj(record)
	if err != nil {
		return fmt.Errorf("error creating Weaviate object: %w", err)
	}

	return d.client.Insert(ctx, obj)
}

func (d *Destination) update(ctx context.Context, record opencdc.Record) error {
	obj, err := d.toWeaviateObj(record)
	if err != nil {
		return fmt.Errorf("error creating Weaviate object: %w", err)
	}

	return d.client.Update(ctx, obj)
}

func (d *Destination) delete(ctx context.Context, record opencdc.Record) error {
	return d.client.Delete(
		ctx,
		&weaviate.Object{
			ID:    d.recordUUID(record),
			Class: d.config.Class,
		},
	)
}

func (d *Destination) toWeaviateObj(record opencdc.Record) (*weaviate.Object, error) {
	properties, err := d.recordProperties(record)
	if err != nil {
		return nil, fmt.Errorf("update property conversion: %w", err)
	}

	class := d.config.Class
	if record.Metadata != nil && record.Metadata[MetadataClass] != "" {
		class = record.Metadata[MetadataClass]
	}

	var vector []float32
	if record.Metadata != nil && record.Metadata[MetadataVector] != "" {
		vector, err = d.recordVector(record.Metadata[MetadataVector])
		if err != nil {
			return nil, fmt.Errorf("failed parsing vector from metadata, input: %v, error: %w", record.Metadata[MetadataVector], err)
		}
	}

	return &weaviate.Object{
		ID:         d.recordUUID(record),
		Class:      class,
		Properties: properties,
		Vector:     vector,
	}, nil
}

func (d *Destination) recordUUID(record opencdc.Record) string {
	key := record.Key.Bytes()
	if !d.config.GenerateUUID {
		return string(key)
	}
	return uuid.NewMD5(uuid.NameSpaceOID, key).String()
}

func (d *Destination) recordProperties(record opencdc.Record) (map[string]interface{}, error) {
	data := record.Payload.After

	if data == nil || len(data.Bytes()) == 0 {
		return nil, errors.New("empty payload")
	}

	properties := make(map[string]interface{})
	err := json.Unmarshal(data.Bytes(), &properties)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal payload to structured data: %w", err)
	}

	return properties, nil
}

func (d *Destination) weaviateConfig() weaviate.Config {
	cfg := weaviate.Config{
		Endpoint: d.config.Endpoint,
		Scheme:   d.config.Scheme,
	}

	if d.config.ModuleHeader.IsValid() {
		cfg.Headers = map[string]string{
			d.config.ModuleHeader.Name: d.config.ModuleHeader.Value,
		}
	}

	if d.config.Auth.APIKey != "" {
		cfg.APIKey = d.config.Auth.APIKey
	} else {
		cfg.WCSAuth = weaviate.WCSAuth{
			Username: d.config.Auth.WCSCredentials.Username,
			Password: d.config.Auth.WCSCredentials.Password,
		}
	}

	return cfg
}

func (d *Destination) recordVector(s string) ([]float32, error) {
	var vector []float32
	for _, vs := range strings.Split(s, ",") {
		if vs == "" {
			return nil, errors.New("got an empty string")
		}

		v, err := strconv.ParseFloat(vs, 32)
		if err != nil {
			return nil, fmt.Errorf("cannot parse %v: %w", vs, err)
		}

		vector = append(vector, float32(v))
	}

	return vector, nil
}
