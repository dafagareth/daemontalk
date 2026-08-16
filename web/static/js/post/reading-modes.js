(function() {
    // A+/A- font size
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
    var saved = parseInt(localStorage.getItem("prose-size"));
    if (!isNaN(saved)) applySize(saved);

    // Serif toggle
    function applySerif(on) {
        var el = document.getElementById("prose-body");
        var btns = document.querySelectorAll(".serif-toggle-btn");
        if (!el) return;
        if (on) {
            el.classList.add("prose-serif");
            btns.forEach(function(btn) {
                btn.classList.add("bg-[var(--c-text)]", "text-[var(--c-surface)]");
                btn.classList.remove("text-muted", "bg-transparent");
            });
        } else {
            el.classList.remove("prose-serif");
            btns.forEach(function(btn) {
                btn.classList.remove("bg-[var(--c-text)]", "text-[var(--c-surface)]");
                btn.classList.add("text-muted", "bg-transparent");
            });
        }
    }
    window.toggleSerif = function() {
        var on = !document.getElementById("prose-body").classList.contains("prose-serif");
        localStorage.setItem("prose-serif", on ? "1" : "0");
        applySerif(on);
    };
    applySerif(localStorage.getItem("prose-serif") === "1");

    // Warm tint overlay
    var warmOverlay = document.createElement("div");
    warmOverlay.id = "warm-tint-overlay";
    warmOverlay.style.cssText = "position:fixed;inset:0;pointer-events:none;z-index:9999;background:rgba(255,170,40,0.08);opacity:0;transition:opacity .15s ease;";
    document.body.appendChild(warmOverlay);

    function applyWarmTint(val) {
        var v = parseInt(val) || 0;
        warmOverlay.style.opacity = (v / 100).toFixed(2);
        // Sync all warm sliders
        document.querySelectorAll(".warm-slider").forEach(function(s) { s.value = v; });
    }
    window.setWarmTint = function(val) {
        localStorage.setItem("warm-tint", val);
        applyWarmTint(val);
    };
    var savedWarm = localStorage.getItem("warm-tint");
    if (savedWarm) applyWarmTint(savedWarm);

    // Reading progress bar
    var bar = document.getElementById("reading-progress");
    function updateProgress() {
        if (!bar) return;
        var el = document.getElementById("prose-body");
        if (!el) return;
        var rect = el.getBoundingClientRect();
        var scrollTop = window.pageYOffset || document.documentElement.scrollTop;
        var elTop = rect.top + scrollTop;
        var elHeight = el.offsetHeight;
        var start = elTop;
        var end = elTop + elHeight - window.innerHeight;
        var pct = 0;
        if (scrollTop > start) {
            if (end > start) {
                pct = ((scrollTop - start) / (end - start)) * 100;
            } else {
                pct = 100;
            }
        }
        bar.style.width = Math.min(100, Math.max(0, pct)) + "%";
    }
    window.addEventListener("scroll", updateProgress, { passive: true });
    window.addEventListener("resize", updateProgress, { passive: true });
    updateProgress();

    // Back to top button
    var btt = document.getElementById("back-to-top");
    window.addEventListener("scroll", function() {
        if (!btt) return;
        if (window.scrollY > 300) {
            btt.classList.remove("opacity-0", "pointer-events-none");
            btt.classList.add("opacity-100");
        } else {
            btt.classList.add("opacity-0", "pointer-events-none");
            btt.classList.remove("opacity-100");
        }
    }, { passive: true });
})();
