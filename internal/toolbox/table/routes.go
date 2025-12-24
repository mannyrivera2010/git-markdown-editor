package table

import (
	"gitwiki/internal/git"
	"gitwiki/internal/store"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Routes struct {
	Store    store.Store
	VCS      git.VCS
	Renderer *Renderer
}

func (rt *Routes) RegisterRoutes(r *gin.RouterGroup) {
	r.POST("/table/edit", rt.handleTableEdit)
	r.POST("/table/row/delete", rt.handleTableRowDelete)
}

func (rt *Routes) handleTableEdit(c *gin.Context) {
	f := c.PostForm("file")
	r1, _ := strconv.Atoi(c.PostForm("row"))
	col, _ := strconv.Atoi(c.PostForm("col"))
	rt.Store.TableEditCell(f, r1, col, c.PostForm("value"))
	rt.VCS.Commit("Edit")
}

func (rt *Routes) handleTableRowDelete(c *gin.Context) {
	f := c.PostForm("file")
	r1, _ := strconv.Atoi(c.PostForm("row"))
	rt.Store.TableRemoveRow(f, r1)
	rt.VCS.Commit("DelRow")
	content, _ := rt.Store.Read(f)
	c.Writer.Write([]byte(rt.Renderer.renderTable(content, f)))
}
