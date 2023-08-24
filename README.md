# Conduit Connector for Weaviate
A [Conduit](https://conduit.io) destination connector for Weaviate.

## How to build?
Run `make build` to build the connector.

## Testing
Run `make test` to run all the unit tests. Run `make test-integration` to run the integration tests.

The Docker compose file at `test/docker-compose.yml` can be used to run an instance of Weaviate locally.

## Destination
A destination connector pushes data from upstream resources to an external resource via Conduit.

### Configuration

| name                       | description                                | required | default value |
|----------------------------|--------------------------------------------|----------|---------------|
| `destination_config_param` | Description of `destination_config_param`. | true     | 1000          |

## Known Issues & Limitations
* Known issue A
* Limitation A

## Planned work
- [ ] Item A
- [ ] Item B
