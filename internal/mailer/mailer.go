package mailer

import (
	"fmt"
	"html/template"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/Stuhub-io/core/domain"
	"github.com/Stuhub-io/core/ports"
	"github.com/Stuhub-io/logger"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type Mailer struct {
	from      string
	address   string
	endpoint  string
	host      string
	clientKey string
	logger    logger.Logger
}

type NewMailerParams struct {
	From      string
	Address   string
	ClientKey string
	Logger    logger.Logger
}

func (m *Mailer) addEndpoint(endpoint string) *Mailer {
	m.endpoint = endpoint
	return m
}

func (m *Mailer) addHost(host string) *Mailer {
	m.host = host
	return m
}

func NewMailer(params NewMailerParams) ports.Mailer {
	return &Mailer{
		from:      params.From,
		address:   params.Address,
		endpoint:  "/v3/mail/send",
		host:      "https://api.sendgrid.com",
		clientKey: params.ClientKey,
		logger:    params.Logger,
	}
}

func (m *Mailer) SendMail(payload ports.SendSendGridMailPayload) *domain.Error {
	from := m.from
	if payload.FromName != nil {
		from = *payload.FromName
	}

	v3Mail := mail.NewV3Mail()
	v3Mail.SetFrom(mail.NewEmail(from, m.address))
	v3Mail.SetTemplateID(payload.TemplateId)
	v3Mail.Subject = payload.Subject

	p := mail.NewPersonalization()
	for name := range payload.Data {
		p.SetDynamicTemplateData(name, payload.Data[name])
	}
	p.AddTos(mail.NewEmail(payload.ToName, payload.ToAddress))
	v3Mail.AddPersonalizations(p)

	request := sendgrid.GetRequest(
		m.clientKey,
		m.endpoint,
		m.host,
	)
	request.Method = "POST"
	request.Body = mail.GetRequestBody(v3Mail)
	_, err := sendgrid.API(request)
	if err != nil {
		m.logger.Error(err, err.Error())
		return domain.ErrSendMail
	}

	return nil
}

func (m *Mailer) SendMailCustomTemplate(
	payload ports.SendSendGridMailCustomTemplatePayload,
) *domain.Error {
	from := m.from
	if payload.FromName != nil {
		from = *payload.FromName
	}

	v3Mail := mail.NewV3Mail()
	v3Mail.SetFrom(mail.NewEmail(from, m.address))
	v3Mail.Subject = payload.Subject

	htmlContent, perr := m.parseHTMLTemplateFile(payload.TemplateHTMLName, payload.Data)
	if perr != nil {
		return perr
	}
	content := mail.NewContent(
		"text/html",
		htmlContent,
	)

	v3Mail.AddContent(content)

	p := mail.NewPersonalization()
	p.AddTos(mail.NewEmail(payload.ToName, payload.ToAddress))
	v3Mail.AddPersonalizations(p)

	request := sendgrid.GetRequest(
		m.clientKey,
		m.endpoint,
		m.host,
	)
	request.Method = "POST"
	request.Body = mail.GetRequestBody(v3Mail)
	resp, err := sendgrid.API(request)
	if err != nil {
		m.logger.Error(err, err.Error())
		return domain.ErrSendMail
	}

	logger.L.Debug("Email sent successfully: " + resp.Body)

	return nil
}

func (m *Mailer) parseHTMLTemplateFile(
	templateName string,
	data interface{},
) (string, *domain.Error) {
	templatesDir := "./templates"

	_, currentFile, _, _ := runtime.Caller(0)
	currentDir := filepath.Dir(currentFile)

	templatePath := filepath.Join(currentDir, templatesDir, fmt.Sprintf("%s.html", templateName))
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return "", domain.ErrInvalidTemplate
	}

	var builder strings.Builder
	err = tmpl.Execute(&builder, data)
	if err != nil {
		return "", domain.ErrInvalidTemplate
	}

	return builder.String(), nil
}
