package controllers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (co *Controller) HandleInsert(c *gin.Context) {
	f := co.GetFile(c)
	k := c.PostForm("kind")
	var t string
	for _, tool := range co.Renderer.Tools {
		if tool.Name() == k {
			t = tool.GetInitialMarkdown()
			break
		}
	}
	if t != "" {
		co.Store.AppendText(f, t)
		co.VCS.Commit("Ins " + k)
	}
	co.RenderResponse(c, f)
}

func (co *Controller) HandleRaw(c *gin.Context) {
	f := co.GetFile(c)
	content, _ := co.Store.Read(f)
	data := struct {
		File    string
		Content string
	}{
		File:    f,
		Content: string(content),
	}
	co.Templates.ExecuteTemplate(c.Writer, "raw_editor.html", data)
}

func (co *Controller) HandleRawSave(c *gin.Context) {
	f := co.GetFile(c)
	co.Store.WriteRaw(f, []byte(c.PostForm("content")))
	co.VCS.Commit("Raw Edit")
	co.RenderResponse(c, f)
}

func (co *Controller) HandleAdd(c *gin.Context) {
	f := co.GetFile(c)
	co.Store.Add(f, c.PostForm("task"))
	co.VCS.Commit("Add to " + f)
	co.RenderResponse(c, f)
}

func (co *Controller) HandleToggle(c *gin.Context) {
	f := co.GetFile(c)
	i, _ := strconv.Atoi(c.PostForm("line"))
	co.Store.Toggle(f, i)
	co.VCS.Commit("Tog in " + f)
	co.RenderResponse(c, f)
}

func (co *Controller) HandleDelete(c *gin.Context) {
	f := co.GetFile(c)
	i, _ := strconv.Atoi(c.PostForm("line"))
	co.Store.Delete(f, i)
	co.VCS.Commit("Del in " + f)
	co.RenderResponse(c, f)
}

func (co *Controller) HandleArchive(c *gin.Context) {
	f := co.GetFile(c)
	co.Store.Archive(f)
	co.VCS.Commit("Arc in " + f)
	co.RenderResponse(c, f)
}

func (co *Controller) HandleUpload(c *gin.Context) {
	file, handler, err := c.Request.FormFile("image")
	if err != nil {
		fmt.Printf("Error Retrieving the File: %v\n", err)
		c.String(http.StatusInternalServerError, "Error retrieving file")
		return
	}
	defer file.Close()
	if handler == nil {
		c.String(http.StatusBadRequest, "Invalid file handler")
		return
	}

	// Sanitize filename to prevent directory traversal or other issues
	filename := filepath.Base(handler.Filename)
	if filename == "" {
		c.String(http.StatusBadRequest, "Invalid filename")
		return
	}

	uploadPath := filepath.Join("../../static/uploads", filename)
	dst, err := os.Create(uploadPath)
	if err != nil {
		fmt.Printf("Error creating file: %v\n", err)
		c.String(http.StatusInternalServerError, "Error saving file")
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		fmt.Printf("Error copying file: %v\n", err)
		c.String(http.StatusInternalServerError, "Error saving file")
		return
	}

	f := co.GetFile(c)
	// The path for the markdown should be relative to the web root, so /static/uploads/filename
	markdownImagePath := filepath.Join("/static", "uploads", filename)
	co.Store.AppendText(f, fmt.Sprintf("\n![](%s)\n", markdownImagePath))
	co.VCS.Commit("Add image: " + filename)
	co.RenderResponse(c, f)
}

func (co *Controller) RenderResponse(c *gin.Context, f string) {
	content, _ := co.Store.Read(f)
	if c.GetHeader("HX-Request") == "true" {
		c.Writer.Write([]byte(co.Renderer.Render(content, f) + co.Renderer.RenderHistory(func() []string {
			l, _ := co.VCS.Log()
			return l
		}()))) 
	} else {
		c.Redirect(http.StatusFound, "/"+f)
	}
}
