(function() {
    // 1. Font Size Adjuster (A- / A+)
    var SIZES = [13, 14, 15, 16, 17, 18, 20];
    var DEFAULT = 3;
    function applySize(idx) {
        var el = document.getElementById("prose-body");
        if (el) {
            el.style.fontSize = SIZES[idx] + "px";
            document.documentElement.style.setProperty("--prose-size", SIZES[idx] + "px");
        }
    }
    window.adjustProseSize = function(dir) {
        var idx = parseInt(localStorage.getItem("prose-size"));
        if (isNaN(idx)) idx = DEFAULT;
        idx = Math.max(0, Math.min(SIZES.length - 1, idx + dir));
        localStorage.setItem("prose-size", idx);
        applySize(idx);
    };
    var savedSize = parseInt(localStorage.getItem("prose-size"));
    if (!isNaN(savedSize)) applySize(savedSize);

    // 2. Serif / Sans Font Toggle (Default: Serif, No Orange Highlight)
    function applySerif(isSerif) {
        var el = document.getElementById("prose-body");
        var label = document.getElementById("serif-toggle-label");
        if (!el) return;
        if (isSerif) {
            el.classList.add("prose-serif");
            el.classList.remove("prose-sans");
            if (label) label.textContent = "Serif";
        } else {
            el.classList.remove("prose-serif");
            el.classList.add("prose-sans");
            if (label) label.textContent = "Sans";
        }
    }

    window.toggleSerif = function() {
        var el = document.getElementById("prose-body");
        var currentIsSerif = el ? !el.classList.contains("prose-sans") : true;
        var nextIsSerif = !currentIsSerif;
        localStorage.setItem("prose-serif", nextIsSerif ? "1" : "0");
        applySerif(nextIsSerif);
    };

    // Initialize state on page load: default to Serif (unless explicitly set to "0" for Sans)
    var savedSerif = localStorage.getItem("prose-serif");
    applySerif(savedSerif !== "0");
})();
