(function() {

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
(function() {
    function updateBookmarkBtnUI(btn, isSaved) {
        var icon = btn.querySelector('.bookmark-icon');
        var text = btn.querySelector('.bookmark-text');
        if (isSaved) {
            btn.classList.add("is-saved", "bg-hover", "text-text", "font-bold");
            btn.classList.remove("border-text/70", "bg-accent/15", "text-accent", "border-accent");
            if (icon) {
                icon.setAttribute("fill", "currentColor");
            }
            if (text) {
                var lang = document.documentElement.lang || "id";
                text.textContent = lang === "id" ? "Tersimpan" : "Saved";
            }
        } else {
            btn.classList.remove("is-saved", "bg-hover", "text-text", "border-text/70", "font-bold", "bg-accent/15", "text-accent", "border-accent");
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
            if (isSaved) btn.classList.add("text-text");
            else btn.classList.remove("text-text", "text-accent");
        }
    }

    function syncPostBookmarks() {
        var bookmarks = [];
        try { bookmarks = JSON.parse(localStorage.getItem("bookmarks") || "[]"); } catch(e) {}
        var saved = {};
        bookmarks.forEach(function(b) { saved[b.slug] = true; });
        document.querySelectorAll(".bookmark-btn").forEach(function(btn) {
            updateBookmarkBtnUI(btn, !!saved[btn.dataset.slug]);
        });
    }
    syncPostBookmarks();
    document.addEventListener("htmx:afterSwap", syncPostBookmarks);

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

        document.querySelectorAll('.bookmark-btn[data-slug="' + slug + '"]').forEach(function(el) {
            updateBookmarkBtnUI(el, isNowSaved);
        });
    };
})();
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
(function() {

    document.querySelectorAll("#prose-body pre").forEach(function(pre) {
        if (pre.closest(".code-tabs-wrap") || pre.closest(".code-output-wrap")) return;
        var wrap = document.createElement("div");
        wrap.className = "code-wrap";
        pre.parentNode.insertBefore(wrap, pre);
        wrap.appendChild(pre);

        var codeEl = pre.querySelector("code");
        var codeText = codeEl ? codeEl.innerText : pre.innerText;

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

        var isDiagram = codeText.indexOf("┌") >= 0 || codeText.indexOf("┼") >= 0 || codeText.indexOf("-->") >= 0 || codeText.indexOf("──►") >= 0 || codeText.indexOf("flowchart") >= 0 || codeText.indexOf("graph TD") >= 0 || lang === "mermaid" || lang === "diagram";
        if (isDiagram && !lang) {
            lang = "diagram";
        }

        var toolbar = document.createElement("div");
        toolbar.className = "code-toolbar";

        if (lang) {
            var badge = document.createElement("span");
            badge.className = "code-lang-badge" + (isDiagram ? " text-link font-bold" : "");
            badge.textContent = isDiagram ? "ARCHITECTURE DIAGRAM" : lang.toUpperCase();
            toolbar.appendChild(badge);
        }

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
    function isMobileDevice() {
        return (
            /Android|webOS|iPhone|iPad|iPod|BlackBerry|IEMobile|Opera Mini/i.test(navigator.userAgent) ||
            (window.matchMedia && window.matchMedia("(max-width: 768px)").matches) ||
            ('ontouchstart' in window && window.innerWidth <= 1024)
        );
    }

    function showShareMenu(menuId) {
        document.querySelectorAll('[id*="more-popover-menu"], #reading-popover-menu').forEach(function(m) {
            m.classList.add('hidden');
        });
        var id = menuId || 'share-popover-menu';
        var menu = document.getElementById(id);
        if (menu) {
            var isClosed = menu.classList.contains('hidden');
            document.querySelectorAll('[id*="share-popover-menu"]').forEach(function(m) {
                m.classList.add('hidden');
            });
            if (isClosed) {
                menu.classList.remove('hidden');
            }
        }
    }

    window.toggleSharePopover = function(e, menuId) {
        if (e) e.stopPropagation();

        var info = getPostInfo();
        if (navigator.share && isMobileDevice()) {
            var shareData = {
                title: info.title,
                url: info.url
            };
            if (info.text) {
                shareData.text = info.text;
            }
            navigator.share(shareData).catch(function(err) {
                if (err && err.name !== 'AbortError') {
                    showShareMenu(menuId);
                }
            });
            return;
        }

        showShareMenu(menuId);
    };

    document.addEventListener('click', function(e) {
        document.querySelectorAll('[id*="share-popover-menu"]').forEach(function(menu) {
            if (!menu.contains(e.target) && !e.target.closest('[onclick*="toggleSharePopover"]')) {
                menu.classList.add('hidden');
            }
        });
    });

    function getPostInfo() {
        var title = (document.querySelector("h1") || {}).textContent || "DaemonTalk Article";
        var descMeta = document.querySelector('meta[name="description"]');
        var text = descMeta ? descMeta.getAttribute("content") : "";
        var url = window.location.href;
        return { title: title.trim(), text: text, url: url };
    }

    window.copyPostUrl = function(btn, e) {
        if (e) e.stopPropagation();
        var info = getPostInfo();
        performCopy(info.url, btn);
    };

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
(function() {

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

    document.querySelectorAll("[data-code-tabs]").forEach(function(wrap) {
        var navTrack = wrap.querySelector(".tabs-nav-track");
        var buttons = wrap.querySelectorAll(".tab-btn");
        var panes = wrap.querySelectorAll(".tab-pane");
        var copyBtn = wrap.querySelector(".copy-tab-code");

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
(function() {
    var path = window.location.pathname;
    var match = path.match(/\/blog\/([^\/]+)$/);
    if (match) {
        var slug = match[1];
        var read = [];
        try { read = JSON.parse(localStorage.getItem('readPosts') || '[]'); } catch(e) {}
        if (read.indexOf(slug) === -1) {
            read.push(slug);

            if (read.length > 200) read = read.slice(-200);
            localStorage.setItem('readPosts', JSON.stringify(read));
        }
    }
})();
document.addEventListener('DOMContentLoaded', () => {
    const wikiLinks = document.querySelectorAll('.prose a[href*="wikipedia.org/wiki/"]');
    if (wikiLinks.length === 0) return;

    const wikiCache = new Map();

    async function getWikiSummary(lang, title) {
        const safeTitle = encodeURIComponent(decodeURIComponent(title));
        const key = `${lang}:${safeTitle}`;
        if (wikiCache.has(key)) return wikiCache.get(key);

        const apiUrl = `https://${lang}.wikipedia.org/api/rest_v1/page/summary/${safeTitle}?redirect=true`;

        const res = await fetch(apiUrl);
        if (!res.ok) throw new Error(`HTTP ${res.status}`);
        const data = await res.json();
        wikiCache.set(key, data);
        return data;
    }

    function parseWikiUrl(href) {
        try {
            const url = new URL(href);
            const langMatch = url.hostname.match(/^([a-z\-]+)\.wikipedia\.org/);
            const lang = langMatch ? langMatch[1] : 'en';
            const parts = url.pathname.split('/wiki/');
            if (parts.length < 2 || !parts[1]) return null;
            return { lang, title: parts[1], href: url.href };
        } catch {
            return null;
        }
    }

    function isMobileOrTouch() {
        return window.innerWidth < 768 || window.matchMedia('(hover: none)').matches || window.matchMedia('(pointer: coarse)').matches;
    }

    const desktopTooltip = document.createElement('div');
    desktopTooltip.className = 'wiki-tooltip';
    document.body.appendChild(desktopTooltip);

    let hoverTimeout;
    let currentDesktopLink = null;

    function showDesktopLoading(targetEl) {
        positionDesktopTooltip(targetEl);
        desktopTooltip.innerHTML = '<div class="wiki-loading">Memuat dari Wikipedia...</div>';
        desktopTooltip.classList.add('visible');
    }

    function hideDesktopTooltip() {
        desktopTooltip.classList.remove('visible');
    }

    function positionDesktopTooltip(targetEl) {
        const rect = targetEl.getBoundingClientRect();
        const tooltipWidth = 320;

        let top = rect.bottom + window.scrollY + 8;
        let left = rect.left + window.scrollX - (tooltipWidth / 2) + (rect.width / 2);

        if (left < 16) left = 16;
        if (left + tooltipWidth > window.innerWidth - 16) {
            left = window.innerWidth - tooltipWidth - 16;
        }

        desktopTooltip.style.top = `${top}px`;
        desktopTooltip.style.left = `${left}px`;
    }

    function renderDesktopTooltip(data) {
        let html = '';
        if (data.thumbnail && data.thumbnail.source) {
            html += `<img src="${data.thumbnail.source}" class="wiki-thumb" alt="">`;
        }
        html += `<div class="wiki-content">`;
        html += `<div class="wiki-title">${data.title}</div>`;

        let extract = data.extract || '';
        if (extract.length > 200) extract = extract.substring(0, 200) + '...';

        html += `<p class="wiki-extract">${extract}</p>`;
        html += `<div class="wiki-footer">W — Disediakan oleh Wikipedia</div>`;
        html += `</div>`;

        desktopTooltip.innerHTML = html;
    }

    desktopTooltip.addEventListener('mouseenter', () => clearTimeout(hoverTimeout));
    desktopTooltip.addEventListener('mouseleave', () => {
        hoverTimeout = setTimeout(hideDesktopTooltip, 300);
    });

    const sheetModal = document.createElement('div');
    sheetModal.id = 'wiki-bottom-sheet';
    sheetModal.className = 'wiki-sheet-modal';
    sheetModal.setAttribute('role', 'dialog');
    sheetModal.setAttribute('aria-modal', 'true');
    sheetModal.setAttribute('aria-hidden', 'true');
    sheetModal.setAttribute('aria-label', 'Wikipedia Preview');
    sheetModal.innerHTML = `
        <div class="wiki-sheet-backdrop"></div>
        <div class="wiki-sheet-container">
            <div class="wiki-sheet-drag-area">
                <div class="wiki-sheet-handle"></div>
            </div>
            <div class="wiki-sheet-header">
                <div class="wiki-sheet-badge">
                    <svg class="wiki-sheet-logo" viewBox="0 0 24 24" width="15" height="15" fill="currentColor">
                        <path d="M12.09 13.124l2.77-7.858h2.02l-4.148 10.978h-1.606l-2.614-7.398-2.615 7.398H4.29L.142 5.266h2.02l2.77 7.858 2.383-6.732h1.666l2.383 6.732 2.726-7.724z"/>
                    </svg>
                    <span>Wikipedia</span>
                </div>
                <button type="button" class="wiki-sheet-close" aria-label="Tutup pratinjau">
                    <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                        <line x1="18" y1="6" x2="6" y2="18"></line>
                        <line x1="6" y1="6" x2="18" y2="18"></line>
                    </svg>
                </button>
            </div>
            <div class="wiki-sheet-body"></div>
            <div class="wiki-sheet-footer">
                <a href="#" target="_blank" rel="noopener noreferrer" class="wiki-sheet-btn">
                    <span>Buka di Wikipedia</span>
                    <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                        <path d="M18 13v6a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V8a2 2 0 0 1 2-2h6"></path>
                        <polyline points="15 3 21 3 21 9"></polyline>
                        <line x1="10" y1="14" x2="21" y2="3"></line>
                    </svg>
                </a>
            </div>
        </div>
    `;
    document.body.appendChild(sheetModal);

    const sheetBackdrop = sheetModal.querySelector('.wiki-sheet-backdrop');
    const sheetContainer = sheetModal.querySelector('.wiki-sheet-container');
    const sheetDragArea = sheetModal.querySelector('.wiki-sheet-drag-area');
    const sheetCloseBtn = sheetModal.querySelector('.wiki-sheet-close');
    const sheetBody = sheetModal.querySelector('.wiki-sheet-body');
    const sheetExtBtn = sheetModal.querySelector('.wiki-sheet-btn');

    function openBottomSheet(lang, title, href) {
        hideDesktopTooltip();
        document.body.classList.add('wiki-sheet-lock');
        sheetExtBtn.href = href;
        sheetContainer.style.transform = '';
        sheetBody.innerHTML = `
            <div class="wiki-sheet-loading">
                <div class="wiki-sheet-spinner"></div>
                <span>Memuat ringkasan Wikipedia...</span>
            </div>
        `;
        sheetModal.classList.add('open');
        sheetModal.setAttribute('aria-hidden', 'false');

        getWikiSummary(lang, title).then(data => {
            if (!sheetModal.classList.contains('open')) return;
            renderBottomSheet(data);
        }).catch(err => {
            if (!sheetModal.classList.contains('open')) return;
            sheetBody.innerHTML = `
                <div class="wiki-sheet-error">
                    <p>Gagal memuat pratinjau Wikipedia.</p>
                </div>
            `;
        });
    }

    function renderBottomSheet(data) {
        let html = '';
        if (data.thumbnail && data.thumbnail.source) {
            html += `<div class="wiki-sheet-thumb-wrap"><img src="${data.thumbnail.source}" class="wiki-sheet-thumb" alt="${data.title || ''}"></div>`;
        }
        html += `<div class="wiki-sheet-text">`;
        html += `<h3 class="wiki-sheet-title">${data.title}</h3>`;
        if (data.description) {
            html += `<div class="wiki-sheet-meta">${data.description}</div>`;
        }
        let extract = data.extract || '';
        if (extract.length > 280) extract = extract.substring(0, 280) + '...';
        html += `<p class="wiki-sheet-extract">${extract}</p>`;
        html += `</div>`;
        sheetBody.innerHTML = html;
    }

    function closeBottomSheet() {
        sheetModal.classList.remove('open');
        sheetModal.setAttribute('aria-hidden', 'true');
        document.body.classList.remove('wiki-sheet-lock');
        sheetContainer.style.transform = '';
    }

    sheetBackdrop.addEventListener('click', closeBottomSheet);
    sheetCloseBtn.addEventListener('click', closeBottomSheet);

    document.addEventListener('keydown', (e) => {
        if (e.key === 'Escape' && sheetModal.classList.contains('open')) {
            closeBottomSheet();
        }
    });

    let startY = 0;
    let isDragging = false;
    let currentDelta = 0;

    sheetDragArea.addEventListener('touchstart', (e) => {
        startY = e.touches[0].clientY;
        isDragging = true;
        currentDelta = 0;
        sheetContainer.style.transition = 'none';
    }, { passive: true });

    sheetDragArea.addEventListener('touchmove', (e) => {
        if (!isDragging) return;
        const touchY = e.touches[0].clientY;
        currentDelta = touchY - startY;
        if (currentDelta > 0) {
            sheetContainer.style.transform = `translateY(${currentDelta}px)`;
        }
    }, { passive: true });

    sheetDragArea.addEventListener('touchend', () => {
        if (!isDragging) return;
        isDragging = false;
        sheetContainer.style.transition = '';
        if (currentDelta > 70) {
            closeBottomSheet();
        } else {
            sheetContainer.style.transform = '';
        }
        startY = 0;
        currentDelta = 0;
    });

    wikiLinks.forEach(link => {
        const parsed = parseWikiUrl(link.href);
        if (!parsed) return;

        link.addEventListener('mouseenter', () => {
            if (isMobileOrTouch()) return;
            clearTimeout(hoverTimeout);
            currentDesktopLink = link;

            hoverTimeout = setTimeout(async () => {
                showDesktopLoading(link);
                try {
                    const data = await getWikiSummary(parsed.lang, parsed.title);
                    if (currentDesktopLink === link && desktopTooltip.classList.contains('visible')) {
                        renderDesktopTooltip(data);
                        positionDesktopTooltip(link);
                    }
                } catch (e) {
                    if (currentDesktopLink === link) {
                        desktopTooltip.innerHTML = `<div class="wiki-error">Gagal: ${e.message}</div>`;
                    }
                }
            }, 300);
        });

        link.addEventListener('mouseleave', () => {
            if (isMobileOrTouch()) return;
            clearTimeout(hoverTimeout);
            if (currentDesktopLink === link) {
                hoverTimeout = setTimeout(hideDesktopTooltip, 300);
            }
        });

        link.addEventListener('click', (e) => {
            if (isMobileOrTouch()) {
                e.preventDefault();
                openBottomSheet(parsed.lang, parsed.title, link.href);
            }
        });
    });
});
(function() {
    window.commentVisibleCount = 3;

    function handleLoadMoreComments() {
        var container = document.getElementById('comment-items-container');
        var btnContainer = document.getElementById('comment-load-more-container');
        if (!container) return;

        var hiddenItems = container.querySelectorAll('.comment-root-item.hidden');
        var step = 10;
        for (var i = 0; i < hiddenItems.length && i < step; i++) {
            hiddenItems[i].classList.remove('hidden');
        }

        var currentlyVisible = container.querySelectorAll('.comment-root-item:not(.hidden)').length;
        window.commentVisibleCount = currentlyVisible;

        var remaining = container.querySelectorAll('.comment-root-item.hidden').length;
        if (remaining === 0) {
            if (btnContainer) btnContainer.classList.add('hidden');
        } else {
            var countEl = document.getElementById('comment-remaining-count');
            if (countEl) countEl.textContent = '(' + remaining + ')';
        }
    }
    window.handleLoadMoreComments = handleLoadMoreComments;

    function syncCommentVisibility() {
        var container = document.getElementById('comment-items-container');
        var btnContainer = document.getElementById('comment-load-more-container');
        if (!container) return;

        var targetVisible = window.commentVisibleCount || 3;
        var items = container.querySelectorAll('.comment-root-item');
        var hiddenCount = 0;
        for (var i = 0; i < items.length; i++) {
            if (i < targetVisible) {
                items[i].classList.remove('hidden');
            } else {
                items[i].classList.add('hidden');
                hiddenCount++;
            }
        }

        if (btnContainer) {
            if (hiddenCount === 0) {
                btnContainer.classList.add('hidden');
            } else {
                btnContainer.classList.remove('hidden');
                var countEl = document.getElementById('comment-remaining-count');
                if (countEl) countEl.textContent = '(' + hiddenCount + ')';
            }
        }
    }
    window.syncCommentVisibility = syncCommentVisibility;

    function onCommentPosted(parentId) {
        if (parentId) {
            var parentEl = document.getElementById('comment-' + parentId);
            if (parentEl) {
                var details = parentEl.querySelector('details.comment-replies-group');
                if (details) details.open = true;
                setTimeout(function() {
                    parentEl.scrollIntoView({ behavior: 'smooth', block: 'center' });
                }, 100);
            }
        } else {
            var container = document.getElementById('comment-items-container');
            if (container) {
                var items = container.querySelectorAll('.comment-root-item');
                if (items.length > 0) {
                    var firstItem = items[0];
                    firstItem.classList.remove('hidden');
                    firstItem.classList.add('bg-hover/60');
                    setTimeout(function() {
                        firstItem.classList.remove('bg-hover/60');
                    }, 2000);
                    setTimeout(function() {
                        firstItem.scrollIntoView({ behavior: 'smooth', block: 'center' });
                    }, 100);
                }
            }
            var btnContainer = document.getElementById('comment-load-more-container');
            if (btnContainer) {
                var remaining = document.querySelectorAll('.comment-root-item.hidden').length;
                if (remaining === 0) {
                    btnContainer.classList.add('hidden');
                } else {
                    var countEl = document.getElementById('comment-remaining-count');
                    if (countEl) countEl.textContent = '(' + remaining + ')';
                }
            }
        }
    }
    window.onCommentPosted = onCommentPosted;

    document.addEventListener('htmx:afterSwap', function(e) {
        if (e.target && (e.target.id === 'comment-list' || (e.target.querySelector && e.target.querySelector('#comment-list')))) {
            syncCommentVisibility();
        }
    });

    document.addEventListener('click', function(e) {
        if (e.target && e.target.closest && e.target.closest('#comment-load-more-btn')) {
            e.preventDefault();
            handleLoadMoreComments();
        }
    });

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
