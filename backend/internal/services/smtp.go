package services

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"mime/multipart"
	"net/smtp"
	"net/textproto"
	"strings"
)

// SMTPService is responsible for sending emails via direct SMTP
type SMTPService struct {
	From     string
	Host     string
	Port     int
	Password string
}

// NewSMTPService creates a new SMTPService instance
func NewSMTPService(from, host string, port int, password string) *SMTPService {
	return &SMTPService{
		From:     from,
		Host:     host,
		Port:     port,
		Password: password,
	}
}

// Send sends a basic plaintext email via SMTP
func (s *SMTPService) Send(to []string, subject string, body string) error {
	auth := smtp.PlainAuth("", s.From, s.Password, s.Host)
	address := fmt.Sprintf("%s:%d", s.Host, s.Port)

	header := make(map[string]string)
	header["From"] = s.From
	header["To"] = strings.Join(to, ", ")
	header["Subject"] = subject
	header["MIME-version"] = "1.0"
	header["Content-Type"] = "text/plain; charset=utf-8"

	var msgBuilder strings.Builder
	for k, v := range header {
		msgBuilder.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
	}
	msgBuilder.WriteString("\r\n") // empty line seperates header from body
	msgBuilder.WriteString(body)

	msg := msgBuilder.String()

	err := smtp.SendMail(
		address,
		auth,
		s.From,
		to,
		[]byte(msg),
	)
	if err != nil {
		return fmt.Errorf("smpt.SendMail failed: %v", err)
	}
	return nil
}

// AttachmentType represents the MIME type of the attachment
type AttachmentType string

const (
	AttachmentTypePDF       AttachmentType = "application/pdf"          // pdf file
	AttachmentTypeZIP       AttachmentType = "application/zip"          // zip file
	AttachmentTypeJSON      AttachmentType = "application/json"         // json file
	AttachmentTypeGeneric   AttachmentType = "application/octet-stream" // generic binary
	AttachmentTypeImageJPEG AttachmentType = "image/jpeg"               // jpeg image
	AttachmentTypeImagePNG  AttachmentType = "image/png"                // png image
)

// SendWithAttachment sends an email with an attachment via SMTP
func (s *SMTPService) SendWithAttachment(
	to []string,
	subject string,
	body string,
	attachmentName string,
	attachmentData []byte,
	attachmentType AttachmentType) error {

	address := fmt.Sprintf("%s:%d", s.Host, s.Port)
	auth := smtp.PlainAuth("", s.From, s.Password, s.Host)

	var emailBuffer bytes.Buffer

	mw := multipart.NewWriter(&emailBuffer)

	header := make(map[string]string)
	header["From"] = s.From
	header["To"] = strings.Join(to, ", ")
	header["Subject"] = subject
	header["MIME-Version"] = "1.0"

	header["Content-Type"] = fmt.Sprintf("multipart/mixed; boundary=%s", mw.Boundary())
	for k, v := range header {
		emailBuffer.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
	}

	emailBuffer.WriteString("\r\n")

	// first part: plain text body
	partHeaders := textproto.MIMEHeader{}
	partHeaders.Set("Content-Type", "text/plain; charset=utf-8")
	part, err := mw.CreatePart(partHeaders)
	if err != nil {
		return err
	}
	part.Write([]byte(body))

	// second part: attachment
	partHeaders = textproto.MIMEHeader{}
	partHeaders.Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", attachmentName))
	partHeaders.Set("Content-Type", string(attachmentType))
	partHeaders.Set("Content-Transfer-Encoding", "base64")

	part, err = mw.CreatePart(partHeaders)
	if err != nil {
		return err
	}

	b64Encoder := base64.NewEncoder(base64.StdEncoding, part)
	b64Encoder.Write(attachmentData)
	b64Encoder.Close()

	// finish multipart message
	if err := mw.Close(); err != nil { // <-- this writes the final boundary
		return err
	}

	// send the email
	err = smtp.SendMail(
		address,
		auth,
		s.From,
		to,
		emailBuffer.Bytes(),
	)

	if err != nil {
		return fmt.Errorf("smtp.SendMail failed: %v", err)
	}

	return nil
}
