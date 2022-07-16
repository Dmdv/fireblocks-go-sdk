package fireblocksdk_test

import (
	sdk "fireblocksdk"
	"testing"

	"github.com/spf13/viper"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

func TestSDKSuite(t *testing.T) {
	suite.Run(t, new(SDKSuite))
}

type SDKSuite struct {
	suite.Suite
	apiKey       string
	apiSecretKey []byte
	baseURL      string
}

//goland:noinspection ALL
func (suite *SDKSuite) SetupTest() {
	viper.AutomaticEnv()
	viper.SetConfigFile("env.yaml")
	viper.ReadInConfig()

	suite.baseURL = "https://api.fireblocks.io"
	suite.apiKey = viper.GetString("apikey")
	suite.apiSecretKey = []byte(viper.GetString("apisecret"))
}

func (suite *SDKSuite) TestSDKSuite() {
	fb, err := sdk.CreateSDK(
		"apiKey",
		[]byte(privateKey),
		suite.baseURL,
		sdk.WithHTTPTimout(111),
	)

	require.NoError(suite.T(), err)
	require.NotNil(suite.T(), fb)
}

func (suite *SDKSuite) TestSupportedAssets() {
	fb, err := sdk.CreateSDK(
		suite.apiKey,
		suite.apiSecretKey,
		suite.baseURL,
		sdk.WithTokenTimeout(5),
	)
	require.NoError(suite.T(), err)

	assets, err := fb.GetSupportedAssets()
	require.NoError(suite.T(), err)
	require.NotNil(suite.T(), assets)
	require.Greater(suite.T(), len(assets), 0)
}
