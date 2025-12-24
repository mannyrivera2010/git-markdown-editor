package main

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"gitwiki/internal/auth"
	"gitwiki/internal/git"
	"gitwiki/internal/renderer"
	"gitwiki/internal/store"

	"github.com/gin-gonic/gin"
)

const (
	defaultFile = "todo.md"
	port        = ":8080"
)

type Server struct {
	Store     store.Store
	VCS       git.VCS
	Renderer  *renderer.Renderer
	Auth      *auth.AuthService
	Msg       string
	templates *template.Template
}

func main() {
	store := &store.FileStore{}
	vcs := &git.GitVCS{File: "../../"}
	renderer := renderer.NewRenderer()
	auth := &auth.AuthService{}
	if err := store.Init(); err == nil {
		vcs.Init()
	}
	auth.Init()
	// Ensure the uploads directory exists
	if err := os.MkdirAll("../../static/uploads", 0755); err != nil {
		log.Fatalf("Failed to create uploads directory: %v", err)
	}
	templates := template.Must(template.ParseGlob("../../templates/*.html"))
	srv := &Server{Store: store, VCS: vcs, Renderer: renderer, Auth: auth, templates: templates}

	r := gin.Default()
	r.Static("/static", "../../static")
	r.GET("/login", srv.handleLogin)
	r.POST("/login", srv.handleLogin)
	r.GET("/logout", srv.handleLogout)

	protected := r.Group("/")
	protected.Use(srv.protect())
	{
		protected.GET("/", srv.rootHandler)
		protected.GET("/:file", srv.rootHandler)
		protected.GET("/tree", srv.handleFileTree)
		protected.POST("/file/create", srv.handleFileCreate)
		protected.POST("/file/delete", srv.handleFileDelete)
		protected.POST("/insert", srv.handleInsert)
		protected.GET("/raw", srv.handleRaw)
		protected.POST("/raw/save", srv.handleRawSave)
		protected.POST("/add", srv.handleAdd)
		protected.POST("/toggle", srv.handleToggle)
		protected.POST("/delete", srv.handleDelete)
		protected.POST("/archive", srv.handleArchive)
		protected.POST("/table/col/add", srv.handleTableAddCol)
		protected.POST("/table/row/add", srv.handleTableAddRow)
		protected.POST("/table/row/delete", srv.handleTableRowDelete)
		protected.POST("/table/edit", srv.handleTableEdit)
		protected.POST("/push", srv.handlePush)
		protected.POST("/pull", srv.handlePull)
		protected.GET("/diff", srv.handleDiff)
		protected.GET("/diff/hide", srv.handleDiffHide)
		protected.POST("/upload", srv.handleUpload)
	}

	fmt.Printf("Server running at http://localhost%s\n", port)
	r.Run(port)
}

func (s *Server) rootHandler(c *gin.Context) {
	// a regex that matches any path ending in .md
	re := regexp.MustCompile(`\.md$`)
	if re.MatchString(c.Request.URL.Path) || c.Request.URL.Path == "/" {
		s.handleIndex(c)
	} else {
		c.Status(http.StatusNotFound)
	}
}
func getFile(c *gin.Context) string {
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
func (s *Server) protect() gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Cookie("session_token")
		if err != nil || cookie != "logged-in" {
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
			return
		}
		c.Next()
	}
}
func (s *Server) handleLogin(c *gin.Context) {
	if c.Request.Method == "GET" {
		msg := ""
		if c.Query("login") == "failed" {
			msg = "Invalid"
		}
		s.templates.ExecuteTemplate(c.Writer, "login.html", msg)
		return
	}
	if c.Request.Method == "POST" && s.Auth.Authenticate(c.PostForm("username"), c.PostForm("password")) {
		c.SetCookie("session_token", "logged-in", 3600*24, "/", "", false, true)
		c.Redirect(http.StatusFound, "/")
	} else {
		c.Redirect(http.StatusFound, "/login?login=failed")
	}
}
func (s *Server) handleLogout(c *gin.Context) {
	c.SetCookie("session_token", "", -1, "/", "", false, true)
	c.Redirect(http.StatusFound, "/login")
}
func (s *Server) renderResponse(c *gin.Context, f string) {
	content, _ := s.Store.Read(f)
	if c.GetHeader("HX-Request") == "true" {
		c.Writer.Write([]byte(s.Renderer.Render(content, f) + s.Renderer.RenderHistory(func() []string {
			l, _ := s.VCS.Log()
			return l
		}())))
	} else {
		c.Redirect(http.StatusFound, "/"+f)
	}
}
func (s *Server) renderTableOnly(c *gin.Context, f string) {
	// TODO: Fix this after the refactoring
	// content, _ := s.Store.Read(f)
	// c.Writer.Write([]byte(s.Renderer.RenderTable(content, f)))
}
func (s *Server) renderTreeOnly(c *gin.Context) {
	files, _ := s.Store.GetFileTree(false)
	c.Writer.Write([]byte(s.Renderer.RenderFileTree(files, getFile(c))))
}
func (s *Server) handleDiff(c *gin.Context) {
	raw, _ := s.VCS.DiffLast()
	diffHtml := s.Renderer.RenderDiff(raw)
	c.Writer.Write([]byte(diffHtml))
	s.templates.ExecuteTemplate(c.Writer, "hide_diff_button.html", nil)
}

func (s *Server) handleDiffHide(c *gin.Context) {
	// The content for #diff-viewer (empty to clear it)
	diffViewerContent := ""
	c.Writer.Write([]byte(diffViewerContent))
	s.templates.ExecuteTemplate(c.Writer, "diff_button.html", nil)
}

