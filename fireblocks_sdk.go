package fireblocksdk

import (
	"time"
)

type SDKOptions struct {
	timeoutMS time.Duration
	auth      IAuthProvider
	client    IAPIClient
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

	if opt.client == nil {
		opt.client = NewAPIClient(opt.auth)
	}

	sdk := &FireblocksSDK{
		baseURL: baseURL,
		client:  opt.client,
		auth:    opt.auth,
	}

	return sdk, nil
}
