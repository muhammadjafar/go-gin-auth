package mail

import (
	"bytes"
	"fmt"
	"golang-api/dto"
	"html/template"
	"net/smtp"
	"os"
	"path"

	"github.com/joho/godotenv"
)

type MailBookData struct {
	Receiver string
}

func MailBook(mailbook MailBookData, book dto.BookCreateDTO) {
	errEnv := godotenv.Load()
	if errEnv != nil {
		panic("Failed to load env file")
	}

	hostURL := os.Getenv("MAIL_HOST")
	hostPort := os.Getenv("MAIL_PORT")
	user := os.Getenv("MAIL_USERNAME")
	password := os.Getenv("MAIL_PASSWORD")
	receiver := mailbook.Receiver
	//set up auth information
	auth := smtp.PlainAuth("", user, password, hostURL)

	fp := path.Join("mail/template", "book-file.html")
	t, errTemp := template.ParseFiles(fp)
	if errTemp != nil {
		fmt.Println(errTemp.Error())
	}
	buff := new(bytes.Buffer)
	t.Execute(buff, book)

	msg := []byte(
		"To: " + receiver + "\r\n" +
			"Subject: New Book!\r\n" +
			"MIME: MIME-version: 1.0\r\n" +
			"Content-Type: text/html; charset=\"UTF-8\";\r\n" +
			"\r\n" +
			buff.String())

	err := smtp.SendMail(hostURL+":"+hostPort, auth, user, []string{receiver}, msg)

	if err != nil {
		fmt.Println(err.Error())
	}
}
