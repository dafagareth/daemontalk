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
