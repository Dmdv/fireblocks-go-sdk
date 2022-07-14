package fireblocksdk

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/pkg/errors"
)

type IAuthProvider interface {
	SignJwt(path string, bodyJson interface{}) (string, error)
	GetApiKey() string
}

type AuthProvider struct {
	apiKey, privateKey string
}

// NewAuthProvider Creates signer using api key and private key from config
func NewAuthProvider(apiKey, privateKey string) *AuthProvider {
	return &AuthProvider{apiKey, privateKey}
}

// SignJwt Creates token using path and payload
func (ap *AuthProvider) SignJwt(path string, bodyJson interface{}) (string, error) {
	bodyBytes, err := json.Marshal(bodyJson)
	if err != nil {
		return "", errors.Wrap(err, "failed to marshal body")
	}

	hash, err := hashBody(bodyBytes)
	if err != nil {
		return "", err
	}

	now := time.Now()
	exp := now.Add(10 * time.Second)

	signJwt, err := ap.signJwt(jwt.MapClaims{
		"uri":      path,
		"nonce":    now.Unix(),
		"now":      now.Unix(),
		"exp":      exp.Unix(),
		"sub":      ap.apiKey,
		"bodyHash": fmt.Sprintf("%x", hash),
	})

	return signJwt, err
}

func (ap *AuthProvider) GetApiKey() string {
	return ap.apiKey
}

func hashBody(body []byte) (string, error) {
	sha := sha256.New()
	_, err := sha.Write(body)
	if err != nil {
		return "", errors.Wrap(err, "failed to hash the body")
	}

	return hex.EncodeToString(sha.Sum(nil)), nil
}

func (ap *AuthProvider) signJwt(claims jwt.MapClaims) (string, error) {
	key, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(ap.privateKey))
	if err != nil {
		return "", errors.Wrap(err, "failed to read token from string")
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodRS256, claims).SignedString(key)
	if err != nil {
		return "", errors.Wrap(err, "failed to create token")
	}

	return token, nil
}
