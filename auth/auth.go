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

func CreatePairTokens(ip, guid string, signature []byte) (accessToken string, refreshToken []byte, err error) {
	accessToken, _ = CreateAccessToken(ip, guid, signature)
	refreshToken, _, _ = CreateRefreshToken()

	//fmt.Println("\nаксес:", accessToken)
	//fmt.Println("\nрефреш:", base64.StdEncoding.EncodeToString(refreshToken))
	//fmt.Println("\nхэш рефреша]:", base64.StdEncoding.EncodeToString(hash))
	return accessToken, refreshToken, nil
}
