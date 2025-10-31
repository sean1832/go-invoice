package storage

import (
	"encoding/json"
	"fmt"
	"invoice/internal/core/model/view"
)

// DataModel represents the storage model for invoices
type DataModel struct {
	Provider Provider `json:"provider"`
	Client   Client   `json:"client"`
}

func NewModel(provider Provider, client Client) *DataModel {
	return &DataModel{
		Provider: provider,
		Client:   client,
	}
}

func (m *DataModel) ToViewModel() *view.ViewModel {
	return view.NewInvoiceData(m.Provider.Party, m.Client.Party, m.Provider.Payment, m.Client.TaxRate)
}

// Client represents the client/customer party
type Client struct {
	view.Party
	TaxRate     float32 `json:"tax_rate"`
	EmailTarget string  `json:"email_target"`
}

func NewClientFromJSON(data []byte) (*Client, error) {
	var c Client
	if err := json.Unmarshal(data, &c); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %v", err)
	}
	return &c, nil
}

// Provider represents the service provider party
type Provider struct {
	view.Party
	Payment view.PaymentInfo `json:"payment_info"`
}

func NewProviderFromJSON(data []byte) (*Provider, error) {
	var p Provider
	if err := json.Unmarshal(data, &p); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %v", err)
	}
	return &p, nil
}
