package handler

import (
	"context"

	sdk "github.com/conduitio/conduit-connector-sdk"
	"github.com/weaviate/weaviate-go-client/v4/weaviate"
)

type RecordHandler struct {
	generateUUID bool
	client       *weaviate.Client
}

func New(client *weaviate.Client, genUUID bool) (*RecordHandler, error) {
	handler := &RecordHandler{
		generateUUID: genUUID,
		client:       client,
	}

	return handler, nil
}

func (h *RecordHandler) Insert(ctx context.Context, record sdk.Record) error {
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
	return nil
}

func (h *RecordHandler) Update(ctx context.Context, record sdk.Record) error {
	return nil
}

func (h *RecordHandler) Delete(ctx context.Context, record sdk.Record) error {
	return nil
}
