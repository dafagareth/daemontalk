(function() {
			var list = document.getElementById("saved-list");
			var bookmarks = [];
			try {
				bookmarks = JSON.parse(localStorage.getItem("bookmarks") || "[]");
			} catch(e) {}

			if (bookmarks.length === 0) {
				list.innerHTML = '<p class="text-sm text-muted">No bookmarks yet. Click the ☆ on any post to save it here.</p>';
				return;
			}

			var html = '<ul>';
			bookmarks.forEach(function(b) {
				html += '<li class="py-4 flex items-center justify-between gap-4">';
				html += '<a href="/blog/' + b.slug + '" class="text-sm text-link hover:underline flex-1 min-w-0 truncate">' + escHtml(b.title) + '</a>';
				html += '<div class="flex items-center gap-3 shrink-0">';
				html += '<span class="text-xs text-muted">' + escHtml(b.date || "") + '</span>';
				html += '<button onclick="removeBookmark(\'' + b.slug + '\')" class="text-xs text-[var(--c-warn)] hover:underline">Remove</button>';
				html += '</div>';
				html += '</li>';
			});
			html += '</ul>';
			list.innerHTML = html;
		})();

		function escHtml(s) {
			return String(s).replace(/&/g,"&amp;").replace(/</g,"&lt;").replace(/>/g,"&gt;").replace(/"/g,"&quot;");
		}

		window.removeBookmark = function(slug) {
			var bookmarks = [];
			try { bookmarks = JSON.parse(localStorage.getItem("bookmarks") || "[]"); } catch(e) {}
			bookmarks = bookmarks.filter(function(b) { return b.slug !== slug; });
			localStorage.setItem("bookmarks", JSON.stringify(bookmarks));
			location.reload();
		};