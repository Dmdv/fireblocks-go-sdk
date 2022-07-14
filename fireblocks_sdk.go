package fireblocksdk

import (
	"time"

	"github.com/hashicorp/go-retryablehttp"
)

type SDKOptions struct {
	timeoutMS time.Duration
	auth      IAuthProvider
}

type FireblocksSDK struct {
	baseURL string
	client  *retryablehttp.Client
	auth    IAuthProvider
}

func WithAuthProvider(auth IAuthProvider) func(o *SDKOptions) {
	return func(o *SDKOptions) {
		o.auth = auth
	}
}

func WithTimout(timeout time.Duration) func(o *SDKOptions) {
	return func(o *SDKOptions) {
		o.timeoutMS = timeout
	}
}

func CreateSDK(apikey, privateKey string, baseURL string, opts ...func(o *SDKOptions)) (*FireblocksSDK, error) {
	opt := &SDKOptions{}

	for _, o := range opts {
		o(opt)
	}

	if opt.auth == nil {
		provider, err := NewAuthProvider(apikey, privateKey)
		if err != nil {
			return nil, err
		}

		opt.auth = provider
	}

	client := retryablehttp.NewClient()

	sdk := &FireblocksSDK{
		client:  client,
		baseURL: baseURL,
		auth:    opt.auth,
	}

	return sdk, nil
}
