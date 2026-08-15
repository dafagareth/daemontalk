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

(function() {
			// A+/A- font size
			var SIZES = [13, 14, 15, 16, 17, 18, 20];
			var DEFAULT = 3;
			function applySize(idx) {
				var el = document.getElementById("prose-body");
				if (el) el.style.fontSize = SIZES[idx] + "px";
			}
			window.adjustProseSize = function(dir) {
				var idx = parseInt(localStorage.getItem("prose-size"));
				if (isNaN(idx)) idx = DEFAULT;
				idx = Math.max(0, Math.min(SIZES.length - 1, idx + dir));
				localStorage.setItem("prose-size", idx);
				applySize(idx);
			};
			var saved = parseInt(localStorage.getItem("prose-size"));
			if (!isNaN(saved)) applySize(saved);

			// Serif toggle
			function applySerif(on) {
				var el = document.getElementById("prose-body");
				var btns = document.querySelectorAll(".serif-toggle-btn");
				if (!el) return;
				if (on) {
					el.classList.add("prose-serif");
					btns.forEach(function(btn) {
						btn.classList.add("bg-[var(--c-text)]", "text-[var(--c-surface)]");
						btn.classList.remove("text-muted", "bg-transparent");
					});
				} else {
					el.classList.remove("prose-serif");
					btns.forEach(function(btn) {
						btn.classList.remove("bg-[var(--c-text)]", "text-[var(--c-surface)]");
						btn.classList.add("text-muted", "bg-transparent");
					});
				}
			}
			window.toggleSerif = function() {
				var on = !document.getElementById("prose-body").classList.contains("prose-serif");
				localStorage.setItem("prose-serif", on ? "1" : "0");
				applySerif(on);
			};
			applySerif(localStorage.getItem("prose-serif") === "1");

			// Reading progress bar
			var bar = document.getElementById("reading-progress");
			function updateProgress() {
				if (!bar) return;
				var el = document.getElementById("prose-body");
				if (!el) return;
				var rect = el.getBoundingClientRect();
				var scrollTop = window.pageYOffset || document.documentElement.scrollTop;
				var elTop = rect.top + scrollTop;
				var elHeight = el.offsetHeight;
				var start = elTop;
				var end = elTop + elHeight - window.innerHeight;
				var pct = 0;
				if (scrollTop > start) {
					if (end > start) {
						pct = ((scrollTop - start) / (end - start)) * 100;
					} else {
						pct = 100;
					}
				}
				bar.style.width = Math.min(100, Math.max(0, pct)) + "%";
			}
			window.addEventListener("scroll", updateProgress, { passive: true });
			window.addEventListener("resize", updateProgress, { passive: true });
			updateProgress();

			// Back to top button
			var btt = document.getElementById("back-to-top");
			window.addEventListener("scroll", function() {
				if (!btt) return;
				if (window.scrollY > 300) {
					btt.classList.remove("opacity-0", "pointer-events-none");
					btt.classList.add("opacity-100");
				} else {
					btt.classList.add("opacity-0", "pointer-events-none");
					btt.classList.remove("opacity-100");
				}
			}, { passive: true });

			// Enhanced Code block wrapper: adds language badges and copy button.
			document.querySelectorAll("#prose-body pre").forEach(function(pre) {
				var wrap = document.createElement("div");
				wrap.className = "code-wrap";
				pre.parentNode.insertBefore(wrap, pre);
				wrap.appendChild(pre);

				var codeEl = pre.querySelector("code");
				var codeText = codeEl ? codeEl.innerText : pre.innerText;

				// Detect language from class or content
				var lang = "";
				var classNames = (pre.className + " " + (codeEl ? codeEl.className : "")).toLowerCase();
				var langMatch = classNames.match(/language-([a-z0-9_-]+)/) || classNames.match(/lang-([a-z0-9_-]+)/);
				if (langMatch) {
					lang = langMatch[1];
				} else if (codeText.indexOf("package main") >= 0 || codeText.indexOf("fmt.Print") >= 0 || codeText.indexOf("func ") >= 0) {
					lang = "go";
				} else if (codeText.indexOf("def ") >= 0 || (codeText.indexOf("import ") >= 0 && codeText.indexOf("from ") >= 0)) {
					lang = "python";
				} else if (codeText.indexOf("fn ") >= 0 || codeText.indexOf("let mut ") >= 0) {
					lang = "rust";
				} else if (codeText.indexOf("const ") >= 0 || codeText.indexOf("console.log") >= 0) {
					lang = "javascript";
				} else if (codeText.startsWith("$ ") || codeText.startsWith("sudo ") || codeText.startsWith("curl ") || codeText.startsWith("docker ")) {
					lang = "bash";
				}

				// Detect diagram / architecture schematics
				var isDiagram = codeText.indexOf("┌") >= 0 || codeText.indexOf("┼") >= 0 || codeText.indexOf("-->") >= 0 || codeText.indexOf("──►") >= 0 || codeText.indexOf("flowchart") >= 0 || codeText.indexOf("graph TD") >= 0 || lang === "mermaid" || lang === "diagram";
				if (isDiagram && !lang) {
					lang = "diagram";
				}

				// Create action toolbar
				var toolbar = document.createElement("div");
				toolbar.className = "code-toolbar";

				// Language or Diagram badge
				if (lang) {
					var badge = document.createElement("span");
					badge.className = "code-lang-badge" + (isDiagram ? " text-link font-bold" : "");
					badge.textContent = isDiagram ? "ARCHITECTURE DIAGRAM" : lang.toUpperCase();
					toolbar.appendChild(badge);
				}

				// Interactive Diagram Zoom/Focus Toggle
				if (isDiagram) {
					var zoomBtn = document.createElement("button");
					zoomBtn.textContent = "expand";
					zoomBtn.className = "copy-btn";
					zoomBtn.setAttribute("aria-label", "Expand diagram");
					zoomBtn.addEventListener("click", function() {
						wrap.classList.toggle("diagram-expanded");
						if (wrap.classList.contains("diagram-expanded")) {
							zoomBtn.textContent = "collapse";
							pre.style.maxHeight = "none";
							pre.style.fontSize = "13px";
						} else {
							zoomBtn.textContent = "expand";
							pre.style.maxHeight = "";
							pre.style.fontSize = "";
						}
					});
					toolbar.appendChild(zoomBtn);
				}

				// Copy button
				var copyBtn = document.createElement("button");
				copyBtn.textContent = "copy";
				copyBtn.className = "copy-btn";
				copyBtn.setAttribute("aria-label", "Copy code");
				copyBtn.addEventListener("click", function() {
					var text = codeEl ? codeEl.innerText : pre.innerText;
					copyText(text, function() {
						copyBtn.textContent = "copied!";
						setTimeout(function() { copyBtn.textContent = "copy"; }, 2000);
					});
				});
				toolbar.appendChild(copyBtn);

				wrap.appendChild(toolbar);
			});

			// Clipboard helper: falls back to execCommand for HTTP (non-secure) contexts
			function copyText(text, done) {
				if (navigator.clipboard && navigator.clipboard.writeText) {
					navigator.clipboard.writeText(text).then(done).catch(function() { fallback(text, done); });
				} else {
					fallback(text, done);
				}
			}
			function fallback(text, done) {
				var ta = document.createElement('textarea');
				ta.value = text;
				ta.style.cssText = 'position:fixed;top:-9999px;left:-9999px;opacity:0';
				document.body.appendChild(ta);
				ta.focus(); ta.select();
				try { document.execCommand('copy'); if (done) done(); } catch(e) {}
				document.body.removeChild(ta);
			}

			// Share: Instagram (Web Share API on mobile, copy link on desktop)
			window.shareToInstagram = function(btn) {
				var url = window.location.href;
				var title = (document.querySelector("h1") || {}).textContent || "daemontalk.com";
				if (navigator.share) {
					navigator.share({ title: title, url: url }).catch(function() {});
				} else {
					copyText(url, function() {
						var old = btn.textContent;
						btn.textContent = "Link copied!";
						setTimeout(function() { btn.textContent = old; }, 2000);
					});
				}
			};

			// Share: copy link
			window.copyPostLink = function(btn) {
				copyText(window.location.href, function() {
					var old = btn.textContent;
					btn.textContent = "Copied!";
					setTimeout(function() { btn.textContent = old; }, 2000);
				});
			};

			// ToC: highlight active section
			var tocLinks = document.querySelectorAll(".toc-link");
			if (tocLinks.length > 0) {
				var headings = Array.from(tocLinks).map(function(a) {
					return document.getElementById(a.getAttribute("href").slice(1));
				}).filter(Boolean);
				function setActive() {
					var scrollY = window.scrollY + 100;
					var active = headings[0];
					for (var i = 0; i < headings.length; i++) {
						if (headings[i].offsetTop <= scrollY) active = headings[i];
					}
					tocLinks.forEach(function(a) {
						var isActive = active && a.getAttribute("href") === "#" + active.id;
						a.classList.toggle("toc-active", isActive);
					});
				}
				tocLinks.forEach(function(a) {
					a.addEventListener("click", function(e) {
						var targetId = a.getAttribute("href").slice(1);
						var targetEl = document.getElementById(targetId);
						if (targetEl) {
							e.preventDefault();
							var navOffset = 80;
							var top = targetEl.getBoundingClientRect().top + window.pageYOffset - navOffset;
							window.scrollTo({ top: Math.max(0, top), behavior: "smooth" });
							history.pushState(null, "", "#" + targetId);
						}
					});
				});

				window.addEventListener("scroll", setActive, { passive: true });
				setActive();
			}

			// Mobile ToC chevron
			var details = document.querySelector("details");
			if (details) {
				details.addEventListener("toggle", function() {
					var chevron = details.querySelector(".toc-chevron");
					if (chevron) chevron.style.transform = details.open ? "rotate(180deg)" : "";
				});
			}

		})();

// Mark post as read in localStorage
		(function() {
			var path = window.location.pathname;
			var match = path.match(/\/blog\/([^\/]+)$/);
			if (match) {
				var slug = match[1];
				var read = [];
				try { read = JSON.parse(localStorage.getItem('readPosts') || '[]'); } catch(e) {}
				if (read.indexOf(slug) === -1) {
					read.push(slug);
					// Keep only the last 200 entries
					if (read.length > 200) read = read.slice(-200);
					localStorage.setItem('readPosts', JSON.stringify(read));
				}
			}
		})();