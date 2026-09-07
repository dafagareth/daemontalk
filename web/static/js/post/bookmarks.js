(function() {
    function updateBookmarkBtnUI(btn, isSaved) {
        var icon = btn.querySelector('.bookmark-icon');
        var text = btn.querySelector('.bookmark-text');
        if (isSaved) {
            btn.classList.add("is-saved", "bg-hover", "text-text", "font-bold");
            btn.classList.remove("border-text/70", "bg-accent/15", "text-accent", "border-accent");
            if (icon) {
                icon.setAttribute("fill", "currentColor");
            }
            if (text) {
                var lang = document.documentElement.lang || "id";
                text.textContent = lang === "id" ? "Tersimpan" : "Saved";
            }
        } else {
            btn.classList.remove("is-saved", "bg-hover", "text-text", "border-text/70", "font-bold", "bg-accent/15", "text-accent", "border-accent");
            if (icon) {
                icon.setAttribute("fill", "none");
            }
            if (text) {
                var lang = document.documentElement.lang || "id";
                text.textContent = lang === "id" ? "Simpan" : "Save";
            }
        }
        if (!icon && !text) {
            btn.textContent = isSaved ? "★" : "☆";
            if (isSaved) btn.classList.add("text-text");
            else btn.classList.remove("text-text", "text-accent");
        }
    }

    function syncPostBookmarks() {
        var bookmarks = [];
        try { bookmarks = JSON.parse(localStorage.getItem("bookmarks") || "[]"); } catch(e) {}
        var saved = {};
        bookmarks.forEach(function(b) { saved[b.slug] = true; });
        document.querySelectorAll(".bookmark-btn").forEach(function(btn) {
            updateBookmarkBtnUI(btn, !!saved[btn.dataset.slug]);
        });
    }
    syncPostBookmarks();
    document.addEventListener("htmx:afterSwap", syncPostBookmarks);

    window.toggleBookmark = function(btn) {
        var slug = btn.dataset.slug;
        var title = btn.dataset.title || "";
        var date = btn.dataset.date || "";
        var bookmarks = [];
        try { bookmarks = JSON.parse(localStorage.getItem("bookmarks") || "[]"); } catch(e) {}
        var idx = bookmarks.findIndex(function(b) { return b.slug === slug; });
        var isNowSaved = false;
        if (idx >= 0) {
            bookmarks.splice(idx, 1);
            isNowSaved = false;
        } else {
            bookmarks.unshift({slug: slug, title: title, date: date});
            isNowSaved = true;
        }
        localStorage.setItem("bookmarks", JSON.stringify(bookmarks));

        document.querySelectorAll('.bookmark-btn[data-slug="' + slug + '"]').forEach(function(el) {
            updateBookmarkBtnUI(el, isNowSaved);
        });
    };
})();
