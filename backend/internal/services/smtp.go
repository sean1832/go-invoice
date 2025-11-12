package services

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"go-invoice/internal/auth"
	"mime/multipart"
	"net/smtp"
	"net/textproto"
	"strings"
)

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

// ========== OAuth2 SMTP Auth ==========

type oauth2Auth struct {
	username, accessToken string
}

func newOAuth2Auth(username, accessToken string) smtp.Auth {
	return &oauth2Auth{username, accessToken}
}
func (a *oauth2Auth) Start(server *smtp.ServerInfo) (string, []byte, error) {
	authStr := fmt.Sprintf("user=%s\x01auth=Bearer %s\x01\x01", a.username, a.accessToken)
	return "XOAUTH2", []byte(authStr), nil
}
func (a *oauth2Auth) Next(fromServer []byte, more bool) ([]byte, error) {
	if more {
		return nil, fmt.Errorf("[auth] unexpected server challenge: %s", fromServer)
	}
	return nil, nil
}

func (a *oauth2Auth) Close() error {
	return nil
}

// ========== Email struct ==========

// SMTPService is responsible for sending emails via direct SMTP
type SMTPService struct {
	from    string
	address string
	auth    smtp.Auth
}

// NewSMTPService creates a new SMTPService instance
func NewSMTPService(from, host string, port int, credential string, authMethod auth.AuthMethod) *SMTPService {
	address := fmt.Sprintf("%s:%d", host, port)
	switch authMethod {
	case auth.AuthMethodPlain:
		auth := smtp.PlainAuth("", from, credential, host)
		return &SMTPService{
			from:    from,
			address: address,
			auth:    auth,
		}
	case auth.AuthMethodOAuth2:
		auth := newOAuth2Auth(from, credential)
		return &SMTPService{
			from:    from,
			address: address,
			auth:    auth,
		}
	default:
		return nil
	}
}

// Send sends a basic plaintext email via SMTP
func (s *SMTPService) Send(to []string, subject string, body string) error {
	header := make(map[string]string)
	header["From"] = s.from
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
		s.address,
		s.auth,
		s.from,
		to,
		[]byte(msg),
	)
	if err != nil {
		return fmt.Errorf("smpt.SendMail failed: %v", err)
	}
	return nil
}

// SendWithAttachment sends an email with an attachment via SMTP
func (s *SMTPService) SendWithAttachment(
	to []string,
	subject string,
	body string,
	attachmentName string,
	attachmentData []byte,
	attachmentType AttachmentType) error {

	var emailBuffer bytes.Buffer

	mw := multipart.NewWriter(&emailBuffer)

	header := make(map[string]string)
	header["From"] = s.from
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
		s.address,
		s.auth,
		s.from,
		to,
		emailBuffer.Bytes(),
	)

	if err != nil {
		return fmt.Errorf("smtp.SendMail failed: %v", err)
	}

	return nil
}
