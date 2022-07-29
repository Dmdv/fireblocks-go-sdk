package fireblocksdk_test

import (
	sdk "fireblocksdk"
	"testing"
	"time"

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
	sdk          *sdk.FireblocksSDK
}

//goland:noinspection ALL
func (suite *SDKSuite) SetupTest() {
	viper.AutomaticEnv()
	viper.SetConfigFile("env.yaml")
	_ = viper.ReadInConfig()

	suite.baseURL = "https://api.fireblocks.io"
	suite.apiKey = viper.GetString("apikey")
	suite.apiSecretKey = []byte(viper.GetString("apisecret"))

	fb, err := sdk.CreateSDK(
		suite.apiKey,
		suite.apiSecretKey,
		suite.baseURL,
		sdk.WithTokenTimeout(5),
	)
	require.NoError(suite.T(), err)

	suite.sdk = fb
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

func (suite *SDKSuite) TestAccounts() {
	time.Sleep(time.Millisecond * 500)
	accounts, err := suite.sdk.GetVaultAccounts(nil)
	require.NoError(suite.T(), err)
	require.NotNil(suite.T(), accounts)
	require.Greater(suite.T(), len(accounts), 0)
}

func (suite *SDKSuite) TestAccountsPrefixFilters() {
	q := &sdk.VaultAccountsFilter{
		NamePrefix: "vault",
	}

	time.Sleep(time.Millisecond * 500)

	accounts, err := suite.sdk.GetVaultAccounts(q)
	require.NoError(suite.T(), err)
	require.NotNil(suite.T(), accounts)
	require.Greater(suite.T(), len(accounts), 0)
}

func (suite *SDKSuite) TestAccountsSuffixFilters() {
	q := &sdk.VaultAccountsFilter{
		NameSuffix: "UTX46",
	}

	time.Sleep(time.Millisecond * 500)

	accounts, err := suite.sdk.GetVaultAccounts(q)
	require.NoError(suite.T(), err)
	require.NotNil(suite.T(), accounts)
	require.GreaterOrEqual(suite.T(), len(accounts), 3)
}

func (suite *SDKSuite) TestAccountsAssetIDFilters() {
	q := &sdk.VaultAccountsFilter{
		AssetID: "DAI_UNI_TEST",
	}

	time.Sleep(time.Millisecond * 500)

	accounts, err := suite.sdk.GetVaultAccounts(q)
	require.NoError(suite.T(), err)
	require.NotNil(suite.T(), accounts)
	require.GreaterOrEqual(suite.T(), len(accounts), 3)
}

func (suite *SDKSuite) TestAccountsTotalAmountFilters() {
	q := &sdk.VaultAccountsFilter{
		MinAmountThreshold: "6.0",
	}

	time.Sleep(time.Millisecond * 500)

	accounts, err := suite.sdk.GetVaultAccounts(q)
	require.NoError(suite.T(), err)
	require.NotNil(suite.T(), accounts)
	require.Greater(suite.T(), len(accounts), 0)
}

func (suite *SDKSuite) TestAccountsPagedFilters() {
	qfirst := &sdk.PagedVaultAccountsRequestFilters{
		Limit: 1,
	}

	time.Sleep(time.Millisecond * 500)

	accounts1, err := suite.sdk.GetVaultAccountsWithPageInfo(qfirst)
	require.NoError(suite.T(), err)
	require.NotNil(suite.T(), accounts1)
	require.Equal(suite.T(), 1, len(accounts1.Accounts))
	require.NotEmpty(suite.T(), accounts1.NextURL)

	qsecond := &sdk.PagedVaultAccountsRequestFilters{
		Limit: 1,
		After: accounts1.Paging.After,
	}

	time.Sleep(time.Millisecond * 500)

	accounts2, err := suite.sdk.GetVaultAccountsWithPageInfo(qsecond)
	require.NotEqual(suite.T(), accounts1.Accounts[0].ID, accounts2.Accounts[0].ID)
	require.NoError(suite.T(), err)
	require.NotNil(suite.T(), accounts2)
	require.Equal(suite.T(), 1, len(accounts2.Accounts))
	require.NotEmpty(suite.T(), accounts2.NextURL)

	require.Contains(suite.T(), accounts1.NextURL, accounts1.Paging.After)
	require.Contains(suite.T(), accounts2.PreviousURL, accounts1.Paging.Before)
	require.Equal(suite.T(), accounts1.Paging.After, accounts1.Paging.After)
	require.Equal(suite.T(), accounts1.Paging.Before, accounts1.Paging.Before)

	qthird := &sdk.PagedVaultAccountsRequestFilters{
		Limit:  1,
		Before: accounts2.Paging.Before,
	}

	time.Sleep(time.Millisecond * 500)

	accounts3, err := suite.sdk.GetVaultAccountsWithPageInfo(qthird)
	require.NoError(suite.T(), err)
	id := accounts3.Accounts[0].ID
	actual := accounts1.Accounts[0].ID
	require.Equal(suite.T(), id, actual)
}

func (suite *SDKSuite) TestGetAccountsByID() {
	time.Sleep(time.Millisecond * 100)

	accounts, err := suite.sdk.GetVaultAccountsByID("256")
	require.NoError(suite.T(), err)
	require.NotNil(suite.T(), accounts)
	require.Equal(suite.T(), "256", accounts.ID)
}

func (suite *SDKSuite) TestGetAccountAsset() {
	time.Sleep(time.Millisecond * 100)

	asset, err := suite.sdk.GetVaultAccountAsset("60", "DAI_UNI_TEST")
	require.NoError(suite.T(), err)
	require.NotNil(suite.T(), asset)

	require.Equal(suite.T(), asset.ID, "DAI_UNI_TEST")
}
