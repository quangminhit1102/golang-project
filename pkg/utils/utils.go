package utils

import (
	"log"
	"net/smtp"

	"github.com/google/uuid"
)

func SendMail(toEmail, subject, body string) {
	// Choose auth method and set it up
	auth := smtp.PlainAuth("", "quangminhit.test01@gmail.com", "czkizmswbiljzmos", "smtp.gmail.com")
	// Here we do it all: connect to our server, set up a message and send it
	to := []string{toEmail}
	msg := []byte(body)
	err := smtp.SendMail("smtp.gmail.com:587", auth, "quangminhit.test01@gmail.com", to, msg)
	if err != nil {
		log.Fatal(err)
	}
}

func RandomTokenGenerator() string {
	// This is not professional Way
	// b := make([]byte, 4)
	// rand.Read(b)
	// return fmt.Sprintf("%x", b)

	// Professional Way
	ranUuid, err := uuid.NewRandom()
	if err != nil {
		return uuid.Nil.String()
	}
	return ranUuid.String()
}
