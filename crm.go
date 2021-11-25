package go_hubspot

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

type IHubspotCRMAPI interface {
	UpdateCompany(companyID string, jsonPayload *bytes.Buffer) error
	GetCompanyForContact(contactID string) (string, error)
	GetDealForCompany(companyID string) (string, error)
	SearchContacts(filterMap map[string]string, properties []string) ([]ContactResult, error)
}

type HubspotCRMAPI struct {
	APIKey     string
	httpClient IHTTPClient
}

type HubSpotContactSearchResponse struct {
	Total   int             `json:"total"`
	Results []ContactResult `json:"results"`
}

type Association struct {
	Id              string `json:"id"`
	AssociationType string `json:"type"`
}

type Associations struct {
	Results []Association `json:"results"`
}

type ContactResult struct {
	Id           string                  `json:"id"`
	Properties   map[string]string       `json:"properties"`
	Associations map[string]Associations `json:"associations"`
}

type filter struct {
	Value        string `json:"value"`
	PropertyName string `json:"propertyName"`
	Operator     string `json:"operator"`
}

type filterGroup struct {
	Filters []filter `json:"filters"`
}

type hubSpotSearchRequest struct {
	FilterGroups []filterGroup `json:"filterGroups"`
	Properties   []string      `json:"properties"`
}

// NewHubspotCRMAPI creates new HubspotCRMAPI with form ID and API key
func NewHubspotCRMAPI(apiKey string) HubspotCRMAPI {
	return HubspotCRMAPI{
		APIKey:     apiKey,
		httpClient: HTTPClient{},
	}
}

// UpdateCompany updates company details in HubSpot CRM
func (api HubspotCRMAPI) UpdateCompany(companyID string, jsonPayload *bytes.Buffer) error {
	url := fmt.Sprintf(
		"https://api.hubapi.com/crm/v3/objects/companies/%s?hapikey=%s",
		companyID,
		api.APIKey,
	)

	req, _ := http.NewRequest("PATCH", url, jsonPayload)

	req.Header.Set("Content-Type", "application/json")

	resp, err := api.httpClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		var bodyString string
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			bodyString = ""
		} else {
			bodyString = string(bodyBytes)
		}

		return errors.New(fmt.Sprintf(
			"Failed to update company with ID '%s': %s",
			companyID,
			bodyString,
		))
	}

	return nil
}

// GetCompanyForContact returns the company id for the contact with the given id
// If no company is found "" is returned, no error is thrown
func (api HubspotCRMAPI) GetCompanyForContact(contactID string) (string, error) {
	url := fmt.Sprintf(
		"https://api.hubapi.com/crm/v3/objects/contacts/%s?associations=company&archived=false&hapikey=%s",
		contactID,
		api.APIKey,
	)

	req, _ := http.NewRequest("GET", url, nil)

	resp, err := api.httpClient.Do(req)
	if err != nil {
		return "", nil
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var hubspotResp ContactResult
	err = json.Unmarshal(body, &hubspotResp)
	if err != nil {
		return "", err
	}

	companies := hubspotResp.Associations["companies"].Results

	if len(companies) == 0 {
		return "", nil
	} else if len(companies) > 1 {
		return "", errors.New(fmt.Sprintf("There are multiple companies associated with contact '%s' there should only be one", contactID))
	} else {
		return companies[0].Id, nil
	}
}

// GetDealForCompany returns the deal id associated with the given companyID
// Returns "" with nil error if no company exists
func (api HubspotCRMAPI) GetDealForCompany(companyID string) (string, error) {
	url := fmt.Sprintf(
		"https://api.hubapi.com/crm/v3/objects/companies/%s/associations/deal?limit=500&hapikey=%s",
		companyID,
		api.APIKey,
	)

	req, _ := http.NewRequest("GET", url, nil)

	resp, err := api.httpClient.Do(req)
	if err != nil {
		return "", nil
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	log.Infof("Raw response: %s", string(body))

	var associationResp Associations
	err = json.Unmarshal(body, &associationResp)
	if err != nil {
		return "", err
	}

	dealAssociations := associationResp.Results

	if len(dealAssociations) == 0 {
		return "", nil
	} else if len(dealAssociations) > 1 {
		return "", errors.New(fmt.Sprintf("There are multiple deals associated with company '%s' there should only be one", companyID))
	} else {
		return dealAssociations[0].Id, nil
	}
}

// SearchContacts searches for contacts with the provided filters and returns properties for the results found
func (api HubspotCRMAPI) SearchContacts(filterMap map[string]string, properties []string) ([]ContactResult, error) {
	censoredUrl := fmt.Sprintf("https://api.hubapi.com/crm/v3/objects/contacts/search?hapikey=%s", "<censored>")
	url := fmt.Sprintf("https://api.hubapi.com/crm/v3/objects/contacts/search?hapikey=%s", api.APIKey)

	var filters = make([]filter, len(filterMap))
	filterIndex := 0
	for propertyName, value := range filterMap {
		filters[filterIndex] = filter{
			Value:        value,
			PropertyName: propertyName,
			Operator:     "EQ",
		}
	}

	searchQuery := hubSpotSearchRequest{
		FilterGroups: []filterGroup{
			{
				Filters: filters,
			},
		},
		Properties: properties,
	}

	log.Infof("Making query to contact search endpoint (%s) with: %#v", censoredUrl, searchQuery)

	payloadBuf := new(bytes.Buffer)
	err := json.NewEncoder(payloadBuf).Encode(searchQuery)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, payloadBuf)
	req.Header.Set("Content-Type", "application/json")

	log.Infof("Query Payload: %s", payloadBuf.String())

	if err != nil {
		return nil, err
	}

	resp, err := api.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var hubspotResp HubSpotContactSearchResponse

	log.Infof("Raw response: %s", string(body))
	err = json.Unmarshal(body, &hubspotResp)

	if err != nil {
		return nil, err
	}

	return hubspotResp.Results, nil
}
