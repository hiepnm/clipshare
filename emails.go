package main

import (
	"net/smtp"
	"bytes"
	"net"
	"text/template"
)

type SmtpTemplateActiveData struct {
	From        string
	To          string
	DisplayName string
	UserName    string
	Password    string
	Uid         string
}
func sendActivateTempUserMail(recipients string, uid string, tempUser *TempUser) error {
	var err error
	var doc bytes.Buffer
	serverMail, _, _ := net.SplitHostPort(config.ServerMail)
	auth := smtp.PlainAuth("", config.activeEmail, config.activeEmailPass, serverMail)
	context := &SmtpTemplateActiveData{config.activeEmail, recipients, tempUser.displayName, tempUser.userName, tempUser.password, uid}
	t := template.New("emailTemplate")
	t, err = t.Parse(config.templateEmail)
	if err != nil {
		return err
	}
	err = t.Execute(&doc, context)
	if err != nil {
		return err
	}
	err = smtp.SendMail(config.ServerMail, auth, config.activeEmail, []string{recipients}, doc.Bytes())
	return err
}
