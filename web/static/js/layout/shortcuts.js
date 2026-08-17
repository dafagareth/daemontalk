// Global keyboard shortcuts (?, /, ~, j/k, g+h/b/p/t/l/u)
window.openShortcuts = function() {
	var el = document.getElementById("shortcuts-overlay");
	if (el) el.classList.add("open");
};

(function() {
	var pendingG = false;
	var gTimer = null;

	function isInputFocused() {
		var el = document.activeElement;
		if (!el) return false;
		var tag = el.tagName;
		return tag === 'INPUT' || tag === 'TEXTAREA' || tag === 'SELECT' || el.isContentEditable;
	}

	document.addEventListener('keydown', function(e) {
		if (e.metaKey || e.ctrlKey || e.altKey) return;

		if (e.key === 'Escape') {
			var shortcuts = document.getElementById('shortcuts-overlay');
			if (shortcuts) shortcuts.classList.remove('open');
			if (document.activeElement) document.activeElement.blur();
			pendingG = false;
			clearTimeout(gTimer);
			return;
		}

		if (isInputFocused()) return;

		if (e.key === '?') {
			e.preventDefault();
			var shortcuts = document.getElementById('shortcuts-overlay');
			if (shortcuts) shortcuts.classList.toggle('open');
			pendingG = false;
			return;
		}

		if (e.key === '/') {
			e.preventDefault();
			pendingG = false;
			var si = document.getElementById('blog-search');
			if (si) { si.focus(); si.select(); }
			else window.location.href = '/search';
			return;
		}

		if (e.key === '`' || e.key === '~') {
			e.preventDefault();
			pendingG = false;
			window.location.href = '/terminal';
			return;
		}

		if (pendingG) {
			pendingG = false;
			clearTimeout(gTimer);
			e.preventDefault();
			var destinations = { h: '/', b: '/blog', p: '/projects', u: '/uses', t: '/terminal', l: '/til' };
			if (destinations[e.key]) window.location.href = destinations[e.key];
			return;
		}

		if (e.key === 'g') {
			e.preventDefault();
			pendingG = true;
			gTimer = setTimeout(function() { pendingG = false; }, 1000);
			return;
		}

		if (e.key === 'j' || e.key === 'k') {
			var items = document.querySelectorAll('#blog-list .blog-item');
			if (!items.length) return;
			e.preventDefault();
			var current = document.querySelector('#blog-list .blog-item.kb-focus');
			var idx = -1;
			for (var i = 0; i < items.length; i++) {
				if (items[i] === current) { idx = i; break; }
			}
			if (e.key === 'j') idx = Math.min(idx + 1, items.length - 1);
			else idx = Math.max(idx - 1, 0);
			for (var j = 0; j < items.length; j++) items[j].classList.remove('kb-focus');
			items[idx].classList.add('kb-focus');
			items[idx].scrollIntoView({ block: 'nearest', behavior: 'smooth' });
			return;
		}

		if (e.key === 'Enter') {
			var focused = document.querySelector('#blog-list .blog-item.kb-focus');
			if (focused) {
				var link = focused.querySelector('a');
				if (link) { e.preventDefault(); window.location.href = link.href; }
			}
		}
	});
})();
