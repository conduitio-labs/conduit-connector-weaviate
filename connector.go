package weaviate

import (
	sdk "github.com/conduitio/conduit-connector-sdk"

	"github.com/conduitio-labs/conduit-connector-weaviate/destination"
)

// Connector combines all constructors for each plugin in one struct.
var Connector = sdk.Connector{
	NewSpecification: Specification,
	NewSource:        nil,
	NewDestination:   destination.New,
}
