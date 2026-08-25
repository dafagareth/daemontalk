(function() {
    // Toggle Share Popover Menu
    window.toggleSharePopover = function(e) {
        if (e) e.stopPropagation();
        var menu = document.getElementById('share-popover-menu');
        if (menu) {
            menu.classList.toggle('hidden');
        }
    };

    // Close share popover when clicking anywhere outside
    document.addEventListener('click', function(e) {
        var menu = document.getElementById('share-popover-menu');
        if (menu && !menu.contains(e.target) && !e.target.closest('[onclick*="toggleSharePopover"]')) {
            menu.classList.add('hidden');
        }
    });

    function getPostInfo() {
        var title = (document.querySelector("h1") || {}).textContent || "DaemonTalk Article";
        var url = window.location.href;
        return { title: title.trim(), url: url };
    }

    // Copy standard URL with badge feedback
    window.copyPostUrl = function(btn, e) {
        if (e) e.stopPropagation();
        var info = getPostInfo();
        performCopy(info.url, btn);
    };

    // Copy Markdown formatted link: [Title](URL)
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

    // Social share triggers
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
