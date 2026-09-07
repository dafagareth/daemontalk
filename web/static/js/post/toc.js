(function() {

    var tocLinks = document.querySelectorAll(".toc-link");
    if (tocLinks.length > 0) {
        var headings = [];
        var seenIds = new Set();
        tocLinks.forEach(function(a) {
            var href = a.getAttribute("href");
            if (href && href.startsWith("#")) {
                var id = href.slice(1);
                if (!seenIds.has(id)) {
                    seenIds.add(id);
                    var el = document.getElementById(id);
                    if (el) headings.push(el);
                }
            }
        });

        function setActive() {
            var navOffset = 120;
            var active = null;

            var isAtBottom = (window.innerHeight + window.scrollY) >= (document.documentElement.scrollHeight - 60);
            if (isAtBottom && headings.length > 0) {
                active = headings[headings.length - 1];
            } else {
                for (var i = 0; i < headings.length; i++) {
                    var rect = headings[i].getBoundingClientRect();
                    if (rect.top <= navOffset) {
                        active = headings[i];
                    }
                }
            }

            tocLinks.forEach(function(a) {
                var isActive = active && a.getAttribute("href") === "#" + active.id;
                a.classList.toggle("toc-active", !!isActive);
            });
        }
        tocLinks.forEach(function(a) {
            a.addEventListener("click", function(e) {
                var targetId = a.getAttribute("href").slice(1);
                var targetEl = document.getElementById(targetId);
                if (targetEl) {
                    e.preventDefault();
                    var navOffset = 80;
                    var top = targetEl.getBoundingClientRect().top + window.pageYOffset - navOffset;
                    window.scrollTo({ top: Math.max(0, top), behavior: "smooth" });
                    history.pushState(null, "", "#" + targetId);
                }
            });
        });

        window.addEventListener("scroll", setActive, { passive: true });
        setActive();
    }

    var details = document.querySelector("details");
    if (details) {
        details.addEventListener("toggle", function() {
            var chevron = details.querySelector(".toc-chevron");
            if (chevron) chevron.style.transform = details.open ? "rotate(180deg)" : "";
        });
    }
})();
