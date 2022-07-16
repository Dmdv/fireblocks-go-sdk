package fireblocksdk_test

import (
	sdk "fireblocksdk"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/suite"
)

var required = []string{
	"iat",
	"uri",
	"nonce",
	"now",
	"exp",
	"sub",
	"bodyHash",
}

type TestClaimsSuite struct {
	suite.Suite
	time   sdk.ITimeProvider
	claims sdk.IFireblocksClaims
}

func TestTestClaimsSuite(t *testing.T) {
	suite.Run(t, new(TestClaimsSuite))
}

func (suite *TestClaimsSuite) SetupTest() {
	suite.time = &testTimeProvider{}
	suite.claims = sdk.NewFireblocksClaimsProvider(suite.time, sdk.DefaultTokenExpiry())
}

func (suite *TestClaimsSuite) TestClaimsHaveAllFields() {
	claims := suite.claims.CreateClaims("path", "apikey", "hash")

	for _, field := range required {
		_, has := claims[field]
		require.True(suite.T(), has)
	}
}

func (suite *TestClaimsSuite) TestClaimsShouldBeOnTimeDependent() {
	claims1 := suite.claims.CreateClaims("path", "apikey", "hash")
	time.Sleep(time.Millisecond)
	claims2 := suite.claims.CreateClaims("path", "apikey", "hash")
	require.Equal(suite.T(), claims2, claims1)
}
