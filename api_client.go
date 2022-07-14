package fireblocksdk

import (
	"net/http"

	"github.com/hashicorp/go-retryablehttp"
)

type IAPIClient interface {
	DoPostRequest(req *http.Request) (*http.Response, error)
	DoGetRequest(req *http.Request) (*http.Response, error)
	DoPutRequest(req *http.Request) (*http.Response, error)
	DoDeleteRequest(req *http.Request) (*http.Response, error)
}

type APIClient struct {
	client *retryablehttp.Client
}

func NewAPIClient() *APIClient {
	client := retryablehttp.NewClient()

	return &APIClient{client}
}

func (client *APIClient) DoPostRequest(req *http.Request) (*http.Response, error) {
	return nil, nil
}

func (client *APIClient) DoGetRequest(req *http.Request) (*http.Response, error) {
	return nil, nil
}

func (client *APIClient) DoPutRequest(req *http.Request) (*http.Response, error) {
	return nil, nil
}

func (client *APIClient) DoDeleteRequest(req *http.Request) (*http.Response, error) {
	return nil, nil
}
