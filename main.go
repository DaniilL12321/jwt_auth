package main

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"os"
	"testTaskBackDev/auth"
	"testTaskBackDev/database"
)

func main() {
	ctx := context.Background()
	dbpool, err := database.InitDBconnection(ctx)

	godotenv.Load()

	if err != nil {
		panic(err)
	} else {
		fmt.Printf("DB connect\n\n")
	}

	defer dbpool.Close()

	ip := "127.0.0.1"
	guid := "123e4567-e89b-12d3-a456-426614174000"
	signature := []byte(os.Getenv("SIGNATURE_SECRET"))
	println("SIGNATURE_SECRET:", signature)

	auth.CreatePairTokens(ip, guid, signature)

}
