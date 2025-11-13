package invoice

import (
	"errors"
	"go-invoice/internal/types"
)

type InvoiceStatus string

const (
	StatusDraft InvoiceStatus = "draft"
	StatusSent  InvoiceStatus = "send"
)

// Invoice represents the core invoice domain model
type Invoice struct {
	ID              string        `json:"id"`                     // invoice number/identifier
	Status          InvoiceStatus `json:"status"`                 // invoice status (draft, sent)
	Date            types.Date    `json:"date"`                   // invoice date
	Due             types.Date    `json:"due"`                    // payment due date
	Provider        Party         `json:"provider"`               // service provider
	Client          Party         `json:"client"`                 // client/customer
	Items           []ServiceItem `json:"items"`                  // list of services/products
	Pricing         Pricing       `json:"pricing"`                // pricing details
	Payment         PaymentInfo   `json:"payment"`                // payment information
	EmailTarget     string        `json:"email_target,omitempty"` // (optional) email target for sending the invoice
	EmailTemplateID string        `json:"email_template_id"`      // email template ID
}

// SetEmailTarget sets the email address to send the invoice to
func (inv *Invoice) SetEmailTarget(email string) {
	inv.EmailTarget = email
}

func (inv *Invoice) SetID(id string) {
	inv.ID = id
}

func (inv *Invoice) HasRequiredFields() bool {
	return inv.Status != "" && inv.Provider.HasRequiredFields() && inv.Client.HasRequiredFields() && inv.Payment.HasRequiredFields()
}

// AddItem adds a service item to the invoice and updates the pricing
func (inv *Invoice) AddItem(item ServiceItem) {
	inv.Items = append(inv.Items, item)
	inv.Pricing.Update(calculateSubtotal(inv.Items))
}

// Party represents either the service provider or the client/customer
type Party struct {
	Id      string `json:"id"`                // (optional) unique identifier
	Name    string `json:"name"`              // party name
	Address string `json:"address,omitempty"` // (optional) address
	Email   string `json:"email,omitempty"`   // (optional) email address
	Phone   string `json:"phone,omitempty"`   // (optional) phone number
	ABN     string `json:"abn,omitempty"`     // (optional) Australian Business Number
}

func (p *Party) HasRequiredFields() bool {
	return p.Name != "" && p.Id != ""
}

// ServiceItem represents a single line item in the invoice
type ServiceItem struct {
	Date              types.Date `json:"date"`                         // date of service/product
	Description       string     `json:"description"`                  // description of service/product
	DescriptionDetail string     `json:"description_detail,omitempty"` // (optional) detailed description
	Quantity          float32    `json:"quantity"`                     // quantity provided
	UnitPrice         float32    `json:"unit_price"`                   // price per unit
	TotalPrice        float32    `json:"total_price"`                  // total price (Quantity * UnitPrice)
}

// NewServiceItem creates a new service item
func NewServiceItem(date types.Date, description string, quantity, unitPrice float32) ServiceItem {
	return ServiceItem{
		Date:        date,
		Description: description,
		Quantity:    quantity,
		UnitPrice:   unitPrice,
		TotalPrice:  quantity * unitPrice,
	}
}

// NewServiceItemWithDetail creates a new service item with detailed description
func NewServiceItemWithDetail(date types.Date, description, descriptionDetail string, quantity, unitPrice float32) ServiceItem {
	return ServiceItem{
		Date:              date,
		Description:       description,
		DescriptionDetail: descriptionDetail,
		Quantity:          quantity,
		UnitPrice:         unitPrice,
		TotalPrice:        quantity * unitPrice,
	}
}

// Pricing holds the pricing breakdown of the invoice
type Pricing struct {
	Subtotal  float32 `json:"subtotal"` // subtotal before tax
	TaxAmount float32 `json:"tax"`      // tax amount
	TaxRate   float32 `json:"tax_rate"` // tax rate (percentage)
	Total     float32 `json:"total"`    // total amount (subtotal + tax)
}

// NewPricing creates a new pricing structure
func NewPricing(subtotal, taxRate float32) (*Pricing, error) {
	if subtotal < 0 {
		return nil, errors.New("subtotal cannot be negative")
	}
	if taxRate < 0 {
		return nil, errors.New("tax rate cannot be negative")
	}

	taxAmount := subtotal * taxRate / 100
	total := subtotal + taxAmount
	return &Pricing{
		Subtotal:  subtotal,
		TaxAmount: taxAmount,
		TaxRate:   taxRate,
		Total:     total,
	}, nil
}

// Update recalculates the pricing based on a new subtotal
func (p *Pricing) Update(subtotal float32) {
	p.Subtotal = subtotal
	p.TaxAmount = subtotal * p.TaxRate / 100
	p.Total = p.Subtotal + p.TaxAmount
}

// calculateSubtotal calculates the sum of all service item prices
func calculateSubtotal(items []ServiceItem) float32 {
	var subtotal float32
	for _, item := range items {
		subtotal += item.TotalPrice
	}
	return subtotal
}

// PaymentInfo holds the payment details for the invoice
type PaymentInfo struct {
	Method        string `json:"method"`         // payment method (e.g., bank transfer, credit card)
	AccountName   string `json:"account_name"`   // account holder name
	BSB           string `json:"bsb"`            // bank state branch number
	AccountNumber string `json:"account_number"` // bank account number
}

func (p *PaymentInfo) HasRequiredFields() bool {
	return p.Method != "" && p.AccountName != "" && p.BSB != "" && p.AccountNumber != ""
}
