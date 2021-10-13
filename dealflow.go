package ehe_hubspot

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	log "github.com/sirupsen/logrus"
)

type IHubspotDealFlowAPI interface {
	AssociateDealFlowCard(dealId, assocId string, assocType CardAssociation) error
	CreateDealFlowCard(
		cardName string,
		contactID string,
		companyID string,
		applicationId string,
	) (*DealCreationResponse, error)
	UpdateDealFlowCard(
		dealId string,
		dealName string,
		dealStage string,
		applicationId string,
		dealValidationCheckFinished bool,
	) error
	UpdateDealFlowCardValidationStatus(
		dealId string,
		dealValidationCheckFinished bool,
	) error
}

type HubspotDealFlowAPI struct {
	APIKey     string
	httpClient IHTTPClient
}

// dealCreationRequestProperties is a representation of the deal creation request to HubSpot
type dealCreationRequestProperties struct {
	DealName                string `json:"dealname"`
	DealStage               string `json:"dealstage"`
	Pipeline                string `json:"pipeline"`
	ApplicationId           string `json:"application_id"`
	HubspotOwnerId          string `json:"hubspot_owner_id"`
	ValidationCheckFinished string `json:"validation_check_finished"`
}

// dealCreationRequest is a representation of the deal creation request to HubSpot
type dealCreationRequest struct {
	Properties dealCreationRequestProperties `json:"properties"`
}

// dealCreationResponseProperties is a representation of the deal creation response from HubSpot
type dealCreationResponseProperties struct {
	Amount                  string `json:"amount"`
	CloseDate               string `json:"closedate"`
	CreateDate              string `json:"createdate"`
	DealName                string `json:"dealname"`
	DealStage               string `json:"dealstage"`
	HsLastModifiedDate      string `json:"hs_lastmodifieddate"`
	HubspotOwnerId          string `json:"hubspot_owner_id"`
	Pipeline                string `json:"pipeline"`
	ValidationCheckFinished string `json:"validation_check_finished"`
}

// DealCreationResponse is a representation of the deal creation response from HubSpot
type DealCreationResponse struct {
	Id         string                         `json:"id"`
	Properties dealCreationResponseProperties `json:"properties"`
	CreatedAt  string                         `json:"createdAt"`
	UpdatedAt  string                         `json:"updatedAt"`
	Archived   bool                           `json:"archived"`
}

type dealUpdateRequestProperties struct {
	DealName                string `json:"dealname"`
	DealStage               string `json:"dealstage"`
	ApplicationId           string `json:"uuid"`
	ValidationCheckFinished string `json:"validation_check_finished"`
}

type dealUpdateRequest struct {
	Properties dealUpdateRequestProperties `json:"properties"`
}

type CardAssociation int64

const (
	Company CardAssociation = iota
	Contact
)

// NewHubspotDealFlowAPI creates new HubspotDealFlowAPI with form ID and API key
func NewHubspotDealFlowAPI(apiKey string) HubspotDealFlowAPI {
	return HubspotDealFlowAPI{
		APIKey:     apiKey,
		httpClient: HTTPClient{},
	}
}

// AssociateDealFlowCard associates a deal flow card with a company or contact using the internal HubSpot dealId and companyId/contactid
// Choose whether to associate a company or contact by setting assocType to "contact" or "company"
func (api HubspotDealFlowAPI) AssociateDealFlowCard(dealId, assocId string, assocType CardAssociation) error {
	var url string
	switch assocType {
	case Company:
		url = fmt.Sprintf(
			"https://api.hubapi.com/crm/v3/objects/deals/%s/associations/company/%s/deal_to_company?hapikey=%s",
			dealId,
			assocId,
			api.APIKey,
		)
	case Contact:
		url = fmt.Sprintf(
			"https://api.hubapi.com/crm/v3/objects/deals/%s/associations/contact/%s/deal_to_contact?hapikey=%s",
			dealId,
			assocId,
			api.APIKey,
		)
	}

	req, err := http.NewRequest("PUT", url, nil)
	if err != nil {
		return err
	}

	_, err = api.httpClient.Do(req)
	if err != nil {
		return err
	}

	return nil
}

