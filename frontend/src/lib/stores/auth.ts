import { writable, derived } from 'svelte/store';

export interface AuthState {
	isAuthenticated: boolean;
	userEmail: string | null;
	userAvatarURL: string | null;
	userName: string | null;
	authMethod: 'oauth2' | 'plain' | 'none' | null;
	loading: boolean;
}

function createAuthStore() {
	const { subscribe, update } = writable<AuthState>({
		isAuthenticated: false,
		userEmail: null,
		userAvatarURL: null,
		userName: null,
		authMethod: null,
		loading: true
	});

	return {
		subscribe,
		setAuthenticated: (email: string, avatarURL: string, name: string) =>
			update((state) => ({
				...state,
				isAuthenticated: true,
				userEmail: email,
				userAvatarURL: avatarURL,
				userName: name,
				loading: false
			})),
		setUnauthenticated: () =>
			update((state) => ({
				...state,
				isAuthenticated: false,
				userEmail: null,
				userAvatarURL: null,
				userName: null,
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
				userAvatarURL: null,
				userName: null,
				authMethod: state.authMethod, // Preserve authMethod so UI doesn't disappear
				loading: false
			}))
	};
}

export const authStore = createAuthStore();

// Derived stores for convenience
export const isAuthenticated = derived(authStore, ($auth) => $auth.isAuthenticated);
export const currentUserEmail = derived(authStore, ($auth) => $auth.userEmail);
export const currentUserAvatarURL = derived(authStore, ($auth) => $auth.userAvatarURL);
export const currentUserName = derived(authStore, ($auth) => $auth.userName);
export const authMethod = derived(authStore, ($auth) => $auth.authMethod);
