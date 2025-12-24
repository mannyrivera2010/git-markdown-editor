package controllers

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

func (co *Controller) HandleTableAddCol(c *gin.Context) {
	f := co.GetFile(c)
	co.Store.TableAddColumn(f, c.PostForm("name"))
	co.VCS.Commit("Col")
	co.RenderTableOnly(c, f)
}

func (co *Controller) HandleTableAddRow(c *gin.Context) {
	f := co.GetFile(c)
	co.Store.TableAddRow(f, []string{})
	co.VCS.Commit("Row")
	co.RenderTableOnly(c, f)
}

func (co *Controller) HandleTableRowDelete(c *gin.Context) {
	f := co.GetFile(c)
	r1, _ := strconv.Atoi(c.PostForm("row"))
	co.Store.TableRemoveRow(f, r1)
	co.VCS.Commit("DelRow")
	co.RenderTableOnly(c, f)
}

func (co *Controller) HandleTableEdit(c *gin.Context) {
	f := co.GetFile(c)
	r1, _ := strconv.Atoi(c.PostForm("row"))
	col, _ := strconv.Atoi(c.PostForm("col"))
	co.Store.TableEditCell(f, r1, col, c.PostForm("value"))
	co.VCS.Commit("Edit")
}
func (co *Controller) RenderTableOnly(c *gin.Context, f string) {
	// TODO: Fix this after the refactoring
	// content, _ := s.Store.Read(f)
	// c.Writer.Write([]byte(s.Renderer.RenderTable(content, f)))
}
