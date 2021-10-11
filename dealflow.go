package ehe_hubspot

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"
)

// dealCreationRequestProperties is a representation of the deal creation request to HubSpot
type dealCreationRequestProperties struct {
	DealName       string `json:"dealname"`
	DealStage      string `json:"dealstage"`
	Pipeline       string `json:"pipeline"`
	ApplicationId  string `json:"application_id"`
	HubspotOwnerId string `json:"hubspot_owner_id"`
}

// dealCreationRequest is a representation of the deal creation request to HubSpot
type dealCreationRequest struct {
	Properties dealCreationRequestProperties `json:"properties"`
}

// dealCreationResponseProperties is a representation of the deal creation response from HubSpot
type dealCreationResponseProperties struct {
	Amount             string `json:"amount"`
	CloseDate          string `json:"closedate"`
	CreateDate         string `json:"createdate"`
	DealName           string `json:"dealname"`
	DealStage          string `json:"dealstage"`
	HsLastModifiedDate string `json:"hs_lastmodifieddate"`
	HubspotOwnerId     string `json:"hubspot_owner_id"`
	Pipeline           string `json:"pipeline"`
}

// dealCreationResponse is a representation of the deal creation response from HubSpot
type dealCreationResponse struct {
	Id         string                         `json:"id"`
	Properties dealCreationResponseProperties `json:"properties"`
	CreatedAt  string                         `json:"createdAt"`
	UpdatedAt  string                         `json:"updatedAt"`
	Archived   bool                           `json:"archived"`
}

type dealUpdateRequestProperties struct {
	DealName      string `json:"dealname"`
	DealStage     string `json:"dealstage"`
	ApplicationId string `json:"uuid"`
}

type dealUpdateRequest struct {
	Properties dealUpdateRequestProperties `json:"properties"`
}

// AssociateDealFlowCardWithCompany associates a deal flow card with a company using the internal HubSpot dealId and companyId
func (api HubspotFormAPI) AssociateDealFlowCardWithCompany(dealId string, companyId string) error {
	url := fmt.Sprintf(
		"https://api.hubapi.com/crm/v3/objects/deals/%s/associations/company/%s/deal_to_company?hapikey=%s",
		dealId,
		companyId,
		api.APIKey,
	)

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
func (api HubspotFormAPI) CreateDealFlowCard(
	cardName string,
	companyID string,
	applicationId string,
) (*dealCreationResponse, error) {

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

	var hubspotResp dealCreationResponse

	log.Infof("Raw response: %s", string(body))
	err = json.Unmarshal(body, &hubspotResp)
	if err != nil {
		return nil, err
	}

	// Associate the deal with a company based on the application id
	err = api.AssociateDealFlowCardWithCompany(hubspotResp.Id, companyID)
	if err != nil {
		return nil, err
	}

	return &hubspotResp, nil

}

// UpdateDealFlowCard updates the deal flow card attached to the given id with the given information
func (api HubspotFormAPI) UpdateDealFlowCard(
	dealId string,
	dealName string,
	dealStage string,
	applicationId string,
) error {

	log.Infof("Updating a deal flow card")

	url := fmt.Sprintf("https://api.hubapi.com/crm/v3/objects/deals/%s?hapikey=%s", dealId, api.APIKey)

	updateRequest := dealUpdateRequest{
		dealUpdateRequestProperties{
			dealName,
			dealStage,
			applicationId,
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
