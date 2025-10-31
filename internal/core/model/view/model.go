package view

import (
	"encoding/json"
	"fmt"
	"go/format"
	"time"
)

// ViewModel represents the invoice data view model for rendering and JSON serialization
type ViewModel struct {
	Id          string         `json:"id"`                     // invoice number/identifier
	Date        time.Time      `json:"date"`                   // invoice date
	Due         time.Time      `json:"due"`                    // payment due date
	Provider    Party          `json:"provider"`               // service provider
	Client      Party          `json:"client"`                 // client/customer
	Items       []ServiceEntry `json:"items"`                  // list of services/products
	Pricing     Pricing        `json:"pricing"`                // pricing details
	Payment     PaymentInfo    `json:"payment"`                // payment information
	EmailTarget string         `json:"email_target,omitempty"` // (optional) email target for sending the invoice
}

func NewInvoiceData(provider, client Party, payment PaymentInfo, taxRate float32) *ViewModel {
	return &ViewModel{
		Id:       createInvoiceNumber(),
		Date:     time.Now(),
		Due:      time.Now().AddDate(0, 0, 30), // due in 30 days
		Provider: provider,
		Client:   client,
		Items:    []ServiceEntry{},
		Pricing:  NewPricing(0, taxRate),
		Payment:  payment,
	}
}

func (i *ViewModel) SetEmailTarget(email string) {
	i.EmailTarget = email
}

func createInvoiceNumber() string {
	today := time.Now()
	formatted, _ := format.Source([]byte(today.Format("20060102")))
	return fmt.Sprintf("INV-%s", string(formatted))
}

func NewInvoiceDataFromJSON(data []byte) (*ViewModel, error) {
	var invoice ViewModel
	err := json.Unmarshal(data, &invoice)
	if err != nil {
		return nil, err
	}
	return &invoice, nil
}

func (i *ViewModel) AddItem(item ServiceEntry) {
	i.Items = append(i.Items, item)
	i.Pricing.Update(CalculateSubtotal(i.Items))
}

func (i *ViewModel) ToJSON() ([]byte, error) {
	return json.Marshal(i)
}

// Party represents either the service provider or the client/customer
type Party struct {
	Name    string `json:"name"`    // service provider name
	Address string `json:"address"` // (optional) address
	Email   string `json:"email"`   // (optional) email address
	Phone   string `json:"phone"`   // (optional) phone number
	ABN     string `json:"abn"`     // (optional) Australian Business Number
}

// ServiceEntry represents a single line item in the invoice
type ServiceEntry struct {
	Date              time.Time `json:"date"`               // date of service/product
	Description       string    `json:"description"`        // description of service/product
	DescriptionDetail string    `json:"description_detail"` // (optional) detailed description
	Quantity          float32   `json:"quantity"`           // quantity provided
	UnitPrice         float32   `json:"unit_price"`         // price per unit
	TotalPrice        float32   `json:"total_price"`        // total price (Quantity * UnitPrice)
}

func NewServiceEntry(date time.Time, description string, quantity, unitPrice float32) ServiceEntry {
	return ServiceEntry{
		Date:        date,
		Description: description,
		Quantity:    quantity,
		UnitPrice:   unitPrice,
		TotalPrice:  quantity * unitPrice,
	}
}

func NewServiceEntryWithDetail(date time.Time, description, descriptionDetail string, quantity, unitPrice float32) ServiceEntry {
	return ServiceEntry{
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

func NewPricing(subtotal, taxRate float32) Pricing {
	taxAmount := subtotal * taxRate / 100
	total := subtotal + taxAmount
	return Pricing{
		Subtotal:  subtotal,
		TaxAmount: taxAmount,
		TaxRate:   taxRate,
		Total:     total,
	}
}
func (p *Pricing) Update(subtotal float32) {
	p.Subtotal = subtotal
	p.TaxAmount = subtotal * p.TaxRate / 100
	p.Total = p.Subtotal + p.TaxAmount
}

func CalculateSubtotal(items []ServiceEntry) float32 {
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
