package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (co *Controller) HandlePush(c *gin.Context) {
	co.VCS.Push()
	c.Redirect(http.StatusFound, "/")
}

func (co *Controller) HandlePull(c *gin.Context) {
	co.VCS.Pull()
	c.Redirect(http.StatusFound, "/")
}

func (co *Controller) HandleDiff(c *gin.Context) {
	raw, _ := co.VCS.DiffLast()
	diffHtml := co.Renderer.RenderDiff(raw)
	c.Writer.Write([]byte(diffHtml))
	co.Templates.ExecuteTemplate(c.Writer, "hide_diff_button.html", nil)
}

func (co *Controller) HandleDiffHide(c *gin.Context) {
	// The content for #diff-viewer (empty to clear it)
	diffViewerContent := ""
	c.Writer.Write([]byte(diffViewerContent))
	co.Templates.ExecuteTemplate(c.Writer, "diff_button.html", nil)
}
