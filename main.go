package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"testTaskBackDev/auth"
	"testTaskBackDev/database"
)

func main() {
	ctx := context.Background()
	dbpool, err := database.InitDBconnection(ctx)

	if err != nil {
		panic(err)
	} else {
		fmt.Printf("DB connect\n\n")
	}

	defer dbpool.Close()

	ip := "127.0.0.1"
	guid := "123e4567-e89b-12d3-a456-426614174000"
	signature := []byte("kmkmewfml")

	accessToken, _ := auth.CreateAccessToken(ip, guid, signature)
	refreshToken, hash, _ := auth.CreateRefreshToken()
	fmt.Println("\nаксес:", accessToken)
	fmt.Println("\nрефреш:", base64.StdEncoding.EncodeToString(refreshToken))
	fmt.Println("\nхэш рефреша]:", base64.StdEncoding.EncodeToString(hash))

}
