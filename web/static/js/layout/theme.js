// Theme toggle (Dark / Light)
window.toggleTheme = function() {
	var html = document.documentElement;
	var current = html.getAttribute("data-theme");
	var isDark = current === "dark" ||
		(!current && window.matchMedia("(prefers-color-scheme: dark)").matches);
	var next = isDark ? "light" : "dark";
	html.setAttribute("data-theme", next);
	localStorage.setItem("theme", next);
};

window.setTheme = function(t) {
	document.documentElement.setAttribute("data-theme", t);
	localStorage.setItem("theme", t);
};
