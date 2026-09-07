(function() {

    var SCALES = [
        { mobile: 15,   desktop: 17.5 },
        { mobile: 17,   desktop: 19.5 },
        { mobile: 19.5, desktop: 22   }
    ];
    var DEFAULT_STEP = 1;

    function applySizeStep(step) {
        var idx = Math.max(0, Math.min(SCALES.length - 1, parseInt(step) || 0));
        var scale = SCALES[idx];

        document.documentElement.style.setProperty("--prose-size-mobile", scale.mobile + "px");
        document.documentElement.style.setProperty("--prose-size-desktop", scale.desktop + "px");

        var el = document.getElementById("prose-body");
        if (el) {
            el.style.fontSize = "";
        }

        var sliders = document.querySelectorAll(".text-size-slider");
        sliders.forEach(function(s) {
            s.value = idx;
        });

        var stepLabels = document.querySelectorAll("[data-size-step]");
        stepLabels.forEach(function(l) {
            var sVal = parseInt(l.getAttribute("data-size-step"));
            if (sVal === idx) {
                l.classList.add("text-text", "font-bold");
                l.classList.remove("text-muted");
            } else {
                l.classList.remove("text-text", "font-bold");
                l.classList.add("text-muted");
            }
        });
    }

    window.setTextSizeIndex = function(step) {
        var idx = Math.max(0, Math.min(SCALES.length - 1, parseInt(step) || 0));
        localStorage.setItem("prose-size-step", idx);
        applySizeStep(idx);
    };

    window.adjustProseSize = function(dir) {
        var current = parseInt(localStorage.getItem("prose-size-step"));
        if (isNaN(current)) current = DEFAULT_STEP;
        window.setTextSizeIndex(current + dir);
    };

    var savedStep = parseInt(localStorage.getItem("prose-size-step"));
    if (isNaN(savedStep) || savedStep < 0 || savedStep >= SCALES.length) {
        savedStep = DEFAULT_STEP;
    }
    applySizeStep(savedStep);

    window.toggleReadingPopover = function(e) {
        if (e) e.stopPropagation();
        document.querySelectorAll('[id*="share-popover-menu"], [id*="more-popover-menu"]').forEach(function(m) {
            m.classList.add('hidden');
        });
        var pop = document.getElementById('reading-popover-menu');
        if (pop) {
            pop.classList.toggle('hidden');
        }
    };

    document.addEventListener('click', function(e) {
        var pop = document.getElementById('reading-popover-menu');
        if (pop && !pop.classList.contains('hidden')) {
            if (!pop.contains(e.target) && !e.target.closest('[onclick*="toggleReadingPopover"]')) {
                pop.classList.add('hidden');
            }
        }
    });

    window.toggleMorePopover = function(e, menuId) {
        if (e) e.stopPropagation();
        document.querySelectorAll('[id*="share-popover-menu"], #reading-popover-menu').forEach(function(m) {
            m.classList.add('hidden');
        });
        var id = menuId || 'header-more-popover-menu';
        var pop = document.getElementById(id);
        if (pop) {
            var isClosed = pop.classList.contains('hidden');
            document.querySelectorAll('[id*="more-popover-menu"]').forEach(function(p) {
                p.classList.add('hidden');
            });
            if (isClosed) {
                pop.classList.remove('hidden');
            }
        }
    };

    document.addEventListener('click', function(e) {
        document.querySelectorAll('[id*="more-popover-menu"]').forEach(function(pop) {
            if (!pop.contains(e.target) && !e.target.closest('[onclick*="toggleMorePopover"]')) {
                pop.classList.add('hidden');
            }
        });
    });

    window.copyBibtex = function(btn) {
        var bibtex = btn.getAttribute('data-bibtex');
        if (!bibtex) return;

        navigator.clipboard.writeText(bibtex).then(function() {
            var status = btn.querySelector('.bibtex-copied-status');
            if (status) {
                status.classList.remove('hidden');
                setTimeout(function() {
                    status.classList.add('hidden');
                }, 2000);
            } else {
                var originalHTML = btn.innerHTML;
                btn.innerHTML = '<span class="font-serif font-bold text-sm text-accent">✓</span>';
                btn.classList.add('text-accent');

                setTimeout(function() {
                    btn.innerHTML = originalHTML;
                    btn.classList.remove('text-accent');
                }, 2000);
            }
        }).catch(function(err) {
            console.error('Failed to copy bibtex: ', err);
        });
    };

    window.toggleMobileReadingModal = function() {
        var modal = document.getElementById("mobile-reading-modal-overlay");
        if (!modal) return;
        if (modal.classList.contains("hidden")) {
            modal.classList.remove("hidden");
            document.body.style.overflow = "hidden";
        } else {
            modal.classList.add("hidden");
            document.body.style.overflow = "";
        }
    };

    document.addEventListener("keydown", function(e) {
        if (e.key === "Escape") {
            var modal = document.getElementById("mobile-reading-modal-overlay");
            if (modal && !modal.classList.contains("hidden")) {
                modal.classList.add("hidden");
                document.body.style.overflow = "";
            }
        }
    });
})();
