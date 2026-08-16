
(function() {
	var DT = window.DaemonTerminal;

	DT.cmd.evalSingleCommand = function (cmdLine, stdin) {
				var args = DT.cmd.parseArgs(cmdLine);
				var cmd = (args[0] || "").toLowerCase();

				// Alias resolution
				if (DT.state.aliases[cmd]) {
					var aliasedLine = DT.state.aliases[cmd] + (args.length > 1 ? " " + args.slice(1).join(" ") : "");
					args = DT.cmd.parseArgs(aliasedLine);
					cmd = (args[0] || "").toLowerCase();
				}

				var out = "";

				switch (cmd) {
					case "help":
						out = "daemontalk Linux Shell - Available Commands:\n\n" +
							"  File & Directory Operations:\n" +
							"    ls [-l|-a|-la] [dir]    List directory contents\n" +
							"    cd [dir]                Change directory (cd posts, cd .., cd ~)\n" +
							"    pwd                     Print current working directory\n" +
							"    mkdir [-p] <dir>        Create new directory\n" +
							"    rm [-r|-rf] <path>      Remove file or directory\n" +
							"    cp [-r] <src> <dest>    Copy files and directories\n" +
							"    mv <src> <dest>         Move or rename files\n" +
							"    touch <file>            Create empty file or update timestamp\n" +
							"    tree [-L N] [dir]       Display directory tree visually\n" +
							"    stat <file>             Display detailed file status and metadata\n" +
							"    file <path>             Determine file type\n" +
							"    diff [-u] <f1> <f2>     Compare two files line by line\n" +
							"    find [dir] [-name pat]  Search for files in a directory tree\n\n" +
							"  Text & Stream Processing:\n" +
							"    cat <file...>           Concatenate and display files\n" +
							"    head [-n N] [file]      Output first N lines of input\n" +
							"    tail [-n N] [file]      Output last N lines of input\n" +
							"    grep [-i|-n|-v] <pat>   Search for pattern in files or stdin\n" +
							"    sort [-r|-n|-u] [file]  Sort lines of text\n" +
							"    uniq [-c|-d] [file]     Report or omit repeated lines\n" +
							"    cut -d <c> -f <num>     Extract fields from delimited text\n" +
							"    tr <set1> <set2>        Translate or replace characters\n" +
							"    sed 's/find/repl/g'     Stream editor regex replacement\n" +
							"    awk '{print $N}'        Extract column fields\n" +
							"    wc [-l|-w|-c] [file]    Count lines, words, and characters\n" +
							"    base64 [-d] [text]      Base64 encode or decode data\n" +
							"    rev [text|file]         Reverse lines character by character\n" +
							"    echo [text]             Print text (supports > and >> redirects)\n" +
							"    seq [first] <last>      Generate integer sequence numbers\n" +
							"    basename / dirname      Extract filename or directory components\n\n" +
							"  System & Process Inspection:\n" +
							"    ps [aux|-ef]            Report active process snapshots\n" +
							"    top / htop              Display system tasks and resource usage\n" +
							"    df [-h]                 Report disk space usage\n" +
							"    du [-h|-sh] [dir]       Estimate file space usage\n" +
							"    free [-h|-m]            Display free and used memory\n" +
							"    uptime, arch, uname -a  System uptime and kernel architecture\n" +
							"    hostname, whoami, id    User identity and host information\n" +
							"    env / printenv [var]    Print environment variables\n" +
							"    export KEY=VAL          Set environment variable\n" +
							"    dmesg                   Display kernel buffer logs\n" +
							"    kill [-9] <PID>         Terminate active processes\n\n" +
							"  Troubleshooting & Linux War Stories:\n" +
							"    incident [list|1-5|hint|verify]  Interactive Linux production RCA incident challenges\n" +
							"    recipes / cheat / ebpf           Curated eBPF tracing & performance recipes\n" +
							"    lsof [-i :port]                  List open files and network sockets\n\n" +
							"  Networking & Utilities:\n" +
							"    ping [-c N] <host>      Send ICMP echo requests\n" +
							"    dig / nslookup <domain> DNS query inspection\n" +
							"    ip [addr|route|link]    Display network interfaces\n" +
							"    ifconfig, netstat, ss   Network statistics & listening ports\n" +
							"    curl [-s] <url>         Fetch URL or API data\n" +
							"    run <lang> <code>       Execute Go, Python, JS, or Bash code\n" +
							"    man <command>           View detailed manual pages\n" +
							"    which <cmd>, alias      Locate binary or manage command DT.state.aliases\n" +
							"    cal, date, history      Calendar, clock, and command history\n" +
							"    clear, exit             Clear screen or return to homepage\n\n" +
							"  Piping & Redirection:\n" +
							"    Supports standard pipes '|', write '>', and append '>>'.\n" +
							"    Example: cat /etc/passwd | cut -d: -f1 | sort | head -n 5";
						break;

					case "man":
						if (args.length < 2) {
							out = "What manual page do you want?\nFor example, try 'man ls', 'man grep', 'man cat', 'man find', 'man ps'.";
						} else {
							var targetCmd = args[1].toLowerCase();
							if (DT.fs.manPages[targetCmd]) {
								out = DT.fs.manPages[targetCmd];
							} else {
								out = "No manual entry for " + targetCmd + "\nType 'help' to see all available commands.";
							}
						}
						break;

					case "pwd":
						out = DT.state.currentPath;
						break;

					case "cd":
						var target = args[1] || "~";
						var resolved = DT.fs.resolvePath(target);
						if (!DT.fs.vfs[resolved]) {
							out = "cd: no such file or directory: " + target;
						} else if (DT.fs.vfs[resolved].type !== "dir") {
							out = "cd: not a directory: " + target;
						} else {
							DT.state.currentPath = resolved;
							DT.ui.updatePrompt();
						}
						break;

					case "ls":
						var longFormat = false;
						var showAll = false;
						var targetDir = DT.state.currentPath;

						for (var i = 1; i < args.length; i++) {
							var a = args[i];
							if (a === "-l") longFormat = true;
							else if (a === "-a") showAll = true;
							else if (a === "-la" || a === "-al") { longFormat = true; showAll = true; }
							else if (a === "-lh" || a === "-lah" || a === "-hal") { longFormat = true; showAll = true; }
							else if (!a.startsWith("-")) { targetDir = DT.fs.resolvePath(a); }
						}

						var node = DT.fs.vfs[targetDir];
						if (!node) {
							out = "ls: cannot access '" + (args[1] || targetDir) + "': No such file or directory";
						} else if (node.type === "file") {
							out = targetDir.split("/").pop();
						} else {
							var entries = (node.children || []).slice();
							entries.sort();
							if (showAll) entries = [".", ".."].concat(entries);

							if (longFormat) {
								var lines = ["total " + entries.length];
								entries.forEach(function(name) {
									var full = targetDir === "/" ? "/" + name : targetDir + "/" + name;
									var isDir = (DT.fs.vfs[full] && DT.fs.vfs[full].type === "dir") || name === "." || name === "..";
									var perms = isDir ? "drwxr-xr-x" : "-rw-r--r--";
									var size = isDir ? "4096" : String((DT.fs.vfs[full] && DT.fs.vfs[full].content ? DT.fs.vfs[full].content.length : 0));
									while (size.length < 6) size = " " + size;
									lines.push(perms + " 1 visitor visitor " + size + " Aug  8 16:00 " + name);
								});
								out = lines.join("\n");
							} else {
								out = entries.join("  ");
							}
						}
						break;

					case "cat":
						if (args.length < 2) {
							if (stdin !== undefined && stdin !== null) {
								out = stdin;
							} else {
								out = "cat: missing file operand";
							}
						} else {
							var outputs = [];
							for (var i = 1; i < args.length; i++) {
								var fileTarget = DT.fs.resolvePath(args[i]);
								var fnode = DT.fs.vfs[fileTarget];
								if (!fnode) {
									outputs.push("cat: " + args[i] + ": No such file or directory");
								} else if (fnode.type === "dir") {
									outputs.push("cat: " + args[i] + ": Is a directory");
								} else {
									outputs.push(fnode.content || "");
								}
							}
							out = outputs.join("\n");
						}
						break;

					case "head":
					case "tail":
						var n = 10;
						var filename = null;
						for (var i = 1; i < args.length; i++) {
							if (args[i] === "-n" && args[i+1]) {
								n = parseInt(args[i+1], 10) || 10;
								i++;
							} else if (!args[i].startsWith("-")) {
								filename = args[i];
							}
						}
						var textData = "";
						if (filename) {
							var fileTarget = DT.fs.resolvePath(filename);
							var fnode = DT.fs.vfs[fileTarget];
							if (!fnode || fnode.type !== "file") {
								out = cmd + ": cannot open '" + filename + "': No such file";
								break;
							}
							textData = fnode.content || "";
						} else if (stdin !== undefined && stdin !== null) {
							textData = stdin;
						} else {
							out = cmd + ": missing file operand";
							break;
						}
						var lines = textData.split("\n");
						if (lines[lines.length - 1] === "") lines.pop();
						if (cmd === "head") out = lines.slice(0, n).join("\n");
						else out = lines.slice(Math.max(0, lines.length - n)).join("\n");
						break;

					case "grep":
						var caseInsensitive = false;
						var invertMatch = false;
						var lineNumbers = false;
						var countOnly = false;
						var pattern = "";
						var targetFile = null;

						for (var i = 1; i < args.length; i++) {
							var a = args[i];
							if (a === "-i") caseInsensitive = true;
							else if (a === "-v") invertMatch = true;
							else if (a === "-n") lineNumbers = true;
							else if (a === "-c") countOnly = true;
							else if (!pattern) pattern = a;
							else if (!targetFile) targetFile = a;
						}

						if (!pattern) {
							out = "grep: missing pattern operand";
							break;
						}

						var rx = new RegExp(DT.ui.escapeRegExp(pattern), caseInsensitive ? "i" : "");
						var matches = [];

						if (stdin !== undefined && stdin !== null && !targetFile) {
							var sLines = stdin.split("\n");
							sLines.forEach(function(l, idx) {
								var isMatch = rx.test(l);
								if (invertMatch) isMatch = !isMatch;
								if (isMatch) matches.push(lineNumbers ? (idx + 1) + ": " + l : l);
							});
							out = countOnly ? String(matches.length) : matches.join("\n");
						} else if (targetFile) {
							var resolvedFile = DT.fs.resolvePath(targetFile);
							var fnode = DT.fs.vfs[resolvedFile];
							if (!fnode || fnode.type !== "file") {
								out = "grep: " + targetFile + ": No such file";
							} else {
								(fnode.content || "").split("\n").forEach(function(l, idx) {
									var isMatch = rx.test(l);
									if (invertMatch) isMatch = !isMatch;
									if (isMatch) matches.push(lineNumbers ? (idx + 1) + ": " + l : l);
								});
								out = countOnly ? String(matches.length) : matches.join("\n");
							}
						} else {
							var cnode = DT.fs.vfs[DT.state.currentPath];
							(cnode.children || []).forEach(function(fname) {
								var full = DT.state.currentPath === "/" ? "/" + fname : DT.state.currentPath + "/" + fname;
								var f = DT.fs.vfs[full];
								if (f && f.type === "file") {
									(f.content || "").split("\n").forEach(function(l, idx) {
										var isMatch = rx.test(l);
										if (invertMatch) isMatch = !isMatch;
										if (isMatch) matches.push(fname + (lineNumbers ? ":" + (idx + 1) : "") + ": " + l);
									});
								}
							});
							out = countOnly ? String(matches.length) : (matches.join("\n") || "grep: no matches found");
						}
						break;

					case "tree":
						var maxDepth = 3;
						var startDir = DT.state.currentPath;
						for (var i = 1; i < args.length; i++) {
							if (args[i] === "-L" && args[i+1]) {
								maxDepth = parseInt(args[i+1], 10) || 3;
								i++;
							} else if (!args[i].startsWith("-")) {
								startDir = DT.fs.resolvePath(args[i]);
							}
						}
						if (!DT.fs.vfs[startDir]) {
							out = "tree: '" + startDir + "': No such file or directory";
							break;
						}
						var treeLines = [startDir === "/home/visitor" ? "." : startDir];
						var totalDirs = 0;
						var totalFiles = 0;

						function buildTree(dir, prefix, depth) {
							if (depth > maxDepth) return;
							var dNode = DT.fs.vfs[dir];
							if (!dNode || dNode.type !== "dir") return;
							var kids = (dNode.children || []).slice().sort();
							kids.forEach(function(childName, idx) {
								var isLast = idx === kids.length - 1;
								var pointer = isLast ? "└── " : "├── ";
								var childPath = dir === "/" ? "/" + childName : dir + "/" + childName;
								var cNode = DT.fs.vfs[childPath];
								var isDir = cNode && cNode.type === "dir";

								if (isDir) {
									totalDirs++;
									treeLines.push(prefix + pointer + childName);
									buildTree(childPath, prefix + (isLast ? "    " : "│   "), depth + 1);
								} else {
									totalFiles++;
									treeLines.push(prefix + pointer + childName);
								}
							});
						}

						buildTree(startDir, "", 1);
						treeLines.push("\n" + totalDirs + " directories, " + totalFiles + " files");
						out = treeLines.join("\n");
						break;

					case "stat":
						if (args.length < 2) {
							out = "stat: missing operand";
						} else {
							var statPath = DT.fs.resolvePath(args[1]);
							var sNode = DT.fs.vfs[statPath];
							if (!sNode) {
								out = "stat: cannot statx '" + args[1] + "': No such file or directory";
							} else {
								var fname = statPath.split("/").pop() || "/";
								var isDir = sNode.type === "dir";
								var size = isDir ? 4096 : (sNode.content ? sNode.content.length : 0);
								var now = new Date().toISOString().replace("T", " ").substring(0, 19) + " +0700";
								out = "  File: " + fname + "\n" +
									"  Size: " + size + "\tBlocks: " + Math.max(8, Math.ceil(size / 512)) + "   IO Block: 4096   " + (isDir ? "directory" : "regular file") + "\n" +
									"Device: nvme0n1p2\tInode: " + (Math.floor(Math.random() * 8000000) + 1000000) + "   Links: " + (isDir ? (sNode.children ? sNode.children.length + 2 : 2) : 1) + "\n" +
									"Access: (" + (isDir ? "0755/drwxr-xr-x" : "0644/-rw-r--r--") + ")  Uid: ( 1000/ visitor)   Gid: ( 1000/ visitor)\n" +
									"Access: " + now + "\n" +
									"Modify: " + now + "\n" +
									"Change: " + now + "\n" +
									" Birth: " + now;
							}
						}
						break;

					case "file":
						if (args.length < 2) {
							out = "file: missing operand";
						} else {
							var fResults = [];
							for (var i = 1; i < args.length; i++) {
								var fp = DT.fs.resolvePath(args[i]);
								var fNode = DT.fs.vfs[fp];
								if (!fNode) {
									fResults.push(args[i] + ": cannot open (No such file or directory)");
								} else if (fNode.type === "dir") {
									fResults.push(args[i] + ": directory");
								} else if (fp.endsWith(".md")) {
									fResults.push(args[i] + ": Markdown document, UTF-8 Unicode text");
								} else if (fp.endsWith(".json")) {
									fResults.push(args[i] + ": JSON data text");
								} else if (fp.startsWith("/bin/")) {
									fResults.push(args[i] + ": ELF 64-bit LSB pie executable, x86-64, dynamic linked");
								} else {
									fResults.push(args[i] + ": ASCII text, with CRLF line terminators");
								}
							}
							out = fResults.join("\n");
						}
						break;

					case "diff":
						if (args.length < 3) {
							out = "diff: missing operand after '" + (args[1] || "") + "'";
						} else {
							var p1 = DT.fs.resolvePath(args[1]);
							var p2 = DT.fs.resolvePath(args[2]);
							var n1 = DT.fs.vfs[p1];
							var n2 = DT.fs.vfs[p2];
							if (!n1) out = "diff: " + args[1] + ": No such file";
							else if (!n2) out = "diff: " + args[2] + ": No such file";
							else {
								var l1 = (n1.content || "").split("\n");
								var l2 = (n2.content || "").split("\n");
								var diffOut = ["--- " + args[1] + "\t2026-08-08 16:00:00", "+++ " + args[2] + "\t2026-08-08 16:00:00", "@@ -1," + l1.length + " +1," + l2.length + " @@"];
								var maxL = Math.max(l1.length, l2.length);
								var hasDiff = false;
								for (var j = 0; j < maxL; j++) {
									if (l1[j] !== l2[j]) {
										hasDiff = true;
										if (l1[j] !== undefined) diffOut.push("-" + l1[j]);
										if (l2[j] !== undefined) diffOut.push("+" + l2[j]);
									} else if (l1[j] !== undefined) {
										diffOut.push(" " + l1[j]);
									}
								}
								out = hasDiff ? diffOut.join("\n") : "";
							}
						}
						break;

					case "sort":
						var reverse = false;
						var numeric = false;
						var unique = false;
						var sortFile = null;
						for (var i = 1; i < args.length; i++) {
							if (args[i] === "-r") reverse = true;
							else if (args[i] === "-n") numeric = true;
							else if (args[i] === "-u") unique = true;
							else if (!args[i].startsWith("-")) sortFile = args[i];
						}
						var sortText = "";
						if (sortFile) {
							var sfNode = DT.fs.vfs[DT.fs.resolvePath(sortFile)];
							if (!sfNode || sfNode.type !== "file") {
								out = "sort: cannot read: " + sortFile;
								break;
							}
							sortText = sfNode.content || "";
						} else if (stdin !== undefined && stdin !== null) {
							sortText = stdin;
						}
						var sList = sortText.split("\n").filter(function(l) { return l.length > 0; });
						sList.sort(function(a, b) {
							if (numeric) {
								return parseFloat(a) - parseFloat(b);
							}
							return a.localeCompare(b);
						});
						if (reverse) sList.reverse();
						if (unique) {
							sList = sList.filter(function(item, pos, self) {
								return self.indexOf(item) === pos;
							});
						}
						out = sList.join("\n");
						break;

					case "uniq":
						var countFlag = false;
						var dupOnly = false;
						var uFile = null;
						for (var i = 1; i < args.length; i++) {
							if (args[i] === "-c") countFlag = true;
							else if (args[i] === "-d") dupOnly = true;
							else if (!args[i].startsWith("-")) uFile = args[i];
						}
						var uText = "";
						if (uFile) {
							var ufNode = DT.fs.vfs[DT.fs.resolvePath(uFile)];
							if (!ufNode) { out = "uniq: " + uFile + ": No such file"; break; }
							uText = ufNode.content || "";
						} else if (stdin !== undefined && stdin !== null) {
							uText = stdin;
						}
						var uLines = uText.split("\n");
						var uResult = [];
						var lastLine = null;
						var runCount = 0;
						uLines.forEach(function(l) {
							if (l === lastLine) {
								runCount++;
							} else {
								if (lastLine !== null) {
									if (!dupOnly || runCount > 1) {
										uResult.push(countFlag ? "  " + runCount + " " + lastLine : lastLine);
									}
								}
								lastLine = l;
								runCount = 1;
							}
						});
						if (lastLine !== null && (!dupOnly || runCount > 1)) {
							uResult.push(countFlag ? "  " + runCount + " " + lastLine : lastLine);
						}
						out = uResult.join("\n");
						break;

					case "cut":
						var delim = "\t";
						var fieldIdx = 1;
						var cutFile = null;
						for (var i = 1; i < args.length; i++) {
							if (args[i] === "-d" && args[i+1]) { delim = args[i+1]; i++; }
							else if (args[i].startsWith("-d")) { delim = args[i].substring(2); }
							else if (args[i] === "-f" && args[i+1]) { fieldIdx = parseInt(args[i+1], 10) || 1; i++; }
							else if (args[i].startsWith("-f")) { fieldIdx = parseInt(args[i].substring(2), 10) || 1; }
							else if (!args[i].startsWith("-")) { cutFile = args[i]; }
						}
						var cutText = "";
						if (cutFile) {
							var cfn = DT.fs.vfs[DT.fs.resolvePath(cutFile)];
							if (!cfn) { out = "cut: " + cutFile + ": No such file"; break; }
							cutText = cfn.content || "";
						} else if (stdin !== undefined && stdin !== null) {
							cutText = stdin;
						}
						var cutLines = cutText.split("\n").map(function(l) {
							var parts = l.split(delim);
							return parts[fieldIdx - 1] !== undefined ? parts[fieldIdx - 1] : l;
						});
						out = cutLines.join("\n");
						break;

					case "tr":
						if (args.length < 3) {
							out = "tr: missing operand";
						} else {
							var set1 = args[1];
							var set2 = args[2];
							var trInput = stdin !== undefined && stdin !== null ? stdin : "";
							if (set1 === "a-z" && set2 === "A-Z") {
								out = trInput.toUpperCase();
							} else if (set1 === "A-Z" && set2 === "a-z") {
								out = trInput.toLowerCase();
							} else {
								out = trInput.split(set1).join(set2);
							}
						}
						break;

					case "sed":
						if (args.length < 2) {
							out = "sed: missing expression";
						} else {
							var expr = args[1];
							var sedFile = args[2];
							var sedText = "";
							if (sedFile) {
								var sn = DT.fs.vfs[DT.fs.resolvePath(sedFile)];
								if (!sn) { out = "sed: " + sedFile + ": No such file"; break; }
								sedText = sn.content || "";
							} else if (stdin !== undefined && stdin !== null) {
								sedText = stdin;
							}
							if (expr.startsWith("s/") || expr.startsWith("s|")) {
								var sParts = expr.split(expr[1]);
								var fPat = sParts[1] || "";
								var rPat = sParts[2] || "";
								var flags = sParts[3] || "g";
								try {
									var rExp = new RegExp(fPat, flags);
									out = sedText.replace(rExp, rPat);
								} catch (e) {
									out = "sed: invalid regular expression: " + fPat;
								}
							} else {
								out = sedText;
							}
						}
						break;

					case "awk":
						var awkCode = args[1] || "{print $0}";
						var awkFile = args[2];
						var awkText = "";
						if (awkFile) {
							var an = DT.fs.vfs[DT.fs.resolvePath(awkFile)];
							if (!an) { out = "awk: " + awkFile + ": No such file"; break; }
							awkText = an.content || "";
						} else if (stdin !== undefined && stdin !== null) {
							awkText = stdin;
						}
						var colMatch = awkCode.match(/\$(\d+)/);
						var colNum = colMatch ? parseInt(colMatch[1], 10) : 0;
						var awkLines = awkText.split("\n").map(function(line) {
							if (!line.trim()) return "";
							if (colNum === 0) return line;
							var cols = line.trim().split(/\s+/);
							return cols[colNum - 1] || "";
						});
						out = awkLines.join("\n");
						break;

					case "wc":
						var targetFile = null;
						var countL = false, countW = false, countC = false;
						for (var i = 1; i < args.length; i++) {
							if (args[i] === "-l") countL = true;
							else if (args[i] === "-w") countW = true;
							else if (args[i] === "-c" || args[i] === "-m") countC = true;
							else if (!args[i].startsWith("-")) targetFile = args[i];
						}
						var wcText = "";
						if (targetFile) {
							var wNode = DT.fs.vfs[DT.fs.resolvePath(targetFile)];
							if (!wNode || wNode.type !== "file") {
								out = "wc: " + targetFile + ": No such file";
								break;
							}
							wcText = wNode.content || "";
						} else if (stdin !== undefined && stdin !== null) {
							wcText = stdin;
						} else {
							out = "wc: missing file operand";
							break;
						}
						var wLines = wcText ? wcText.split("\n").length : 0;
						var wWords = wcText ? wcText.trim().split(/\s+/).filter(Boolean).length : 0;
						var wChars = wcText.length;

						if (!countL && !countW && !countC) {
							out = "  " + wLines + "  " + wWords + "  " + wChars + (targetFile ? " " + targetFile : "");
						} else {
							var res = [];
							if (countL) res.push(wLines);
							if (countW) res.push(wWords);
							if (countC) res.push(wChars);
							out = "  " + res.join("  ") + (targetFile ? " " + targetFile : "");
						}
						break;

					case "base64":
						var decode = false;
						var b64Payload = "";
						for (var i = 1; i < args.length; i++) {
							if (args[i] === "-d" || args[i] === "--decode") decode = true;
							else if (!b64Payload) b64Payload = args[i];
						}
						if (!b64Payload && stdin) b64Payload = stdin.trim();
						if (DT.fs.vfs[DT.fs.resolvePath(b64Payload)]) {
							b64Payload = DT.fs.vfs[DT.fs.resolvePath(b64Payload)].content || "";
						}
						try {
							if (decode) out = atob(b64Payload);
							else out = btoa(b64Payload);
						} catch (e) {
							out = "base64: invalid input";
						}
						break;

					case "rev":
						var revInput = "";
						if (args[1]) {
							var rn = DT.fs.vfs[DT.fs.resolvePath(args[1])];
							revInput = rn ? (rn.content || "") : args.slice(1).join(" ");
						} else if (stdin) {
							revInput = stdin;
						}
						out = revInput.split("\n").map(function(l) { return l.split("").reverse().join(""); }).join("\n");
						break;

					case "echo":
						out = args.slice(1).join(" ");
						break;

					case "seq":
						if (args.length < 2) {
							out = "seq: missing operand";
						} else if (args.length === 2) {
							var end = parseInt(args[1], 10) || 1;
							var arr = [];
							for (var i = 1; i <= Math.min(100, end); i++) arr.push(i);
							out = arr.join("\n");
						} else if (args.length === 3) {
							var start = parseInt(args[1], 10) || 1;
							var end = parseInt(args[2], 10) || 1;
							var arr = [];
							for (var i = start; i <= Math.min(start + 100, end); i++) arr.push(i);
							out = arr.join("\n");
						}
						break;

					case "basename":
						if (args.length < 2) out = "basename: missing operand";
						else {
							var bname = args[1].replace(/\/+$/, "").split("/").pop();
							if (args[2] && bname.endsWith(args[2])) {
								bname = bname.substring(0, bname.length - args[2].length);
							}
							out = bname;
						}
						break;

					case "dirname":
						if (args.length < 2) out = "dirname: missing operand";
						else {
							var dname = args[1].replace(/\/+$/, "");
							var lastSlash = dname.lastIndexOf("/");
							out = lastSlash === -1 ? "." : (lastSlash === 0 ? "/" : dname.substring(0, lastSlash));
						}
						break;

					case "cp":
						if (args.length < 3) {
							out = "cp: missing destination file operand";
						} else {
							var srcPath = DT.fs.resolvePath(args[1]);
							var destPath = DT.fs.resolvePath(args[2]);
							var srcNode = DT.fs.vfs[srcPath];
							if (!srcNode) {
								out = "cp: cannot stat '" + args[1] + "': No such file or directory";
							} else {
								var destParent = destPath.substring(0, destPath.lastIndexOf("/")) || "/";
								var destName = destPath.split("/").pop();
								if (DT.fs.vfs[destPath] && DT.fs.vfs[destPath].type === "dir") {
									destPath = destPath === "/" ? "/" + args[1].split("/").pop() : destPath + "/" + args[1].split("/").pop();
									destParent = destPath.substring(0, destPath.lastIndexOf("/")) || "/";
									destName = destPath.split("/").pop();
								}
								DT.fs.vfs[destPath] = JSON.parse(JSON.stringify(srcNode));
								if (DT.fs.vfs[destParent] && DT.fs.vfs[destParent].children.indexOf(destName) === -1) {
									DT.fs.vfs[destParent].children.push(destName);
								}
							}
						}
						break;

					case "mv":
						if (args.length < 3) {
							out = "mv: missing destination file operand";
						} else {
							var srcPath = DT.fs.resolvePath(args[1]);
							var destPath = DT.fs.resolvePath(args[2]);
							var srcNode = DT.fs.vfs[srcPath];
							if (!srcNode) {
								out = "mv: cannot stat '" + args[1] + "': No such file or directory";
							} else {
								var srcParent = srcPath.substring(0, srcPath.lastIndexOf("/")) || "/";
								var srcName = srcPath.split("/").pop();
								var destParent = destPath.substring(0, destPath.lastIndexOf("/")) || "/";
								var destName = destPath.split("/").pop();

								if (DT.fs.vfs[destPath] && DT.fs.vfs[destPath].type === "dir") {
									destPath = destPath === "/" ? "/" + srcName : destPath + "/" + srcName;
									destParent = destPath.substring(0, destPath.lastIndexOf("/")) || "/";
									destName = destPath.split("/").pop();
								}

								DT.fs.vfs[destPath] = srcNode;
								delete DT.fs.vfs[srcPath];

								if (DT.fs.vfs[srcParent]) {
									var sIdx = DT.fs.vfs[srcParent].children.indexOf(srcName);
									if (sIdx !== -1) DT.fs.vfs[srcParent].children.splice(sIdx, 1);
								}
								if (DT.fs.vfs[destParent] && DT.fs.vfs[destParent].children.indexOf(destName) === -1) {
									DT.fs.vfs[destParent].children.push(destName);
								}
							}
						}
						break;

					case "touch":
						if (args.length < 2) {
							out = "touch: missing file operand";
						} else {
							for (var i = 1; i < args.length; i++) {
								var newPath = DT.fs.resolvePath(args[i]);
								var parentDir = newPath.substring(0, newPath.lastIndexOf("/")) || "/";
								var fname = newPath.split("/").pop();
								if (!DT.fs.vfs[parentDir] || DT.fs.vfs[parentDir].type !== "dir") {
									out = "touch: cannot touch '" + args[i] + "': No such directory";
								} else if (!DT.fs.vfs[newPath]) {
									DT.fs.vfs[newPath] = { type: "file", content: "" };
									if (DT.fs.vfs[parentDir].children.indexOf(fname) === -1) {
										DT.fs.vfs[parentDir].children.push(fname);
									}
								}
							}
						}
						break;

					case "mkdir":
						if (args.length < 2) {
							out = "mkdir: missing operand";
						} else {
							for (var i = 1; i < args.length; i++) {
								if (args[i] === "-p") continue;
								var newDir = DT.fs.resolvePath(args[i]);
								var parentDir = newDir.substring(0, newDir.lastIndexOf("/")) || "/";
								var dname = newDir.split("/").pop();
								if (!DT.fs.vfs[parentDir] || DT.fs.vfs[parentDir].type !== "dir") {
									out = "mkdir: cannot create directory '" + args[i] + "': No such parent directory";
								} else if (DT.fs.vfs[newDir]) {
									out = "mkdir: cannot create directory '" + args[i] + "': File exists";
								} else {
									DT.fs.vfs[newDir] = { type: "dir", children: [] };
									if (DT.fs.vfs[parentDir].children.indexOf(dname) === -1) {
										DT.fs.vfs[parentDir].children.push(dname);
									}
								}
							}
						}
						break;

					case "rm":
						if (args.length < 2) {
							out = "rm: missing operand";
						} else {
							var rmArgs = args.slice(1).filter(function(a) { return !a.startsWith("-"); });
							rmArgs.forEach(function(target) {
								var rmPath = DT.fs.resolvePath(target);
								if (!DT.fs.vfs[rmPath]) {
									out = "rm: cannot remove '" + target + "': No such file or directory";
								} else {
									var parentDir = rmPath.substring(0, rmPath.lastIndexOf("/")) || "/";
									var name = rmPath.split("/").pop();
									delete DT.fs.vfs[rmPath];
									if (DT.fs.vfs[parentDir]) {
										var idx = DT.fs.vfs[parentDir].children.indexOf(name);
										if (idx !== -1) DT.fs.vfs[parentDir].children.splice(idx, 1);
									}
								}
							});
						}
						break;

					case "find":
						var startDir = args[1] && !args[1].startsWith("-") ? DT.fs.resolvePath(args[1]) : DT.state.currentPath;
						var namePattern = null;
						for (var i = 1; i < args.length; i++) {
							if (args[i] === "-name" && args[i+1]) { namePattern = args[i+1]; i++; }
						}
						var results = [];
						function walk(dir) {
							results.push(dir);
							var node = DT.fs.vfs[dir];
							if (node && node.type === "dir") {
								(node.children || []).forEach(function(c) {
									var childPath = dir === "/" ? "/" + c : dir + "/" + c;
									if (DT.fs.vfs[childPath] && DT.fs.vfs[childPath].type === "dir") {
										walk(childPath);
									} else {
										if (!namePattern || childPath.indexOf(namePattern.replace(/\*/g, "")) !== -1) {
											results.push(childPath);
										}
									}
								});
							}
						}
						if (DT.fs.vfs[startDir]) walk(startDir);
						else out = "find: '" + startDir + "': No such file or directory";
						if (!out) out = results.join("\n");
						break;

					case "ps":
						var psLines = [
							"PID  TTY          TIME CMD",
							"  1  ?        00:00:02 systemd",
							"  2  ?        00:00:00 kthreadd",
							" 14  ?        00:00:01 rcu_preempt",
							"104  ?        00:00:00 systemd-journald",
							"240  ?        00:00:00 postgres: 16/main",
							"312  ?        00:00:01 caddy -config /etc/caddy/Caddyfile",
							"450  ?        00:00:08 ./portfolio (daemontalk web)",
							"892  pts/0    00:00:00 bash"
						];
						if (DT.state.activeIncidentState.port8080PID !== null) {
							psLines.push("4192 ?        00:04:12 python3 -m http.server 8080");
						}
						if (DT.state.activeIncidentState.zombieParentPID !== null) {
							psLines.push("1337 ?        00:00:15 python3 bad_supervisor.py");
							psLines.push("1338 ?        00:00:00 [worker.py] <defunct>");
							psLines.push("1339 ?        00:00:00 [worker.py] <defunct>");
							psLines.push("1340 ?        00:00:00 [worker.py] <defunct>");
							psLines.push("1341 ?        00:00:00 [worker.py] <defunct>");
							psLines.push("1342 ?        00:00:00 [worker.py] <defunct>");
						}
						psLines.push("941  pts/0    00:00:00 ps " + args.slice(1).join(" "));
						out = psLines.join("\n");
						break;

					case "kill":
						if (args.length < 2) {
							out = "kill: usage: kill [-s sigspec | -n signum | -sigspec] pid | jobspec ... or kill -l [sigspec]";
						} else {
							var targetPID = parseInt(args[args.length - 1], 10);
							if (targetPID === 4192 && DT.state.activeIncidentState.port8080PID === 4192) {
								DT.state.activeIncidentState.port8080PID = null;
								out = "[4192]+  Terminated              python3 -m http.server 8080";
							} else if (targetPID === 1337 && DT.state.activeIncidentState.zombieParentPID === 1337) {
								DT.state.activeIncidentState.zombieParentPID = null;
								out = "[1337]+  Killed                  python3 bad_supervisor.py\nInit (PID 1) reaped 5 defunct child processes.";
							} else if (targetPID === 1) {
								out = "kill: (1) - Operation not permitted";
							} else if (!targetPID) {
								out = "kill: invalid process id";
							} else {
								out = "kill: (" + targetPID + ") - No such process";
							}
						}
						break;

					case "top":
					case "htop":
						out = "top - 16:08:00 up 42 days,  1 user,  load average: 0.14, 0.08, 0.05\n" +
							"Tasks: 184 total,   1 running, 183 sleeping,   0 stopped,   0 zombie\n" +
							"%Cpu(s):  1.2 us,  0.4 sy,  0.0 ni, 98.2 id,  0.1 wa,  0.0 hi,  0.1 si\n" +
							"MiB Mem :  32000.0 total,  17812.4 free,   7654.1 used,   6533.5 buff/cache\n" +
							"MiB Swap:   8192.0 total,   8192.0 free,      0.0 used.  24345.9 avail Mem\n\n" +
							"  PID USER      PR  NI    VIRT    RES    SHR S  %CPU  %MEM     TIME+ COMMAND\n" +
							"  450 visitor   20   0  712480  42120  18400 S   0.8   0.1   0:08.42 portfolio\n" +
							"  312 caddy     20   0   45120  12450   8920 S   0.2   0.0   0:01.18 caddy\n" +
							"  240 postgres  20   0  289400  38920  24500 S   0.1   0.1   0:00.94 postgres\n" +
							"    1 root      20   0  168240  14120   9200 S   0.0   0.0   0:02.14 systemd";
						break;

					case "df":
						var showInodes = false;
						for (var i = 1; i < args.length; i++) {
							if (args[i] === "-i" || args[i] === "--inodes") showInodes = true;
						}
						if (showInodes) {
							if (DT.state.activeIncidentState.inodesFull) {
								out = "Filesystem       Inodes   IUsed   IFree IUse% Mounted on\n" +
									"devtmpfs        2048000     450 2047550    1% /dev\n" +
									"tmpfs           2048000       1 2047999    1% /run\n" +
									"/dev/nvme0n1p2  6553600 6553600       0  100% /\n" +
									"tmpfs           2048000       4 2047996    1% /dev/shm\n" +
									"/dev/nvme0n1p1   131072     320  130752    1% /boot/efi";
							} else {
								out = "Filesystem       Inodes   IUsed   IFree IUse% Mounted on\n" +
									"devtmpfs        2048000     450 2047550    1% /dev\n" +
									"tmpfs           2048000       1 2047999    1% /run\n" +
									"/dev/nvme0n1p2  6553600  421800 6131800    7% /\n" +
									"tmpfs           2048000       4 2047996    1% /dev/shm\n" +
									"/dev/nvme0n1p1   131072     320  130752    1% /boot/efi";
							}
						} else {
							if (DT.state.activeIncidentState.inodesFull) {
								out = "Filesystem      Size  Used Avail Use% Mounted on\n" +
									"devtmpfs         16G     0   16G   0% /dev\n" +
									"tmpfs            16G  240M   16G   2% /run\n" +
									"/dev/nvme0n1p2  468G  178G  267G  38% /\n" +
									"tmpfs            16G  4.0K   16G   1% /dev/shm\n" +
									"/dev/nvme0n1p1  511M   32M  480M   7% /boot/efi";
							} else {
								out = "Filesystem      Size  Used Avail Use% Mounted on\n" +
									"devtmpfs         16G     0   16G   0% /dev\n" +
									"tmpfs            16G  240M   16G   2% /run\n" +
									"/dev/nvme0n1p2  468G   48G  397G  11% /\n" +
									"tmpfs            16G  4.0K   16G   1% /dev/shm\n" +
									"/dev/nvme0n1p1  511M   32M  480M   7% /boot/efi\n" +
									"tmpfs           3.2G   64K  3.2G   1% /run/user/1000";
							}
						}
						break;

					case "du":
						var duDir = args[1] ? DT.fs.resolvePath(args[1]) : DT.state.currentPath;
						out = "4.0K\t" + duDir + "/README.md\n" +
							"4.0K\t" + duDir + "/about.txt\n" +
							"4.0K\t" + duDir + "/contact.txt\n" +
							"348K\t" + duDir + "/posts\n" +
							"48K\t" + duDir + "/projects\n" +
							"412K\t" + duDir;
						break;

					case "free":
						out = "               total        used        free      shared  buff/cache   available\n" +
							"Mem:           31.2G        7.4G       17.3G        240M        6.5G       23.7G\n" +
							"Swap:           8.0G          0B        8.0G";
						break;

					case "env":
					case "printenv":
						if (args[1]) {
							out = DT.state.envVars[args[1]] || "";
						} else {
							var envList = [];
							for (var k in DT.state.envVars) envList.push(k + "=" + DT.state.envVars[k]);
							out = envList.join("\n");
						}
						break;

					case "export":
						if (args.length < 2) {
							for (var k in DT.state.envVars) out += "declare -x " + k + "=\"" + DT.state.envVars[k] + "\"\n";
						} else {
							var pair = args.slice(1).join(" ");
							var eqIdx = pair.indexOf("=");
							if (eqIdx !== -1) {
								var varK = pair.substring(0, eqIdx).trim();
								var varV = pair.substring(eqIdx + 1).trim().replace(/^["']|["']$/g, "");
								DT.state.envVars[varK] = varV;
							}
						}
						break;

					case "id":
						out = "uid=1000(visitor) gid=1000(visitor) groups=1000(visitor),4(adm),24(cdrom),27(sudo),30(dip),46(plugdev),100(users)";
						break;

					case "groups":
						out = "visitor adm cdrom sudo dip plugdev users";
						break;

					case "whoami":
						out = DT.state.envVars["USER"] || "visitor";
						break;

					case "hostname":
						out = DT.state.envVars["HOSTNAME"] || "daemontalk.local";
						break;

					case "uname":
						var all = false;
						for (var i = 1; i < args.length; i++) if (args[i] === "-a") all = true;
						out = all ? "Linux daemontalk 6.12.8-generic #1 SMP PREEMPT_DYNAMIC Sat Aug 8 16:00:00 WIB 2026 x86_64 GNU/Linux" : "Linux";
						break;

					case "arch":
						out = DT.state.envVars["ARCH"] || "x86_64";
						break;

					case "date":
						out = new Date().toString();
						break;

					case "uptime":
						out = " 16:08:12 up 42 days,  3:14,  1 user,  load average: 0.14, 0.08, 0.05";
						break;

					case "dmesg":
						out = DT.fs.vfs["/var/log/dmesg.log"].content || "";
						break;

					case "ping":
						if (args.length < 2) {
							out = "ping: usage: ping [-c count] destination";
						} else {
							var host = args[args.length - 1];
							var count = 4;
							for (var i = 1; i < args.length; i++) {
								if (args[i] === "-c" && args[i+1]) { count = parseInt(args[i+1], 10) || 4; i++; }
							}
							var pLines = ["PING " + host + " (1.1.1.1) 56(84) bytes of data."];
							for (var seq = 1; seq <= Math.min(6, count); seq++) {
								var rtt = (12 + Math.random() * 6).toFixed(2);
								pLines.push("64 bytes from " + host + ": icmp_seq=" + seq + " ttl=56 time=" + rtt + " ms");
							}
							pLines.push("\n--- " + host + " ping statistics ---");
							pLines.push(count + " packets transmitted, " + count + " received, 0% packet loss, time " + (count * 1000) + "ms");
							out = pLines.join("\n");
						}
						break;

					case "dig":
					case "nslookup":
						var domain = args[1] || "daemontalk.com";
						out = "; <<>> DiG 9.18.28 <<>> " + domain + "\n" +
							";; global options: +cmd\n" +
							";; Got answer:\n" +
							";; ->>HEADER<<- opcode: QUERY, status: NOERROR, id: " + Math.floor(Math.random()*60000) + "\n" +
							";; flags: qr rd ra; QUERY: 1, ANSWER: 2, AUTHORITY: 0, ADDITIONAL: 1\n\n" +
							";; QUESTION SECTION:\n" +
							";" + domain + ".\t\t\tIN\tA\n\n" +
							";; ANSWER SECTION:\n" +
							domain + ".\t\t300\tIN\tA\t104.21.72.19\n" +
							domain + ".\t\t300\tIN\tA\t172.67.182.41\n\n" +
							";; Query time: 18 msec\n" +
							";; SERVER: 1.1.1.1#53(1.1.1.1)\n" +
							";; WHEN: Sat Aug 08 16:08:12 WIB 2026\n" +
							";; MSG SIZE  rcvd: 88";
						break;

					case "ip":
					case "ifconfig":
						out = "1: lo: <LOOPBACK,UP,LOWER_UP> mtu 65536 qdisc noqueue state UNKNOWN group default qlen 1000\n" +
							"    link/loopback 00:00:00:00:00:00 brd 00:00:00:00:00:00\n" +
							"    inet 127.0.0.1/8 scope host lo\n" +
							"       valid_lft forever preferred_lft forever\n" +
							"2: eth0: <BROADCAST,MULTICAST,UP,LOWER_UP> mtu 1500 qdisc fq_codel state UP group default qlen 1000\n" +
							"    link/ether 52:54:00:89:a1:fe brd ff:ff:ff:ff:ff:ff\n" +
							"    inet 192.168.1.105/24 brd 192.168.1.255 scope global dynamic eth0\n" +
							"       valid_lft 86240sec preferred_lft 86240sec\n" +
							"3: tailscale0: <POINTOPOINT,MULTICAST,NOARP,UP,LOWER_UP> mtu 1280 qdisc fq_codel state UNKNOWN\n" +
							"    link/none \n" +
							"    inet 100.84.120.45/32 scope global tailscale0";
						break;

					case "ss":
					case "netstat":
					case "lsof":
						if (DT.state.activeIncidentState.port8080PID !== null) {
							out = "Netid State  Recv-Q Send-Q Local Address:Port  Peer Address:Port Process\n" +
								"tcp   LISTEN 0      128          0.0.0.0:22         0.0.0.0:*    users:((\"sshd\",pid=104,fd=3))\n" +
								"tcp   LISTEN 0      511          0.0.0.0:80         0.0.0.0:*    users:((\"caddy\",pid=312,fd=7))\n" +
								"tcp   LISTEN 0      511          0.0.0.0:443        0.0.0.0:*    users:((\"caddy\",pid=312,fd=8))\n" +
								"tcp   LISTEN 0      128          0.0.0.0:5432       0.0.0.0:*    users:((\"postgres\",pid=240,fd=5))\n" +
								"tcp   LISTEN 0      128        127.0.0.1:8080       0.0.0.0:*    users:((\"python3\",pid=" + DT.state.activeIncidentState.port8080PID + ",fd=4))";
						} else {
							out = "Netid State  Recv-Q Send-Q Local Address:Port  Peer Address:Port Process\n" +
								"tcp   LISTEN 0      128          0.0.0.0:22         0.0.0.0:*    users:((\"sshd\",pid=104,fd=3))\n" +
								"tcp   LISTEN 0      511          0.0.0.0:80         0.0.0.0:*    users:((\"caddy\",pid=312,fd=7))\n" +
								"tcp   LISTEN 0      511          0.0.0.0:443        0.0.0.0:*    users:((\"caddy\",pid=312,fd=8))\n" +
								"tcp   LISTEN 0      128          0.0.0.0:5432       0.0.0.0:*    users:((\"postgres\",pid=240,fd=5))\n" +
								"tcp   LISTEN 0      128        127.0.0.1:3000       0.0.0.0:*    users:((\"portfolio\",pid=450,fd=3))";
						}
						break;

					case "incident":
					case "warstory":
					case "rca":
					case "challenge":
						var sub = args[1] ? args[1].toLowerCase() : "list";
						if (sub === "list" || sub === "ls") {
							var listOut = "╔══════════════════════════════════════════════════════════════════════════════╗\n" +
								"║               LINUX WAR STORIES & INCIDENT RCA CHALLENGES                    ║\n" +
								"║  Investigate real production failures using standard Linux diagnostic tools. ║\n" +
								"╚══════════════════════════════════════════════════════════════════════════════╝\n\n";
							DT.fs.incidents.forEach(function(inc) {
								var isSolved = DT.state.activeIncidentState.solvedList.indexOf(inc.id) !== -1;
								var statusBadge = isSolved ? "[✓ SOLVED]" : "[PENDING ]";
								listOut += "  " + statusBadge + " incident " + inc.id + " (" + inc.difficulty + ")\n" +
									"             Title: " + inc.title + "\n\n";
							});
							listOut += "Commands:\n" +
								"  incident <id>     : Start an incident scenario (e.g. `incident 1`)\n" +
								"  incident hint     : Get troubleshooting hint for current incident\n" +
								"  incident verify   : Verify fix and complete challenge\n" +
								"  incident status   : Show score and progress\n";
							out = listOut;
						} else if (sub === "status") {
							var solvedCount = DT.state.activeIncidentState.solvedList.length;
							out = "Incident Progress: " + solvedCount + " / " + DT.fs.incidents.length + " Solved\n" +
								"Active Incident: " + (DT.state.activeIncidentState.currentId ? "Incident " + DT.state.activeIncidentState.currentId : "None (type `incident 1` to start)");
						} else if (sub === "hint") {
							if (!DT.state.activeIncidentState.currentId) {
								out = "No incident active. Type `incident 1` to start a challenge.";
							} else {
								var curInc = DT.fs.incidents.find(function(i) { return i.id === DT.state.activeIncidentState.currentId; });
								out = "💡 HINT (Incident " + curInc.id + "):\n" + curInc.hint;
							}
						} else if (sub === "verify" || sub === "solve" || sub === "check") {
							if (!DT.state.activeIncidentState.currentId) {
								out = "No active incident. Type `incident 1` to start.";
							} else {
								var activeInc = DT.fs.incidents.find(function(i) { return i.id === DT.state.activeIncidentState.currentId; });
								var answerParam = args.slice(2).join(" ");
								var checkRes = activeInc.verify(answerParam);
								if (checkRes.solved) {
									if (DT.state.activeIncidentState.solvedList.indexOf(activeInc.id) === -1) {
										DT.state.activeIncidentState.solvedList.push(activeInc.id);
										localStorage.setItem("daemontalk_solved_incidents", JSON.stringify(DT.state.activeIncidentState.solvedList));
									}
									out = "╔══════════════════════════════════════════════════════════════════════════════╗\n" +
										"║  🎉 CHALLENGE SOLVED: " + activeInc.title.toUpperCase() + "\n" +
										"╚══════════════════════════════════════════════════════════════════════════════╝\n\n" +
										checkRes.message + "\n\n" +
										"Progress: " + DT.state.activeIncidentState.solvedList.length + "/" + DT.fs.incidents.length + " completed.\n" +
										"Type `incident list` to select the next incident!";
								} else {
									out = "❌ NOT RESOLVED YET\n" + checkRes.message + "\n\nType `incident hint` if you need guidance.";
								}
							}
						} else {
							var incId = parseInt(sub, 10);
							var targetInc = DT.fs.incidents.find(function(i) { return i.id === incId; });
							if (targetInc) {
								DT.state.activeIncidentState.currentId = targetInc.id;
								targetInc.setup();
								out = "╔══════════════════════════════════════════════════════════════════════════════╗\n" +
									"║  🚨 INCIDENT #" + targetInc.id + ": " + targetInc.title.toUpperCase() + "\n" +
									"║  Difficulty: " + targetInc.difficulty + "\n" +
									"╚══════════════════════════════════════════════════════════════════════════════╝\n\n" +
									"SYMPTOM:\n" + targetInc.symptom + "\n\n" +
									"YOUR MISSION:\nInvestigate with terminal tools, resolve the root cause, then run `incident verify`.\n(Type `incident hint` for guidance).";
							} else {
								out = "incident: unknown argument '" + sub + "'. Type `incident list` to view all challenges.";
							}
						}
						break;

					case "recipes":
					case "cheat":
					case "ebpf":
						out = "╔══════════════════════════════════════════════════════════════════════════════╗\n" +
							"║                  PRODUCTION LINUX & eBPF ONE-LINER RECIPES                   ║\n" +
							"╚══════════════════════════════════════════════════════════════════════════════╝\n\n" +
							"1. Trace Top System Calls by Process in Real-Time (bpftrace)\n" +
							"   sudo bpftrace -e 'tracepoint:raw_syscalls:sys_enter { @[comm] = count(); }'\n\n" +
							"2. Monitor Every File Opened Across Entire Linux System\n" +
							"   sudo bpftrace -e 'tracepoint:syscalls:sys_enter_openat { printf(\"%s opened %s\\n\", comm, str(args->filename)); }'\n\n" +
							"3. Identify Which Process is Generating Heavy Disk Write I/O\n" +
							"   sudo pidstat -d 1 5   ||   sudo iotop -oPa\n\n" +
							"4. Trace Dropped TCP Packets Inside Kernel Subsystems\n" +
							"   sudo bpftrace -e 'kprobe:kfree_skb { @[kstack] = count(); }'\n\n" +
							"5. Find Large Deleted Files Still Held Open by Processes (Reclaim Space)\n" +
							"   sudo lsof +L1 | grep deleted\n\n" +
							"6. Inspect Inode Starvation (when 'No space left on device' but disk 50% free)\n" +
							"   df -i\n" +
							"   sudo find /var -xdev -printf '%h\\n' | sort | uniq -c | sort -k1 -nr | head -n 10\n\n" +
							"7. Inspect Real Socket Listeners & Buffer Drops\n" +
							"   ss -tulpn   ||   netstat -s | grep -i drop\n\n" +
							"Tip: You can also access these from your local terminal with `curl -sL " + window.location.host + "/recipes`";
						break;

					case "which":
						if (args.length < 2) out = "which: missing argument";
						else {
							var cmdName = args[1];
							if (DT.fs.vfs["/bin"].children.indexOf(cmdName) !== -1) {
								out = "/bin/" + cmdName;
							} else {
								out = cmdName + " not found in PATH";
							}
						}
						break;

					case "alias":
						if (args.length < 2) {
							var aList = [];
							for (var k in DT.state.aliases) aList.push("alias " + k + "='" + DT.state.aliases[k] + "'");
							out = aList.join("\n") || "alias: no DT.state.aliases defined";
						} else {
							var aStr = args.slice(1).join(" ");
							var aEq = aStr.indexOf("=");
							if (aEq !== -1) {
								var aK = aStr.substring(0, aEq).trim();
								var aV = aStr.substring(aEq + 1).trim().replace(/^['"]|['"]$/g, "");
								DT.state.aliases[aK] = aV;
							}
						}
						break;

					case "unalias":
						if (args[1] && DT.state.aliases[args[1]]) {
							delete DT.state.aliases[args[1]];
						}
						break;

					case "cal":
						var d = new Date();
						var months = ["January", "February", "March", "April", "May", "June", "July", "August", "September", "October", "November", "December"];
						var monthName = months[d.getMonth()];
						var year = d.getFullYear();
						var header = "     " + monthName + " " + year;
						var calGrid = [
							header,
							"Su Mo Tu We Th Fr Sa",
							"                   1",
							" 2  3  4  5  6  7  8",
							" 9 10 11 12 13 14 15",
							"16 17 18 19 20 21 22",
							"23 24 25 26 27 28 29",
							"30 31"
						];
						out = calGrid.join("\n");
						break;

					case "yes":
						var yText = args[1] || "y";
						var yLines = [];
						for (var i = 0; i < 12; i++) yLines.push(yText);
						out = yLines.join("\n");
						break;

					case "history":
						out = DT.state.cmdHistory.map(function(c, i) { return "  " + (i + 1) + "  " + c; }).join("\n");
						break;

					case "clear":
						window.termClear();
						return null;

					case "exit":
						window.location.href = "/";
						return null;

					default:
						out = cmd + ": command not found. Type 'help' for a list of available commands.";
						break;
				}

				return out;
			}

			// Main Command Dispatcher with Piping and Redirection
			
	DT.cmd.execute = function (rawCmd) {
				var trimmed = rawCmd.trim();
				if (!trimmed) {
					DT.ui.appendHistory(rawCmd, "");
					return;
				}

				DT.state.cmdHistory.push(rawCmd);
				DT.state.histIndex = DT.state.cmdHistory.length;

				// Handle async special commands like curl / run
				var firstWord = trimmed.split(/\s+/)[0].toLowerCase();
				if (firstWord === "curl") {
					var cArgs = DT.cmd.parseArgs(trimmed);
					if (cArgs.length < 2) {
						DT.ui.appendHistory(rawCmd, "curl: try 'curl <url>' (e.g. curl /api/terminal/data)");
						return;
					}
					var url = cArgs[1];
					DT.ui.appendHistory(rawCmd, "Connecting to " + url + "...");
					fetch(url)
						.then(function(r) { return r.text(); })
						.then(function(text) {
							DT.ui.appendHistory("", text.substring(0, 1500) + (text.length > 1500 ? "\n... [truncated]" : ""));
						})
						.catch(function(err) {
							DT.ui.appendHistory("", "curl: (7) Failed to connect: " + err.toString());
						});
					return;
				}

				if (firstWord === "run") {
					var rArgs = DT.cmd.parseArgs(trimmed);
					if (rArgs.length < 3) {
						DT.ui.appendHistory(rawCmd, "Usage: run <go|js|python|bash> <code>\nExample: run go fmt.Println(\"Hello from daemontalk\")");
						return;
					}
					var lang = rArgs[1];
					var code = rArgs.slice(2).join(" ");
					DT.ui.appendHistory(rawCmd, "Executing code (" + lang + ")...");
					fetch("/api/run", {
						method: "POST",
						headers: { "Content-Type": "application/json" },
						body: JSON.stringify({ language: lang, code: code })
					})
					.then(function(r) { return r.json(); })
					.then(function(data) {
						var result = "";
						if (data.stdout) result += data.stdout;
						if (data.stderr) result += (result ? "\n" : "") + "[stderr] " + data.stderr;
						if (data.error) result += (result ? "\n" : "") + "[error] " + data.error;
						if (!result) result = "Program completed with exit code 0 (no output)";
						DT.ui.appendHistory("", result);
					})
					.catch(function(err) {
						DT.ui.appendHistory("", "Execution error: " + err.toString());
					});
					return;
				}

				// Check Redirection
				var redirectAppend = false;
				var redirectFile = null;
				var workCmd = trimmed;

				if (workCmd.indexOf(">>") !== -1) {
					var rParts = workCmd.split(">>");
					workCmd = rParts[0].trim();
					redirectFile = rParts[1].trim();
					redirectAppend = true;
				} else if (workCmd.indexOf(">") !== -1) {
					var rParts = workCmd.split(">");
					workCmd = rParts[0].trim();
					redirectFile = rParts[1].trim();
					redirectAppend = false;
				}

				// Check Piping
				var pipeSegments = workCmd.split("|").map(function(s) { return s.trim(); });
				var currentPipeOut = null;

				for (var i = 0; i < pipeSegments.length; i++) {
					var segment = pipeSegments[i];
					var segOut = DT.cmd.evalSingleCommand(segment, currentPipeOut);
					if (segOut === null) return; // handled by clear/exit
					currentPipeOut = segOut;
				}

				var finalOut = currentPipeOut || "";

				// Handle File Redirection
				if (redirectFile && finalOut !== null) {
					var rfPath = DT.fs.resolvePath(redirectFile);
					var rfParent = rfPath.substring(0, rfPath.lastIndexOf("/")) || "/";
					var rfName = rfPath.split("/").pop();
					if (DT.fs.vfs[rfParent] && DT.fs.vfs[rfParent].type === "dir") {
						if (!DT.fs.vfs[rfPath]) {
							DT.fs.vfs[rfPath] = { type: "file", content: finalOut + "\n" };
							if (DT.fs.vfs[rfParent].children.indexOf(rfName) === -1) DT.fs.vfs[rfParent].children.push(rfName);
						} else if (redirectAppend) {
							DT.fs.vfs[rfPath].content = (DT.fs.vfs[rfPath].content || "") + finalOut + "\n";
						} else {
							DT.fs.vfs[rfPath].content = finalOut + "\n";
						}
						finalOut = "";
					} else {
						finalOut = "bash: " + redirectFile + ": No such directory";
					}
				}

				DT.ui.appendHistory(rawCmd, finalOut);
			}

			
	DT.cmd.parseArgs = function (str) {
				var match = str.match(/(?:[^\s"']+|"[^"]*"|'[^']*')+/g);
				if (!match) return [];
				return match.map(function(m) {
					if ((m.startsWith('"') && m.endsWith('"')) || (m.startsWith("'") && m.endsWith("'"))) {
						return m.substring(1, m.length - 1);
					}
					return m;
				});
			}

			
})();
