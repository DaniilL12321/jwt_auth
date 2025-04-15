package auth

import (
	"crypto/rand"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type claims struct {
	Ip string `json:"ip"`
	jwt.RegisteredClaims
}

func CreateAccessToken(ip, guid string, signature []byte) (accessToken string, err error) {
	claimAcc := claims{
		ip,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
			Subject:   guid,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claimAcc)

	return token.SignedString(signature)
}

func CreateRefreshToken() (refreshToken []byte, hash []byte, err error) {
	refreshToken = make([]byte, 72)

	rand.Read(refreshToken)

	hash, _ = bcrypt.GenerateFromPassword(refreshToken, bcrypt.DefaultCost)

	return refreshToken, hash, nil
}
