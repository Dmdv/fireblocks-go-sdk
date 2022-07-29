package fireblocksdk

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/pkg/errors"
)

const (
	APIVERSION string = "v1"
)

type IAPIClient interface {
	DoPostRequest(path string, body interface{}) ([]byte, int, error)
	DoGetRequest(path string, q url.Values) ([]byte, int, error)
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
	var (
		status   = http.StatusInternalServerError
		bodyJSON = []byte("")
	)

	if method != http.MethodGet && body != nil {
		var err error
		bodyJSON, err = json.Marshal(body)
		if err != nil {
			return nil, status, errors.Wrap(err, "failed to marshal body")
		}
	}

	jwtToken, err := api.auth.SignJwt(path, bodyJSON)
	if err != nil {
		return nil, status, errors.Wrap(err, "failed to sign request")
	}

	path = fmt.Sprintf("%s%s", api.baseURL, path)

	req, err := retryablehttp.NewRequest(method, path, prepareBody(bodyJSON))
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

	var result []byte
	if resp != nil && resp.Body != nil {
		//goland:noinspection GoUnhandledErrorResult
		defer resp.Body.Close()

		status = resp.StatusCode
		log.Printf(
			"[%s] %s: status: %s, code: %d, content-length: %d",
			method,
			path,
			resp.Status,
			resp.StatusCode,
			resp.ContentLength,
		)

		responseBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return result, status, errors.Wrap(err, "failed to parse response body")
		}

		result = responseBody
	}

	return result, status, nil
}

func (api *APIClient) DoPostRequest(path string, body interface{}) ([]byte, int, error) {
	path = api.GetRelativePath(path)

	return api.makeRequest(http.MethodPost, path, body)
}

func (api *APIClient) DoGetRequest(path string, q url.Values) ([]byte, int, error) {
	query := ""
	path = api.GetRelativePath(path)

	if q != nil {
		query = q.Encode()
		path = fmt.Sprintf(`%s?%s`, path, query)
	}

	return api.makeRequest(http.MethodGet, path, []byte(query))
}

func (api *APIClient) DoPutRequest(path string, body interface{}) ([]byte, int, error) {
	path = api.GetRelativePath(path)

	return api.makeRequest(http.MethodPut, path, body)
}

func (api *APIClient) DoDeleteRequest(path string) ([]byte, int, error) {
	path = api.GetRelativePath(path)

	return api.makeRequest(http.MethodDelete, path, nil)
}

// GetRelativePath returns path without baseURL
func (api *APIClient) GetRelativePath(path string) string {
	return fmt.Sprintf(`/%s%s`, APIVERSION, path)
}

func prepareBody(encodedBody []byte) io.ReadCloser {
	if string(encodedBody) == "{}" {
		encodedBody = []byte("")
	}

	return ioutil.NopCloser(
		strings.NewReader(
			string(encodedBody),
		),
	)
}
