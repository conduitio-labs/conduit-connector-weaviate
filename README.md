# Conduit Connector for <!-- readmegen:name -->Weaviate<!-- /readmegen:name -->

[Conduit](https://conduit.io) connector for <!-- readmegen:name -->Weaviate<!-- /readmegen:name -->.

<!-- readmegen:description -->
A Conduit destination connector for Weaviate, written in Go<!-- /readmegen:description -->

## How to build?
Run `make build` to build the connector. For instructions on how to use built connector
with Conduit, check the [Conduit documentation](https://conduit.io/docs/connectors/installing).

## Testing
To run the unit test, execute `make test`. 

The integration tests require a running Weaviate instance. The provided `make` targets run
a Weaviate instance using Docker. To run the integration tests, you need to:
1. set the `OPENAI_APIKEY` environment variable, which stores an OpenAI API key.
2. execute `make test-integration`.

The Docker compose file at `test/docker-compose.yml` can be used to run an instance of Weaviate locally.

## Destination

The Weaviate destination connectors handles all the changes supported by Conduit, 
which are: inserts, updates, and deletes. 

### Configuration

<!-- readmegen:destination.parameters.yaml -->
```yaml
version: 2.2
pipelines:
  - id: example
    status: running
    connectors:
      - id: example
        plugin: "weaviate"
        settings:
          # The class name as defined in the schema. A record will be saved
          # under this class unless it has the `weaviate.class` metadata field.
          # Type: string
          # Required: yes
          class: ""
          # Host of the Weaviate instance.
          # Type: string
          # Required: yes
          endpoint: ""
          # A Weaviate API key.
          # Type: string
          # Required: no
          auth.apiKey: ""
          # Mechanism specifies in which way the connector will authenticate to
          # Weaviate.
          # Type: string
          # Required: no
          auth.mechanism: "none"
          # WCS password
          # Type: string
          # Required: no
          auth.wcsCreds.password: ""
          # WCS username
          # Type: string
          # Required: no
          auth.wcsCreds.username: ""
          # Whether a UUID for records should be automatically generated. The
          # generated UUIDs are MD5 sums of record keys.
          # Type: bool
          # Required: no
          generateUUID: "false"
          # Name of the header configuring a module (e.g. `X-OpenAI-Api-Key`)
          # Type: string
          # Required: no
          moduleHeader.name: ""
          # Value for header given in `moduleHeader.name`.
          # Type: string
          # Required: no
          moduleHeader.value: ""
          # Scheme of the Weaviate instance.
          # Type: string
          # Required: no
          scheme: "https"
          # Maximum delay before an incomplete batch is written to the
          # destination.
          # Type: duration
          # Required: no
          sdk.batch.delay: "0"
          # Maximum size of batch before it gets written to the destination.
          # Type: int
          # Required: no
          sdk.batch.size: "0"
          # Allow bursts of at most X records (0 or less means that bursts are
          # not limited). Only takes effect if a rate limit per second is set.
          # Note that if `sdk.batch.size` is bigger than `sdk.rate.burst`, the
          # effective batch size will be equal to `sdk.rate.burst`.
          # Type: int
          # Required: no
          sdk.rate.burst: "0"
          # Maximum number of records written per second (0 means no rate
          # limit).
          # Type: float
          # Required: no
          sdk.rate.perSecond: "0"
          # The format of the output record. See the Conduit documentation for a
          # full list of supported formats
          # (https://conduit.io/docs/using/connectors/configuration-parameters/output-format).
          # Type: string
          # Required: no
          sdk.record.format: "opencdc/json"
          # Options to configure the chosen output record format. Options are
          # normally key=value pairs separated with comma (e.g.
          # opt1=val2,opt2=val2), except for the `template` record format, where
          # options are a Go template.
          # Type: string
          # Required: no
          sdk.record.format.options: ""
          # Whether to extract and decode the record key with a schema.
          # Type: bool
          # Required: no
          sdk.schema.extract.key.enabled: "true"
          # Whether to extract and decode the record payload with a schema.
          # Type: bool
          # Required: no
          sdk.schema.extract.payload.enabled: "true"
```
<!-- /readmegen:destination.parameters.yaml -->

(For more information about authentication in Weaviate, refer to the Weaviate https://weaviate.io/developers/wcs/guides/authentication.)
![scarf pixel](https://static.scarf.sh/a.png?x-pxid=3864585b-04e5-4a20-aa86-5bc4751f61b4)
