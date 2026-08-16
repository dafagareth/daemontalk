(function() {
    // Clipboard helper: falls back to execCommand for HTTP (non-secure) contexts
    window.copyText = function(text, done) {
        if (navigator.clipboard && navigator.clipboard.writeText) {
            navigator.clipboard.writeText(text).then(done).catch(function() { fallback(text, done); });
        } else {
            fallback(text, done);
        }
    };
    function fallback(text, done) {
        var ta = document.createElement('textarea');
        ta.value = text;
        ta.style.cssText = 'position:fixed;top:-9999px;left:-9999px;opacity:0';
        document.body.appendChild(ta);
        ta.focus(); ta.select();
        try { document.execCommand('copy'); if (done) done(); } catch(e) {}
        document.body.removeChild(ta);
    }
})();
