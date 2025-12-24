package table

import (
	"bytes"
	"fmt"
	"gitwiki/internal/store"
	"html/template"
	"strings"
)

type Renderer struct {
	Store store.Store
}

func (r *Renderer) GetButton(currentFile string) string {
	tmpl, err := template.ParseFiles("internal/toolbox/table/table_button.html")
	if err != nil {
		return ""
	}
	var b bytes.Buffer
	tmpl.ExecuteTemplate(&b, "table_button.html", currentFile)
	return b.String()
}

func (r *Renderer) GetInitialMarkdown() string {
	return "```table\n```\n"
}

func (r *Renderer) Render(content []byte, currentFile string) string {
	// content is the content of the block, without the ```table and ```
	tableRows := strings.Split(string(content), "\n")
	if len(tableRows) < 2 {
		return ""
	}
	htmlStr := `<div class="overflow-x-auto" id="table-container">`
	htmlStr += r.renderTable(content, currentFile)
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

func (r *Renderer) renderTable(content []byte, currentFile string) string {
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
