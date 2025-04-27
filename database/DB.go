package database

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v4"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"os"
	"testTaskBackDev/auth"
)

type User struct {
	Guid         uuid.UUID `db:"guid"`
	Email        string    `db:"email"`
	Password     string    `db:"password"`
	RefreshToken []byte    `db:"refresh_token"`
}

func InitDBconnection() (*pgx.Conn, error) {
	conn, err := pgx.Connect(context.Background(), os.Getenv("DB_URL"))
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func SaveDataUser(conn *pgx.Conn, email string, password string, r *http.Request) (string, string, error) {
	var hash []byte

	id, _ := uuid.NewV4()
	accessToken, refreshToken, hash, _ := auth.CreatePairTokens(auth.GetIpUser(r), id.String(), []byte(os.Getenv("SIGNATURE_SECRET")))
	passwordHash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	rows, err := conn.Query(context.Background(), "INSERT into users(guid, email, password, refresh_token) VALUES ($1, $2, $3, $4)", id, email, passwordHash, hash)
	if err != nil {
		return "", "", err
	}

	for rows.Next() {
		rows.Scan(&email, &password, &hash)
	}

	defer rows.Close()
	//fmt.Println("data user successfully saved")
	return accessToken, base64.StdEncoding.EncodeToString(refreshToken), nil
}

func UpdateRefreshToken(conn *pgx.Conn, refreshToken []byte, accessToken string, guid string, r *http.Request) (string, string, error) {
	acceptRefresh, _ := CheckRefreshToken(conn, refreshToken, accessToken, guid)

	if acceptRefresh {
		ip := auth.GetIpUser(r)
		signature := []byte(os.Getenv("SIGNATURE_SECRET"))

		accessToken, refreshToken, hash, _ := auth.CreatePairTokens(ip, guid, signature)
		query, err := conn.Query(context.Background(), "UPDATE users set refresh_token = $1 WHERE users.guid = $2", hash, guid)
		if err != nil {
			return "", "", err
		}

		defer query.Close()
		//fmt.Println("user refresh token successfully updated")
		//fmt.Println("new access: ", accessToken)
		//fmt.Println("new refresh: ", base64.StdEncoding.EncodeToString(refreshToken))
		//fmt.Println("new hash: ", hash)

		return accessToken, base64.StdEncoding.EncodeToString(refreshToken), nil
	}

	return "", "", nil
}

func CheckRefreshToken(conn *pgx.Conn, refreshToken []byte, accessToken string, guid string) (bool, error) {
	savedHashRefreshToken, _ := FindLastRefreshToken(conn, guid)

	if err := bcrypt.CompareHashAndPassword(savedHashRefreshToken, refreshToken); err != nil {
		//fmt.Println("not correct refresh token")
		//println(base64.StdEncoding.EncodeToString(refreshToken))
		//println(base64.StdEncoding.EncodeToString(savedHashRefreshToken))
		return false, err
	}

	partAccessToken := refreshToken[len(refreshToken)-5:]

	currentHashPart := sha256.Sum256([]byte(accessToken))
	currentPartAccessToken := currentHashPart[:5]

	if !bytes.Equal(partAccessToken, currentPartAccessToken) {
		return false, errors.New("refresh token not linked to current access token")
	}

	//fmt.Println("token ok")
	return true, nil
}

func FindLastRefreshToken(conn *pgx.Conn, guid string) ([]byte, error) {
	query, _ := conn.Query(context.Background(), "SELECT refresh_token FROM users WHERE guid = $1", guid)
	defer query.Close()

	var refreshToken []byte
	if query.Next() {
		query.Scan(&refreshToken)
	}

	return refreshToken, nil
}

// for updating with ID by parametrs Get request
func UpdateRefreshTokenById(conn *pgx.Conn, guid string, r *http.Request) (string, string, error) {

	ip := auth.GetIpUser(r)
	signature := []byte(os.Getenv("SIGNATURE_SECRET"))

	accessToken, refreshToken, hash, _ := auth.CreatePairTokens(ip, guid, signature)
	query, err := conn.Query(context.Background(), "UPDATE users set refresh_token = $1 WHERE users.guid = $2", hash, guid)
	if err != nil {
		return "", "", err
	}

	defer query.Close()
	//fmt.Println("user refresh token successfully updated")
	//fmt.Println("new access: ", accessToken)
	//fmt.Println("new refresh: ", base64.StdEncoding.EncodeToString(refreshToken))
	//fmt.Println("new hash: ", hash)

	return accessToken, base64.StdEncoding.EncodeToString(refreshToken), nil

}

func CheckGuid(conn *pgx.Conn, guid string) (bool, error) {
	query, err := conn.Query(context.Background(), "SELECT guid FROM users WHERE guid = $1", guid)
	if err != nil {
		return false, err
	}
	defer query.Close()

	if query.Next() {
		return true, nil
	}

	return false, nil
}

func CheckEmail(conn *pgx.Conn, email string) (bool, error) {
	query, err := conn.Query(context.Background(), "SELECT email FROM users WHERE email = $1", email)
	if err != nil {
		return false, err
	}
	defer query.Close()
	if query.Next() {
		return false, nil
	}
	return true, nil
}

func FindEmailFromId(conn *pgx.Conn, id string) (string, error) {
	query, err := conn.Query(context.Background(), "SELECT email FROM users WHERE guid = $1", id)
	if err != nil {
		return "", err
	}
	defer query.Close()
	var email string
	if query.Next() {
		query.Scan(&email)
	}

	return email, nil
}
