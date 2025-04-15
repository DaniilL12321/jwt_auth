package auth

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type claims struct {
	Ip string `json:"ip"`
	jwt.RegisteredClaims
}

func CreateAccessToken(ip, guid string, signature []byte) (string, error) {
	claimAcc := claims{
		ip,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
			Subject:   guid,
		},
	}

	return jwt.NewWithClaims(jwt.SigningMethodHS512, claimAcc).SignedString(signature)
}
