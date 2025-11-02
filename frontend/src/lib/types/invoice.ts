// Invoice status types - matching Go backend
export type InvoiceStatus = 'draft' | 'send';

// Party represents either the service provider or the client/customer
export interface Party {
	name: string;
	address?: string;
	email?: string;
	phone?: string;
	abn?: string; // Australian Business Number
}

// Service item represents a single line item in the invoice
export interface ServiceItem {
	date: string; // ISO date string
	description: string;
	descriptionDetail?: string;
	quantity: number;
	unitPrice: number;
	totalPrice: number; // quantity * unitPrice
}

// Pricing holds the pricing breakdown of the invoice
export interface Pricing {
	subtotal: number;
	tax: number; // renamed from taxAmount in display
	taxRate: number;
	total: number;
}

// Payment information
export interface PaymentInfo {
	method: string; // e.g., "Bank Transfer"
	accountName: string;
	bsb: string; // Bank State Branch number
	accountNumber: string;
}

// Main invoice type - matching Go backend
export interface Invoice {
	id: string; // invoice number/identifier
	status: InvoiceStatus;
	date: string; // ISO date string - invoice date
	due: string; // ISO date string - payment due date
	provider: Party;
	client: Party;
	items: ServiceItem[];
	pricing: Pricing;
	payment: PaymentInfo;
	emailTarget?: string; // optional email target for sending
}

// Form data for creating/editing invoices
export interface InvoiceFormData {
	id: string;
	provider: Party;
	client: Party;
	date: string;
	due: string;
	items: Omit<ServiceItem, 'totalPrice'>[];
	payment: PaymentInfo;
	taxRate: number;
}

// Filter options for the shelf
export interface InvoiceFilters {
	status?: InvoiceStatus | 'all';
	searchQuery?: string;
	sortBy?: 'date' | 'amount' | 'client';
	sortOrder?: 'asc' | 'desc';
}