// CreateDealFlowCard creates a deal flow card with the given parameters in HubSpot,
// and associates it with a company and contact based on applicationId
func (api HubspotDealFlowAPI) CreateDealFlowCard(
	cardName string,
	contactID string,
	companyID string,
	applicationId string,
) (*DealCreationResponse, error) {

	log.Infof("Creating a deal flow card")

	stageName, exists := os.LookupEnv("DEALFLOW_STARTING_STAGE")
	if !exists {
		return nil, errors.New("DEALFLOW_STARTING_STAGE environment variable is unset")
	}

	pipeline, exists := os.LookupEnv("DEALFLOW_PIPELINE_NAME")
	if !exists {
		return nil, errors.New("DEALFLOW_PIPELINE_NAME environment variable is unset")
	}

	ownerId, exists := os.LookupEnv("DEALFLOW_OWNER_ID")
	if !exists {
		return nil, errors.New("DEALFLOW_OWNER_ID environment variable is unset")
	}

	url := fmt.Sprintf("https://api.hubapi.com/crm/v3/objects/deals?hapikey=%s", api.APIKey)

	creationRequest := dealCreationRequest{
		dealCreationRequestProperties{
			cardName,
			stageName,
			pipeline,
			applicationId,
			ownerId,
			"false",
		},
	}

	payloadBuf := new(bytes.Buffer)
	err := json.NewEncoder(payloadBuf).Encode(creationRequest)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, payloadBuf)
	req.Header.Set("Content-Type", "application/json")
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

	var hubspotResp DealCreationResponse

	log.Infof("Raw response: %s", string(body))
	err = json.Unmarshal(body, &hubspotResp)
	if err != nil {
		return nil, err
	}

	// Associate the deal with a company based on the application id
	err = api.AssociateDealFlowCard(hubspotResp.Id, companyID, Company)
	if err != nil {
		return nil, err
	}

	// Associate the deal with a contact based on the application id
	err = api.AssociateDealFlowCard(hubspotResp.Id, contactID, Contact)
	if err != nil {
		return nil, err
	}

	return &hubspotResp, nil

}

// UpdateDealFlowCard updates the deal flow card attached to the given id with the given information
func (api HubspotDealFlowAPI) UpdateDealFlowCard(
	dealId string,
	dealName string,
	dealStage string,
	applicationId string,
	dealValidationCheckFinished bool,
) error {

	log.Infof("Updating a deal flow card")

	url := fmt.Sprintf("https://api.hubapi.com/crm/v3/objects/deals/%s?hapikey=%s", dealId, api.APIKey)

	updateRequest := dealUpdateRequest{
		dealUpdateRequestProperties{
			dealName,
			dealStage,
			applicationId,
			strconv.FormatBool(dealValidationCheckFinished),
		},
	}

	payloadBuf := new(bytes.Buffer)
	err := json.NewEncoder(payloadBuf).Encode(updateRequest)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("PATCH", url, payloadBuf)
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		return err
	}

	_, err = api.httpClient.Do(req)
	if err != nil {
		return err
	}

	return nil
}

type dealUpdateValidationCheckDoneRequestProperties struct {
	ValidationCheckFinished string `json:"validation_check_finished"`
}

type dealUpdateValidationCheckDoneRequest struct {
	Properties dealUpdateValidationCheckDoneRequestProperties `json:"properties"`
}

// UpdateDealFlowCardValidationStatus updates the deal flow card attached to the given id with the given information
func (api HubspotDealFlowAPI) UpdateDealFlowCardValidationStatus(
	dealId string,
	dealValidationCheckFinished bool,
) error {

	log.Infof("Updating a deal flow card")

	url := fmt.Sprintf("https://api.hubapi.com/crm/v3/objects/deals/%s?hapikey=%s", dealId, api.APIKey)

	strconv.FormatBool(dealValidationCheckFinished)

	updateRequest := dealUpdateValidationCheckDoneRequest{
		dealUpdateValidationCheckDoneRequestProperties{
			strconv.FormatBool(dealValidationCheckFinished),
		},
	}

	payloadBuf := new(bytes.Buffer)
	err := json.NewEncoder(payloadBuf).Encode(updateRequest)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("PATCH", url, payloadBuf)
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		return err
	}

	_, err = api.httpClient.Do(req)
	if err != nil {
		return err
	}

	return nil
}
