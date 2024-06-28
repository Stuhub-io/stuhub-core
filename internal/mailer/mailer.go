package mailer

import (
	"github.com/Stuhub-io/config"
	"github.com/Stuhub-io/core/ports"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type Mailer struct {
	name      string
	address   string
	clientKey string
	config    config.Config
}

type NewMailerParams struct {
	Name      string
	Address   string
	ClientKey string
	config.Config
}

func NewMailer(params NewMailerParams) *Mailer {
	return &Mailer{
		name:      params.Name,
		address:   params.Address,
		clientKey: params.ClientKey,
		config:    params.Config,
	}
}

func (m *Mailer) SendMail(payload ports.SendSendGridMailPayload) error {
	from := mail.NewEmail(payload.FromName, payload.FromAddress)
	subject := payload.Subject
	to := mail.NewEmail(payload.ToName, payload.ToAddress)
	content := mail.NewContent("text/html", payload.Content)

	sender := mail.NewV3MailInit(from, subject, to, content)
	for name := range payload.Data {
		sender.Personalizations[0].SetSubstitution(name, payload.Data[name])
	}
	sender.Personalizations[0].SetSubstitution("-name-", "Example User")
	sender.Personalizations[0].SetSubstitution("-city-", "Denver")
	sender.SetTemplateID(payload.TemplateId)

	request := sendgrid.GetRequest(m.config.SendgridKey, "/v3/mail/send", "https://api.sendgrid.com")
	request.Method = "POST"
	request.Body = mail.GetRequestBody(sender)
	_, err := sendgrid.API(request)
	return err
}
