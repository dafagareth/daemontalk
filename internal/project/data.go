package project

var All = []Project{
	{
		Name:          "svault",
		Slug:          "svault",
		Description:   "Local encrypted secret vault for developers. AES-256-GCM, Argon2id, session-based unlocking, git-aware namespaces. Single binary, zero dependencies.",
		DescriptionID: "Vault secret terenkripsi lokal untuk developer. AES-256-GCM, Argon2id, sesi berbatas waktu, namespace sadar-git. Single binary, tanpa dependensi.",
		TechStack:     []string{"Go", "AES-256-GCM", "Argon2id", "Cobra"},
		RepoURL:       "https://github.com/dafagareth/svault",
		DemoURL:       "https://github.com/dafagareth/svault/releases",
		Status:        StatusActive,
		Tags:          []string{"go", "cli", "security"},
		Featured:      true,
		Order:         1,
	},
	{
		Name:          "daemontalk",
		Slug:          "daemontalk",
		Description:   "This site. Personal technology notebook and reading platform built with Go, Chi, templ, HTMX, and Tailwind. Fast SSR and zero heavy JS.",
		DescriptionID: "Situs ini. Buku catatan eksplorasi teknologi dan platform membaca dengan Go, Chi, templ, HTMX, dan Tailwind. Render cepat di server tanpa JS berat.",
		TechStack:     []string{"Go", "Chi", "templ", "HTMX", "Tailwind", "SQLite"},
		RepoURL:       "https://github.com/dafagareth/daemontalk",
		Status:        StatusActive,
		Tags:          []string{"go", "web", "ssr"},
		Featured:      true,
		Order:         2,
	},
}
