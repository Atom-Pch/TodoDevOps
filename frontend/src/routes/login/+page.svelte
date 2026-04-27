<script lang="ts">
	// Access the environment variable. It defaults to localhost if not set.
	const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080';

	let username = $state('');
	let password = $state('');
	let errorMessage = $state('');

	async function handleLogin(event: Event) {
        if (event) event.preventDefault();
        errorMessage = '';

        try {
            const res = await fetch(`${API_URL}/api/login`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                // CRITICAL DEVOPS DETAIL: This tells the browser to accept the Set-Cookie header 
                // from the Go backend and store the session token securely.
                credentials: 'include',
                body: JSON.stringify({ username, password })
            });

            if (res.ok) {
                console.log("Login successful! Cookie saved.");
                // TODO: Redirect the user to the To-Do dashboard
                window.location.href = '/todos';
            } else {
                errorMessage = 'Invalid username or password.';
            }
        } catch (err) {
            errorMessage = 'Network error. You\'re offline or server not running.';
            console.error(err);
        }
    }
</script>

<div class="min-h-[80vh] flex items-center justify-center px-4 py-12 sm:px-6 lg:px-8">
	<div class="w-full max-w-md bg-gray-800 rounded-2xl shadow-2xl border border-gray-700 p-8 sm:p-10">
		<h1 class="text-3xl font-bold text-center text-white mb-8">Welcome Back</h1>

		{#if errorMessage}
			<div class="mb-6 bg-red-900/50 border border-red-500/50 text-red-200 px-4 py-3 rounded-lg text-sm text-center shadow-sm flex items-center justify-center gap-2">
				<svg class="w-5 h-5 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path></svg>
				{errorMessage}
			</div>
		{/if}

		<form onsubmit={handleLogin} class="space-y-6">
			<div>
				<label for="username" class="block text-sm font-medium text-gray-300 mb-2">Username</label>
				<input
					type="text"
					id="username"
					bind:value={username}
					placeholder="Enter your username"
					required
					class="w-full px-4 py-3 bg-gray-900 border border-gray-700 rounded-xl text-white placeholder-gray-500 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-transparent transition"
				/>
			</div>

			<div>
				<label for="password" class="block text-sm font-medium text-gray-300 mb-2">Password</label>
				<input
					type="password"
					id="password"
					bind:value={password}
					placeholder="Enter your password"
					required
					class="w-full px-4 py-3 bg-gray-900 border border-gray-700 rounded-xl text-white placeholder-gray-500 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-transparent transition"
				/>
			</div>

			<button type="submit" class="w-full flex justify-center py-3.5 px-4 border border-transparent rounded-xl shadow-md text-base font-medium text-white bg-gradient-to-r from-indigo-600 to-purple-600 hover:opacity-90 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 focus:ring-offset-gray-900 transition-all transform hover:-translate-y-0.5">
				Log in
			</button>
		</form>
		
		<p class="mt-8 text-center text-sm text-gray-400">
			Don't have an account? 
			<a href="/register" class="font-medium text-indigo-400 hover:text-indigo-300 transition hover:underline">Register here</a>
		</p>
	</div>
</div>