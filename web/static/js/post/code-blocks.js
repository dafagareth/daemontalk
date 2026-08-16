(function() {
    // Enhanced Code block wrapper: adds language badges and copy button.
    document.querySelectorAll("#prose-body pre").forEach(function(pre) {
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
