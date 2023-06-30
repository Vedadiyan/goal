package smtpclient

import (
	"bytes"
	"fmt"
	"net/smtp"
	"strings"
)

type SmtpClient struct {
	server   string
	port     int
	username string
	password string
}

func New(server string, port int, username string, password string) *SmtpClient {
	return &SmtpClient{
		server:   server,
		port:     port,
		username: username,
		password: password,
	}
}

func (smtpClient SmtpClient) sendBase(from string, subject string, to string, cc []string, bcc []string, body []byte, options func(msg *bytes.Buffer)) error {
	recepients := make([]string, 0, len(cc)+len(bcc))
	recepients = append(recepients, to)
	recepients = append(recepients, cc...)
	recepients = append(recepients, bcc...)
	var msg bytes.Buffer
	msg.WriteString(fmt.Sprintf("From: %s\r\n", from))
	msg.WriteString(fmt.Sprintf("To: %s\r\n", to))
	if len(cc) > 0 {
		msg.WriteString(fmt.Sprintf("Cc: %s\r\n", strings.Join(cc, ";")))
	}
	if len(bcc) > 0 {
		msg.WriteString(fmt.Sprintf("Bcc: %s\r\n", strings.Join(bcc, ";")))
	}
	msg.WriteString(fmt.Sprintf("Subject: %s\r\n", subject))
	if options != nil {
		options(&msg)
	}
	msg.Write(body)
	return smtp.SendMail(
		fmt.Sprintf("%s:%d", smtpClient.server, smtpClient.port),
		smtp.PlainAuth("", smtpClient.username, smtpClient.password, smtpClient.server),
		from,
		recepients,
		msg.Bytes(),
	)
}

func (smtpClient SmtpClient) Send(from string, subject string, to string, cc []string, bcc []string, body []byte) error {
	return smtpClient.sendBase(from, subject, to, cc, bcc, body, nil)
}
func (smtpClient SmtpClient) SendHtml(from string, subject string, to string, cc []string, bcc []string, body []byte) error {
	return smtpClient.sendBase(from, subject, to, cc, bcc, body, func(msg *bytes.Buffer) {
		msg.WriteString("MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n")
	})
}
