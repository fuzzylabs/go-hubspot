package ehe_hubspot

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/fuzzylabs/ehe-hubspot/schema"
	"io/ioutil"

	log "github.com/sirupsen/logrus"
)

// HubspotFormAPI is the structure to interact with Hubspot Form API
type HubspotFormAPI struct {
	URLTemplate string
	FormID      string
	APIKey      string
	httpClient  IHTTPClient
}

// FormValue form value
type FormValue struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// Submission form submission
type Submission struct {
	SubmittedAt int64       `json:"submittedAt"`
	Values      []FormValue `json:"values"`
}

// Paging optional paging information
type Paging struct {
	Next map[string]string `json:"next"`
}

// HubspotResponse response of the form API
type HubspotResponse struct {
	Results []Submission `json:"results"`
	Paging  *Paging      `json:"paging"`
}

// NewHubspotFormAPI creates new HubspotFormAPI with form ID and API key
func NewHubspotFormAPI(formID string, apiKey string) HubspotFormAPI {
	return HubspotFormAPI{
		URLTemplate: "https://api.hubapi.com/form-integrations/v1/submissions/forms/%s?hapikey=%s&limit=50&after=%s",
		FormID:      formID,
		APIKey:      apiKey,
		httpClient:  HTTPClient{},
	}
}

// GetSubmissionMap transforms Hubspot form submission into a string map
func GetSubmissionMap(submission Submission) map[string]string {
	submissionMap := map[string]string{}
	for _, value := range submission.Values {
		submissionMap[value.Name] = value.Value
	}
	return submissionMap
}

// GetPageURL gets query URL for a page of results
func (api HubspotFormAPI) GetPageURL(after string) string {
	return fmt.Sprintf(
		api.URLTemplate,
		api.FormID,
		api.APIKey,
		after,
	)
}

// Query queries Hubspot for a page of form results
func (api HubspotFormAPI) Query(after string) (*HubspotResponse, error) {
	url := api.GetPageURL(after)

	log.Println(url)

	resp, err := api.httpClient.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var hubspotResp HubspotResponse
	err = json.Unmarshal(body, &hubspotResp)

	if err != nil {
		return nil, err
	}

	return &hubspotResp, nil
}

// GetNextAfter get next page from the response
func (r HubspotResponse) GetNextAfter() (string, error) {
	if r.Paging != nil {
		log.Infof("Try next %s", r.Paging.Next["after"])
		return r.Paging.Next["after"], nil
	}
	return "", errors.New("There is no next page")
}

// SearchForApplicationID searches for a submission on Hubspot for a given Application ID
func (api HubspotFormAPI) SearchForApplicationID(applicationId string) (*schema.ApplicationForm, error) {
	log.Printf("Searching for submission with Application ID %s\n", applicationId)

	after := ""

	for {
		hubspotResp, err := api.Query(after)

		if err != nil {
			return nil, err
		}

		submissionMap, err := hubspotResp.GetByApplicationID(applicationId)

		if err != nil {
			if err.Error() == fmt.Sprintf("Submission with applicationId `%s` not found", applicationId) {
				after, err = hubspotResp.GetNextAfter()
				if err != nil {
					return nil, err
				}
			} else {
				return nil, err
			}
		} else {

			var submission schema.ApplicationForm
			submissionJSON, err := json.Marshal(submissionMap)
			if err != nil {
				return nil, err
			}

			err = json.Unmarshal(submissionJSON, &submission)
			if err != nil {
				return nil, err
			}

			// Manually set company name since protobuf sets the companyName json field to `company_name` when actually
			// it is `company`
			companyName, ok := submissionMap["company"]

			if !ok {
				return nil, errors.New(fmt.Sprintf("The submission '%s' found does not have a company number", applicationId))
			}

			submission.CompanyName = companyName

			return &submission, err
		}
	}

}

// GetByApplicationID searches Hubspot results for a form with a given application ID
func (r HubspotResponse) GetByApplicationID(applicationId string) (map[string]string, error) {
	for _, result := range r.Results {
		submission := GetSubmissionMap(result)

		if submission["application_id"] == applicationId {
			log.Printf("Found: %#v", submission)
			return submission, nil
		}
	}

	return nil, errors.New(fmt.Sprintf("Submission with applicationId `%s` not found", applicationId))
}
