# EHE HubSpot library

This library provides generic methods for the interation with HubSpot [Forms](https://legacydocs.hubspot.com/docs/methods/forms/forms_overview), CRM ([Contacts](https://developers.hubspot.com/docs/api/crm/contacts) and [Companies](https://developers.hubspot.com/docs/api/crm/companies)) and [DealFlow](https://developers.hubspot.com/docs/api/crm/deals) APIs

## Usage
You can install the module directly from GitHub

```shell
go get -u github.com/fuzzylabs/ehe-hubspot@<version>
```

Where version can point to a commit hash or a branch, for example:

```shell
go get -u github.com/fuzzylabs/ehe-hubspot@9302e1d
```

or 

```shell
go get -u github.com/fuzzylabs/ehe-hubspot@master
```

You can then import the library as follows:
```go
import (
	hubspot "github.com/fuzzylabs/ehe-hubspot"
)
```

## Mocking
`moq` is used to generate mocks:
* Mocks for external interfaces to use within unit tests
* Mocks for `ehe-hubspot` API interfaces, for to make testing of applications that use the library easier

```
go generate
```

## Testing
```
go vet
go test -coverprofile=coverage.out
go tool cover -html=coverage.out # To view test coverage
```