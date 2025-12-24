package toolbox

type Tool interface {
	Name() string
	GetButton(currentFile string) string
	GetInitialMarkdown() string
	Render(content []byte, currentFile string) string
}
