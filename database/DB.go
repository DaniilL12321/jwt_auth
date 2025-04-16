package database

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"log"
)

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
