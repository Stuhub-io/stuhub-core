package ports

import "github.com/Stuhub-io/core/domain"

type SendSendGridMailPayload struct {
	FromName    string
	FromAddress string
	ToName      string
	ToAddress   string
	TemplateId  string
	Data        map[string]string
	Subject     string
	Content     string
}

type Mailer interface {
	SendMail(payload SendSendGridMailPayload) *domain.Error
}
