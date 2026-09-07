(function() {
	var list = document.getElementById("saved-list");
	if (!list) return;

	var box = document.getElementById("saved-box");
	var i18nCount = (box && box.dataset.i18nCount) || "articles";
	var i18nEmptyHeading = (box && box.dataset.i18nEmptyHeading) || "Reading List is Empty";
	var i18nEmptyBody = (box && box.dataset.i18nEmptyBody) || "Click the bookmark icon on any article to save it here.";
	var i18nRemove = (box && box.dataset.i18nRemove) || "REMOVE";
	var lang = (box && box.dataset.lang) || "";
	var prefix = lang === "id" ? "/id" : "";

	var bookmarks = [];
	try {
		bookmarks = JSON.parse(localStorage.getItem("bookmarks") || "[]");
	} catch(e) {}

	var countDisplay = document.getElementById("saved-count-display");
	if (countDisplay) {
		countDisplay.textContent = bookmarks.length + " " + i18nCount;
	}

	if (bookmarks.length === 0) {
		list.innerHTML = '<div class="p-12 text-center font-sans"><p class="text-sm font-bold text-text mb-1">' + escHtml(i18nEmptyHeading) + '</p><p class="text-xs text-muted">' + escHtml(i18nEmptyBody) + '</p></div>';
		return;
	}

	var html = '';
	bookmarks.forEach(function(b) {
		html += '<div class="px-3 sm:px-6 py-4 sm:py-4.5 flex items-center justify-between gap-3 sm:gap-4 hover:bg-hover transition-colors group">';
		html += '  <div class="flex-1 min-w-0 pr-1 sm:pr-4">';
		html += '    <a href="' + prefix + '/blog/' + escHtml(b.slug) + '" class="block font-bold text-text text-base sm:text-lg group-hover:text-link leading-snug transition-colors line-clamp-2 sm:line-clamp-1">';
		html +=        escHtml(b.title);
		html += '    </a>';
		if (b.date) {
			html += '    <time datetime="' + escHtml(b.date) + '" class="block sm:hidden text-xs font-sans text-muted mt-1.5">' + escHtml(b.date) + '</time>';
		}
		html += '  </div>';
		html += '  <div class="flex items-center gap-4 shrink-0 text-xs font-sans text-muted">';
		if (b.date) {
			html += '    <time datetime="' + escHtml(b.date) + '" class="hidden sm:inline-block text-xs text-muted shrink-0 font-sans">' + escHtml(b.date) + '</time>';
		}
		html += '    <button onclick="removeBookmark(\'' + escHtml(b.slug) + '\')" class="text-xs sm:text-[13px] font-sans font-semibold text-muted hover:text-rose-500 transition-colors uppercase cursor-pointer" title="' + escHtml(i18nRemove) + '">' + escHtml(i18nRemove) + '</button>';
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
