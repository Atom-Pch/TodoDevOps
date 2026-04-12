<script>
	import { onMount } from 'svelte';

	let message = $state('loading...');

	onMount(async () => {
		try {
			// In docker-compose, we will map port 8080
			const res = await fetch('http://localhost:8080/api/hello');
			const data = await res.json();
			message = data.message;
		} catch (error) {
			message = 'Error connecting to backend.';
			console.error(error);
		}
	});
</script>

<main style="font-family: sans-serif; text-align: center; margin-top: 50px;">
	<h1>DevOps Practice Web App</h1>
	<div style="padding: 20px; background: #f0f0f0; border-radius: 8px; display: inline-block;">
		<p>Backend says: <strong>{message}</strong></p>
	</div>
</main>