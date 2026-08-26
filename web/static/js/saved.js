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
				var num = (i + 1).toString().padStart(2, '0');
				
				html += '<div class="group flex flex-col sm:flex-row items-stretch hover:bg-surface transition-colors">';
				
				// Col 1: Huge Ledger Number
				html += '	<div class="hidden sm:flex w-20 shrink-0 items-center justify-center border-r border-border font-mono text-2xl font-black text-muted opacity-30 group-hover:opacity-100 transition-opacity bg-surface/50">';
				html += '		' + num;
				html += '	</div>';
				
				// Col 2: Title & Context
				html += '	<a href="/blog/' + b.slug + '" class="flex-1 min-w-0 p-5 sm:p-6 sm:px-8 flex flex-col justify-center border-b sm:border-b-0 border-border">';
				html += '		<h3 class="display text-base sm:text-lg font-bold text-text group-hover:text-link transition-colors leading-snug line-clamp-2 mb-1.5">' + escHtml(b.title) + '</h3>';
				html += '		<p class="text-xs font-mono text-muted uppercase tracking-wider">ENTRY / ' + escHtml(b.slug).split('-')[0] + '</p>';
				html += '	</a>';
				
				// Col 3: Meta & Actions
				html += '	<div class="sm:w-48 shrink-0 p-4 sm:p-6 flex sm:flex-col items-center sm:items-end justify-between sm:justify-center border-l-0 sm:border-l border-border gap-2 bg-surface/30 group-hover:bg-transparent transition-colors">';
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