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

| name                 | description                                                                                                                         | required | default value |
|----------------------|-------------------------------------------------------------------------------------------------------------------------------------|----------|---------------|
| `endpoint`           | Host of the Weaviate instance.                                                                                                      | true     | ""            |
| `scheme`             | Scheme of the Weaviate instance. Values: https, http.                                                                               | false    | "https"       |
| `class`              | The class name as defined in the schema. A record will be saved under this class unless it has the `weaviate.class` metadata field. | true     | ""            |
| `moduleHeader.name`  | Name of the header configuring a module (e.g. `X-OpenAI-Api-Key`).                                                                  | false    | ""            |
| `moduleHeader.value` | API key for the module defined above.                                                                                               | false    | ""            |
| `generateUUID`       | Generate a UUID for records (an MD5 sum of a record's key).                                                                         | false    | "false"       |

#### Authentication

In addition to the configuration mentioned above, this connector supports the following options for authentication 
into a Weaviate Cloud Services (WCS) instance:
* Using the account owner's WCS username and password, or
* Using an API key (recommended method).

Only one of them can be used at a time.

(For more information about authentication in Weaviate, refer to the Weaviate https://weaviate.io/developers/wcs/guides/authentication.)

##### Authentication using the WCS username and password
* `wcs.username` (string, required)
* `wcs.password` (string, required)

##### Authentication using an API key
* `apiKey` (string, required)
