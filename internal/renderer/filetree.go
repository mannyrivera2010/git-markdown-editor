package renderer

import (
	"fmt"
)

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