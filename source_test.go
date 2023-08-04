package weaviate_test

import (
	"context"
	"testing"

	weaviate "github.com/conduitio-labs/conduit-connector-weaviate"
	"github.com/matryer/is"
)

func TestTeardownSource_NoOpen(t *testing.T) {
	is := is.New(t)
	con := weaviate.NewSource()
	err := con.Teardown(context.Background())
	is.NoErr(err)
}
