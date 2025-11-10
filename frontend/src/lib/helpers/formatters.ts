/**
 * Formatters - Pure formatting functions for display
 * All functions are pure (no side effects) and easily testable
 */

/**
 * Format a number as Australian currency
 * @param amount - The amount to format
 * @param currency - Currency code (default: AUD)
 * @returns Formatted currency string (e.g., "$1,234.56")
 */
export function formatCurrency(amount: number, currency = 'AUD'): string {
	return new Intl.NumberFormat('en-AU', {
		style: 'currency',
		currency
	}).format(amount);
}

/**
 * Format a date string to long format
 * @param dateString - ISO date string (YYYY-MM-DD)
 * @returns Formatted date (e.g., "November 3, 2025")
 */
export function formatDate(dateString: string): string {
	return new Date(dateString).toLocaleDateString('en-US', {
		year: 'numeric',
		month: 'long',
		day: 'numeric'
	});
}

/**
 * Format a date string to short format
 * @param dateString - ISO date string (YYYY-MM-DD)
 * @returns Formatted date (e.g., "Nov 3, 2025")
 */
export function formatDateShort(dateString: string): string {
	return new Date(dateString).toLocaleDateString('en-US', {
		month: 'short',
		day: 'numeric',
		year: 'numeric'
	});
}

/**
 * Format a date string to numeric format
 * @param dateString - ISO date string (YYYY-MM-DD)
 * @returns Formatted date (e.g., "11/03/2025")
 */
export function formatDateNumeric(dateString: string): string {
	return new Date(dateString).toLocaleDateString('en-US', {
		month: '2-digit',
		day: '2-digit',
		year: 'numeric'
	});
}

/**
 * Format a phone number (basic Australian format)
 * @param phone - Phone number string
 * @returns Formatted phone number
 */
export function formatPhone(phone: string): string {
	// clean the phone number, but keep the leading '+'
	const cleanPhone = phone.replace(/[^\d+]/g, '');

	// International Format (+61)
	if (cleanPhone.startsWith('+61')) {
		const nationalPart = cleanPhone.substring(3); // Get the number without +61

		// Mobile: +61 4xx xxx xxx (9 digits starting with 4)
		if (nationalPart.startsWith('4') && nationalPart.length === 9) {
			return `+61 ${nationalPart.slice(0, 3)} ${nationalPart.slice(3, 6)} ${nationalPart.slice(6)}`;
		}

		// Landline: +61 x xxxx xxxx (9 digits starting with 2, 3, 7, or 8)
		if (['2', '3', '7', '8'].includes(nationalPart.charAt(0)) && nationalPart.length === 9) {
			return `+61 ${nationalPart.slice(0, 1)} ${nationalPart.slice(1, 5)} ${nationalPart.slice(5)}`;
		}
	}

	// Domestic Format (0x)
	// strip all non-digits
	const digits = phone.replace(/\D/g, '');

	if (digits.length === 10) {
		// Mobile: 04xx xxx xxx
		if (digits.startsWith('04')) {
			return `${digits.slice(0, 4)} ${digits.slice(4, 7)} ${digits.slice(7)}`;
		}

		// Landline: 0x xxxx xxxx
		if (
			digits.startsWith('02') ||
			digits.startsWith('03') ||
			digits.startsWith('07') ||
			digits.startsWith('08')
		) {
			return `${digits.slice(0, 2)} ${digits.slice(2, 6)} ${digits.slice(6)}`;
		}
	}

	// Fallback: Return the original string if no format matches
	return phone;
}

/**
 * Format an ABN (Australian Business Number)
 * @param abn - ABN string
 * @returns Formatted ABN (e.g., "12 345 678 901")
 */
export function formatABN(abn: string): string {
	const digits = abn.replace(/\D/g, '');
	if (digits.length === 11) {
		return `${digits.slice(0, 2)} ${digits.slice(2, 5)} ${digits.slice(5, 8)} ${digits.slice(8)}`;
	}
	return abn;
}

/**
 * Format a BSB (Bank State Branch) number
 * @param bsb - BSB string
 * @returns Formatted BSB (e.g., "123-456")
 */
export function formatBSB(bsb: string): string {
	const digits = bsb.replace(/\D/g, '');
	if (digits.length === 6) {
		return `${digits.slice(0, 3)}-${digits.slice(3)}`;
	}
	return bsb;
}

/**
 * Format a percentage
 * @param rate - Percentage rate as number (e.g., 10 for 10%)
 * @param decimals - Number of decimal places
 * @returns Formatted percentage string (e.g., "10%")
 */
export function formatPercentage(rate: number, decimals: number = 0): string {
	return `${rate.toFixed(decimals)}%`;
}

/**
 * Format an email template with the given variables
 * @param template - The email template string
 * @param variables - The variables to replace in the template
 * @returns The formatted email string
 */
export function formatEmailTemplate(template: string, variables: Record<string, string>): string {
	let formatted = template;
	for (const [key, value] of Object.entries(variables)) {
		const placeholder = `{{${key}}}`;
		formatted = formatted.replaceAll(placeholder, value);
	}
	return formatted;
}
