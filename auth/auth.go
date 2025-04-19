package auth

import (
	"crypto/rand"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"os"
	"strings"
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

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claimAcc)

	return token.SignedString(signature)
}

func CreateRefreshToken() (refreshToken []byte, hash []byte, err error) {
	refreshToken = make([]byte, 72)

	rand.Read(refreshToken)

	hash, _ = bcrypt.GenerateFromPassword(refreshToken, bcrypt.DefaultCost)

	return refreshToken, hash, nil
}

func CreatePairTokens(ip, guid string, signature []byte) (accessToken string, refreshToken []byte, hash []byte, err error) {
	accessToken, _ = CreateAccessToken(ip, guid, signature)
	refreshToken, hash, _ = CreateRefreshToken()

	//fmt.Println("\nаксес:", accessToken)
	//fmt.Println("\nрефреш:", base64.StdEncoding.EncodeToString(refreshToken))
	//fmt.Println("\nхэш рефреша]:", base64.StdEncoding.EncodeToString(hash))
	return accessToken, refreshToken, hash, nil
}

func GetIpUser(r *http.Request) string {
	ip := r.Header.Get("X-Forwarded-For")
	if ip == "" {
		return ""
	}
	ips := strings.Split(ip, ",")
	ip = strings.TrimSpace(ips[0])
	return ip
}

type ParsedClaims struct {
	Ip  string `json:"ip"`
	Exp int64  `json:"exp"`
	Sub string `json:"sub"`
	jwt.RegisteredClaims
}

func ParseToken(accessToken string) (*ParsedClaims, error) {
	token, err := jwt.ParseWithClaims(accessToken, &ParsedClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("SIGNATURE_SECRET")), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*ParsedClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
