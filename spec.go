package weaviate

import (
	sdk "github.com/conduitio/conduit-connector-sdk"
)

// version is set during the build process with ldflags (see Makefile).
// Default version matches default from runtime/debug.
var version = "(devel)"

// Specification returns the connector's specification.
func Specification() sdk.Specification {
	return sdk.Specification{
		Name:        "weaviate",
		Summary:     "Conduit Connector for Weaviate",
		Description: "Conduit Connector for Weaviate",
		Version:     version,
		Author:      "Meroxa, Inc.",
	}
}
