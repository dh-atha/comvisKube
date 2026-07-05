package dashboard

const pageTemplate = `<!doctype html>
<html lang="en">
<head>
	<meta charset="utf-8">
	<meta name="viewport" content="width=device-width, initial-scale=1">
	<title>Cloud Native Pod Dashboard</title>
	<style>
		:root {
			color-scheme: dark;
			--bg: #07111f;
			--panel: rgba(10, 18, 35, 0.88);
			--panel-border: rgba(148, 163, 184, 0.16);
			--text: #e2e8f0;
			--muted: #94a3b8;
			--accent: #38bdf8;
			--accent-2: #34d399;
			--danger: #fb7185;
		}
		* { box-sizing: border-box; }
		body {
			margin: 0;
			min-height: 100vh;
			font-family: Inter, ui-sans-serif, system-ui, -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif;
			color: var(--text);
			background:
				radial-gradient(circle at top left, rgba(56, 189, 248, 0.22), transparent 36%),
				radial-gradient(circle at bottom right, rgba(52, 211, 153, 0.18), transparent 32%),
				linear-gradient(180deg, #040816 0%, #081121 100%);
		}
		.wrap {
			max-width: 1100px;
			margin: 0 auto;
			padding: 40px 20px 56px;
		}
		.hero {
			display: grid;
			gap: 20px;
			grid-template-columns: 1.4fr 0.9fr;
			align-items: start;
			margin-bottom: 24px;
		}
		.title {
			margin: 0;
			font-size: clamp(2.4rem, 4vw, 4.8rem);
			line-height: 0.95;
			letter-spacing: -0.05em;
		}
		.subtitle {
			margin: 12px 0 0;
			max-width: 60ch;
			color: var(--muted);
			font-size: 1.02rem;
			line-height: 1.7;
		}
		.badge {
			display: inline-flex;
			align-items: center;
			gap: 8px;
			padding: 8px 12px;
			border-radius: 999px;
			background: rgba(56, 189, 248, 0.12);
			color: #bae6fd;
			border: 1px solid rgba(56, 189, 248, 0.2);
			font-size: 0.85rem;
			text-transform: uppercase;
			letter-spacing: 0.12em;
			width: fit-content;
		}
		.panel {
			background: var(--panel);
			border: 1px solid var(--panel-border);
			border-radius: 24px;
			box-shadow: 0 24px 80px rgba(0, 0, 0, 0.35);
			backdrop-filter: blur(18px);
		}
		.summary {
			padding: 22px;
			display: grid;
			gap: 14px;
		}
		.summary strong { font-size: 1.4rem; }
		.metrics {
			display: grid;
			grid-template-columns: repeat(2, minmax(0, 1fr));
			gap: 14px;
		}
		.metric {
			padding: 14px 16px;
			border-radius: 18px;
			background: rgba(148, 163, 184, 0.08);
			border: 1px solid rgba(148, 163, 184, 0.12);
		}
		.metric span {
			display: block;
			color: var(--muted);
			font-size: 0.8rem;
			margin-bottom: 6px;
			text-transform: uppercase;
			letter-spacing: 0.08em;
		}
		.metric strong {
			font-size: 1rem;
			word-break: break-word;
		}
		.grid {
			display: grid;
			grid-template-columns: repeat(12, minmax(0, 1fr));
			gap: 18px;
		}
		.card { padding: 20px; }
		.card h2 {
			margin: 0 0 10px;
			font-size: 1.1rem;
		}
		.controls {
			display: flex;
			flex-wrap: wrap;
			gap: 12px;
			margin-top: 18px;
		}
		button, a.button {
			appearance: none;
			border: 0;
			border-radius: 14px;
			padding: 12px 16px;
			font: inherit;
			font-weight: 600;
			color: #06111f;
			background: linear-gradient(135deg, var(--accent), #7dd3fc);
			cursor: pointer;
			text-decoration: none;
			display: inline-flex;
			align-items: center;
			justify-content: center;
		}
		button.secondary, a.secondary {
			background: rgba(148, 163, 184, 0.14);
			color: var(--text);
			border: 1px solid rgba(148, 163, 184, 0.18);
		}
		.status {
			margin-top: 14px;
			min-height: 24px;
			color: var(--muted);
		}
		pre {
			margin: 0;
			padding: 18px;
			border-radius: 18px;
			background: rgba(2, 6, 23, 0.9);
			overflow-x: auto;
			color: #dbeafe;
			font-size: 0.88rem;
			line-height: 1.65;
			border: 1px solid rgba(148, 163, 184, 0.16);
		}
		.span-7 { grid-column: span 7; }
		.span-5 { grid-column: span 5; }
		.span-12 { grid-column: span 12; }
		@media (max-width: 900px) {
			.hero, .grid { grid-template-columns: 1fr; }
			.span-7, .span-5, .span-12 { grid-column: span 1; }
		}
	</style>
</head>
<body>
	<div class="wrap">
		<div class="hero">
			<div>
				<div class="badge">Kubernetes pod dashboard</div>
				<h1 class="title">Inspect the pod, trigger a crash, and verify restarts from the browser.</h1>
				<p class="subtitle">This app shows the pod hostname and IP address, writes metadata to a persistent volume, and exposes crash and load buttons so you can force a restart or push the HPA toward the second pod.</p>
				<div class="controls">
					<button id="refreshBtn">Refresh status</button>
					<button id="podBtn" class="secondary">Show pod</button>
					<button id="workBtn" class="secondary">Run long load</button>
					<button id="panicBtn" class="secondary">Trigger panic</button>
					<a class="button secondary" href="/state" target="_blank" rel="noreferrer">View raw state</a>
				</div>
				<div id="status" class="status">Loading status...</div>
			</div>
			<div class="panel summary">
				<strong>Live pod facts</strong>
				<div class="metrics">
					<div class="metric"><span>Hostname</span><strong id="hostname">-</strong></div>
					<div class="metric"><span>IP address</span><strong id="ipAddress">-</strong></div>
					<div class="metric"><span>Pod name</span><strong id="podName">-</strong></div>
					<div class="metric"><span>Uptime</span><strong id="uptime">-</strong></div>
				</div>
			</div>
		</div>

		<div class="grid">
			<section class="panel card span-12">
				<details open>
					<summary><h2>Request log</h2></summary>
					<pre id="logOutput">No requests yet.</pre>
				</details>
			</section>
			<section class="panel card span-7">
				<details open>
					<summary><h2>Machine state</h2></summary>
					<pre id="stateOutput">Waiting for data...</pre>
				</details>
			</section>
			<section class="panel card span-5">
				<h2>How to test restart behavior</h2>
				<pre>1. Open / and click "Run long load" a few times.
2. Kubernetes should push toward the second pod when CPU stays high.
3. Click "Trigger panic" to force a restart.
4. Reopen /state to confirm metadata survived on the volume.</pre>
			</section>
		</div>
	</div>

	<script>
		async function loadState() {
			const status = document.getElementById('status');
			const stateOutput = document.getElementById('stateOutput');
			const logOutput = document.getElementById('logOutput');

			status.textContent = 'Refreshing status...';

			try {
				const response = await fetch('/api/info');
				if (!response.ok) {
					throw new Error('status request failed');
				}

				const info = await response.json();
				document.getElementById('hostname').textContent = info.hostname || '-';
				document.getElementById('ipAddress').textContent = info.ip_address || '-';
				document.getElementById('podName').textContent = info.pod_name || '-';
				document.getElementById('uptime').textContent = info.uptime || '-';

				const stateResponse = await fetch('/state');
				if (stateResponse.ok) {
					const state = await stateResponse.json();
					stateOutput.textContent = JSON.stringify(state, null, 2);

					const recent = (state.events || []).slice(0, 12);
					logOutput.textContent = recent.length ? recent.map((event) => JSON.stringify(event, null, 2)).join('\n\n') : 'No events recorded yet.';
				} else {
					stateOutput.textContent = 'Unable to load persisted state.';
				}

				status.textContent = 'Status updated successfully.';
			} catch (error) {
				status.textContent = 'Unable to load status right now.';
				stateOutput.textContent = String(error);
			}
		}

		async function runWork() {
			const status = document.getElementById('status');
			status.textContent = 'Running long load...';

			try {
				const response = await fetch('/work');
				const payload = await response.json();
				status.textContent = payload.message || 'Load completed.';
			} catch (error) {
				status.textContent = 'Load request failed.';
			}
		}

		async function showPod() {
			const status = document.getElementById('status');
			status.textContent = 'Loading pod identity...';

			try {
				const response = await fetch('/pod');
				const payload = await response.json();
				status.textContent = payload.message + ': ' + payload.pod_name;
				await loadState();
			} catch (error) {
				status.textContent = 'Pod request failed.';
			}
		}

		async function runPanic() {
			const status = document.getElementById('status');
			status.textContent = 'Triggering panic...';

			try {
				const response = await fetch('/panic');
				const payload = await response.json();
				status.textContent = payload.message || 'Panic scheduled.';
			} catch (error) {
				status.textContent = 'Panic request failed.';
			}
		}

		document.getElementById('refreshBtn').addEventListener('click', loadState);
		document.getElementById('podBtn').addEventListener('click', showPod);
		document.getElementById('workBtn').addEventListener('click', runWork);
		document.getElementById('panicBtn').addEventListener('click', runPanic);
		loadState();
		setInterval(loadState, 3000);
	</script>
</body>
</html>`
