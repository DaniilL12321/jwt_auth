package auth

import (
	"crypto/rand"
	"crypto/sha256"
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

func CreateRefreshToken(accessToken string) (refreshToken []byte, hash []byte, err error) {
	hashPart := sha256.Sum256([]byte(accessToken))
	partAccessToken := hashPart[:5]

	randomPart := make([]byte, 67)
	rand.Read(randomPart)

	refreshToken = append(randomPart, partAccessToken...)

	hash, _ = bcrypt.GenerateFromPassword(refreshToken, bcrypt.DefaultCost)

	return refreshToken, hash, nil
}

func CreatePairTokens(ip, guid string, signature []byte) (accessToken string, refreshToken []byte, hash []byte, err error) {
	accessToken, _ = CreateAccessToken(ip, guid, signature)
	refreshToken, hash, _ = CreateRefreshToken(accessToken)

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
	for _, ip := range ips {
		ip = strings.TrimSpace(ip)
		if ip != "" {
			return ip
		}
	}
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
