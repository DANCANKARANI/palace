package utilities

import (
	"math/rand"
	"net/mail"
	"net/smtp"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

func GenerateCode()(string,time.Time) {
	rand.Seed(time.Now().UnixNano())
	//generate a 4 random code
	random_code := rand.Intn(9000)+1000
    exp_time := time.Now().Add(time.Minute*15)
	return strconv.Itoa(random_code),exp_time
}


func SendEmail(email, resetCode string, exp_time time.Time) error {
    err := godotenv.Load(".env")
	if err != nil {
    panic(err.Error())
	}
    
    from := os.Getenv("EMAIL")
    password := os.Getenv("SMTP_PASSWORD")
    smtpHost := "smtp.gmail.com"
    smtpPort := "587"
    // Compose the email
    subject := "Password Reset Code"
    body := "Your password reset code is: " + resetCode+"\n\n"+"Code expires at "+exp_time.Format("15:04")
    msg := []byte("Subject: " + subject + "\r\n" +
        "To: " + email + "\r\n" +
        "\r\n" +
        body)

    // Create the "from" address
   // fromAddr := mail.Address{Name: "Dancan", Address: from}
   fromAddr := mail.Address{ Address: from}
    // Establish a connection to the SMTP server
    auth := smtp.PlainAuth("", from, password, smtpHost)
    err = smtp.SendMail(smtpHost+":"+smtpPort, auth, fromAddr.Address, []string{email}, msg)
    if err != nil {
        return err
    }
	return nil
}