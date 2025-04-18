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
	WCSAuth  WCSAuth
	Endpoint string
	Scheme   string
	Headers  map[string]string
}

type WCSAuth struct {
	Username string
	Password string
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
	var authConfig auth.Config
	if config.APIKey != "" {
		authConfig = auth.ApiKey{Value: config.APIKey}
	} else {
		authConfig = auth.ResourceOwnerPasswordFlow{
			Username: config.WCSAuth.Username,
			Password: config.WCSAuth.Password,
		}
	}

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
