package fireblocksdk

import (
	"crypto/sha256"
	"encoding/hex"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/pkg/errors"
)

type ITimeProvider interface {
	Now() time.Time
}

type TimeProvider struct{}

func (tp *TimeProvider) Now() time.Time {
	return time.Now()
}

func DefaultTimeProvider() ITimeProvider {
	return &TimeProvider{}
}

func DefaultClaimProvider(time ITimeProvider) IFireblocksClaims {
	return NewClaimsProvider(time)
}

type IAuthProvider interface {
	SignJwt(path string, bodyJSON []byte) (string, error)
	GetApiKey() string
}

type AuthProviderConfig struct {
	timeProvider ITimeProvider
}

type AuthProvider struct {
	apiKey        string
	apiSecretKey  []byte
	claimProvider IFireblocksClaims
}

func WithTimeProvider(tp ITimeProvider) func(c *AuthProviderConfig) error {
	return func(c *AuthProviderConfig) error {
		c.timeProvider = tp

		return nil
	}
}

// NewAuthProvider Creates signer using api key and private key from config
func NewAuthProvider(apiKey string, apiSecretKey []byte, configs ...func(*AuthProviderConfig) error) (*AuthProvider, error) {
	cfg := &AuthProviderConfig{DefaultTimeProvider()}
	for _, conf := range configs {
		err := conf(cfg)
		if err != nil {
			return nil, errors.Wrap(err, "invalid/unsupported options")
		}
	}

	auth := &AuthProvider{
		apiKey,
		apiSecretKey,
		DefaultClaimProvider(cfg.timeProvider),
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
	signJwt, err := ap.signJwt(claims)

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
