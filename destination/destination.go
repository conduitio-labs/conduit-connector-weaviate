package destination

//go:generate paramgen -output=paramgen_dest.go DestinationConfig
//go:generate mockgen -source=destination.go -package=mock -destination=mock/client_mock.go -mock_names=weaviateClient=WeaviateClient . weaviateClient

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/conduitio-labs/conduit-connector-weaviate/config"
	"github.com/conduitio-labs/conduit-connector-weaviate/destination/weaviate"
	sdk "github.com/conduitio/conduit-connector-sdk"
	"github.com/google/uuid"
)

var (
	metadataClass = "weaviate.class"
)

type weaviateClient interface {
	Open(weaviate.Config) error

	Insert(context.Context, *weaviate.Object) error
	Update(context.Context, *weaviate.Object) error
	Delete(context.Context, *weaviate.Object) error
}

type Destination struct {
	sdk.UnimplementedDestination

	config DestinationConfig
	client weaviateClient
}

type ModuleApiKey struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

func (m ModuleApiKey) IsValid() bool {
	return (m.Name == "" && m.Value == "") ||
		(m.Name != "" && m.Value != "")
}

type DestinationConfig struct {
	config.Config
	//TODO: better naming for this value __sL__
	ModuleAPIKey ModuleApiKey `json:"moduleAPIKey"`
	GenerateUUID bool         `json:"generateUUID"`
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

func (d *Destination) Parameters() map[string]sdk.Parameter {
	return d.config.Parameters()
}

func (d *Destination) Configure(ctx context.Context, cfg map[string]string) error {
	sdk.Logger(ctx).Info().Msg("Configuring Destination...")
	err := sdk.Util.ParseConfig(cfg, &d.config)
	if err != nil {
		return fmt.Errorf("invalid config: %w", err)
	}

	if !d.config.ModuleAPIKey.IsValid() {
		return errors.New("invalid module configuration")
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

func (d *Destination) Write(ctx context.Context, records []sdk.Record) (int, error) {
	//TODO: will need differential handling of insert/update/delete __sL__
	// weaviate has id field that is required to be UUID, if not provided it will
	// generate one itself. Issue here is how to handle update/delete if we don't know
	// the id.

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

func (d *Destination) insert(ctx context.Context, record sdk.Record) error {
	obj, err := d.toWeaviateObj(record)
	if err != nil {
		return fmt.Errorf("error creating Weaviate object: %w", err)
	}

	return d.client.Insert(ctx, obj)
}

func (d *Destination) update(ctx context.Context, record sdk.Record) error {
	obj, err := d.toWeaviateObj(record)
	if err != nil {
		return fmt.Errorf("error creating Weaviate object: %w", err)
	}

	return d.client.Update(ctx, obj)
}

func (d *Destination) delete(ctx context.Context, record sdk.Record) error {
	return d.client.Delete(
		ctx,
		&weaviate.Object{
			ID:    d.recordUUID(record),
			Class: d.config.Class,
		},
	)
}

func (d *Destination) toWeaviateObj(record sdk.Record) (*weaviate.Object, error) {
	properties, err := d.recordProperties(record)

	if err != nil {
		return nil, fmt.Errorf("update property conversion: %w", err)
	}

	class := d.config.Class
	if record.Metadata != nil && record.Metadata[metadataClass] != "" {
		class = record.Metadata[metadataClass]
	}

	return &weaviate.Object{
		ID:         d.recordUUID(record),
		Class:      class,
		Properties: properties,
	}, nil
}

func (d *Destination) recordUUID(record sdk.Record) string {
	key := record.Key.Bytes()
	if !d.config.GenerateUUID {
		return string(key)
	}
	return uuid.NewMD5(uuid.NameSpaceOID, key).String()
}

func (d *Destination) recordProperties(record sdk.Record) (sdk.StructuredData, error) {
	data := record.Payload.After

	if data == nil || len(data.Bytes()) == 0 {
		return nil, errors.New("empty payload")
	}

	properties := make(sdk.StructuredData)
	err := json.Unmarshal(data.Bytes(), &properties)

	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal payload to structured data: %w", err)
	}

	return properties, nil
}

func (d *Destination) weaviateConfig() weaviate.Config {
	var headers map[string]string
	if d.config.ModuleAPIKey.IsValid() {
		headers = map[string]string{
			d.config.ModuleAPIKey.Name: d.config.ModuleAPIKey.Value,
		}
	}

	return weaviate.Config{
		APIKey:   d.config.APIKey,
		Endpoint: d.config.Endpoint,
		Scheme:   d.config.Scheme,
		Headers:  headers,
	}
}
