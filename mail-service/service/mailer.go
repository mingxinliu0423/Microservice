package main

import (
	"bytes"
	"html/template"
	"log"
	"time"

	mail "github.com/xhit/go-simple-mail/v2"
)

// Message struct represents an email message.
type Message struct {
	From        string
	FromName    string
	To          string
	Subject     string
	Attachments []string
	Data        any
	DataMap     map[string]any
}

// Mail struct holds the configuration for sending emails.
type Mail struct {
	Domain      string
	Host        string
	Port        int
	Username    string
	Password    string
	Encryption  string
	FromAddress string
	FromName    string
}

// buildPlainTextMessage generates the plain text version of the email message using a template.
func (m *Mail) buildPlainTextMessage(msg Message) (string, error) {
	templateToRender := "./templates/mail.plain.gohtml"

	t, err := template.New("email-plain").ParseFiles(templateToRender)
	if err != nil {
		return "", err
	}

	var tpl bytes.Buffer
	if err = t.ExecuteTemplate(&tpl, "body", msg.DataMap); err != nil {
		return "", err
	}

	plainMessage := tpl.String()

	return plainMessage, nil
}

// SendSMTPMessage sends an email message using SMTP.
func (m *Mail) SendSMTPMessage(msg Message) error {
	// Set default 'From' and 'FromName' if not provided in the message.
	if msg.From == "" {
		msg.From = m.FromAddress
	}

	if msg.FromName == "" {
		msg.FromName = m.FromName
	}

	// Prepare data for the email template.
	data := map[string]any{
		"message": msg.Data,
	}

	msg.DataMap = data

	plainMessage, err := m.buildPlainTextMessage(msg)
	if err != nil {
		return err
	}

	// Configure the SMTP client.
	server := mail.NewSMTPClient()
	server.Host = m.Host
	server.Port = m.Port
	server.Username = m.Username
	server.Password = m.Password
	server.Encryption = mail.EncryptionNone
	server.KeepAlive = false
	server.ConnectTimeout = 10 * time.Second
	server.SendTimeout = 10 * time.Second

	// Connect to the SMTP server.
	smtpClient, err := server.Connect()
	if err != nil {
		log.Println(err)
		return err
	}

	// Create the email message.
	email := mail.NewMSG()
	email.SetFrom(msg.From).
		AddTo(msg.To).
		SetSubject(msg.Subject)

	email.SetBody(mail.TextPlain, plainMessage)
	// email.AddAlternative(mail.TextHTML, formattedMessage)

	// Add attachments, if any.
	if len(msg.Attachments) > 0 {
		for _, x := range msg.Attachments {
			email.AddAttachment(x)
		}
	}

	// Send the email.
	err = email.Send(smtpClient)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
