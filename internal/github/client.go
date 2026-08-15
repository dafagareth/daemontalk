package github

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

type Repo struct {
	Name        string `json:"name"`
	URL         string `json:"html_url"`
	Description string `json:"description"`
	Stars       int    `json:"stargazers_count"`
	Language    string `json:"language"`
	Fork        bool   `json:"fork"`
}

type ContribDay struct {
	Date  string
	Count int
}

type ContribWeek struct {
	Days []ContribDay
}

type Stats struct {
	Login         string
	Name          string
	AvatarURL     string
	Followers     int
	PublicRepos   int
	TopRepos      []Repo
	Contributions []ContribWeek
	TotalContribs int
}

var cache struct {
	sync.RWMutex
	val Stats
	at  time.Time
}

const cacheTTL = time.Hour

// Fetch returns GitHub stats for the given user, cached for 1 hour.
// If the API is unavailable, returns the last cached value (or empty Stats).
func Fetch(user, token string) Stats {
	cache.RLock()
	if time.Since(cache.at) < cacheTTL && cache.val.Login != "" {
		val := cache.val
		cache.RUnlock()
		return val
	}
	cache.RUnlock()

	stats, err := fetch(user, token)
	if err != nil {
		cache.RLock()
		val := cache.val
		cache.RUnlock()
		return val
	}

	cache.Lock()
	cache.val = stats
	cache.at = time.Now()
	cache.Unlock()
	return stats
}

func fetch(user, token string) (Stats, error) {
	do := func(url string, out interface{}) error {
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return err
		}
		req.Header.Set("Accept", "application/vnd.github+json")
		if token != "" {
			req.Header.Set("Authorization", "Bearer "+token)
		}
		req.Header.Set("User-Agent", "daemontalk.com/1.0")
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("github api: %s", resp.Status)
		}
		return json.NewDecoder(resp.Body).Decode(out)
	}

	var u struct {
		Login       string `json:"login"`
		Name        string `json:"name"`
		AvatarURL   string `json:"avatar_url"`
		Followers   int    `json:"followers"`
		PublicRepos int    `json:"public_repos"`
	}
	if err := do("https://api.github.com/users/"+user, &u); err != nil {
		return Stats{}, err
	}

	var repos []Repo
	if err := do(fmt.Sprintf("https://api.github.com/users/%s/repos?sort=stars&per_page=6&type=owner", user), &repos); err != nil {
		return Stats{}, err
	}

	// Filter forks
	var topRepos []Repo
	for _, r := range repos {
		if !r.Fork {
			topRepos = append(topRepos, r)
		}
		if len(topRepos) == 6 {
			break
		}
	}

	weeks, total := fetchContributions(user, token)

	return Stats{
		Login:         u.Login,
		Name:          u.Name,
		AvatarURL:     u.AvatarURL,
		Followers:     u.Followers,
		PublicRepos:   u.PublicRepos,
		TopRepos:      topRepos,
		Contributions: weeks,
		TotalContribs: total,
	}, nil
}

func fetchContributions(user, token string) ([]ContribWeek, int) {
	if token == "" {
		return nil, 0
	}
	query := `{"query":"{ user(login: \"` + user + `\") { contributionsCollection { contributionCalendar { totalContributions weeks { contributionDays { contributionCount date } } } } } }"}`
	req, err := http.NewRequest("POST", "https://api.github.com/graphql", bytes.NewBufferString(query))
	if err != nil {
		return nil, 0
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "daemontalk.com/1.0")

	resp, err := http.DefaultClient.Do(req)
	if err != nil || resp.StatusCode != 200 {
		return nil, 0
	}
	defer resp.Body.Close()

	var out struct {
		Data struct {
			User struct {
				ContributionsCollection struct {
					ContributionCalendar struct {
						TotalContributions int `json:"totalContributions"`
						Weeks              []struct {
							ContributionDays []struct {
								ContributionCount int    `json:"contributionCount"`
								Date              string `json:"date"`
							} `json:"contributionDays"`
						} `json:"weeks"`
					} `json:"contributionCalendar"`
				} `json:"contributionsCollection"`
			} `json:"user"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, 0
	}

	cal := out.Data.User.ContributionsCollection.ContributionCalendar
	weeks := make([]ContribWeek, len(cal.Weeks))
	for i, w := range cal.Weeks {
		days := make([]ContribDay, len(w.ContributionDays))
		for j, d := range w.ContributionDays {
			days[j] = ContribDay{Date: d.Date, Count: d.ContributionCount}
		}
		weeks[i] = ContribWeek{Days: days}
	}
	return weeks, cal.TotalContributions
}
