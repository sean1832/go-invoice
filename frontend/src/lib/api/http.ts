import { apiClient, type ApiClientOptions } from './apiClient';

type KitFetch = typeof window.fetch;

export const http = {
	/**
	 * Perform a `GET` request.
	 */
	get: <T>(fetch: KitFetch, path: string, options: ApiClientOptions = {}) => {
		return apiClient<T>(fetch, path, options); // no need to do anything, GET request by default
	},

	/**
	 * Perform a `POST` request
	 */
	post: <T>(fetch: KitFetch, path: string, data: unknown, options: ApiClientOptions = {}) => {
		return apiClient<T>(fetch, path, {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json'
			},
			body: JSON.stringify(data),
			...options
		});
	},

	/**
	 * Perform a `PUT` request
	 */
	put: <T>(fetch: KitFetch, path: string, data: unknown, options: ApiClientOptions = {}) => {
		return apiClient<T>(fetch, path, {
			method: 'PUT',
			headers: {
				'Content-Type': 'application/json'
			},
			body: JSON.stringify(data),
			...options
		});
	},

	/**
	 * Perform a `DELETE` request
	 */
	delete: <T = Record<string, never>>(
		fetch: KitFetch,
		path: string,
		options: ApiClientOptions = {}
	) => {
		return apiClient<T>(fetch, path, {
			method: 'DELETE',
			...options
		});
	}
};
