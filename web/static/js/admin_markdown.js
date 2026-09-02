// Markdown Raw Editor: Tab indentation, live statistics, draft recovery, and image upload
(function() {
	var tx = document.getElementById('md-editor-textarea');
	var charCount = document.getElementById('md-char-count');
	var wordCount = document.getElementById('md-word-count');

	function updateCounts() {
		if (!tx) return;
		var val = tx.value || "";
		if (charCount) charCount.textContent = val.length + " chars";
		var words = val.trim().split(/\s+/).filter(Boolean).length;
		if (wordCount) wordCount.textContent = words + " words";
	}

	if (tx) {
		var slugEl = document.getElementById('md-slug-val');
		var draftKey = 'daemontalk_draft_' + ((slugEl && slugEl.value) ? slugEl.value : 'new');
		
		var savedDraft = localStorage.getItem(draftKey);
		if (savedDraft && savedDraft !== tx.value) {
			if (confirm("Found an unsaved draft in your browser. Do you want to restore it?")) {
				tx.value = savedDraft;
			} else {
				localStorage.removeItem(draftKey);
			}
		}

		var saveTimeout;
		tx.addEventListener('input', function() {
			clearTimeout(saveTimeout);
			saveTimeout = setTimeout(function() {
				localStorage.setItem(draftKey, tx.value);
			}, 3000);
			updateCounts();
		});

		var form = tx.closest('form');
		if (form) {
			form.addEventListener('submit', function() {
				localStorage.removeItem(draftKey);
			});
		}

		tx.addEventListener('keydown', function(e) {
			if (e.key === 'Tab') {
				e.preventDefault();
				var start = this.selectionStart;
				var end = this.selectionEnd;
				this.value = this.value.substring(0, start) + "\t" + this.value.substring(end);
				this.selectionStart = this.selectionEnd = start + 1;
			}
		});
		tx.addEventListener('input', updateCounts);
		updateCounts();
	}

	window.uploadMarkdownImage = function(input) {
		if (!input.files || !input.files[0]) return;
		var file = input.files[0];
		var status = document.getElementById('img-upload-status');
		if (status) status.textContent = "Uploading " + file.name + "...";

		var targetSlug = "";
		if (tx && tx.value) {
			var match = tx.value.match(/^slug:\s*["']?([a-zA-Z0-9_-]+)["']?/m);
			if (match && match[1]) {
				targetSlug = match[1];
			}
		}

		if (!targetSlug) {
			var rootEl = document.getElementById('md-editor-root');
			targetSlug = (rootEl ? rootEl.getAttribute('data-slug') : "") || 
						 (document.getElementById('md-slug-val') || {}).value || 
						 new URLSearchParams(window.location.search).get('slug') || 
						 "uploads";
		}

		var formData = new FormData();
		formData.append("image", file);
		formData.append("slug", targetSlug);

		fetch("/admin/upload-image", {
			method: "POST",
			body: formData
		})
		.then(function(res) {
			if (!res.ok) throw new Error("Upload failed");
			return res.json();
		})
		.then(function(data) {
			if (data.markdown && tx) {
				var start = tx.selectionStart;
				var end = tx.selectionEnd;
				var text = tx.value;
				var insert = "\n" + data.markdown + "\n";
				tx.value = text.substring(0, start) + insert + text.substring(end);
				tx.selectionStart = tx.selectionEnd = start + insert.length;
				if (status) status.textContent = "✓ Uploaded to /" + targetSlug + ": " + data.filename;
				updateCounts();
			}
		})
		.catch(function(err) {
			if (status) status.textContent = "Error: " + err.message;
		});
	};
})();
