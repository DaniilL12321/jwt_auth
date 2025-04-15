package main

import (
	"context"
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
		fmt.Printf("DB connect")
	}

	defer dbpool.Close()

	ip := "127.0.0.1"
	guid := "123e4567-e89b-12d3-a456-426614174000"
	signature := []byte("kmkmewfml")

	token, _ := auth.CreateAccessToken(ip, guid, signature)
	fmt.Println("\n\nТокен:", token)
}
