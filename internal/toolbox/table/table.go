package table

import (
	"bytes"
	"fmt"
	"gitwiki/internal/git"
	"gitwiki/internal/store"
	"html/template"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type TableTool struct {
	Store store.Store
	VCS   git.VCS
}

func NewTableTool(store store.Store, vcs git.VCS) *TableTool {
	return &TableTool{Store: store, VCS: vcs}
}

func (t *TableTool) Name() string {
	return "table"
}

func (t *TableTool) GetButton(currentFile string) string {
	tmpl, err := template.ParseFiles("internal/toolbox/table/table_button.html")
	if err != nil {
		return ""
	}
	var b bytes.Buffer
	tmpl.ExecuteTemplate(&b, "table_button.html", currentFile)
	return b.String()
}

func (t *TableTool) GetInitialMarkdown() string {
	return "\n\n```table\n| A | B |\n|---|---|\n| 1 | 2 |\n```\n"
}

func (t *TableTool) Render(content []byte, currentFile string) string {
	// content is the content of the block, without the ```table and ```
	tableRows := strings.Split(string(content), "\n")
	if len(tableRows) < 2 {
		return ""
	}
	htmlStr := `<div class="overflow-x-auto" id="table-container">`
	htmlStr += t.renderTable(content, currentFile)
	htmlStr += `</div>`
	return htmlStr
}

type TemplateData struct {
	Headers     []string
	Rows        []Row
	CurrentFile string
}
type Row struct {
	Cells []string
}

func (t *TableTool) renderTable(content []byte, currentFile string) string {
	tableRows := strings.Split(string(content), "\n")
	if len(tableRows) < 2 {
		return ""
	}

	headers := strings.Split(strings.Trim(tableRows[0], "|"), "|")
	for i, h := range headers {
		headers[i] = strings.TrimSpace(h)
	}

	var rows []Row
	for i := 2; i < len(tableRows); i++ {
		var row Row
		cells := strings.Split(strings.Trim(tableRows[i], "|"), "|")
		for _, c := range cells {
			row.Cells = append(row.Cells, strings.TrimSpace(c))
		}
		rows = append(rows, row)
	}

	data := TemplateData{
		Headers:     headers,
		Rows:        rows,
		CurrentFile: currentFile,
	}

	tmpl, err := template.ParseFiles("internal/toolbox/table/table.html")
	if err != nil {
		return fmt.Sprintf("Error parsing template: %v", err)
	}

	var b bytes.Buffer
	if err := tmpl.Execute(&b, data); err != nil {
		return fmt.Sprintf("Error executing template: %v", err)
	}

	return b.String()
}

func (t *TableTool) RegisterRoutes(r *gin.RouterGroup) {
	r.POST("/table/edit", t.handleTableEdit)
	r.POST("/table/row/delete", t.handleTableRowDelete)
}

func (t *TableTool) handleTableEdit(c *gin.Context) {
	f := c.PostForm("file")
	r1, _ := strconv.Atoi(c.PostForm("row"))
	col, _ := strconv.Atoi(c.PostForm("col"))
	t.Store.TableEditCell(f, r1, col, c.PostForm("value"))
	t.VCS.Commit("Edit")
}

func (t *TableTool) handleTableRowDelete(c *gin.Context) {
	f := c.PostForm("file")
	r1, _ := strconv.Atoi(c.PostForm("row"))
	t.Store.TableRemoveRow(f, r1)
	t.VCS.Commit("DelRow")
	content, _ := t.Store.Read(f)
	c.Writer.Write([]byte(t.renderTable(content, f)))
}
