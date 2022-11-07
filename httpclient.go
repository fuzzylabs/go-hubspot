package go_hubspot

import "net/http"

// IHTTPClient HTTP client interface to be used with HubSpot API clients
type IHTTPClient interface {
	Do(req *http.Request) (resp *http.Response, err error)
}

// HTTPClient HTTP client to be used with HubSpot API clients
type HTTPClient struct{}

// Do performs a request to a url
func (c HTTPClient) Do(req *http.Request) (resp *http.Response, err error) {
	client := http.Client{}
	return client.Do(req)
}
