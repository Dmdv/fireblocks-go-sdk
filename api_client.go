package fireblocksdk

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/pkg/errors"
)

const (
	APIVERSION string = "v1"
)

type IAPIClient interface {
	DoPostRequest(path string, body interface{}) ([]byte, int, error)
	DoGetRequest(path string) ([]byte, int, error)
	DoPutRequest(path string, body interface{}) ([]byte, int, error)
	DoDeleteRequest(path string) ([]byte, int, error)
}

type APIClient struct {
	httpClient *retryablehttp.Client
	auth       IAuthProvider
	baseURL    string
}

func NewAPIClient(auth IAuthProvider, baseURL string) *APIClient {
	client := retryablehttp.NewClient()

	return &APIClient{client, auth, baseURL}
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
		return nil, status, errors.Wrapf(err, "failed to do request: %s", path)
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

	var result []byte
	if resp != nil && resp.Body != nil {
		//goland:noinspection GoUnhandledErrorResult
		defer resp.Body.Close()

		responseBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return result, status, errors.Wrap(err, "failed to parse response body")
		}

		result = responseBody
	}

	return result, status, nil
}

func (api *APIClient) DoPostRequest(path string, body interface{}) ([]byte, int, error) {
	path = GetPath(path)
	return api.makeRequest(http.MethodPost, path, body)
}

func (api *APIClient) DoGetRequest(path string) ([]byte, int, error) {
	query := ""
	path = GetPath(path)
	return api.makeRequest(http.MethodGet, path, []byte(query))
}

func (api *APIClient) DoPutRequest(path string, body interface{}) ([]byte, int, error) {
	path = GetPath(path)
	return api.makeRequest(http.MethodPut, path, body)
}

func (api *APIClient) DoDeleteRequest(path string) ([]byte, int, error) {
	path = GetPath(path)
	return api.makeRequest(http.MethodDelete, path, []byte(""))
}

func GetPath(path string) string {
	return fmt.Sprintf(`/%s%s`, APIVERSION, path)
}
