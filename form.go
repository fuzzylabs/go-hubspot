package go_hubspot

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"

	log "github.com/sirupsen/logrus"
)

type IHubspotFormAPI interface {
	GetPageURL(after string) string
	Query(after string) (*HubspotResponse, error)
	SearchForKeyValue(key string, value string) (map[string]HubspotFormField, error)
}

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
		URLTemplate: "https://api.hubapi.com/form-integrations/v1/submissions/forms/%s?limit=50&after=%s",
		FormID:      formID,
		APIKey:      apiKey,
		httpClient:  HTTPClient{},
	}
}

type HubspotFieldType int64

const (
	SingleValue HubspotFieldType = iota
	MultipleValues
)

type HubspotFormField struct {
	Type           HubspotFieldType
	SingleValue    string
	MultipleValues []string
}

func (field HubspotFormField) ForceMultipleValues() []string {
	if field.Type == SingleValue {
		return []string{field.SingleValue}
	} else {
		return field.MultipleValues
	}
}

// GetSubmissionMap transforms Hubspot form submission into a string map
func GetSubmissionMap(submission Submission) map[string]HubspotFormField {
	submissionMap := map[string]HubspotFormField{}
	for _, value := range submission.Values {
		if field, ok := submissionMap[value.Name]; ok {
			if field.Type == SingleValue {
				field.Type = MultipleValues
				field.MultipleValues = make([]string, 1)
				field.MultipleValues[0] = field.SingleValue
				field.SingleValue = ""
			}
			field.MultipleValues = append(field.MultipleValues, value.Value)
			submissionMap[value.Name] = field
		} else {
			field := HubspotFormField{
				Type:        SingleValue,
				SingleValue: value.Value,
			}
			submissionMap[value.Name] = field
		}
	}
	return submissionMap
}

// GetPageURL gets query URL for a page of results
func (api HubspotFormAPI) GetPageURL(after string) string {
	return fmt.Sprintf(
		api.URLTemplate,
		api.FormID,
		after,
	)
}

// Query queries Hubspot for a page of form results
func (api HubspotFormAPI) Query(after string) (*HubspotResponse, error) {
	url := api.GetPageURL(after)

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", api.APIKey))
	req.Header.Set("Content-Type", "application/json")

	resp, err := api.httpClient.Do(req)
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

// SearchForKeyValue searches for a form submission on Hubspot for a given key-value pair
func (api HubspotFormAPI) SearchForKeyValue(key string, value string) (map[string]HubspotFormField, error) {
	log.Printf("Searching for submission with %s = %s\n", key, value)

	after := ""

	for {
		hubspotResp, err := api.Query(after)

		if err != nil {
			return nil, err
		}

		submissionMap, err := hubspotResp.GetByKeyValue(key, value)

		if err != nil {
			if err.Error() == fmt.Sprintf("Submission with %s `%s` not found", key, value) {
				after, err = hubspotResp.GetNextAfter()
				if err != nil {
					return nil, err
				}
			} else {
				return nil, err
			}
		} else {
			return submissionMap, nil
		}
	}

}

// GetByKeyValue searches Hubspot results for a form with a given key-value pair
func (r HubspotResponse) GetByKeyValue(key string, value string) (map[string]HubspotFormField, error) {
	for _, result := range r.Results {
		submission := GetSubmissionMap(result)

		if submission[key].Type == SingleValue && submission[key].SingleValue == value {
			log.Printf("Found: %#v", submission)
			return submission, nil
		}
	}

	return nil, errors.New(fmt.Sprintf("Submission with %s `%s` not found", key, value))
}
