package ehe_hubspot

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"

	log "github.com/sirupsen/logrus"
)

type IHubspotFileAPI interface {
	GetPageURL() string
	UploadFile(file, folderPath, fileName, options string) error
}

// HubspotFileAPI is the structure to interact with Hubspot File API
type HubspotFileAPI struct {
	URLTemplate string
	APIKey      string
	httpClient  IHTTPClient
}

// HubspotResponse response of the file API
type FileUploadResponse struct {
	Results []Submission `json:"results"`
	Paging  *Paging      `json:"paging"`
}

// NewHubspotFormAPI creates new HubspotFormAPI with form ID and API key
func NewHubspotFormAPI(formID string, apiKey string) HubspotFormAPI {
	return HubspotFormAPI{
		URLTemplate: "https://api.hubapi.com/files/v3/files?hapikey=%s",
		APIKey:      apiKey,
		httpClient:  HTTPClient{},
	}
}

// GetPageURL gets query URL for a page of results
func (api HubspotFormAPI) GetPageURL() string {
	return fmt.Sprintf(
		api.URLTemplate,
		api.APIKey,
	)
}

// SearchForKeyValue searches for a form submission on Hubspot for a given key-value pair
func (api HubspotFileAPI) UploadFile(file , folderPath, fileName, options string) error {
	log.Printf("Uploading file to HubSpot")

	req, _ := http.NewRequest("POST", api.GetPageURL(), jsonPayload)

	req.Header.Set("Content-Type", "application/json")

	resp, err := api.httpClient.Do(req)
	if err != nil {
		return "", nil
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var hubspotResp FileUploadResponse
	err = json.Unmarshal(body, &hubspotResp)
	if err != nil {
		return "", err
	}

}
