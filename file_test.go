package go_hubspot

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// getFormAPI Get default HubSpot Form API client
func getFileAPI() HubspotFileAPI {
	return NewHubspotFileAPI("apiKey", "portalId")
}

func getMockFileAPI(mockClient *IHTTPClientMock) HubspotFileAPI {
	return HubspotFileAPI{
		URLTemplate: "urltemplate",
		APIKey:      "apiKey",
		PortalID:    "portalId",
		httpClient:  mockClient,
	}
}

func TestUploadFile(t *testing.T){

  expectedOptions := options := FileUploadOptions{
		Access:                      "PRIVATE",
		Overwrite:                   true,
		DuplicateValidationStrategy: "NONE",
		DuplicateValidationScope:    "EXACT_FOLDER",
	}

	expectedLink := "link"

  mockHubspotHTTPClient := IHTTPClientMock{
    GetFunc: func(url string) (resp *http.Response, err error) { return nil, nil },
		DoFunc: func(req *http.Request) (resp *http.Response, err error) {
		  w := httptest.NewRecorder()
			if url == "https://example.com/form_id?hapikey=api_key&limit=50&after=" {

			  if req.Method != "POST" {
			    t.Errorf("Unexpected method: expected POST, got: %s", req.Method)
			  }

			  expectedRequest := dealCreationRequest{
					map[string]string{
						"dealname":                  "cardName",
						"dealstage":                 "stageName",
						"pipeline":                  "pipeline",
						"hubspot_owner_id":          "HubspotOwnerId",
						"application_id":            "applicationId",
						"validation_check_finished": "false",
					},
				}

			  w.WriteHeader(200)
				w.Write([]byte(`{
				  "id":"59358383010",
				  "createdAt":"2021-11-09T16:23:40.166Z",
				  "updatedAt":"2021-11-09T16:23:40.166Z",
				  "archived":false,
				  "parentFolderId":"59348125351",
				  "name":"fileName",
				  "path":"/folderPath/Test123",
				  "size":16,
				  "type":"OTHER",
				  "extension":"",
				  "defaultHostingUrl":"https://f.hubspotusercontent10.net/hubfs/8458264/folderPath/Test123",
				  "url":"https://f.hubspotusercontent10.net/hubfs/8458264/folderPath/Test123",
				  "isUsableInContent":true,
				  "access":"PRIVATE"
				}`))

			} else {
			  t.Errorf("Unexpected URL: %s", url)
			  return nil, nil
			}

      return w.Result(), nil
		},
	}

	api := getMockFileAPI(mockHubspotHTTPClient)

	fileContent := "Hello World File Content"

	gotLink, err := api.UploadFile([]byte(fileContent), "folderPath", "fileName")
	if err != nil {
	  t.Errorf("Error uploading file: %s", err.Error())
	}

}