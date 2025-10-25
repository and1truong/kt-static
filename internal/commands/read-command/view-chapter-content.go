package read_command

import (
	"fmt"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
)

func chapterContentEnter(m modal) modal {
	m.chapterPath = filepath.Join(m.bookPath, m.chapters[m.cursor].Name())
	m.state = viewChapterContent
	m.chapterCursor = m.cursor
	m.cursor = 0

	if m.chapterContent == "" {
		contentBytes, err := os.ReadFile(m.chapterPath)
		if err != nil {
			m.err = fmt.Errorf("failed to read chapter file %s: %v", m.chapterPath, err)
		}
		m.chapterContent = string(contentBytes)
	}

	return m
}

func chapterContentDisplay(m modal) string {
	if m.err != nil {
		return "Error: " + m.err.Error() + "\n"
	}

	return m.chapterContent
}

func chapterContentUpdate(m modal, msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.state = viewChapters
		m.chapterContent = ""
		m.cursor = m.chapterCursor
		return m, nil
	}
	return m, nil
}
