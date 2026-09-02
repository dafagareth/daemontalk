(function() {
	var list = document.getElementById("saved-list");
	if (!list) return;

	var box = document.getElementById("saved-box");
	var i18nCount = (box && box.dataset.i18nCount) || "articles";
	var i18nEmptyHeading = (box && box.dataset.i18nEmptyHeading) || "Reading List is Empty";
	var i18nEmptyBody = (box && box.dataset.i18nEmptyBody) || "Click the bookmark icon on any article to save it here.";
	var i18nRemove = (box && box.dataset.i18nRemove) || "REMOVE";

	var bookmarks = [];
	try {
		bookmarks = JSON.parse(localStorage.getItem("bookmarks") || "[]");
	} catch(e) {}

	var countDisplay = document.getElementById("saved-count-display");
	if (countDisplay) {
		countDisplay.textContent = bookmarks.length + " " + i18nCount;
	}

	if (bookmarks.length === 0) {
		list.innerHTML = '<div class="p-12 text-center font-mono"><p class="text-sm font-bold text-text mb-1">' + escHtml(i18nEmptyHeading) + '</p><p class="text-xs text-muted">' + escHtml(i18nEmptyBody) + '</p></div>';
		return;
	}

	var html = '';
	bookmarks.forEach(function(b) {
		html += '<div class="px-5 sm:px-7 py-5 sm:py-6 flex items-baseline justify-between gap-4 hover:bg-hover transition-colors group">';
		html += '  <a href="/blog/' + escHtml(b.slug) + '" class="flex-1 min-w-0 font-bold text-text text-sm sm:text-base group-hover:text-link leading-snug transition-colors truncate">';
		html +=      escHtml(b.title);
		html += '  </a>';
		html += '  <div class="flex items-center gap-3 shrink-0 text-xs font-mono text-muted">';
		if (b.date) {
			html += '    <time datetime="' + escHtml(b.date) + '" class="text-xs text-muted shrink-0">' + escHtml(b.date) + '</time>';
		}
		html += '    <button onclick="removeBookmark(\'' + escHtml(b.slug) + '\')" class="text-[11px] font-bold text-[var(--c-warn)] hover:underline ml-1 uppercase cursor-pointer" title="' + escHtml(i18nRemove) + '">[ ' + escHtml(i18nRemove) + ' ]</button>';
		html += '  </div>';
		html += '</div>';
	});
	list.innerHTML = html;
})();

function escHtml(s) {
	return String(s).replace(/&/g,"&amp;").replace(/</g,"&lt;").replace(/>/g,"&gt;").replace(/"/g,"&quot;");
}

window.removeBookmark = function(slug) {
	var bookmarks = [];
	try {
		bookmarks = JSON.parse(localStorage.getItem("bookmarks") || "[]");
	} catch(e) {}
	bookmarks = bookmarks.filter(function(b) { return b.slug !== slug; });
	localStorage.setItem("bookmarks", JSON.stringify(bookmarks));
	location.reload();
};