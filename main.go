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
	guid := "1d51b15b-9ca5-4a7e-8803-507152ff7003"
	signature := []byte(os.Getenv("SIGNATURE_SECRET"))
	println("SIGNATURE_SECRET:", signature)

	accessToken, _, _, _ := auth.CreatePairTokens(ip, guid, signature)
	{
		fmt.Println("\nаксес:", accessToken)
		//fmt.Println("\nрефреш:", refreshToken)
		//fmt.Println("\nхэш:", hash)
		//database.SaveDataUser(conn, "lkqemflkmqelkfmqelkmf", "1212ljkefw", hash)
		//database.FindLastRefreshToken(conn, guid)
	}

	refreshToken := []byte{
		107, 59, 95, 156, 176, 106, 149, 219, 120, 206, 97, 135, 174, 162, 106, 207,
		88, 157, 182, 89, 66, 110, 130, 239, 105, 117, 248, 237, 14, 198, 35, 116,
		251, 35, 173, 76, 254, 216, 62, 28, 120, 245, 178, 145, 86, 172, 187, 170,
		97, 247, 131, 185, 108, 22, 6, 69, 219, 26, 223, 11, 194, 109, 82, 250,
		103, 123, 141, 196, 248, 64, 82, 10,
	}

	database.UpdateRefreshToken(conn, refreshToken, guid)

}
