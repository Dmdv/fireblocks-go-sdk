package fireblocksdk

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/pkg/errors"
)

type PostRequestOption struct {
	idempotencyKey string
}

type GetRequestOption struct {
	paged bool
}

func WithIdempotencyKey(idempotencyKey string) func(*PostRequestOption) {
	return func(option *PostRequestOption) {
		option.idempotencyKey = idempotencyKey
	}
}

func WithPaged(paged bool) func(*GetRequestOption) {
	return func(option *GetRequestOption) {
		option.paged = paged
	}
}

type SDKOptions struct {
	HTTPTimeoutMilliseconds time.Duration
	tokenExpirySeconds      int64
	auth                    IAuthProvider
	client                  IAPIClient
}

type FireblocksSDK struct {
	baseURL string
	client  IAPIClient
	auth    IAuthProvider
}

func WithAuthProvider(auth IAuthProvider) func(o *SDKOptions) {
	return func(o *SDKOptions) {
		o.auth = auth
	}
}

func WithAPIClient(client IAPIClient) func(o *SDKOptions) {
	return func(o *SDKOptions) {
		o.client = client
	}
}

func WithHTTPTimout(timeout time.Duration) func(o *SDKOptions) {
	return func(o *SDKOptions) {
		o.HTTPTimeoutMilliseconds = timeout
	}
}

func WithTokenTimeout(exp int64) func(o *SDKOptions) {
	return func(o *SDKOptions) {
		o.tokenExpirySeconds = exp
	}
}

func CreateSDK(apikey string, privateKey []byte, baseURL string, opts ...func(o *SDKOptions)) (*FireblocksSDK, error) {
	opt := &SDKOptions{tokenExpirySeconds: DefaultTokenExpiry()}

	for _, o := range opts {
		o(opt)
	}

	if opt.auth == nil {
		provider, err := NewAuthProvider(apikey, privateKey, WithTokenExpiry(opt.tokenExpirySeconds))
		if err != nil {
			return nil, err
		}

		opt.auth = provider
	}

	if opt.client == nil {
		opt.client = NewAPIClient(opt.auth, baseURL)
	}

	sdk := &FireblocksSDK{
		baseURL: baseURL,
		client:  opt.client,
		auth:    opt.auth,
	}

	return sdk, nil
}

// GetSupportedAssets Gets all assets that are currently supported by Fireblocks
func (sdk *FireblocksSDK) GetSupportedAssets() (resp []*AssetTypeResponse, err error) {
	body, status, err := sdk.client.DoGetRequest("/supported_assets", nil)
	if err == nil && status == http.StatusOK {
		err = json.Unmarshal(body, &resp)
		return
	}

	return resp, errors.Wrap(err, "failed to make request")
}
