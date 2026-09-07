(function () {
	'use strict';

	const canvas = document.getElementById('knowledge-canvas');
	if (!canvas) return;

	const lang = document.documentElement.lang === 'id' ? 'id' : 'en';
	const blogPrefix = lang === 'id' ? '/id' : '';

	let rawData = null;
	const graphAttr = canvas.getAttribute('data-graph');
	if (graphAttr) {
		try {
			rawData = JSON.parse(graphAttr);
		} catch (e) {
			console.error('Failed to parse canvas data-graph:', e);
		}
	}

	async function boot() {
		if (!rawData || !rawData.nodes || rawData.nodes.length === 0) {
			try {
				const res = await fetch(blogPrefix + '/api/graph');
				if (res.ok) {
					rawData = await res.json();
				}
			} catch (err) {
				console.error('Failed to fetch /api/graph:', err);
			}
		}

		if (!rawData || !rawData.nodes || rawData.nodes.length === 0) {
			console.warn('No knowledge graph data available.');
			return;
		}

		initGraph(rawData);
	}

	function initGraph(rawData) {
		const ctx = canvas.getContext('2d');

		const tagPalette = {
			linux: '#38bdf8',
			kernel: '#0284c7',
			ebpf: '#a855f7',
			go: '#06b6d4',
			golang: '#06b6d4',
			rust: '#f97316',
			architecture: '#6366f1',
			backend: '#3b82f6',
			security: '#f43f5e',
			storage: '#10b981',
			database: '#14b8a6',
			performance: '#eab308',
			memory: '#ec4899',
			concurrency: '#8b5cf6',
			networking: '#0ea5e9'
		};

		function hexToRgba(hex, alpha = 1.0) {
			if (!hex || hex.startsWith('rgb')) return hex || `rgba(148, 163, 184, ${alpha})`;
			let clean = hex.replace('#', '');
			if (clean.length === 3) clean = clean.split('').map(c => c + c).join('');
			const num = parseInt(clean, 16);
			const r = (num >> 16) & 255;
			const g = (num >> 8) & 255;
			const b = num & 255;
			return `rgba(${r}, ${g}, ${b}, ${alpha})`;
		}

		const adjacency = new Map();
		rawData.nodes.forEach(n => adjacency.set(n.id, new Set()));
		(rawData.links || []).forEach(l => {
			if (adjacency.has(l.source) && adjacency.has(l.target)) {
				adjacency.get(l.source).add(l.target);
				adjacency.get(l.target).add(l.source);
			}
		});

		const nodes = rawData.nodes.map(n => {
			const isTag = n.type === 'tag';
			const deg = adjacency.get(n.id) ? adjacency.get(n.id).size : 0;
			const cleanTag = (n.tag || '').toLowerCase();
			const customColor = tagPalette[cleanTag] || n.color;

			const r = isTag
				? Math.min(10.5, 5.5 + Math.sqrt(deg) * 1.5)
				: Math.min(6.2, 3.2 + Math.sqrt(deg) * 0.9);

			return {
				...n,
				x: 0,
				y: 0,
				vx: 0,
				vy: 0,
				fx: null,
				fy: null,
				r: r,
				color: customColor || (isTag ? '#38bdf8' : null)
			};
		});

		const nodeMap = new Map();
		nodes.forEach(n => nodeMap.set(n.id, n));

		const links = (rawData.links || [])
			.map(l => ({
				source: nodeMap.get(l.source),
				target: nodeMap.get(l.target),
				weight: l.weight || 1.0,
				kind: l.kind || 'tag'
			}))
			.filter(l => l.source && l.target);

		const transform = {
			x: 0,
			y: 0,
			scale: 1.0
		};

		let width = 800;
		let height = 560;

		function resize() {
			const rect = canvas.getBoundingClientRect();
			width = rect.width;
			height = rect.height;
			canvas.width = width * window.devicePixelRatio;
			canvas.height = height * window.devicePixelRatio;
			ctx.setTransform(1, 0, 0, 1, 0, 0);
			ctx.scale(window.devicePixelRatio, window.devicePixelRatio);
		}
		window.addEventListener('resize', resize);
		resize();

		const tagNodes = nodes.filter(n => n.type === 'tag');
		const tagAngleMap = new Map();
		const tagRadius = Math.min(width, height) * 0.32;

		tagNodes.forEach((tn, i) => {
			const angle = (i / Math.max(1, tagNodes.length)) * Math.PI * 2 - Math.PI / 2;
			tn.x = width / 2 + Math.cos(angle) * tagRadius;
			tn.y = height / 2 + Math.sin(angle) * tagRadius;
			tagAngleMap.set(tn.tag, { x: tn.x, y: tn.y, angle });
		});

		nodes.filter(n => n.type === 'post').forEach((pn, i) => {
			const parentTag = tagAngleMap.get(pn.tag);
			if (parentTag) {
				const jitterAngle = parentTag.angle + (Math.random() - 0.5) * 1.5;
				const dist = 50 + Math.random() * 75;
				pn.x = width / 2 + Math.cos(jitterAngle) * (tagRadius + dist * 0.7);
				pn.y = height / 2 + Math.sin(jitterAngle) * (tagRadius + dist * 0.7);
			} else {
				const angle = (i / Math.max(1, nodes.length)) * Math.PI * 2;
				pn.x = width / 2 + Math.cos(angle) * 130 + (Math.random() - 0.5) * 30;
				pn.y = height / 2 + Math.sin(angle) * 130 + (Math.random() - 0.5) * 30;
			}
		});

		let alpha = 1.0;
		const alphaMin = 0.0005;
		const alphaDecay = 0.994;
		const springK = 0.032;
		const repulsionK = 1150;
		const centerGravity = 0.012;
		const friction = 0.88;

		function tickPhysics() {
			if (alpha < alphaMin) return;

			const cx = width / 2;
			const cy = height / 2;

			for (let i = 0; i < links.length; i++) {
				const l = links[i];
				const s = l.source;
				const t = l.target;

				const dx = t.x - s.x;
				const dy = t.y - s.y;
				const dist = Math.sqrt(dx * dx + dy * dy) || 1;

				const targetDist = l.kind === 'crosslink' ? 80 : 60;
				const force = (dist - targetDist) * springK * l.weight;

				const fx = (dx / dist) * force;
				const fy = (dy / dist) * force;

				if (s.fx === null) { s.vx += fx; s.vy += fy; }
				if (t.fx === null) { t.vx -= fx; t.vy -= fy; }
			}

			for (let i = 0; i < nodes.length; i++) {
				const n1 = nodes[i];
				for (let j = i + 1; j < nodes.length; j++) {
					const n2 = nodes[j];
					const dx = n2.x - n1.x;
					const dy = n2.y - n1.y;
					const distSq = dx * dx + dy * dy || 1;
					const dist = Math.sqrt(distSq);

					const repulse = repulsionK / (distSq + 90);
					let rx = (dx / dist) * repulse;
					let ry = (dy / dist) * repulse;

					const minDist = n1.r + n2.r + 7;
					if (dist < minDist) {
						const overlap = (minDist - dist) * 0.5;
						rx += (dx / dist) * overlap * 0.85;
						ry += (dy / dist) * overlap * 0.85;
					}

					if (n1.fx === null) { n1.vx -= rx; n1.vy -= ry; }
					if (n2.fx === null) { n2.vx += rx; n2.vy += ry; }
				}
			}

			for (let i = 0; i < nodes.length; i++) {
				const n = nodes[i];
				if (n.fx !== null) {
					n.x = n.fx;
					n.y = n.fy;
					n.vx = 0;
					n.vy = 0;
					continue;
				}

				n.vx += (cx - n.x) * centerGravity;
				n.vy += (cy - n.y) * centerGravity;

				n.vx *= friction;
				n.vy *= friction;

				n.x += n.vx * alpha;
				n.y += n.vy * alpha;
			}

			alpha *= alphaDecay;
		}

		function wakeSimulation(val) {
			alpha = Math.max(alpha, val || 0.4);
		}

		function screenToWorld(sx, sy) {
			return {
				x: (sx - transform.x) / transform.scale,
				y: (sy - transform.y) / transform.scale
			};
		}

		function getNodeAtScreenPos(sx, sy) {
			const wPos = screenToWorld(sx, sy);

			for (let i = 0; i < nodes.length; i++) {
				const n = nodes[i];
				const dx = wPos.x - n.x;
				const dy = wPos.y - n.y;
				const hitRadius = Math.max(n.r + 6, 14);
				if (dx * dx + dy * dy <= hitRadius * hitRadius) {
					return n;
				}
			}
			return null;
		}

		let hoveredNode = null;
		let selectedNode = null;
		let isDraggingNode = false;
		let draggedNode = null;
		let isPanning = false;
		let panStartX = 0;
		let initialTransformX = 0;
		let initialTransformY = 0;
		let searchQuery = '';
		let clickStartX = 0;
		let clickStartY = 0;

		canvas.addEventListener('mousedown', (e) => {
			const rect = canvas.getBoundingClientRect();
			const sx = e.clientX - rect.left;
			const sy = e.clientY - rect.top;

			clickStartX = e.clientX;
			clickStartY = e.clientY;

			const hit = getNodeAtScreenPos(sx, sy);
			if (hit) {
				isDraggingNode = true;
				draggedNode = hit;
				const wPos = screenToWorld(sx, sy);
				hit.fx = wPos.x;
				hit.fy = wPos.y;
				selectNode(hit);
				wakeSimulation(0.35);
				canvas.style.cursor = 'grabbing';
			} else {
				isPanning = true;
				panStartX = e.clientX;
				panStartY = e.clientY;
				initialTransformX = transform.x;
				initialTransformY = transform.y;
				canvas.style.cursor = 'grabbing';
			}
		});

		window.addEventListener('mousemove', (e) => {
			const rect = canvas.getBoundingClientRect();
			const sx = e.clientX - rect.left;
			const sy = e.clientY - rect.top;

			if (isDraggingNode && draggedNode) {
				const wPos = screenToWorld(sx, sy);
				draggedNode.fx = wPos.x;
				draggedNode.fy = wPos.y;
				wakeSimulation(0.2);
			} else if (isPanning) {
				transform.x = initialTransformX + (e.clientX - panStartX);
				transform.y = initialTransformY + (e.clientY - panStartY);
			} else {
				if (sx >= 0 && sx <= width && sy >= 0 && sy <= height) {
					const hit = getNodeAtScreenPos(sx, sy);
					if (hoveredNode !== hit) {
						hoveredNode = hit;
						canvas.style.cursor = hit ? 'pointer' : 'grab';
					}
				} else if (hoveredNode) {
					hoveredNode = null;
				}
			}
		});

		window.addEventListener('mouseup', (e) => {
			const distMoved = Math.hypot(e.clientX - clickStartX, e.clientY - clickStartY);

			if (isDraggingNode && draggedNode) {
				draggedNode.fx = null;
				draggedNode.fy = null;
				isDraggingNode = false;
				draggedNode = null;
				wakeSimulation(0.15);
			}

			if (isPanning) {
				isPanning = false;

				if (distMoved < 5) {
					const rect = canvas.getBoundingClientRect();
					const sx = e.clientX - rect.left;
					const sy = e.clientY - rect.top;
					if (sx >= 0 && sx <= width && sy >= 0 && sy <= height) {
						const hit = getNodeAtScreenPos(sx, sy);
						if (!hit && selectedNode) {
							selectNode(null, false);
						}
					}
				}
			}
			canvas.style.cursor = hoveredNode ? 'pointer' : 'grab';
		});

		canvas.addEventListener('wheel', (e) => {
			e.preventDefault();
			const rect = canvas.getBoundingClientRect();
			const sx = e.clientX - rect.left;
			const sy = e.clientY - rect.top;

			const zoomFactor = e.deltaY < 0 ? 1.12 : 0.89;
			const newScale = Math.max(0.4, Math.min(3.2, transform.scale * zoomFactor));

			transform.x = sx - (sx - transform.x) * (newScale / transform.scale);
			transform.y = sy - (sy - transform.y) * (newScale / transform.scale);
			transform.scale = newScale;
		}, { passive: false });

		canvas.addEventListener('dblclick', (e) => {
			const rect = canvas.getBoundingClientRect();
			const sx = e.clientX - rect.left;
			const sy = e.clientY - rect.top;
			const hit = getNodeAtScreenPos(sx, sy);
			if (hit) {
				focusNode(hit);
			} else {
				resetCamera();
			}
		});

		let lastTouchDist = 0;
		let lastTouchX = 0;
		let lastTouchY = 0;
		let touchStartX = 0;
		let touchStartY = 0;
		let touchMoved = false;

		canvas.addEventListener('touchstart', (e) => {
			if (e.touches.length === 1) {
				const rect = canvas.getBoundingClientRect();
				const sx = e.touches[0].clientX - rect.left;
				const sy = e.touches[0].clientY - rect.top;
				touchStartX = e.touches[0].clientX;
				touchStartY = e.touches[0].clientY;
				touchMoved = false;
				const hit = getNodeAtScreenPos(sx, sy);
				if (hit) {
					isDraggingNode = true;
					draggedNode = hit;
					const wPos = screenToWorld(sx, sy);
					hit.fx = wPos.x;
					hit.fy = wPos.y;
					selectNode(hit);
					wakeSimulation(0.35);
				} else {
					isPanning = true;
					lastTouchX = e.touches[0].clientX;
					lastTouchY = e.touches[0].clientY;
					initialTransformX = transform.x;
					initialTransformY = transform.y;
				}
			} else if (e.touches.length === 2) {
				isDraggingNode = false;
				isPanning = false;
				touchMoved = true;
				const dx = e.touches[0].clientX - e.touches[1].clientX;
				const dy = e.touches[0].clientY - e.touches[1].clientY;
				lastTouchDist = Math.sqrt(dx * dx + dy * dy);
			}
		}, { passive: true });

		canvas.addEventListener('touchmove', (e) => {
			if (e.touches.length === 1) {
				const dist = Math.hypot(e.touches[0].clientX - touchStartX, e.touches[0].clientY - touchStartY);
				if (dist > 8) touchMoved = true;
				const rect = canvas.getBoundingClientRect();
				const sx = e.touches[0].clientX - rect.left;
				const sy = e.touches[0].clientY - rect.top;
				if (isDraggingNode && draggedNode) {
					const wPos = screenToWorld(sx, sy);
					draggedNode.fx = wPos.x;
					draggedNode.fy = wPos.y;
					wakeSimulation(0.2);
				} else if (isPanning) {
					transform.x = initialTransformX + (e.touches[0].clientX - lastTouchX);
					transform.y = initialTransformY + (e.touches[0].clientY - lastTouchY);
				}
			} else if (e.touches.length === 2) {
				touchMoved = true;
				const dx = e.touches[0].clientX - e.touches[1].clientX;
				const dy = e.touches[0].clientY - e.touches[1].clientY;
				const dist = Math.sqrt(dx * dx + dy * dy);
				if (lastTouchDist > 0) {
					const factor = dist / lastTouchDist;
					const newScale = Math.max(0.4, Math.min(3.2, transform.scale * factor));
					const midX = (e.touches[0].clientX + e.touches[1].clientX) / 2;
					const midY = (e.touches[0].clientY + e.touches[1].clientY) / 2;
					const rect = canvas.getBoundingClientRect();
					const sx = midX - rect.left;
					const sy = midY - rect.top;
					transform.x = sx - (sx - transform.x) * (newScale / transform.scale);
					transform.y = sy - (sy - transform.y) * (newScale / transform.scale);
					transform.scale = newScale;
				}
				lastTouchDist = dist;
			}
		}, { passive: true });

		canvas.addEventListener('touchend', () => {
			if (draggedNode) {
				draggedNode.fx = null;
				draggedNode.fy = null;
				draggedNode = null;
				isDraggingNode = false;
			}
			if (!touchMoved && isPanning && selectedNode) {
				selectNode(null, false);
			}
			isPanning = false;
			lastTouchDist = 0;
		});

		function resetCamera() {
			let minX = Infinity, maxX = -Infinity, minY = Infinity, maxY = -Infinity;
			nodes.forEach(n => {
				if (n.x < minX) minX = n.x;
				if (n.x > maxX) maxX = n.x;
				if (n.y < minY) minY = n.y;
				if (n.y > maxY) maxY = n.y;
			});

			const graphW = (maxX - minX) || 400;
			const graphH = (maxY - minY) || 300;
			const padding = 70;

			const scaleX = (width - padding * 2) / graphW;
			const scaleY = (height - padding * 2) / graphH;
			const targetScale = Math.max(0.45, Math.min(1.15, Math.min(scaleX, scaleY)));

			const targetCenterX = (minX + maxX) / 2;
			const targetCenterY = (minY + maxY) / 2;

			animateCamera(
				width / 2 - targetCenterX * targetScale,
				height / 2 - targetCenterY * targetScale,
				targetScale
			);
		}

		function focusNode(node) {
			const targetScale = Math.max(transform.scale, 1.25);
			const targetX = width / 2 - node.x * targetScale;
			const targetY = height / 2 - node.y * targetScale;
			animateCamera(targetX, targetY, targetScale);
		}

		function animateCamera(toX, toY, toScale, duration = 320) {
			const startX = transform.x;
			const startY = transform.y;
			const startScale = transform.scale;
			const startTime = performance.now();

			function step(now) {
				const elapsed = now - startTime;
				const progress = Math.min(1, elapsed / duration);
				const ease = 1 - Math.pow(1 - progress, 3);

				transform.x = startX + (toX - startX) * ease;
				transform.y = startY + (toY - startY) * ease;
				transform.scale = startScale + (toScale - startScale) * ease;

				if (progress < 1) {
					requestAnimationFrame(step);
				}
			}
			requestAnimationFrame(step);
		}

		const btnReset = document.getElementById('btn-reset-view');
		if (btnReset) btnReset.addEventListener('click', resetCamera);

		const btnZoomIn = document.getElementById('btn-zoom-in');
		if (btnZoomIn) {
			btnZoomIn.addEventListener('click', () => {
				const newScale = Math.min(3.2, transform.scale * 1.25);
				transform.x = width / 2 - (width / 2 - transform.x) * (newScale / transform.scale);
				transform.y = height / 2 - (height / 2 - transform.y) * (newScale / transform.scale);
				transform.scale = newScale;
			});
		}

		const btnZoomOut = document.getElementById('btn-zoom-out');
		if (btnZoomOut) {
			btnZoomOut.addEventListener('click', () => {
				const newScale = Math.max(0.4, transform.scale * 0.8);
				transform.x = width / 2 - (width / 2 - transform.x) * (newScale / transform.scale);
				transform.y = height / 2 - (height / 2 - transform.y) * (newScale / transform.scale);
				transform.scale = newScale;
			});
		}

		const nodeRows = document.querySelectorAll('.graph-node-row');
		const contextBar = document.getElementById('graph-context-bar');
		const emptyBar = document.getElementById('graph-empty-bar');
		const inspectorBadge = document.getElementById('inspector-badge');
		const inspectorMetaTag = document.getElementById('inspector-meta-tag');
		const inspectorMetaInfo = document.getElementById('inspector-meta-info');
		const inspectorTitle = document.getElementById('inspector-title');
		const inspectorCta = document.getElementById('inspector-cta');
		const inspectorCtaLabel = document.getElementById('inspector-cta-label');
		const btnClearSelection = document.getElementById('btn-clear-selection');

		function selectNode(node, scrollList = true) {
			selectedNode = node;

			nodeRows.forEach(r => r.classList.remove('bg-bg'));

			if (!node) {
				if (contextBar) contextBar.classList.add('hidden');
				if (emptyBar) emptyBar.classList.remove('hidden');
				return;
			}

			const matchingRow = document.querySelector(`.graph-node-row[data-node-id="${node.id}"]`);
			if (matchingRow) {
				matchingRow.classList.add('bg-bg');
				if (scrollList) {
					matchingRow.scrollIntoView({ block: 'nearest', behavior: 'smooth' });
				}
			}

			if (contextBar) contextBar.classList.remove('hidden');
			if (emptyBar) emptyBar.classList.add('hidden');

			const isTag = node.type === 'tag';

			if (inspectorBadge) {
				inspectorBadge.textContent = isTag ? 'DOMAIN' : 'DISPATCH';
			}

			if (inspectorMetaTag) {
				inspectorMetaTag.textContent = isTag ? ('#' + node.tag) : (node.tag ? '#' + node.tag.toUpperCase() : 'POST');
			}

			if (inspectorMetaInfo) {
				if (!isTag) {
					const parts = [];
					if (node.date) parts.push(node.date);
					if (node.readTime) parts.push(node.readTime + ' min read');
					inspectorMetaInfo.textContent = parts.join(' · ');
				} else {
					const connCount = adjacency.get(node.id) ? adjacency.get(node.id).size : 0;
					inspectorMetaInfo.textContent = connCount + (lang === 'id' ? ' tulisan terhubung' : ' dispatches connected');
				}
			}

			if (inspectorTitle) {
				inspectorTitle.textContent = node.label;
			}

			if (inspectorCta) {
				if (!isTag) {
					inspectorCta.href = blogPrefix + '/blog/' + node.slug;
					if (inspectorCtaLabel) inspectorCtaLabel.textContent = lang === 'id' ? 'BACA ARTIKEL' : 'READ DISPATCH';
				} else {
					inspectorCta.href = blogPrefix + '/blog/tag/' + node.tag;
					if (inspectorCtaLabel) inspectorCtaLabel.textContent = lang === 'id' ? ('JELAJAHI #' + node.tag.toUpperCase()) : ('EXPLORE #' + node.tag.toUpperCase());
				}
			}
		}

		nodeRows.forEach(row => {
			row.addEventListener('click', () => {
				const nodeId = row.getAttribute('data-node-id');
				const node = nodeMap.get(nodeId);
				if (node) {
					selectNode(node, false);
					focusNode(node);
					wakeSimulation(0.25);
				}
			});
		});

		if (btnClearSelection) {
			btnClearSelection.addEventListener('click', () => {
				selectNode(null);
				resetCamera();
			});
		}

		let currentFilter = 'all';
		const filterTabs = document.querySelectorAll('.graph-filter-tab');
		const listCountEl = document.getElementById('graph-list-count');

		function applyFilter(filterType) {
			currentFilter = filterType;
			let visibleCount = 0;

			filterTabs.forEach(tab => {
				const isTarget = tab.getAttribute('data-filter') === filterType;
				if (isTarget) {
					tab.className = 'graph-filter-tab transition-colors tracking-wider text-text font-bold underline underline-offset-4 decoration-text cursor-pointer';
				} else {
					tab.className = 'graph-filter-tab transition-colors tracking-wider text-muted hover:text-text cursor-pointer';
				}
			});

			nodeRows.forEach(row => {
				const rowType = row.getAttribute('data-node-type');
				const rowTitle = (row.getAttribute('data-node-title') || '').toLowerCase();
				const rowTag = (row.getAttribute('data-node-tag') || '').toLowerCase();

				const matchesType = (filterType === 'all') ||
					(filterType === 'posts' && rowType === 'post') ||
					(filterType === 'domains' && rowType === 'tag');

				const matchesSearch = !searchQuery ||
					rowTitle.includes(searchQuery) ||
					rowTag.includes(searchQuery);

				if (matchesType && matchesSearch) {
					row.style.display = '';
					visibleCount++;
				} else {
					row.style.display = 'none';
				}
			});

			if (listCountEl) {
				listCountEl.textContent = visibleCount + ' ' + (lang === 'id' ? 'simpul' : 'nodes');
			}
		}

		filterTabs.forEach(tab => {
			tab.addEventListener('click', () => {
				applyFilter(tab.getAttribute('data-filter'));
			});
		});

		const searchInput = document.getElementById('graph-search');
		if (searchInput) {
			searchInput.addEventListener('input', (e) => {
				searchQuery = e.target.value.trim().toLowerCase();
				applyFilter(currentFilter);

				if (searchQuery) {
					const match = nodes.find(n =>
						n.label.toLowerCase().includes(searchQuery) ||
						(n.tag && n.tag.toLowerCase().includes(searchQuery))
					);
					if (match) {
						selectNode(match, true);
						focusNode(match);
					}
				}
			});
		}

		const graphShortcutsModal = document.getElementById('graph-shortcuts-modal');
		const btnShortcuts = document.getElementById('btn-graph-shortcuts');

		function toggleShortcutsModal(show) {
			if (!graphShortcutsModal) return;
			if (typeof show === 'boolean') {
				graphShortcutsModal.classList.toggle('open', show);
			} else {
				graphShortcutsModal.classList.toggle('open');
			}
		}

		if (btnShortcuts) {
			btnShortcuts.addEventListener('click', () => toggleShortcutsModal(true));
		}

		window.addEventListener('keydown', (e) => {
			const modalOpen = graphShortcutsModal && graphShortcutsModal.classList.contains('open');

			if (e.target.matches('input, textarea, select')) {
				if (e.key === 'Escape') {
					e.target.blur();
					if (searchQuery) {
						searchQuery = '';
						if (searchInput) searchInput.value = '';
						applyFilter(currentFilter);
					}
				} else if (e.key === 'Enter') {
					e.target.blur();
					const visibleRows = Array.from(nodeRows).filter(r => r.style.display !== 'none');
					if (visibleRows.length > 0 && !selectedNode) {
						visibleRows[0].click();
					}
				}
				return;
			}

			if (e.metaKey || e.ctrlKey || e.altKey) return;

			if (e.key === 'Escape') {
				if (modalOpen) {
					toggleShortcutsModal(false);
					return;
				}
				if (selectedNode) {
					selectNode(null);
				} else {
					resetCamera();
				}
				return;
			}

			if (e.key === '?') {
				e.preventDefault();
				e.stopPropagation();
				toggleShortcutsModal();
				return;
			}

			if (modalOpen) return;

			if (e.key === '/') {
				if (searchInput) {
					e.preventDefault();
					e.stopPropagation();
					searchInput.focus();
					searchInput.select();
				}
				return;
			}

			if (e.key === 'Enter' || e.key === 'o') {
				if (selectedNode) {
					e.preventDefault();
					if (selectedNode.type === 'post') {
						window.location.href = blogPrefix + '/blog/' + selectedNode.slug;
					} else if (selectedNode.type === 'tag') {
						window.location.href = blogPrefix + '/blog/tag/' + selectedNode.tag;
					}
				}
				return;
			}

			if (e.key === '1' || e.key === 'a') {
				e.preventDefault();
				applyFilter('all');
				return;
			} else if (e.key === '2' || e.key === 'p') {
				e.preventDefault();
				applyFilter('posts');
				return;
			} else if (e.key === '3' || e.key === 'd' || e.key === 't') {
				e.preventDefault();
				applyFilter('domains');
				return;
			}

			if (e.key === '+' || e.key === '=' || e.code === 'NumpadAdd') {
				e.preventDefault();
				if (btnZoomIn) btnZoomIn.click();
				return;
			} else if (e.key === '-' || e.key === '_' || e.code === 'NumpadSubtract') {
				e.preventDefault();
				if (btnZoomOut) btnZoomOut.click();
				return;
			} else if (e.key === '0' || e.key === 'r') {
				e.preventDefault();
				resetCamera();
				return;
			} else if (e.key === 'f') {
				e.preventDefault();
				if (selectedNode) {
					focusNode(selectedNode);
				} else {
					resetCamera();
				}
				return;
			}

			if (e.key === 'ArrowLeft') {
				e.preventDefault();
				transform.x += 45;
				return;
			} else if (e.key === 'ArrowRight') {
				e.preventDefault();
				transform.x -= 45;
				return;
			}

			if (e.key === 'j' || e.key === 'ArrowDown' || e.key === 'k' || e.key === 'ArrowUp') {
				const visibleRows = Array.from(nodeRows).filter(r => r.style.display !== 'none');
				if (visibleRows.length === 0) return;

				const isNext = e.key === 'j' || e.key === 'ArrowDown';
				let activeIdx = visibleRows.findIndex(r => r.classList.contains('bg-bg'));
				if (activeIdx === -1) {
					activeIdx = isNext ? 0 : visibleRows.length - 1;
				} else {
					activeIdx = isNext ? activeIdx + 1 : activeIdx - 1;
				}

				if (activeIdx >= 0 && activeIdx < visibleRows.length) {
					e.preventDefault();
					const targetRow = visibleRows[activeIdx];
					targetRow.click();
					targetRow.scrollIntoView({ block: 'nearest', behavior: 'smooth' });
				}
				return;
			}

			if (e.key === 'Home') {
				const visibleRows = Array.from(nodeRows).filter(r => r.style.display !== 'none');
				if (visibleRows.length > 0) {
					e.preventDefault();
					visibleRows[0].click();
					visibleRows[0].scrollIntoView({ block: 'nearest', behavior: 'smooth' });
				}
				return;
			} else if (e.key === 'End') {
				const visibleRows = Array.from(nodeRows).filter(r => r.style.display !== 'none');
				if (visibleRows.length > 0) {
					e.preventDefault();
					const last = visibleRows[visibleRows.length - 1];
					last.click();
					last.scrollIntoView({ block: 'nearest', behavior: 'smooth' });
				}
				return;
			}
		});

		function draw() {
			tickPhysics();

			ctx.save();
			ctx.setTransform(1, 0, 0, 1, 0, 0);
			ctx.clearRect(0, 0, canvas.width, canvas.height);

			ctx.scale(window.devicePixelRatio, window.devicePixelRatio);
			ctx.translate(transform.x, transform.y);
			ctx.scale(transform.scale, transform.scale);

			const isDark = document.documentElement.getAttribute('data-theme') !== 'light';
			const activeFocusNode = hoveredNode || selectedNode;
			const activeNeighbors = activeFocusNode ? (adjacency.get(activeFocusNode.id) || new Set()) : null;

			for (let i = 0; i < links.length; i++) {
				const l = links[i];
				const s = l.source;
				const t = l.target;

				let isConnectedToFocus = false;
				if (activeFocusNode) {
					if ((s.id === activeFocusNode.id && activeNeighbors.has(t.id)) ||
						(t.id === activeFocusNode.id && activeNeighbors.has(s.id))) {
						isConnectedToFocus = true;
					}
				}

				ctx.beginPath();
				ctx.moveTo(s.x, s.y);
				ctx.lineTo(t.x, t.y);

				if (activeFocusNode) {
					if (isConnectedToFocus) {
						ctx.lineWidth = l.kind === 'crosslink' ? 1.8 : 1.4;
						ctx.strokeStyle = activeFocusNode.color || (isDark ? '#38bdf8' : '#0284c7');
					} else {
						ctx.lineWidth = 0.4;
						ctx.strokeStyle = isDark ? 'rgba(255, 255, 255, 0.02)' : 'rgba(0, 0, 0, 0.025)';
					}
				} else {
					if (l.kind === 'crosslink') {
						ctx.lineWidth = 0.9;
						ctx.strokeStyle = isDark ? 'rgba(56, 189, 248, 0.35)' : 'rgba(2, 132, 199, 0.3)';
					} else {
						ctx.lineWidth = 0.6;
						ctx.strokeStyle = isDark ? 'rgba(255, 255, 255, 0.12)' : 'rgba(15, 23, 42, 0.10)';
					}
				}
				ctx.stroke();
			}

			for (let i = 0; i < nodes.length; i++) {
				const n = nodes[i];
				const isTag = n.type === 'tag';
				const isSelected = selectedNode && selectedNode.id === n.id;
				const isHovered = hoveredNode && hoveredNode.id === n.id;
				const isNeighbor = activeFocusNode && activeNeighbors && activeNeighbors.has(n.id);
				const isDimmed = activeFocusNode && !isSelected && !isHovered && !isNeighbor;

				const isMatchSearch = searchQuery && (
					n.label.toLowerCase().includes(searchQuery) ||
					(n.tag && n.tag.toLowerCase().includes(searchQuery))
				);

				const isFilteredOut = (currentFilter === 'posts' && isTag) ||
					(currentFilter === 'domains' && !isTag);

				let r = n.r;
				if (isSelected || isHovered) r += 1.8;

				if (isDimmed || isFilteredOut) {

					ctx.beginPath();
					ctx.arc(n.x, n.y, Math.max(2.2, r - 1.2), 0, Math.PI * 2);
					ctx.fillStyle = isDark ? 'rgba(255, 255, 255, 0.07)' : 'rgba(15, 23, 42, 0.07)';
					ctx.fill();
					continue;
				}

				const baseColor = n.color || (isDark ? '#cbd5e1' : '#64748b');

				if (isSelected || isHovered) {
					const glowRadius = r * 3.5;
					const glowGrad = ctx.createRadialGradient(n.x, n.y, r * 0.4, n.x, n.y, glowRadius);
					glowGrad.addColorStop(0, hexToRgba(baseColor, 0.48));
					glowGrad.addColorStop(0.55, hexToRgba(baseColor, 0.16));
					glowGrad.addColorStop(1, hexToRgba(baseColor, 0));

					ctx.beginPath();
					ctx.arc(n.x, n.y, glowRadius, 0, Math.PI * 2);
					ctx.fillStyle = glowGrad;
					ctx.fill();
				} else if (isNeighbor) {

					const neighborGlow = ctx.createRadialGradient(n.x, n.y, r * 0.5, n.x, n.y, r * 2.2);
					neighborGlow.addColorStop(0, hexToRgba(baseColor, 0.28));
					neighborGlow.addColorStop(1, hexToRgba(baseColor, 0));

					ctx.beginPath();
					ctx.arc(n.x, n.y, r * 2.2, 0, Math.PI * 2);
					ctx.fillStyle = neighborGlow;
					ctx.fill();
				}

				ctx.beginPath();
				ctx.arc(n.x, n.y, r, 0, Math.PI * 2);
				ctx.fillStyle = baseColor;
				ctx.fill();

				if (isSelected) {
					ctx.lineWidth = 1.6;
					ctx.strokeStyle = isDark ? '#ffffff' : '#0f172a';
					ctx.stroke();
				} else if (isHovered) {
					ctx.lineWidth = 1.2;
					ctx.strokeStyle = isDark ? 'rgba(255, 255, 255, 0.8)' : 'rgba(15, 23, 42, 0.8)';
					ctx.stroke();
				}

				if (isMatchSearch) {
					ctx.beginPath();
					ctx.arc(n.x, n.y, r + 4 + Math.sin(performance.now() / 160) * 2.5, 0, Math.PI * 2);
					ctx.lineWidth = 1.6;
					ctx.strokeStyle = '#f59e0b';
					ctx.stroke();
				}

				const isHighPriority = isSelected || isHovered || isNeighbor || isMatchSearch;
				const isZoomedIn = transform.scale >= 0.85;
				const isTagHubVisible = isTag && transform.scale >= 0.65;
				const shouldDrawLabel = isHighPriority || isZoomedIn || isTagHubVisible;

				if (shouldDrawLabel) {
					ctx.save();
					ctx.textAlign = 'center';
					ctx.textBaseline = 'middle';

					ctx.shadowColor = isDark ? 'rgba(0, 0, 0, 0.92)' : 'rgba(255, 255, 255, 0.95)';
					ctx.shadowBlur = 4;
					ctx.shadowOffsetX = 0;
					ctx.shadowOffsetY = 1;

					const labelY = n.y + r + 9;
					if (isTag) {
						ctx.font = '700 10px system-ui, -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif';
						ctx.fillStyle = isSelected || isHovered ? (isDark ? '#ffffff' : '#0f172a') : baseColor;
						ctx.fillText('#' + n.label.toUpperCase(), n.x, labelY);
					} else {
						ctx.font = isSelected || isHovered
							? '600 10.5px system-ui, -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif'
							: '500 10px system-ui, -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif';

						let title = n.label;
						if (title.length > 26 && !isHighPriority) {
							title = title.substring(0, 24) + '…';
						}
						ctx.fillStyle = isHighPriority
							? (isDark ? '#ffffff' : '#0f172a')
							: (isDark ? '#cbd5e1' : '#334155');

						ctx.fillText(title, n.x, labelY);
					}
					ctx.restore();
				}
			}

			ctx.restore();
			requestAnimationFrame(draw);
		}

		resetCamera();
		requestAnimationFrame(draw);
	}

	boot();
})();
