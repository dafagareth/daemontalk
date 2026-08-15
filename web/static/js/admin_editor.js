(function() {
			var root = document.getElementById("ed-root");
			var postId = parseInt(root.dataset.postId, 10) || 0;
			var isDraft = root.dataset.draft === "true";
			var statusEl = document.getElementById("ed-status");
			var titleEl = document.getElementById("ed-title");
			var modal = document.getElementById("ed-modal");
			var form = document.getElementById("ed-publish-form");
			var slugEl = document.getElementById("pub-slug");
			var hiddenTitle = document.getElementById("pub-hidden-title");
			var hiddenBody = document.getElementById("pub-hidden-body");
			var slugTouched = !isDraft;

			var quill = new Quill("#ed-body", {
				theme: "bubble",
				placeholder: "Write your story...",
				modules: {
					toolbar: [
						["bold", "italic", "link"],
						[{ header: 1 }, { header: 2 }, { header: 3 }],
						["blockquote", "code-block"]
					]
				}
			});

			var td = new TurndownService({
				headingStyle: "atx",
				codeBlockStyle: "fenced",
				bulletListMarker: "-"
			});

			td.addRule("qlCodeBlock", {
				filter: function(node) {
					return node.nodeName === "PRE" ||
						(node.nodeName === "DIV" && node.classList.contains("ql-code-block-container"));
				},
				replacement: function(content, node) {
					var text;
					var lines = node.querySelectorAll(".ql-code-block");
					if (lines.length > 0) {
						text = Array.prototype.map.call(lines, function(l) { return l.textContent; }).join("\n");
					} else {
						text = node.textContent.replace(/^\n+|\n+$/g, "");
					}
					return "\n```\n" + text + "\n```\n";
				}
			});

			function toMarkdown() {
				return td.turndown(quill.getSemanticHTML()).replace(/\u00A0/g, " ");
			}

			var timer = null;
			var saving = null;
			var dirty = false;

			function save() {
				dirty = false;
				statusEl.textContent = "Saving...";
				var titleVal = (titleEl.value || "").trim();
				var bodyVal = toMarkdown();
				if (hiddenTitle) hiddenTitle.value = titleVal;
				if (hiddenBody) hiddenBody.value = bodyVal;

				saving = fetch("/admin/posts/autosave", {
					method: "POST",
					headers: { "Content-Type": "application/json" },
					body: JSON.stringify({ id: postId, title: titleVal, body: bodyVal, slug: (slugEl ? slugEl.value : "") })
				}).then(function(res) {
					if (!res.ok) throw new Error("HTTP " + res.status);
					return res.json();
				}).then(function(data) {
					postId = data.id;
					root.dataset.postId = data.id;
					form.action = "/admin/posts/" + postId + "/publish";
					if (isDraft && slugEl && !slugTouched && data.slug) slugEl.value = data.slug;
					statusEl.textContent = "Saved at " + data.savedAt;
					return data;
				}).catch(function(err) {
					dirty = true;
					statusEl.textContent = "Failed to save: please try again";
					throw err;
				}).finally(function() {
					saving = null;
				});
				return saving;
			}

			function scheduleSave() {
				dirty = true;
				clearTimeout(timer);
				timer = setTimeout(save, 2000);
			}

			quill.on("text-change", function(_d, _o, source) {
				if (source === "user") scheduleSave();
			});
			titleEl.addEventListener("input", scheduleSave);
			if (slugEl && isDraft) {
				slugEl.addEventListener("input", function() { slugTouched = true; });
			}

			function hideModal() {
				modal.classList.add("hidden");
				modal.classList.remove("flex");
			}

			document.getElementById("ed-publish-btn").addEventListener("click", function() {
				clearTimeout(timer);
				var pending = (dirty || !postId) ? save() : (saving || Promise.resolve());
				Promise.resolve(pending).then(function() {
					if (!postId) return;
					if (hiddenTitle) hiddenTitle.value = (titleEl.value || "").trim();
					if (hiddenBody) hiddenBody.value = toMarkdown();
					modal.classList.remove("hidden");
					modal.classList.add("flex");
				}).catch(function(err) {
					console.error("Save before publish failed:", err);
				});
			});
			document.getElementById("ed-modal-close").addEventListener("click", hideModal);
			modal.addEventListener("click", function(e) {
				if (e.target === modal) hideModal();
			});
			document.addEventListener("keydown", function(e) {
				if (e.key === "Escape") hideModal();
			});
			form.addEventListener("submit", function() {
				if (hiddenTitle) hiddenTitle.value = (titleEl.value || "").trim();
				if (hiddenBody) hiddenBody.value = toMarkdown();
			});
		})();