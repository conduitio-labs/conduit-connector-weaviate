package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	sdk "github.com/conduitio/conduit-connector-sdk"
	"github.com/google/uuid"
	"github.com/weaviate/weaviate-go-client/v4/weaviate"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/data/replication"
)

type RecordHandler struct {
	client       *weaviate.Client
	class        string
	generateUUID bool
}

func New(client *weaviate.Client, class string, genUUID bool) (*RecordHandler, error) {
	handler := &RecordHandler{
		client:       client,
		class:        class,
		generateUUID: genUUID,
	}

	return handler, nil
}

func recordUUID(record sdk.Record) string {
	key := record.Key.Bytes()
	return uuid.NewMD5(uuid.NameSpaceOID, key).String()
}

func recordProperties(record sdk.Record) (sdk.StructuredData, error) {
	data := record.Payload.After

	if data == nil || len(data.Bytes()) == 0 {
		return nil, errors.New("Empty payload")
	}

	properties := make(sdk.StructuredData)
	err := json.Unmarshal(data.Bytes(), &properties)

	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal data to structured data: %w", err)
	}

	return properties, nil
}

func (h *RecordHandler) Insert(ctx context.Context, record sdk.Record) error {

	properties, err := recordProperties(record)

	if err != nil {
		return fmt.Errorf("insert property coversion: %w", err)
	}

	id := string(record.Key.Bytes())
	if h.generateUUID {
		id = recordUUID(record)
	}

	_, err = h.client.Data().Creator().
		WithClassName(h.class).
		WithID(id).
		WithProperties(properties).
		WithConsistencyLevel(replication.ConsistencyLevel.ALL).
		Do(ctx)

	if err != nil {
		return fmt.Errorf("error creating object: %w", err)
	}

	return nil
}

func (h *RecordHandler) Update(ctx context.Context, record sdk.Record) error {

	properties, err := recordProperties(record)

	if err != nil {
		return fmt.Errorf("update property coversion: %w", err)
	}

	id := string(record.Key.Bytes())
	if h.generateUUID {
		id = recordUUID(record)
	}

	err = h.client.Data().Updater().
		WithID(id).
		WithClassName(h.class).
		WithProperties(properties).
		WithConsistencyLevel(replication.ConsistencyLevel.ALL).
		Do(ctx)

	if err != nil {
		return fmt.Errorf("error update object: %w", err)
	}

	return nil
}

func (h *RecordHandler) Delete(ctx context.Context, record sdk.Record) error {

	id := string(record.Key.Bytes())
	if h.generateUUID {
		id = recordUUID(record)
	}

	err := h.client.Data().Deleter().
		WithClassName(h.class).
		WithID(id).
		WithConsistencyLevel(replication.ConsistencyLevel.ALL).
		Do(ctx)

	if err != nil {
		return fmt.Errorf("error deleting object: %w", err)
	}

	return nil
}
