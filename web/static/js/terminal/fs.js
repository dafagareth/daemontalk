
(function() {
	var DT = window.DaemonTerminal;
	DT.fs.incidents = [
				{
					id: 1,
					title: "The Mysterious 'No Space Left on Device'",
					difficulty: "Medium",
					symptom: "Production log pipeline fails with 'ENOSPC: No space left on device'. However, `df -h` shows the disk is only 38% full! Find the root cause and fix it.",
					hint: "Check filesystem metadata limits: compare `df -h` (block usage) with `df -i` (inode usage). Inspect /var/spool/clientmqueue for abandoned zero-byte files.",
					setup: function() {
						DT.fs.vfs["/var/spool"] = { type: "dir", children: ["clientmqueue"] };
						DT.fs.vfs["/var/spool/clientmqueue"] = { type: "dir", children: ["msg_10482.tmp", "msg_10483.tmp", "msg_10484.tmp", "msg_10485.tmp", "msg_10486.tmp", "abandoned_tokens.lock"] };
						DT.fs.vfs["/var/spool/clientmqueue/msg_10482.tmp"] = { type: "file", content: "" };
						DT.fs.vfs["/var/spool/clientmqueue/abandoned_tokens.lock"] = { type: "file", content: "LOCKED_INODES=6553600" };
						DT.state.activeIncidentState.inodesFull = true;
					},
					verify: function() {
						if (!DT.state.activeIncidentState.inodesFull || !DT.fs.vfs["/var/spool/clientmqueue"] || DT.fs.vfs["/var/spool/clientmqueue"].children.length === 0) {
							DT.state.activeIncidentState.inodesFull = false;
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
						DT.state.activeIncidentState.port8080PID = 4192;
					},
					verify: function() {
						if (DT.state.activeIncidentState.port8080PID === null) {
							return { solved: true, message: "PORT RELEASED: PID 4192 terminated. Port 8080 is now free for new deployments." };
						}
						return { solved: false, message: "Port 8080 is still locked by PID " + DT.state.activeIncidentState.port8080PID + ". Use `ss -tulpn` or `lsof -i :8080` to find it, then `kill " + DT.state.activeIncidentState.port8080PID + "`." };
					}
				},
				{
					id: 3,
					title: "The Silent Midnight Crash (OOM-Killer)",
					difficulty: "Medium",
					symptom: "Redis cache server abruptly disappeared at 03:14:02 UTC with exit code 137. Application logs show no stack trace. Identify why Linux terminated the process.",
					hint: "Kernel-level kills do not appear in user application logs. Check kernel ring buffer logs with `dmesg` or `dmesg | grep -i oom`.",
					setup: function() {
						DT.fs.vfs["/var/log/dmesg.log"].content += "\n[14829.102390] Out of memory: Killed process 8192 (redis-server) total-vm:4194304kB, anon-rss:2048512kB, file-rss:0kB\n[14829.102401] oom_reaper: reaped process 8192 (redis-server), now anon-rss:0kB\n";
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
						DT.state.activeIncidentState.zombieParentPID = 1337;
					},
					verify: function() {
						if (DT.state.activeIncidentState.zombieParentPID === null) {
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
						DT.fs.vfs["/etc/resolv.conf"].content = "# Misconfigured by faulty provisioning script\nnameserver 192.168.1.999\nnameserver 0.0.0.0\n";
					},
					verify: function() {
						var conf = DT.fs.vfs["/etc/resolv.conf"] ? DT.fs.vfs["/etc/resolv.conf"].content : "";
						if (conf.indexOf("1.1.1.1") !== -1 || conf.indexOf("8.8.8.8") !== -1 || conf.indexOf("127.0.0.53") !== -1) {
							return { solved: true, message: "DNS RESOLUTION RESTORED: /etc/resolv.conf repaired with valid nameservers. Hostname resolution functional." };
						}
						return { solved: false, message: "Resolver configuration in /etc/resolv.conf is still invalid. Inspect with `cat /etc/resolv.conf` and update it with `echo 'nameserver 1.1.1.1' > /etc/resolv.conf`." };
					}
				}
			];

			;
	DT.fs.vfs = {
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

			;
	DT.fs.manPages = {
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

			;

	DT.fs.getDisplayPath = function (path) {
				if (path === "/home/visitor") return "~";
				if (path.indexOf("/home/visitor/") === 0) return "~/" + path.substring(14);
				return path;
			}

			
	DT.fs.resolvePath = function (target) {
				if (!target || target === "~") return "/home/visitor";
				if (target.indexOf("~/") === 0) target = "/home/visitor/" + target.substring(2);
				
				var segments;
				if (target.startsWith("/")) {
					segments = target.split("/");
				} else {
					segments = (DT.state.currentPath + "/" + target).split("/");
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

			
	DT.fs.formatBytes = function (bytes) {
				if (bytes < 1024) return bytes + "B";
				if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(1) + "K";
				if (bytes < 1024 * 1024 * 1024) return (bytes / (1024 * 1024)).toFixed(1) + "M";
				return (bytes / (1024 * 1024 * 1024)).toFixed(1) + "G";
			}

			// Single Command Evaluator (supports stdin from pipes)
			
	
	DT.fs.init = function() {
		// Any setup logic
	};
})();
