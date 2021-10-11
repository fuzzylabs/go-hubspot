package ehe_hubspot

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	log "github.com/sirupsen/logrus"
)

type HubspotCRMAPI struct {
	APIKey     string
	httpClient IHTTPClient
}

type HubSpotContactSearchResponse struct {
	Total   int             `json:"total"`
	Results []ContactResult `json:"results"`
}

type association struct {
	Id              string `json:"id"`
	AssociationType string `json:"type"`
}

type associations struct {
	Results []association `json:"results"`
}

type ContactResult struct {
	Id           string                  `json:"id"`
	Properties   map[string]string       `json:"properties"`
	Associations map[string]associations `json:"associations"`
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

	var associationResp associations
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

// searchContactsForApplicationId searches for contacts with the given application ID and returns every result it finds
func (api HubspotCRMAPI) searchContactsForApplicationId(application_id string) ([]ContactResult, error) {
	url := fmt.Sprintf("https://api.hubapi.com/crm/v3/objects/contacts/search?hapikey=%s", api.APIKey)

	searchQuery := hubSpotSearchRequest{
		FilterGroups: []filterGroup{
			{
				Filters: []filter{
					{
						Value:        application_id,
						PropertyName: "application_id",
						Operator:     "EQ",
					},
				},
			},
		},
		Properties: []string{
			"contact_id",
			"company_number",
		},
	}

	log.Infof("Making query to contact search endpoint (%s) with: %#v", url, searchQuery)

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

// GetContactID searches for contacts on HubSpot with a specific company number
func (api HubspotCRMAPI) GetContactID(applicationId string, companyNumber string) (string, error) {
	contacts, err := api.searchContactsForApplicationId(applicationId)
	if err != nil {
		return "", err
	}

	if len(contacts) == 0 {
		return "", errors.New(fmt.Sprintf("Could not find a contact with application ID '%s'", applicationId))
	} else if len(contacts) > 1 {
		return "", errors.New(fmt.Sprintf("Multiple contacts found for application ID '%s' there should only be one", applicationId))
	} else {
		contact := contacts[0]
		if contact.Properties["company_number"] != companyNumber {
			return "", errors.New(fmt.Sprintf("Contact with application ID '%s' has company number '%s', but we expected '%s'", applicationId, contact.Properties["company_number"], companyNumber))
		}
		return contact.Id, nil
	}
}
