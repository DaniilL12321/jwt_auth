package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
	"testTaskBackDev/auth"
	"testTaskBackDev/database"
)

func main() {
	conn, err := database.InitDBconnection()

	if err != nil {
		panic(err)
	} else {
		fmt.Println("DB connect\n\n")
	}
	defer conn.PgConn()

	godotenv.Load()

	ip := auth.GetIpUser()
	guid := "123e4567-e89b-12d3-a456-426614174000"
	signature := []byte(os.Getenv("SIGNATURE_SECRET"))
	println("SIGNATURE_SECRET:", signature)

	accessToken, refreshToken, _ := auth.CreatePairTokens(ip, guid, signature)
	{
		fmt.Println("\nаксес:", accessToken)
		//fmt.Println("\nрефреш:", base64.StdEncoding.EncodeToString(refreshToken))
		database.SaveDataUser(conn, "kwemfkwe", "ljkefw", refreshToken)
	}

}
