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

// Live Time in Subnav
(function() {
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
})();

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
