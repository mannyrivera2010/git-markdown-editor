package store

type Store interface {
	Init() error
	Read(path string) ([]byte, error)
	Add(path, task string) error
	Toggle(path string, index int) error
	Delete(path string, index int) error
	Archive(path string) error
	TableAddColumn(path, header string) error
	TableAddRow(path string, data []string) error
	TableRemoveRow(path string, rowIndex int) error
	TableEditCell(path string, row, col int, value string) error
	GetFileTree(recursive bool) ([]string, error)
	CreateFile(name string) error
	DeleteFile(path string) error
	AppendText(path, text string) error
	WriteRaw(path string, content []byte) error
}

// New returns a new Store implementation.
func New() Store {
	return NewFileStore()
}
