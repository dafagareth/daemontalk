// Custom confirmation modal (replaces window.confirm)
(function() {
	var overlay = document.getElementById('confirm-overlay');
	var msgEl = document.getElementById('confirm-message');
	if (!overlay || !msgEl) return;
	var pendingResolve = null;

	var okBtn = overlay.querySelector('[data-confirm-ok]');
	var cancelBtn = overlay.querySelector('[data-confirm-cancel]');

	function humanConfirm(message) {
		return new Promise(function(resolve) {
			msgEl.textContent = message;
			var isDelete = /hapus|delete|remove|destroy/i.test(message);
			if (okBtn) {
				if (isDelete) {
					okBtn.textContent = 'Hapus';
					okBtn.className = 'text-sm px-4 py-2 rounded-none bg-[var(--c-warn)] text-white font-medium hover:brightness-110 transition-all cursor-pointer';
				} else {
					okBtn.textContent = 'OK';
					okBtn.className = 'text-sm px-4 py-2 rounded-none bg-[var(--c-link)] text-white font-medium hover:brightness-110 transition-all cursor-pointer';
				}
			}
			overlay.classList.add('open');
			pendingResolve = resolve;
		});
	}
	function close(result) {
		overlay.classList.remove('open');
		if (pendingResolve) { pendingResolve(result); pendingResolve = null; }
	}
	if (okBtn) okBtn.addEventListener('click', function() { close(true); });
	if (cancelBtn) cancelBtn.addEventListener('click', function() { close(false); });
	document.addEventListener('keydown', function(e) {
		if (e.key === 'Escape' && overlay.classList.contains('open')) close(false);
	});

	// HTMX hx-confirm -> custom modal
	document.body.addEventListener('htmx:confirm', function(e) {
		if (!e.detail.question) return;
		e.preventDefault();
		humanConfirm(e.detail.question).then(function(ok) {
			if (ok) e.detail.issueRequest(true);
		});
	});

	// Standard form data-confirm attribute
	document.addEventListener('submit', function(e) {
		var form = e.target;
		if (!(form instanceof HTMLFormElement)) return;
		var msg = form.getAttribute('data-confirm');
		if (!msg || form.dataset.confirmed === 'true') return;
		e.preventDefault();
		humanConfirm(msg).then(function(ok) {
			if (ok) { form.dataset.confirmed = 'true'; form.submit(); }
		});
	});
})();
