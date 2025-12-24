package main

import (
	"fmt"
	"html/template"
	"log"
	"os"

	"gitwiki/internal/auth"
	"gitwiki/internal/controllers"
	"gitwiki/internal/git"
	"gitwiki/internal/renderer"
	"gitwiki/internal/store"

	"github.com/gin-gonic/gin"
)

const (
	port = ":8080"
)

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
	co := &controllers.Controller{Store: store, VCS: vcs, Renderer: renderer, Auth: auth, Templates: templates}

	r := gin.Default()
	r.Static("/static", "../../static")
	r.GET("/login", co.HandleLogin)
	r.POST("/login", co.HandleLogin)
	r.GET("/logout", co.HandleLogout)

	protected := r.Group("/")
	protected.Use(co.Protect())
	{
		protected.GET("/", co.RootHandler)
		protected.GET("/:file", co.RootHandler)
		protected.GET("/tree", co.HandleFileTree)
		protected.POST("/file/create", co.HandleFileCreate)
		protected.POST("/file/delete", co.HandleFileDelete)
		protected.POST("/insert", co.HandleInsert)
		protected.GET("/raw", co.HandleRaw)
		protected.POST("/raw/save", co.HandleRawSave)
		protected.POST("/add", co.HandleAdd)
		protected.POST("/toggle", co.HandleToggle)
		protected.POST("/delete", co.HandleDelete)
		protected.POST("/archive", co.HandleArchive)
		protected.POST("/table/col/add", co.HandleTableAddCol)
		protected.POST("/table/row/add", co.HandleTableAddRow)
		protected.POST("/table/row/delete", co.HandleTableRowDelete)
		protected.POST("/table/edit", co.HandleTableEdit)
		protected.POST("/push", co.HandlePush)
		protected.POST("/pull", co.HandlePull)
		protected.GET("/diff", co.HandleDiff)
		protected.GET("/diff/hide", co.HandleDiffHide)
		protected.POST("/upload", co.HandleUpload)
	}

	fmt.Printf("Server running at http://localhost%s\n", port)
	r.Run(port)
}
