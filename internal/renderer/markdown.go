package renderer

import (
	"bytes"
	"regexp"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
)

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
