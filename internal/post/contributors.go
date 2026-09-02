package post

import "strings"

// ContributorStats holds contribution metrics for a user.
type ContributorStats struct {
	Username    string   `json:"username"`
	AvatarURL   string   `json:"avatar_url"`
	GitHubURL   string   `json:"github_url"`
	Articles    []string `json:"articles"`
	TotalCounts int      `json:"total_counts"`
}

// GetAllContributors extracts all unique contributors and authors across all parsed posts.
func GetAllContributors(posts []Post) []ContributorStats {
	statsMap := make(map[string]*ContributorStats)

	for _, p := range posts {
		if p.Draft {
			continue
		}

		// 1. Author GitHub
		if p.AuthorGitHub != "" {
			u := strings.ToLower(strings.TrimSpace(p.AuthorGitHub))
			if _, ok := statsMap[u]; !ok {
				avatar := p.AuthorAvatar
				if avatar == "" {
					avatar = "https://github.com/" + u + ".png?size=80"
				}
				statsMap[u] = &ContributorStats{
					Username:  u,
					AvatarURL: avatar,
					GitHubURL: "https://github.com/" + u,
					Articles:  []string{},
				}
			}
			statsMap[u].Articles = append(statsMap[u].Articles, p.Title)
			statsMap[u].TotalCounts++
		}

		// 2. Contributors List
		for _, contrib := range p.Contributors {
			u := strings.ToLower(strings.TrimSpace(contrib))
			if u == "" {
				continue
			}
			if _, ok := statsMap[u]; !ok {
				statsMap[u] = &ContributorStats{
					Username:  u,
					AvatarURL: "https://github.com/" + u + ".png?size=80",
					GitHubURL: "https://github.com/" + u,
					Articles:  []string{},
				}
			}
			statsMap[u].Articles = append(statsMap[u].Articles, p.Title)
			statsMap[u].TotalCounts++
		}
	}

	var result []ContributorStats
	for _, stat := range statsMap {
		result = append(result, *stat)
	}
	return result
}

// IsContributor returns true if the username matches any author or contributor.
func IsContributor(posts []Post, username string) bool {
	u := strings.ToLower(strings.TrimSpace(username))
	if u == "" {
		return false
	}
	for _, p := range posts {
		if strings.ToLower(strings.TrimSpace(p.AuthorGitHub)) == u {
			return true
		}
		for _, c := range p.Contributors {
			if strings.ToLower(strings.TrimSpace(c)) == u {
				return true
			}
		}
	}
	return false
}
