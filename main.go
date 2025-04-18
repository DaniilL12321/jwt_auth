package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"testTaskBackDev/auth"
	"testTaskBackDev/database"
)

type Request struct {
	RefreshToken string `json:"refresh_token"`
}

type Response struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func main() {
	conn := connectToDb()
	godotenv.Load()

	ip := auth.GetIpUser()
	guid := "5d424e86-29f7-4f0e-9d23-c621694cb938"
	signature := []byte(os.Getenv("SIGNATURE_SECRET"))
	println("SIGNATURE_SECRET:", signature)

	accessToken, refreshToken, hash, _ := auth.CreatePairTokens(ip, guid, signature)
	{
		//baseRefresh := base64.StdEncoding.EncodeToString(refreshToken)
		//noBaseRefresh, _ := base64.StdEncoding.DecodeString(baseRefresh)
		fmt.Println("\nаксес:", accessToken)

		//fmt.Println("\nбейс64:", []byte(baseRefresh))
		//fmt.Println("\nобратно:", noBaseRefresh)

		fmt.Println("\nрефреш:", base64.StdEncoding.EncodeToString(refreshToken))
		fmt.Println("\nхэш:", base64.StdEncoding.EncodeToString(hash))
		database.SaveDataUser(conn, "email.com", "alkmfklmeklmef", hash)
		//database.FindLastRefreshToken(conn, guid)
	}

	refreshToken, _ = base64.StdEncoding.DecodeString("HYJcGEBPUBcmrL8BuY9TFVQVml96/JENe4JT8yiCo9vekLVU4pSoZaWf2dYwhCve5ni0yLM9Of6Ja0W5UdcnyEafMxMpY9V3")

	fmt.Println("kjwenfkjnmwefjwjeknfjkwenjkwnegjknwegjkn", refreshToken)
	database.UpdateRefreshToken(conn, refreshToken, guid)

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

	decodedRefreshToken, _ := base64.StdEncoding.DecodeString(re.RefreshToken)
	fmt.Println("Decode refreshToken:", decodedRefreshToken)
	isOkToken, err := database.CheckRefreshToken(conn, decodedRefreshToken, guid)
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

	accessToken, newRefreshToken := database.UpdateRefreshToken(conn, decodedRefreshToken, guid)

	response := Response{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(response)

}
