package fireblocksdk

import (
	"crypto/rsa"
	"crypto/sha256"
	"encoding/hex"

	"github.com/golang-jwt/jwt"
	"github.com/pkg/errors"
)

func DefaultTimeProvider() ITimeProvider {
	return &TimeProvider{}
}

func DefaultClaimProvider(time ITimeProvider, exp int64) IFireblocksClaims {
	return NewFireblocksClaimsProvider(time, exp)
}

func DefaultTokenExpiry() int64 {
	return 10
}

type IAuthProvider interface {
	SignJwt(path string, bodyJSON []byte) (string, error)
	GetApiKey() string
}

type AuthProviderConfig struct {
	timeProvider  ITimeProvider
	expirySeconds int64
}

type AuthProvider struct {
	apiKey        string
	privateKey    *rsa.PrivateKey
	claimProvider IFireblocksClaims
}

func WithTokenExpiry(exp int64) func(c *AuthProviderConfig) error {
	return func(c *AuthProviderConfig) error {
		c.expirySeconds = exp

		return nil
	}
}

func WithTimeProvider(tp ITimeProvider) func(c *AuthProviderConfig) error {
	return func(c *AuthProviderConfig) error {
		c.timeProvider = tp

		return nil
	}
}

// NewAuthProvider Creates signer using api key and private key from config
func NewAuthProvider(apiKey string, apiSecretKey []byte, configs ...func(*AuthProviderConfig) error) (*AuthProvider, error) {
	cfg := &AuthProviderConfig{DefaultTimeProvider(), DefaultTokenExpiry()}
	for _, conf := range configs {
		err := conf(cfg)
		if err != nil {
			return nil, errors.Wrap(err, "invalid/unsupported options")
		}
	}

	key, err := jwt.ParseRSAPrivateKeyFromPEM(apiSecretKey)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read token from string")
	}

	auth := &AuthProvider{
		apiKey,
		key,
		DefaultClaimProvider(cfg.timeProvider, cfg.expirySeconds),
	}

	return auth, nil
}

// SignJwt Creates token using path and payload
func (ap *AuthProvider) SignJwt(path string, bodyJSON []byte) (string, error) {
	hash, err := hashBody(bodyJSON)
	if err != nil {
		return "", err
	}

	claims := ap.claimProvider.CreateClaims(path, ap.apiKey, hash)

	jwtToken, err := jwt.NewWithClaims(jwt.SigningMethodRS256, claims).SignedString(ap.privateKey)
	if err != nil {
		return "", errors.Wrap(err, "failed to create token")
	}

	return jwtToken, err
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
