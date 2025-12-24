package renderer

import (
	"fmt"
	"html"
	"strings"
)

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
