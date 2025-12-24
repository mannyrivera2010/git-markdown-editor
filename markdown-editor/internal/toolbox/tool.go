package toolbox

import "html/template"

type Tool interface {
	Name() string
	GetButton(currentFile string, tmpl *template.Template) string
	GetInitialMarkdown() string
	Render(content []byte, currentFile string) string
}
