function filterPosts(q) {
	q = (q || "").trim().toLowerCase();
	var items = document.querySelectorAll(".blog-item");
	var shown = 0;
	items.forEach(function(el) {
		var hay = (el.getAttribute("data-search") || "").toLowerCase();
		var match = q === "" || hay.indexOf(q) !== -1;
		el.style.display = match ? "" : "none";
		if (match) shown++;
	});
	var none = document.getElementById("blog-noresults");
	if (none) none.classList.toggle("hidden", shown !== 0);
}

function setTagActive(btn) {
	btn.parentElement.querySelectorAll('a').forEach(function(el) {
		el.className = 'text-xs font-mono px-3 py-1.5 rounded-none transition-all bg-chip hover:bg-hover text-muted hover:text-text';
	});
	btn.className = 'text-xs font-mono px-3 py-1.5 rounded-none transition-all bg-[var(--c-link)] text-white font-semibold shadow-sm';
}

function initLoadMore(btn) {
	// HTMX handles swap
}

// Bookmark state: set ★/☆ on all bookmark buttons based on localStorage.
function syncBookmarks() {
	var bookmarks = [];
	try { bookmarks = JSON.parse(localStorage.getItem("bookmarks") || "[]"); } catch(e) {}
	var saved = {};
	bookmarks.forEach(function(b) { saved[b.slug] = true; });
	document.querySelectorAll(".bookmark-btn").forEach(function(btn) {
		if (saved[btn.dataset.slug]) { btn.textContent = "★"; btn.style.color = "#f59e0b"; }
	});
}
syncBookmarks();
document.addEventListener("htmx:afterSwap", syncBookmarks);

// Live Masthead Clock
(function() {
	var clock = document.getElementById('masthead-clock');
	if (clock) {
		function updateTime() {
			var now = new Date();
			var d = now.toLocaleDateString('id-ID', { day: '2-digit', month: 'short', year: 'numeric' }).toUpperCase();
			var t = now.toTimeString().split(' ')[0];
			clock.textContent = d + ' // ' + t + ' WIB';
		}
		updateTime();
		setInterval(updateTime, 1000);
	}
})();

// Live Dispatch Ticker Slider
(function() {
	var cur = 0;
	var total = 5;
	var timer = null;
	function show(idx) {
		for (var i = 0; i < total; i++) {
			var el = document.getElementById('ticker-item-' + i);
			if (!el) continue;
			if (i === idx) {
				el.className = 'ticker-item absolute inset-0 flex items-center gap-2 text-xs text-text transition-all duration-300 opacity-100 transform translate-y-0';
			} else if (i < idx) {
				el.className = 'ticker-item absolute inset-0 flex items-center gap-2 text-xs text-text transition-all duration-300 opacity-0 transform -translate-y-4';
			} else {
				el.className = 'ticker-item absolute inset-0 flex items-center gap-2 text-xs text-text transition-all duration-300 opacity-0 transform translate-y-4';
			}
		}
		cur = idx;
	}
	window.nextTicker = function() {
		show((cur + 1) % total);
	};
	window.prevTicker = function() {
		show((cur - 1 + total) % total);
	};
	function startTimer() {
		if (timer) clearInterval(timer);
		timer = setInterval(function() {
			window.nextTicker();
		}, 4500);
	}
	var slider = document.getElementById('ticker-slider');
	if (slider) {
		startTimer();
		slider.parentElement.addEventListener('mouseenter', function() { if (timer) clearInterval(timer); });
		slider.parentElement.addEventListener('mouseleave', function() { startTimer(); });
	}
})();

// Portal Hero Slider Controller
(function() {
	var curHero = 0;
	var heroTimer = null;

	window.setHeroSlide = function(idx) {
		var slides = document.querySelectorAll('.hero-slide');
		if (!slides.length) return;
		for (var i = 0; i < slides.length; i++) {
			var s = document.getElementById('hero-slide-' + i);
			var d = document.getElementById('hero-dot-' + i);
			if (s) {
				if (i === idx) {
					s.classList.remove('opacity-0', 'pointer-events-none', 'z-0');
					s.classList.add('opacity-100', 'z-10');
				} else {
					s.classList.remove('opacity-100', 'z-10');
					s.classList.add('opacity-0', 'pointer-events-none', 'z-0');
				}
			}
			if (d) {
				if (i === idx) {
					d.className = 'w-2 h-2 rounded-none transition-all cursor-pointer bg-white scale-125';
				} else {
					d.className = 'w-2 h-2 rounded-none transition-all cursor-pointer bg-white/40 hover:bg-white/70';
				}
			}
		}
		curHero = idx;
	};

	function startHeroTimer() {
		if (heroTimer) clearInterval(heroTimer);
		heroTimer = setInterval(function() {
			var slides = document.querySelectorAll('.hero-slide');
			if (slides.length > 1) {
				window.setHeroSlide((curHero + 1) % slides.length);
			}
		}, 5500);
	}

	var heroContainer = document.getElementById('hero-slider-container');
	if (heroContainer) {
		startHeroTimer();
		heroContainer.addEventListener('mouseenter', function() { if (heroTimer) clearInterval(heroTimer); });
		heroContainer.addEventListener('mouseleave', function() { startHeroTimer(); });

		// Touch swipe support for mobile
		var touchStartX = 0;
		var touchStartY = 0;
		heroContainer.addEventListener('touchstart', function(e) {
			if (e.touches.length === 1) {
				touchStartX = e.touches[0].clientX;
				touchStartY = e.touches[0].clientY;
			}
		}, { passive: true });

		heroContainer.addEventListener('touchend', function(e) {
			if (e.changedTouches.length === 1) {
				var diffX = e.changedTouches[0].clientX - touchStartX;
				var diffY = e.changedTouches[0].clientY - touchStartY;
				if (Math.abs(diffX) > 40 && Math.abs(diffX) > Math.abs(diffY)) {
					var slides = document.querySelectorAll('.hero-slide');
					if (slides.length > 1) {
						if (diffX < 0) {
							window.setHeroSlide((curHero + 1) % slides.length);
						} else {
							window.setHeroSlide((curHero - 1 + slides.length) % slides.length);
						}
						startHeroTimer();
					}
				}
			}
		}, { passive: true });
	}
})();

window.toggleBookmark = function(btn) {
	var slug = btn.dataset.slug;
	var title = btn.dataset.title || "";
	var date = btn.dataset.date || "";
	var bookmarks = [];
	try { bookmarks = JSON.parse(localStorage.getItem("bookmarks") || "[]"); } catch(e) {}
	var idx = bookmarks.findIndex(function(b) { return b.slug === slug; });
	if (idx >= 0) {
		bookmarks.splice(idx, 1);
		btn.textContent = "☆";
		btn.style.color = "";
	} else {
		bookmarks.unshift({slug: slug, title: title, date: date});
		btn.textContent = "★";
		btn.style.color = "#f59e0b";
	}
	localStorage.setItem("bookmarks", JSON.stringify(bookmarks));
};
