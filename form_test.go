package go_hubspot

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// getFormAPI Get default HubSpot Form API client
func getFormAPI() HubspotFormAPI {
	return NewHubspotFormAPI("form", "key")
}

func getMockFormAPI(mockClient *IHTTPClientMock) HubspotFormAPI {
	return HubspotFormAPI{
		URLTemplate: "",
		FormID:      "",
		APIKey:      "api_key",
		httpClient:  mockClient,
	}
}

func TestGetSubmissionMap(t *testing.T) {
	submission := Submission{
		SubmittedAt: 0,
		Values: []FormValue{
			{
				Name:  "name1",
				Value: "value1",
			},
			{
				Name:  "name2",
				Value: "value2",
			},
		},
	}

	expected := map[string]string{
		"name1": "value1",
		"name2": "value2",
	}

	got := GetSubmissionMap(submission)

	for key := range expected {
		if got[key] != expected[key] {
			t.Errorf("Expected: %#v, got: %#v", expected, got)
			break
		}
	}

	for key := range got {
		if got[key] != expected[key] {
			t.Errorf("Expected: %#v, got: %#v", expected, got)
			break
		}
	}
}

func TestGetPageURL(t *testing.T) {
	expected := "https://api.hubapi.com/form-integrations/v1/submissions/forms/form?hapikey=key&limit=50&after=some-application-id"
	api := getFormAPI()
	got := api.GetPageURL("some-application-id")
	if expected != got {
		t.Errorf("Expected: %#v, got: %#v", expected, got)
	}
}

func TestGetByApplicationID(t *testing.T) {
	response := HubspotResponse{
		Results: []Submission{
			{
				SubmittedAt: 0,
				Values: []FormValue{
					{
						Name:  "application_id",
						Value: "some-application_id",
					},
					{
						Name:  "other",
						Value: "other-value",
					},
				},
			},
			{
				SubmittedAt: 0,
				Values: []FormValue{
					{
						Name:  "application_id",
						Value: "other-application_id",
					},
					{
						Name:  "other",
						Value: "other-value",
					},
				},
			},
		},
		Paging: nil,
	}

	_, err := response.GetByKeyValue("application_id", "application_id")

	if err == nil || err.Error() != "Submission with application_id `application_id` not found" {
		t.Errorf("Expected application_id not found")
	}

	submission, err := response.GetByKeyValue("application_id", "other-application_id")

	if err != nil || submission == nil {
		t.Errorf("Expected other-application_id found")
	}
}

func TestGetNextAfter(t *testing.T) {
	hasNext := HubspotResponse{
		Results: []Submission{},
		Paging: &Paging{
			Next: map[string]string{
				"after": "page-id",
			},
		},
	}

	after, err := hasNext.GetNextAfter()
	if err != nil || after != "page-id" {
		t.Errorf("Expected next `after` value to be `page-id`")
	}

	hasNoNext := HubspotResponse{
		Results: []Submission{},
		Paging:  nil,
	}

	after, err = hasNoNext.GetNextAfter()
	if err == nil || after != "" {
		t.Errorf("Expected next `after` to not exist")
	}

}

func TestSearchForApplicationID(t *testing.T) {
	mockHubspotHTTPClient := IHTTPClientMock{
		DoFunc: func(req *http.Request) (resp *http.Response, err error) { return nil, nil },
		GetFunc: func(url string) (resp *http.Response, err error) {
			w := httptest.NewRecorder()
			if url == "https://example.com/form_id?hapikey=api_key&limit=50&after=" {
				w.WriteHeader(200)
				w.Write([]byte(`
				{
					"results": [
						{
							"submittedAt":1611226634790,
							"values":[
								{"name":"company","value":"company1","objectTypeId":"0-1"},
								{"name":"application_id","value":"application_id1","objectTypeId":"0-1"},
								{"name":"company_number","value":"1","objectTypeId":"0-1"}
							]
						},
						{
							"submittedAt":1611226634790,
							"values":[
								{"name":"company","value":"company11","objectTypeId":"0-1"},
								{"name":"application_id","value":"application_id11","objectTypeId":"0-1"},
								{"name":"company_number","value":"11","objectTypeId":"0-1"}
							]
						}
					],
					"paging": {
						"next": {
							"after": "first"
						}
					}
				}
				`))
			} else if url == "https://example.com/form_id?hapikey=api_key&limit=50&after=first" {
				w.WriteHeader(200)
				w.Write([]byte(`
				{
					"results": [
						{
							"submittedAt":1611226634790,
							"values":[
								{"name":"company","value":"company22","objectTypeId":"0-1"},
								{"name":"application_id","value":"application_id22","objectTypeId":"0-1"},
								{"name":"company_number","value":"22","objectTypeId":"0-1"}
							]
						},
						{
							"submittedAt":1611226634790,
							"values":[
								{"name":"company","value":"company2","objectTypeId":"0-1"},
								{"name":"application_id","value":"application_id2","objectTypeId":"0-1"},
								{"name":"company_number","value":"2","objectTypeId":"0-1"}
							]
						}
					],
					"paging": {
						"next": {
							"after": "second"
						}
					}
				}
				`))
			} else if url == "https://example.com/form_id?hapikey=api_key&limit=50&after=second" {
				w.WriteHeader(200)
				w.Write([]byte(`
				{
					"results": []
				}
				`))
			}

			return w.Result(), nil
		},
	}

	api := HubspotFormAPI{
		URLTemplate: "https://example.com/%s?hapikey=%s&limit=50&after=%s",
		FormID:      "form_id",
		APIKey:      "api_key",
		httpClient:  &mockHubspotHTTPClient,
	}

	form, err := api.SearchForKeyValue("application_id", "application_id1")

	if err != nil || form["application_id"] != "application_id1" || form["company"] != "company1" {
		t.Errorf("Expected to find form with application_id1 on page 1")
	}

	if len(mockHubspotHTTPClient.GetCalls()) != 1 {
		t.Errorf("Expected 1 call to HubSpot API")
	}

	form, err = api.SearchForKeyValue("application_id", "application_id2")

	if err != nil || form["application_id"] != "application_id2" || form["company"] != "company2" {
		t.Errorf("Expected to find form with application_id2 on page 2")
	}

	if len(mockHubspotHTTPClient.GetCalls()) != 3 {
		t.Errorf("Expected 2 call to HubSpot API")
	}

	form, err = api.SearchForKeyValue("application_id", "none")

	if err == nil {
		t.Errorf("Expected to not find form with application_id=none")
	}

	if len(mockHubspotHTTPClient.GetCalls()) != 6 {
		t.Errorf("Expected 3 call to HubSpot API")
	}

}
