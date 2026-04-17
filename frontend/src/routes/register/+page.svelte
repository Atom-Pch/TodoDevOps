<script lang="ts">
	// Pull the API URL from your environment variables
	const API_URL = import.meta.env.VITE_PUBLIC_API_URL || 'http://localhost:8080';

	// State variables for your UI to bind to (adjust if you named yours differently)
	let username = $state('');
	let email = $state('');
	let password = $state('');
	let errorMessage = $state('');

	// --- REGISTRATION LOGIC ---
	async function handleRegister(event: Event) {
		if (event) event.preventDefault();
		errorMessage = '';

		try {
			const res = await fetch(`${API_URL}/api/register`, {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json'
				},
				// We don't strictly need credentials: 'include' for registration,
				// but it's good practice to keep the payload consistent.
				body: JSON.stringify({ username, email, password })
			});

			if (res.ok) {
				console.log('Registration successful! You can now log in.');
				// Optional: Automatically trigger handleLogin here, or redirect to a login view
				window.location.href = '/login'; // Redirect to login page after successful registration
			} else {
				const data = await res.text();
				errorMessage = `Registration failed: ${data}`;
			}
		} catch (err) {
			errorMessage = 'Network error. Is the Go server running?';
			console.error(err);
		}
	}
</script>

<div class="reg-container">
	<div class="reg-card">
		<h1>Register</h1>

		{#if errorMessage}
			<div class="error-message">{errorMessage}</div>
		{/if}

		<form onsubmit={handleRegister}>
			<div class="form-group">
				<label for="email">Email</label>
				<input type="email" id="email" bind:value={email} placeholder="Enter your email" required />
			</div>

			<div class="form-group">
				<label for="username">Username</label>
				<input
					type="text"
					id="username"
					bind:value={username}
					placeholder="Enter your username"
					required
				/>
			</div>

			<div class="form-group">
				<label for="password">Password</label>
				<input
					type="password"
					id="password"
					bind:value={password}
					placeholder="Enter your password"
					required
				/>
			</div>

			<button type="submit" class="reg-btn">Register</button>
		</form>
	</div>
</div>

<style>
	.reg-container {
		display: flex;
		justify-content: center;
		align-items: center;
		min-height: 60vh;
	}

	.reg-card {
		background: #363636;
		padding: 2rem;
		border-radius: 8px;
		box-shadow: 0 10px 25px rgba(0, 0, 0, 0.2);
		width: 100%;
		max-width: 400px;
	}

	h1 {
		text-align: center;
		color: #e0e0e0;
		margin-bottom: 1.5rem;
		font-size: 1.5rem;
	}

	.error-message {
		background-color: #fee;
		color: #c33;
		padding: 0.75rem;
		border-radius: 4px;
		margin-bottom: 1rem;
		text-align: center;
	}

	.form-group {
		margin-bottom: 1.5rem;
	}

	label {
		display: block;
		margin-bottom: 0.5rem;
		color: #e0e0e0;
		font-weight: 500;
	}

	input {
		width: 100%;
		padding: 0.75rem;
		border: 1px solid #ddd;
		border-radius: 4px;
		font-size: 1rem;
		box-sizing: border-box;
		transition: border-color 0.3s;
		color: #363636;
	}

	input:focus {
		outline: none;
		border-color: #667eea;
		box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
	}

	.reg-btn {
		width: 100%;
		padding: 0.75rem;
		background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
		color: white;
		border: none;
		border-radius: 4px;
		font-size: 1rem;
		font-weight: 600;
		cursor: pointer;
		transition:
			transform 0.2s,
			box-shadow 0.2s;
	}

	.reg-btn:hover {
		transform: translateY(-2px);
		box-shadow: 0 5px 15px rgba(102, 126, 234, 0.4);
	}

	.reg-btn:active {
		transform: translateY(0);
	}
</style>
