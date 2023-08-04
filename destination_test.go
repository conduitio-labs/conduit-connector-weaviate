package weaviate_test

import (
	"context"
	"testing"

	weaviate "github.com/conduitio-labs/conduit-connector-weaviate"
	"github.com/matryer/is"
)

func TestTeardown_NoOpen(t *testing.T) {
	is := is.New(t)
	con := weaviate.NewDestination()
	err := con.Teardown(context.Background())
	is.NoErr(err)
}
