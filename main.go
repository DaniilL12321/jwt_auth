package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"testTaskBackDev/auth"
	"testTaskBackDev/database"
	"time"
)

type Request struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type Response struct {
	Guid         string    `json:"guid"`
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
}

func main() {
	godotenv.Load()

	http.HandleFunc("GET /", createPairById)
	http.HandleFunc("POST /", createPairByTokens)
	http.ListenAndServe(":8080", nil)
}

func connectToDb() (conn *pgx.Conn) {
	conn, err := database.InitDBconnection()

	if err != nil {
		panic(err)
	} else {
		fmt.Println("DB connect\n\n")
	}
	defer conn.PgConn()
	return conn
}

func createPairByTokens(w http.ResponseWriter, r *http.Request) {
	conn := connectToDb()

	var re Request
	json.NewDecoder(r.Body).Decode(&re)

	accessToken := re.AccessToken

	claims, _ := auth.ParseToken(accessToken)
	println(claims.Sub)

	ippp := auth.GetIpUser(r)

	println(ippp)

	if claims.Ip != auth.GetIpUser(r) {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("ip not correct"))
		return
	}

	decodedRefreshToken, _ := base64.StdEncoding.DecodeString(re.RefreshToken)
	fmt.Println("Decode refreshToken:", decodedRefreshToken)
	isOkToken, err := database.CheckRefreshToken(conn, decodedRefreshToken, claims.Sub)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(err.Error()))
		return
	}

	if !isOkToken {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("token not corrected"))
		return
	}

	newAccessToken, newRefreshToken := database.UpdateRefreshToken(conn, decodedRefreshToken, claims.Sub, r)

	response := Response{
		Guid:         claims.Sub,
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
		ExpiresAt:    time.Now().Add(time.Hour),
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(response)

}

// for updating with ID by parametrs Get request
func createPairById(w http.ResponseWriter, r *http.Request) {
	conn := connectToDb()
	guid := r.URL.Query().Get("guid")
	if guid == "" {
		log.Println("guid not found")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("guid not found"))
		return
	}

	isGuid, _ := database.CheckGuid(conn, guid)
	if !isGuid {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("guid not found"))
		return
	}

	newAccessToken, newRefreshToken, _ := database.UpdateRefreshTokenById(conn, guid, r)

	response := Response{
		Guid:         guid,
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
		ExpiresAt:    time.Now().Add(time.Hour),
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(response)

}
