package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"invoice/internal/invoice"
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
	Id string `json:"id,omitempty"`
	invoice.Party
	TaxRate     float32 `json:"tax_rate"`
	EmailTarget string  `json:"email_target"`
}

// FromJSON deserializes client data from JSON
func NewClientDataFromJSON(data []byte) (*ClientData, error) {
	var c ClientData
	if err := json.Unmarshal(data, &c); err != nil {
		return nil, fmt.Errorf("failed to unmarshal client JSON: %v", err)
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
	Id string `json:"id,omitempty"`
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

// EmailConfig represents email configuration as stored on disk
type EmailConfig struct {
	ID         string `json:"id"`
	SMTPConfig struct {
		Host     string `json:"host"`
		Port     string `json:"port"`
		Username string `json:"username"`
		Password string `json:"password"`
	} `json:"smtp"`
	DefaultMessage struct {
		Subject string `json:"subject"`
		Body    string `json:"body"`
		Footer  string `json:"footer"`
	} `json:"default_message"`
	Recipients []string `json:"recipients"`
}

// FromJSON deserializes email config from JSON
func (e *EmailConfig) FromJSON(data []byte) error {
	if err := json.Unmarshal(data, e); err != nil {
		return fmt.Errorf("failed to unmarshal email config JSON: %v", err)
	}
	return nil
}
