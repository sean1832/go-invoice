import { http } from '@/api/http';
import type { ProviderData } from '@/types/invoice';

/**
 *  retrieves all provider data using the provided `KitFetch` function.
 * @param KitFetch - `KitFetch` is SvelteKit's built-in fetch function that is used to make HTTP requests.
 * @returns A Promise that resolves to an array of `ProviderData` objects.
 */
export async function getAllProviders(KitFetch: typeof fetch): Promise<ProviderData[]> {
	return http.get<ProviderData[]>(KitFetch, '/providers');
}

/**
 * retrieves a single provider by its ID using the provided fetch function.
 * @param KitFetch - `KitFetch` is a parameter that represents the fetch function provided by SvelteKit.
 * @param id - The `id` parameter is a string that represents the unique identifier of the provider you want to retrieve.
 * @returns A Promise that resolves to the ProviderData object with the specified ID.
 */
export async function getProvider(KitFetch: typeof fetch, id: string): Promise<ProviderData> {
	return http.get<ProviderData>(KitFetch, `/providers/${id}`);
}

/**
 * creates a new provider using the provided fetch function.
 * @param KitFetch - `KitFetch` is a parameter that represents the fetch function provided by SvelteKit.
 * @param data - The `data` parameter is a ProviderData object that contains the details of the provider to be created.
 * @returns A Promise that resolves to the created ProviderData object.
 */
export async function createProvider(
	KitFetch: typeof fetch,
	data: ProviderData
): Promise<ProviderData> {
	return http.post<ProviderData>(KitFetch, '/providers', data);
}

/**
 * updates an existing provider using the provided fetch function.
 * @param KitFetch - `KitFetch` is a parameter that represents the fetch function provided by SvelteKit.
 * @param id - The `id` parameter is a string that represents the unique identifier of the provider to be updated.
 * @param data - The `data` parameter is a ProviderData object that contains the updated details of the provider.
 * @returns A Promise that resolves to the updated ProviderData object.
 */
export async function updateProvider(
	KitFetch: typeof fetch,
	id: string,
	data: ProviderData
): Promise<ProviderData> {
	return http.put<ProviderData>(KitFetch, `/providers/${id}`, data);
}

/**
 * deletes a provider by its ID using the provided fetch function.
 * @param KitFetch - `KitFetch` is a parameter that represents the fetch function provided by SvelteKit.
 * @param id - The `id` parameter is a string that represents the unique identifier of the provider to be deleted.
 * @returns A Promise that resolves to the response from the server.
 */
export async function deleteProvider(KitFetch: typeof fetch, id: string) {
	return http.delete(KitFetch, `/providers/${id}`);
}
