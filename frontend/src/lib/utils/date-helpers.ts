/**
 * Date Helpers - Pure date utility functions
 * All dates are in ISO format (YYYY-MM-DD)
 */

import { getLocalTimeZone, parseDate, today, type DateValue } from '@internationalized/date';

/**
 * Get today's date as ISO string
 * @returns Today's date in YYYY-MM-DD format
 */
export function getTodayISOString(): string {
	return today(getLocalTimeZone()).toString();
}

/**
 * Get a date N days from today
 * @param daysFromNow - Number of days to add (positive) or subtract (negative)
 * @returns Date in YYYY-MM-DD format
 */
export function getDateFromToday(daysFromNow: number): string {
	const startDate = today(getLocalTimeZone());
	const targetDate = startDate.add({ days: daysFromNow });
	return targetDate.toString();
}

/**
 * Get default due date (30 days from today)
 * @returns Due date in YYYY-MM-DD format
 */
export function getDefaultDueDate(): string {
	return getDateFromToday(30);
}

/**
 * Get default invoice issue date (today)
 * @returns Issue date in YYYY-MM-DD format
 */
export function getDefaultIssueDate(): string {
	return getTodayISOString();
}

/**
 * Safely parse an ISO date string to DateValue
 * @param dateString - ISO date string (YYYY-MM-DD) or (YYYY-MM-DDTHH:mm:ssZ)
 * @returns DateValue or undefined if invalid
 */
export function safeParseDate(dateString: string | undefined | null): DateValue | undefined {
	if (!dateString) return undefined;

	// Extract only date part (YYYY-MM-DD) before parsing
	const dateOnlyString = dateString.split('T')[0];
	try {
		return parseDate(dateOnlyString);
	} catch (error) {
		console.error('Failed to parse date:', dateString, error);
		return undefined;
	}
}

/**
 * Convert DateValue to ISO string
 * @param dateValue - DateValue object
 * @returns ISO date string (YYYY-MM-DD)
 */
export function dateValueToISOString(dateValue: DateValue): string {
	return dateValue.toString();
}

/**
 * Check if a date string is valid ISO format
 * @param dateString - Date string to validate
 * @returns True if valid ISO date
 */
export function isValidISODate(dateString: string): boolean {
	const isoDateRegex = /^\d{4}-\d{2}-\d{2}$/;
	if (!isoDateRegex.test(dateString)) return false;

	const date = new Date(dateString);
	return date instanceof Date && !isNaN(date.getTime());
}

/**
 * Compare two ISO date strings
 * @param date1 - First date string
 * @param date2 - Second date string
 * @returns Negative if date1 < date2, 0 if equal, positive if date1 > date2
 */
export function compareDates(date1: string, date2: string): number {
	if (date1 < date2) return -1;
	if (date1 > date2) return 1;
	return 0;
}

/**
 * Check if a date is in the past
 * @param dateString - ISO date string
 * @returns True if date is before today
 */
export function isDateInPast(dateString: string): boolean {
	return compareDates(dateString, getTodayISOString()) < 0;
}

/**
 * Check if a date is in the future
 * @param dateString - ISO date string
 * @returns True if date is after today
 */
export function isDateInFuture(dateString: string): boolean {
	return compareDates(dateString, getTodayISOString()) > 0;
}
