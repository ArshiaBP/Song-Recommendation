package services

import (
	"fmt"
	"net/smtp"
)

var (
	host = "smtp.gmail.com"
	port = "587"
	from = "planverse.companies@gmail.com"
	//password = os.Getenv("GmailPassword")
	password = "kihhxbbdoinhdwql"
)

func sendMail(subject, body string, to []string) error {
	auth := smtp.PlainAuth("", from, password, host)
	err := smtp.SendMail(fmt.Sprint(host+":"+port), auth, from, to, []byte(fmt.Sprint("Subject: "+subject+"\n"+body)))
	if err != nil {
		return err
	}
	return nil
}
