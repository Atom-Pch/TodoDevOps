<script lang="ts">
	import { goto } from '$app/navigation';
	import { onMount } from 'svelte';

	// Access the environment variable. It defaults to localhost if not set.
	const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080';

	let todos = $state<any[]>([]);
	let newTitle = $state('');
	let newDescription = $state('');
	let error = $state('');
	let imageFile = $state(<FileList | null>(null));
	let isUploading = $state(false);

	// Fetch the To-Dos as soon as the page loads
	onMount(async () => {
		await fetchTodos();
	});

	async function fetchTodos() {
		try {
			const res = await fetch(`${API_URL}/api/todos`, {
				credentials: 'include' // Important for session cookie
			});
			if (res.status === 401) {
				// THE BOUNCER: If Go says unauthorized, kick the user to the login page!
				goto('/login');
				console.log('Not logged in. Redirecting...');
				return; // Stop running this function
			}
			if (res.ok) {
				// If the DB is empty, the API might return null, so we fallback to an empty array
				todos = (await res.json()) || [];
			} else {
				error = 'Failed to load To-Dos from the server.';
			}
		} catch (err) {
			error = 'Could not connect to the API. Is the Go backend running?';
			console.error(err);
		}
	}

	async function addTodo(event: Event) {
		event.preventDefault(); // Prevent the form from refreshing the page
		error = '';
		isUploading = true;
		let finalImageUrl = '';

		try {
			// 1. If an image is selected, handle the S3 upload first
			if (imageFile && imageFile.length > 0) {
				const file = imageFile[0];

				// Get the presigned URL from Go
				const presignRes = await fetch(
					`${API_URL}/api/todos/s3-presign?filename=${encodeURIComponent(file.name)}`,
					{
						credentials: 'include'
					}
				);
				const presignData = await presignRes.json();

				// Upload the file directly to AWS S3
				const uploadRes = await fetch(presignData.upload_url, {
					method: 'PUT',
					body: file,
					headers: {
						'Content-Type': file.type
					}
				});

				if (!uploadRes.ok) throw new Error('Failed to upload image to S3');
				finalImageUrl = presignData.image_url;
			}

			// 2. Save the To-Do item to the Go backend
			const res = await fetch(`${API_URL}/api/todos`, {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				credentials: 'include',
				body: JSON.stringify({
					title: newTitle,
					description: newDescription,
					image_url: finalImageUrl // Send the S3 URL to your DB
				})
			});

			if (res.ok) {
				const newTodo = await res.json();
				todos = [...todos, newTodo];
				newTitle = '';
				newDescription = '';
				imageFile = null; // Clear the file input
			} else {
				error = 'Failed to create To-Do.';
			}
		} catch (err) {
			console.error("Error creating To-Do:", err);
		} finally {
			isUploading = false;
		}
	}

	async function deleteTodo(id: number) {
		try {
			const res = await fetch(`${API_URL}/api/todos/${id}`, {
				method: 'DELETE',
				credentials: 'include' // Must send the cookie!
			});

			if (res.ok) {
				// Instantly remove the deleted item from the UI
				todos = todos.filter((todo) => todo.id !== id);
			} else {
				console.error('Failed to delete task');
			}
		} catch (err) {
			console.error('Could not connect to the API to delete the To-Do.', err);
		}
	}
</script>

<main>
	{#if error}
		<div class="error-banner">{error}</div>
	{/if}

	<form onsubmit={addTodo}>
		<input type="text" placeholder="To-Do Title" bind:value={newTitle} required />
		<input type="text" placeholder="Description (Optional)" bind:value={newDescription} />
		<div class="flex flex-col space-y-1 w-1/4">
			<input type="file" accept="image/*" bind:files={imageFile} />
			<button type="submit">Add</button>
		</div>
	</form>

	<ul>
		{#each todos.toReversed() as todo}
			<li>
				<strong>{todo.title}</strong>
				{#if todo.description}
					<p>{todo.description}</p>
				{/if}
				{#if todo.image_url}
                    <img src={todo.image_url} alt="Task attachment" style="max-width: 200px; margin-top: 10px; border-radius: 4px;" />
                {/if}
			</li>
			<div class="del-block">
				<button class="delete-btn" onclick={() => deleteTodo(todo.id)}>Delete </button>
			</div>
		{/each}
		{#if todos.length === 0 && !error}
			<p class="empty-state">No tasks yet. Create one above!</p>
		{/if}
	</ul>
</main>

<style>
	main {
		max-width: 50vw;
		margin: 3rem auto;
		padding: 0 1rem;
	}
	.error-banner {
		background-color: #fee2e2;
		color: #991b1b;
		padding: 1rem;
		border-radius: 6px;
		margin-bottom: 1rem;
		text-align: center;
	}
	form {
		display: flex;
		gap: 0.5rem;
		margin-bottom: 2rem;
	}
	input {
		flex: 1;
		padding: 0.75rem;
		border: 1px solid #d1d5db;
		border-radius: 6px;
		color: #363636;
		background-color: #e0e0e0;
	}
	button {
		padding: 0.75rem 1.5rem;
		background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
		color: white;
		border: none;
		border-radius: 6px;
		cursor: pointer;
		font-weight: bold;
	}
	ul {
		list-style: none;
		padding: 0;
	}
	li {
		background: #575757;
		padding: 1.5rem;
		border-radius: 8px;
		box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
		border: 1px solid #e5e7eb;
	}
	li strong {
		display: block;
		font-size: 1.1rem;
	}
	li p {
		margin: 0.5rem 0 0 0;
		color: #e0e0e0;
	}
	.empty-state {
		text-align: center;
		color: #97a2b6;
		font-style: italic;
	}
	.delete-btn {
		background: #e0e0e0;
		color: #ef4444;
		border: 1px solid #ef4444;
		padding: 0.5rem 1rem;
		margin-bottom: 1rem;
	}
	.del-block {
		display: flex;
		justify-content: flex-end;
	}
</style>
