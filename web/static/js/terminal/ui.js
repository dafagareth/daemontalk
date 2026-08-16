
(function() {
	var DT = window.DaemonTerminal;

	DT.ui.updatePrompt = function () {
				var disp = DT.fs.getDisplayPath(DT.state.currentPath);
				DT.state.envVars["PWD"] = DT.state.currentPath;
				DT.state.elements.promptLabel.textContent = "visitor@daemontalk:" + disp + "$";
				DT.state.elements.titleText.textContent = "visitor@daemontalk: " + disp;
			}

			
	DT.ui.appendHistory = function (cmd, output) {
				var item = document.createElement("div");
				item.className = "space-y-1";

				if (cmd !== "") {
					var pLine = document.createElement("div");
					pLine.className = "flex items-baseline gap-2";
					var disp = DT.fs.getDisplayPath(DT.state.currentPath);
					pLine.innerHTML = '<span class="text-[var(--term-prompt)] font-semibold select-none">visitor@daemontalk:' + disp + '$</span> ' +
						'<span class="text-[var(--term-text)]">' + DT.ui.escapeHTML(cmd) + '</span>';
					item.appendChild(pLine);
				}

				if (output) {
					var outEl = document.createElement("pre");
					outEl.className = "text-[var(--term-text)] whitespace-pre-wrap break-words leading-relaxed pl-2 text-xs sm:text-sm font-mono select-text";
					outEl.style.background = "transparent";
					outEl.style.border = "none";
					outEl.style.padding = "0";
					outEl.style.margin = "0";
					outEl.textContent = output;
					item.appendChild(outEl);
				}

				DT.state.elements.historyEl.appendChild(item);
				DT.state.elements.screenEl.scrollTop = DT.state.elements.screenEl.scrollHeight;
			}

			
	DT.ui.escapeRegExp = function (string) {
				return string.replace(/[.*+?^${}()|[\]\\]/g, '\\$&');
			}

			
	DT.ui.escapeHTML = function (str) {
				return (str || "").replace(/&/g, "&amp;").replace(/</g, "&lt;").replace(/>/g, "&gt;").replace(/"/g, "&quot;");
			}

			// Tab Auto-Completion Engine
			
	DT.ui.handleTabComplete = function () {
				var val = DT.state.elements.inputEl.value;
				var cursor = DT.state.elements.inputEl.selectionStart;
				var beforeCursor = val.substring(0, cursor);
				var tokens = beforeCursor.split(/\s+/);
				var lastToken = tokens[tokens.length - 1] || "";

				// If completing command name
				if (tokens.length <= 1) {
					var cmds = DT.fs.vfs["/bin"].children || [];
					var matches = cmds.filter(function(c) { return c.startsWith(lastToken); });
					if (matches.length === 1) {
						DT.state.elements.inputEl.value = matches[0] + " ";
					} else if (matches.length > 1) {
						DT.ui.appendHistory(val, matches.join("  "));
					}
					return;
				}

				// Completing files or directories
				var dirPath = DT.state.currentPath;
				var prefix = lastToken;
				if (lastToken.indexOf("/") !== -1) {
					var lastSlash = lastToken.lastIndexOf("/");
					var dirPart = lastToken.substring(0, lastSlash + 1);
					prefix = lastToken.substring(lastSlash + 1);
					dirPath = DT.fs.resolvePath(dirPart);
				}

				var node = DT.fs.vfs[dirPath];
				if (node && node.type === "dir") {
					var entries = (node.children || []).filter(function(c) { return c.startsWith(prefix); });
					if (entries.length === 1) {
						var match = entries[0];
						var full = dirPath === "/" ? "/" + match : dirPath + "/" + match;
						var isDir = DT.fs.vfs[full] && DT.fs.vfs[full].type === "dir";
						var replaceStr = (lastToken.indexOf("/") !== -1 ? lastToken.substring(0, lastToken.lastIndexOf("/") + 1) : "") + match + (isDir ? "/" : " ");
						tokens[tokens.length - 1] = replaceStr;
						DT.state.elements.inputEl.value = tokens.join(" ");
					} else if (entries.length > 1) {
						DT.ui.appendHistory(val, entries.join("  "));
					}
				}
			}

			// Key Event Listeners
			
	DT.ui.renderMotd = function () {
				var motd = document.createElement("div");
				motd.className = "mb-4 text-xs text-[var(--term-muted)] space-y-1 font-mono pb-3 border-b border-[var(--term-border)]/50";
				motd.innerHTML = "<div class='text-[var(--term-text)] font-semibold'>daemontalk Linux 6.12.8-generic (x86_64) &bull; bash 5.2.26</div>" +
					"<div>Press <kbd class='px-1.5 py-0.5 rounded border border-[var(--term-border)] bg-[var(--term-chip-bg)] text-[var(--term-text)] text-[11px]'>Tab</kbd> to auto-complete &bull; Type <span class='text-[var(--term-link)] font-semibold cursor-pointer underline hover:text-[var(--term-text)]' onclick=\"window.termExec('help')\">help</span> &bull; Try <span class='text-[var(--term-prompt)] font-bold cursor-pointer underline hover:text-[var(--term-text)]' onclick=\"window.termExec('incident list')\">incident list</span></div>";
				DT.state.elements.historyEl.appendChild(motd);
			}

			

	window.termClear = function() {
				DT.state.elements.historyEl.innerHTML = "";
				DT.state.elements.inputEl.value = "";
				DT.state.elements.inputEl.focus();
			};

			window.termExec = function(cmd) {
				DT.state.elements.inputEl.value = cmd;
				DT.state.elements.inputEl.focus();
				DT.cmd.execute(cmd);
			};

			window.termToggleFullscreen = function() {
				var win = document.getElementById("term-window");
				if (!document.fullscreenElement) {
					win.requestFullscreen().catch(function(e) { console.warn(e); });
				} else {
					document.exitFullscreen();
				}
			};

			

	DT.ui.init = function() {
		DT.state.elements.inputEl.addEventListener("keydown", function(e) {
				if (e.key === "Enter") {
					e.preventDefault();
					var val = DT.state.elements.inputEl.value;
					DT.state.elements.inputEl.value = "";
					DT.cmd.execute(val);
				} else if (e.key === "Tab") {
					e.preventDefault();
					DT.ui.handleTabComplete();
				} else if (e.key === "ArrowUp") {
					e.preventDefault();
					if (DT.state.cmdHistory.length === 0) return;
					if (DT.state.histIndex === DT.state.cmdHistory.length) {
						DT.state.currentInputBuffer = DT.state.elements.inputEl.value;
					}
					if (DT.state.histIndex > 0) {
						DT.state.histIndex--;
						DT.state.elements.inputEl.value = DT.state.cmdHistory[DT.state.histIndex];
					}
				} else if (e.key === "ArrowDown") {
					e.preventDefault();
					if (DT.state.histIndex < DT.state.cmdHistory.length - 1) {
						DT.state.histIndex++;
						DT.state.elements.inputEl.value = DT.state.cmdHistory[DT.state.histIndex];
					} else if (DT.state.histIndex === DT.state.cmdHistory.length - 1) {
						DT.state.histIndex = DT.state.cmdHistory.length;
						DT.state.elements.inputEl.value = DT.state.currentInputBuffer;
					}
				} else if (e.ctrlKey && e.key === "l") {
					e.preventDefault();
					window.termClear();
				} else if (e.ctrlKey && e.key === "c") {
					e.preventDefault();
					var val = DT.state.elements.inputEl.value;
					DT.state.elements.inputEl.value = "";
					DT.ui.appendHistory(val + "^C", "");
				} else if (e.ctrlKey && e.key === "u") {
					e.preventDefault();
					DT.state.elements.inputEl.value = "";
				}
			});

			
		
		// Initial Setup
		DT.ui.renderMotd();
		DT.ui.updatePrompt();
		setTimeout(function() { DT.state.elements.inputEl.focus(); }, 100);
	};
})();
