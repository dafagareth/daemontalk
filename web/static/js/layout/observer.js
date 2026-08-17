// Scroll reveal observer & comment datetime formatting
(function() {
	var els = document.querySelectorAll('.reveal');
	if (!els.length) return;
	var io = new IntersectionObserver(function(entries) {
		entries.forEach(function(e) {
			if (e.isIntersecting) {
				e.target.classList.add('revealed');
				io.unobserve(e.target);
			}
		});
	}, { threshold: 0.08 });
	els.forEach(function(el) { io.observe(el); });
})();

(function() {
	var months = ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec'];
	function formatTimes() {
		document.querySelectorAll('time.comment-time').forEach(function(el) {
			var iso = el.getAttribute('datetime');
			if (!iso) return;
			var d = new Date(iso);
			if (isNaN(d.getTime())) return;
			var day = d.getDate();
			var month = months[d.getMonth()];
			var year = d.getFullYear();
			var hours = String(d.getHours()).padStart(2, '0');
			var minutes = String(d.getMinutes()).padStart(2, '0');
			el.textContent = day + ' ' + month + ' ' + year + ', ' + hours + ':' + minutes;
		});
	}
	if (document.readyState === 'loading') {
		document.addEventListener('DOMContentLoaded', formatTimes);
	} else {
		formatTimes();
	}
	document.addEventListener('htmx:afterSwap', formatTimes);
	document.addEventListener('htmx:afterRequest', function(e) {
		if (e.detail && e.detail.successful && e.target && e.target.tagName === 'FORM') {
			e.target.reset();
			var ta = e.target.querySelector('textarea');
			if (ta) {
				ta.value = '';
			}
			if (window.cancelCommentReply) {
				window.cancelCommentReply();
			}
		}
	});

	function setCommentReply(id, author) {
		var parentInput = document.getElementById('comment-parent-id');
		var indicator = document.getElementById('reply-indicator');
		var authorEl = document.getElementById('reply-author-name');
		var cancelBtn = document.getElementById('reply-cancel-btn');
		var textarea = document.getElementById('comment-textarea');
		var formContainer = document.getElementById('comment-form-container');

		if (parentInput) parentInput.value = id;
		if (authorEl) authorEl.textContent = '@' + author;
		if (indicator) {
			indicator.classList.remove('hidden');
			indicator.classList.add('flex');
		}
		if (cancelBtn) {
			cancelBtn.classList.remove('hidden');
		}
		if (textarea) {
			if (!textarea.getAttribute('data-original-placeholder')) {
				textarea.setAttribute('data-original-placeholder', textarea.placeholder);
			}
			textarea.placeholder = 'Reply to @' + author + '...';
		}
		if (formContainer) {
			formContainer.scrollIntoView({ behavior: 'smooth', block: 'center' });
		}
		if (textarea) {
			setTimeout(function() {
				textarea.focus();
			}, 150);
		}
	}

	function cancelCommentReply() {
		var parentInput = document.getElementById('comment-parent-id');
		var indicator = document.getElementById('reply-indicator');
		var cancelBtn = document.getElementById('reply-cancel-btn');
		var textarea = document.getElementById('comment-textarea');

		if (parentInput) parentInput.value = '';
		if (indicator) {
			indicator.classList.add('hidden');
			indicator.classList.remove('flex');
		}
		if (cancelBtn) {
			cancelBtn.classList.add('hidden');
		}
		if (textarea) {
			var orig = textarea.getAttribute('data-original-placeholder');
			if (orig) textarea.placeholder = orig;
		}
	}

	window.setCommentReply = setCommentReply;
	window.cancelCommentReply = cancelCommentReply;

	// Delegated touch and click handler for instant responsiveness on mobile
	document.addEventListener('click', function(e) {
		var replyBtn = e.target.closest('[data-reply-btn]');
		if (replyBtn) {
			e.preventDefault();
			var id = replyBtn.getAttribute('data-reply-id');
			var author = replyBtn.getAttribute('data-reply-author');
			setCommentReply(id, author);
			return;
		}
		var cancelBtn = e.target.closest('[data-reply-cancel-btn]');
		if (cancelBtn) {
			e.preventDefault();
			cancelCommentReply();
			return;
		}
	});
})();
