package main

import (
	"encoding/base64"
	"encoding/json"
	"github.com/jackc/pgx/v4"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
	log "github.com/sirupsen/logrus"
	"github.com/swaggo/http-swagger"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"net/mail"
	"testTaskBackDev/auth"
	"testTaskBackDev/database"
	_ "testTaskBackDev/docs"
	"testTaskBackDev/smtp"
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

// @title [JWT tokens service] Swagger API
// @version 0.0.1
// @host jwt-auth-4tmd.onrender.com
// @BasePath /

// ** \\ при запросе со свагера в локалке используетс http на хост, где https
// ** \\ из-за этого запрос не проходит, для решения проблемы нужно добавить // @schemes https
// ** \\ либо поменять на локальный хост (localhost:8080)

// @description This documentation describes [JWT tokens service] Swagger API
// @contact.name github Open Source Code
// @contact.url https://github.com/DaniilL12321/jwt_auth
// @license.name MIT License
func main() {
	godotenv.Load()
	startTime = time.Now()

	http.HandleFunc("/swagger/", httpSwagger.WrapHandler)

	http.HandleFunc("GET /uptime", uptimeCheck)
	http.HandleFunc("GET /tokens", createPairById)
	http.HandleFunc("POST /refresh", createPairByTokens)
	http.HandleFunc("POST /register", createUser)

	handler := cors.Default().Handler(http.DefaultServeMux)

	http.ListenAndServe(":8080", handler)
}

var startTime time.Time

// @uptime godoc
//
// @Summary get time work server
// @Tags default
// @Success 200
// @Router /uptime [get]
func uptimeCheck(w http.ResponseWriter, request *http.Request) {

	uptime := time.Since(startTime).Milliseconds()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(uptime)
	log.Println("check uptime: ", uptime)
	return
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

// @createPairByTokens godoc
//
// @Summary get pair new access+refresh tokens by old pair tokens
// @Tags auth
// @Accept json
// @Produce json
// @Param request body Request true "old tokens"
// @Success 200 {array} Response
// @Failure 400 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /refresh [post]
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

	var email string
	ip := auth.GetIpUser(r)
	if claims.Ip != ip {
		email, err = database.FindEmailFromId(conn, claims.Sub)
		if err != nil {
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
			log.Error("alert message not send to " + email + " by new IP request with token: " + ip)
			return
		}

		log.Println("alert message send to " + email + " by new IP request with token: " + ip)

		_, err := smtp.SendIpMessage(ip, email)
		if err != nil {
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "ip not correct"})
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
	if err == bcrypt.ErrMismatchedHashAndPassword {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "refresh_token hash and password not same"})
		return
	}
	if err != nil {
		log.Error(err)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
		return
	}

	if !isOkToken {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "token not valid"})
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

// @createPairById godoc
//
// @Summary get pair access+refresh tokens by ID
// @Tags auth
// @Accept json
// @Produce json
// @Param guid query string true "User ID" format(uuid)
// @Success 200 {array} Response
// @Failure 400 {object} ErrorResponse
// @Router /tokens [get]
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

// @createUser godoc
//
// @Summary register user in DB and get start pair tokens and ID
// @Tags auth
// @Accept json
// @Produce json
// @Param user body User true "user register data"
// @Success 201 {array} Response
// @Failure 400 {object} ErrorResponse
// @Router /register [post]
func createUser(w http.ResponseWriter, r *http.Request) {
	conn := connectToDb()

	var user User
	json.NewDecoder(r.Body).Decode(&user)

	email := user.Email
	password := user.Password

	if !validEmail(email) {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "email address is invalid"})
		return
	}

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

func validEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}
