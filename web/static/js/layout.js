(function() {
			var overlay = document.getElementById('confirm-overlay');
			var msgEl = document.getElementById('confirm-message');
			var pendingResolve = null;

			function humanConfirm(message) {
				return new Promise(function(resolve) {
					msgEl.textContent = message;
					overlay.classList.add('open');
					pendingResolve = resolve;
				});
			}
			function close(result) {
				overlay.classList.remove('open');
				if (pendingResolve) { pendingResolve(result); pendingResolve = null; }
			}
			overlay.querySelector('[data-confirm-ok]').addEventListener('click', function() { close(true); });
			overlay.querySelector('[data-confirm-cancel]').addEventListener('click', function() { close(false); });
			document.addEventListener('keydown', function(e) {
				if (e.key === 'Escape' && overlay.classList.contains('open')) close(false);
			});

			// Tombol hx-confirm (htmx) -> modal custom, bukan browser confirm().
			document.body.addEventListener('htmx:confirm', function(e) {
				if (!e.detail.question) return;
				e.preventDefault();
				humanConfirm(e.detail.question).then(function(ok) {
					if (ok) e.detail.issueRequest(true);
				});
			});

			// Form biasa: atribut data-confirm="pesan" pengganti onsubmit="return confirm(...)".
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

function openShortcuts() {
			var el = document.getElementById("shortcuts-overlay");
			if (el) el.classList.add("open");
		}

		function toggleTheme() {
			var html = document.documentElement;
			var current = html.getAttribute("data-theme");
			var isDark = current === "dark" ||
				(!current && window.matchMedia("(prefers-color-scheme: dark)").matches);
			var next = isDark ? "light" : "dark";
			html.setAttribute("data-theme", next);
			localStorage.setItem("theme", next);
		}

		function setTheme(t) {
			document.documentElement.setAttribute("data-theme", t);
			localStorage.setItem("theme", t);
		}

		function toggleMobileMenu() {
			var menu = document.getElementById("mobile-menu-overlay");
			if (!menu) return;
			if (menu.classList.contains("translate-x-full")) {
				menu.classList.remove("translate-x-full");
				menu.classList.add("translate-x-0");
				document.body.style.overflow = "hidden";
			} else {
				menu.classList.remove("translate-x-0");
				menu.classList.add("translate-x-full");
				document.body.style.overflow = "";
			}
		}

		function updateLiveTime() {
			try {
				var now = new Date();
				var parts = new Intl.DateTimeFormat(undefined, {
					weekday: 'short',
					day: 'numeric',
					month: 'short',
					year: 'numeric',
					hour: '2-digit',
					minute: '2-digit',
					second: '2-digit',
					hour12: false
				}).formatToParts(now);
				var weekday = "", day = "", month = "", year = "", hour = "", minute = "", second = "";
				for (var i = 0; i < parts.length; i++) {
					if (parts[i].type === "weekday") weekday = parts[i].value;
					if (parts[i].type === "day") day = parts[i].value;
					if (parts[i].type === "month") month = parts[i].value;
					if (parts[i].type === "year") year = parts[i].value;
					if (parts[i].type === "hour") hour = parts[i].value;
					if (parts[i].type === "minute") minute = parts[i].value;
					if (parts[i].type === "second") second = parts[i].value;
				}
				var str = (weekday ? weekday + ", " : "") + day + " " + month + " " + year + " " + hour + ":" + minute + ":" + second;
				var els = document.querySelectorAll(".user-live-time");
				for (var j = 0; j < els.length; j++) {
					els[j].textContent = str;
				}
			} catch(e) {
				var d = new Date();
				var pad = function(n) { return String(n).padStart(2, '0'); };
				var days = ['Sun','Mon','Tue','Wed','Thu','Fri','Sat'];
				var months = ['Jan','Feb','Mar','Apr','May','Jun','Jul','Aug','Sep','Oct','Nov','Dec'];
				var strFallback = days[d.getDay()] + ", " + d.getDate() + " " + months[d.getMonth()] + " " + d.getFullYear() + " " + pad(d.getHours()) + ":" + pad(d.getMinutes()) + ":" + pad(d.getSeconds());
				var elsFallback = document.querySelectorAll(".user-live-time");
				for (var k = 0; k < elsFallback.length; k++) {
					elsFallback[k].textContent = strFallback;
				}
			}
		}
		updateLiveTime();
		setInterval(updateLiveTime, 1000);

		// Highlight active sub-nav category tag based on URL
		(function() {
			var params = new URLSearchParams(window.location.search);
			var tag = params.get("tag") || "";
			var links = document.querySelectorAll("#subnav-tags .subnav-tag");
			for (var i = 0; i < links.length; i++) {
				if (links[i].getAttribute("data-tag") === tag) {
					links[i].classList.add("active");
				}
			}
		})();

		window.addEventListener("scroll", function() {
			var btt = document.getElementById("back-to-top");
			if (btt) {
				if (window.scrollY > 250) {
					btt.classList.remove("opacity-0", "pointer-events-none");
					btt.classList.add("opacity-100");
				} else {
					btt.classList.add("opacity-0", "pointer-events-none");
					btt.classList.remove("opacity-100");
				}
			}
		}, { passive: true });

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
		}
	});
})();