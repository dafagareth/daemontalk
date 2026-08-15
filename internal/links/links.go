package links

// Link is a curated external resource.
type Link struct {
	Title       string
	URL         string
	Description string
	Category    string
}

// All is the hand-curated list of useful links.
var All = []Link{
	// Go
	{Title: "pkg.go.dev", URL: "https://pkg.go.dev", Description: "Official Go package documentation.", Category: "Go"},
	{Title: "chi router", URL: "https://github.com/go-chi/chi", Description: "Lightweight, idiomatic HTTP router for Go.", Category: "Go"},
	{Title: "goldmark", URL: "https://github.com/yuin/goldmark", Description: "A Markdown parser written in Go. CommonMark compliant, extensible.", Category: "Go"},
	{Title: "templ", URL: "https://templ.guide", Description: "Type-safe HTML templating language for Go.", Category: "Go"},
	{Title: "The Go Blog", URL: "https://go.dev/blog", Description: "Official articles and announcements from the Go team.", Category: "Go"},

	// Tools
	{Title: "Excalidraw", URL: "https://excalidraw.com", Description: "Virtual whiteboard for hand-drawn-looking diagrams.", Category: "Tools"},
	{Title: "Vale", URL: "https://vale.sh", Description: "Prose linter for technical writing — style as code.", Category: "Tools"},
	{Title: "Tailwind CSS", URL: "https://tailwindcss.com", Description: "Utility-first CSS framework. Build UIs without leaving your HTML.", Category: "Tools"},
	{Title: "HTMX", URL: "https://htmx.org", Description: "Access modern browser features directly from HTML attributes.", Category: "Tools"},
	{Title: "Mermaid", URL: "https://mermaid.js.org", Description: "Diagram-as-code tool that renders in Markdown.", Category: "Tools"},

	// Reading
	{Title: "The Grug Brained Developer", URL: "https://grugbrain.dev", Description: "A meditation on complexity and why simple is better.", Category: "Reading"},
	{Title: "Suckless philosophy", URL: "https://suckless.org/philosophy/", Description: "The case for simple, minimal software.", Category: "Reading"},
	{Title: "Joel on Software", URL: "https://www.joelonsoftware.com", Description: "Classic essays on software development by Joel Spolsky.", Category: "Reading"},
	{Title: "The Architecture of Open Source Applications", URL: "https://aosabook.org/en/", Description: "How real-world open source systems are designed.", Category: "Reading"},

	// Design
	{Title: "Refactoring UI", URL: "https://www.refactoringui.com", Description: "Design tips specifically written for developers.", Category: "Design"},
	{Title: "Heroicons", URL: "https://heroicons.com", Description: "Beautiful hand-crafted SVG icons from the Tailwind team.", Category: "Design"},
	{Title: "Radix Colors", URL: "https://www.radix-ui.com/colors", Description: "An open-source color system designed for building UIs.", Category: "Design"},
}

// CategoryOrder is the preferred display order of link categories.
var CategoryOrder = []string{"Go", "Tools", "Reading", "Design"}

// ByCategory groups links into a map keyed by category.
func ByCategory(links []Link) map[string][]Link {
	m := make(map[string][]Link)
	for _, l := range links {
		m[l.Category] = append(m[l.Category], l)
	}
	return m
}
