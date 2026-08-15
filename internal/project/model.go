package project

import "html/template"

type Status string

const (
	StatusActive    Status = "active"
	StatusCompleted Status = "completed"
	StatusArchived  Status = "archived"
)

type Project struct {
	Name          string
	Slug          string
	Description   string
	DescriptionID string
	TechStack     []string
	RepoURL       string
	DemoURL       string
	Status        Status
	Tags          []string
	Featured      bool
	Order         int
	Body          template.HTML
}
