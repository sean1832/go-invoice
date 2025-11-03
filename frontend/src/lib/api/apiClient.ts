const BASE_URL = '/api/v1';

export async function apiClient<T>(
	KitFetch: typeof fetch,
	path: string,
	options?: RequestInit
): Promise<T> {
	const url = `${BASE_URL}${path}`;

	const response = await KitFetch(url, options);

	// CRITICAL: HTTP Error (400, 500)
	if (!response.ok) {
		const errBody = await response.text();
		throw new Error(`API Error [${response.status} ${response.statusText}] for ${url}: ${errBody}`);
	}

	// 204: No Content
	// .json will fail so return an empty object
	if (response.status == 204) {
		return {} as T;
	}

	return response.json() as Promise<T>;
}
