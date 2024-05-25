package email

import (
	"github.com/wneessen/go-mail"
	"log"
)

type EmailClient struct {
	c         *mail.Client
	thisEmail string
}

func NewEmailClient(server string, port int, username string, password string) (*EmailClient, error) {
	obj := &EmailClient{
		thisEmail: username,
	}
	cc, err := mail.NewClient(server, mail.WithPort(port), mail.WithSMTPAuth(mail.SMTPAuthPlain),
		mail.WithUsername(username), mail.WithPassword(password))
	if err != nil {
		log.Fatalf("failed to create mail client: %s", err)
	}
	obj.c = cc
	return obj, err
}

func (s *EmailClient) SendEmail(to, subject, body string) error {
	m := mail.NewMsg()
	if err := m.From(s.thisEmail); err != nil {
		return err
	}
	if err := m.To(to); err != nil {
		return err
	}
	m.Subject(subject)
	m.SetBodyString(mail.TypeTextHTML, body)

	if err := s.c.DialAndSend(m); err != nil {
		return err
	}
	return nil
}
