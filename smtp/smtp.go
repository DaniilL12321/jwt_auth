package smtp

import (
	"fmt"
	"net/smtp"
	"os"
)

func SendIpMessage(ip string, email string) (bool, error) {
	username := os.Getenv("SMTP_USERNAME")
	password := os.Getenv("SMTP_PASSWORD")
	host := os.Getenv("SMTP_HOST")
	port := os.Getenv("SMTP_PORT")

	auth := smtp.PlainAuth("", username, password, host)

	subject := "Security alert"
	body := fmt.Sprintln("🚨maybe you token stolen, new IP (" + ip + ") request with your token")
	message := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s",
		os.Getenv("SMTP_USERNAME")+"@ya.ru",
		email,
		subject,
		body)

	err := smtp.SendMail(host+":"+port, auth, os.Getenv("SMTP_USERNAME")+"@ya.ru", []string{email}, []byte(message))
	if err != nil {
		return false, err
	}

	return true, nil
}
