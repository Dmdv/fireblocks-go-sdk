package fireblocksdk_test

import (
	"crypto/rsa"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/go-chi/jwtauth"
	"github.com/golang-jwt/jwt"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type SignerSuite struct {
	suite.Suite
	url          string
	apiKey       string
	apiSecretKey *rsa.PrivateKey
	baseURL      string
}

func TestSignerSuite(t *testing.T) {
	suite.Run(t, new(SignerSuite))
}

//goland:noinspection ALL
func (suite *SignerSuite) SetupTest() {
	viper.AutomaticEnv()
	viper.SetConfigFile("env.yaml")
	_ = viper.ReadInConfig()

	suite.url = "/v1/supported_assets"
	suite.baseURL = "https://api.fireblocks.io"
	suite.apiKey = viper.GetString("apikey")
	suite.apiSecretKey = GetPrivateKeyFromFile("fireblocks_secret.key")
}

func (suite *SignerSuite) TestSignerSuiteUsingRetryableHTTP() {
	/////////// Signing part

	now := time.Now()
	nowUnix := now.Unix()
	exp := now.Add(10 * time.Second)

	payload := jwt.MapClaims{
		"uri":      suite.url,
		"nonce":    nowUnix,
		"iat":      nowUnix,
		"now":      nowUnix,
		"exp":      exp.Unix(),
		"sub":      suite.apiKey,
		"bodyHash": hashBody([]byte("")),
	}

	token, err := suite.signJwt(payload)
	require.NoError(suite.T(), err)
	require.NotEmpty(suite.T(), token)

	//////////// Request part

	path := fmt.Sprintf("%s%s", suite.baseURL, suite.url)
	suite.T().Logf("GET %s", path)

	request, err := retryablehttp.NewRequest(http.MethodGet, path, prepareBody([]byte("")))
	require.NoError(suite.T(), err)

	request.Header.Add("X-API-Key", suite.apiKey)
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	request.Header.Add("Content-Type", "application/json")

	cl := retryablehttp.NewClient()
	resp, err := cl.Do(request)
	require.NoError(suite.T(), err)
	require.NotNil(suite.T(), resp)

	body, err := ioutil.ReadAll(resp.Body)
	require.NoError(suite.T(), err)
	suite.T().Log(string(body))

	require.Equal(suite.T(), 200, resp.StatusCode)
}

func (suite *SignerSuite) TestSignerSuiteUsingBasicHTTP() {
	auth := jwtauth.New("RS256", suite.apiSecretKey, suite.apiKey)

	payload := map[string]interface{}{
		"uri":      suite.url,
		"nonce":    time.Now().Unix(),
		"iat":      time.Now().Unix(),
		"exp":      time.Now().Unix() + int64(10),
		"sub":      suite.apiKey,
		"bodyHash": hashBody([]byte("")),
	}

	_, token, err := auth.Encode(payload)
	require.NoError(suite.T(), err)
	require.NotEmpty(suite.T(), token)

	headers := http.Header{}

	headers.Set("X-API-Key", suite.apiKey)
	headers.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	headers.Set("Content-Type", "application/json")

	path := fmt.Sprintf("%s%s", suite.baseURL, suite.url)
	suite.T().Logf("GET %s", path)

	request, err := http.NewRequest(http.MethodGet, path, prepareBody([]byte("")))
	require.NoError(suite.T(), err)
	request.Header = headers

	cl := &http.Client{}
	resp, err := cl.Do(request)
	require.NoError(suite.T(), err)
	require.NotNil(suite.T(), resp)

	body, err := ioutil.ReadAll(resp.Body)
	require.NoError(suite.T(), err)
	suite.T().Log(string(body))

	require.Equal(suite.T(), 200, resp.StatusCode)
}

func hashBody(body []byte) string {
	sha := sha256.New()
	_, _ = sha.Write(body)

	return hex.EncodeToString(sha.Sum(nil))
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

func GetPrivateKeyFromFile(path string) *rsa.PrivateKey {
	caPrivateKey, _ := ioutil.ReadFile(path)
	key, _ := jwt.ParseRSAPrivateKeyFromPEM(caPrivateKey)
	return key
}

func (suite *SignerSuite) signJwt(claims jwt.MapClaims) (string, error) {
	token, err := jwt.NewWithClaims(jwt.SigningMethodRS256, claims).SignedString(suite.apiSecretKey)
	if err != nil {
		return "", err
	}

	return token, nil
}
