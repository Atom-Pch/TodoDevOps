<script lang="ts">
	import './layout.css';
	import favicon from '$lib/assets/favicon.svg';
	import { onMount } from 'svelte';

	let { children } = $props();

	const API_URL = import.meta.env.VITE_PUBLIC_API_URL || 'http://localhost:8080';

	let currentUser: string | null = $state(null);
	async function handleLogout() {
		try {
			await fetch(`${API_URL}/api/logout`, {
				method: 'POST',
				// CRITICAL: Tells the browser to send the cookie so Go can invalidate it
				credentials: 'include'
			});
			console.log('Logged out successfully.');
			// TODO: Clear local UI state and redirect to login page
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
				currentUser = data.username; // We found the user!
			} else {
				currentUser = null; // No valid session
			}
		} catch (err) {
			currentUser = null;
		}
	}
</script>

<nav>
	<div class="grid grid-cols-3">
		<div class="flex items-center justify-start">
			<a href="/">Home</a>
		</div>
		<div>
			<div class="relative"></div>
		</div>
		<div class="flex items-center justify-end space-x-4">
			{#if currentUser}
				<form onsubmit={handleLogout}>
					<button type="submit">Logout ({currentUser})</button>
				</form>
			{/if}
			<div class="h-6 border-l border-gray-300"></div>
			<a href="/register">Register</a>
			<a href="/login">Login</a>
		</div>
	</div>
</nav>

<svelte:head><link rel="icon" href={favicon} /></svelte:head>
{@render children()}

<style>
	/* Clean, minimal styling for the UI */
	:global(body) {
		font-family:
			system-ui,
			-apple-system,
			sans-serif;
		background-color: #1a1a1a;
		color: #e0e0e0;
	}
	nav {
		background-color: #363636;
		border-bottom: 1px solid #dee2e6;
		padding: 1rem 5rem;
	}
	button, a {
		cursor: pointer;
		color: #e0e0e0;
	}
</style>
