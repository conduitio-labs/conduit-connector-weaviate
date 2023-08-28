# Conduit Connector for Weaviate
A [Conduit](https://conduit.io) destination connector for Weaviate.

## How to build?
Run `make build` to build the connector.

## Testing
To run the unit test, execute `make test`. 

To run the integration tests, you need to:
1. set the `OPENAI_APIKEY` environment variable, which stores an OpenAI API key.
2. execute `make test-integration`.

The Docker compose file at `test/docker-compose.yml` can be used to run an instance of Weaviate locally.

## Destination
A destination connector pushes data from upstream resources to an external resource via Conduit.

### Configuration

| name           | description                                                                                                                         | required | default value |
|----------------|-------------------------------------------------------------------------------------------------------------------------------------|----------|---------------|
| `endpoint`     | Host of the Weaviate instance.                                                                                                      | true     | ""            |
| `scheme`       | Scheme of the Weaviate instance. Values: https, http.                                                                               | false    | "https"       |
| `apiKey`       | A Weaviate API key.                                                                                                                 | false    | ""            |
| `class`        | The class name as defined in the schema. A record will be saved under this class unless it has the `weaviate.class` metadata field. | true     | ""            |
| `generateUUID` | Generate a UUID for records (an MD5 sum of a record's key).                                                                         | false    | "false"       |