func (s *Server) handleIndex(c *gin.Context) {
	f := getFile(c)
	content, err := s.Store.Read(f)
	if err != nil {
		s.Store.CreateFile(f)
		content, _ = s.Store.Read(f)
	}
	l, _ := s.VCS.Log()
	files, _ := s.Store.GetFileTree(false)

	toolboxHtml := ""
	for _, tool := range s.Renderer.Tools {
		toolboxHtml += tool.GetButton(f, s.templates)
	}

	tmpl := s.templates.Lookup("index.html")
	renderedContent := s.Renderer.Render(content, f)

	tmpl.Execute(c.Writer, struct {
		Body, FileTree, Toolbox template.HTML
		GitLogs                 []string
		CurrentFile             string
	}{
		Body:        template.HTML(renderedContent),
		FileTree:    template.HTML(s.Renderer.RenderFileTree(files, f)),
		Toolbox:     template.HTML(toolboxHtml),
		GitLogs:     l,
		CurrentFile: f,
	})
}
func (s *Server) handleFileTree(c *gin.Context) {
	f := getFile(c)
	rec := c.Query("recursive") == "true"
	files, _ := s.Store.GetFileTree(rec)
	c.Writer.Write([]byte(s.Renderer.RenderFileTree(files, f)))
}
func (s *Server) handleFileCreate(c *gin.Context) {
	s.Store.CreateFile(c.PostForm("filename"))
	s.VCS.Commit("New File")
	s.renderTreeOnly(c)
}
func (s *Server) handleFileDelete(c *gin.Context) {
	s.Store.DeleteFile(c.PostForm("name"))
	s.VCS.Commit("Del File")
	s.renderTreeOnly(c)
}
func (s *Server) handleInsert(c *gin.Context) {
	f := getFile(c)
	k := c.PostForm("kind")
	var t string
	for _, tool := range s.Renderer.Tools {
		if tool.Name() == k {
			t = tool.GetInitialMarkdown()
			break
		}
	}
	if t != "" {
		s.Store.AppendText(f, t)
		s.VCS.Commit("Ins " + k)
	}
	s.renderResponse(c, f)
}
func (s *Server) handleRaw(c *gin.Context) {
	f := getFile(c)
	content, _ := s.Store.Read(f)
	data := struct {
		File    string
		Content string
	}{
		File:    f,
		Content: string(content),
	}
	s.templates.ExecuteTemplate(c.Writer, "raw_editor.html", data)
}
func (s *Server) handleRawSave(c *gin.Context) {
	f := getFile(c)
	s.Store.WriteRaw(f, []byte(c.PostForm("content")))
	s.VCS.Commit("Raw Edit")
	s.renderResponse(c, f)
}
func (s *Server) handleAdd(c *gin.Context) {
	f := getFile(c)
	s.Store.Add(f, c.PostForm("task"))
	s.VCS.Commit("Add to " + f)
	s.renderResponse(c, f)
}
func (s *Server) handleToggle(c *gin.Context) {
	f := getFile(c)
	i, _ := strconv.Atoi(c.PostForm("line"))
	s.Store.Toggle(f, i)
	s.VCS.Commit("Tog in " + f)
	s.renderResponse(c, f)
}
func (s *Server) handleDelete(c *gin.Context) {
	f := getFile(c)
	i, _ := strconv.Atoi(c.PostForm("line"))
	s.Store.Delete(f, i)
	s.VCS.Commit("Del in " + f)
	s.renderResponse(c, f)
}
func (s *Server) handleArchive(c *gin.Context) {
	f := getFile(c)
	s.Store.Archive(f)
	s.VCS.Commit("Arc in " + f)
	s.renderResponse(c, f)
}
func (s *Server) handleTableAddCol(c *gin.Context) {
	f := getFile(c)
	s.Store.TableAddColumn(f, c.PostForm("name"))
	s.VCS.Commit("Col")
	s.renderTableOnly(c, f)
}
func (s *Server) handleTableAddRow(c *gin.Context) {
	f := getFile(c)
	s.Store.TableAddRow(f, []string{})
	s.VCS.Commit("Row")
	s.renderTableOnly(c, f)
}
func (s *Server) handleTableRowDelete(c *gin.Context) {
	f := getFile(c)
	r1, _ := strconv.Atoi(c.PostForm("row"))
	s.Store.TableRemoveRow(f, r1)
	s.VCS.Commit("DelRow")
	s.renderTableOnly(c, f)
}
func (s *Server) handleTableEdit(c *gin.Context) {
	f := getFile(c)
	r1, _ := strconv.Atoi(c.PostForm("row"))
	col, _ := strconv.Atoi(c.PostForm("col"))
	s.Store.TableEditCell(f, r1, col, c.PostForm("value"))
	s.VCS.Commit("Edit")
}
func (s *Server) handleUpload(c *gin.Context) {
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

	uploadPath := filepath.Join("static", "uploads", filename)
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

	f := getFile(c)
	// The path for the markdown should be relative to the web root, so /static/uploads/filename
	markdownImagePath := filepath.Join("/static", "uploads", filename)
	s.Store.AppendText(f, fmt.Sprintf("\n![](%s)\n", markdownImagePath))
	s.VCS.Commit("Add image: " + filename)
	s.renderResponse(c, f)
}
func (s *Server) handlePush(c *gin.Context) {
	s.VCS.Push()
	c.Redirect(http.StatusFound, "/")
}
func (s *Server) handlePull(c *gin.Context) {
	s.VCS.Pull()
	c.Redirect(http.StatusFound, "/")
}
