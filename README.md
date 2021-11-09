# Go HubSpot library

This library provides generic methods for interaction with HubSpot [Forms](https://legacydocs.hubspot.com/docs/methods/forms/forms_overview), CRM ([Contacts](https://developers.hubspot.com/docs/api/crm/contacts), [Companies](https://developers.hubspot.com/docs/api/crm/companies)) and [DealFlow](https://developers.hubspot.com/docs/api/crm/deals) APIs

## Usage
You can install the module directly from GitHub

```shell
go get -u github.com/fuzzylabs/go-hubspot@<version>
```

Where version can point to a commit hash or a branch, for example:

```shell
go get -u github.com/fuzzylabs/go-hubspot@9302e1d
```

or 

```shell
go get -u github.com/fuzzylabs/go-hubspot@master
```

You can then import the library as follows:
```go
import (
	hubspot "github.com/fuzzylabs/go-hubspot"
)
```

## Examples
Search for form submissions with the first name John:
```go
package main

import (
	hubspot "github.com/fuzzylabs/go-hubspot"
)

func main() {
	api := hubspot.NewHubspotFormAPI("form-id", "hapikey")
	_, _ = api.SearchForKeyValue("firstname", "John")
}
```

Get company ID associated with a contact:
```go
package main

import (
	hubspot "github.com/fuzzylabs/go-hubspot"
)

func main() {
	api := hubspot.NewHubspotCRMAPI("hapikey")
	_, _ = api.GetCompanyForContact("123456")
}
```

Update move a DealFlow card to another column (i.e. update its `dealstage` property):
```go
package main

import (
	hubspot "github.com/fuzzylabs/go-hubspot"
)

func main() {
	api := hubspot.NewHubspotDealFlowAPI("hapikey")
	_ = api.UpdateDealFlowCard(
		"123456",
		map[string]string{
			"dealstage": "another-stage",
        },
    )
}
```

Upload a file to the HubSpot CRM
```go
package main

import (
	hubspot "github.com/fuzzylabs/go-hubspot"
)

func main() {
	fileApi := hubspot.NewHubspotFileAPI("hapikey", "portalId")
	fileUrl, err := fileApi.UploadFile(bytes, "folder path", "file name")
}
```

## Mocking
`moq` is used to generate mocks:
* Mocks for external interfaces to use within unit tests
* Mocks for `go-hubspot` API interfaces, to make testing of applications that use the library easier

```
go generate
```

## Testing
```
go vet
go test -coverprofile=coverage.out
go tool cover -html=coverage.out # To view test coverage
```
