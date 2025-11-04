/**
 * Invoice Calculations - Pure calculation functions
 * All functions are pure (no side effects) and easily testable
 */

import type { ServiceItem, Pricing } from '@/types/invoice';

/**
 * Calculate the total price for a single line item
 * @param quantity - Number of units
 * @param unitPrice - Price per unit
 * @returns Total price (quantity × unitPrice)
 */
export function calculateLineItemTotal(quantity: number, unitPrice: number): number {
	return quantity * unitPrice;
}

/**
 * Calculate subtotal from an array of service items
 * @param items - Array of service items
 * @returns Sum of all item total prices
 */
export function calculateSubtotal(items: ServiceItem[]): number {
	return items.reduce((sum, item) => sum + item.total_price, 0);
}

/**
 * Calculate tax amount from subtotal and tax rate
 * @param subtotal - Subtotal amount
 * @param taxRate - Tax rate as percentage (e.g., 10 for 10%)
 * @returns Tax amount
 */
export function calculateTax(subtotal: number, taxRate: number): number {
	return subtotal * (taxRate / 100);
}

/**
 * Calculate total (subtotal + tax)
 * @param subtotal - Subtotal amount
 * @param tax - Tax amount
 * @returns Total amount
 */
export function calculateTotal(subtotal: number, tax: number): number {
	return subtotal + tax;
}

/**
 * Calculate complete pricing breakdown from items and tax rate
 * @param items - Array of service items
 * @param taxRate - Tax rate as percentage (e.g., 10 for 10%)
 * @returns Complete pricing object with subtotal, tax, taxRate, and total
 */
export function calculatePricing(items: ServiceItem[], taxRate: number): Pricing {
	const subtotal = calculateSubtotal(items);
	const tax = calculateTax(subtotal, taxRate);
	const total = calculateTotal(subtotal, tax);

	return {
		subtotal,
		tax,
		tax_rate: taxRate,
		total
	};
}

/**
 * Round currency to 2 decimal places
 * @param amount - Amount to round
 * @returns Rounded amount
 */
export function roundCurrency(amount: number): number {
	return Math.round(amount * 100) / 100;
}
