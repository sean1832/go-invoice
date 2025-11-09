import { apiClient } from '@/api/apiClient';

export interface VersionResponse {
	version: string;
}

/**
 * Get the application version from the backend
 */
export async function getVersion(fetchFn: typeof fetch): Promise<string> {
	const response = await apiClient<VersionResponse>(fetchFn, '/version', {
		method: 'GET'
	});
	return response.version;
}
