package database

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v4"
	"golang.org/x/crypto/bcrypt"
	"log"
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
	conn, err := pgx.Connect(context.Background(), "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable")
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func SaveDataUser(conn *pgx.Conn, email string, password string, refreshToken []byte) {
	rows, err := conn.Query(context.Background(), "INSERT into users(email, password, refresh_token) VALUES ($1, $2, $3)", email, password, refreshToken)
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()
	fmt.Println("data user successfully saved")
}

func UpdateRefreshToken(conn *pgx.Conn, refreshToken []byte, guid string) (string, []byte) {
	acceptRefresh := CheckRefreshToken(conn, refreshToken, guid)
	{
		if acceptRefresh {
			ip := auth.GetIpUser()
			signature := []byte(os.Getenv("SIGNATURE_SECRET"))

			accessToken, refreshToken, hash, _ := auth.CreatePairTokens(ip, guid, signature)
			query, err := conn.Query(context.Background(), "UPDATE users set refresh_token = $1 WHERE users.guid = $2", hash, guid)
			if err != nil {
				log.Fatal(err)
				return "", nil
			}

			defer query.Close()
			fmt.Println("user refresh token successfully updated")
			fmt.Println("new access: ", accessToken)
			fmt.Println("new refresh: ", refreshToken)
			fmt.Println("new hash: ", hash)

			return accessToken, refreshToken
		}
	}

	return "", nil
}

func CheckRefreshToken(conn *pgx.Conn, refreshToken []byte, guid string) bool {
	savedHashRefreshToken, _ := FindLastRefreshToken(conn, guid)
	if err := bcrypt.CompareHashAndPassword(savedHashRefreshToken, refreshToken); err != nil {
		fmt.Println("not correct refresh token")
		println(base64.StdEncoding.EncodeToString(refreshToken))
		println(base64.StdEncoding.EncodeToString(savedHashRefreshToken))
		return false
	}
	fmt.Println("token ok")
	return true
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
