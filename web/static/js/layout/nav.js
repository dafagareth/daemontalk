window.toggleMobileMenu = function() {
	var menu = document.getElementById("mobile-menu-overlay");
	if (!menu) return;
	if (menu.classList.contains("-translate-x-full")) {
		menu.classList.remove("-translate-x-full");
		menu.classList.add("translate-x-0");
		document.body.style.overflow = "hidden";
	} else {
		menu.classList.remove("translate-x-0");
		menu.classList.add("-translate-x-full");
		document.body.style.overflow = "";
	}
};

(function() {
	var params = new URLSearchParams(window.location.search);
	var tag = params.get("tag") || "";
	var isTagPage = false;
	if (!tag) {
		var match = window.location.pathname.match(/\/blog\/tag\/([^\/]+)/);
		if (match) {
			tag = decodeURIComponent(match[1]).toLowerCase();
			isTagPage = true;
		}
	} else {
		isTagPage = true;
	}
	var isHome = window.location.pathname === "/" || window.location.pathname === "/id" || window.location.pathname === "/id/";
	var links = document.querySelectorAll("#subnav-tags .subnav-tag");
	for (var i = 0; i < links.length; i++) {
		var linkTag = links[i].getAttribute("data-tag");
		if (isTagPage && linkTag === tag) {
			links[i].classList.add("active");
		} else if (isHome && linkTag === "") {
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
	var header = document.getElementById("site-header-wrapper");
	if (!header) return;

	var lastScrollY = window.scrollY;
	var threshold = 10;
	var topThreshold = 70;

	window.addEventListener("scroll", function() {

		var isBlogPost = !!document.getElementById("reading-progress-bar") || !!document.getElementById("nav-post-title");
		if (isBlogPost) {
			header.classList.remove("-translate-y-full");
			return;
		}

		var currentScrollY = window.scrollY;

		if (currentScrollY <= topThreshold) {

			header.classList.remove("-translate-y-full");
			lastScrollY = currentScrollY;
			return;
		}

		var diff = currentScrollY - lastScrollY;

		if (Math.abs(diff) < threshold) {
			return;
		}

		if (diff > 0 && currentScrollY > topThreshold) {

			header.classList.add("-translate-y-full");
		} else if (diff < 0) {

			header.classList.remove("-translate-y-full");
		}

		lastScrollY = currentScrollY;
	}, { passive: true });
})();

(function() {
	var navTitle = document.getElementById("nav-post-title");
	var progressBar = document.getElementById("reading-progress-bar");
	if (!navTitle && !progressBar) return;

	var articleHeader = document.querySelector("article header h1") || document.querySelector("article h1") || document.querySelector("h1");
	var blogContent = document.getElementById("prose-body") || document.getElementById("blog-content") || document.querySelector("article");

	function updateTitleAndProgress() {

		if (navTitle && articleHeader) {
			var rect = articleHeader.getBoundingClientRect();
			if (rect.bottom < 56) {
				navTitle.classList.remove("opacity-0", "-translate-y-1", "pointer-events-none");
				navTitle.classList.add("opacity-100", "translate-y-0", "pointer-events-auto");
			} else {
				navTitle.classList.add("opacity-0", "-translate-y-1", "pointer-events-none");
				navTitle.classList.remove("opacity-100", "translate-y-0", "pointer-events-auto");
			}
		}

		if (progressBar && blogContent) {
			var rect = blogContent.getBoundingClientRect();
			var contentTop = rect.top + window.scrollY;
			var contentHeight = blogContent.offsetHeight;
			var windowHeight = window.innerHeight;

			var startY = 0;
			var endY = contentTop + contentHeight - (windowHeight * 0.75);

			var progress = 0;
			if (endY > startY) {
				if (window.scrollY >= endY) {
					progress = 100;
				} else if (window.scrollY > startY) {
					progress = ((window.scrollY - startY) / (endY - startY)) * 100;
				}
			} else {
				progress = window.scrollY > 0 ? 100 : 0;
			}
			progressBar.style.width = Math.min(Math.max(progress, 0), 100).toFixed(1) + "%";
		}
	}

	window.addEventListener("scroll", updateTitleAndProgress, { passive: true });
	window.addEventListener("resize", updateTitleAndProgress, { passive: true });
	if (document.readyState === "loading") {
		document.addEventListener("DOMContentLoaded", updateTitleAndProgress);
	} else {
		updateTitleAndProgress();
	}
})();

window.openSearchPopper = function() {
	var popper = document.getElementById("search-dropdown-popper");
	var input = document.getElementById("header-search-input");
	if (popper && input && input.value.trim().length > 0) {
		popper.classList.remove("hidden");
	}
};

window.closeSearchPopper = function() {
	var popper = document.getElementById("search-dropdown-popper");
	if (popper) {
		popper.classList.add("hidden");
	}
};

document.addEventListener("htmx:afterSwap", function(evt) {
	if (evt.detail && evt.detail.target && evt.detail.target.id === "search-dropdown-results") {
		var popper = document.getElementById("search-dropdown-popper");
		var input = document.getElementById("header-search-input");
		if (popper) {
			if (input && input.value.trim().length > 0 && evt.detail.target.innerHTML.trim().length > 0) {
				popper.classList.remove("hidden");
			} else {
				popper.classList.add("hidden");
			}
		}
	}
});

document.addEventListener("input", function(e) {
	if (e.target && e.target.id === "header-search-input") {
		if (e.target.value.trim().length === 0) {
			window.closeSearchPopper();
		}
	}
});

window.handleSearchKeyNav = function(e) {
	var popper = document.getElementById("search-dropdown-popper");
	if (!popper || popper.classList.contains("hidden")) return;

	var items = popper.querySelectorAll(".search-popper-item");
	if (items.length === 0) return;

	var activeIndex = -1;
	for (var i = 0; i < items.length; i++) {
		if (document.activeElement === items[i]) {
			activeIndex = i;
			break;
		}
	}

	if (e.key === "ArrowDown") {
		e.preventDefault();
		var nextIndex = (activeIndex + 1) % items.length;
		items[nextIndex].focus();
	} else if (e.key === "ArrowUp") {
		e.preventDefault();
		if (activeIndex <= 0) {
			var input = document.getElementById("header-search-input");
			if (input) input.focus();
		} else {
			items[activeIndex - 1].focus();
		}
	} else if (e.key === "Escape") {
		e.preventDefault();
		window.closeSearchPopper();
		var input = document.getElementById("header-search-input");
		if (input) input.blur();
	}
};

window.handleNavSearchSubmit = function(e) {
	var input = document.getElementById("header-search-input");
	if (!input || !input.value.trim()) {
		e.preventDefault();
	}
};

document.addEventListener("click", function(e) {
	var container = document.getElementById("nav-search-container");
	if (container && !container.contains(e.target)) {
		window.closeSearchPopper();
	}
});

window.clearSearchInput = function() {
	var input = document.getElementById('search-input');
	if (input) {
		input.value = '';
		input.focus();
		if (typeof htmx !== 'undefined') {
			htmx.trigger(input, 'search');
		}
	}
};

document.addEventListener('click', function(e) {
	var userContainer = document.getElementById('user-menu-container');
	var userDropdown = document.getElementById('user-menu-dropdown');
	if (userDropdown && userContainer && !userContainer.contains(e.target)) {
		userDropdown.classList.add('hidden');
	}

	var guestContainer = document.getElementById('guest-menu-container');
	var guestDropdown = document.getElementById('guest-menu-dropdown');
	if (guestDropdown && guestContainer && !guestContainer.contains(e.target)) {
		guestDropdown.classList.add('hidden');
	}
});
