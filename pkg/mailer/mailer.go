package mailer

import (
	"context"
	"crypto/tls"
	"fmt"

	gomail "gopkg.in/gomail.v2"
)

type Mailer struct {
	Dialer   *gomail.Dialer
	From     string // no-reply@mediawatch.io
	FromName string // MediaWatch
	ReplyTo  string // press@mediawatch.io
}

func New(smtp string, port int, username, password string, From, FromName, ReplyTo string) *Mailer {
	d := gomail.NewDialer(smtp, port, username, password)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	m := &Mailer{Dialer: d}
	m.From = From
	m.FromName = FromName
	m.ReplyTo = ReplyTo
	return m
}

func Message(ctx context.Context, m *Mailer, To, subject, body string) error {
	gm := gomail.NewMessage()
	// gm.SetAddressHeader("Cc", "","")
	gm.SetAddressHeader("From", m.From, m.FromName)
	gm.SetHeader("To", To)
	gm.SetHeader("Reply-To", m.ReplyTo)
	gm.SetHeader("Subject", subject)
	gm.SetBody("text/html", body)

	return m.Dialer.DialAndSend(gm)
}

func MessageSimple(ctx context.Context, m *Mailer, To, subject, body string) error {
	msgBody := fmt.Sprintf(msgDefault, body)
	return Message(ctx, m, To, subject, msgBody)
}

func SendInvite(ctx context.Context, m *Mailer, To, First, Last, email string) error {
	msgBody := fmt.Sprintf(msgInvitation, First, Last, email)
	return Message(ctx, m, To, fmt.Sprintf("Join MediaWatch (Invitation by \"%s\")", First), msgBody)
}

func SendNewPass(ctx context.Context, m *Mailer, To, First, pass string) error {
	msgBody := fmt.Sprintf(msgNewPass, First, To)
	return Message(ctx, m, To, "Password Reset", msgBody)
}

func SendReset(ctx context.Context, m *Mailer, To, First, pin, id string) error {
	msgBody := fmt.Sprintf(msgReset, First, pin, id, id)
	return Message(ctx, m, To, "Your Verification Code", msgBody)
}

func SendPin(ctx context.Context, m *Mailer, To, First, pin, id string) error {
	msgBody := fmt.Sprintf(msgPin, First, pin, id, id)
	return Message(ctx, m, To, "Your Verification Code", msgBody)
}
