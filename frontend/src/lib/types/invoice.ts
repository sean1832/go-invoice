// Invoice status types - matching Go backend
export type InvoiceStatus = 'draft' | 'send';

export interface ClientData extends Party {
	tax_rate: number;
	email_target?: string;
	email_template_id: string;
}

export interface ProviderData extends Party {
	payment_info: PaymentInfo;
}

// Party represents either the service provider or the client/customer
export interface Party {
	id: string;
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
	description_detail?: string;
	quantity: number;
	unit_price: number;
	total_price: number; // quantity * unitPrice
}

// Pricing holds the pricing breakdown of the invoice
export interface Pricing {
	subtotal: number;
	tax: number; // renamed from taxAmount in display
	tax_rate: number;
	total: number;
}

// Payment information
export interface PaymentInfo {
	method: string; // e.g., "Bank Transfer"
	account_name: string;
	bsb: string; // Bank State Branch number
	account_number: string;
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
	email_target?: string; // optional email target for sending
	email_template_id?: string; // optional email template ID
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
	tax_rate: number;
}

// Filter options for the shelf
export interface InvoiceFilters {
	status?: InvoiceStatus | 'all';
	search_query?: string;
	sort_by?: 'date' | 'amount' | 'client';
	sort_order?: 'asc' | 'desc';
}

export interface EmailContent {
	subject: string;
	body: string;
}

export interface EmailTemplate extends EmailContent {
	id: string;
	name: string;
}

export interface EmailConfig extends EmailContent {
	to: string[];
}
