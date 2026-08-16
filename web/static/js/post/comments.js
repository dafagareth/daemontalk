// Real-time comments SSE
(function() {
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
