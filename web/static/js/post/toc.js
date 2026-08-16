(function() {
    // ToC: highlight active section
    var tocLinks = document.querySelectorAll(".toc-link");
    if (tocLinks.length > 0) {
        var headings = Array.from(tocLinks).map(function(a) {
            return document.getElementById(a.getAttribute("href").slice(1));
        }).filter(Boolean);
        function setActive() {
            var scrollY = window.scrollY + 100;
            var active = headings[0];
            for (var i = 0; i < headings.length; i++) {
                if (headings[i].offsetTop <= scrollY) active = headings[i];
            }
            tocLinks.forEach(function(a) {
                var isActive = active && a.getAttribute("href") === "#" + active.id;
                a.classList.toggle("toc-active", isActive);
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

    // Mobile ToC chevron
    var details = document.querySelector("details");
    if (details) {
        details.addEventListener("toggle", function() {
            var chevron = details.querySelector(".toc-chevron");
            if (chevron) chevron.style.transform = details.open ? "rotate(180deg)" : "";
        });
    }
})();
