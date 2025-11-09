package types

type EmailMessage struct {
	To      []string `json:"to"`
	Subject string   `json:"subject"`
	Body    string   `json:"body"`
}

func NewEmailMessage(to []string, subject, body string) EmailMessage {
	return EmailMessage{
		To:      to,
		Subject: subject,
		Body:    body,
	}
}
