package controllers

import (
	"html/template"
	"net/http"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	defaultFile = "todo.md"
)

func (controller *Controller) RootHandler(c *gin.Context) {
	// a regex that matches any path ending in .md
	re := regexp.MustCompile(`\.md$`)
	if re.MatchString(c.Request.URL.Path) || c.Request.URL.Path == "/" {
		controller.HandleIndex(c)
	} else {
		c.Status(http.StatusNotFound)
	}
}

func (controller *Controller) HandleIndex(c *gin.Context) {
	f := controller.GetFile(c)
	content, err := controller.Store.Read(f)
	if err != nil {
		controller.Store.CreateFile(f)
		content, _ = controller.Store.Read(f)
	}
	l, _ := controller.VCS.Log()
	files, _ := controller.Store.GetFileTree(false)

	toolboxHtml := ""
	for _, tool := range controller.Renderer.Tools {
		toolboxHtml += tool.GetButton(f)
	}

	tmpl := controller.Templates.Lookup("index.html")
	renderedContent := controller.Renderer.Render(content, f)

	tmpl.Execute(c.Writer, struct {
		Body, FileTree, Toolbox template.HTML
		GitLogs                 []string
		CurrentFile             string
	}{
		Body:        template.HTML(renderedContent),
		FileTree:    template.HTML(controller.Renderer.RenderFileTree(files, f)),
		Toolbox:     template.HTML(toolboxHtml),
		GitLogs:     l,
		CurrentFile: f,
	})
}

func (controller *Controller) HandleFileTree(c *gin.Context) {
	f := controller.GetFile(c)
	rec := c.Query("recursive") == "true"
	files, _ := controller.Store.GetFileTree(rec)
	c.Writer.Write([]byte(controller.Renderer.RenderFileTree(files, f)))
}

func (controller *Controller) HandleFileCreate(c *gin.Context) {
	controller.Store.CreateFile(c.PostForm("filename"))
	controller.VCS.Commit("New File")
	controller.RenderTreeOnly(c)
}

func (controller *Controller) HandleFileDelete(c *gin.Context) {
	controller.Store.DeleteFile(c.PostForm("name"))
	controller.VCS.Commit("Del File")
	controller.RenderTreeOnly(c)
}

func (controller *Controller) GetFile(c *gin.Context) string {
	f := c.Param("file")
	if f != "" && strings.HasSuffix(f, ".md") {
		return f
	}
	f = c.Query("file")
	if f == "" {
		f = c.PostForm("file")
	}
	if f == "" {
		return defaultFile
	}
	return f
}

func (controller *Controller) RenderTreeOnly(c *gin.Context) {
	files, _ := controller.Store.GetFileTree(false)
	c.Writer.Write([]byte(controller.Renderer.RenderFileTree(files, controller.GetFile(c))))
}
