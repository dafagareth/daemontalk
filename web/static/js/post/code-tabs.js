(function() {
    // Multi-File Code Tabs Interactive Switching & Touch/Drag Swiping
    document.querySelectorAll("[data-code-tabs]").forEach(function(wrap) {
        var navTrack = wrap.querySelector(".tabs-nav-track");
        var buttons = wrap.querySelectorAll(".tab-btn");
        var panes = wrap.querySelectorAll(".tab-pane");
        var copyBtn = wrap.querySelector(".copy-tab-code");

        // Mouse drag-to-scroll for horizontal tabs header
        if (navTrack) {
            var isDown = false;
            var startX, scrollLeft;

            navTrack.addEventListener("mousedown", function(e) {
                isDown = true;
                startX = e.pageX - navTrack.offsetLeft;
                scrollLeft = navTrack.scrollLeft;
            });
            navTrack.addEventListener("mouseleave", function() {
                isDown = false;
            });
            navTrack.addEventListener("mouseup", function() {
                isDown = false;
            });
            navTrack.addEventListener("mousemove", function(e) {
                if (!isDown) return;
                e.preventDefault();
                var x = e.pageX - navTrack.offsetLeft;
                var walk = (x - startX) * 1.5;
                navTrack.scrollLeft = scrollLeft - walk;
            });

            // Touch swipe for mobile
            var touchStartX;
            navTrack.addEventListener("touchstart", function(e) {
                touchStartX = e.touches[0].pageX;
                scrollLeft = navTrack.scrollLeft;
            }, {passive: true});
            navTrack.addEventListener("touchmove", function(e) {
                var x = e.touches[0].pageX;
                navTrack.scrollLeft = scrollLeft - (x - touchStartX);
            }, {passive: true});
        }

        buttons.forEach(function(btn) {
            btn.addEventListener("click", function() {
                var idx = btn.getAttribute("data-tab-index");
                buttons.forEach(function(b) {
                    b.classList.remove("active", "border-link", "text-text", "bg-bg/80", "font-bold");
                    b.classList.add("border-transparent", "text-muted");
                });
                btn.classList.add("active", "border-link", "text-text", "bg-bg/80", "font-bold");
                btn.classList.remove("border-transparent", "text-muted");

                panes.forEach(function(p) {
                    if (p.getAttribute("data-tab-pane") === idx) {
                        p.classList.remove("hidden");
                        p.classList.add("active");
                    } else {
                        p.classList.add("hidden");
                        p.classList.remove("active");
                    }
                });
                btn.scrollIntoView({ behavior: "smooth", block: "nearest", inline: "nearest" });
            });
        });

        if (copyBtn) {
            copyBtn.addEventListener("click", function() {
                var activePane = wrap.querySelector(".tab-pane.active code") || wrap.querySelector(".tab-pane:not(.hidden) code") || wrap.querySelector(".tab-pane.active pre") || wrap.querySelector(".tab-pane:not(.hidden) pre");
                if (activePane) {
                    navigator.clipboard.writeText(activePane.textContent || "");
                    var label = copyBtn.querySelector(".copy-label");
                    if (label) {
                        var oldText = label.textContent;
                        label.textContent = "Copied!";
                        setTimeout(function() { label.textContent = oldText; }, 1500);
                    }
                }
            });
        }
    });
})();
