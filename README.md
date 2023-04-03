# Go GraphQL Server Example

This service acts as a GraphQL interface for querying the Hearthstone card related APIs. 

## Running the Service

Ensure you have Go 1.20 installed. Simply execute the command `go run app/*.go`. 

## Using the Service

While the service is running, you can open http://localhost:8080/sandbox in your browser to easily view the schema and interact with the API. In the explorer section, you will need to select the Headers section on the bottom. You must enter a header key of `Authorization` and a value of `Bearer {token_value}`. The token value must be a valid Battle.net API token. In structions for acquiring such tokens can be found here: https://develop.battle.net/documentation/guides/getting-started

Metadata from the Hearthstone API will be cached indefiitely after the first request is processed. In order to refresh the cache, the server must be restarted.

## Possible Improvements
- add more filtering options
- add more card response fields
- time out the cache
- improve error handling
- write tests
