package read_command

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

func bookListEnter(m modal) modal {
	m.translationPath = filepath.Join(m.baseDir, m.translations[m.cursor].Name())
	m.state = viewBooks
	m.cursor = 0

	books, err := os.ReadDir(m.translationPath)
	if err != nil {
		m.err = fmt.Errorf("failed to read translation directory %s: %v", m.translationPath, err)
	}

	for _, bookEntry := range books {
		if bookEntry.IsDir() {
			categoryPath := filepath.Join(m.translationPath, bookEntry.Name(), "_category_.json")

			data, err := os.ReadFile(categoryPath)
			if err != nil {
				// Skip if _category_.json is not found or unreadable
				continue
			}

			var category bookCategory
			if err := json.Unmarshal(data, &category); err != nil {
				// Skip if parsing fails
				continue
			}
			m.books = append(m.books, bookInfo{
				Label:    category.Label,
				DirEntry: bookEntry,
				Position: category.Position,
			})
		}
	}

	sort.Slice(m.books, func(i, j int) bool {
		return m.books[i].Position < m.books[j].Position
	})

	return m
}

func bookListDisplay(m modal) string {
	if m.err != nil {
		return "Error: " + m.err.Error() + "\n"
	}

	if len(m.books) == 0 {
		return "No books found in " + filepath.Base(m.translationPath) + "\n"
	}

	tpl := "Books in " + filepath.Base(m.translationPath) + ":\n\n%s\n\n"
	tpl += subtleStyle.Render("j/k, up/down: select") + dotStyle +
		subtleStyle.Render("enter: choose") + dotStyle +
		subtleStyle.Render("q, esc: quit")

	var checkboxes []string
	c := m.cursor
	for i, b := range m.books {
		checkboxes = append(checkboxes, checkbox(b.Label, c == i))
	}

	choices := strings.Join(checkboxes, "\n")

	return fmt.Sprintf(tpl, choices)
}

func (m modal) updateBooks(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		if m.cursor > 0 {
			m.cursor--
		} else {
			m.cursor = len(m.books) - 1
		}
	case "down", "j":
		if m.cursor < len(m.books)-1 {
			m.cursor++
		} else {
			m.cursor = 0
		}
	case "enter":
		return chapterListEnter(m), nil
	}

	return m, nil
}
