package smtp

type recipient struct {
	Name    string `json:"name,omitempty"`
	Address string `json:"address"`
}

type EmailMessage struct {
	From    recipient   `json:"from"`
	To      []recipient `json:"to"`
	Cc      []recipient `json:"cc,omitempty"`
	Bcc     []recipient `json:"bcc,omitempty"`
	Subject string      `json:"subject"`
	Body    string      `json:"body"`
}

type SMTPSender struct {
	Host     string
	Port     string
	Username string
	Password string
}

func NewSmtpSender(host, port, username, password string) *SMTPSender {
	return &SMTPSender{
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,
	}
}
