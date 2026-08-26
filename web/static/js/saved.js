(function() {
			var list = document.getElementById("saved-list");
			var bookmarks = [];
			try {
				bookmarks = JSON.parse(localStorage.getItem("bookmarks") || "[]");
			} catch(e) {}
			var countDisplay = document.getElementById("saved-count-display");
			if (countDisplay) {
				countDisplay.textContent = bookmarks.length + " Entries";
			}

			if (bookmarks.length === 0) {
				list.innerHTML = '<div class="p-8 sm:p-16 text-center flex flex-col items-center justify-center"><p class="text-[10px] uppercase tracking-widest font-mono text-muted mb-3 border border-border px-3 py-1 bg-surface">LEDGER EMPTY</p><h3 class="text-base sm:text-lg text-text font-bold mb-2">No dispatches bookmarked yet.</h3><p class="text-sm text-muted">Click the bookmark icon on any post to save it to your local ledger.</p></div>';
				return;
			}

			var html = '';
			bookmarks.forEach(function(b, i) {
				var bgClass = i % 2 !== 0 ? 'bg-surface/30' : 'bg-bg';
				
				html += '<div class="group flex flex-col sm:flex-row sm:items-center justify-between p-5 sm:p-6 sm:px-8 gap-4 hover:bg-hover/20 transition-colors ' + bgClass + '">';
				
				// Left: Title & Context
				html += '	<a href="/blog/' + b.slug + '" class="flex-1 min-w-0 flex flex-col justify-center">';
				html += '		<h3 class="display text-base sm:text-lg font-bold text-text group-hover:text-link transition-colors leading-snug line-clamp-2 mb-1.5">' + escHtml(b.title) + '</h3>';
				html += '		<p class="text-[11px] font-mono text-muted uppercase tracking-wider">ENTRY / ' + escHtml(b.slug).split('-')[0] + '</p>';
				html += '	</a>';
				
				// Right: Meta & Actions
				html += '	<div class="flex items-center sm:flex-col sm:items-end justify-between sm:justify-center shrink-0 gap-3 sm:gap-2">';
				html += '		<div class="font-mono text-[10px] text-muted tracking-wider uppercase">' + escHtml(b.date || "UNKNOWN DATE") + '</div>';
				html += '		<button onclick="removeBookmark(\'' + b.slug + '\')" class="font-mono text-[10px] font-bold text-[var(--c-warn)] hover:underline uppercase tracking-wider" title="Remove from ledger">[ DEL ]</button>';
				html += '	</div>';
				
				html += '</div>';
			});
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