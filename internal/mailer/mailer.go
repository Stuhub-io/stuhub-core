package mailer

import (
	"os"

	"github.com/Stuhub-io/config"
	"github.com/Stuhub-io/core/domain"
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

func NewMailer(params NewMailerParams) ports.Mailer {
	return &Mailer{
		name:      params.Name,
		address:   params.Address,
		clientKey: params.ClientKey,
		config:    params.Config,
	}
}

func (m *Mailer) SendMail(payload ports.SendSendGridMailPayload) *domain.Error {
	v3Mail := mail.NewV3Mail()
	from := mail.NewEmail(payload.FromName, payload.FromAddress)
	v3Mail.SetFrom(from)
	v3Mail.SetTemplateID(payload.TemplateId)

	p := mail.NewPersonalization()
	for name := range payload.Data {
		p.SetDynamicTemplateData(name, payload.Data[name])
	}
	p.AddTos(mail.NewEmail(payload.ToName, payload.ToAddress))
	v3Mail.AddPersonalizations(p)

	request := sendgrid.GetRequest(os.Getenv("SENDGRID_API_KEY"), "/v3/mail/send", "https://api.sendgrid.com")
	request.Method = "POST"
	request.Body = mail.GetRequestBody(v3Mail)
	_, err := sendgrid.API(request)
	if err != nil {
		return domain.ErrSendMail
	}

	return nil
}
