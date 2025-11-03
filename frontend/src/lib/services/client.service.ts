import type { ClientData } from '@/types/invoice';
import { http } from '@/api/http';

/**
 * retrieves all clients using the provided fetch function.
 * @param KitFetch - `KitFetch` is a parameter that represents the fetch function provided by SvelteKit.
 * @returns A Promise that resolves to an array of ClientData objects.
 */
export async function getAllClients(KitFetch: typeof fetch): Promise<ClientData[]> {
	return http.get<ClientData[]>(KitFetch, '/clients');
}

/**
 * retrieves a single client by its ID using the provided fetch function.
 * @param KitFetch - `KitFetch` is a parameter that represents the fetch function provided by SvelteKit.
 * @param id - The `id` parameter is a string that represents the unique identifier of the client you want to retrieve.
 * @returns A Promise that resolves to the ClientData object with the specified ID.
 */
export async function getClient(KitFetch: typeof fetch, id: string): Promise<ClientData> {
	return http.get<ClientData>(KitFetch, `/clients/${id}`);
}

/**
 * creates a new client using the provided fetch function.
 * @param KitFetch - `KitFetch` is a parameter that represents the fetch function provided by SvelteKit.
 * @param data - The `data` parameter is a ClientData object that contains the details of the client to be created.
 * @returns A Promise that resolves to the created ClientData object.
 */
export async function createClient(KitFetch: typeof fetch, data: ClientData): Promise<ClientData> {
	return http.post<ClientData>(KitFetch, '/clients', data);
}

/**
 * updates an existing client using the provided fetch function.
 * @param KitFetch - `KitFetch` is a parameter that represents the fetch function provided by SvelteKit.
 * @param id - The `id` parameter is a string that represents the unique identifier of the client to be updated.
 * @param data - The `data` parameter is a ClientData object that contains the updated details of the client.
 * @returns A Promise that resolves to the updated ClientData object.
 */
export async function updateClient(
	KitFetch: typeof fetch,
	id: string,
	data: ClientData
): Promise<ClientData> {
	return http.put<ClientData>(KitFetch, `/clients/${id}`, data);
}

/**
 * deletes an existing client using the provided fetch function.
 * @param KitFetch - `KitFetch` is a parameter that represents the fetch function provided by SvelteKit.
 * @param id - The `id` parameter is a string that represents the unique identifier of the client to be deleted.
 * @returns A Promise that resolves to the response from the server.
 */
export async function deleteClient(KitFetch: typeof fetch, id: string) {
	return http.delete(KitFetch, `/clients/${id}`);
}
