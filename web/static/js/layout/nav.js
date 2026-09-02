// Navigation helpers: mobile menu, active subnav, back-to-top, live clock
window.toggleMobileMenu = function() {
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
};


// Highlight active sub-nav category tag based on URL
(function() {
	var params = new URLSearchParams(window.location.search);
	var tag = params.get("tag") || "";
	if (!tag) {
		var match = window.location.pathname.match(/\/blog\/tag\/([^\/]+)/);
		if (match) {
			tag = decodeURIComponent(match[1]).toLowerCase();
		}
	}
	var links = document.querySelectorAll("#subnav-tags .subnav-tag");
	for (var i = 0; i < links.length; i++) {
		if (links[i].getAttribute("data-tag") === tag) {
			links[i].classList.add("active");
		}
	}
})();

// Back to Top Button
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

// Smart Auto-Hide Header on Scroll Down, Reveal on Scroll Up (Header + Subheader together)
(function() {
	var header = document.getElementById("site-header-wrapper");
	if (!header) return;

	var lastScrollY = window.scrollY;
	var threshold = 10;
	var topThreshold = 70;

	window.addEventListener("scroll", function() {
		var currentScrollY = window.scrollY;
		
		if (currentScrollY <= topThreshold) {
			// At the top of the page: always visible
			header.classList.remove("-translate-y-full");
			lastScrollY = currentScrollY;
			return;
		}

		var diff = currentScrollY - lastScrollY;

		if (Math.abs(diff) < threshold) {
			return;
		}

		if (diff > 0 && currentScrollY > topThreshold) {
			// Scrolling down -> hide both header & subheader
			header.classList.add("-translate-y-full");
		} else if (diff < 0) {
			// Scrolling up -> reveal both header & subheader
			header.classList.remove("-translate-y-full");
		}

		lastScrollY = currentScrollY;
	}, { passive: true });
})();

// Instant Search Dropdown Popper Logic
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

// Monitor HTMX content swap inside #search-dropdown-results
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

// Hide popper when input is cleared
document.addEventListener("input", function(e) {
	if (e.target && e.target.id === "header-search-input") {
		if (e.target.value.trim().length === 0) {
			window.closeSearchPopper();
		}
	}
});

// Keyboard navigation in search dropdown
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

// Form submit handler
window.handleNavSearchSubmit = function(e) {
	var input = document.getElementById("header-search-input");
	if (!input || !input.value.trim()) {
		e.preventDefault();
	}
};

// Click outside listener to dismiss search dropdown
document.addEventListener("click", function(e) {
	var container = document.getElementById("nav-search-container");
	if (container && !container.contains(e.target)) {
		window.closeSearchPopper();
	}
});

// Clear search input helper (search page)
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

// User & Guest menu dismiss on outside click
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
