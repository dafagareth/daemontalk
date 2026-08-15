(function() {
	function initDaemonTalkTerminal() {
		var screenEl = document.getElementById("term-screen");
		var historyEl = document.getElementById("term-history");
		var inputEl = document.getElementById("term-input");
		var promptLabel = document.getElementById("term-prompt-label");
		var titleText = document.getElementById("term-title-text");

		if (!screenEl || !inputEl || !historyEl) return;
		if (inputEl.dataset.initialized === "true") return;
		inputEl.dataset.initialized = "true";

			var cmdHistory = [];
			var histIndex = -1;
			var currentInputBuffer = "";

			var currentPath = "/home/visitor";
			var envVars = {
				"USER": "visitor",
				"HOME": "/home/visitor",
				"SHELL": "/bin/bash",
				"PATH": "/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin",
				"TERM": "xterm-256color",
				"LANG": "en_US.UTF-8",
				"HOSTNAME": "daemontalk.local",
				"PWD": "/home/visitor",
				"EDITOR": "nano",
				"ARCH": "x86_64"
			};
			var aliases = {};

			var activeIncidentState = {
				currentId: null,
				inodesFull: false,
				port8080PID: null,
				zombieParentPID: null,
				solvedList: JSON.parse(localStorage.getItem("daemontalk_solved_incidents") || "[]")
			};

			var incidents = [
				{
					id: 1,
					title: "The Mysterious 'No Space Left on Device'",
					difficulty: "Medium",
					symptom: "Production log pipeline fails with 'ENOSPC: No space left on device'. However, `df -h` shows the disk is only 38% full! Find the root cause and fix it.",
					hint: "Check filesystem metadata limits: compare `df -h` (block usage) with `df -i` (inode usage). Inspect /var/spool/clientmqueue for abandoned zero-byte files.",
					setup: function() {
						vfs["/var/spool"] = { type: "dir", children: ["clientmqueue"] };
						vfs["/var/spool/clientmqueue"] = { type: "dir", children: ["msg_10482.tmp", "msg_10483.tmp", "msg_10484.tmp", "msg_10485.tmp", "msg_10486.tmp", "abandoned_tokens.lock"] };
						vfs["/var/spool/clientmqueue/msg_10482.tmp"] = { type: "file", content: "" };
						vfs["/var/spool/clientmqueue/abandoned_tokens.lock"] = { type: "file", content: "LOCKED_INODES=6553600" };
						activeIncidentState.inodesFull = true;
					},
					verify: function() {
						if (!activeIncidentState.inodesFull || !vfs["/var/spool/clientmqueue"] || vfs["/var/spool/clientmqueue"].children.length === 0) {
							activeIncidentState.inodesFull = false;
							return { solved: true, message: "ROOT CAUSE CONFIRMED: Inode exhaustion! You cleaned up abandoned queue files, freeing 6.5M inodes. Filesystem writes restored." };
						}
						return { solved: false, message: "Issue still persists: Inodes on /dev/nvme0n1p2 are still 100% full. Try inspecting /var/spool/clientmqueue and cleaning up files with `rm -rf /var/spool/clientmqueue/*` or `incident solve inode`." };
					}
				},
				{
					id: 2,
					title: "The Rogue Ghost Daemon on Port 8080",
					difficulty: "Easy",
					symptom: "New API deployment crashed on startup with `listen tcp :8080: bind: address already in use`. Find the culprit process holding port 8080 and terminate it.",
					hint: "Use socket inspection tools: `ss -tulpn` or `lsof -i :8080` to locate the PID, then use `kill <PID>`.",
					setup: function() {
						activeIncidentState.port8080PID = 4192;
					},
					verify: function() {
						if (activeIncidentState.port8080PID === null) {
							return { solved: true, message: "PORT RELEASED: PID 4192 terminated. Port 8080 is now free for new deployments." };
						}
						return { solved: false, message: "Port 8080 is still locked by PID " + activeIncidentState.port8080PID + ". Use `ss -tulpn` or `lsof -i :8080` to find it, then `kill " + activeIncidentState.port8080PID + "`." };
					}
				},
				{
					id: 3,
					title: "The Silent Midnight Crash (OOM-Killer)",
					difficulty: "Medium",
					symptom: "Redis cache server abruptly disappeared at 03:14:02 UTC with exit code 137. Application logs show no stack trace. Identify why Linux terminated the process.",
					hint: "Kernel-level kills do not appear in user application logs. Check kernel ring buffer logs with `dmesg` or `dmesg | grep -i oom`.",
					setup: function() {
						vfs["/var/log/dmesg.log"].content += "\n[14829.102390] Out of memory: Killed process 8192 (redis-server) total-vm:4194304kB, anon-rss:2048512kB, file-rss:0kB\n[14829.102401] oom_reaper: reaped process 8192 (redis-server), now anon-rss:0kB\n";
					},
					verify: function(answer) {
						if (answer && (answer.toLowerCase().indexOf("oom") !== -1 || answer.toLowerCase().indexOf("out of memory") !== -1 || answer.toLowerCase().indexOf("memory") !== -1)) {
							return { solved: true, message: "DIAGNOSIS CORRECT: Linux kernel OOM-Killer sacrificed redis-server (PID 8192) due to cgroup memory exhaustion. Recommendation: Adjust vm.overcommit_memory=1 and configure maxmemory in redis.conf." };
						}
						return { solved: false, message: "Check kernel buffer messages with `dmesg` or `cat /var/log/dmesg.log | grep -i oom`, then type `incident solve oom` to submit your diagnosis." };
					}
				},
				{
					id: 4,
					title: "Zombie Apocalypse & Fork Starvation",
					difficulty: "Hard",
					symptom: "System throwing `fork: retry: Resource temporarily unavailable`. The PID space is flooded with defunct child processes. Identify and kill the rogue parent supervisor.",
					hint: "Inspect process states with `ps aux` or `ps -ef`. Notice which parent PID (PPID) spawned all the defunct workers, then kill that parent PID.",
					setup: function() {
						activeIncidentState.zombieParentPID = 1337;
					},
					verify: function() {
						if (activeIncidentState.zombieParentPID === null) {
							return { solved: true, message: "ZOMBIES REAPED: Parent supervisor PID 1337 killed. Init (PID 1) adopted and successfully reaped all defunct children. PID table capacity restored." };
						}
						return { solved: false, message: "Dozens of zombie processes still exist. Run `ps aux` to locate the rogue parent process (PID 1337 bad_supervisor.py) and execute `kill -9 1337`." };
					}
				},
				{
					id: 5,
					title: "The DNS Blackhole",
					difficulty: "Easy",
					symptom: "Microservices cannot communicate: `curl api.internal` or `ping backend.local` fails with `Could not resolve host`. Fix the resolver configuration.",
					hint: "Inspect `/etc/resolv.conf`. Look for invalid IP addresses, and write valid DNS servers like `nameserver 1.1.1.1` or `nameserver 8.8.8.8`.",
					setup: function() {
						vfs["/etc/resolv.conf"].content = "# Misconfigured by faulty provisioning script\nnameserver 192.168.1.999\nnameserver 0.0.0.0\n";
					},
					verify: function() {
						var conf = vfs["/etc/resolv.conf"] ? vfs["/etc/resolv.conf"].content : "";
						if (conf.indexOf("1.1.1.1") !== -1 || conf.indexOf("8.8.8.8") !== -1 || conf.indexOf("127.0.0.53") !== -1) {
							return { solved: true, message: "DNS RESOLUTION RESTORED: /etc/resolv.conf repaired with valid nameservers. Hostname resolution functional." };
						}
						return { solved: false, message: "Resolver configuration in /etc/resolv.conf is still invalid. Inspect with `cat /etc/resolv.conf` and update it with `echo 'nameserver 1.1.1.1' > /etc/resolv.conf`." };
					}
				}
			];

			var vfs = {
				"/": { type: "dir", children: ["home", "bin", "etc", "var", "proc", "tmp", "dev"] },
				"/bin": {
					type: "dir",
					children: [
						"alias", "arch", "awk", "base64", "basename", "cal", "cat", "cd", "challenge", "cheat", "clear", "cp", "curl", "cut",
						"date", "df", "diff", "dig", "dirname", "dmesg", "du", "ebpf", "echo", "env", "exit", "export", "file",
						"find", "free", "grep", "groups", "head", "help", "history", "hostname", "htop", "id", "ifconfig", "incident",
						"ip", "kill", "less", "ls", "lsof", "man", "mkdir", "more", "mv", "netstat", "nslookup", "ping", "printenv", "ps",
						"pwd", "rca", "recipes", "rev", "rm", "run", "sed", "seq", "sleep", "sort", "ss", "stat", "tail", "tar", "top",
						"touch", "tr", "tree", "uname", "unalias", "uniq", "uptime", "warstory", "wc", "which", "whoami", "yes"
					]
				},
				"/etc": { type: "dir", children: ["os-release", "hostname", "passwd", "hosts", "resolv.conf"] },
				"/etc/os-release": { type: "file", content: "NAME=\"daemontalk Linux\"\nVERSION=\"2026.8 LTS (Codename Tumbleweed)\"\nID=daemontalk\nID_LIKE=arch\nPRETTY_NAME=\"daemontalk Linux 6.12.8 LTS\"\nHOME_URL=\"https://daemontalk.com\"\n" },
				"/etc/hostname": { type: "file", content: "daemontalk.local\n" },
				"/etc/passwd": { type: "file", content: "root:x:0:0:root:/root:/bin/bash\ndaemon:x:1:1:daemon:/usr/sbin:/usr/sbin/nologin\nbin:x:2:2:bin:/bin:/usr/sbin/nologin\nsys:x:3:3:sys:/dev:/usr/sbin/nologin\nsync:x:4:65534:sync:/bin:/bin/sync\nvisitor:x:1000:1000:Visitor User,,,:/home/visitor:/bin/bash\npostgres:x:114:120:PostgreSQL administrator,,,:/var/lib/postgresql:/bin/bash\ncaddy:x:998:998:Caddy web server:/var/lib/caddy:/usr/sbin/nologin\n" },
				"/etc/hosts": { type: "file", content: "127.0.0.1\tlocalhost\n127.0.1.1\tdaemontalk.local daemontalk\n::1\t\tlocalhost ip6-localhost ip6-loopback\nff02::1\t\tip6-allnodes\nff02::2\t\tip6-allrouters\n" },
				"/etc/resolv.conf": { type: "file", content: "# Generated by systemd-resolved\nnameserver 1.1.1.1\nnameserver 8.8.8.8\nsearch daemontalk.local\n" },
				"/proc": { type: "dir", children: ["cpuinfo", "meminfo", "version", "uptime", "loadavg"] },
				"/proc/cpuinfo": { type: "file", content: "processor\t: 0\nvendor_id\t: AuthenticAMD\ncpu family\t: 25\nmodel\t\t: 33\nmodel name\t: AMD Ryzen 9 5950X 16-Core Processor\nstepping\t: 0\nmicrocode\t: 0xa201016\ncpu MHz\t\t: 3400.000\ncache size\t: 512 KB\nflags\t\t: fpu vme de pse tsc msr pae mce cx8 apic sep mtrr pge mca cmov pat pse36 clflush mmx fxsr sse sse2 ht syscall nx mmxext fxsr_opt pdpe1gb rdtscp lm constant_tsc rep_good nopl nonstop_tsc cpuid extd_apicid aperfmperf rapl pni pclmulqdq monitor ssse3 fma cx16 sse4_1 sse4_2 movbe popcnt aes xsave avx f16c rdrand hypervisor lahf_lm cmp_legacy svm cr8_legacy abm sse4a misalignsse 3dnowprefetch osvw topoext perfctr_core invpcid ibpb vmmcall fsgsbase bmi1 avx2 smep bmi2 erms invpcid_single rdseed adx smap clflushopt clwb sha_ni xsaveopt xsavec xgetbv1 xsaves clzero xsaveerptr rdpru wbnoinvd arat\nbogomips\t: 6799.98\nTLB size\t: 2560 4K pages\nclflush size\t: 64\ncache_alignment\t: 64\naddress sizes\t: 48 bits physical, 48 bits virtual\n" },
				"/proc/meminfo": { type: "file", content: "MemTotal:       32768000 kB\nMemFree:        18245120 kB\nMemAvailable:   24912040 kB\nBuffers:          541200 kB\nCached:          6125720 kB\nSwapCached:            0 kB\nActive:          7845120 kB\nInactive:        4821000 kB\nSwapTotal:       8388608 kB\nSwapFree:        8388608 kB\nDirty:               120 kB\nWriteback:             0 kB\nAnonPages:       6000400 kB\nMapped:           892400 kB\nShmem:            240100 kB\nSlab:             720400 kB\n" },
				"/proc/version": { type: "file", content: "Linux version 6.12.8-generic (buildd@buildhost) (gcc version 14.2.0 (Debian 14.2.0-8)) #1 SMP PREEMPT_DYNAMIC\n" },
				"/proc/uptime": { type: "file", content: "3648120.45 58369927.20\n" },
				"/proc/loadavg": { type: "file", content: "0.14 0.08 0.05 1/412 18920\n" },
				"/var": { type: "dir", children: ["log"] },
				"/var/log": { type: "dir", children: ["syslog.log", "dmesg.log"] },
				"/var/log/syslog.log": { type: "file", content: "systemd[1]: Started daemontalk web server daemon.\nsystemd[1]: Started PostgreSQL database server.\nsystemd[1]: Started Caddy reverse proxy & TLS manager.\nkernel: [ 0.000000] Linux version 6.12.8-generic (root@buildhost) #1 SMP PREEMPT\nkernel: [ 0.000004] Command line: BOOT_IMAGE=/vmlinuz-6.12.8-generic root=UUID=5e18b273 ro quiet splash\n" },
				"/var/log/dmesg.log": { type: "file", content: "[    0.000000] Linux version 6.12.8-generic (gcc version 14.2.0)\n[    0.000000] Command line: BOOT_IMAGE=/boot/vmlinuz-6.12.8-generic root=/dev/nvme0n1p2 rw\n[    0.000000] BIOS-provided physical RAM map:\n[    0.000000] BIOS-e820: [mem 0x0000000000000000-0x000000000009fbff] usable\n[    0.000000] BIOS-e820: [mem 0x0000000000100000-0x00000007ffffffff] usable\n[    0.001420] smpboot: CPU0: AMD Ryzen 9 5950X 16-Core Processor (family: 0x19, model: 0x21, stepping: 0x0)\n[    0.021940] io_uring: enabling fast unprivileged asynchronous I/O subsystem\n[    0.038100] Landlock: LSM initialized successfully\n[    0.042100] Btrfs loaded, crc32c=crc32c-generic\n[    0.081200] Netfilter: nf_conntrack version 0.5.0 (65536 buckets, 262144 max)\n" },
				"/tmp": { type: "dir", children: [] },
				"/dev": { type: "dir", children: ["null", "zero", "urandom", "stdout", "stderr", "stdin"] },
				"/home": { type: "dir", children: ["visitor"] },
				"/home/visitor": { type: "dir", children: ["posts", "projects", "about.txt", "contact.txt", "README.md", ".bashrc"] },
				"/home/visitor/about.txt": { type: "file", content: "Dafa Gareth\nSoftware Developer & Systems Enthusiast\nFocus: Linux internals, Go concurrency, eBPF, distributed storage, and high-performance minimalist web systems.\n" },
				"/home/visitor/contact.txt": { type: "file", content: "GitHub:   https://github.com/dafagareth\nDiscord:  https://discord.gg/daemontalk\nWebsite:  https://daemontalk.com\nEmail:    dafagareth@gmail.com\n" },
				"/home/visitor/README.md": { type: "file", content: "# Welcome to daemontalk Shell\nA fully interactive, educational UNIX sandbox running in your browser.\n\nQuick Tips:\n- Browse articles: `cd posts && ls -l`\n- Read an article: `cat posts/ghostty-terminal-emulator.md`\n- Search topics: `grep -i \"ebpf\" posts/*`\n- Try piping: `ls posts | sort -r | head -n 5`\n- Run code: `run go fmt.Println(\"Hello from daemontalk\")`\n- View manual: `man grep` or `man find`\n" },
				"/home/visitor/.bashrc": { type: "file", content: "# ~/.bashrc: executed by bash(1) for non-login shells.\nexport PS1='\\u@\\h:\\w\\$ '\nalias ll='ls -la'\nalias la='ls -A'\nalias l='ls -CF'\n" },
				"/home/visitor/posts": { type: "dir", children: [] },
				"/home/visitor/projects": { type: "dir", children: [] }
			};

			var rawPosts = [];
			var rawProjects = [];

			// Manuals Dictionary
			var manPages = {
				"ls": "NAME\n    ls - list directory contents\n\nSYNOPSIS\n    ls [OPTION]... [FILE]...\n\nDESCRIPTION\n    List information about the FILEs (the current directory by default).\n    Sort entries alphabetically.\n\nOPTIONS\n    -l     use a long listing format (permissions, size, date)\n    -a     do not ignore entries starting with .\n    -la    combination of -l and -a\n\nEXAMPLES\n    ls\n    ls -l posts\n    ls -la ~",
				"cd": "NAME\n    cd - change the shell working directory\n\nSYNOPSIS\n    cd [DIRECTORY]\n\nDESCRIPTION\n    Change the current directory to DIRECTORY. The default DIRECTORY is the value of the HOME shell variable (~).\n\nEXAMPLES\n    cd posts\n    cd ..\n    cd ~/\n    cd /etc",
				"cat": "NAME\n    cat - concatenate files and print on the standard output\n\nSYNOPSIS\n    cat [OPTION]... [FILE]...\n\nDESCRIPTION\n    Concatenate FILE(s) to standard output. With no FILE, or when FILE is -, read standard input.\n\nEXAMPLES\n    cat about.txt\n    cat posts/ghostty-terminal-emulator.md\n    cat /etc/passwd | grep visitor",
				"grep": "NAME\n    grep - print lines that match patterns\n\nSYNOPSIS\n    grep [OPTION]... PATTERNS [FILE]...\n\nOPTIONS\n    -i     ignore case distinctions in patterns and input data\n    -n     prefix each line of output with the 1-based line number\n    -v     invert the sense of matching, to select non-matching lines\n    -c     print only a count of selected lines per FILE\n\nEXAMPLES\n    grep -i \"ebpf\" posts/*\n    grep visitor /etc/passwd\n    ps aux | grep daemontalk",
				"head": "NAME\n    head - output the first part of files\n\nSYNOPSIS\n    head [OPTION]... [FILE]...\n\nOPTIONS\n    -n NUM     print the first NUM lines instead of the first 10\n\nEXAMPLES\n    head -n 5 posts/duckdb-analytical-queries.md\n    ls posts | head -n 8",
				"tail": "NAME\n    tail - output the last part of files\n\nSYNOPSIS\n    tail [OPTION]... [FILE]...\n\nOPTIONS\n    -n NUM     output the last NUM lines, instead of the last 10\n\nEXAMPLES\n    tail -n 10 /var/log/syslog.log\n    cat posts/jujutsu-vcs-git-compatible.md | tail -n 5",
				"find": "NAME\n    find - search for files in a directory hierarchy\n\nSYNOPSIS\n    find [path...] [expression]\n\nOPTIONS\n    -name PATTERN     Base of file name matches shell pattern PATTERN\n    -type f|d         File is of type: f (file), d (directory)\n\nEXAMPLES\n    find .\n    find posts -name \"*.md\"\n    find /etc",
				"tree": "NAME\n    tree - list contents of directories in a tree-like format\n\nSYNOPSIS\n    tree [OPTION]... [DIRECTORY]\n\nOPTIONS\n    -L LEVEL     Max display depth of the directory tree\n\nEXAMPLES\n    tree\n    tree -L 2 /home/visitor",
				"sort": "NAME\n    sort - sort lines of text files\n\nSYNOPSIS\n    sort [OPTION]... [FILE]...\n\nOPTIONS\n    -r     reverse the result of comparisons\n    -n     compare according to numerical value\n    -u     output only the first of an equal run\n\nEXAMPLES\n    ls posts | sort -r\n    cat numbers.txt | sort -n",
				"uniq": "NAME\n    uniq - report or omit repeated lines\n\nSYNOPSIS\n    uniq [OPTION]... [INPUT [OUTPUT]]\n\nOPTIONS\n    -c     prefix lines by the number of occurrences\n    -d     only print duplicate lines, one for each group\n\nEXAMPLES\n    cat list.txt | sort | uniq -c",
				"cut": "NAME\n    cut - remove sections from each line of files\n\nSYNOPSIS\n    cut OPTION... [FILE]...\n\nOPTIONS\n    -d DELIM     use DELIM instead of TAB for field delimiter\n    -f LIST      select only these fields\n\nEXAMPLES\n    cut -d: -f1 /etc/passwd\n    cat /etc/passwd | cut -d: -f1,7",
				"tr": "NAME\n    tr - translate or delete characters\n\nSYNOPSIS\n    tr [OPTION]... SET1 [SET2]\n\nEXAMPLES\n    echo \"hello world\" | tr a-z A-Z\n    cat file.txt | tr ' ' '_'",
				"diff": "NAME\n    diff - compare files line by line\n\nSYNOPSIS\n    diff [OPTION]... FILES\n\nOPTIONS\n    -u     output unified context diff\n\nEXAMPLES\n    diff file1.txt file2.txt\n    diff -u old.txt new.txt",
				"ps": "NAME\n    ps - report a snapshot of the current processes\n\nSYNOPSIS\n    ps [OPTIONS]\n\nOPTIONS\n    aux    display all processes on the system with user, cpu, and memory stats\n    -ef    standard full format listing\n\nEXAMPLES\n    ps\n    ps aux\n    ps aux | grep daemontalk",
				"df": "NAME\n    df - report file system disk space usage\n\nSYNOPSIS\n    df [OPTION]... [FILE]...\n\nOPTIONS\n    -h     print sizes in powers of 1024 (e.g., 1023M, 14G)\n\nEXAMPLES\n    df -h",
				"du": "NAME\n    du - estimate file space usage\n\nSYNOPSIS\n    du [OPTION]... [FILE]...\n\nOPTIONS\n    -h     print sizes in human readable format\n    -s     display only a total for each argument\n    -sh    summary human-readable\n\nEXAMPLES\n    du -sh *\n    du -h posts",
				"free": "NAME\n    free - display amount of free and used memory in the system\n\nSYNOPSIS\n    free [OPTION]\n\nOPTIONS\n    -h     show all output fields automatically scaled to shortest three digit unit\n    -m     show output in mebibytes\n\nEXAMPLES\n    free -h\n    free -m",
				"run": "NAME\n    run - compile & execute code via backend runner\n\nSYNOPSIS\n    run <LANGUAGE> <CODE>\n\nSUPPORTED LANGUAGES\n    go, python, js, bash\n\nEXAMPLES\n    run go fmt.Println(\"Hello from Go\")\n    run js console.log(Array.from({length: 5}, (_, i) => i * 2))\n    run python print(sum([x**2 for x in range(10)]))\n    run bash echo \"Kernel: $(uname -r)\""
			};

			// Load real post and project data from backend
			fetch("/api/terminal/data")
				.then(function(r) { return r.json(); })
				.then(function(data) {
					if (!data) return;
					if (data.posts && Array.isArray(data.posts)) {
						rawPosts = data.posts;
						data.posts.forEach(function(p) {
							var fname = p.slug + ".md";
							if (vfs["/home/visitor/posts"].children.indexOf(fname) === -1) {
								vfs["/home/visitor/posts"].children.push(fname);
							}
							var content = "---\ntitle: " + p.title + "\ndate: " + p.date + "\nlang: " + p.lang + "\ntags: [" + (p.tags ? p.tags.join(", ") : "") + "]\nread_time: " + (p.min_read || 5) + " min\n---\n\n" + (p.description || "") + "\n\n" + (p.body_snippet || "");
							vfs["/home/visitor/posts/" + fname] = { type: "file", content: content, slug: p.slug, title: p.title, date: p.date, tags: p.tags || [] };
						});
					}
					if (data.projects && Array.isArray(data.projects)) {
						rawProjects = data.projects;
						data.projects.forEach(function(pr) {
							var fname = pr.slug + ".txt";
							if (vfs["/home/visitor/projects"].children.indexOf(fname) === -1) {
								vfs["/home/visitor/projects"].children.push(fname);
							}
							var content = "Project: " + pr.title + "\nSlug: " + pr.slug + "\nURL: " + (pr.url || "-") + "\nTech: " + (pr.tech ? pr.tech.join(", ") : "-") + "\n\n" + (pr.description || "");
							vfs["/home/visitor/projects/" + fname] = { type: "file", content: content };
						});
					}
				})
				.catch(function(err) {
					console.warn("Failed to pre-populate terminal data:", err);
				});

			// Path helpers
			function getDisplayPath(path) {
				if (path === "/home/visitor") return "~";
				if (path.indexOf("/home/visitor/") === 0) return "~/" + path.substring(14);
				return path;
			}

			function updatePrompt() {
				var disp = getDisplayPath(currentPath);
				envVars["PWD"] = currentPath;
				promptLabel.textContent = "visitor@daemontalk:" + disp + "$";
				titleText.textContent = "visitor@daemontalk: " + disp;
			}

			function resolvePath(target) {
				if (!target || target === "~") return "/home/visitor";
				if (target.indexOf("~/") === 0) target = "/home/visitor/" + target.substring(2);
				
				var segments;
				if (target.startsWith("/")) {
					segments = target.split("/");
				} else {
					segments = (currentPath + "/" + target).split("/");
				}

				var stack = [];
				for (var i = 0; i < segments.length; i++) {
					var s = segments[i];
					if (s === "" || s === ".") continue;
					if (s === "..") {
						if (stack.length > 0) stack.pop();
					} else {
						stack.push(s);
					}
				}
				return "/" + stack.join("/");
			}

			function formatBytes(bytes) {
				if (bytes < 1024) return bytes + "B";
				if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(1) + "K";
				if (bytes < 1024 * 1024 * 1024) return (bytes / (1024 * 1024)).toFixed(1) + "M";
				return (bytes / (1024 * 1024 * 1024)).toFixed(1) + "G";
			}

			// Single Command Evaluator (supports stdin from pipes)
			function evalSingleCommand(cmdLine, stdin) {
				var args = parseArgs(cmdLine);
				var cmd = (args[0] || "").toLowerCase();

				// Alias resolution
				if (aliases[cmd]) {
					var aliasedLine = aliases[cmd] + (args.length > 1 ? " " + args.slice(1).join(" ") : "");
					args = parseArgs(aliasedLine);
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
							"    which <cmd>, alias      Locate binary or manage command aliases\n" +
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
							if (manPages[targetCmd]) {
								out = manPages[targetCmd];
							} else {
								out = "No manual entry for " + targetCmd + "\nType 'help' to see all available commands.";
							}
						}
						break;

					case "pwd":
						out = currentPath;
						break;

					case "cd":
						var target = args[1] || "~";
						var resolved = resolvePath(target);
						if (!vfs[resolved]) {
							out = "cd: no such file or directory: " + target;
						} else if (vfs[resolved].type !== "dir") {
							out = "cd: not a directory: " + target;
						} else {
							currentPath = resolved;
							updatePrompt();
						}
						break;

					case "ls":
						var longFormat = false;
						var showAll = false;
						var targetDir = currentPath;

						for (var i = 1; i < args.length; i++) {
							var a = args[i];
							if (a === "-l") longFormat = true;
							else if (a === "-a") showAll = true;
							else if (a === "-la" || a === "-al") { longFormat = true; showAll = true; }
							else if (a === "-lh" || a === "-lah" || a === "-hal") { longFormat = true; showAll = true; }
							else if (!a.startsWith("-")) { targetDir = resolvePath(a); }
						}

						var node = vfs[targetDir];
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
									var isDir = (vfs[full] && vfs[full].type === "dir") || name === "." || name === "..";
									var perms = isDir ? "drwxr-xr-x" : "-rw-r--r--";
									var size = isDir ? "4096" : String((vfs[full] && vfs[full].content ? vfs[full].content.length : 0));
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
								var fileTarget = resolvePath(args[i]);
								var fnode = vfs[fileTarget];
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
							var fileTarget = resolvePath(filename);
							var fnode = vfs[fileTarget];
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

						var rx = new RegExp(escapeRegExp(pattern), caseInsensitive ? "i" : "");
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
							var resolvedFile = resolvePath(targetFile);
							var fnode = vfs[resolvedFile];
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
							var cnode = vfs[currentPath];
							(cnode.children || []).forEach(function(fname) {
								var full = currentPath === "/" ? "/" + fname : currentPath + "/" + fname;
								var f = vfs[full];
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
						var startDir = currentPath;
						for (var i = 1; i < args.length; i++) {
							if (args[i] === "-L" && args[i+1]) {
								maxDepth = parseInt(args[i+1], 10) || 3;
								i++;
							} else if (!args[i].startsWith("-")) {
								startDir = resolvePath(args[i]);
							}
						}
						if (!vfs[startDir]) {
							out = "tree: '" + startDir + "': No such file or directory";
							break;
						}
						var treeLines = [startDir === "/home/visitor" ? "." : startDir];
						var totalDirs = 0;
						var totalFiles = 0;

						function buildTree(dir, prefix, depth) {
							if (depth > maxDepth) return;
							var dNode = vfs[dir];
							if (!dNode || dNode.type !== "dir") return;
							var kids = (dNode.children || []).slice().sort();
							kids.forEach(function(childName, idx) {
								var isLast = idx === kids.length - 1;
								var pointer = isLast ? "└── " : "├── ";
								var childPath = dir === "/" ? "/" + childName : dir + "/" + childName;
								var cNode = vfs[childPath];
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
							var statPath = resolvePath(args[1]);
							var sNode = vfs[statPath];
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
								var fp = resolvePath(args[i]);
								var fNode = vfs[fp];
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
							var p1 = resolvePath(args[1]);
							var p2 = resolvePath(args[2]);
							var n1 = vfs[p1];
							var n2 = vfs[p2];
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
							var sfNode = vfs[resolvePath(sortFile)];
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
							var ufNode = vfs[resolvePath(uFile)];
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
							var cfn = vfs[resolvePath(cutFile)];
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
								var sn = vfs[resolvePath(sedFile)];
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
							var an = vfs[resolvePath(awkFile)];
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
							var wNode = vfs[resolvePath(targetFile)];
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
						if (vfs[resolvePath(b64Payload)]) {
							b64Payload = vfs[resolvePath(b64Payload)].content || "";
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
							var rn = vfs[resolvePath(args[1])];
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
							var srcPath = resolvePath(args[1]);
							var destPath = resolvePath(args[2]);
							var srcNode = vfs[srcPath];
							if (!srcNode) {
								out = "cp: cannot stat '" + args[1] + "': No such file or directory";
							} else {
								var destParent = destPath.substring(0, destPath.lastIndexOf("/")) || "/";
								var destName = destPath.split("/").pop();
								if (vfs[destPath] && vfs[destPath].type === "dir") {
									destPath = destPath === "/" ? "/" + args[1].split("/").pop() : destPath + "/" + args[1].split("/").pop();
									destParent = destPath.substring(0, destPath.lastIndexOf("/")) || "/";
									destName = destPath.split("/").pop();
								}
								vfs[destPath] = JSON.parse(JSON.stringify(srcNode));
								if (vfs[destParent] && vfs[destParent].children.indexOf(destName) === -1) {
									vfs[destParent].children.push(destName);
								}
							}
						}
						break;

					case "mv":
						if (args.length < 3) {
							out = "mv: missing destination file operand";
						} else {
							var srcPath = resolvePath(args[1]);
							var destPath = resolvePath(args[2]);
							var srcNode = vfs[srcPath];
							if (!srcNode) {
								out = "mv: cannot stat '" + args[1] + "': No such file or directory";
							} else {
								var srcParent = srcPath.substring(0, srcPath.lastIndexOf("/")) || "/";
								var srcName = srcPath.split("/").pop();
								var destParent = destPath.substring(0, destPath.lastIndexOf("/")) || "/";
								var destName = destPath.split("/").pop();

								if (vfs[destPath] && vfs[destPath].type === "dir") {
									destPath = destPath === "/" ? "/" + srcName : destPath + "/" + srcName;
									destParent = destPath.substring(0, destPath.lastIndexOf("/")) || "/";
									destName = destPath.split("/").pop();
								}

								vfs[destPath] = srcNode;
								delete vfs[srcPath];

								if (vfs[srcParent]) {
									var sIdx = vfs[srcParent].children.indexOf(srcName);
									if (sIdx !== -1) vfs[srcParent].children.splice(sIdx, 1);
								}
								if (vfs[destParent] && vfs[destParent].children.indexOf(destName) === -1) {
									vfs[destParent].children.push(destName);
								}
							}
						}
						break;

					case "touch":
						if (args.length < 2) {
							out = "touch: missing file operand";
						} else {
							for (var i = 1; i < args.length; i++) {
								var newPath = resolvePath(args[i]);
								var parentDir = newPath.substring(0, newPath.lastIndexOf("/")) || "/";
								var fname = newPath.split("/").pop();
								if (!vfs[parentDir] || vfs[parentDir].type !== "dir") {
									out = "touch: cannot touch '" + args[i] + "': No such directory";
								} else if (!vfs[newPath]) {
									vfs[newPath] = { type: "file", content: "" };
									if (vfs[parentDir].children.indexOf(fname) === -1) {
										vfs[parentDir].children.push(fname);
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
								var newDir = resolvePath(args[i]);
								var parentDir = newDir.substring(0, newDir.lastIndexOf("/")) || "/";
								var dname = newDir.split("/").pop();
								if (!vfs[parentDir] || vfs[parentDir].type !== "dir") {
									out = "mkdir: cannot create directory '" + args[i] + "': No such parent directory";
								} else if (vfs[newDir]) {
									out = "mkdir: cannot create directory '" + args[i] + "': File exists";
								} else {
									vfs[newDir] = { type: "dir", children: [] };
									if (vfs[parentDir].children.indexOf(dname) === -1) {
										vfs[parentDir].children.push(dname);
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
								var rmPath = resolvePath(target);
								if (!vfs[rmPath]) {
									out = "rm: cannot remove '" + target + "': No such file or directory";
								} else {
									var parentDir = rmPath.substring(0, rmPath.lastIndexOf("/")) || "/";
									var name = rmPath.split("/").pop();
									delete vfs[rmPath];
									if (vfs[parentDir]) {
										var idx = vfs[parentDir].children.indexOf(name);
										if (idx !== -1) vfs[parentDir].children.splice(idx, 1);
									}
								}
							});
						}
						break;

					case "find":
						var startDir = args[1] && !args[1].startsWith("-") ? resolvePath(args[1]) : currentPath;
						var namePattern = null;
						for (var i = 1; i < args.length; i++) {
							if (args[i] === "-name" && args[i+1]) { namePattern = args[i+1]; i++; }
						}
						var results = [];
						function walk(dir) {
							results.push(dir);
							var node = vfs[dir];
							if (node && node.type === "dir") {
								(node.children || []).forEach(function(c) {
									var childPath = dir === "/" ? "/" + c : dir + "/" + c;
									if (vfs[childPath] && vfs[childPath].type === "dir") {
										walk(childPath);
									} else {
										if (!namePattern || childPath.indexOf(namePattern.replace(/\*/g, "")) !== -1) {
											results.push(childPath);
										}
									}
								});
							}
						}
						if (vfs[startDir]) walk(startDir);
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
						if (activeIncidentState.port8080PID !== null) {
							psLines.push("4192 ?        00:04:12 python3 -m http.server 8080");
						}
						if (activeIncidentState.zombieParentPID !== null) {
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
							if (targetPID === 4192 && activeIncidentState.port8080PID === 4192) {
								activeIncidentState.port8080PID = null;
								out = "[4192]+  Terminated              python3 -m http.server 8080";
							} else if (targetPID === 1337 && activeIncidentState.zombieParentPID === 1337) {
								activeIncidentState.zombieParentPID = null;
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
							if (activeIncidentState.inodesFull) {
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
							if (activeIncidentState.inodesFull) {
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
						var duDir = args[1] ? resolvePath(args[1]) : currentPath;
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
							out = envVars[args[1]] || "";
						} else {
							var envList = [];
							for (var k in envVars) envList.push(k + "=" + envVars[k]);
							out = envList.join("\n");
						}
						break;

					case "export":
						if (args.length < 2) {
							for (var k in envVars) out += "declare -x " + k + "=\"" + envVars[k] + "\"\n";
						} else {
							var pair = args.slice(1).join(" ");
							var eqIdx = pair.indexOf("=");
							if (eqIdx !== -1) {
								var varK = pair.substring(0, eqIdx).trim();
								var varV = pair.substring(eqIdx + 1).trim().replace(/^["']|["']$/g, "");
								envVars[varK] = varV;
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
						out = envVars["USER"] || "visitor";
						break;

					case "hostname":
						out = envVars["HOSTNAME"] || "daemontalk.local";
						break;

					case "uname":
						var all = false;
						for (var i = 1; i < args.length; i++) if (args[i] === "-a") all = true;
						out = all ? "Linux daemontalk 6.12.8-generic #1 SMP PREEMPT_DYNAMIC Sat Aug 8 16:00:00 WIB 2026 x86_64 GNU/Linux" : "Linux";
						break;

					case "arch":
						out = envVars["ARCH"] || "x86_64";
						break;

					case "date":
						out = new Date().toString();
						break;

					case "uptime":
						out = " 16:08:12 up 42 days,  3:14,  1 user,  load average: 0.14, 0.08, 0.05";
						break;

					case "dmesg":
						out = vfs["/var/log/dmesg.log"].content || "";
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
						if (activeIncidentState.port8080PID !== null) {
							out = "Netid State  Recv-Q Send-Q Local Address:Port  Peer Address:Port Process\n" +
								"tcp   LISTEN 0      128          0.0.0.0:22         0.0.0.0:*    users:((\"sshd\",pid=104,fd=3))\n" +
								"tcp   LISTEN 0      511          0.0.0.0:80         0.0.0.0:*    users:((\"caddy\",pid=312,fd=7))\n" +
								"tcp   LISTEN 0      511          0.0.0.0:443        0.0.0.0:*    users:((\"caddy\",pid=312,fd=8))\n" +
								"tcp   LISTEN 0      128          0.0.0.0:5432       0.0.0.0:*    users:((\"postgres\",pid=240,fd=5))\n" +
								"tcp   LISTEN 0      128        127.0.0.1:8080       0.0.0.0:*    users:((\"python3\",pid=" + activeIncidentState.port8080PID + ",fd=4))";
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
							incidents.forEach(function(inc) {
								var isSolved = activeIncidentState.solvedList.indexOf(inc.id) !== -1;
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
							var solvedCount = activeIncidentState.solvedList.length;
							out = "Incident Progress: " + solvedCount + " / " + incidents.length + " Solved\n" +
								"Active Incident: " + (activeIncidentState.currentId ? "Incident " + activeIncidentState.currentId : "None (type `incident 1` to start)");
						} else if (sub === "hint") {
							if (!activeIncidentState.currentId) {
								out = "No incident active. Type `incident 1` to start a challenge.";
							} else {
								var curInc = incidents.find(function(i) { return i.id === activeIncidentState.currentId; });
								out = "💡 HINT (Incident " + curInc.id + "):\n" + curInc.hint;
							}
						} else if (sub === "verify" || sub === "solve" || sub === "check") {
							if (!activeIncidentState.currentId) {
								out = "No active incident. Type `incident 1` to start.";
							} else {
								var activeInc = incidents.find(function(i) { return i.id === activeIncidentState.currentId; });
								var answerParam = args.slice(2).join(" ");
								var checkRes = activeInc.verify(answerParam);
								if (checkRes.solved) {
									if (activeIncidentState.solvedList.indexOf(activeInc.id) === -1) {
										activeIncidentState.solvedList.push(activeInc.id);
										localStorage.setItem("daemontalk_solved_incidents", JSON.stringify(activeIncidentState.solvedList));
									}
									out = "╔══════════════════════════════════════════════════════════════════════════════╗\n" +
										"║  🎉 CHALLENGE SOLVED: " + activeInc.title.toUpperCase() + "\n" +
										"╚══════════════════════════════════════════════════════════════════════════════╝\n\n" +
										checkRes.message + "\n\n" +
										"Progress: " + activeIncidentState.solvedList.length + "/" + incidents.length + " completed.\n" +
										"Type `incident list` to select the next incident!";
								} else {
									out = "❌ NOT RESOLVED YET\n" + checkRes.message + "\n\nType `incident hint` if you need guidance.";
								}
							}
						} else {
							var incId = parseInt(sub, 10);
							var targetInc = incidents.find(function(i) { return i.id === incId; });
							if (targetInc) {
								activeIncidentState.currentId = targetInc.id;
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
							if (vfs["/bin"].children.indexOf(cmdName) !== -1) {
								out = "/bin/" + cmdName;
							} else {
								out = cmdName + " not found in PATH";
							}
						}
						break;

					case "alias":
						if (args.length < 2) {
							var aList = [];
							for (var k in aliases) aList.push("alias " + k + "='" + aliases[k] + "'");
							out = aList.join("\n") || "alias: no aliases defined";
						} else {
							var aStr = args.slice(1).join(" ");
							var aEq = aStr.indexOf("=");
							if (aEq !== -1) {
								var aK = aStr.substring(0, aEq).trim();
								var aV = aStr.substring(aEq + 1).trim().replace(/^['"]|['"]$/g, "");
								aliases[aK] = aV;
							}
						}
						break;

					case "unalias":
						if (args[1] && aliases[args[1]]) {
							delete aliases[args[1]];
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
						out = cmdHistory.map(function(c, i) { return "  " + (i + 1) + "  " + c; }).join("\n");
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
			function execute(rawCmd) {
				var trimmed = rawCmd.trim();
				if (!trimmed) {
					appendHistory(rawCmd, "");
					return;
				}

				cmdHistory.push(rawCmd);
				histIndex = cmdHistory.length;

				// Handle async special commands like curl / run
				var firstWord = trimmed.split(/\s+/)[0].toLowerCase();
				if (firstWord === "curl") {
					var cArgs = parseArgs(trimmed);
					if (cArgs.length < 2) {
						appendHistory(rawCmd, "curl: try 'curl <url>' (e.g. curl /api/terminal/data)");
						return;
					}
					var url = cArgs[1];
					appendHistory(rawCmd, "Connecting to " + url + "...");
					fetch(url)
						.then(function(r) { return r.text(); })
						.then(function(text) {
							appendHistory("", text.substring(0, 1500) + (text.length > 1500 ? "\n... [truncated]" : ""));
						})
						.catch(function(err) {
							appendHistory("", "curl: (7) Failed to connect: " + err.toString());
						});
					return;
				}

				if (firstWord === "run") {
					var rArgs = parseArgs(trimmed);
					if (rArgs.length < 3) {
						appendHistory(rawCmd, "Usage: run <go|js|python|bash> <code>\nExample: run go fmt.Println(\"Hello from daemontalk\")");
						return;
					}
					var lang = rArgs[1];
					var code = rArgs.slice(2).join(" ");
					appendHistory(rawCmd, "Executing code (" + lang + ")...");
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
						appendHistory("", result);
					})
					.catch(function(err) {
						appendHistory("", "Execution error: " + err.toString());
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
					var segOut = evalSingleCommand(segment, currentPipeOut);
					if (segOut === null) return; // handled by clear/exit
					currentPipeOut = segOut;
				}

				var finalOut = currentPipeOut || "";

				// Handle File Redirection
				if (redirectFile && finalOut !== null) {
					var rfPath = resolvePath(redirectFile);
					var rfParent = rfPath.substring(0, rfPath.lastIndexOf("/")) || "/";
					var rfName = rfPath.split("/").pop();
					if (vfs[rfParent] && vfs[rfParent].type === "dir") {
						if (!vfs[rfPath]) {
							vfs[rfPath] = { type: "file", content: finalOut + "\n" };
							if (vfs[rfParent].children.indexOf(rfName) === -1) vfs[rfParent].children.push(rfName);
						} else if (redirectAppend) {
							vfs[rfPath].content = (vfs[rfPath].content || "") + finalOut + "\n";
						} else {
							vfs[rfPath].content = finalOut + "\n";
						}
						finalOut = "";
					} else {
						finalOut = "bash: " + redirectFile + ": No such directory";
					}
				}

				appendHistory(rawCmd, finalOut);
			}

			function appendHistory(cmd, output) {
				var item = document.createElement("div");
				item.className = "space-y-1";

				if (cmd !== "") {
					var pLine = document.createElement("div");
					pLine.className = "flex items-baseline gap-2";
					var disp = getDisplayPath(currentPath);
					pLine.innerHTML = '<span class="text-[var(--term-prompt)] font-semibold select-none">visitor@daemontalk:' + disp + '$</span> ' +
						'<span class="text-[var(--term-text)]">' + escapeHTML(cmd) + '</span>';
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

				historyEl.appendChild(item);
				screenEl.scrollTop = screenEl.scrollHeight;
			}

			function parseArgs(str) {
				var match = str.match(/(?:[^\s"']+|"[^"]*"|'[^']*')+/g);
				if (!match) return [];
				return match.map(function(m) {
					if ((m.startsWith('"') && m.endsWith('"')) || (m.startsWith("'") && m.endsWith("'"))) {
						return m.substring(1, m.length - 1);
					}
					return m;
				});
			}

			function escapeRegExp(string) {
				return string.replace(/[.*+?^${}()|[\]\\]/g, '\\$&');
			}

			function escapeHTML(str) {
				return (str || "").replace(/&/g, "&amp;").replace(/</g, "&lt;").replace(/>/g, "&gt;").replace(/"/g, "&quot;");
			}

			// Tab Auto-Completion Engine
			function handleTabComplete() {
				var val = inputEl.value;
				var cursor = inputEl.selectionStart;
				var beforeCursor = val.substring(0, cursor);
				var tokens = beforeCursor.split(/\s+/);
				var lastToken = tokens[tokens.length - 1] || "";

				// If completing command name
				if (tokens.length <= 1) {
					var cmds = vfs["/bin"].children || [];
					var matches = cmds.filter(function(c) { return c.startsWith(lastToken); });
					if (matches.length === 1) {
						inputEl.value = matches[0] + " ";
					} else if (matches.length > 1) {
						appendHistory(val, matches.join("  "));
					}
					return;
				}

				// Completing files or directories
				var dirPath = currentPath;
				var prefix = lastToken;
				if (lastToken.indexOf("/") !== -1) {
					var lastSlash = lastToken.lastIndexOf("/");
					var dirPart = lastToken.substring(0, lastSlash + 1);
					prefix = lastToken.substring(lastSlash + 1);
					dirPath = resolvePath(dirPart);
				}

				var node = vfs[dirPath];
				if (node && node.type === "dir") {
					var entries = (node.children || []).filter(function(c) { return c.startsWith(prefix); });
					if (entries.length === 1) {
						var match = entries[0];
						var full = dirPath === "/" ? "/" + match : dirPath + "/" + match;
						var isDir = vfs[full] && vfs[full].type === "dir";
						var replaceStr = (lastToken.indexOf("/") !== -1 ? lastToken.substring(0, lastToken.lastIndexOf("/") + 1) : "") + match + (isDir ? "/" : " ");
						tokens[tokens.length - 1] = replaceStr;
						inputEl.value = tokens.join(" ");
					} else if (entries.length > 1) {
						appendHistory(val, entries.join("  "));
					}
				}
			}

			// Key Event Listeners
			inputEl.addEventListener("keydown", function(e) {
				if (e.key === "Enter") {
					e.preventDefault();
					var val = inputEl.value;
					inputEl.value = "";
					execute(val);
				} else if (e.key === "Tab") {
					e.preventDefault();
					handleTabComplete();
				} else if (e.key === "ArrowUp") {
					e.preventDefault();
					if (cmdHistory.length === 0) return;
					if (histIndex === cmdHistory.length) {
						currentInputBuffer = inputEl.value;
					}
					if (histIndex > 0) {
						histIndex--;
						inputEl.value = cmdHistory[histIndex];
					}
				} else if (e.key === "ArrowDown") {
					e.preventDefault();
					if (histIndex < cmdHistory.length - 1) {
						histIndex++;
						inputEl.value = cmdHistory[histIndex];
					} else if (histIndex === cmdHistory.length - 1) {
						histIndex = cmdHistory.length;
						inputEl.value = currentInputBuffer;
					}
				} else if (e.ctrlKey && e.key === "l") {
					e.preventDefault();
					window.termClear();
				} else if (e.ctrlKey && e.key === "c") {
					e.preventDefault();
					var val = inputEl.value;
					inputEl.value = "";
					appendHistory(val + "^C", "");
				} else if (e.ctrlKey && e.key === "u") {
					e.preventDefault();
					inputEl.value = "";
				}
			});

			window.termClear = function() {
				historyEl.innerHTML = "";
				inputEl.value = "";
				inputEl.focus();
			};

			window.termExec = function(cmd) {
				inputEl.value = cmd;
				inputEl.focus();
				execute(cmd);
			};

			window.termToggleFullscreen = function() {
				var win = document.getElementById("term-window");
				if (!document.fullscreenElement) {
					win.requestFullscreen().catch(function(e) { console.warn(e); });
				} else {
					document.exitFullscreen();
				}
			};

			function renderMotd() {
				var motd = document.createElement("div");
				motd.className = "mb-4 text-xs text-[var(--term-muted)] space-y-1 font-mono pb-3 border-b border-[var(--term-border)]/50";
				motd.innerHTML = "<div class='text-[var(--term-text)] font-semibold'>daemontalk Linux 6.12.8-generic (x86_64) &bull; bash 5.2.26</div>" +
					"<div>Press <kbd class='px-1.5 py-0.5 rounded border border-[var(--term-border)] bg-[var(--term-chip-bg)] text-[var(--term-text)] text-[11px]'>Tab</kbd> to auto-complete &bull; Type <span class='text-[var(--term-link)] font-semibold cursor-pointer underline hover:text-[var(--term-text)]' onclick=\"window.termExec('help')\">help</span> &bull; Try <span class='text-[var(--term-prompt)] font-bold cursor-pointer underline hover:text-[var(--term-text)]' onclick=\"window.termExec('incident list')\">incident list</span></div>";
				historyEl.appendChild(motd);
			}

			// Initial Setup
			renderMotd();
			updatePrompt();
			setTimeout(function() { inputEl.focus(); }, 100);
		}

		if (document.readyState === "loading") {
			document.addEventListener("DOMContentLoaded", initDaemonTalkTerminal);
		} else {
			initDaemonTalkTerminal();
		}
		document.addEventListener("htmx:afterSwap", initDaemonTalkTerminal);
		document.addEventListener("htmx:load", initDaemonTalkTerminal);
	})();
