// Admin Panel Management: Tab routing, Theme toggle, Digest copy
function toggleTheme() {
	var html = document.documentElement;
	var current = html.getAttribute("data-theme");
	var isDark = current === "dark" ||
		(!current && window.matchMedia("(prefers-color-scheme: dark)").matches);
	var next = isDark ? "light" : "dark";
	html.setAttribute("data-theme", next);
	localStorage.setItem("theme", next);
}

function switchAdminTab(tabName, event) {
	if (event) {
		event.preventDefault();
	}
	var tabs = ["dashboard", "content", "analytics", "comments", "digest"];
	var targetEl = document.getElementById("tab-" + tabName);
	if (!targetEl) {
		if (document.getElementById("admin-forbidden")) {
			return;
		}
		window.location.href = tabName === "dashboard" ? "/admin" : "/admin/" + tabName;
		return;
	}

	tabs.forEach(function(t) {
		var el = document.getElementById("tab-" + t);
		if (el) {
			if (t === tabName) {
				el.classList.remove("hidden");
			} else {
				el.classList.add("hidden");
			}
		}
	});

	var navBtns = document.querySelectorAll("[data-tab-btn]");
	navBtns.forEach(function(btn) {
		var btnTab = btn.getAttribute("data-tab-btn");
		if (btnTab === tabName) {
			btn.classList.add("active");
		} else {
			btn.classList.remove("active");
		}
	});

	var targetPath = tabName === "dashboard" ? "/admin" : "/admin/" + tabName;
	if (window.location.pathname !== targetPath) {
		history.pushState({ tab: tabName }, "", targetPath);
	}
}

function initAdminTabFromURL() {
	var path = window.location.pathname.replace(/\/$/, "");
	if (path.indexOf("/admin/posts") === 0 || path.indexOf("/admin/post") === 0) {
		return;
	}
	var tabName = "dashboard";
	if (path === "/admin/content") tabName = "content";
	else if (path === "/admin/analytics") tabName = "analytics";
	else if (path === "/admin/comments") tabName = "comments";
	else if (path === "/admin/digest") tabName = "digest";
	else {
		var hash = window.location.hash.replace("#", "");
		if (hash && ["dashboard", "content", "analytics", "comments", "digest"].indexOf(hash) !== -1) {
			tabName = hash;
		}
	}
	switchAdminTab(tabName);
}

function copyDigestMarkdown() {
	var copyText = document.getElementById("digest-markdown");
	if (!copyText) return;
	copyText.select();
	navigator.clipboard.writeText(copyText.value);
	var btn = document.getElementById("copy-digest-btn");
	if (btn) {
		var old = btn.innerText;
		btn.innerText = "✓ Copied!";
		setTimeout(function() { btn.innerText = old; }, 2000);
	}
}

window.toggleTheme = toggleTheme;
window.switchAdminTab = switchAdminTab;
window.copyDigestMarkdown = copyDigestMarkdown;

window.addEventListener("popstate", function() {
	initAdminTabFromURL();
});

document.addEventListener("DOMContentLoaded", function() {
	initAdminTabFromURL();
});
