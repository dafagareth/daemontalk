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

// Visited / Read post styling marker
(function() {
	function markVisited() {
		var read = [];
		try { read = JSON.parse(localStorage.getItem('readPosts') || '[]'); } catch(e) {}
		if (!read.length) return;
		var set = {};
		read.forEach(function(s) { set[s] = true; });

		document.querySelectorAll('a[data-slug], a[href*="/blog/"]').forEach(function(a) {
			var slug = a.dataset.slug;
			if (!slug) {
				var href = a.getAttribute('href') || '';
				var match = href.match(/\/blog\/([^\/\?#]+)$/);
				if (match) slug = match[1];
			}
			if (slug && set[slug]) {
				var h = a.querySelector('h1,h2,h3,h4,h5,h6');
				if (h) {
					h.classList.add('post-visited');
				} else {
					a.classList.add('post-visited');
				}
			}
		});
	}

	if (document.readyState === 'loading') {
		document.addEventListener('DOMContentLoaded', markVisited);
	} else {
		markVisited();
	}
	document.addEventListener('htmx:afterSwap', markVisited);
})();

(function() {
	var monthsID = ['Januari', 'Februari', 'Maret', 'April', 'Mei', 'Juni', 'Juli', 'Agustus', 'September', 'Oktober', 'November', 'Desember'];
	var monthsEN = ['January', 'February', 'March', 'April', 'May', 'June', 'July', 'August', 'September', 'October', 'November', 'December'];

	function formatCommentDate(d, lang) {
		var now = new Date();
		var diffMs = now.getTime() - d.getTime();
		if (diffMs < 0) diffMs = 0;
		var mins = Math.floor(diffMs / 60000);
		var isID = (lang === 'id');

		if (mins < 1) {
			return isID ? 'baru saja' : 'just now';
		}
		if (mins < 60) {
			if (isID) {
				return mins + ' menit yang lalu';
			}
			return mins === 1 ? '1 minute ago' : mins + ' minutes ago';
		}
		var hours = Math.floor(diffMs / 3600000);
		if (hours < 24) {
			if (isID) {
				return hours + ' jam yang lalu';
			}
			return hours === 1 ? '1 hour ago' : hours + ' hours ago';
		}
		var days = Math.floor(diffMs / 86400000);
		if (days <= 7) {
			if (isID) {
				return days + ' hari yang lalu';
			}
			return days === 1 ? '1 day ago' : days + ' days ago';
		}
		var months = isID ? monthsID : monthsEN;
		var day = d.getDate();
		var monthName = months[d.getMonth()];
		if (d.getFullYear() === now.getFullYear()) {
			return day + ' ' + monthName;
		}
		return day + ' ' + monthName + ' ' + d.getFullYear();
	}

	function formatTimes() {
		var lang = document.documentElement.lang || 'en';
		document.querySelectorAll('time.comment-time').forEach(function(el) {
			var iso = el.getAttribute('datetime');
			if (!iso) return;
			var d = new Date(iso);
			if (isNaN(d.getTime())) return;
			el.textContent = formatCommentDate(d, lang);
		});
	}

	// Close comment-menu dropdowns when clicking or touching outside
	document.addEventListener('click', function(e) {
		document.querySelectorAll('details.comment-menu[open]').forEach(function(el) {
			if (!el.contains(e.target)) {
				el.removeAttribute('open');
			}
		});
	});
	document.addEventListener('touchstart', function(e) {
		document.querySelectorAll('details.comment-menu[open]').forEach(function(el) {
			if (!el.contains(e.target)) {
				el.removeAttribute('open');
			}
		});
	}, { passive: true });
	document.addEventListener('toggle', function(e) {
		if (e.target && e.target.open && e.target.classList.contains('comment-menu')) {
			document.querySelectorAll('details.comment-menu[open]').forEach(function(el) {
				if (el !== e.target) {
					el.removeAttribute('open');
				}
			});
		}
	}, true);

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
