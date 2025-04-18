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

type ErrorResponse struct {
	Error string `json:"error"`
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
	refreshToken := re.RefreshToken

	if accessToken == "" && refreshToken == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "access_token and refresh_token empty"})
		return
	}

	if accessToken == "" {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "access_token empty"})
		return
	}

	if refreshToken == "" {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "refresh_token empty"})
		return
	}

	claims, err := auth.ParseToken(accessToken)
	if err != nil {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
		return
	}
	//log.Println(claims.Sub)

	if claims.Ip != auth.GetIpUser(r) {
		errorResponse := ErrorResponse{
			Error: "ip not correct",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	decodedRefreshToken, err := base64.StdEncoding.DecodeString(refreshToken)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "invalid refresh token: " + err.Error()})
		return
	}
	//log.Println("Decode refreshToken:", decodedRefreshToken)
	isOkToken, err := database.CheckRefreshToken(conn, decodedRefreshToken, claims.Sub)
	if err != nil {
		log.Error(err)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
		return
	}

	if !isOkToken {
		errorResponse := ErrorResponse{
			Error: "token not valid",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	var newAccessToken, newRefreshToken string
	newAccessToken, newRefreshToken, err = database.UpdateRefreshToken(conn, decodedRefreshToken, claims.Sub, r)

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
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "guid not found"})
		return
	}

	isGuid, err := database.CheckGuid(conn, guid)
	if err != nil {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
		return
	}

	if !isGuid {
		log.Warn("guid not found")
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "guid not found"})
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
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "email already used"})
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
