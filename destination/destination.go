package destination

//go:generate paramgen -output=paramgen_dest.go DestinationConfig

import (
	"context"
	"fmt"

	"github.com/conduitio-labs/conduit-connector-weaviate/config"
	sdk "github.com/conduitio/conduit-connector-sdk"
	"github.com/weaviate/weaviate-go-client/v4/weaviate"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/auth"
)

type Destination struct {
	sdk.UnimplementedDestination

	config DestinationConfig
	client *weaviate.Client
}

type ModuleApiKey struct {
	name  string
	value string
}

type DestinationConfig struct {
	config.Config
	moduleApiKey ModuleApiKey `json:"module_api_key"`
	generateUUID string       `json:"generate_uuid"`
}

func New() sdk.Destination {
	// Create Destination and wrap it in the default middleware.
	return sdk.DestinationWithMiddleware(&Destination{}, sdk.DefaultDestinationMiddleware()...)
}

func (d *Destination) Parameters() map[string]sdk.Parameter {
	return d.config.Parameters()
}

func (d *Destination) Configure(ctx context.Context, cfg map[string]string) error {
	// Configure is the first function to be called in a connector. It provides
	// the connector with the configuration that can be validated and stored.
	// In case the configuration is not valid it should return an error.
	// Testing if your connector can reach the configured data source should be
	// done in Open, not in Configure.
	// The SDK will validate the configuration and populate default values
	// before calling Configure. If you need to do more complex validations you
	// can do them manually here.

	var authConfig auth.Config
	var clientHeaders map[string]string

	sdk.Logger(ctx).Info().Msg("Configuring Destination...")
	err := sdk.Util.ParseConfig(cfg, &d.config)
	if err != nil {
		return fmt.Errorf("invalid config: %w", err)
	}

	//TODO: support additional auth schemes __sL__
	if d.config.ApiKey != "" {
		authConfig = auth.ApiKey{Value: d.config.ApiKey}
	}

	//TODO: better naming for this value __sL__
	if d.config.moduleApiKey.name != "" && d.config.moduleApiKey.value != "" {
		clientHeaders = map[string]string{
			d.config.moduleApiKey.name: d.config.moduleApiKey.value,
		}
	}

	wcfg := weaviate.Config{
		Host:       d.config.Endpoint,
		Scheme:     d.config.Scheme,
		AuthConfig: authConfig,
		Headers:    clientHeaders,
	}

	//TODO: need to look into this is actually creating connection and thus should be in open func __sL__
	d.client, err = weaviate.NewClient(wcfg)
	if err != nil {
		return fmt.Errorf("Error creating client: %w", err)
	}

	return nil
}

func (d *Destination) Open(ctx context.Context) error {
	return nil
}

func insertRecord(ctx context.Context, record sdk.Record) error {
	return nil
}

func updateRecord(ctx context.Context, record sdk.Record) error {
	return nil
}

func deleteRecord(ctx context.Context, record sdk.Record) error {
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
			insertRecord,
			updateRecord,
			deleteRecord,
			insertRecord,
		)

		if err != nil {
			return i, fmt.Errorf("Error routing %s: %w", record.Operation.String(), err)
		}
	}

	//objects := make([]*models.Object, len(records))
	//for i := range records {
	//	objects[i] = &models.Object{
	//		Class: d.Config.Class,
	//		Properties: XXX
	//	}
	//}

	//result, err := client.Batch().ObjectsBatcher().WithObjects(objects...).Do(context.Background())
	//if err != nil {
	//	return 0, nil
	//}

	return len(records), nil
}

func (d *Destination) Teardown(ctx context.Context) error {
	// Teardown signals to the plugin that all records were written and there
	// will be no more calls to any other function. After Teardown returns, the
	// plugin should be ready for a graceful shutdown.
	return nil
}
