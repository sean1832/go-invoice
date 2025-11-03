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
	// Remove all non-digit characters
	const digits = phone.replace(/\D/g, '');

	// Format as: 0412 345 678 or 02 1234 5678
	if (digits.length === 10) {
		if (digits.startsWith('04')) {
			return `${digits.slice(0, 4)} ${digits.slice(4, 7)} ${digits.slice(7)}`;
		} else {
			return `${digits.slice(0, 2)} ${digits.slice(2, 6)} ${digits.slice(6)}`;
		}
	}

	return phone; // Return as-is if format doesn't match
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
 * Format a percentage
 * @param rate - Percentage rate as number (e.g., 10 for 10%)
 * @param decimals - Number of decimal places
 * @returns Formatted percentage string (e.g., "10%")
 */
export function formatPercentage(rate: number, decimals: number = 0): string {
	return `${rate.toFixed(decimals)}%`;
}
