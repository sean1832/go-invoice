package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"go-invoice/internal/invoice"
)

type StorageDir struct {
	Root           string
	Clients        string
	Providers      string
	Invoices       string
	Config         string
	EmailTemplates string
}

// NewStorageDir initializes the storage directory structure.
// It takes a rootDir as an argument, making the function testable and configurable.
func NewStorageDir(rootDir string) (*StorageDir, error) {
	// Define all paths relative to the configurable root.
	storage := &StorageDir{
		Root:           rootDir,
		Clients:        filepath.Join(rootDir, "clients"),
		Providers:      filepath.Join(rootDir, "providers"),
		Invoices:       filepath.Join(rootDir, "invoices"),
		EmailTemplates: filepath.Join(rootDir, "email_templates"),
	}

	// Create a list of all paths that must exist.
	paths := []string{
		storage.Root,
		storage.Clients,
		storage.Providers,
		storage.Invoices,
		storage.EmailTemplates,
	}

	// Loop and create each one, using the correct tool (MkdirAll).
	// Use 0755 (rwxr-xr-x) instead of 0777 (os.ModePerm).
	const perm = 0755

	for _, path := range paths {
		// os.MkdirAll is idempotent: it does nothing if the path already exists.
		if err := os.MkdirAll(path, perm); err != nil {
			return nil, fmt.Errorf("failed to create storage directory %q: %v", path, err)
		}
	}

	// Create default email template if it doesn't exist
	defaultTemplatePath := filepath.Join(storage.EmailTemplates, "default.json")
	if _, err := os.Stat(defaultTemplatePath); os.IsNotExist(err) {
		defaultTemplate := NewDefaultEmailTemplateData()
		if err := defaultTemplate.SaveToFile(defaultTemplatePath); err != nil {
			return nil, fmt.Errorf("failed to create default email template: %v", err)
		}
	}

	return storage, nil
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

type EmailTemplate struct {
	Id      string `json:"id"`
	Name    string `json:"name"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

func NewDefaultEmailTemplateData() *EmailTemplate {
	return &EmailTemplate{
		Id:      "default",
		Name:    "Default Invoice Email",
		Subject: "Invoice from {{PROVIDER_NAME}} ({{INVOICE_ID}})",
		Body:    "Please find the attached invoice for the services rendered.\n\nClient name: {{CLIENT_NAME}}\nSubcontractor email: {{PROVIDER_EMAIL}}\nService type: {{SERVICE_TYPE}}\n\nKind regards,\n{{PROVIDER_NAME}}",
	}
}

// SaveToFile saves the email template to a JSON file
func (et *EmailTemplate) SaveToFile(filePath string) error {
	data, err := json.MarshalIndent(et, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal email template to JSON: %v", err)
	}
	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write email template to file: %v", err)
	}

	return nil
}
