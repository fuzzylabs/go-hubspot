package ehe_hubspot

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"
)

// getDealFlowAPI Get default HubSpot DealFlow API client
func getDealFlowAPI() HubspotDealFlowAPI {
	return NewHubspotDealFlowAPI("key")
}

func getMockDealFlowAPI(mockClient *IHTTPClientMock) HubspotDealFlowAPI {
	return HubspotDealFlowAPI{
		APIKey:     "api_key",
		httpClient: mockClient,
	}
}

func TestCreateDealFlowCard(t *testing.T) {

	expected := dealCreationResponse{
		"dealid",
		dealCreationResponseProperties{
			"Amount",
			"CloseDate",
			"CreateDate",
			"DealName",
			"DealStage",
			"HsLastModifiedDate",
			"HubspotOwnerId",
			"Pipeline",
			"false",
		},
		"CreatedAt",
		"UpdatedAt",
		false,
	}

	// Set to true when the correct association call is made
	companyAssociation := false
	contactAssociation := false

	mockHubspotHTTPClient := IHTTPClientMock{
		GetFunc: func(url string) (resp *http.Response, err error) { return nil, nil },
		DoFunc: func(req *http.Request) (resp *http.Response, err error) {
			url := fmt.Sprintf("%s", req.URL)

			w := httptest.NewRecorder()
			if url == "https://api.hubapi.com/crm/v3/objects/deals?hapikey=api_key" {
				// This is a deal flow creation call
				// Test the body

				expectedRequest := dealCreationRequest{
					dealCreationRequestProperties{
						"cardName",
						"stageName",
						"pipeline",
						"applicationId",
						"HubspotOwnerId",
						"false",
					},
				}

				body, err := ioutil.ReadAll(req.Body)
				defer func(Body io.ReadCloser) {
					err := Body.Close()
					if err != nil {
						t.Errorf("Error closing mock request body: %s", err.Error())
					}
				}(req.Body)

				if err != nil {
					t.Errorf("Error reading CreateDealFlowCard request body: %s", err.Error())
				}

				var request dealCreationRequest
				err = json.Unmarshal(body, &request)
				if err != nil {
					t.Errorf("Error unmarshalling CreateDealFlowCard request: %s", err.Error())
				}

				if !reflect.DeepEqual(expectedRequest, request) {
					t.Errorf("Unexpected CreateDealFlowCard request, expected:\n%s\ngot:\n%s", expectedRequest, request)
				}

				// Send the expected response
				response, err := json.Marshal(expected)
				if err != nil {
					t.Errorf("Error marshalling expected response: %s", err.Error())
				}

				w.WriteHeader(200)
				_, err = w.Write(response)
				if err != nil {
					t.Errorf("Error writing response in mock: %s", err.Error())
				}
			} else if url == "https://api.hubapi.com/crm/v3/objects/contacts/search?hapikey=api_key" {
				// Searching for contact
				searchResponse := HubSpotContactSearchResponse{
					1,
					[]ContactResult{
						{
							"contactid",
							map[string]string{
								"company_number": "companyNumber",
							},
							map[string]associations{},
						},
					},
				}

				// Send the expected response
				response, err := json.Marshal(searchResponse)
				if err != nil {
					t.Errorf("Error marshalling search response: %s", err.Error())
				}

				w.WriteHeader(200)
				_, err = w.Write(response)
				if err != nil {
					t.Errorf("Error writing response in mock: %s", err.Error())
				}
			} else if url == "https://api.hubapi.com/crm/v3/objects/contacts/contactid?associations=company&archived=false&hapikey=api_key" {
				// Searching for company
				searchResponse := ContactResult{
					"contactid",
					map[string]string{
						"company_number": "companyNumber",
					},
					map[string]associations{
						"companies": {
							[]association{
								{
									"companyId",
									"type",
								},
							},
						},
					},
				}

				// Send the expected response
				response, err := json.Marshal(searchResponse)
				if err != nil {
					t.Errorf("Error marshalling search response: %s", err.Error())
				}

				_, err = w.Write(response)
				if err != nil {
					t.Errorf("Error writing response in mock: %s", err.Error())
				}
			} else if url == "https://api.hubapi.com/crm/v3/objects/deals/dealid/associations/company/companyId/deal_to_company?hapikey=api_key" {
				// Created an association between the deal and correct company
				if req.Method != "PUT" {
					t.Errorf("Deal association used %s, instead of PUT", req.Method)
				}

				companyAssociation = true
				w.WriteHeader(200)
			} else if url == "https://api.hubapi.com/crm/v3/objects/deals/dealid/associations/contact/contactId/deal_to_contact?hapikey=api_key" {
				// Created an association between the deal and correct company
				if req.Method != "PUT" {
					t.Errorf("Deal association used %s, instead of PUT", req.Method)
				}

				contactAssociation = true
				w.WriteHeader(200)
			} else {
				t.Errorf("Unexpected url %s", url)
			}

			return w.Result(), nil
		},
	}

	api := getMockDealFlowAPI(&mockHubspotHTTPClient)

	// Set environment variables required for CreateDealFlowCard
	err := os.Setenv("DEALFLOW_STARTING_STAGE", "stageName")
	if err != nil {
		t.Errorf("Error setting environment variable DEALFLOW_STARTING_STAGE: %s", err.Error())
	}

	err = os.Setenv("DEALFLOW_PIPELINE_NAME", "pipeline")
	if err != nil {
		t.Errorf("Error setting environment variable DEALFLOW_PIPELINE_NAME: %s", err.Error())
	}

	err = os.Setenv("DEALFLOW_OWNER_ID", "HubspotOwnerId")
	if err != nil {
		t.Errorf("Error setting environment variable DEALFLOW_OWNER_ID: %s", err.Error())
	}

	response, err := api.CreateDealFlowCard("cardName", "contactId", "companyId", "applicationId")
	if err != nil {
		t.Errorf("CreateDealFlowCard returned an error: %s", err.Error())
		return
	}

	if !reflect.DeepEqual(expected, *response) {
		t.Errorf("CreateDealFlowCard returned incorrect response, expected:\n%v\ngot:\n%v", expected, response)
		return
	}

	if !companyAssociation {
		// The correct association PUT call was not made
		t.Errorf("Expected correct call to the association api for a company, did not receive it")
	}

	if !contactAssociation {
		// The correct association PUT call was not made
		t.Errorf("Expected correct call to the association api for a contact, did not receive it")
	}

	if len(mockHubspotHTTPClient.DoCalls()) != 3 {
		t.Errorf("Expected 3 calls to HubSpot API")
		return
	}
}

