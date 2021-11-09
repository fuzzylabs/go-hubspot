package ehe_hubspot

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
)

type IHubspotFileAPI interface {
	GetPageURL() string
	UploadFile(file, folderPath, fileName, options string) error
}

// HubspotFileAPI is the structure to interact with Hubspot File API
type HubspotFileAPI struct {
	URLTemplate string
	APIKey      string
	PortalID    string
	httpClient  IHTTPClient
}

// FileUploadResponse response of the file API
type FileUploadResponse struct {
	Id string `json:"id"`
}

// NewHubspotFileAPI creates new HubspotFileAPI and API key
func NewHubspotFileAPI(apiKey string, portalId string) HubspotFileAPI {
	return HubspotFileAPI{
		URLTemplate: "https://api.hubapi.com/files/v3/files?hapikey=%s",
		APIKey:      apiKey,
		PortalID:    portalId,
		httpClient:  HTTPClient{},
	}
}

// GetPageURL gets query URL for a page of results
func (api HubspotFileAPI) GetPageURL() string {
	return fmt.Sprintf(
		api.URLTemplate,
		api.APIKey,
	)
}

type FileUploadOptions struct {
	Access                      string `json:"access"`
	Overwrite                   bool   `json:"overwrite"`
	DuplicateValidationStrategy string `json:"duplicateValidationStrategy"`
	DuplicateValidationScope    string `json:"duplicateValidationScope"`
}

// TODO
func (api HubspotFileAPI) UploadFile(file []byte, folderPath, fileName string) (string, error) {
	var data bytes.Buffer
	w := multipart.NewWriter(&data)

	fileWriter, err := w.CreateFormFile("file", fileName)
	if err != nil {
		return "", errors.New(fmt.Sprintf("Error while creating a file writer: %s", err.Error()))
	}

	_, err = fileWriter.Write(file)
	if err != nil {
		return "", errors.New(fmt.Sprintf("Error while writing a file: %s", err.Error()))
	}

	err = w.WriteField("folderPath", folderPath)
	if err != nil {
		return "", errors.New(fmt.Sprintf("Error while writing folder name: %s", err.Error()))
	}

	options := FileUploadOptions{
		Access:                      "PRIVATE",
		Overwrite:                   true,
		DuplicateValidationStrategy: "NONE",
		DuplicateValidationScope:    "EXACT_FOLDER",
	}

	optionsBytes, err := json.Marshal(options)
	if err != nil {
		return "", errors.New(fmt.Sprintf("Error while marshaling options: %s", err.Error()))
	}

	err = w.WriteField("options", string(optionsBytes))
	if err != nil {
		return "", errors.New(fmt.Sprintf("Error while writing options: %s", err.Error()))
	}

	err = w.Close()
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", api.GetPageURL(), &data)
	if err != nil {
		return "", errors.New(fmt.Sprintf("Error while constructing a request: %s", err.Error()))
	}

	req.Header.Set("Content-Type", w.FormDataContentType())

	resp, err := api.httpClient.Do(req)
	if err != nil {
		return "", errors.New(fmt.Sprintf("Error while making a request: %s", err.Error()))
	}

	if resp.StatusCode != 201 {
		return "", errors.New(fmt.Sprintf("Request to HubSpot File API failed: %s", resp.Status))
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", errors.New(fmt.Sprintf("Error while reading response body: %s", err.Error()))
	}

	var hubspotResp FileUploadResponse
	err = json.Unmarshal(body, &hubspotResp)
	if err != nil {
		return "", errors.New(fmt.Sprintf("Error while unmarshaling HubSpot response: %s", err.Error()))
	}

	url := fmt.Sprintf(
		"https://app.hubspot.com/file-preview/%s/file/%s",
		api.PortalID,
		hubspotResp.Id,
	)

	return url, nil
}
