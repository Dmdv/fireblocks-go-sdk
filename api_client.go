package fireblocksdk

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/pkg/errors"

	"github.com/hashicorp/go-retryablehttp"
)

type IAPIClient interface {
	DoPostRequest(path string, body interface{}) ([]byte, int, error)
	DoGetRequest(path string, body interface{}) ([]byte, int, error)
	DoPutRequest(path string, body interface{}) ([]byte, int, error)
	DoDeleteRequest(path string, body interface{}) ([]byte, int, error)
}

type APIClient struct {
	httpClient *retryablehttp.Client
	auth       IAuthProvider
}

func NewAPIClient(auth IAuthProvider) *APIClient {
	client := retryablehttp.NewClient()

	return &APIClient{client, auth}
}

func (api *APIClient) makeRequest(method, path string, body interface{}) ([]byte, int, error) {
	var status = http.StatusInternalServerError
	bodyJson, err := json.Marshal(body)
	if err != nil {
		return nil, status, errors.Wrap(err, "failed to marshal body")
	}

	jwtToken, err := api.auth.SignJwt(path, bodyJson)
	if err != nil {
		return nil, status, errors.Wrap(err, "failed to sign request")
	}

	req, err := retryablehttp.NewRequest(method, path, bodyJson)
	if err != nil {
		return nil, status, errors.Wrap(err, "failed to create request")
	}

	req.Header.Add("X-API-Key", api.auth.GetApiKey())
	req.Header.Add("Authorization", fmt.Sprintf(`Bearer %s`, jwtToken))
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")

	resp, err := api.httpClient.Do(req)
	if err != nil {
		return nil, status, err
	}

	status = resp.StatusCode
	log.Printf(
		"[%s] Path %s: status: %s, code: %d, content-length: %d",
		method,
		path,
		resp.Status,
		resp.StatusCode,
		resp.ContentLength,
	)

	//goland:noinspection GoUnhandledErrorResult
	defer resp.Body.Close()

	var result []byte
	if resp != nil {
		responseBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return result, status, errors.Wrap(err, "failed to parse response body")
		}
		result = responseBody
	}

	return result, status, nil
}

func (api *APIClient) DoPostRequest(path string, body interface{}) ([]byte, int, error) {
	return api.makeRequest("POST", path, body)
}

func (api *APIClient) DoGetRequest(path string, body interface{}) ([]byte, int, error) {
	return api.makeRequest("GET", path, body)
}

func (api *APIClient) DoPutRequest(path string, body interface{}) ([]byte, int, error) {
	return api.makeRequest("PUT", path, body)
}

func (api *APIClient) DoDeleteRequest(path string, body interface{}) ([]byte, int, error) {
	return api.makeRequest("DELETE", path, body)
}
