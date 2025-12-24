package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (controller *Controller) HandlePush(c *gin.Context) {
	controller.VCS.Push()
	c.Redirect(http.StatusFound, "/")
}

func (controller *Controller) HandlePull(c *gin.Context) {
	controller.VCS.Pull()
	c.Redirect(http.StatusFound, "/")
}

func (controller *Controller) HandleDiff(c *gin.Context) {
	raw, _ := controller.VCS.DiffLast()
	diffHtml := controller.Renderer.RenderDiff(raw)
	c.Writer.Write([]byte(diffHtml))
	controller.Templates.ExecuteTemplate(c.Writer, "hide_diff_button.html", nil)
}

func (controller *Controller) HandleDiffHide(c *gin.Context) {
	// The content for #diff-viewer (empty to clear it)
	diffViewerContent := ""
	c.Writer.Write([]byte(diffViewerContent))
	controller.Templates.ExecuteTemplate(c.Writer, "diff_button.html", nil)
}
