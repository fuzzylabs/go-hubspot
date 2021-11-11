package go_hubspot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/google/go-cmp/cmp"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func getMockFileAPI(mockClient *IHTTPClientMock) HubspotFileAPI {
	return HubspotFileAPI{
		URLTemplate: "https://api.hubapi.com/files/v3/files?hapikey=%s",
		APIKey:      "apiKey",
		PortalID:    "portalId",
		httpClient:  mockClient,
	}
}

func TestUploadFile(t *testing.T) {

	expectedLink := "https://app.hubspot.com/file-preview/portalId/file/59358383010"

	fileContent := []byte("Hello World File Content")

	mockHubspotHTTPClient := IHTTPClientMock{
		GetFunc: func(url string) (resp *http.Response, err error) { return nil, nil },
		DoFunc: func(req *http.Request) (resp *http.Response, err error) {
			url := fmt.Sprintf("%s", req.URL)

			w := httptest.NewRecorder()
			if url == "https://api.hubapi.com/files/v3/files?hapikey=apiKey" {

				if req.Method != "POST" {
					t.Errorf("Unexpected method: expected POST, got: %s", req.Method)
				}

				expectedOptions := FileUploadOptions{
					Access:                      "PRIVATE",
					Overwrite:                   true,
					DuplicateValidationStrategy: "NONE",
					DuplicateValidationScope:    "EXACT_FOLDER",
				}

				var gotOptions FileUploadOptions
				err := json.Unmarshal([]byte(req.FormValue("options")), &gotOptions)
				if err != nil {
					t.Errorf("Error unmarshalling request options to a FileUploadOptions struct: %s" + err.Error())
				}

				if !cmp.Equal(expectedOptions, gotOptions) {
					t.Errorf("Unexpected options for file upload expected:\n %v\ngot:\n%v", expectedOptions, gotOptions)
				}

				file, _, err := req.FormFile("file")
				defer file.Close()
				if err != nil {
					t.Errorf("Error getting file from form request: %s", err.Error())
				}

				buf := bytes.NewBuffer(nil)
				if _, err := io.Copy(buf, file); err != nil {
					return nil, err
				}

				if !cmp.Equal(buf.Bytes(), fileContent) {
					t.Errorf("Incorrect file information uploaded")
				}

				gotFolderPath := req.FormValue("folderPath")
				if gotFolderPath != "folderPath" {
					t.Errorf("Incorrect folderPath used, expected: folderPath, got: %s", req.FormValue("folderPath"))
				}

				w.WriteHeader(201)
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

	api := getMockFileAPI(&mockHubspotHTTPClient)

	gotLink, err := api.UploadFile(fileContent, "folderPath", "fileName")
	if err != nil {
		t.Errorf("Error uploading file: %s", err.Error())
	}

	if expectedLink != gotLink {
		t.Errorf("File upload returned unexpected link, expected: %s got: %s", expectedLink, gotLink)
	}

}
