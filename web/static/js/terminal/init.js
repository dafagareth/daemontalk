
window.DaemonTerminal = {
	state: {
		cmdHistory: [],
		histIndex: -1,
		currentInputBuffer: "",
		currentPath: "/home/visitor",
		envVars: {
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
		},
		aliases: {},
		activeIncidentState: {
			currentId: null,
			inodesFull: false,
			port8080PID: null,
			zombieParentPID: null,
			solvedList: JSON.parse(localStorage.getItem("daemontalk_solved_incidents") || "[]")
		},
		rawPosts: [],
		rawProjects: [],
		elements: {
			screenEl: null,
			historyEl: null,
			inputEl: null,
			promptLabel: null,
			titleText: null
		}
	},
	fs: {},
	cmd: {},
	ui: {},
	api: {},

	init: function() {
		this.state.elements.screenEl = document.getElementById("term-screen");
		this.state.elements.historyEl = document.getElementById("term-history");
		this.state.elements.inputEl = document.getElementById("term-input");
		this.state.elements.promptLabel = document.getElementById("term-prompt-label");
		this.state.elements.titleText = document.getElementById("term-title-text");

		if (!this.state.elements.screenEl || !this.state.elements.inputEl || !this.state.elements.historyEl) return;
		if (!this.ui || !this.ui.init || !this.cmd || !this.cmd.execute) return;
		if (this.state.elements.inputEl.dataset.initialized === "true") return;
		this.state.elements.inputEl.dataset.initialized = "true";

		if (this.fs.init) this.fs.init();
		if (this.ui.init) this.ui.init();
		if (this.api.init) this.api.init();
	}
};

(function() {
	function initDaemonTalkTerminal() {
		if (window.DaemonTerminal && window.DaemonTerminal.init) {
			window.DaemonTerminal.init();
		}
	}
	if (document.readyState === "loading") {
		document.addEventListener("DOMContentLoaded", initDaemonTalkTerminal);
	} else {
		initDaemonTalkTerminal();
	}
	document.addEventListener("htmx:afterSwap", initDaemonTalkTerminal);
	document.addEventListener("htmx:load", initDaemonTalkTerminal);
})();
