// Bookmark toggle (shared with blog list)
(function() {
    var bookmarks = [];
    try { bookmarks = JSON.parse(localStorage.getItem("bookmarks") || "[]"); } catch(e) {}
    var saved = {};
    bookmarks.forEach(function(b) { saved[b.slug] = true; });
    document.querySelectorAll(".bookmark-btn").forEach(function(btn) {
        if (saved[btn.dataset.slug]) {
            btn.textContent = "★";
            btn.classList.add("text-link");
        }
    });
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
        btn.classList.remove("text-link");
    } else {
        bookmarks.unshift({slug: slug, title: title, date: date});
        btn.textContent = "★";
        btn.classList.add("text-link");
    }
    localStorage.setItem("bookmarks", JSON.stringify(bookmarks));
};
