package fireblocksdk_test

import (
	sdk "fireblocksdk"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/suite"
)

func TestSDKSuite(t *testing.T) {
	suite.Run(t, new(SDKSuite))
}

type SDKSuite struct {
	suite.Suite
	apiKey       string
	apiSecretKey string
	baseURL      string
}

func (suite *SDKSuite) SetupTest() {
	suite.apiKey = "apiKey"
	suite.apiSecretKey = privateKey
	suite.baseURL = "api.fireblocks.io"
}

func (suite *SDKSuite) TestSDKSuite() {
	fb, err := sdk.CreateSDK(
		suite.apiKey,
		suite.apiSecretKey,
		suite.baseURL,
		sdk.WithTimout(111),
	)

	require.NoError(suite.T(), err)
	require.NotNil(suite.T(), fb)
}
