package utils

import (
	"fmt"
	"os"

	"gopkg.in/gomail.v2"
)

// this send email by smtp  may take 1-3 second , so its not go fault if register take 1-3 second
func SendEmail(toEmail, token string) {
	mailer := gomail.NewMessage()
	mailer.SetHeader("From", "rifqiadlihernawan@gmail.com")
	mailer.SetHeader("To", toEmail)
	mailer.SetHeader("Subject", "Verify Your Account")
	mailer.SetBody("text/html", fmt.Sprintf(`<a href="http://localhost:8080/verification?email=%s&token=%s">Klik di sini untuk verifikasi</a>`, toEmail, token))
	dialer := gomail.NewDialer("smtp.gmail.com", 587, os.Getenv("EMAIL_SENDER"), os.Getenv("APP_PASSWORD"))

	if err := dialer.DialAndSend(mailer); err != nil {
		fmt.Println("Error sending email:", err)
	}
}
