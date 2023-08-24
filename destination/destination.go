package destination

//go:generate paramgen -output=paramgen_dest.go DestinationConfig

import (
	"context"
	"errors"
	"fmt"
	"github.com/conduitio-labs/conduit-connector-weaviate/destination/handler"

	"github.com/conduitio-labs/conduit-connector-weaviate/config"
	sdk "github.com/conduitio/conduit-connector-sdk"
)

type recordHandler interface {
	Open(DestinationConfig) error

	Insert(context.Context, sdk.Record) error
	Update(context.Context, sdk.Record) error
	Delete(context.Context, sdk.Record) error
}

type Destination struct {
	sdk.UnimplementedDestination

	config  DestinationConfig
	handler recordHandler
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
	ModuleAPIKey ModuleApiKey `json:"module_api_key"`
	GenerateUUID bool         `json:"generate_uuid"`
}

func New() sdk.Destination {
	return NewWithHandler(&handler.RecordHandler{})
}

func NewWithHandler(h recordHandler) sdk.Destination {
	return sdk.DestinationWithMiddleware(
		&Destination{handler: h},
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
	err := d.handler.Open(d.config)
	if err != nil {
		return fmt.Errorf("error creating handler: %w}", err)
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
			d.handler.Insert,
			d.handler.Update,
			d.handler.Delete,
			d.handler.Insert,
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
