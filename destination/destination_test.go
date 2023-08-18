package destination_test

import (
	"context"
	"testing"

	"github.com/conduitio-labs/conduit-connector-weaviate/destination"
	"github.com/matryer/is"
)

func TestTeardown_NoOpen(t *testing.T) {
	is := is.New(t)
	con := destination.New()
	err := con.Teardown(context.Background())
	is.NoErr(err)
}
