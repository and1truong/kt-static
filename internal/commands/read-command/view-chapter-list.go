package read_command

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

func chapterListEnter(m modal) modal {
	m.bookPath = filepath.Join(m.translationPath, m.books[m.cursor].DirEntry.Name())
	m.state = viewChapters
	m.bookCursor = m.cursor
	m.cursor = 0

	if m.chapters == nil {
		chapters, err := os.ReadDir(m.bookPath)
		if err != nil {
			m.err = fmt.Errorf("failed to read book directory %s: %v", m.bookPath, err)
		}

		for _, chapterEntry := range chapters {
			if chapterEntry.Name() != "_category_.json" {
				m.chapters = append(m.chapters, chapterEntry)
			}
		}

		sort.Slice(m.chapters, func(i, j int) bool {
			nameI := m.chapters[i].Name()
			nameJ := m.chapters[j].Name()

			numStrI := strings.TrimSuffix(nameI, filepath.Ext(nameI))
			numStrJ := strings.TrimSuffix(nameJ, filepath.Ext(nameJ))

			numI, errI := strconv.Atoi(numStrI)
			numJ, errJ := strconv.Atoi(numStrJ)

			if errI != nil && errJ != nil {
				return nameI < nameJ // both non-numeric, sort alphabetically
			}
			if errI != nil {
				return true // i is non-numeric, j is numeric, non-numeric comes first
			}
			if errJ != nil {
				return false // j is non-numeric, i is numeric, non-numeric comes first
			}

			return numI < numJ // both are numeric
		})

		m.chapterLabels = make([]string, len(m.chapters))
		for i, chapterEntry := range m.chapters {
			name := chapterEntry.Name()
			numStr := strings.TrimSuffix(name, filepath.Ext(name))
			m.chapterLabels[i] = fmt.Sprintf("Chapter %s", numStr)
		}
		m.chapterTotal = len(m.chapters)
	}

	return m
}

func chapterListDisplay(m modal) string {
	if m.err != nil {
		return "Error: " + m.err.Error() + "\n"
	}

	if len(m.chapters) == 0 {
		return "No chapters found in " + m.books[m.cursor].Label + "\n"
	}

	tpl := "Chapters in " + m.books[m.bookCursor].Label + ":\n\n%s\n\n"
	tpl += subtleStyle.Render("j/k, up/down: select") + dotStyle +
		subtleStyle.Render("enter: choose") + dotStyle +
		subtleStyle.Render("q, esc: quit")

	checkboxes := []string{}
	c := m.cursor
	for i, label := range m.chapterLabels {
		checkboxes = append(checkboxes, checkbox(label, c == i))
	}

	choices := strings.Join(checkboxes, "\n")

	return fmt.Sprintf(tpl, choices)
}

func chapterListUpdate(m modal, msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		if m.cursor > 0 {
			m.cursor--
		} else {
			m.cursor = len(m.chapters) - 1
		}
	case "down", "j":
		if m.cursor < len(m.chapters)-1 {
			m.cursor++
		} else {
			m.cursor = 0
		}
	case "enter":
		return chapterContentEnter(m), nil
	}
	return m, nil
}
