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

| name                 | description                                                                                                                         | required                                       | default value |
|----------------------|-------------------------------------------------------------------------------------------------------------------------------------|------------------------------------------------|---------------|
| `endpoint`           | Host of the Weaviate instance.                                                                                                      | true                                           | ""            |
| `scheme`             | Scheme of the Weaviate instance. Values: `https`, `http`.                                                                           | false                                          | "https"       |
| `class`              | The class name as defined in the schema. A record will be saved under this class unless it has the `weaviate.class` metadata field. | true                                           | ""            |
| `moduleHeader.name`  | Name of the header configuring a module (e.g. `X-OpenAI-Api-Key`).                                                                  | false                                          | ""            |
| `moduleHeader.value` | API key for the module defined above.                                                                                               | false                                          | ""            |
| `generateUUID`       | Generate a UUID for records (an MD5 sum of a record's key).                                                                         | false                                          | "false"       |
| `auth.mechanism`     | Specifies in which way the connector will authenticate to Weaviate. Values: `none`, `apiKey`, `wcsCredentials`.                     | false                                          | "none"        |
| `apiKey`             | A Weaviate API key.                                                                                                                 | Required if `auth.mechanism = apiKey`.         | ""            |
| `wcs.username`       | Weaviate Cloud Services (WCS) username.                                                                                             | Required if `auth.mechanism = wcsCredentials`. | ""            |
| `wcs.password`       | Weaviate Cloud Services (WCS) password.                                                                                             | Required if `auth.mechanism = wcsCredentials`. | ""            |

(For more information about authentication in Weaviate, refer to the Weaviate https://weaviate.io/developers/wcs/guides/authentication.)
