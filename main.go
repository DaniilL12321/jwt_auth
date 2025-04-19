package main

import (
	"encoding/base64"
	"encoding/json"
	"github.com/jackc/pgx/v4"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
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

type User struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func main() {
	godotenv.Load()

	http.HandleFunc("GET /tokens", createPairById)
	http.HandleFunc("POST /refresh", createPairByTokens)
	http.HandleFunc("POST /register", createUser)
	http.ListenAndServe(":8080", nil)
}

func connectToDb() (conn *pgx.Conn) {
	conn, err := database.InitDBconnection()

	if err != nil {
		panic(err)
	} else {
		//log.Print("DB connect\n")
	}
	defer conn.PgConn()
	return conn
}

// for updating with tokens by body Post request
func createPairByTokens(w http.ResponseWriter, r *http.Request) {
	conn := connectToDb()

	var re Request
	json.NewDecoder(r.Body).Decode(&re)

	accessToken := re.AccessToken

	claims, _ := auth.ParseToken(accessToken)
	//log.Println(claims.Sub)

	if claims.Ip != auth.GetIpUser(r) {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("ip not correct"))
		return
	}

	decodedRefreshToken, _ := base64.StdEncoding.DecodeString(re.RefreshToken)
	//log.Println("Decode refreshToken:", decodedRefreshToken)
	isOkToken, err := database.CheckRefreshToken(conn, decodedRefreshToken, claims.Sub)
	if err != nil {
		log.Error(err)
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

	log.Print("tokens refresh with used pair access and refresh token for ", claims.Sub)
}

// for updating with ID by parametrs Get request
func createPairById(w http.ResponseWriter, r *http.Request) {
	conn := connectToDb()
	guid := r.URL.Query().Get("guid")
	if guid == "" {
		log.Warn("guid not found")
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

	log.Print("tokens refresh with used ID parameter for ", guid)
}

func createUser(w http.ResponseWriter, r *http.Request) {
	conn := connectToDb()

	var user User
	json.NewDecoder(r.Body).Decode(&user)

	email := user.Email
	password := user.Password

	EmailIsOk, _ := database.CheckEmail(conn, email)
	if !EmailIsOk {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("email already used"))
		log.Warn("user with email: ", email, " already exists")
		return
	}

	accessToken, refreshToken, _ := database.SaveDataUser(conn, email, password, r)

	claims, _ := auth.ParseToken(accessToken)

	response := Response{
		Guid:         claims.Sub,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    time.Now().Add(time.Hour),
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)

	log.Print("register new user: ", email)
}
