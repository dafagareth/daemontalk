(function() {
    function isMobileDevice() {
        return (
            /Android|webOS|iPhone|iPad|iPod|BlackBerry|IEMobile|Opera Mini/i.test(navigator.userAgent) ||
            (window.matchMedia && window.matchMedia("(max-width: 768px)").matches) ||
            ('ontouchstart' in window && window.innerWidth <= 1024)
        );
    }

    function showShareMenu(menuId) {
        document.querySelectorAll('[id*="more-popover-menu"], #reading-popover-menu').forEach(function(m) {
            m.classList.add('hidden');
        });
        var id = menuId || 'share-popover-menu';
        var menu = document.getElementById(id);
        if (menu) {
            var isClosed = menu.classList.contains('hidden');
            document.querySelectorAll('[id*="share-popover-menu"]').forEach(function(m) {
                m.classList.add('hidden');
            });
            if (isClosed) {
                menu.classList.remove('hidden');
            }
        }
    }

    window.toggleSharePopover = function(e, menuId) {
        if (e) e.stopPropagation();

        var info = getPostInfo();
        if (navigator.share && isMobileDevice()) {
            var shareData = {
                title: info.title,
                url: info.url
            };
            if (info.text) {
                shareData.text = info.text;
            }
            navigator.share(shareData).catch(function(err) {
                if (err && err.name !== 'AbortError') {
                    showShareMenu(menuId);
                }
            });
            return;
        }

        showShareMenu(menuId);
    };

    document.addEventListener('click', function(e) {
        document.querySelectorAll('[id*="share-popover-menu"]').forEach(function(menu) {
            if (!menu.contains(e.target) && !e.target.closest('[onclick*="toggleSharePopover"]')) {
                menu.classList.add('hidden');
            }
        });
    });

    function getPostInfo() {
        var title = (document.querySelector("h1") || {}).textContent || "DaemonTalk Article";
        var descMeta = document.querySelector('meta[name="description"]');
        var text = descMeta ? descMeta.getAttribute("content") : "";
        var url = window.location.href;
        return { title: title.trim(), text: text, url: url };
    }

    window.copyPostUrl = function(btn, e) {
        if (e) e.stopPropagation();
        var info = getPostInfo();
        performCopy(info.url, btn);
    };

    window.copyMarkdownLink = function(btn, e) {
        if (e) e.stopPropagation();
        var info = getPostInfo();
        var md = '[' + info.title + '](' + info.url + ')';
        performCopy(md, btn);
    };

    function performCopy(text, btn) {
        var copySuccess = function() {
            var label = btn.querySelector('.share-label-text');
            var status = btn.querySelector('.share-copied-status');
            if (label && status) {
                label.classList.add('hidden');
                status.classList.remove('hidden');
                setTimeout(function() {
                    label.classList.remove('hidden');
                    status.classList.add('hidden');
                }, 2000);
            }
        };

        if (navigator.clipboard && navigator.clipboard.writeText) {
            navigator.clipboard.writeText(text).then(copySuccess).catch(function() {
                fallbackCopy(text, copySuccess);
            });
        } else {
            fallbackCopy(text, copySuccess);
        }
    }

    function fallbackCopy(text, cb) {
        var ta = document.createElement('textarea');
        ta.value = text;
        ta.style.position = 'fixed';
        ta.style.opacity = '0';
        document.body.appendChild(ta);
        ta.select();
        try {
            document.execCommand('copy');
            if (cb) cb();
        } catch (err) {}
        document.body.removeChild(ta);
    }

    window.shareToBluesky = function(e) {
        if (e) e.stopPropagation();
        var info = getPostInfo();
        window.open('https://bsky.app/intent/compose?text=' + encodeURIComponent(info.title + '\n' + info.url), '_blank', 'width=600,height=450');
    };

    window.shareToTwitter = function(e) {
        if (e) e.stopPropagation();
        var info = getPostInfo();
        window.open('https://twitter.com/intent/tweet?text=' + encodeURIComponent(info.title) + '&url=' + encodeURIComponent(info.url), '_blank', 'width=600,height=450');
    };

    window.shareToLinkedIn = function(e) {
        if (e) e.stopPropagation();
        var info = getPostInfo();
        window.open('https://www.linkedin.com/sharing/share-offsite/?url=' + encodeURIComponent(info.url), '_blank', 'width=600,height=500');
    };

    window.shareToThreads = function(e) {
        if (e) e.stopPropagation();
        var info = getPostInfo();
        window.open('https://www.threads.net/intent/post?text=' + encodeURIComponent(info.title + ' ' + info.url), '_blank', 'width=600,height=500');
    };
})();
