<script lang="ts">
	import './layout.css';
	import favicon from '$lib/assets/favicon.svg';
	import { onMount } from 'svelte';

	let { children } = $props();

	const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080';

	let currentUser: string | null = $state(null);
	// NEW: Add the loading state, default to true so it blocks the UI immediately
	let isCheckingAuth = $state(true);

	async function handleLogout() {
		try {
			await fetch(`${API_URL}/api/logout`, {
				method: 'POST',
				credentials: 'include'
			});
			console.log('Logged out successfully.');
			window.location.href = '/';
		} catch (err) {
			console.error('Failed to log out', err);
		}
	}

	onMount(async () => {
		// As soon as the app loads, check if we are logged in
		await checkSession();
	});

	// --- CHECK WHO IS LOGGED IN ---
	async function checkSession() {
		try {
			const res = await fetch(`${API_URL}/api/me`, { credentials: 'include' });
			if (res.ok) {
				const data = await res.json();
				currentUser = data.username;
			} else {
				currentUser = null; 
                if (res.status === 404 || res.status === 401) {
                    await fetch(`${API_URL}/api/logout`, { method: 'POST', credentials: 'include' });
                }
			}
		} catch (err) {
            console.error("Auth check failed (offline or backend down?):", err);
            currentUser = null;
        } finally {
			// NEW: Unblock the UI whether the check succeeded or failed
			isCheckingAuth = false;
		}
	}
</script>

<svelte:head><link rel="icon" href={favicon} /></svelte:head>

<style>
	:global(body) {
		background-color: #111827; 
		margin: 0;
	}
</style>

<div class="min-h-screen bg-gray-900 text-gray-100 font-sans selection:bg-indigo-500/30 flex flex-col">
	<nav class="bg-gray-800 border-b border-gray-700 px-4 sm:px-8 py-4 sticky top-0 z-50 shadow-md">
		<div class="max-w-7xl mx-auto flex justify-between items-center">
			<div class="flex items-center">
				<a href="/" class="text-xl font-bold text-transparent bg-clip-text bg-gradient-to-r from-indigo-400 to-purple-500 hover:opacity-80 transition">
					DevOps App
				</a>
			</div>

			<div class="flex items-center space-x-4 sm:space-x-6">
				{#if currentUser}
					<form onsubmit={handleLogout} class="m-0">
						<button type="submit" class="text-sm font-medium text-gray-300 hover:text-white transition bg-gray-700/50 hover:bg-gray-600 px-4 py-2 rounded-lg">
							Logout <span class="opacity-60">({currentUser})</span>
						</button>
					</form>
				{:else}
					<a href="/register" class="text-sm font-medium text-gray-300 hover:text-white transition hidden sm:inline">Register</a>
					<a href="/login" class="text-sm font-medium px-4 py-2 rounded-lg bg-indigo-600 hover:bg-indigo-500 text-white transition shadow-sm">
						Login
					</a>
				{/if}
			</div>
		</div>
	</nav>

	<div class="flex-grow w-full">
		{#if isCheckingAuth}
			<div class="flex flex-col items-center justify-center min-h-[60vh] gap-4">
				<svg class="animate-spin h-10 w-10 text-indigo-500" fill="none" viewBox="0 0 24 24">
					<circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
					<path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
				</svg>
				<p class="text-gray-400 font-medium">Connecting...</p>
			</div>
		{:else}
			{@render children()}
		{/if}
	</div>
</div>