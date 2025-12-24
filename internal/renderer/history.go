package renderer

import (
	"fmt"
)

func (r *Renderer) RenderHistory(logs []string) string {
	htmlStr := `<div id="git-history" hx-swap-oob="true" class="bg-gray-50 rounded-lg p-4 border border-gray-200 text-xs font-mono text-gray-600 h-screen overflow-y-auto shadow-inner"><h3 class="font-bold text-gray-400 mb-2 uppercase tracking-wider">Activity</h3><ul class="space-y-2">`
	for _, l := range logs {
		htmlStr += fmt.Sprintf("<li>%s</li>", l)
	}
	htmlStr += `</ul></div>`
	return htmlStr
}