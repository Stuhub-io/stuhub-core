package ports

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
	SendMail(payload SendSendGridMailPayload) error
}
