package toolbox

import (
	"bytes"
	"fmt"
	"html/template"
	"net/url"
	"regexp"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
)

type ListTool struct{}

func (t *ListTool) Name() string {
	return "list"
}

func (t *ListTool) GetButton(currentFile string, tmpl *template.Template) string {
	var b bytes.Buffer
	tmpl.ExecuteTemplate(&b, "list_button.html", currentFile)
	return b.String()
}

func (t *ListTool) GetInitialMarkdown() string {
	return "\n\n```list\n- [ ] Item\n```\n"
}

func (t *ListTool) Render(content []byte, currentFile string) string {
	// content is the content of the block, without the ```list and ```
	text := string(content)
	wikiRe := regexp.MustCompile(`\[\[([^\]]+)\]\]`)
	text = wikiRe.ReplaceAllStringFunc(text, func(match string) string {
		name := match[2 : len(match)-2]
		target := url.QueryEscape(name + ".md")
		return fmt.Sprintf("[%s](/%s)", name, target)
	})

	var buf bytes.Buffer
	md := goldmark.New(goldmark.WithExtensions(extension.GFM))
	md.Convert([]byte(text), &buf)
	return buf.String()
}
