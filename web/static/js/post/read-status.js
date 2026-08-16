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
