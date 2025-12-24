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

func (co *Controller) RootHandler(c *gin.Context) {
	// a regex that matches any path ending in .md
	re := regexp.MustCompile(`\.md$`)
	if re.MatchString(c.Request.URL.Path) || c.Request.URL.Path == "/" {
		co.HandleIndex(c)
	} else {
		c.Status(http.StatusNotFound)
	}
}
func (co *Controller) HandleIndex(c *gin.Context) {
	f := co.GetFile(c)
	content, err := co.Store.Read(f)
	if err != nil {
		co.Store.CreateFile(f)
		content, _ = co.Store.Read(f)
	}
	l, _ := co.VCS.Log()
	files, _ := co.Store.GetFileTree(false)

	toolboxHtml := ""
	for _, tool := range co.Renderer.Tools {
		toolboxHtml += tool.GetButton(f)
	}

	tmpl := co.Templates.Lookup("index.html")
	renderedContent := co.Renderer.Render(content, f)

	tmpl.Execute(c.Writer, struct {
		Body, FileTree, Toolbox template.HTML
		GitLogs                 []string
		CurrentFile             string
	}{
		Body:        template.HTML(renderedContent),
		FileTree:    template.HTML(co.Renderer.RenderFileTree(files, f)),
		Toolbox:     template.HTML(toolboxHtml),
		GitLogs:     l,
		CurrentFile: f,
	})
}
func (co *Controller) HandleFileTree(c *gin.Context) {
	f := co.GetFile(c)
	rec := c.Query("recursive") == "true"
	files, _ := co.Store.GetFileTree(rec)
	c.Writer.Write([]byte(co.Renderer.RenderFileTree(files, f)))
}

func (co *Controller) HandleFileCreate(c *gin.Context) {
	co.Store.CreateFile(c.PostForm("filename"))
	co.VCS.Commit("New File")
	co.RenderTreeOnly(c)
}

func (co *Controller) HandleFileDelete(c *gin.Context) {
	co.Store.DeleteFile(c.PostForm("name"))
	co.VCS.Commit("Del File")
	co.RenderTreeOnly(c)
}

func (co *Controller) GetFile(c *gin.Context) string {
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
func (co *Controller) RenderTreeOnly(c *gin.Context) {
	files, _ := co.Store.GetFileTree(false)
	c.Writer.Write([]byte(co.Renderer.RenderFileTree(files, co.GetFile(c))))
}
