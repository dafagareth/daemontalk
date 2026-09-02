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
        var labelDesktop = document.getElementById("serif-toggle-label");
        var labelModal = document.getElementById("serif-toggle-label-modal");
        if (!el) return;
        if (isSerif) {
            el.classList.add("prose-serif");
            el.classList.remove("prose-sans");
            if (labelDesktop) labelDesktop.textContent = "Serif";
            if (labelModal) labelModal.textContent = "Serif";
        } else {
            el.classList.remove("prose-serif");
            el.classList.add("prose-sans");
            if (labelDesktop) labelDesktop.textContent = "Sans";
            if (labelModal) labelModal.textContent = "Sans";
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

    // 3. Zen / Focus Mode (Shortcut: Z - Desktop only)
    document.addEventListener("keydown", function(e) {
        // Abaikan jika user sedang mengetik di input/textarea
        if (e.target.tagName === "INPUT" || e.target.tagName === "TEXTAREA") return;
        
        if (e.key === 'z' || e.key === 'Z') {
            document.body.classList.toggle('zen-mode');
        }
    });

    // 4. Copy BibTeX Citation
    window.copyBibtex = function(btn) {
        var bibtex = btn.getAttribute('data-bibtex');
        if (!bibtex) return;
        
        navigator.clipboard.writeText(bibtex).then(function() {
            var originalHTML = btn.innerHTML;
            btn.innerHTML = '<span class="font-mono font-bold tracking-widest">[✓]</span><span>Copied</span>';
            btn.classList.add('!text-[var(--c-link)]');
            
            setTimeout(function() {
                btn.innerHTML = originalHTML;
                btn.classList.remove('!text-[var(--c-link)]');
            }, 2000);
        }).catch(function(err) {
            console.error('Failed to copy bibtex: ', err);
        });
    };

    // 5. Mobile Center Reading Modal
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
