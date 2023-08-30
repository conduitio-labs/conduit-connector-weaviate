# Conduit Connector for Weaviate
A [Conduit](https://conduit.io) destination connector for Weaviate.

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

| name                     | description                                                                                                                             | required                                       | default value |
|--------------------------|-----------------------------------------------------------------------------------------------------------------------------------------|------------------------------------------------|---------------|
| `endpoint`               | Host of the Weaviate instance.                                                                                                          | true                                           | ""            |
| `scheme`                 | Scheme of the Weaviate instance.<br/> Values: `https`, `http`.                                                                          | false                                          | "https"       |
| `class`                  | The class name as defined in the schema.<br/>A record will be saved under this class unless it has the `weaviate.class` metadata field. | true                                           | ""            |
| `moduleHeader.name`      | Name of the header configuring a module (e.g. `X-OpenAI-Api-Key`).                                                                      | false                                          | ""            |
| `moduleHeader.value`     | API key for the module defined above.                                                                                                   | false                                          | ""            |
| `generateUUID`           | Generate a UUID for records (an MD5 sum of a record's key).                                                                             | false                                          | "false"       |
| `auth.mechanism`         | Specifies in which way the connector will authenticate to Weaviate. <br/>Values: `none`, `apiKey`, `wcsCreds`.                    | false                                          | "none"        |
| `auth.apiKey`            | A Weaviate API key.                                                                                                                     | Required if `auth.mechanism = apiKey`.         | ""            |
| `auth.wcsCreds.username` | Weaviate Cloud Services (WCS) username.                                                                                                 | Required if `auth.mechanism = wcsCreds`. | ""            |
| `auth.wcsCreds.password` | Weaviate Cloud Services (WCS) password.                                                                                                 | Required if `auth.mechanism = wcsCreds`. | ""            |

(For more information about authentication in Weaviate, refer to the Weaviate https://weaviate.io/developers/wcs/guides/authentication.)
