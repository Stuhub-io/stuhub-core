package mailer

import (
	"github.com/Stuhub-io/core/ports"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type Mailer struct {
	name      string
	address   string
	clientKey string
}

type NewMailerParams struct {
	Name      string
	Address   string
	ClientKey string
}

func NewMailer(params NewMailerParams) *Mailer {
	return &Mailer{
		name:      params.Name,
		address:   params.Address,
		clientKey: params.ClientKey,
	}
}

func (m *Mailer) SendMail(payload ports.SendMailPayload) error {
	from := mail.NewEmail(m.name, m.address)
	subject := payload.Subject
	to := mail.NewEmail(payload.To, payload.Address)
	plainTextContent := payload.PlainText
	htmlContent := payload.HTMLContent
	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
	client := sendgrid.NewSendClient(m.clientKey)
	_, err := client.Send(message)

	return err
}
