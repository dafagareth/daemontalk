(function() {
    // Share: Instagram (Web Share API on mobile, copy link on desktop)
    window.shareToInstagram = function(btn) {
        var url = window.location.href;
        var title = (document.querySelector("h1") || {}).textContent || "daemontalk.com";
        if (navigator.share) {
            navigator.share({ title: title, url: url }).catch(function() {});
        } else {
            if (window.copyText) {
                window.copyText(url, function() {
                    var old = btn.textContent;
                    btn.textContent = "Link copied!";
                    setTimeout(function() { btn.textContent = old; }, 2000);
                });
            }
        }
    };

    // Share: copy link
    window.copyPostLink = function(btn) {
        if (window.copyText) {
            window.copyText(window.location.href, function() {
                var old = btn.textContent;
                btn.textContent = "Copied!";
                setTimeout(function() { btn.textContent = old; }, 2000);
            });
        }
    };
})();
