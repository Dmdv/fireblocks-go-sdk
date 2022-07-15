package fireblocksdk

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/pkg/errors"
)

type ITimeProvider interface {
	Now() time.Time
}

type TimeProvider struct {
}

func (tp *TimeProvider) Now() time.Time {
	return time.Now()
}

func DefaultTimeProvider() *TimeProvider {
	return &TimeProvider{}
}

type IAuthProvider interface {
	SignJwt(path string, bodyJSON []byte) (string, error)
	GetApiKey() string
}

type AuthProvider struct {
	apiKey       string
	apiSecretKey []byte
	timeProvider ITimeProvider
}

func WithTimeProvider(tp ITimeProvider) func(c *AuthProvider) error {
	return func(c *AuthProvider) error {
		c.timeProvider = tp

		return nil
	}
}

// NewAuthProvider Creates signer using api key and private key from config
func NewAuthProvider(apiKey string, apiSecretKey []byte, configs ...func(*AuthProvider) error) (*AuthProvider, error) {
	auth := &AuthProvider{apiKey, apiSecretKey, DefaultTimeProvider()}

	for _, conf := range configs {
		err := conf(auth)
		if err != nil {
			return nil, errors.Wrap(err, "invalid/unsupported options")
		}
	}

	return auth, nil
}

// SignJwt Creates token using path and payload
func (ap *AuthProvider) SignJwt(path string, bodyJSON []byte) (string, error) {
	hash, err := hashBody(bodyJSON)
	if err != nil {
		return "", err
	}

	now := ap.timeProvider.Now()
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
	key, err := jwt.ParseRSAPrivateKeyFromPEM(ap.apiSecretKey)
	if err != nil {
		return "", errors.Wrap(err, "failed to read token from string")
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodRS256, claims).SignedString(key)
	if err != nil {
		return "", errors.Wrap(err, "failed to create token")
	}

	return token, nil
}
