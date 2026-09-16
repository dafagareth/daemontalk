package handler

import (
	"net/http"

	"daemontalk/internal/handler/distribution"
)

const (
	seoBaseURL = distribution.SeoBaseURL
)

func (h *Handler) DistributionHandler() *distribution.Handler {
	h.distOnce.Do(func() {
		h.distHandler = &distribution.Handler{
			AllPosts:        h.AllPosts,
			ReloadFilePosts: h.ReloadFilePosts,
			RefreshPosts:    h.RefreshPosts,
		}
	})
	return h.distHandler
}

func (h *Handler) RSS(w http.ResponseWriter, r *http.Request) {
	h.DistributionHandler().RSS(w, r)
}

func (h *Handler) JSONFeed(w http.ResponseWriter, r *http.Request) {
	h.DistributionHandler().JSONFeed(w, r)
}

func (h *Handler) Manifest(w http.ResponseWriter, r *http.Request) {
	h.DistributionHandler().Manifest(w, r)
}

func (h *Handler) Robots(w http.ResponseWriter, r *http.Request) {
	h.DistributionHandler().Robots(w, r)
}

func (h *Handler) Sitemap(w http.ResponseWriter, r *http.Request) {
	h.DistributionHandler().Sitemap(w, r)
}

func (h *Handler) SitemapIndex(w http.ResponseWriter, r *http.Request) {
	h.DistributionHandler().SitemapIndex(w, r)
}

func (h *Handler) SitemapEN(w http.ResponseWriter, r *http.Request) {
	h.DistributionHandler().SitemapEN(w, r)
}

func (h *Handler) SitemapID(w http.ResponseWriter, r *http.Request) {
	h.DistributionHandler().SitemapID(w, r)
}

func (h *Handler) SiteOGImage(w http.ResponseWriter, r *http.Request) {
	h.DistributionHandler().SiteOGImage(w, r)
}

func (h *Handler) OGImage(w http.ResponseWriter, r *http.Request) {
	h.DistributionHandler().OGImage(w, r)
}

func (h *Handler) GitHubWebhook(w http.ResponseWriter, r *http.Request) {
	h.DistributionHandler().GitHubWebhook(w, r)
}
