package renderer

import (
	"gitwiki/internal/git"
	"gitwiki/internal/store"
	"gitwiki/internal/toolbox"
	"gitwiki/internal/toolbox/list"
	"gitwiki/internal/toolbox/table"

	"github.com/gin-gonic/gin"
)

type Renderer struct {
	Tools []toolbox.Tool
}

func NewRenderer(store store.Store, vcs git.VCS) *Renderer {
	r := &Renderer{}
	r.Tools = []toolbox.Tool{
		&list.ListTool{},
		table.NewTableTool(store, vcs),
	}
	return r
}

func (r *Renderer) RegisterToolRoutes(rg *gin.RouterGroup) {
	for _, tool := range r.Tools {
		if tool, ok := tool.(interface {
			RegisterRoutes(rg *gin.RouterGroup)
		}); ok {
			tool.RegisterRoutes(rg)
		}
	}
}
