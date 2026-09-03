// Interactive Knowledge Graph Canvas Visualization
(function() {
	const canvas = document.getElementById('knowledge-canvas');
	if (!canvas) return;
	const ctx = canvas.getContext('2d');

	function resize() {
		const rect = canvas.getBoundingClientRect();
		canvas.width = rect.width * window.devicePixelRatio;
		canvas.height = rect.height * window.devicePixelRatio;
		ctx.scale(window.devicePixelRatio, window.devicePixelRatio);
	}
	window.addEventListener('resize', resize);
	resize();

	const nodes = [
		{ id: 'linux', label: 'Linux Kernel & eBPF', x: 180, y: 140, vx: 0, vy: 0, r: 24, tag: 'linux', color: '#88c0d0' },
		{ id: 'go', label: 'Go Runtimes & GC', x: 380, y: 120, vx: 0, vy: 0, r: 22, tag: 'go', color: '#81a1c1' },
		{ id: 'rust', label: 'Rust Memory Safety', x: 340, y: 280, vx: 0, vy: 0, r: 20, tag: 'rust', color: '#d08770' },
		{ id: 'storage', label: 'Storage & Btrfs/DB', x: 160, y: 320, vx: 0, vy: 0, r: 20, tag: 'storage', color: '#a3be8c' },
		{ id: 'security', label: 'Zero Trust & LSM', x: 500, y: 220, vx: 0, vy: 0, r: 22, tag: 'security', color: '#bf616a' },
		{ id: 'docker', label: 'Containers & OCI', x: 260, y: 440, vx: 0, vy: 0, r: 18, tag: 'docker', color: '#b48ead' },
		{ id: 'ai', label: 'AI Open-Weight', x: 520, y: 380, vx: 0, vy: 0, r: 18, tag: 'ai', color: '#ebcb8b' }
	];

	const links = [
		{ source: 'linux', target: 'go' },
		{ source: 'linux', target: 'rust' },
		{ source: 'linux', target: 'storage' },
		{ source: 'linux', target: 'security' },
		{ source: 'linux', target: 'docker' },
		{ source: 'go', target: 'security' },
		{ source: 'rust', target: 'security' },
		{ source: 'rust', target: 'linux' },
		{ source: 'docker', target: 'storage' },
		{ source: 'ai', target: 'go' },
		{ source: 'ai', target: 'rust' }
	];

	let selectedNode = null;
	let hoveredNode = null;
	let isDragging = false;
	let dragNode = null;

	function getMousePos(evt) {
		const rect = canvas.getBoundingClientRect();
		return {
			x: evt.clientX - rect.left,
			y: evt.clientY - rect.top
		};
	}

	canvas.addEventListener('mousedown', (e) => {
		const pos = getMousePos(e);
		for (const n of nodes) {
			const dx = pos.x - n.x;
			const dy = pos.y - n.y;
			if (dx * dx + dy * dy < n.r * n.r) {
				isDragging = true;
				dragNode = n;
				selectNode(n);
				return;
			}
		}
	});

	window.addEventListener('mousemove', (e) => {
		const pos = getMousePos(e);
		if (isDragging && dragNode) {
			dragNode.x = pos.x;
			dragNode.y = pos.y;
		} else {
			let found = null;
			for (const n of nodes) {
				const dx = pos.x - n.x;
				const dy = pos.y - n.y;
				if (dx * dx + dy * dy < n.r * n.r) {
					found = n;
					break;
				}
			}
			if (hoveredNode !== found) {
				hoveredNode = found;
				canvas.style.cursor = found ? 'pointer' : 'grab';
			}
		}
	});

	window.addEventListener('mouseup', () => {
		isDragging = false;
		dragNode = null;
	});

	function selectNode(node) {
		selectedNode = node;
		var tagEl = document.getElementById('inspector-tag');
		var titleEl = document.getElementById('inspector-title');
		var descEl = document.getElementById('inspector-desc');
		if (tagEl) tagEl.innerText = 'Domain: ' + node.tag.toUpperCase();
		if (titleEl) titleEl.innerText = node.label;
		if (descEl) descEl.innerText = 'Filtered tech dispatches under #' + node.tag;
	}

	function draw() {
		const rect = canvas.getBoundingClientRect();
		const w = rect.width;
		const h = rect.height;

		ctx.clearRect(0, 0, w, h);

		const isDark = document.documentElement.getAttribute('data-theme') !== 'light';
		const lineColor = isDark ? 'rgba(255, 255, 255, 0.12)' : 'rgba(0, 0, 0, 0.1)';
		const textColor = isDark ? '#eceff4' : '#1a202c';

		ctx.lineWidth = 1.5;
		for (const link of links) {
			const src = nodes.find(n => n.id === link.source);
			const tgt = nodes.find(n => n.id === link.target);
			if (!src || !tgt) continue;

			ctx.beginPath();
			ctx.moveTo(src.x, src.y);
			ctx.lineTo(tgt.x, tgt.y);
			ctx.strokeStyle = lineColor;
			ctx.stroke();
		}

		for (const n of nodes) {
			ctx.beginPath();
			ctx.arc(n.x, n.y, n.r, 0, Math.PI * 2);
			ctx.fillStyle = n.color;
			ctx.fill();

			if (n === selectedNode || n === hoveredNode) {
				ctx.lineWidth = 3;
				ctx.strokeStyle = '#ffffff';
				ctx.stroke();
			}

			ctx.font = '600 11px system-ui, sans-serif';
			ctx.fillStyle = textColor;
			ctx.textAlign = 'center';
			ctx.fillText(n.label, n.x, n.y + n.r + 14);
		}

		requestAnimationFrame(draw);
	}

	draw();
})();
