/**
 * Validators - Pure validation functions
 * All functions return validation results with error messages
 */

import type { ServiceItem, Party, Invoice, EmailConfig } from '@/types/invoice';

export interface ValidationResult {
	isValid: boolean;
	errors: string[];
}

/**
 * Validate email format
 * @param email - Email address to validate
 * @returns True if valid email format
 */
export function isValidEmail(email: string): boolean {
	const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
	return emailRegex.test(email);
}

/**
 * Validate ABN format (11 digits)
 * @param abn - ABN to validate
 * @returns True if valid ABN format
 */
export function isValidABN(abn: string): boolean {
	const digits = abn.replace(/\D/g, '');
	return digits.length === 11;
}

/**
 * Validate phone number (basic check for Australian numbers)
 * @param phone - Phone number to validate
 * @returns True if valid phone format
 */
export function isValidPhone(phone: string): boolean {
	const digits = phone.replace(/\D/g, '');
	// Australian numbers are typically 10 digits (including area code)
	return digits.length >= 10;
}

/**
 * Validate a party (provider or client) object
 * @param party - Party object to validate
 * @param type - Type of party for error messages
 * @returns Validation result with errors
 */
export function validateParty(party: Party, type: 'provider' | 'client'): ValidationResult {
	const errors: string[] = [];

	if (!party.name || party.name.trim() === '') {
		errors.push(`${type} name is required`);
	}

	if (party.email && !isValidEmail(party.email)) {
		errors.push(`${type} email is invalid`);
	}

	if (party.abn && !isValidABN(party.abn)) {
		errors.push(`${type} ABN must be 11 digits`);
	}

	if (party.phone && !isValidPhone(party.phone)) {
		errors.push(`${type} phone number is invalid`);
	}

	return {
		isValid: errors.length === 0,
		errors
	};
}

/**
 * Validate a service line item
 * @param item - Service item to validate
 * @returns Validation result with errors
 */
export function validateLineItem(item: ServiceItem): ValidationResult {
	const errors: string[] = [];

	if (!item.description || item.description.trim() === '') {
		errors.push(`Description is required`);
	}

	if (item.quantity <= 0) {
		errors.push(`Quantity must be greater than 0`);
	}

	if (item.unit_price < 0) {
		errors.push(`Unit price cannot be negative`);
	}

	return {
		isValid: errors.length === 0,
		errors
	};
}

/**
 * Validate an entire invoice
 * @param invoice - Partial invoice object to validate
 * @returns Validation result with all errors
 */
export function validateInvoice(invoice: Partial<Invoice>): ValidationResult {
	const errors: string[] = [];
	// TODO: update validations requirements
	// Validate provider
	if (invoice.provider) {
		const providerResult = validateParty(invoice.provider, 'provider');
		errors.push(...providerResult.errors);
	} else {
		errors.push('Provider is required');
	}

	// Validate client
	if (invoice.client) {
		const clientResult = validateParty(invoice.client, 'client');
		errors.push(...clientResult.errors);
	} else {
		errors.push('Client is required');
	}

	// Validate items
	if (!invoice.items || invoice.items.length === 0) {
		errors.push('At least one service item is required');
	} else {
		invoice.items.forEach((item) => {
			const itemResult = validateLineItem(item);
			errors.push(...itemResult.errors);
		});
	}

	// Validate dates
	if (!invoice.date) {
		errors.push('Issue date is required');
	}

	if (!invoice.due) {
		errors.push('Due date is required');
	}

	return {
		isValid: errors.length === 0,
		errors
	};
}

/**
 * Check if a string is empty or only whitespace
 * @param value - String to check
 * @returns True if empty or whitespace
 */
export function isEmpty(value: string | undefined | null): boolean {
	return !value || value.trim() === '';
}

export function isListEmpty<T>(list: T[] | undefined | null): boolean {
	return !list || list.length === 0;
}

/**
 * Validate required field
 * @param value - Value to check
 * @param fieldName - Name of field for error message
 * @returns Error message or null if valid
 */
export function validateRequired(
	value: string | undefined | null,
	fieldName: string
): string | null {
	return isEmpty(value) ? `${fieldName} is required` : null;
}

export function validateEmailConfig(emailConfig: EmailConfig): ValidationResult {
	const errors: string[] = [];

	if (isListEmpty(emailConfig.to)) {
		errors.push('Email recipient (To) is required');
	} else {
		// Validate each email address in the array
		const invalidEmails = emailConfig.to!.filter((email) => !isValidEmail(email));
		if (invalidEmails.length > 0) {
			errors.push(`Invalid email address(es): ${invalidEmails.join(', ')}`);
		}
	}

	if (isEmpty(emailConfig.subject)) {
		errors.push('Email subject is required');
	}
	if (isEmpty(emailConfig.body)) {
		errors.push('Email body is required');
	}
	return {
		isValid: errors.length === 0,
		errors
	};
}
