package main

import (
	"log"
	"time"

	"github.com/fangjjcs/bookings-app/pkg/models"

	mail "github.com/xhit/go-simple-mail/v2"
)

func listenForMail(){
	// execute in the background
	go func(){
		for{
			msg := <- app.MailChan
			sendMsg(msg)
		}
	}() 
}

func sendMsg(m models.MailData){

	// setting up a local mail server
	server := mail.NewSMTPClient()
	server.Host = "localhost"
	server.Port = 1025
	server.KeepAlive =  false
	server.ConnectTimeout = 10 * time.Second
	server.SendTimeout = 10 * time.Second
	
	// client
	client, err := server.Connect()
	if err != nil {
		errorLog.Println(err)
	}

	email := mail.NewMSG()
	email.SetFrom(m.From).AddTo(m.To).SetSubject(m.Subject)
	email.SetBody(mail.TextHTML,m.Content)
	
	err = email.Send(client)
	if err != nil{
		log.Println(err)
	}else{
		log.Println("Email Send!")
	}



}