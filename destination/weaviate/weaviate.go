package weaviate

import (
	"context"
	"fmt"

	"github.com/weaviate/weaviate-go-client/v4/weaviate"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/auth"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/data/replication"
)

type Config struct {
	APIKey   string
	Endpoint string
	Scheme   string
	Headers  map[string]string
}

type Object struct {
	ID         string
	Class      string
	Properties map[string]interface{}
	Vector     []float32
}

type Client struct {
	client *weaviate.Client
}

func (c *Client) Open(config Config) error {
	authConfig := auth.ApiKey{Value: config.APIKey}

	wcfg := weaviate.Config{
		Host:       config.Endpoint,
		Scheme:     config.Scheme,
		AuthConfig: authConfig,
		Headers:    config.Headers,
	}

	client, err := weaviate.NewClient(wcfg)
	if err != nil {
		return fmt.Errorf("error creating client: %w", err)
	}

	c.client = client

	return nil
}

func (c *Client) Insert(ctx context.Context, obj *Object) error {
	//TODO: We should handle case where "vector" is in the payload.
	// you'd need to pull it out and add it on higher level __sL__
	// https://github.com/conduitio-labs/conduit-connector-weaviate/issues/3
	_, err := c.client.Data().Creator().
		WithClassName(obj.Class).
		WithID(obj.ID).
		WithProperties(obj.Properties).
		WithVector(obj.Vector).
		WithConsistencyLevel(replication.ConsistencyLevel.ALL).
		Do(ctx)

	if err != nil {
		return fmt.Errorf("error creating object: %w", err)
	}

	return nil
}

func (c *Client) Update(ctx context.Context, obj *Object) error {
	err := c.client.Data().Updater().
		WithID(obj.ID).
		WithClassName(obj.Class).
		WithProperties(obj.Properties).
		WithConsistencyLevel(replication.ConsistencyLevel.ALL).
		Do(ctx)

	if err != nil {
		return fmt.Errorf("error update object: %w", err)
	}

	return nil
}

func (c *Client) Delete(ctx context.Context, obj *Object) error {
	err := c.client.Data().Deleter().
		WithClassName(obj.Class).
		WithID(obj.ID).
		WithConsistencyLevel(replication.ConsistencyLevel.ALL).
		Do(ctx)

	if err != nil {
		return fmt.Errorf("error deleting object: %w", err)
	}

	return nil
}
