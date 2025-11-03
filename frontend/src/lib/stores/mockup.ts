/**
 * Centralized mock data for development
 * When switching to API, this file can be safely deleted
 */

import type { Invoice, ClientData, ProviderData } from '@/types/invoice';

export const mockProviders: ProviderData[] = [
	{
		id: 'zeke_zhang',
		name: 'Zeke Zhang',
		email: 'zeke@example.com',
		address: '123 Provider St',
		phone: '+61 2 1234 5678',
		abn: '12 345 678 901',
		payment_info: {
			method: 'Bank Transfer',
			account_name: 'Zeke Zhang',
			bsb: '123-456',
			account_number: '12345678'
		}
	},
	{
		id: 'lan_zhang',
		name: 'Lan Zhang',
		email: 'lan@example.com',
		address: '456 Business Ave',
		phone: '+61 3 9876 5432',
		abn: '98 765 432 109',
		payment_info: {
			method: 'Bank Transfer',
			account_name: 'Lan Zhang',
			bsb: '654-321',
			account_number: '87654321'
		}
	}
];

export const mockClients: ClientData[] = [
	{
		id: 'dingyu_xu',
		name: 'Dingyu Xu',
		email: 'dingyu@example.com',
		address: '789 Client Rd',
		phone: '+61 7 5555 1234',
		abn: '11 222 333 444',
		tax_rate: 10,
		target_email: 'dingyu@example.com'
	},
	{
		id: 'acme_corp',
		name: 'Acme Corporation',
		email: 'contact@acme.com',
		address: '123 Business St',
		phone: '+61 2 9876 5432',
		abn: '98 765 432 109',
		tax_rate: 10
	},
	{
		id: 'tech_startup',
		name: 'Tech Startup Inc',
		email: 'info@techstartup.com',
		address: '456 Innovation Way',
		phone: '+61 3 1234 5678',
		abn: '12 345 678 901',
		tax_rate: 10
	},
	{
		id: 'global_solutions',
		name: 'Global Solutions Ltd',
		email: 'contact@globalsolutions.com',
		address: '789 Enterprise Rd',
		phone: '+61 7 8765 4321',
		abn: '23 456 789 012',
		tax_rate: 10
	}
];

export const mockInvoices: Invoice[] = [
	{
		id: 'INV-251103001',
		status: 'draft',
		date: '2025-11-01',
		due: '2025-11-30',
		client: mockClients[1], // Acme Corporation
		provider: mockProviders[0], // Zeke Zhang
		items: [
			{
				date: '2025-11-01',
				description: 'Web Development Services',
				description_detail: 'Frontend development and UI/UX improvements',
				quantity: 40,
				unit_price: 125,
				total_price: 5000
			}
		],
		pricing: {
			subtotal: 5000,
			tax: 500,
			tax_rate: 10,
			total: 5500
		},
		payment: mockProviders[0].payment_info,
		email_target: 'contact@acme.com'
	},
	{
		id: 'INV-251103002',
		status: 'send',
		date: '2025-10-28',
		due: '2025-11-15',
		client: mockClients[2], // Tech Startup Inc
		provider: mockProviders[0], // Zeke Zhang
		items: [
			{
				date: '2025-10-28',
				description: 'Consulting Services',
				description_detail: 'Technical architecture consultation',
				quantity: 20,
				unit_price: 160,
				total_price: 3200
			}
		],
		pricing: {
			subtotal: 3200,
			tax: 320,
			tax_rate: 10,
			total: 3520
		},
		payment: mockProviders[0].payment_info,
		email_target: 'billing@techstartup.com'
	},
	{
		id: 'INV-251103003',
		status: 'draft',
		date: '2025-10-15',
		due: '2025-11-01',
		client: mockClients[3], // Global Solutions Ltd
		provider: mockProviders[1], // Lan Zhang
		items: [
			{
				date: '2025-10-15',
				description: 'API Development',
				description_detail: 'RESTful API design and implementation',
				quantity: 30,
				unit_price: 150,
				total_price: 4500
			}
		],
		pricing: {
			subtotal: 4500,
			tax: 450,
			tax_rate: 10,
			total: 4950
		},
		payment: mockProviders[1].payment_info,
		email_target: 'accounts@global.com'
	}
];
