package portal

import (
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type CategoryTagLink struct {
	Name string
	Slug string
}

func getCategoryTagLinks(key string) []CategoryTagLink {
	switch strings.ToLower(key) {
	case "software":
		return []CategoryTagLink{
			{Name: "Software", Slug: "software"},
			{Name: "Development", Slug: "development"},
		}
	case "linux":
		return []CategoryTagLink{
			{Name: "Linux", Slug: "linux"},
			{Name: "OS", Slug: "os"},
		}
	case "ai":
		return []CategoryTagLink{
			{Name: "AI", Slug: "ai"},
			{Name: "Machine Learning", Slug: "machine-learning"},
		}
	case "security":
		return []CategoryTagLink{
			{Name: "Security", Slug: "security"},
			{Name: "Privacy", Slug: "privacy"},
		}
	case "networking":
		return []CategoryTagLink{
			{Name: "Networking", Slug: "networking"},
			{Name: "Protocols", Slug: "protocols"},
		}
	case "gaming":
		return []CategoryTagLink{
			{Name: "Gaming", Slug: "gaming"},
			{Name: "Graphics", Slug: "graphics"},
		}
	case "tools":
		return []CategoryTagLink{
			{Name: "Tools", Slug: "tools"},
			{Name: "Workflow", Slug: "workflow"},
		}
	case "science":
		return []CategoryTagLink{
			{Name: "Science", Slug: "science"},
			{Name: "Research", Slug: "research"},
		}
	case "policy":
		return []CategoryTagLink{
			{Name: "Tech Policy", Slug: "policy"},
			{Name: "Law", Slug: "law"},
		}
	case "backend-architecture", "backend":
		return []CategoryTagLink{
			{Name: "Backend", Slug: "backend"},
			{Name: "Architecture", Slug: "architecture"},
		}
	case "systems":
		return []CategoryTagLink{
			{Name: "Systems", Slug: "systems"},
			{Name: "Low-Level", Slug: "low-level"},
		}
	case "container-internals", "containers":
		return []CategoryTagLink{
			{Name: "Containers", Slug: "containers"},
			{Name: "Internals", Slug: "container-internals"},
		}
	case "terminal":
		return []CategoryTagLink{
			{Name: "Terminal", Slug: "terminal"},
			{Name: "Shell", Slug: "shell"},
		}
	case "database":
		return []CategoryTagLink{
			{Name: "Database", Slug: "database"},
			{Name: "Storage", Slug: "storage"},
		}
	case "ebpf":
		return []CategoryTagLink{
			{Name: "eBPF", Slug: "ebpf"},
			{Name: "Observability", Slug: "observability"},
		}
	case "crypto":
		return []CategoryTagLink{
			{Name: "Privacy", Slug: "privacy"},
			{Name: "Crypto", Slug: "crypto"},
		}
	default:
		return []CategoryTagLink{
			{Name: cases.Title(language.English).String(key), Slug: key},
		}
	}
}
