package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/joho/godotenv"
	"log"
	"net/http"
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

	http.HandleFunc("POST /", createPair)
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

func createPair(w http.ResponseWriter, r *http.Request) {
	conn := connectToDb()
	guid := r.URL.Query().Get("guid")
	if guid == "" {
		log.Println("guid not found")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("guid not found"))
		return
	}

	var re Request
	json.NewDecoder(r.Body).Decode(&re)

	accessToken := re.AccessToken
	decodedRefreshToken, _ := base64.StdEncoding.DecodeString(re.RefreshToken)
	fmt.Println("Decode refreshToken:", decodedRefreshToken)
	isOkToken, err := database.CheckRefreshToken(conn, accessToken, decodedRefreshToken, guid)
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

	newAccessToken, newRefreshToken := database.UpdateRefreshToken(conn, accessToken, decodedRefreshToken, guid)

	response := Response{
		Guid:         guid,
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
		ExpiresAt:    time.Now().Add(time.Hour),
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(response)

}
