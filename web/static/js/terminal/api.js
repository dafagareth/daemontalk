
(function() {
	var DT = window.DaemonTerminal;

	DT.api.init = function() {
		fetch("/api/terminal/data")
				.then(function(r) { return r.json(); })
				.then(function(data) {
					if (!data) return;
					if (data.posts && Array.isArray(data.posts)) {
						DT.state.rawPosts = data.posts;
						data.posts.forEach(function(p) {
							var fname = p.slug + ".md";
							if (DT.fs.vfs["/home/visitor/posts"].children.indexOf(fname) === -1) {
								DT.fs.vfs["/home/visitor/posts"].children.push(fname);
							}
							var content = "---\ntitle: " + p.title + "\ndate: " + p.date + "\nlang: " + p.lang + "\ntags: [" + (p.tags ? p.tags.join(", ") : "") + "]\nread_time: " + (p.min_read || 5) + " min\n---\n\n" + (p.description || "") + "\n\n" + (p.body_snippet || "");
							DT.fs.vfs["/home/visitor/posts/" + fname] = { type: "file", content: content, slug: p.slug, title: p.title, date: p.date, tags: p.tags || [] };
						});
					}
					if (data.projects && Array.isArray(data.projects)) {
						DT.state.rawProjects = data.projects;
						data.projects.forEach(function(pr) {
							var fname = pr.slug + ".txt";
							if (DT.fs.vfs["/home/visitor/projects"].children.indexOf(fname) === -1) {
								DT.fs.vfs["/home/visitor/projects"].children.push(fname);
							}
							var content = "Project: " + pr.title + "\nSlug: " + pr.slug + "\nURL: " + (pr.url || "-") + "\nTech: " + (pr.tech ? pr.tech.join(", ") : "-") + "\n\n" + (pr.description || "");
							DT.fs.vfs["/home/visitor/projects/" + fname] = { type: "file", content: content };
						});
					}
				})
				.catch(function(err) {
					console.warn("Failed to pre-populate terminal data:", err);
	};

	if (window.DaemonTerminal && window.DaemonTerminal.init) {
		window.DaemonTerminal.init();
	}
})();
