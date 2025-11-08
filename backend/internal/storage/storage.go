package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"go-invoice/internal/invoice"
)

type StorageDir struct {
	Root      string
	Clients   string
	Providers string
	Invoices  string
	Config    string
}

func NewStorageDir() (*StorageDir, error) {
	exePath, err := os.Executable()
	if err != nil {
		return nil, fmt.Errorf("unable to obtain executable path: %v", err)
	}
	dir := filepath.Dir(exePath)
	root := filepath.Join(dir, "db")

	clientsDir := filepath.Join(root, "clients")
	providersDir := filepath.Join(root, "providers")
	invoicesDir := filepath.Join(root, "invoices")
	configDir := filepath.Join(root, "smtp")
	if err := ensurePathExist(clientsDir); err != nil {
		return nil, err
	}
	if err := ensurePathExist(providersDir); err != nil {
		return nil, err
	}
	if err := ensurePathExist(invoicesDir); err != nil {
		return nil, err
	}
	if err := ensurePathExist(configDir); err != nil {
		return nil, err
	}

	return &StorageDir{
		Root:      root,
		Clients:   clientsDir,
		Providers: providersDir,
		Invoices:  invoicesDir,
		Config:    configDir,
	}, nil
}

func ensurePathExist(path string) error {
	if _, err := os.Stat(path); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("path not exist '%s': %v", path, err)
		}
	}
	return nil
}

// ClientData represents client/customer data as stored on disk
type ClientData struct {
	invoice.Party
	TaxRate         float32 `json:"tax_rate"`
	EmailTarget     string  `json:"email_target"`
	EmailTemplateId string  `json:"email_template_id"`
}

// FromJSON deserializes client data from JSON
func NewClientDataFromJSON(data []byte) (*ClientData, error) {
	var c ClientData
	if err := json.Unmarshal(data, &c); err != nil {
		return nil, fmt.Errorf("failed to unmarshal client JSON: %v", err)
	}
	if c.EmailTemplateId == "" {
		c.EmailTemplateId = "default"
	}
	return &c, nil
}

func (c *ClientData) SetID(id string) {
	c.Id = id
}

func (c *ClientData) HasRequiredFields() bool {
	return c.Party.HasRequiredFields()
}

// ProviderData represents service provider data as stored on disk
type ProviderData struct {
	invoice.Party
	Payment invoice.PaymentInfo `json:"payment_info"`
}

// NewProviderDataFromJSON deserializes provider data from JSON
func NewProviderDataFromJSON(data []byte) (*ProviderData, error) {
	var p ProviderData
	if err := json.Unmarshal(data, &p); err != nil {
		return nil, fmt.Errorf("failed to unmarshal provider JSON: %v", err)
	}
	return &p, nil
}

func (p *ProviderData) SetID(id string) {
	p.Id = id
}

func (p *ProviderData) HasRequiredFields() bool {
	return p.Party.HasRequiredFields() && p.Payment.HasRequiredFields()
}

type EmailTemplateData struct {
	Id      string `json:"id"`
	Name    string `json:"name"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

func NewDefaultEmailTemplateData() *EmailTemplateData {
	return &EmailTemplateData{
		Id:      "default",
		Name:    "Default Invoice Email",
		Subject: "Invoice from ${PROVIDER_NAME} - ${INVOICE_ID}",
		Body:    "Please find the attached invoice for the services rendered.\n\nClient name: ${CLIENT_NAME}\nSubcontractor email: ${PROVIDER_EMAIL}\nService type: ${SERVICE_TYPE}\n\nKind regards,\n${PROVIDER_NAME}",
	}
}
