(function() {
    // Clipboard helper: falls back to execCommand for HTTP (non-secure) contexts
    window.copyText = function(text, done) {
        if (navigator.clipboard && navigator.clipboard.writeText) {
            navigator.clipboard.writeText(text).then(done).catch(function() { fallback(text, done); });
        } else {
            fallback(text, done);
        }
    };
    function fallback(text, done) {
        var ta = document.createElement('textarea');
        ta.value = text;
        ta.style.cssText = 'position:fixed;top:-9999px;left:-9999px;opacity:0';
        document.body.appendChild(ta);
        ta.focus(); ta.select();
        try { document.execCommand('copy'); if (done) done(); } catch(e) {}
        document.body.removeChild(ta);
    }
})();
// Bookmark toggle (supports SVG icons & text state)
(function() {
    function updateBookmarkBtnUI(btn, isSaved) {
        var icon = btn.querySelector('.bookmark-icon');
        var text = btn.querySelector('.bookmark-text');
        if (isSaved) {
            btn.classList.add("bg-accent/15", "text-accent", "border-accent");
            if (icon) {
                icon.setAttribute("fill", "currentColor");
            }
            if (text) {
                var lang = document.documentElement.lang || "id";
                text.textContent = lang === "id" ? "Tersimpan" : "Saved";
            }
        } else {
            btn.classList.remove("bg-accent/15", "text-accent", "border-accent");
            if (icon) {
                icon.setAttribute("fill", "none");
            }
            if (text) {
                var lang = document.documentElement.lang || "id";
                text.textContent = lang === "id" ? "Simpan" : "Save";
            }
        }
        if (!icon && !text) {
            btn.textContent = isSaved ? "★" : "☆";
            if (isSaved) btn.classList.add("text-accent");
            else btn.classList.remove("text-accent");
        }
    }

    var bookmarks = [];
    try { bookmarks = JSON.parse(localStorage.getItem("bookmarks") || "[]"); } catch(e) {}
    var saved = {};
    bookmarks.forEach(function(b) { saved[b.slug] = true; });
    document.querySelectorAll(".bookmark-btn").forEach(function(btn) {
        if (saved[btn.dataset.slug]) {
            updateBookmarkBtnUI(btn, true);
        }
    });

    window.toggleBookmark = function(btn) {
        var slug = btn.dataset.slug;
        var title = btn.dataset.title || "";
        var date = btn.dataset.date || "";
        var bookmarks = [];
        try { bookmarks = JSON.parse(localStorage.getItem("bookmarks") || "[]"); } catch(e) {}
        var idx = bookmarks.findIndex(function(b) { return b.slug === slug; });
        var isNowSaved = false;
        if (idx >= 0) {
            bookmarks.splice(idx, 1);
            isNowSaved = false;
        } else {
            bookmarks.unshift({slug: slug, title: title, date: date});
            isNowSaved = true;
        }
        localStorage.setItem("bookmarks", JSON.stringify(bookmarks));

        // Update all buttons for this slug
        document.querySelectorAll('.bookmark-btn[data-slug="' + slug + '"]').forEach(function(el) {
            updateBookmarkBtnUI(el, isNowSaved);
        });
    };
})();
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
(function() {
    // Enhanced Code block wrapper: adds language badges and copy button.
    document.querySelectorAll("#prose-body pre").forEach(function(pre) {
        if (pre.closest(".code-tabs-wrap") || pre.closest(".code-output-wrap")) return;
        var wrap = document.createElement("div");
        wrap.className = "code-wrap";
        pre.parentNode.insertBefore(wrap, pre);
        wrap.appendChild(pre);

        var codeEl = pre.querySelector("code");
        var codeText = codeEl ? codeEl.innerText : pre.innerText;

        // Detect language from class or content
        var lang = "";
        var classNames = (pre.className + " " + (codeEl ? codeEl.className : "")).toLowerCase();
        var langMatch = classNames.match(/language-([a-z0-9_-]+)/) || classNames.match(/lang-([a-z0-9_-]+)/);
        if (langMatch) {
            lang = langMatch[1];
        } else if (codeText.indexOf("package main") >= 0 || codeText.indexOf("fmt.Print") >= 0 || codeText.indexOf("func ") >= 0) {
            lang = "go";
        } else if (codeText.indexOf("def ") >= 0 || (codeText.indexOf("import ") >= 0 && codeText.indexOf("from ") >= 0)) {
            lang = "python";
        } else if (codeText.indexOf("fn ") >= 0 || codeText.indexOf("let mut ") >= 0) {
            lang = "rust";
        } else if (codeText.indexOf("const ") >= 0 || codeText.indexOf("console.log") >= 0) {
            lang = "javascript";
        } else if (codeText.startsWith("$ ") || codeText.startsWith("sudo ") || codeText.startsWith("curl ") || codeText.startsWith("docker ")) {
            lang = "bash";
        }

        // Detect diagram / architecture schematics
        var isDiagram = codeText.indexOf("┌") >= 0 || codeText.indexOf("┼") >= 0 || codeText.indexOf("-->") >= 0 || codeText.indexOf("──►") >= 0 || codeText.indexOf("flowchart") >= 0 || codeText.indexOf("graph TD") >= 0 || lang === "mermaid" || lang === "diagram";
        if (isDiagram && !lang) {
            lang = "diagram";
        }

        // Create action toolbar
        var toolbar = document.createElement("div");
        toolbar.className = "code-toolbar";

        // Language or Diagram badge
        if (lang) {
            var badge = document.createElement("span");
            badge.className = "code-lang-badge" + (isDiagram ? " text-link font-bold" : "");
            badge.textContent = isDiagram ? "ARCHITECTURE DIAGRAM" : lang.toUpperCase();
            toolbar.appendChild(badge);
        }

        // Interactive Diagram Zoom/Focus Toggle
        if (isDiagram) {
            var zoomBtn = document.createElement("button");
            zoomBtn.textContent = "expand";
            zoomBtn.className = "copy-btn";
            zoomBtn.setAttribute("aria-label", "Expand diagram");
            zoomBtn.addEventListener("click", function() {
                wrap.classList.toggle("diagram-expanded");
                if (wrap.classList.contains("diagram-expanded")) {
                    zoomBtn.textContent = "collapse";
                    pre.style.maxHeight = "none";
                    pre.style.fontSize = "13px";
                } else {
                    zoomBtn.textContent = "expand";
                    pre.style.maxHeight = "";
                    pre.style.fontSize = "";
                }
            });
            toolbar.appendChild(zoomBtn);
        }

        // Copy button
        var copyBtn = document.createElement("button");
        copyBtn.textContent = "copy";
        copyBtn.className = "copy-btn";
        copyBtn.setAttribute("aria-label", "Copy code");
        copyBtn.addEventListener("click", function() {
            var text = codeEl ? codeEl.innerText : pre.innerText;
            if (window.copyText) {
                window.copyText(text, function() {
                    copyBtn.textContent = "copied!";
                    setTimeout(function() { copyBtn.textContent = "copy"; }, 2000);
                });
            } else if (navigator.clipboard) {
                navigator.clipboard.writeText(text).then(function() {
                    copyBtn.textContent = "copied!";
                    setTimeout(function() { copyBtn.textContent = "copy"; }, 2000);
                });
            }
        });
        toolbar.appendChild(copyBtn);

        wrap.appendChild(toolbar);
    });

    // Code Line & Diff Highlighting Post-processor
    document.querySelectorAll("#prose-body pre code").forEach(function(block) {
        if (block.dataset.hlProcessed) return;
        block.dataset.hlProcessed = "true";

        var html = block.innerHTML;
        if (!html.includes("[!code") && !html.includes("&gt; [!code") && !block.closest(".language-diff")) return;

        var lines = html.split("\n");
        var newLines = lines.map(function(line) {
            if (line.includes("[!code ++]") || (line.startsWith("+") && block.closest(".language-diff"))) {
                var clean = line.replace(/(\/\/|\/\*|#|--)\s*\[\!code \+\+\]\s*(\*\/)?/, "");
                return '<span class="line-add">' + clean + '</span>';
            }
            if (line.includes("[!code --]") || (line.startsWith("-") && block.closest(".language-diff"))) {
                var clean = line.replace(/(\/\/|\/\*|#|--)\s*\[\!code \-\-\]\s*(\*\/)?/, "");
                return '<span class="line-del">' + clean + '</span>';
            }
            if (line.includes("[!code hl]") || line.includes("[!code highlight]")) {
                var clean = line.replace(/(\/\/|\/\*|#|--)\s*\[\!code (hl|highlight)\]\s*(\*\/)?/, "");
                return '<span class="line-hl">' + clean + '</span>';
            }
            return line;
        });
        block.innerHTML = newLines.join("\n");
    });
})();
(function() {
    // Toggle Share Popover Menu
    window.toggleSharePopover = function(e) {
        if (e) e.stopPropagation();
        var menu = document.getElementById('share-popover-menu');
        if (menu) {
            menu.classList.toggle('hidden');
        }
    };

    // Close share popover when clicking anywhere outside
    document.addEventListener('click', function(e) {
        var menu = document.getElementById('share-popover-menu');
        if (menu && !menu.contains(e.target) && !e.target.closest('[onclick*="toggleSharePopover"]')) {
            menu.classList.add('hidden');
        }
    });

    function getPostInfo() {
        var title = (document.querySelector("h1") || {}).textContent || "DaemonTalk Article";
        var url = window.location.href;
        return { title: title.trim(), url: url };
    }

    // Copy standard URL with badge feedback
    window.copyPostUrl = function(btn, e) {
        if (e) e.stopPropagation();
        var info = getPostInfo();
        performCopy(info.url, btn);
    };

    // Copy Markdown formatted link: [Title](URL)
    window.copyMarkdownLink = function(btn, e) {
        if (e) e.stopPropagation();
        var info = getPostInfo();
        var md = '[' + info.title + '](' + info.url + ')';
        performCopy(md, btn);
    };

    function performCopy(text, btn) {
        var copySuccess = function() {
            var label = btn.querySelector('.share-label-text');
            var status = btn.querySelector('.share-copied-status');
            if (label && status) {
                label.classList.add('hidden');
                status.classList.remove('hidden');
                setTimeout(function() {
                    label.classList.remove('hidden');
                    status.classList.add('hidden');
                }, 2000);
            }
        };

        if (navigator.clipboard && navigator.clipboard.writeText) {
            navigator.clipboard.writeText(text).then(copySuccess).catch(function() {
                fallbackCopy(text, copySuccess);
            });
        } else {
            fallbackCopy(text, copySuccess);
        }
    }

    function fallbackCopy(text, cb) {
        var ta = document.createElement('textarea');
        ta.value = text;
        ta.style.position = 'fixed';
        ta.style.opacity = '0';
        document.body.appendChild(ta);
        ta.select();
        try {
            document.execCommand('copy');
            if (cb) cb();
        } catch (err) {}
        document.body.removeChild(ta);
    }

    // Social share triggers
    window.shareToBluesky = function(e) {
        if (e) e.stopPropagation();
        var info = getPostInfo();
        window.open('https://bsky.app/intent/compose?text=' + encodeURIComponent(info.title + '\n' + info.url), '_blank', 'width=600,height=450');
    };

    window.shareToTwitter = function(e) {
        if (e) e.stopPropagation();
        var info = getPostInfo();
        window.open('https://twitter.com/intent/tweet?text=' + encodeURIComponent(info.title) + '&url=' + encodeURIComponent(info.url), '_blank', 'width=600,height=450');
    };

    window.shareToLinkedIn = function(e) {
        if (e) e.stopPropagation();
        var info = getPostInfo();
        window.open('https://www.linkedin.com/sharing/share-offsite/?url=' + encodeURIComponent(info.url), '_blank', 'width=600,height=500');
    };

    window.shareToThreads = function(e) {
        if (e) e.stopPropagation();
        var info = getPostInfo();
        window.open('https://www.threads.net/intent/post?text=' + encodeURIComponent(info.title + ' ' + info.url), '_blank', 'width=600,height=500');
    };
})();
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
(function() {
    // Interactive task lists / checklists with local persistence
    var path = window.location.pathname;
    var match = path.match(/\/blog\/([^\/]+)$/);
    var postSlug = match ? match[1] : (window.location.pathname.split("/").filter(Boolean).pop() || "default");
    var storageKey = "dt_checklist_" + postSlug;
    var state = {};
    try {
        state = JSON.parse(localStorage.getItem(storageKey) || "{}");
    } catch(e) {}

    var checkboxes = document.querySelectorAll("#prose-body input[type='checkbox']");
    checkboxes.forEach(function(cb, idx) {
        cb.disabled = false;
        cb.removeAttribute("disabled");
        cb.style.cursor = "pointer";
        cb.style.pointerEvents = "auto";
        cb.style.touchAction = "manipulation";

        var parentLi = cb.closest("li");
        if (parentLi) {
            parentLi.style.listStyle = "none";
            var parentUl = parentLi.parentElement;
            if (parentUl) parentUl.style.listStyle = "none";
        }

        // Restore saved state
        if (state[idx]) {
            cb.checked = true;
            if (parentLi) parentLi.classList.add("task-checked");
        }

        cb.addEventListener("change", function() {
            state[idx] = cb.checked;
            if (parentLi) {
                if (cb.checked) parentLi.classList.add("task-checked");
                else parentLi.classList.remove("task-checked");
            }
            try {
                localStorage.setItem(storageKey, JSON.stringify(state));
            } catch(e) {}
        });
    });
})();
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

        wrap.querySelectorAll(".copy-tab-code").forEach(function(copyBtn) {
            copyBtn.addEventListener("click", function() {
                var pane = copyBtn.closest(".tab-pane");
                var codeEl = pane ? (pane.querySelector("code") || pane.querySelector("pre")) : null;
                if (codeEl) {
                    var text = codeEl.innerText || codeEl.textContent || "";
                    if (window.copyText) {
                        window.copyText(text, function() {
                            var oldText = copyBtn.textContent;
                            copyBtn.textContent = "copied!";
                            setTimeout(function() { copyBtn.textContent = oldText; }, 2000);
                        });
                    } else if (navigator.clipboard) {
                        navigator.clipboard.writeText(text).then(function() {
                            var oldText = copyBtn.textContent;
                            copyBtn.textContent = "copied!";
                            setTimeout(function() { copyBtn.textContent = oldText; }, 2000);
                        });
                    }
                }
            });
        });
    });
})();
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
(function() {
    // Image Lightbox / Full Photo Viewer for Articles & Carousels
    var lightbox = null;
    var imgEl = null;
    var captionEl = null;
    var counterEl = null;
    var bottomCountEl = null;
    var bottomBarEl = null;
    var prevBtn = null;
    var nextBtn = null;

    var galleryItems = [];
    var currentIndex = 0;

    function createLightbox() {
        if (lightbox) return;

        lightbox = document.createElement("div");
        lightbox.className = "lightbox-overlay";
        lightbox.innerHTML = `
            <div class="lightbox-backdrop"></div>
            <div class="lightbox-toolbar">
                <div class="lightbox-caption"></div>
                <div class="lightbox-actions">
                    <span class="lightbox-counter"></span>
                    <button type="button" class="lightbox-close" aria-label="Close photo viewer">
                        <svg class="w-5 h-5" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                            <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" />
                        </svg>
                    </button>
                </div>
            </div>

            <button type="button" class="lightbox-nav-btn lightbox-prev" aria-label="Previous photo">
                <svg class="w-6 h-6" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M15 19l-7-7 7-7" />
                </svg>
            </button>

            <div class="lightbox-content">
                <img class="lightbox-img" src="" alt="" />
            </div>

            <button type="button" class="lightbox-nav-btn lightbox-next" aria-label="Next photo">
                <svg class="w-6 h-6" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M9 5l7 7-7 7" />
                </svg>
            </button>

            <div class="lightbox-bottom-bar">
                <button type="button" class="lightbox-bottom-btn lightbox-bottom-prev" aria-label="Previous slide">
                    <svg class="w-4 h-4" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                        <path stroke-linecap="round" stroke-linejoin="round" d="M15 19l-7-7 7-7" />
                    </svg>
                    <span>PREV</span>
                </button>
                <span class="lightbox-bottom-count"></span>
                <button type="button" class="lightbox-bottom-btn lightbox-bottom-next" aria-label="Next slide">
                    <span>NEXT</span>
                    <svg class="w-4 h-4" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                        <path stroke-linecap="round" stroke-linejoin="round" d="M9 5l7 7-7 7" />
                    </svg>
                </button>
            </div>
        `;
        document.body.appendChild(lightbox);

        imgEl = lightbox.querySelector(".lightbox-img");
        captionEl = lightbox.querySelector(".lightbox-caption");
        counterEl = lightbox.querySelector(".lightbox-counter");
        bottomCountEl = lightbox.querySelector(".lightbox-bottom-count");
        bottomBarEl = lightbox.querySelector(".lightbox-bottom-bar");

        prevBtn = lightbox.querySelector(".lightbox-prev");
        nextBtn = lightbox.querySelector(".lightbox-next");
        var bPrevBtn = lightbox.querySelector(".lightbox-bottom-prev");
        var bNextBtn = lightbox.querySelector(".lightbox-bottom-next");

        // Click event listeners
        lightbox.querySelector(".lightbox-close").addEventListener("click", closeLightbox);
        lightbox.querySelector(".lightbox-backdrop").addEventListener("click", closeLightbox);

        prevBtn.addEventListener("click", prevPhoto);
        nextBtn.addEventListener("click", nextPhoto);
        bPrevBtn.addEventListener("click", prevPhoto);
        bNextBtn.addEventListener("click", nextPhoto);

        lightbox.querySelector(".lightbox-content").addEventListener("click", function(e) {
            if (e.target === this) {
                closeLightbox();
            }
        });

        // Touch swipe support on mobile
        var touchStartX = 0;
        var touchStartY = 0;
        var contentEl = lightbox.querySelector(".lightbox-content");

        contentEl.addEventListener("touchstart", function(e) {
            touchStartX = e.touches[0].clientX;
            touchStartY = e.touches[0].clientY;
        }, { passive: true });

        contentEl.addEventListener("touchend", function(e) {
            if (!touchStartX || !touchStartY) return;
            var diffX = e.changedTouches[0].clientX - touchStartX;
            var diffY = e.changedTouches[0].clientY - touchStartY;

            if (Math.abs(diffX) > 40 && Math.abs(diffX) > Math.abs(diffY)) {
                if (diffX > 0) {
                    prevPhoto();
                } else {
                    nextPhoto();
                }
            } else if (Math.abs(diffY) > 80 && Math.abs(diffY) > Math.abs(diffX)) {
                closeLightbox();
            }

            touchStartX = 0;
            touchStartY = 0;
        }, { passive: true });

        // Keyboard controls
        document.addEventListener("keydown", function(e) {
            if (!lightbox || !lightbox.classList.contains("open")) return;
            if (e.key === "Escape") {
                closeLightbox();
            } else if (e.key === "ArrowLeft" || e.key === "a" || e.key === "A") {
                prevPhoto();
            } else if (e.key === "ArrowRight" || e.key === "d" || e.key === "D") {
                nextPhoto();
            }
        });
    }

    function renderActivePhoto() {
        if (!galleryItems.length || currentIndex < 0 || currentIndex >= galleryItems.length) return;

        var item = galleryItems[currentIndex];
        imgEl.style.opacity = "0.4";
        imgEl.style.transform = "scale(0.96)";

        setTimeout(function() {
            imgEl.src = item.src;
            imgEl.alt = item.alt || "";

            if (item.caption) {
                captionEl.textContent = item.caption;
                captionEl.style.display = "block";
            } else {
                captionEl.textContent = "";
                captionEl.style.display = "none";
            }

            var showNav = galleryItems.length > 1;
            var countStr = showNav ? ((currentIndex + 1) + " / " + galleryItems.length) : "";
            counterEl.textContent = countStr;
            counterEl.style.display = showNav ? "block" : "none";
            bottomCountEl.textContent = countStr;

            prevBtn.style.display = showNav ? "flex" : "none";
            nextBtn.style.display = showNav ? "flex" : "none";
            bottomBarEl.style.display = showNav ? "flex" : "none";

            imgEl.style.opacity = "1";
            imgEl.style.transform = "scale(1)";
        }, 80);
    }

    function nextPhoto() {
        if (galleryItems.length <= 1) return;
        currentIndex = (currentIndex + 1) % galleryItems.length;
        renderActivePhoto();
    }

    function prevPhoto() {
        if (galleryItems.length <= 1) return;
        currentIndex = (currentIndex - 1 + galleryItems.length) % galleryItems.length;
        renderActivePhoto();
    }

    function openGallery(items, startIdx) {
        createLightbox();
        galleryItems = items;
        currentIndex = startIdx || 0;

        renderActivePhoto();
        lightbox.classList.add("open");
        document.body.style.overflow = "hidden";
    }

    function closeLightbox() {
        if (!lightbox) return;
        lightbox.classList.remove("open");
        document.body.style.overflow = "";
        setTimeout(function() {
            if (imgEl && !lightbox.classList.contains("open")) {
                imgEl.src = "";
            }
        }, 250);
    }

    function initLightbox() {
        // Collect all images in articles, covers, galleries, and carousels
        var containerSelectors = [
            ".post-carousel-wrap",
            ".post-gallery-wrap",
            "figure",
            ".post-cover",
            "#prose-body"
        ];

        // Also bind standalone cover & content images
        var images = document.querySelectorAll("#prose-body img, .post-cover img, figure img, .post-gallery-wrap img, .post-carousel-wrap img, article header img");

        images.forEach(function(img) {
            if (img.dataset.lightboxBound) return;
            img.dataset.lightboxBound = "true";
            img.classList.add("lightbox-trigger");

            img.addEventListener("click", function(e) {
                var parentLink = img.closest("a");
                if (parentLink && parentLink.href && !parentLink.href.match(/\.(jpg|jpeg|png|webp|gif|svg)(\?.*)?$/i)) {
                    return; // Standard web links navigate normally
                }

                e.preventDefault();

                // Determine gallery scope: ONLY carousel or gallery containers
                var carouselScope = img.closest(".post-carousel-wrap") || img.closest(".post-gallery-wrap");
                var scopeImages = [];

                if (carouselScope) {
                    scopeImages = Array.from(carouselScope.querySelectorAll("img")).filter(function(i) {
                        return i.src && !i.closest(".post-author-card") && !i.closest(".no-lightbox");
                    });
                }

                if (!scopeImages.length) {
                    scopeImages = [img];
                }

                var items = scopeImages.map(function(i) {
                    var cap = "";
                    var fig = i.closest("figure");
                    if (fig) {
                        var fc = fig.querySelector("figcaption");
                        if (fc) cap = fc.innerText;
                    }
                    if (!cap && i.alt) cap = i.alt;
                    if (!cap && i.title) cap = i.title;

                    var pLink = i.closest("a");
                    var src = (pLink && pLink.href && pLink.href.match(/\.(jpg|jpeg|png|webp|gif|svg)(\?.*)?$/i)) ? pLink.href : i.src;
                    return { src: src, alt: i.alt || "", caption: cap };
                });

                var activeIndex = scopeImages.indexOf(img);
                if (activeIndex < 0) activeIndex = 0;

                openGallery(items, activeIndex);
            });
        });
    }

    if (document.readyState === "loading") {
        document.addEventListener("DOMContentLoaded", initLightbox);
    } else {
        initLightbox();
    }

    document.addEventListener("htmx:afterSwap", initLightbox);
})();
document.addEventListener('DOMContentLoaded', () => {
    const wikiLinks = document.querySelectorAll('.prose a[href*="wikipedia.org/wiki/"]');
    if (wikiLinks.length === 0) return;

    const tooltip = document.createElement('div');
    tooltip.className = 'wiki-tooltip';
    document.body.appendChild(tooltip);

    let hoverTimeout;
    let currentLink = null;

    wikiLinks.forEach(link => {
        link.addEventListener('mouseenter', () => {
            clearTimeout(hoverTimeout);
            currentLink = link;

            const url = new URL(link.href);
            const langMatch = url.hostname.match(/^([a-z\-]+)\.wikipedia\.org/);
            const lang = langMatch ? langMatch[1] : 'en';
            const title = url.pathname.split('/wiki/')[1];
            
            if (!title) return;

            hoverTimeout = setTimeout(() => {
                showLoading(link);
                fetchWikiData(lang, title, link);
            }, 300);
        });

        link.addEventListener('mouseleave', () => {
            clearTimeout(hoverTimeout);
            if (currentLink === link) {
                hoverTimeout = setTimeout(hideTooltip, 300);
            }
        });
    });

    tooltip.addEventListener('mouseenter', () => clearTimeout(hoverTimeout));
    tooltip.addEventListener('mouseleave', () => {
        hoverTimeout = setTimeout(hideTooltip, 300);
    });

    function showLoading(targetEl) {
        positionTooltip(targetEl);
        tooltip.innerHTML = '<div class="wiki-loading">Memuat dari Wikipedia...</div>';
        tooltip.classList.add('visible');
    }

    function hideTooltip() {
        tooltip.classList.remove('visible');
    }

    function positionTooltip(targetEl) {
        const rect = targetEl.getBoundingClientRect();
        const tooltipWidth = 320;
        
        let top = rect.bottom + window.scrollY + 8;
        let left = rect.left + window.scrollX - (tooltipWidth / 2) + (rect.width / 2);

        if (left < 16) left = 16;
        if (left + tooltipWidth > window.innerWidth - 16) {
            left = window.innerWidth - tooltipWidth - 16;
        }

        tooltip.style.top = `${top}px`;
        tooltip.style.left = `${left}px`;
    }

    async function fetchWikiData(lang, title, targetEl) {
        try {
            const safeTitle = encodeURIComponent(decodeURIComponent(title));
            const apiUrl = `https://${lang}.wikipedia.org/api/rest_v1/page/summary/${safeTitle}?redirect=true`;
            
            // Mengirim request polos (tanpa header khusus) agar menjadi CORS "Simple Request".
            // Jika ada header tambahan, browser akan mengirim preflight (OPTIONS).
            // Preflight yang berujung pada Redirect (301) dari Wikipedia akan diblokir browser (NetworkError).
            const res = await fetch(apiUrl);
            
            if (!res.ok) throw new Error(`HTTP ${res.status}`);
            const data = await res.json();
            
            if (currentLink !== targetEl || !tooltip.classList.contains('visible')) return;

            renderTooltip(data);
            positionTooltip(targetEl);
        } catch (e) {
            if (currentLink === targetEl) {
                tooltip.innerHTML = `<div class="wiki-error">Gagal: ${e.message}</div>`;
            }
        }
    }

    function renderTooltip(data) {
        let html = '';
        if (data.thumbnail && data.thumbnail.source) {
            html += `<img src="${data.thumbnail.source}" class="wiki-thumb" alt="">`;
        }
        html += `<div class="wiki-content">`;
        html += `<div class="wiki-title">${data.title}</div>`;
        
        let extract = data.extract;
        if (extract.length > 200) extract = extract.substring(0, 200) + '...';
        
        html += `<p class="wiki-extract">${extract}</p>`;
        html += `<div class="wiki-footer">W — Disediakan oleh Wikipedia</div>`;
        html += `</div>`;
        
        tooltip.innerHTML = html;
    }
});
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
