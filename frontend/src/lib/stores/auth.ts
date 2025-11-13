import { writable, derived } from 'svelte/store';

export interface AuthState {
	isAuthenticated: boolean;
	userEmail: string | null;
	authMethod: 'oauth2' | 'plain' | 'none' | null;
	loading: boolean;
}

function createAuthStore() {
	const { subscribe, update } = writable<AuthState>({
		isAuthenticated: false,
		userEmail: null,
		authMethod: null,
		loading: true
	});

	return {
		subscribe,
		setAuthenticated: (email: string) =>
			update((state) => ({
				...state,
				isAuthenticated: true,
				userEmail: email,
				loading: false
			})),
		setUnauthenticated: () =>
			update((state) => ({
				...state,
				isAuthenticated: false,
				userEmail: null,
				loading: false
			})),
		setAuthMethod: (method: 'oauth2' | 'plain' | 'none') =>
			update((state) => ({
				...state,
				authMethod: method
			})),
		setLoading: (loading: boolean) =>
			update((state) => ({
				...state,
				loading
			})),
		reset: () =>
			update((state) => ({
				isAuthenticated: false,
				userEmail: null,
				authMethod: state.authMethod, // Preserve authMethod so UI doesn't disappear
				loading: false
			}))
	};
}

export const authStore = createAuthStore();

// Derived stores for convenience
export const isAuthenticated = derived(authStore, ($auth) => $auth.isAuthenticated);
export const currentUserEmail = derived(authStore, ($auth) => $auth.userEmail);
export const requiresOAuth = derived(authStore, ($auth) => $auth.authMethod === 'oauth2');
