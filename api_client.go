package fireblocksdk

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pkg/errors"

	"github.com/hashicorp/go-retryablehttp"
)

type IAPIClient interface {
	DoPostRequest(path string, body interface{}) (*http.Response, error)
	DoGetRequest(path string, body interface{}) (*http.Response, error)
	DoPutRequest(path string, body interface{}) (*http.Response, error)
	DoDeleteRequest(path string, body interface{}) (*http.Response, error)
}

type APIClient struct {
	client *retryablehttp.Client
	auth   IAuthProvider
}

func NewAPIClient(auth IAuthProvider) *APIClient {
	client := retryablehttp.NewClient()

	return &APIClient{client, auth}
}

func (client *APIClient) DoPostRequest(path string, body interface{}) (*http.Response, error) {
	bodyJson, err := json.Marshal(body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal body")
	}

	jwtToken, err := client.auth.SignJwt(path, bodyJson)
	if err != nil {
		return nil, errors.Wrap(err, "failed to sign request")
	}

	req, _ := http.NewRequest("POST", path, bytes.NewBuffer(bodyJson))
	req.Header.Add("X-API-Key", client.auth.GetApiKey())
	req.Header.Add("Authorization", fmt.Sprintf(`Bearer %s`, jwtToken))
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")

	return nil, nil
}

func (client *APIClient) DoGetRequest(path string, body interface{}) (*http.Response, error) {
	return nil, nil
}

func (client *APIClient) DoPutRequest(path string, body interface{}) (*http.Response, error) {
	return nil, nil
}

func (client *APIClient) DoDeleteRequest(path string, body interface{}) (*http.Response, error) {
	return nil, nil
}
