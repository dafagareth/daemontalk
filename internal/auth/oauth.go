package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

// GitHubOAuth manages GitHub OAuth authentication.
type GitHubOAuth struct {
	config *oauth2.Config
}

type githubUserResponse struct {
	ID        int64  `json:"id"`
	Login     string `json:"login"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	AvatarURL string `json:"avatar_url"`
	HTMLURL   string `json:"html_url"`
	Bio       string `json:"bio"`
}

type githubEmailResponse struct {
	Email    string `json:"email"`
	Primary  bool   `json:"primary"`
	Verified bool   `json:"verified"`
}

// NewGitHubOAuth creates a new GitHub OAuth service.
func NewGitHubOAuth(clientID, clientSecret, redirectURL string) *GitHubOAuth {
	return &GitHubOAuth{
		config: &oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			RedirectURL:  redirectURL,
			Scopes:       []string{"read:user", "user:email"},
			Endpoint:     github.Endpoint,
		},
	}
}

// AuthCodeURL returns the authorization URL to redirect the user to.
func (g *GitHubOAuth) AuthCodeURL(state string) string {
	return g.config.AuthCodeURL(state, oauth2.AccessTypeOnline)
}

// ExchangeToken exchanges an authorization code for an OAuth2 access token and retrieves user info.
func (g *GitHubOAuth) ExchangeToken(ctx context.Context, code string) (*User, error) {
	token, err := g.config.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("exchange code for token: %w", err)
	}

	client := g.config.Client(ctx, token)
	client.Timeout = 10 * time.Second

	// 1. Fetch GitHub user profile
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.github.com/user", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetch github user: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("github api status: %d", resp.StatusCode)
	}

	var ghUser githubUserResponse
	if err := json.NewDecoder(resp.Body).Decode(&ghUser); err != nil {
		return nil, fmt.Errorf("decode github user: %w", err)
	}

	email := ghUser.Email

	// 2. If email is not in public profile, fetch from emails endpoint
	if email == "" {
		emailReq, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.github.com/user/emails", nil)
		if err == nil {
			emailReq.Header.Set("Accept", "application/vnd.github.v3+json")
			if emailResp, err := client.Do(emailReq); err == nil {
				defer emailResp.Body.Close()
				var emails []githubEmailResponse
				if json.NewDecoder(emailResp.Body).Decode(&emails) == nil {
					for _, e := range emails {
						if e.Primary && e.Verified {
							email = e.Email
							break
						}
					}
				}
			}
		}
	}

	displayName := ghUser.Name
	if displayName == "" {
		displayName = ghUser.Login
	}

	return &User{
		Provider:    "github",
		ProviderID:  fmt.Sprintf("%d", ghUser.ID),
		Username:    ghUser.Login,
		DisplayName: displayName,
		Email:       email,
		AvatarURL:   ghUser.AvatarURL,
		GitHubURL:   ghUser.HTMLURL,
		Role:        "member",
	}, nil
}
