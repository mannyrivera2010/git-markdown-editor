package table

import (
	"gitwiki/internal/git"
	"gitwiki/internal/store"

	"github.com/gin-gonic/gin"
)

type TableTool struct {
	Store    store.Store
	VCS      git.VCS
	Routes   *Routes
	Renderer *Renderer
}

func NewTableTool(store store.Store, vcs git.VCS) *TableTool {
	renderer := &Renderer{Store: store}
	routes := &Routes{Store: store, VCS: vcs, Renderer: renderer}
	return &TableTool{Store: store, VCS: vcs, Routes: routes, Renderer: renderer}
}

func (t *TableTool) Name() string {
	return "table"
}

func (t *TableTool) GetButton(currentFile string) string {
	return t.Renderer.GetButton(currentFile)
}

func (t *TableTool) GetInitialMarkdown() string {
	return t.Renderer.GetInitialMarkdown()
}

func (t *TableTool) Render(content []byte, currentFile string) string {
	return t.Renderer.Render(content, currentFile)
}

func (t *TableTool) RegisterRoutes(r *gin.RouterGroup) {
	t.Routes.RegisterRoutes(r)
}
