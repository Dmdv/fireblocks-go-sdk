package fireblocksdk

import (
	"time"

	"github.com/golang-jwt/jwt"
)

type IFireblocksClaims interface {
	CreateClaims(path, apiKey, hash string) jwt.MapClaims
}

type FireblocksClaims struct {
	time ITimeProvider
}

func NewClaimsProvider(timeProvider ITimeProvider) IFireblocksClaims {
	return &FireblocksClaims{timeProvider}
}

func (client *FireblocksClaims) CreateClaims(path, apiKey, hash string) jwt.MapClaims {
	now := time.Now()
	nowUnix := now.Unix()
	exp := now.Add(10 * time.Second)

	return jwt.MapClaims{
		"uri":      path,
		"nonce":    nowUnix,
		"iat":      nowUnix,
		"now":      nowUnix,
		"exp":      exp.Unix(),
		"sub":      apiKey,
		"bodyHash": hash,
	}
}
