package store

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type FileStore struct{ mu sync.Mutex }

func NewFileStore() *FileStore {
	return &FileStore{}
}

func (s *FileStore) Init() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, err := os.Stat("todo.md"); os.IsNotExist(err) {
		return ioutil.WriteFile("todo.md", []byte("# My Wiki\n\nWelcome.\nUse the **Toolbox** on the left to add features.\n"), 0644)
	}
	return nil
}
func (s *FileStore) Read(path string) ([]byte, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if strings.Contains(path, "..") {
		return nil, fmt.Errorf("bad")
	}
	return ioutil.ReadFile(path)
}
func (s *FileStore) WriteRaw(path string, content []byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if strings.Contains(path, "..") {
		return fmt.Errorf("bad")
	}
	return ioutil.WriteFile(path, content, 0644)
}
func (s *FileStore) AppendText(path, text string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.WriteString(text)
	return err
}
func (s *FileStore) Add(path, task string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	lines, _ := s.readLines(path)
	insertionLine, inListBlock, hasListBlock := -1, false, false
	for i, line := range lines {
		t := strings.TrimSpace(line)
		if t == "```list" {
			inListBlock = true
			hasListBlock = true
		} else if inListBlock && t == "```" {
			insertionLine = i
			inListBlock = false
		}
	}
	newTask := fmt.Sprintf("- [ ] %s", task)
	if !hasListBlock {
		lines = append(lines, "", "```list", newTask, "```")
	} else if insertionLine != -1 {
		lines = append(lines[:insertionLine+1], lines[insertionLine:]...)
		lines[insertionLine] = newTask
	} else {
		lines = append(lines, newTask)
	}
	return s.writeLines(path, lines)
}
func (s *FileStore) Toggle(path string, index int) error {
	return s.processLine(path, index, func(l string) string {
		if strings.Contains(l, "[ ]") {
			return strings.Replace(l, "[ ]", "[x]", 1)
		}
		return strings.Replace(l, "[x]", "[ ]", 1)
	})
}
func (s *FileStore) Delete(path string, index int) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	lines, _ := s.readLines(path)
	realIndex := s.findRealIndex(lines, index)
	if realIndex != -1 {
		lines = append(lines[:realIndex], lines[realIndex+1:]...)
		return s.writeLines(path, lines)
	}
	return nil
}
func (s *FileStore) Archive(path string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	lines, _ := s.readLines(path)
	var newLines []string
	inListBlock := false
	for _, line := range lines {
		t := strings.TrimSpace(line)
		if t == "```list" {
			inListBlock = true
		}
		if inListBlock && t == "```" {
			inListBlock = false
		}
		if inListBlock && strings.Contains(line, "- [x]") {
			continue
		}
		newLines = append(newLines, line)
	}
	return s.writeLines(path, newLines)
}
func (s *FileStore) TableAddColumn(path, header string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	lines, _ := s.readLines(path)
	newLines := s.mapTableLines(lines, func(rType string, p []string) []string {
		if rType == "header" {
			return append(p, header)
		}
		if rType == "separator" {
			return append(p, "---")
		}
		if rType == "row" {
			return append(p, "")
		}
		return p
	})
	return s.writeLines(path, newLines)
}
func (s *FileStore) TableAddRow(path string, data []string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	lines, _ := s.readLines(path)
	insertIdx, inTable, colCount := -1, false, 0
	for i, line := range lines {
		if strings.TrimSpace(line) == "```table" {
			inTable = true
		}
		if inTable && strings.HasPrefix(line, "|") && strings.Contains(line, "---") {
			colCount = strings.Count(line, "|") - 1
		}
		if inTable && strings.TrimSpace(line) == "```" {
			insertIdx = i
			break
		}
	}
	if insertIdx != -1 {
		for len(data) < colCount {
			data = append(data, "")
		}
		newRow := "| " + strings.Join(data, " | ") + " |"
		lines = append(lines[:insertIdx+1], lines[insertIdx:]...)
		lines[insertIdx] = newRow
		return s.writeLines(path, lines)
	}
	return nil
}
func (s *FileStore) TableRemoveRow(path string, rowIndex int) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	lines, _ := s.readLines(path)
	start, _, _ := s.locateTable(lines)
	realIndex := start + 2 + rowIndex
	if realIndex < len(lines) {
		lines = append(lines[:realIndex], lines[realIndex+1:]...)
		return s.writeLines(path, lines)
	}
	return nil
}
func (s *FileStore) TableEditCell(path string, row, col int, value string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	lines, _ := s.readLines(path)
	start, _, _ := s.locateTable(lines)
	realIndex := start + 2 + row
	if realIndex < len(lines) {
		cells := strings.Split(strings.Trim(lines[realIndex], "|"), "|")
		if col < len(cells) {
			cells[col] = " " + value + " "
			lines[realIndex] = "|" + strings.Join(cells, "|") + "|"
			return s.writeLines(path, lines)
		}
	}
	return nil
}
func (s *FileStore) GetFileTree(recursive bool) ([]string, error) {
	var files []string
	root := "."
	if recursive {
		filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() && strings.HasSuffix(info.Name(), ".md") && !strings.Contains(path, ".git") {
				files = append(files, path)
			}
			return nil
		})
	} else {
		entries, _ := ioutil.ReadDir(root)
		for _, e := range entries {
			if !e.IsDir() && strings.HasSuffix(e.Name(), ".md") {
				files = append(files, e.Name())
			}
		}
	}
	return files, nil
}
func (s *FileStore) CreateFile(name string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	name = strings.TrimSpace(name)
	if name == "" {
		return fmt.Errorf("empty")
	}
	if !strings.HasSuffix(name, ".md") {
		name += ".md"
	}
	if _, err := os.Stat(name); !os.IsNotExist(err) {
		return fmt.Errorf("exists")
	}
	return ioutil.WriteFile(name, []byte(fmt.Sprintf("# %s\n\nLinked.\n", name)), 0644)
}
func (s *FileStore) DeleteFile(path string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if path == "todo.md" {
		return fmt.Errorf("no")
	}
	return os.Remove(path)
}
func (s *FileStore) readLines(path string) ([]string, error) {
	c, _ := ioutil.ReadFile(path)
	return strings.Split(string(c), "\n"), nil
}
func (s *FileStore) writeLines(path string, l []string) error {
	return ioutil.WriteFile(path, []byte(strings.Join(l, "\n")), 0644)
}
func (s *FileStore) locateTable(lines []string) (start, end int, rows []int) {
	inTable := false
	start = -1
	for i, line := range lines {
		t := strings.TrimSpace(line)
		if t == "```table" {
			inTable = true
			start = i
			continue
		}
		if inTable && t == "```" {
			end = i
			break
		}
		if inTable && start != -1 && i > start+2 {
			rows = append(rows, i)
		}
	}
	return
}
func (s *FileStore) mapTableLines(lines []string, mapper func(string, []string) []string) []string {
	inTable, rowIdx := false, 0
	for i, line := range lines {
		t := strings.TrimSpace(line)
		if t == "```table" {
			inTable = true
			rowIdx = 0
			continue
		}
		if inTable && t == "```" {
			inTable = false
			continue
		}
		if inTable {
			cells := strings.Split(strings.Trim(line, "|"), "|")
			rType := "row"
			if rowIdx == 0 {
				rType = "header"
			}
			if rowIdx == 1 {
				rType = "separator"
			}
			newCells := mapper(rType, cells)
			lines[i] = "|" + strings.Join(newCells, "|") + "|"
			rowIdx++
		}
	}
	return lines
}
func (s *FileStore) findRealIndex(lines []string, targetIdx int) int {
	listCount, inListBlock := 0, false
	for i, line := range lines {
		t := strings.TrimSpace(line)
		if t == "```list" {
			inListBlock = true
			continue
		}
		if inListBlock && t == "```" {
			inListBlock = false
			continue
		}
		if inListBlock && strings.HasPrefix(t, "- [") {
			if listCount == targetIdx {
				return i
			}
			listCount++
		}
	}
	return -1
}
func (s *FileStore) processLine(path string, targetIdx int, modifier func(string) string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	lines, _ := s.readLines(path)
	realIndex := s.findRealIndex(lines, targetIdx)
	if realIndex != -1 {
		lines[realIndex] = modifier(lines[realIndex])
		return s.writeLines(path, lines)
	}
	return nil
}
