package ports

import "github.com/Stuhub-io/core/domain"

type SendSendGridMailPayload struct {
	FromName   string
	ToName     string
	ToAddress  string
	TemplateId string
	Data       map[string]string
	Subject    string
	Content    string
}

type SendSendGridMailCustomTemplatePayload struct {
	FromName         string
	ToName           string
	ToAddress        string
	TemplateHTMLName string
	Data             map[string]string
	Subject          string
}

type Mailer interface {
	SendMail(payload SendSendGridMailPayload) *domain.Error
	SendMailCustomTemplate(payload SendSendGridMailCustomTemplatePayload) *domain.Error
}
