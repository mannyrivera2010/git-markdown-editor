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

func (controller *Controller) HandleInsert(c *gin.Context) {
	f := controller.GetFile(c)
	k := c.PostForm("kind")
	var t string
	for _, tool := range controller.Renderer.Tools {
		if tool.Name() == k {
			t = tool.GetInitialMarkdown()
			break
		}
	}
	if t != "" {
		controller.Store.AppendText(f, t)
		controller.VCS.Commit("Ins " + k)
	}
	controller.RenderResponse(c, f)
}

func (controller *Controller) HandleRaw(c *gin.Context) {
	f := controller.GetFile(c)
	content, _ := controller.Store.Read(f)
	data := struct {
		File    string
		Content string
	}{
		File:    f,
		Content: string(content),
	}
	controller.Templates.ExecuteTemplate(c.Writer, "raw_editor.html", data)
}

func (controller *Controller) HandleRawSave(c *gin.Context) {
	f := controller.GetFile(c)
	controller.Store.WriteRaw(f, []byte(c.PostForm("content")))
	controller.VCS.Commit("Raw Edit")
	controller.RenderResponse(c, f)
}

func (controller *Controller) HandleAdd(c *gin.Context) {
	f := controller.GetFile(c)
	controller.Store.Add(f, c.PostForm("task"))
	controller.VCS.Commit("Add to " + f)
	controller.RenderResponse(c, f)
}

func (controller *Controller) HandleToggle(c *gin.Context) {
	f := controller.GetFile(c)
	i, _ := strconv.Atoi(c.PostForm("line"))
	controller.Store.Toggle(f, i)
	controller.VCS.Commit("Tog in " + f)
	controller.RenderResponse(c, f)
}

func (controller *Controller) HandleDelete(c *gin.Context) {
	f := controller.GetFile(c)
	i, _ := strconv.Atoi(c.PostForm("line"))
	controller.Store.Delete(f, i)
	controller.VCS.Commit("Del in " + f)
	controller.RenderResponse(c, f)
}

func (controller *Controller) HandleArchive(c *gin.Context) {
	f := controller.GetFile(c)
	controller.Store.Archive(f)
	controller.VCS.Commit("Arc in " + f)
	controller.RenderResponse(c, f)
}

func (controller *Controller) HandleUpload(c *gin.Context) {
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

	f := controller.GetFile(c)
	// The path for the markdown should be relative to the web root, so /static/uploads/filename
	markdownImagePath := filepath.Join("/static", "uploads", filename)
	controller.Store.AppendText(f, fmt.Sprintf("\n![](%s)\n", markdownImagePath))
	controller.VCS.Commit("Add image: " + filename)
	controller.RenderResponse(c, f)
}

func (controller *Controller) RenderResponse(c *gin.Context, f string) {
	content, _ := controller.Store.Read(f)
	if c.GetHeader("HX-Request") == "true" {
		c.Writer.Write([]byte(controller.Renderer.Render(content, f) + controller.Renderer.RenderHistory(func() []string {
			l, _ := controller.VCS.Log()
			return l
		}())))
	} else {
		c.Redirect(http.StatusFound, "/"+f)
	}
}
