package utils

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/soch-fit/GraphGenerator/pkg/configuration"
	"net/smtp"
)

var errorEmail = `Subject: Server failed
From: %s
To: %s

Hello Server unexpectedly died!
Here comes the error causing the problem 

%v

You will need to restart this instance.
`

func sendError(r any) {
	def := configuration.Default()
	if def.ServerMail == nil || def.AdminMail == nil || def.MailServer == nil {
		log.Warningf("Email, server or sender not specified, %v %v %v", configuration.Default().ServerMail, configuration.Default().AdminMail, configuration.Default().MailServer)
		return
	}
	srv, err := smtp.Dial(*def.MailServer)
	if err != nil {
		log.Error(err)
		return
	}
	defer srv.Close()

	err = srv.Mail(*def.ServerMail)
	if err != nil {
		log.Error(err)
		return
	}
	srv.Rcpt(*def.AdminMail)
	if err != nil {

	}
	writer, err := srv.Data()
	if err != nil {
		log.Error(err)
		return
	}
	defer writer.Close()
	_, err = writer.Write([]byte(fmt.Sprintf(errorEmail, *def.ServerMail, *def.AdminMail, r)))
	if err != nil {
		log.Error(err)
		return
	}
}

func StopProcessOnUnhandledPanic() {
	if r := recover(); r != nil {
		sendError(r)
		log.Panicf("Panicikng aftrer delivered panic in function which doesn't permit server continuation, cause %v", r)
	}
}
