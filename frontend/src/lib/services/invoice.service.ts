import { http } from '@/api/http';
import type { EmailConfig, Invoice } from '@/types/invoice';

/**
 * retrieves all invoices using the provided fetch function.
 * @param KitFetch - `KitFetch` is a parameter that represents the fetch function provided by SvelteKit.
 * @returns A Promise that resolves to an array of Invoice objects.
 */
export async function getAllInvoices(KitFetch: typeof fetch): Promise<Invoice[]> {
	return http.get<Invoice[]>(KitFetch, '/invoices');
}

/**
 * retrieves a single invoice by its ID using the provided fetch function.
 * @param KitFetch - `KitFetch` is a parameter that represents the fetch function provided by SvelteKit.
 * @param id - The `id` parameter is a string that represents the unique identifier of the invoice you want to retrieve.
 * @returns A Promise that resolves to the Invoice object with the specified ID.
 */
export async function getInvoice(KitFetch: typeof fetch, id: string): Promise<Invoice> {
	return http.get<Invoice>(KitFetch, `/invoices/${id}`);
}

/**
 * creates a new invoice using the provided fetch function.
 * @param KitFetch - `KitFetch` is a parameter that represents the fetch function provided by SvelteKit.
 * @param data - The `data` parameter is an Invoice object that contains the details of the invoice to be created.
 * @returns A Promise that resolves to the created Invoice object.
 */
export async function createInvoice(KitFetch: typeof fetch, data: Invoice): Promise<Invoice> {
	return http.post<Invoice>(KitFetch, '/invoices', data);
}

/**
 * updates an existing invoice using the provided fetch function.
 * @param KitFetch - `KitFetch` is a parameter that represents the fetch function provided by SvelteKit.
 * @param id - The `id` parameter is a string that represents the unique identifier of the invoice to be updated.
 * @param data - The `data` parameter is an Invoice object that contains the updated details of the invoice.
 * @returns A Promise that resolves to the updated Invoice object.
 */
export async function updateInvoice(
	KitFetch: typeof fetch,
	id: string,
	data: Invoice
): Promise<Invoice> {
	return http.put<Invoice>(KitFetch, `/invoices/${id}`, data);
}

/**
 * deletes an existing invoice using the provided fetch function.
 * @param KitFetch - `KitFetch` is a parameter that represents the fetch function provided by SvelteKit.
 * @param id - The `id` parameter is a string that represents the unique identifier of the invoice to be deleted.
 * @returns A Promise that resolves to the response from the server.
 */
export async function deleteInvoice(KitFetch: typeof fetch, id: string) {
	return http.delete(KitFetch, `/invoices/${id}`);
}

/**
 * downloads the PDF of an invoice using the provided fetch function.
 * @param KitFetch - `KitFetch` is a parameter that represents the fetch function provided by SvelteKit.
 * @param id - The `id` parameter is a string that represents the unique identifier of the invoice whose PDF you want to download.
 * @returns A Promise that resolves to the response containing the PDF data.
 */
export async function downloadPdf(KitFetch: typeof fetch, id: string): Promise<Blob> {
	return http.get<Blob>(KitFetch, `/invoices/${id}/pdf`, {
		responseType: 'blob'
	});
}

/**
 * Sends an email for a specific invoice using the provided fetch function. (this will take some time to process, best to add loading indicator)
 * @param KitFetch - `KitFetch` is a parameter that represents the fetch function provided by SvelteKit.
 * @param id - The `id` parameter is a string that represents the unique identifier of the invoice to send.
 * @param emailConfig - The `emailConfig` parameter is an EmailConfig object that contains the email details.
 * @returns A Promise that resolves to the response from the server.
 */
export async function sendInvoiceEmail(
	KitFetch: typeof fetch,
	id: string,
	emailConfig: EmailConfig
) {
	return http.post<void>(KitFetch, `/invoices/${id}/email`, emailConfig, {
		timeout: 10000 // set timeout to 10 seconds
	});
}
