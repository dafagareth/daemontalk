(function() {
    window.commentVisibleCount = 3;

    function handleLoadMoreComments() {
        var container = document.getElementById('comment-items-container');
        var btnContainer = document.getElementById('comment-load-more-container');
        if (!container) return;

        var hiddenItems = container.querySelectorAll('.comment-root-item.hidden');
        var step = 10;
        for (var i = 0; i < hiddenItems.length && i < step; i++) {
            hiddenItems[i].classList.remove('hidden');
        }

        var currentlyVisible = container.querySelectorAll('.comment-root-item:not(.hidden)').length;
        window.commentVisibleCount = currentlyVisible;

        var remaining = container.querySelectorAll('.comment-root-item.hidden').length;
        if (remaining === 0) {
            if (btnContainer) btnContainer.classList.add('hidden');
        } else {
            var countEl = document.getElementById('comment-remaining-count');
            if (countEl) countEl.textContent = '(' + remaining + ')';
        }
    }
    window.handleLoadMoreComments = handleLoadMoreComments;

    function syncCommentVisibility() {
        var container = document.getElementById('comment-items-container');
        var btnContainer = document.getElementById('comment-load-more-container');
        if (!container) return;

        var targetVisible = window.commentVisibleCount || 3;
        var items = container.querySelectorAll('.comment-root-item');
        var hiddenCount = 0;
        for (var i = 0; i < items.length; i++) {
            if (i < targetVisible) {
                items[i].classList.remove('hidden');
            } else {
                items[i].classList.add('hidden');
                hiddenCount++;
            }
        }

        if (btnContainer) {
            if (hiddenCount === 0) {
                btnContainer.classList.add('hidden');
            } else {
                btnContainer.classList.remove('hidden');
                var countEl = document.getElementById('comment-remaining-count');
                if (countEl) countEl.textContent = '(' + hiddenCount + ')';
            }
        }
    }
    window.syncCommentVisibility = syncCommentVisibility;

    function onCommentPosted(parentId) {
        if (parentId) {
            var parentEl = document.getElementById('comment-' + parentId);
            if (parentEl) {
                var details = parentEl.querySelector('details.comment-replies-group');
                if (details) details.open = true;
                setTimeout(function() {
                    parentEl.scrollIntoView({ behavior: 'smooth', block: 'center' });
                }, 100);
            }
        } else {
            var container = document.getElementById('comment-items-container');
            if (container) {
                var items = container.querySelectorAll('.comment-root-item');
                if (items.length > 0) {
                    var firstItem = items[0];
                    firstItem.classList.remove('hidden');
                    firstItem.classList.add('bg-hover/60');
                    setTimeout(function() {
                        firstItem.classList.remove('bg-hover/60');
                    }, 2000);
                    setTimeout(function() {
                        firstItem.scrollIntoView({ behavior: 'smooth', block: 'center' });
                    }, 100);
                }
            }
            var btnContainer = document.getElementById('comment-load-more-container');
            if (btnContainer) {
                var remaining = document.querySelectorAll('.comment-root-item.hidden').length;
                if (remaining === 0) {
                    btnContainer.classList.add('hidden');
                } else {
                    var countEl = document.getElementById('comment-remaining-count');
                    if (countEl) countEl.textContent = '(' + remaining + ')';
                }
            }
        }
    }
    window.onCommentPosted = onCommentPosted;

    document.addEventListener('htmx:afterSwap', function(e) {
        if (e.target && (e.target.id === 'comment-list' || (e.target.querySelector && e.target.querySelector('#comment-list')))) {
            syncCommentVisibility();
        }
    });

    document.addEventListener('click', function(e) {
        if (e.target && e.target.closest && e.target.closest('#comment-load-more-btn')) {
            e.preventDefault();
            handleLoadMoreComments();
        }
    });

    var cl = document.getElementById("comment-list");
    if (cl) {
        var s = cl.getAttribute("data-slug");
        var prefix = window.location.pathname.startsWith("/id/") ? "/id" : "";
        if (s) {
            var evtSource = new EventSource(prefix + "/blog/" + s + "/comments/stream");
            evtSource.addEventListener("new_comment", function(e) {
                htmx.ajax('GET', prefix + '/blog/' + s + '/comments', '#comment-list');
            });
        }
    }
})();
