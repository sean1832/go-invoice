import { http } from '@/api/http';
import { authStore } from '@/stores/auth';

interface SessionResponse {
	authenticated: boolean;
	email?: string;
	method: 'oauth2' | 'plain' | 'none';
}

/**
 * Check current authentication session
 */
export async function checkSession(fetchFn: typeof fetch): Promise<SessionResponse> {
	try {
		const response = await http.get<SessionResponse>(fetchFn, '/mailer/session');

		if (response.authenticated && response.email) {
			authStore.setAuthenticated(response.email);
		} else {
			authStore.setUnauthenticated();
		}

		authStore.setAuthMethod(response.method);
		return response;
	} catch (error) {
		authStore.setUnauthenticated();
		throw error;
	}
}

/**
 * Open OAuth login popup and wait for completion
 */
export function loginWithGoogle(): Promise<void> {
	return new Promise((resolve, reject) => {
		const width = 500;
		const height = 600;
		const left = window.screenX + (window.outerWidth - width) / 2;
		const top = window.screenY + (window.outerHeight - height) / 2;

		const popup = window.open(
			'/api/v1/mailer/auth/google',
			'oauth-login',
			`width=${width},height=${height},left=${left},top=${top},toolbar=0,menubar=0,location=0`
		);

		if (!popup) {
			reject(new Error('Popup blocked. Please allow popups for this site.'));
			return;
		}

		// Listen for message from popup
		const messageHandler = (event: MessageEvent) => {
			if (event.origin !== window.location.origin) return;
			if (event.data.type === 'AUTH_SUCCESS') {
				window.removeEventListener('message', messageHandler);
				popup?.close();

				// Refresh session to update store
				checkSession(fetch)
					.then(() => resolve())
					.catch(reject);
			}
		};

		window.addEventListener('message', messageHandler);

		// Check if popup was closed without completing auth
		const checkClosed = setInterval(() => {
			if (popup.closed) {
				clearInterval(checkClosed);
				window.removeEventListener('message', messageHandler);
				reject(new Error('Authentication cancelled'));
			}
		}, 500);
	});
}

/**
 * Logout and clear session
 */
export async function logout(fetchFn: typeof fetch): Promise<void> {
	try {
		await http.post(fetchFn, '/mailer/logout', {});
		authStore.reset();
	} catch (error) {
		console.error('Logout failed:', error);
		throw error;
	}
}
