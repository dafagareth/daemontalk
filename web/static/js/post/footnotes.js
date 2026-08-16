(function() {
    // Interactive Footnotes Popover Preview
    var activePopover = null;

    function removePopover() {
        if (activePopover) {
            activePopover.remove();
            activePopover = null;
        }
    }

    document.addEventListener("click", function(e) {
        if (activePopover && !activePopover.contains(e.target) && !e.target.closest(".footnote-ref")) {
            removePopover();
        }
    });

    document.querySelectorAll("#prose-body .footnote-ref").forEach(function(ref) {
        var link = ref.tagName === "A" ? ref : ref.querySelector("a") || ref;
        var href = link.getAttribute("href") || "";
        if (!href.startsWith("#fn:") && !href.startsWith("#fnref:")) return;

        var fnId = href.startsWith("#fn:") ? href.slice(1) : href.slice(1);
        var targetFn = document.getElementById(fnId);
        if (!targetFn) return;

        function showPopover(e) {
            removePopover();
            var clone = targetFn.cloneNode(true);
            clone.querySelectorAll(".footnote-backref").forEach(function(b) { b.remove(); });
            var htmlContent = clone.innerHTML.trim();

            var popover = document.createElement("div");
            popover.className = "footnote-popover";
            popover.innerHTML = '<div class="flex items-center justify-between gap-2 pb-1.5 mb-1.5 border-b border-border text-[10px] font-mono font-bold text-muted uppercase tracking-wider"><span>Footnote ' + link.textContent.trim() + '</span><span class="cursor-pointer hover:text-text text-xs" id="fn-pop-close">✕</span></div><div class="text-xs">' + htmlContent + '</div>';

            document.body.appendChild(popover);
            activePopover = popover;

            var rect = ref.getBoundingClientRect();
            var scrollTop = window.pageYOffset || document.documentElement.scrollTop;
            var scrollLeft = window.pageXOffset || document.documentElement.scrollLeft;

            var top = rect.bottom + scrollTop + 6;
            var left = rect.left + scrollLeft;

            if (left + 350 > window.innerWidth) {
                left = Math.max(10, window.innerWidth - 360);
            }

            popover.style.top = top + "px";
            popover.style.left = left + "px";

            var closeBtn = popover.querySelector("#fn-pop-close");
            if (closeBtn) {
                closeBtn.addEventListener("click", function(ev) {
                    ev.stopPropagation();
                    removePopover();
                });
            }
        }

        ref.addEventListener("click", function(e) {
            e.preventDefault();
            if (activePopover && activePopover.dataset.fnId === fnId) {
                removePopover();
                return;
            }
            showPopover(e);
            if (activePopover) activePopover.dataset.fnId = fnId;
        });

        ref.addEventListener("mouseenter", function(e) {
            if (!activePopover) {
                showPopover(e);
                if (activePopover) activePopover.dataset.fnId = fnId;
            }
        });
    });
})();
