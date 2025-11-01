package repository

import (
	"encoding/json"
	"fmt"

	"invoice/internal/invoice"
)

// ClientData represents client/customer data as stored on disk
type ClientData struct {
	invoice.Party
	TaxRate     float32 `json:"tax_rate"`
	EmailTarget string  `json:"email_target"`
}

// FromJSON deserializes client data from JSON
func (c *ClientData) FromJSON(data []byte) error {
	if err := json.Unmarshal(data, c); err != nil {
		return fmt.Errorf("failed to unmarshal client JSON: %v", err)
	}
	return nil
}

// ProviderData represents service provider data as stored on disk
type ProviderData struct {
	invoice.Party
	Payment invoice.PaymentInfo `json:"payment_info"`
}

// FromJSON deserializes provider data from JSON
func (p *ProviderData) FromJSON(data []byte) error {
	if err := json.Unmarshal(data, p); err != nil {
		return fmt.Errorf("failed to unmarshal provider JSON: %v", err)
	}
	return nil
}

// InvoiceData combines provider and client data for creating invoices
type InvoiceData struct {
	Provider ProviderData `json:"provider"`
	Client   ClientData   `json:"client"`
}

// ToInvoice creates a new Invoice from the stored data
func (d *InvoiceData) ToInvoice() *invoice.Invoice {
	inv := invoice.New(d.Provider.Party, d.Client.Party, d.Provider.Payment, d.Client.TaxRate)
	inv.SetEmailTarget(d.Client.EmailTarget)
	return inv
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
