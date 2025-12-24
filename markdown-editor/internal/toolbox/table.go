package toolbox

import (
	"bytes"
	"fmt"
	"html/template"
	"strings"
)

type TableTool struct{}

func (t *TableTool) Name() string {
	return "table"
}

func (t *TableTool) GetButton(currentFile string, tmpl *template.Template) string {
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
	htmlStr := `<div class="overflow-x-auto"><table class="w-full text-sm text-left text-gray-500 border rounded-lg">`
	htmlStr += `<thead class="text-xs text-gray-700 uppercase bg-gray-50"><tr>`
	headers := strings.Split(strings.Trim(tableRows[0], "|"), "|")
	for _, h := range headers {
		htmlStr += fmt.Sprintf(`<th class="px-6 py-3 border-b">%s</th>`, strings.TrimSpace(h))
	}
	htmlStr += `<th class="px-6 py-3 border-b">Action</th></tr></thead><tbody>`
	for i := 2; i < len(tableRows); i++ {
		htmlStr += `<tr class="bg-white border-b hover:bg-gray-50">`
		cells := strings.Split(strings.Trim(tableRows[i], "|"), "|")
		for cIndex, cell := range cells {
			val := strings.TrimSpace(cell)
			inputHtml := fmt.Sprintf(`<input type="text" value="%s" name="value" class="bg-transparent border-none w-full focus:ring-0 p-0" hx-post="/table/edit" hx-vals="{\"row\": %d, \"col\": %d, \"file\": \"%s\"}" hx-trigger="change" hx-swap="none">`, val, i-2, cIndex, currentFile)
			htmlStr += fmt.Sprintf(`<td class="px-6 py-4">%s</td>`, inputHtml)
		}
		btnHtml := fmt.Sprintf(`<button hx-post="/table/row/delete" hx-vals="{\"row\": %d, \"file\": \"%s\"}" hx-target="#table-container" class="text-red-500 hover:text-red-700">×</button>`, i-2, currentFile)
		htmlStr += fmt.Sprintf(`<td class="px-6 py-4">%s</td></tr>`, btnHtml)
	}
	htmlStr += `</tbody></table></div>`
	return htmlStr
}
