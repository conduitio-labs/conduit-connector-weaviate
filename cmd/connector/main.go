package main

import (
	sdk "github.com/conduitio/conduit-connector-sdk"

	weaviate "github.com/conduitio-labs/conduit-connector-weaviate"
)

func main() {
	sdk.Serve(weaviate.Connector)
}
