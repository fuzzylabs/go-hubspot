package go_hubspot

import (
	"bytes"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

type IHubspotDealFlowAPI interface {
	AssociateDealFlowCard(dealId, assocId, objectType, assocType string) error
	CreateDealFlowCard(
		cardName string,
		contactID string,
		contactAssocType string,
		companyID string,
		stageName string,
		pipeline string,
		ownerId string,
		otherProperties map[string]string,
	) (*DealCreationResponse, error)
	UpdateDealFlowCard(
		dealId string,
		properties map[string]string,
	) error
}

type HubspotDealFlowAPI struct {
	APIKey     string
	httpClient IHTTPClient
}

// dealCreationRequest is a representation of the deal creation request to HubSpot
type dealCreationRequest struct {
	Properties map[string]string `json:"properties"`
}

// DealCreationResponseProperties is a representation of the deal creation response from HubSpot
type DealCreationResponseProperties struct {
	Amount             string `json:"amount"`
	CloseDate          string `json:"closedate"`
	CreateDate         string `json:"createdate"`
	DealName           string `json:"dealname"`
	DealStage          string `json:"dealstage"`
	HsLastModifiedDate string `json:"hs_lastmodifieddate"`
	HubspotOwnerId     string `json:"hubspot_owner_id"`
	Pipeline           string `json:"pipeline"`
}

// DealCreationResponse is a representation of the deal creation response from HubSpot
type DealCreationResponse struct {
	Id         string                         `json:"id"`
	Properties DealCreationResponseProperties `json:"properties"`
	CreatedAt  string                         `json:"createdAt"`
	UpdatedAt  string                         `json:"updatedAt"`
	Archived   bool                           `json:"archived"`
}

type dealUpdateRequest struct {
	Properties map[string]string `json:"properties"`
}

type CardAssociation int64

type DealAssociationFromTo struct {
	Id string `json:"id"`
}

type DealAssociation struct {
	From DealAssociationFromTo `json:"from"`
	To   DealAssociationFromTo `json:"to"`
	Type string                `json:"type"`
}

type DealAssociationBatchRequest struct {
	Inputs []DealAssociation `json:"inputs"`
}

// NewHubspotDealFlowAPI creates new HubspotDealFlowAPI with form ID and API key
func NewHubspotDealFlowAPI(apiKey string) HubspotDealFlowAPI {
	return HubspotDealFlowAPI{
		APIKey:     apiKey,
		httpClient: HTTPClient{},
	}
}

// AssociateDealFlowCard associates a deal flow card with a company or contact using the internal HubSpot dealId and companyId/contactId
// Choose whether to associate a company or contact by setting assocType to "contact" or "company"
func (api HubspotDealFlowAPI) AssociateDealFlowCard(dealId, assocId, objectType, assocType string) error {
	url := fmt.Sprintf("https://api.hubapi.com/crm/v3/associations/deal/%s/batch/create?hapikey=%s",
		objectType,
		api.APIKey,
	)

	associationRequest := DealAssociationBatchRequest{
		Inputs: []DealAssociation{
			{
				From: DealAssociationFromTo{Id: dealId},
				To:   DealAssociationFromTo{Id: assocId},
				Type: assocType,
			},
		},
	}

	payloadBuf := new(bytes.Buffer)
	err := json.NewEncoder(payloadBuf).Encode(associationRequest)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, payloadBuf)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	_, err = api.httpClient.Do(req)
	if err != nil {
		return err
	}

	return nil
}

// CreateDealFlowCard creates a deal flow card with the given parameters in HubSpot,
// and associates it with a company and contact
func (api HubspotDealFlowAPI) CreateDealFlowCard(
	cardName string,
	contactID string,
	contactAssocType string,
	companyID string,
	stageName string,
	pipeline string,
	ownerId string,
	otherProperties map[string]string,
) (*DealCreationResponse, error) {

	log.Infof("Creating a deal flow card")

	url := fmt.Sprintf("https://api.hubapi.com/crm/v3/objects/deals?hapikey=%s", api.APIKey)

	creationRequest := dealCreationRequest{
		map[string]string{
			"dealname":         cardName,
			"dealstage":        stageName,
			"pipeline":         pipeline,
			"hubspot_owner_id": ownerId,
		},
	}

	for key, value := range otherProperties {
		creationRequest.Properties[key] = value
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
	err = api.AssociateDealFlowCard(hubspotResp.Id, companyID, "company", "deal_to_company")
	if err != nil {
		return nil, err
	}

	// Associate the deal with a contact based on the application id
	err = api.AssociateDealFlowCard(hubspotResp.Id, contactID, "contact", contactAssocType)
	if err != nil {
		return nil, err
	}

	return &hubspotResp, nil

}

// UpdateDealFlowCard updates the deal flow card attached to the given id with the given information
func (api HubspotDealFlowAPI) UpdateDealFlowCard(
	dealId string,
	properties map[string]string,
) error {

	log.Infof("Updating a deal flow card")

	url := fmt.Sprintf("https://api.hubapi.com/crm/v3/objects/deals/%s?hapikey=%s", dealId, api.APIKey)

	updateRequest := dealUpdateRequest{
		properties,
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
