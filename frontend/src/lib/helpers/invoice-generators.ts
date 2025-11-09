/**
 * Invoice Generators - Functions to generate IDs and create default objects
 */

import type { ServiceItem, Party } from '@/types/invoice';
import { getTodayISOString } from './date-helpers';


/**
 * Create an empty service line item with default values
 * @returns Empty ServiceItem object
 */
export function createEmptyLineItem(): ServiceItem {
	return {
		date: getTodayISOString(),
		description: '',
		description_detail: '',
		quantity: 1,
		unit_price: 0,
		total_price: 0
	};
}

/**
 * Create an empty party (provider or client) object
 * @returns Empty Party object
 */
export function createEmptyParty(): Party {
	return {
		id: '',
		name: '',
		address: '',
		email: '',
		phone: '',
		abn: ''
	};
}

/**
 * Generate a unique party ID
 * Format: party-{timestamp}-{random}
 * @param type - Type of party ('provider' or 'client')
 * @returns Generated party ID
 */
export function generatePartyId(type: 'provider' | 'client'): string {
	const timestamp = Date.now();
	const random = Math.floor(Math.random() * 1000);
	return `${type}-${timestamp}-${random}`;
}

/**
 * Sanitize a string to create a valid filename-safe ID
 * @param name - Name to sanitize
 * @returns Sanitized string (lowercase, no spaces, alphanumeric + hyphen)
 */
export function sanitizeForId(name: string): string {
	return name
		.toLowerCase()
		.replace(/\s+/g, '_')
		.replace(/[^a-z0-9_-]/g, '')
		.replace(/_+/g, '_');
}
