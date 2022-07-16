package fireblocksdk_test

import (
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/spf13/viper"

	"github.com/go-chi/jwtauth"
	//"github.com/lestrrat-go/jwx/jwt"
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
	viper.ReadInConfig()

	suite.url = "/v1/supported_assets"
	suite.baseURL = "https://api.fireblocks.io"
	suite.apiKey = viper.GetString("apikey")
	suite.apiSecretKey = GetPrivateKeyFromFile("fireblocks_secret.key")
}

func (suite *SignerSuite) TestSignerSuite() {
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
	headers.Set("Authorization", "Bearer "+token)
	headers.Set("Content-Type", "application/json")

	path := fmt.Sprintf("%s%s", suite.baseURL, suite.url)
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

func getPrivateKeyFromString(value []byte) *rsa.PrivateKey {
	block, _ := pem.Decode(value)
	if block == nil {
		return nil
	}

	parseResult, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	key := parseResult.(*rsa.PrivateKey)
	if err != nil {
		return nil
	}
	return key
}

func GetPrivateKeyFromFile(path string) *rsa.PrivateKey {
	caPrivateKey, _ := ioutil.ReadFile(path)
	return getPrivateKeyFromString(caPrivateKey)
}
