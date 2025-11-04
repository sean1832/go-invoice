package invoice

import (
	"errors"
	"fmt"
	"go/format"
	"time"
)

type InvoiceStatus string

const (
	StatusDraft InvoiceStatus = "draft"
	StatusSend  InvoiceStatus = "send"
)

// Invoice represents the core invoice domain model
type Invoice struct {
	ID          string        `json:"id"`                     // invoice number/identifier
	Status      InvoiceStatus `json:"status"`                 // invoice status (draft, sent)
	Date        time.Time     `json:"date"`                   // invoice date
	Due         time.Time     `json:"due"`                    // payment due date
	Provider    Party         `json:"provider"`               // service provider
	Client      Party         `json:"client"`                 // client/customer
	Items       []ServiceItem `json:"items"`                  // list of services/products
	Pricing     Pricing       `json:"pricing"`                // pricing details
	Payment     PaymentInfo   `json:"payment"`                // payment information
	EmailTarget string        `json:"email_target,omitempty"` // (optional) email target for sending the invoice
	count       int           `json:"-"`                      // internal count, not serialized
}

// New creates a new invoice with the given provider, client, payment info, and tax rate
func New(provider, client Party, payment PaymentInfo, taxRate float32, count int) (*Invoice, error) {
	pricing, err := NewPricing(0, taxRate)
	if err != nil {
		return nil, err
	}
	return &Invoice{
		ID:       generateInvoiceID(count),
		Status:   StatusDraft,
		Date:     time.Now(),
		Due:      time.Now().AddDate(0, 0, 30), // due in 30 days
		Provider: provider,
		Client:   client,
		Items:    []ServiceItem{},
		Pricing:  *pricing,
		Payment:  payment,
		count:    count,
	}, nil
}

// SetEmailTarget sets the email address to send the invoice to
func (inv *Invoice) SetEmailTarget(email string) {
	inv.EmailTarget = email
}

func (inv *Invoice) SetID(id string) {
	inv.ID = id
}

func (inv *Invoice) HasRequiredFields() bool {
	return inv.ID != "" && inv.Status != "" && inv.Provider.HasRequiredFields() && inv.Client.HasRequiredFields() && inv.Payment.HasRequiredFields()
}

// AddItem adds a service item to the invoice and updates the pricing
func (inv *Invoice) AddItem(item ServiceItem) {
	inv.Items = append(inv.Items, item)
	inv.Pricing.Update(calculateSubtotal(inv.Items))
}

// generateInvoiceID creates a unique invoice ID based on the current date
func generateInvoiceID(count int) string {
	today := time.Now()
	formatted, _ := format.Source([]byte(today.Format("060102")))
	return fmt.Sprintf("INV-%s%03d", string(formatted), count)
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
	Date              time.Time `json:"date"`                         // date of service/product
	Description       string    `json:"description"`                  // description of service/product
	DescriptionDetail string    `json:"description_detail,omitempty"` // (optional) detailed description
	Quantity          float32   `json:"quantity"`                     // quantity provided
	UnitPrice         float32   `json:"unit_price"`                   // price per unit
	TotalPrice        float32   `json:"total_price"`                  // total price (Quantity * UnitPrice)
}

// NewServiceItem creates a new service item
func NewServiceItem(date time.Time, description string, quantity, unitPrice float32) ServiceItem {
	return ServiceItem{
		Date:        date,
		Description: description,
		Quantity:    quantity,
		UnitPrice:   unitPrice,
		TotalPrice:  quantity * unitPrice,
	}
}

// NewServiceItemWithDetail creates a new service item with detailed description
func NewServiceItemWithDetail(date time.Time, description, descriptionDetail string, quantity, unitPrice float32) ServiceItem {
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
