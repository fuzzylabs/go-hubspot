package ehe_hubspot

//go:generate moq -out httpclient_mock_test.go . IHTTPClient
//go:generate moq -out crm_mock.go . IHubspotCRMAPI
//go:generate moq -out dealflow_mock.go . IHubspotDealFlowAPI
//go:generate moq -out form_mock.go . IHubspotFormAPI
