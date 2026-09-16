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
                text.textContent = "Saved";
            }
        } else {
            btn.classList.remove("is-saved", "bg-hover", "text-text", "border-text/70", "font-bold", "bg-accent/15", "text-accent", "border-accent");
            if (icon) {
                icon.setAttribute("fill", "none");
            }
            if (text) {
                text.textContent = "Save";
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
