package fireblocksdk

import (
	"time"

	"github.com/golang-jwt/jwt"
)

type IFireblocksClaims interface {
	CreateClaims(path, apiKey, hash string) jwt.MapClaims
}

type FireblocksClaims struct {
	time          ITimeProvider
	expirySeconds int64
}

func NewFireblocksClaimsProvider(timeProvider ITimeProvider, exp int64) IFireblocksClaims {
	return &FireblocksClaims{timeProvider, exp}
}

func (client *FireblocksClaims) CreateClaims(path, apiKey, hash string) jwt.MapClaims {
	now := client.time.Now()
	nowUnix := now.Unix()
	exp := now.Add(time.Second * time.Duration(client.expirySeconds))

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