func TestAssociateDealFlowCard(t *testing.T) {
	mockHubspotHTTPClient := IHTTPClientMock{
		GetFunc: func(url string) (resp *http.Response, err error) { return nil, nil },
		DoFunc: func(req *http.Request) (resp *http.Response, err error) {
			url := fmt.Sprintf("%s", req.URL)
			w := httptest.NewRecorder()

			if url == "https://api.hubapi.com/crm/v3/objects/deals/dealid/associations/company/companyid/deal_to_company?hapikey=api_key" {
				// Created an association between the deal and correct company
				if req.Method != "PUT" {
					t.Errorf("Deal association used %s, instead of PUT", req.Method)
				}

				w.WriteHeader(200)
			} else if url == "https://api.hubapi.com/crm/v3/objects/deals/dealid/associations/contact/contactid/deal_to_contact?hapikey=api_key" {
				// Created an association between the deal and correct company
				if req.Method != "PUT" {
					t.Errorf("Deal association used %s, instead of PUT", req.Method)
				}

				w.WriteHeader(200)
			} else {
				t.Errorf("Unexpected url %s", url)
			}

			return w.Result(), nil
		},
	}

	api := getMockDealFlowAPI(&mockHubspotHTTPClient)

	err := api.AssociateDealFlowCard("dealid", "companyid", Company)

	if err != nil {
		t.Errorf("Error on AssociateDealFlowCard: %s", err.Error())
	}

	err = api.AssociateDealFlowCard("dealid", "contactid", Contact)

	if err != nil {
		t.Errorf("Error on AssociateDealFlowCard: %s", err.Error())
	}

	if len(mockHubspotHTTPClient.DoCalls()) != 2 {
		t.Errorf("Expected 2 calls to HubSpot API")
		return
	}
}

func TestUpdateDealFlowCard(t *testing.T) {
	mockHubspotHTTPClient := IHTTPClientMock{
		GetFunc: func(url string) (resp *http.Response, err error) { return nil, nil },
		DoFunc: func(req *http.Request) (resp *http.Response, err error) {
			url := fmt.Sprintf("%s", req.URL)

			w := httptest.NewRecorder()
			if url == "https://api.hubapi.com/crm/v3/objects/deals/dealid?hapikey=api_key" {
				// This is a deal flow creation call
				// Test the body

				if req.Method != "PATCH" {
					t.Errorf("UpdateDealFlowCardValidationStatus used incorrect request method, expected: PATCH, got: %s", req.Method)
				}

				expectedRequest := dealUpdateRequest{
					dealUpdateRequestProperties{
						"dealName",
						"stageName",
						"applicationId",
						"false",
					},
				}

				body, err := ioutil.ReadAll(req.Body)
				defer func(Body io.ReadCloser) {
					err := Body.Close()
					if err != nil {
						t.Errorf("Error closing mock request body: %s", err.Error())
					}
				}(req.Body)

				if err != nil {
					t.Errorf("Error reading UpdateDealFlowCardValidationStatus request body: %s", err.Error())
				}

				var request dealUpdateRequest
				err = json.Unmarshal(body, &request)
				if err != nil {
					t.Errorf("Error unmarshalling UpdateDealFlowCardValidationStatus request: %s", err.Error())
				}

				if !reflect.DeepEqual(expectedRequest, request) {
					t.Errorf("Unexpected UpdateDealFlowCardValidationStatus request, expected:\n%s\ngot:\n%s", expectedRequest, request)
				}

				w.WriteHeader(200)

				// Request normally responds with json about updated deal, this is not used in
				// UpdateDealFlowCardValidationStatus, so it is omitted from the test
				// See: https://developers.hubspot.com/docs/api/crm/deals
			} else {
				t.Errorf("Unexpected url %s", url)
			}

			return w.Result(), nil
		},
	}

	api := getMockDealFlowAPI(&mockHubspotHTTPClient)

	err := api.UpdateDealFlowCard("dealid", "dealName", "stageName", "applicationId", false)
	if err != nil {
		t.Errorf("Error on UpdateDealFlowCardValidationStatus: %s", err.Error())
	}

	if len(mockHubspotHTTPClient.DoCalls()) != 1 {
		t.Errorf("Expected 1 call to HubSpot API")
		return
	}

}
