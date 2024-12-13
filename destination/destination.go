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

//go:generate paramgen -output=paramgen_dest.go Config
//go:generate mockgen -source=destination.go -package=mock -destination=mock/client_mock.go -mock_names=weaviateClient=WeaviateClient . weaviateClient

import (
	"context"
	"encoding/base64"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	sdkconfig "github.com/conduitio/conduit-commons/config"
	"github.com/conduitio/conduit-commons/opencdc"

	"github.com/conduitio-labs/conduit-connector-weaviate/destination/weaviate"
	sdk "github.com/conduitio/conduit-connector-sdk"
	"github.com/google/uuid"
)

var (
	metadataClass  = "weaviate.class"
	metadataVector = "weaviate.vector"
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
		sdk.DefaultDestinationMiddleware()...,
	)
}

func (d *Destination) Parameters() sdkconfig.Parameters {
	return d.config.Parameters()
}

func (d *Destination) Configure(ctx context.Context, cfg sdkconfig.Config) error {
	sdk.Logger(ctx).Info().Msg("Configuring Destination...")
	err := sdk.Util.ParseConfig(ctx, cfg, &d.config, New().Parameters())
	if err != nil {
		return fmt.Errorf("invalid config: %w", err)
	}

	if !d.config.ModuleHeader.IsValid() {
		return errors.New("invalid module configuration")
	}

	err = d.config.Validate()
	if err != nil {
		return fmt.Errorf("invalid configuration: %w", err)
	}

	return nil
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
	if record.Metadata != nil && record.Metadata[metadataClass] != "" {
		class = record.Metadata[metadataClass]
	}

	var vector []float32
	if record.Metadata != nil && record.Metadata[metadataVector] != "" {
		vector, err = parseEmbeddings(record.Metadata[metadataVector])
		if err != nil {
			return nil, fmt.Errorf("failed parsing vector from metadata, input: %v, error: %w", record.Metadata[metadataVector], err)
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

// ParseEmbeddings takes a string input and attempts to parse it into a slice of float32 embeddings.
// It supports two input formats:
// 1. Base64 encoded embeddings
// 2. CSV-formatted embeddings
func parseEmbeddings(input string) ([]float32, error) {
	input = strings.TrimSpace(input)

	if isBase64(input) {
		return decodeBase64Embeddings(input)
	}

	if isCSV(input) {
		return readCSVVector(input)
	}

	// If not base64 or CSV, return an error
	return nil, fmt.Errorf("unsupported input format: must be base64 or CSV")
}

func isBase64(input string) bool {
	_, err := base64.StdEncoding.DecodeString(input)
	return err == nil
}

func decodeBase64Embeddings(input string) ([]float32, error) {
	decodedBytes, err := base64.StdEncoding.DecodeString(input)
	if err != nil {
		return nil, fmt.Errorf("failed to decode base64: %v", err)
	}

	// Note: This assumes the base64 represents a slice of float32
	embeddings := make([]float32, len(decodedBytes)/4)
	for i := range embeddings {
		bits := uint32(decodedBytes[i*4]) |
			uint32(decodedBytes[i*4+1])<<8 |
			uint32(decodedBytes[i*4+2])<<16 |
			uint32(decodedBytes[i*4+3])<<24
		embeddings[i] = float32(bits)
	}

	return embeddings, nil
}

func isCSV(input string) bool {
	return strings.Contains(input, ",")
}

// readCSVVector parses a CSV string into a slice of float32
func readCSVVector(input string) ([]float32, error) {
	// Use csv.Reader to handle potential quotes and escaped characters
	reader := csv.NewReader(strings.NewReader(input))

	// Read the first (and assumed only) record
	record, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("error reading CSV: %v", err)
	}

	// Convert CSV fields to float32
	embeddings := make([]float32, len(record))
	for i, val := range record {
		// Trim any whitespace
		val = strings.TrimSpace(val)

		// Convert to float64 first to handle scientific notation
		floatVal, err := strconv.ParseFloat(val, 32)
		if err != nil {
			return nil, fmt.Errorf("error converting %s to float: %v", val, err)
		}

		embeddings[i] = float32(floatVal)
	}

	return embeddings, nil
}
