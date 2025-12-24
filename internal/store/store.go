package store

// This import is needed if NewFileStore is in a separate package within store
// If it's in the same package (i.e., file_store.go is in internal/store), this import is not needed.
// For now, assuming it's in the same package.
// If it were in a subpackage like internal/store/filestore, then it would be "gitwiki/internal/store/filestore"

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
