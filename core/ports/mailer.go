package ports

type SendMailPayload struct {
	To          string
	Address     string
	Subject     string
	PlainText   string
	HTMLContent string
}

type Mailer interface {
	SendMail(payload SendMailPayload) error
}

//mailTempl := ""
