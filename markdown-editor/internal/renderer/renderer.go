package renderer

import (
	"bytes"
	"fmt"
	"html"
	"regexp"
	"strings"

	"gitwiki/internal/toolbox"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
)

type Renderer struct {
	Tools []toolbox.Tool
}

func NewRenderer() *Renderer {
	r := &Renderer{}
	r.Tools = []toolbox.Tool{
		&toolbox.ListTool{},
		&toolbox.TableTool{},
	}
	return r
}

func (r *Renderer) Render(content []byte, currentFile string) string {
	// a regex to find all ```tool_name ... ``` blocks
	re := regexp.MustCompile("(?s)```(list|table)(.+?)```")

	// Split the content by the blocks
	parts := re.Split(string(content), -1)
	matches := re.FindAllStringSubmatch(string(content), -1)

	var finalHtml strings.Builder

	for i, part := range parts {
		// render the part as markdown
		var buf bytes.Buffer
		md := goldmark.New(goldmark.WithExtensions(extension.GFM))
		md.Convert([]byte(part), &buf)
		finalHtml.WriteString(buf.String())

		// render the block
		if i < len(matches) {
			match := matches[i]
			toolName := match[1]
			blockContent := match[2]
			for _, tool := range r.Tools {
				if tool.Name() == toolName {
					finalHtml.WriteString(tool.Render([]byte(blockContent), currentFile))
					break
				}
			}
		}
	}

	return "<div class=\"markdown-body\">" + finalHtml.String() + "</div>"
}

func (r *Renderer) RenderHistory(logs []string) string {
	htmlStr := `<div id="git-history" hx-swap-oob="true" class="bg-gray-50 rounded-lg p-4 border border-gray-200 text-xs font-mono text-gray-600 h-screen overflow-y-auto shadow-inner"><h3 class="font-bold text-gray-400 mb-2 uppercase tracking-wider">Activity</h3><ul class="space-y-2">`
	for _, l := range logs {
		htmlStr += fmt.Sprintf("<li>%s</li>", l)
	}
	htmlStr += `</ul></div>`
	return htmlStr
}
func (r *Renderer) RenderDiff(raw string) string {
	lines := strings.Split(raw, "\n")
	output := `<div class="font-mono text-xs overflow-x-auto whitespace-pre">`
	for _, line := range lines {
		escaped := html.EscapeString(line)
		if strings.HasPrefix(line, "+") {
			output += fmt.Sprintf(`<div class="bg-green-100 text-green-800 w-full px-2">%s</div>`, escaped)
		} else if strings.HasPrefix(line, "-") {
			output += fmt.Sprintf(`<div class="bg-red-100 text-red-800 w-full px-2">%s</div>`, escaped)
		} else if strings.HasPrefix(line, "@@") {
			output += fmt.Sprintf(`<div class="bg-indigo-50 text-indigo-500 w-full px-2 mt-2 border-t border-b border-indigo-100 py-1">%s</div>`, escaped)
		} else if strings.HasPrefix(line, "diff") || strings.HasPrefix(line, "index") {
			output += fmt.Sprintf(`<div class="text-gray-400 w-full px-2 font-bold">%s</div>`, escaped)
		} else {
			output += fmt.Sprintf(`<div class="text-gray-600 w-full px-2">%s</div>`, escaped)
		}
	}
	output += `</div>`
	return output
}
func (r *Renderer) RenderFileTree(files []string, current string) string {
	if len(files) == 0 {
		return `<div class="text-gray-400 text-xs italic p-2">No markdown files found.</div>`
	}
	html := `<ul class="space-y-1 text-sm text-gray-600">`
	for _, f := range files {
		active := ""
		if f == current {
			active = "bg-indigo-50 text-indigo-700 font-bold"
		}
		html += fmt.Sprintf(`<li class="group flex items-center justify-between p-2 hover:bg-gray-100 rounded cursor-pointer transition %s"><a href="/%s" class="flex items-center truncate flex-grow"><svg class="w-4 h-4 mr-2 text-indigo-400" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"></path></svg>%s</a><button hx-post="/file/delete" hx-vals='{"name": "%s"}' hx-target="#file-tree" hx-confirm="Delete %s?" class="opacity-0 group-hover:opacity-100 text-gray-400 hover:text-red-500 font-bold px-2">&minus;</button></li>`, active, f, f, f, f)
	}
	html += `</ul>`
	return html
}
