package controllers

import (
	"gitwiki/internal/auth"
	"gitwiki/internal/git"
	"gitwiki/internal/renderer"
	"gitwiki/internal/store"
	"html/template"
)

type Controller struct {
	Store     store.Store
	VCS       git.VCS
	Renderer  *renderer.Renderer
	Auth      *auth.AuthService
	Templates *template.Template
}
